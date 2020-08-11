package mongo

import (
	"context"
	"fmt"
	"github.com/Juno-chat-app/user-service/config"
	"github.com/Juno-chat-app/user-service/domain/entity"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type iUserRepository struct {
	conf       config.PersistConfig
	logger     logger.ILogger
	connection *mongo.Client
}

func (ur *iUserRepository) Save(ctx context.Context, user *entity.User) (userId *entity.User, err error) {
	err = ur.establishConnection(ctx)
	if err != nil {
		return nil, err
	}

	if exist, err := ur.isUserNameOrEmailDuplicated(ctx, user.UserName, user.ContactInfo.Email); err != nil {
		return nil, err
	} else if !exist {
		dbContext, cancel := context.WithTimeout(ctx, time.Duration(ur.conf.ConnectionTimeout)*time.Second)
		defer cancel()

		dbRes, err := ur.connection.Database(ur.conf.UserDatabase, nil).
			Collection(ur.conf.UserCollection, nil).
			InsertOne(dbContext, user)

		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		user.UserId = dbRes.InsertedID.(primitive.ObjectID).Hex()
		_, err = ur.connection.Database(ur.conf.UserDatabase, nil).
			Collection(ur.conf.UserCollection, nil).
			UpdateOne(dbContext, bson.M{"_id": dbRes.InsertedID}, bson.M{"$set": user})

		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		return user, nil
	}

	return nil, nil
}

func (ur *iUserRepository) FindWithUserNamePassword(ctx context.Context, userName string, password string) (user *entity.User, err error) {
	err = ur.establishConnection(ctx)
	if err != nil {
		return nil, err
	}

	dbContext, cancel := context.WithTimeout(ctx, time.Duration(ur.conf.ConnectionTimeout)*time.Second)
	defer cancel()
	query := bson.M{
		"$and": []bson.M{
			bson.M{string(entity.UserNamePath): userName},
			bson.M{string(entity.PasswordPath): password},
			bson.M{string(entity.DeletedAtPath): nil},
		},
	}

	res := ur.connection.Database(ur.conf.UserDatabase, nil).
		Collection(ur.conf.UserCollection, nil).
		FindOne(dbContext, query)

	usr := entity.User{}
	err = res.Decode(&usr)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &usr, nil
}

func (ur *iUserRepository) FindWithUserId(ctx context.Context, userId string) (user *entity.User, err error) {
	err = ur.establishConnection(ctx)
	if err != nil {
		return nil, err
	}

	dbContext, cancel := context.WithTimeout(ctx, time.Duration(ur.conf.ConnectionTimeout)*time.Second)
	defer cancel()
	query := bson.M{
		string(entity.UserIdPath): userId,
	}
	res := ur.connection.Database(ur.conf.UserDatabase, nil).
		Collection(ur.conf.UserCollection, nil).
		FindOne(dbContext, query)

	usr := entity.User{}
	err = res.Decode(&usr)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &usr, nil
}

func (ur *iUserRepository) Remove(ctx context.Context, user *entity.User) (user_ *entity.User, err error) {
	err = ur.establishConnection(ctx)
	if err != nil {
		return nil, err
	}

	dbContext, cancel := context.WithTimeout(ctx, time.Duration(ur.conf.ConnectionTimeout)*time.Second)
	defer cancel()

	tm := time.Now().UTC()

	user.Status.UpdatedAt = &tm
	user.Status.Status = entity.Inactive
	user.UpdatedAt = &tm
	user.DeletedAt = &tm

	id, _ := primitive.ObjectIDFromHex(user.UserId)

	_, err = ur.connection.Database(ur.conf.UserDatabase, nil).
		Collection(ur.conf.UserCollection, nil).
		UpdateOne(dbContext, bson.M{"_id": id}, bson.M{"$set": user}, nil)

	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return user, nil
}

func (ur *iUserRepository) Ping(ctx context.Context) (err error) {
	err = ur.establishConnection(ctx)
	if err != nil {
		return err
	}

	err = ur.connection.Ping(ctx, nil)
	if err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (ur *iUserRepository) isUserNameOrEmailDuplicated(ctx context.Context, userName string, email string) (isDuplicated bool, err error) {
	dbContext, cancel := context.WithTimeout(ctx, time.Duration(ur.conf.ConnectionTimeout)*time.Second)
	defer cancel()
	query := bson.M{
		"$or": []bson.M{
			bson.M{string(entity.UserNamePath): userName},
			bson.M{string(entity.EmailPath): email},
		},
	}

	res := ur.connection.Database(ur.conf.UserDatabase, nil).
		Collection(ur.conf.UserCollection, nil).
		FindOne(dbContext, query)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, status.Error(http.StatusInternalServerError, res.Err().Error())
	}

	data := entity.User{}
	err = res.Decode(&data)
	if err != nil {
		return false, status.Error(http.StatusInternalServerError, err.Error())
	}

	if data.UserName == userName && data.DeletedAt == nil {
		return true, status.Error(http.StatusConflict, "User name exist")
	} else if data.ContactInfo.Email == email && data.DeletedAt == nil {
		return true, status.Error(http.StatusConflict, "Email already registered")
	}

	return false, nil
}

func (ur *iUserRepository) establishConnection(ctx context.Context) (err error) {
	if ur.connection == nil {
		var (
			auth string
			uri  string
		)

		if ur.conf.UserName != "" && ur.conf.Password != "" {
			auth = fmt.Sprintf("%v:%v@", ur.conf.UserName, ur.conf.Password)
		}
		if ur.conf.ConnectionUri != "" {
			uri = fmt.Sprintf("mongodb://%v%v", auth, ur.conf.ConnectionUri)
		} else if ur.conf.Host == "" || ur.conf.Port == 0 {
			return status.Error(http.StatusInternalServerError, "mongo connection is not specified")
		} else {
			uri = fmt.Sprintf("mongodb://%v%v:%d", auth, ur.conf.Host, ur.conf.Port)
		}

		connectionCtx, _ := context.WithTimeout(ctx, time.Duration(ur.conf.ConnectionTimeout)*time.Second)
		clientOptions := options.Client().ApplyURI(uri)

		ur.connection, err = mongo.Connect(connectionCtx, clientOptions)
		if err != nil {
			return status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}
