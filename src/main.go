package main

import (
	"context"
	"fmt"
	"github.com/Juno-chat-app/user-service/config"
	"github.com/Juno-chat-app/user-service/domain/model/authorization"
	"github.com/Juno-chat-app/user-service/domain/model/services"
	"github.com/Juno-chat-app/user-service/domain/repository/mongo"
	"github.com/Juno-chat-app/user-service/domain/repository/redis"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/Juno-chat-app/user-service/server/grpc"
	"os"
	"time"
)

func main() {
	mode := os.Getenv("APP_MODE")
	var conf *config.Configuration
	log, err := logger.NewLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if mode == "dev" {
		conf, err = config.LoadConfiguration("./user-service_test.yml")
		if err != nil {
			log.Error("got error on loading configuration", "err", err)
			os.Exit(1)
		}
	} else if mode == "run" {
		conf, err = config.LoadConfiguration("./user-service.yml")
		if err != nil {
			log.Error("got error on loading configuration", "err", err)
			os.Exit(1)
		}
	} else {
		log.Error("APP_MODE is not dev or run")
		os.Exit(1)
	}

	auth := authorization.NewJwtHandler(int64(conf.AuthConfig.AccessTTL*time.Minute), int64(conf.AuthConfig.RefreshTTL*time.Hour), log)

	cacheConfig := conf.CQRSConfig.CacheConfig
	cache := redis.NewCache(cacheConfig.Host, cacheConfig.Port, cacheConfig.Password, cacheConfig.Db, cacheConfig.Retry, log)
	err = cache.Ping(context.Background())
	if err != nil {
		log.Error("got error on redis connection", "err", err)
		os.Exit(1)
	}

	repo := mongo.NewUserRepository(conf.CQRSConfig.PersistConfig, log)
	err = repo.Ping(context.Background())
	if err != nil {
		log.Error("got error on mongo connection", "err", err)
		os.Exit(1)
	}

	service := services.NewUserService(log, repo, cache, auth)
	grpcServer := grpc.NewServer(conf.GRPCConfig.Host, conf.GRPCConfig.Port, service, log)
	err = grpcServer.Start()
	if err != nil {
		log.Error("error on running grpc server", "err", err)
		os.Exit(1)
	}
}
