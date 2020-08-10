package redis

import (
	"context"
	"github.com/Juno-chat-app/user-service/config"
	"github.com/Juno-chat-app/user-service/domain/repository/redis"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	log   logger.ILogger
	conf  *config.Configuration
	cache redis.ICache
)

func TestMain(m *testing.M) {
	conf, err := config.LoadConfiguration("./cache_test_config.yml")
	if err != nil {
		os.Exit(1)
	}

	log, err = logger.NewLogger()
	if err != nil {
		os.Exit(1)
	}

	cacheConfig := conf.CQRSConfig.CacheConfig
	cache = redis.NewCache(cacheConfig.Host, cacheConfig.Port, cacheConfig.Password, cacheConfig.Db, cacheConfig.Retry, log)

	code := m.Run()
	os.Exit(code)
}

func Test_Ping_Set_Get(t *testing.T) {
	ctx := context.Background()

	err := cache.Ping(ctx)
	require.Nil(t, err)

	err = cache.Set(ctx, "test", "value", time.Second)
	require.Nil(t, err)

	value, err := cache.Get(ctx, "test")
	require.Nil(t, err)
	require.Equal(t, "value", value)

	err = cache.Remove(ctx, "test")
	require.Nil(t, err)

	_, err = cache.Get(ctx, "test")
	stat, ok := status.FromError(err)
	require.Equal(t, true, ok)
	require.Equal(t, codes.Code(http.StatusNotFound), stat.Code())
}
