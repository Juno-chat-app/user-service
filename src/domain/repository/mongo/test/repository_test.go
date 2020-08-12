package mongo

import (
	"context"
	"github.com/Juno-chat-app/user-service/config"
	"github.com/Juno-chat-app/user-service/domain/entity"
	"github.com/Juno-chat-app/user-service/domain/repository/mongo"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"os"
	"testing"
	"time"
)

var (
	log  logger.ILogger
	conf *config.Configuration
	repo mongo.IUserRepository
)

func TestMain(m *testing.M) {
	conf, err := config.LoadConfiguration("./repository_test_config.yml")
	if err != nil {
		os.Exit(1)
	}

	log, err = logger.NewLogger()
	if err != nil {
		os.Exit(1)
	}

	repo = mongo.NewUserRepository(conf.CQRSConfig.PersistConfig, log)

	code := m.Run()
	os.Exit(code)
}

func Test_Ping(t *testing.T) {
	ctx := context.Background()
	err := repo.Ping(ctx)
	require.Nil(t, err)
}

func Test_Save_FindWithUserNamePassword_FindWithUserId_Remove(t *testing.T) {
	user := newUser()
	ctx := context.Background()

	user, err := repo.Save(ctx, user)
	require.Nil(t, err)
	require.NotEqual(t, "", user.UserId)

	usr, err := repo.FindWithUserName(ctx, user.UserName)
	require.Nil(t, err)
	require.Equal(t, usr.UserId, user.UserId)

	usr2, err := repo.FindWithUserId(ctx, usr.UserId)
	require.Nil(t, err)
	require.Equal(t, user.UserId, usr2.UserId)

	user, err = repo.Remove(ctx, user)
	require.Nil(t, err)
	require.NotNil(t, user.DeletedAt)
}

func newUser() *entity.User {
	ti := time.Now().UTC()

	user := entity.User{
		UserName: "test",
		Password: "test",
		UserId:   uuid.NewV4().String(),
		Status: &entity.UserStatus{
			Status:         entity.Active,
			ActivationCode: "",
			UpdatedAt:      &ti,
		},
		ContactInfo: &entity.ContactInfo{
			Mobile: "",
			Phone:  "",
			Email:  "test@juno.com",
		},
		Permissions:     nil,
		CreatedAt:       &ti,
		UpdatedAt:       &ti,
		DeletedAt:       nil,
		DocumentVersion: entity.DocumentVersion,
	}

	return &user
}
