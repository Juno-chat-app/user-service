package redis

import (
	"context"
	"fmt"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type ICache interface {
	Ping(ctx context.Context) (err error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) (err error)
	Get(ctx context.Context, key string) (value string, err error)
	Remove(ctx context.Context, key string) (err error)
}

func NewCache(address string, port int32, password string, db int, retry int32, logger logger.ILogger) ICache {
	logger.Info("Create cache-client",
		"method", "NewCache",
		"port", port,
		"host", address,
		"db", db,
		"retry", retry)

	connection := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", address, port),
		Password: password,
		DB:       db,
	})

	cache := iCache{
		address:    address,
		port:       port,
		password:   password,
		db:         db,
		logger:     logger,
		connection: connection,
		retry:      retry,
		sync:       sync.Mutex{},
	}

	return &cache
}
