package authorization

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"
)

type (
	key    string
	secret string
)

const (
	accessSecret  secret = "absdnmqpowiey7219gh214kr3o2p[   :)) 	l]asdkjn21]21n1"
	refreshSecret secret = "askldnqoiwejqopkpglm;werq/weqlkqwna 		alksndasd asdlakdml;aw'qwqeui21y8192y123m1"

	// keys :: -- key are as follow
	authorized  string = "authorized"
	accessUUid  string = "access_uuid"
	refreshUUid string = "refresh_uuid"
	usrId       string = "user_id"
	exp         string = "exp"
	algorithm   string = "alg"
)

type iJwtHandler struct {
	accessTokenTtl  int64
	refreshTokenTtl int64
}

func (j *iJwtHandler) CreateAccessToken(userId string) (tokenDetail *TokenDetail, err error) {
	td := TokenDetail{
		AccessUUid:      uuid.NewV4().String(),
		RefreshUUid:     uuid.NewV4().String(),
		UserId:          userId,
		ExpireAt:        time.Now().Add(time.Duration(j.accessTokenTtl)).Unix(),
		RefreshExpireAt: time.Now().Add(time.Duration(j.refreshTokenTtl)).Unix(),
	}

	//========= Create access token with access claims
	accessClaim := jwt.MapClaims{
		authorized: true,
		accessUUid: td.AccessUUid,
		usrId:      userId,
		exp:        td.ExpireAt,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaim)
	td.AccessToken, err = accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	//========= Create refresh token with refresh claims
	refreshClaim := jwt.MapClaims{
		authorized:  true,
		refreshUUid: td.RefreshUUid,
		usrId:       userId,
		exp:         td.RefreshExpireAt,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	td.RefreshToke, err = refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &td, nil
}

func (j *iJwtHandler) ValidateAccessToken(accessToken string) (tokenDetail *TokenDetail, err error) {
	if strings.HasPrefix(accessToken, "Bearer ") {
		accessToken = strings.Split(accessToken, " ")[1]
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("unexpected signing method: %v", token.Header[algorithm]))
		}
		return []byte(accessSecret), nil
	})

	if err != nil {
		return nil, status.Error(http.StatusBadRequest, err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, status.Error(http.StatusBadRequest, "invalid claim")
	} else {
		accessUuid, ok := claims[accessUUid].(string)
		if !ok {
			return nil, status.Error(http.StatusUnauthorized, " invalid claim")
		}

		userId, ok := claims[usrId].(string)
		if !ok {
			return nil, status.Error(http.StatusUnauthorized, " invalid claim")
		}

		detail := TokenDetail{
			AccessToken: accessToken,
			AccessUUid:  accessUuid,
			UserId:      userId,
		}

		return &detail, nil
	}
}

func (j *iJwtHandler) ValidateRefreshToken(refreshToken string) (tokenDetail *TokenDetail, err error) {
	if strings.HasPrefix(refreshToken, "Bearer ") {
		refreshToken = strings.Split(refreshToken, " ")[1]
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(http.StatusInternalServerError, fmt.Sprintf("unexpected signing method: %v", token.Header[algorithm]))
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, status.Error(http.StatusBadRequest, "invalid claim")
	} else {
		refreshUuid, ok := claims[refreshUUid].(string)
		if !ok {
			return nil, status.Error(http.StatusUnauthorized, " invalid claim")
		}

		userId, ok := claims[usrId].(string)
		if !ok {
			return nil, status.Error(http.StatusUnauthorized, " invalid claim")
		}

		detail := TokenDetail{
			RefreshToke: refreshToken,
			RefreshUUid: refreshUuid,
			UserId:      userId,
		}

		return &detail, nil
	}
}
