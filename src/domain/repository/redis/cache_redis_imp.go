package redis

import (
	"context"
	"fmt"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/status"
	"net/http"
	"sync"
	"time"
)

const (
	Retry string = "retry"
)

type iRedisCache struct {
	address    string
	port       int32
	password   string
	db         int
	logger     logger.ILogger
	connection *redis.Client
	retry      int32
	sync       sync.Mutex
}

func (c *iRedisCache) Ping(ctx context.Context) (err error) {
	_, err = c.connection.Ping(ctx).Result()

	if err != nil {
		err := c.reconnect(ctx)
		if err != nil {
			return err
		} else {
			return c.Ping(ctx)
		}
	}

	return nil
}

func (c *iRedisCache) Remove(ctx context.Context, key string) (err error) {
	err = c.connection.Del(ctx, key).Err()

	if err != nil {
		err := c.reconnect(ctx)
		if err != nil {
			return err
		} else {
			return c.Remove(ctx, key)
		}
	}

	return nil
}

func (c *iRedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) (err error) {
	err = c.connection.Set(ctx, key, value, expiration).Err()

	if err != nil {
		err := c.reconnect(ctx)
		if err != nil {
			return err
		} else {
			return c.Set(ctx, key, value, expiration)
		}
	}

	return nil
}

func (c *iRedisCache) Get(ctx context.Context, key string) (value string, err error) {
	value, err = c.connection.Get(ctx, key).Result()
	if err != redis.Nil && err != nil {
		err := c.reconnect(ctx)
		if err != nil {
			return "", err
		} else {
			return c.Get(ctx, key)
		}
	} else if err == redis.Nil {
		return "", status.Error(http.StatusNotFound, "key not found")
	}

	return value, nil
}

func (c *iRedisCache) reconnect(ctx context.Context) (err error) {
	val := ctx.Value(Retry)
	if val == nil {
		c.connection = nil
		if c.connection == nil {
			c.sync.Lock()
			defer c.sync.Unlock()
			if c.connection == nil {
				retryNContext := context.WithValue(ctx, Retry, int32(0))
				return c.reconnect(retryNContext)
			}
		}
	} else {
		retry, ok := val.(int32)
		if !ok {
			c.logger.Error("Got error on retry key-type",
				"method", "reconnect",
				"retry", retry)

			return status.Error(http.StatusInternalServerError, "Got error in cache layer")
		}

		if retry < c.retry {
			c.connection = redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", c.address, c.port),
				Password: c.password,
				DB:       c.db,
			})

			err := c.connection.Ping(ctx).Err()
			if err != nil {
				c.connection = nil
				retry = retry + 1
				retryNPContext := context.WithValue(ctx, Retry, retry)

				c.logger.Error("Retry for reconnection",
					"method", "reconnect",
					"retry", retry)

				return c.reconnect(retryNPContext)
			}

			return nil
		} else {
			c.logger.Error("Retry context deadline exceeded",
				"method", "retry",
				"retry", retry)

			return status.Error(http.StatusInternalServerError, "Retry context deadline exceeded")
		}
	}

	return
}
