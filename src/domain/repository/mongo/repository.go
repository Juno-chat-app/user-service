package mongo

import (
	"context"
	"github.com/Juno-chat-app/user-service/config"
	"github.com/Juno-chat-app/user-service/domain/entity"
	"github.com/Juno-chat-app/user-service/infra/logger"
)

type IUserRepository interface {
	Save(ctx context.Context, user *entity.User) (user_ *entity.User, err error)
	FindWithUserName(ctx context.Context, userName string) (user *entity.User, err error)
	FindWithUserId(ctx context.Context, userId string) (user *entity.User, err error)
	Remove(ctx context.Context, user *entity.User) (user_ *entity.User, err error)
	Ping(ctx context.Context) (err error)
}

func NewUserRepository(conf config.PersistConfig, logger logger.ILogger) IUserRepository {
	logger.Info("create mongo client",
		"method", "NewUserRepository",
		"host", conf.Host,
		"port", conf.Port,
		"uri", conf.ConnectionUri)

	repo := iUserRepository{
		conf:       conf,
		logger:     logger,
		connection: nil,
	}

	return &repo
}
