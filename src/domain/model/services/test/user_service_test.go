package test

import (
	"context"
	"github.com/Juno-chat-app/user-service/config"
	"github.com/Juno-chat-app/user-service/domain/entity"
	"github.com/Juno-chat-app/user-service/domain/model/authorization"
	"github.com/Juno-chat-app/user-service/domain/model/services"
	"github.com/Juno-chat-app/user-service/domain/repository/mongo"
	"github.com/Juno-chat-app/user-service/domain/repository/redis"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	log     logger.ILogger
	conf    *config.Configuration
	cache   redis.ICache
	auth    authorization.IJwtHandler
	repo    mongo.IUserRepository
	service services.IUserService
)

func TestMain(m *testing.M) {
	conf, err := config.LoadConfiguration("./user_service_test_config.yml")
	if err != nil {
		os.Exit(1)
	}

	log, err = logger.NewLogger()
	if err != nil {
		os.Exit(1)
	}

	repo = mongo.NewUserRepository(conf.CQRSConfig.PersistConfig, log)
	cache = redis.NewCache(conf.CQRSConfig.CacheConfig.Host, conf.CQRSConfig.CacheConfig.Port, conf.CQRSConfig.CacheConfig.Password, conf.CQRSConfig.CacheConfig.Db, conf.CQRSConfig.CacheConfig.Retry, log)
	auth = authorization.NewJwtHandler(int64(conf.AuthConfig.AccessTTL*time.Minute), int64(conf.AuthConfig.RefreshTTL*time.Hour), log)
	service = services.NewUserService(log, repo, cache, auth)

	code := m.Run()
	os.Exit(code)
}

func Test_User_Service(t *testing.T) {
	ctx := context.Background()
	user := NewUser()

	signUpResult, err := service.SignUp(ctx, user)
	user.Password = "test"

	require.Nil(t, err)
	require.Equal(t, signUpResult.Status.Status, entity.Active)

	signInResult, err := service.SignIn(ctx, user)
	require.Nil(t, err)
	require.NotEqual(t, signInResult.AccessToken, "")

	_, err = service.Validate(ctx, signInResult)
	require.Nil(t, err)

	refreshResult, err := service.RefreshToken(ctx, signInResult)
	require.Nil(t, err)
	require.NotEqual(t, signInResult.AccessUUid, refreshResult.AccessUUid)
	require.NotEqual(t, signInResult.AccessToken, refreshResult.AccessToken)

	_, err = repo.Remove(ctx, user)
	require.Nil(t, err)
}

func NewUser() *entity.User {
	user := entity.User{
		UserName: "test",
		Password: "test",
		UserId:   "",
		Status:   nil,
		ContactInfo: &entity.ContactInfo{
			Mobile: "",
			Phone:  "",
			Email:  "test@juno.com",
		},
		Permissions:     nil,
		CreatedAt:       nil,
		UpdatedAt:       nil,
		DeletedAt:       nil,
		DocumentVersion: "",
	}

	return &user
}
