package grpc

import (
	"context"
	userproto "github.com/Juno-chat-app/user-proto"
)

func (s *Server) SignIn(ctx context.Context, req *userproto.RequestMessage, ) (*userproto.ResponseMessage, error) {
	panic("not implemented")
}

func (s *Server) SignUp(ctx context.Context, req *userproto.RequestMessage) (*userproto.ResponseMessage, error) {
	panic("not implemented")
}

func (s *Server) Verify(ctx context.Context, req *userproto.RequestMessage) (*userproto.ResponseMessage, error) {
	panic("not implemented")
}

func (s *Server) Refresh(ctx context.Context, req *userproto.RequestMessage) (*userproto.ResponseMessage, error) {
	panic("not implemented")
}
