package services

import (
	"context"
	"github.com/Juno-chat-app/user-service/domain/entity"
	"github.com/Juno-chat-app/user-service/domain/model/authorization"
	"github.com/Juno-chat-app/user-service/domain/repository/mongo"
	"github.com/Juno-chat-app/user-service/domain/repository/redis"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/twinj/uuid"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type iUserService struct {
	logger     logger.ILogger
	cache      redis.ICache
	repository mongo.IUserRepository
	auth       authorization.IJwtHandler
}

func (i *iUserService) SignUp(ctx context.Context, user *entity.User) (user_ *entity.User, err error) {
	i.logger.Info("SignUp request",
		"method", "SignUp",
		"user-name", user.UserName,
		"user-contact", user.ContactInfo)

	if user.UserName == "" || user.Password == "" || user.ContactInfo.Email == "" {
		return nil, status.Error(http.StatusBadRequest, "invalid value for user-name or password")
	}

	create := time.Now().UTC()

	hashPass, err := generatePasswordOneWayHash(user.Password)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	user.Password = hashPass
	user.UserId = uuid.NewV4().String()
	user.CreatedAt = &create
	user.UpdatedAt = &create

	if user.Status == nil {
		user.Status = &entity.UserStatus{
			Status:         entity.Active,
			ActivationCode: "",
			UpdatedAt:      &create,
		}
	}

	user_, err = i.repository.Save(ctx, user)

	if err != nil {
		i.logger.Error("got error on saving user register request",
			"method", "SignUp",
			"user-id", user.UserId,
			"user-name", user.UserName,
			"user-contact", user.ContactInfo,
			"user-status", user.Status)

		return nil, err
	}

	i.logger.Info("user register completed",
		"method", "SignUp",
		"user-id", user_.UserId,
		"user-name", user_.UserName,
		"user-contact", user_.ContactInfo,
		"user-status", user_.Status)

	user_.Password = "--secret--"
	return user_, nil
}

func (i *iUserService) SignIn(ctx context.Context, user *entity.User) (auth *authorization.TokenDetail, err error) {
	i.logger.Info("SignIn request",
		"method", "SignUp",
		"user-name", user.UserName)

	user_, err := i.repository.FindWithUserName(ctx, user.UserName)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(user.Password, user_.Password) {
		return nil, status.Error(http.StatusUnauthorized, "invalid user-name or password")
	}

	token, err := i.auth.CreateAccessToken(user_.UserId)
	if err != nil {
		return nil, err
	}

	err = i.cache.Set(ctx, token.AccessUUid, token.AccessToken, time.Duration(token.ExpireAt))
	if err != nil {
		return nil, err
	}
	err = i.cache.Set(ctx, token.RefreshUUid, token.RefreshToke, time.Duration(token.RefreshExpireAt))
	if err != nil {
		return nil, err
	}

	i.logger.Info("Signed in successfully",
		"method", "SignUp",
		"user-name", user.UserName)
	return token, nil
}

func (i *iUserService) RefreshToken(ctx context.Context, token *authorization.TokenDetail) (token_ *authorization.TokenDetail, err error) {
	i.logger.Info("Refresh token request",
		"method", "refresh")

	if token.RefreshToke == "" {
		return nil, status.Error(http.StatusBadRequest, "unspecified refresh token")
	}

	token_, err = i.auth.ValidateRefreshToken(token.RefreshToke)
	if err != nil {
		return nil, err
	}

	refreshToken, err := i.cache.Get(ctx, token_.RefreshUUid)
	if err != nil {
		return nil, err
	}

	if refreshToken != token_.RefreshToke {
		return nil, status.Error(http.StatusConflict, "invalid token value")
	}

	if token.UserId != token_.UserId {
		return nil, status.Error(http.StatusConflict, "invalid user-id for token")
	}

	newToken, err := i.auth.CreateAccessToken(token_.UserId)
	if err != nil {
		return nil, err
	}

	err = i.cache.Set(ctx, newToken.AccessUUid, newToken.AccessToken, time.Duration(newToken.ExpireAt))
	if err != nil {
		return nil, err
	}

	token_.AccessToken = newToken.AccessToken
	token_.AccessUUid = newToken.AccessUUid
	token_.ExpireAt = newToken.ExpireAt

	i.logger.Info("token refreshed successfully",
		"method", "RefreshToken",
		"user-id", token_.UserId)
	return token_, nil
}

func (i *iUserService) Validate(ctx context.Context, token *authorization.TokenDetail) (token_ *authorization.TokenDetail, err error) {
	i.logger.Info("token validation",
		"method", "Validate")

	token_, err = i.auth.ValidateAccessToken(token.AccessToken)
	if err != nil {
		return nil, err
	}

	if token.UserId != token_.UserId {
		return nil, status.Error(http.StatusConflict, "invalid access-token")
	}

	if token.AccessToken != token_.AccessToken {
		return nil, status.Error(http.StatusConflict, "invalid access-token")
	}

	accessToken, err := i.cache.Get(ctx, token_.AccessUUid)
	if err != nil {
		return nil, err
	}

	if token_.AccessToken != accessToken {
		return nil, status.Error(http.StatusConflict, "invalid access-token")
	}

	return token_, nil
}

func (i *iUserService) GetUser(ctx context.Context, info entity.ContactInfo) (*entity.User, error) {
	panic("not implemented")
}
