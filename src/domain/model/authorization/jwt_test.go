package authorization

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	handler IJwtHandler
	token   *TokenDetail
)

func TestMain(m *testing.M) {
	handler = NewJwtHandler(int64(time.Second), int64(time.Minute))

	code := m.Run()
	os.Exit(code)
}

func TestIJwtHandler(t *testing.T) {
	var err error

	token, err = handler.CreateAccessToken("123455")
	require.Nil(t, err)

	token1, err := handler.ValidateAccessToken(token.AccessToken)
	require.Nil(t, err)
	require.Equal(t, token.UserId, token1.UserId)
	require.Equal(t, token.AccessUUid, token1.AccessUUid)

	token2, err := handler.ValidateRefreshToken(token.RefreshToke)
	require.Nil(t, err)
	require.Equal(t, token.UserId, token2.UserId)
	require.Equal(t, token.RefreshUUid, token2.RefreshUUid)
}
