package grpc

import (
	"fmt"
	userproto "github.com/Juno-chat-app/user-proto"
	"github.com/Juno-chat-app/user-service/domain/model/services"
	"github.com/Juno-chat-app/user-service/infra/logger"
	"google.golang.org/grpc"
	"net"
)

const (
	SignInRequestMethod   string = "SignInRequest"
	SignUpRequestMethod   string = "SignUpRequest"
	ValidateRequestMethod string = "ValidateRequest"
	RefreshRequestMethod  string = "RefreshRequest"
)

type Server struct {
	address     string
	port        int32
	userService services.IUserService
	logger      logger.ILogger
}

func NewServer(address string, port int32, userService services.IUserService, logger logger.ILogger) *Server {
	server := Server{
		address:     address,
		port:        port,
		userService: userService,
		logger:      logger,
	}

	return &server
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d", s.address, s.port)
	lis, err := net.Listen("tcp", address)

	if err != nil {
		s.logger.Error("Got error on creating listener",
			"method", "Start",
			"host", s.address,
			"port", s.port,
			"error", err)

		return err
	}

	s.logger.Info("gRPC Server started",
		"method", "Start",
		"host", s.address,
		"port", s.port)

	grpcServer := grpc.NewServer()
	userproto.RegisterUserServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Error("Got error on listening",
			"method", "Start",
			"host", s.address,
			"port", s.port,
			"error", err)

		return err
	}

	return nil
}
