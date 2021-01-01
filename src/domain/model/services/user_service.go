package services

import (
	"context"
	"github.com/Juno-chat-app/user-service/domain/entity"
	"github.com/Juno-chat-app/user-service/domain/model/authorization"
	"github.com/Juno-chat-app/user-service/domain/repository/mongo"
	"github.com/Juno-chat-app/user-service/domain/repository/redis"
	"github.com/Juno-chat-app/user-service/infra/logger"
)

type IUserService interface {
	SignUp(ctx context.Context, user *entity.User) (user_ *entity.User, err error)
	SignIn(ctx context.Context, user *entity.User) (auth *authorization.TokenDetail, err error)
	RefreshToken(ctx context.Context, token *authorization.TokenDetail) (token_ *authorization.TokenDetail, err error)
	Validate(ctx context.Context, token *authorization.TokenDetail) (token_ *authorization.TokenDetail, err error)
	GetUser(ctx context.Context, info entity.ContactInfo) (usr *entity.User, err error)
}

func NewUserService(logger logger.ILogger, repo mongo.IUserRepository, cache redis.ICache, auth authorization.IJwtHandler) IUserService {
	userService := iUserService{
		logger:     logger,
		cache:      cache,
		repository: repo,
		auth:       auth,
	}

	return &userService
}
