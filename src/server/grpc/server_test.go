package grpc

import (
	"context"
	"fmt"
	userproto "github.com/Juno-chat-app/user-proto"
	"github.com/Juno-chat-app/user-service/config"
	"github.com/Juno-chat-app/user-service/domain/model/authorization"
	"github.com/Juno-chat-app/user-service/domain/model/services"
	"github.com/Juno-chat-app/user-service/domain/repository/mongo"
	"github.com/Juno-chat-app/user-service/domain/repository/redis"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

var (
	server *Server
	client userproto.UserServiceClient
)

func TestMain(m *testing.M) {
	var conf *config.Configuration
	log, err := logger.NewLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conf, err = config.LoadConfiguration("../../user-service_test.yml")
	if err != nil {
		log.Error("got error on loading configuration", "err", err)
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
	server = NewServer(conf.GRPCConfig.Host, conf.GRPCConfig.Port, service, log)
	go func() {
		err := server.Start()
		if err != nil {
			os.Exit(1)
		}
	}()

	// warm-up server
	time.Sleep(time.Second)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.GRPCConfig.Host, conf.GRPCConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	client = userproto.NewUserServiceClient(conn)

	code := m.Run()
	os.Exit(code)
}

func TestServer_SignUp(t *testing.T) {
	req := userproto.SignUpRequest{
		UserName: "test",
		Email:    "test",
		Password: "test",
	}
	reqBody, err := proto.Marshal(&req)
	require.Nil(t, err)

	request := userproto.RequestMessage{
		Name:   "",
		Type:   "",
		Time:   "",
		Header: nil,
		Body: &userproto.Any{
			TypeUrl: "SignUpRequest",
			Value:   reqBody,
		},
	}

	_, err = client.SignUp(context.Background(), &request)
	require.Nil(t, err)
}
