package authorization

import "github.com/Juno-chat-app/user-service/infra/logger"

type TokenDetail struct {
	AccessToken     string
	RefreshToke     string
	AccessUUid      string
	RefreshUUid     string
	UserId          string
	ExpireAt        int64
	RefreshExpireAt int64
}

type IJwtHandler interface {
	CreateAccessToken(userId string) (tokenDetail *TokenDetail, err error)
	ValidateAccessToken(accessToken string) (tokenDetail *TokenDetail, err error)
	ValidateRefreshToken(refreshToken string) (tokenDetail *TokenDetail, err error)
}

func NewJwtHandler(accessTtl int64, refreshTtl int64, logger logger.ILogger) IJwtHandler {
	logger.Info("initial jwt-handler",
		"method", "NewJwtHandler",
		"access-ttl", accessTtl,
		"refresh-ttl", refreshTtl)

	handler := iJwtHandler{
		accessTokenTtl:  accessTtl,
		refreshTokenTtl: refreshTtl,
	}

	return &handler
}
