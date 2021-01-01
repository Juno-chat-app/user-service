package grpc

import (
	"context"
	userproto "github.com/Juno-chat-app/user-proto"
	"github.com/Juno-chat-app/user-service/domain/entity"
	"github.com/Juno-chat-app/user-service/domain/model/authorization"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *Server) SignIn(ctx context.Context, req *userproto.RequestMessage, ) (*userproto.ResponseMessage, error) {
	if req.Body.TypeUrl != SignInRequestMethod {
		return nil, status.Error(http.StatusBadRequest, "request content type must be SignInRequest")
	}

	body := userproto.SignInRequest{}
	err := proto.Unmarshal(req.Body.Value, &body)
	if err != nil {
		return nil, status.Error(http.StatusBadRequest, "request body type is not SignInRequest")
	}

	user := entity.User{
		UserName: body.UserName,
		Password: body.Password,
	}

	signInResult, err := s.userService.SignIn(ctx, &user)
	if err != nil {
		return nil, err
	}

	responseBody := userproto.Response{
		BearerToken:  signInResult.AccessToken,
		Duration:     signInResult.ExpireAt,
		RefreshToken: signInResult.RefreshToke,
	}

	responseBodyByte, err := proto.Marshal(&responseBody)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "got error on generating response")
	}

	response := userproto.ResponseMessage{
		Entity: "Response",
		Meta:   nil,
		Data: &userproto.Any{
			TypeUrl: "Response",
			Value:   responseBodyByte,
		},
	}

	return &response, nil
}

func (s *Server) SignUp(ctx context.Context, req *userproto.RequestMessage) (*userproto.ResponseMessage, error) {
	if req.Body.TypeUrl != SignUpRequestMethod {
		return nil, status.Error(http.StatusBadRequest, "request content type must be SingUpRequest")
	}

	body := userproto.SignUpRequest{}
	err := proto.Unmarshal(req.Body.Value, &body)
	if err != nil {
		return nil, status.Error(http.StatusBadRequest, "request body type is not SignUpRequest")
	}

	user := entity.User{
		UserName: body.UserName,
		Password: body.Password,
		ContactInfo: &entity.ContactInfo{
			Email: body.Email,
		},
	}

	_, err = s.userService.SignUp(ctx, &user)
	if err != nil {
		return nil, err
	}

	response := userproto.ResponseMessage{
		Entity: "SignUpResponse",
		Meta:   nil,
		Data:   nil,
	}

	return &response, nil
}

func (s *Server) Verify(ctx context.Context, req *userproto.RequestMessage) (*userproto.ResponseMessage, error) {
	if req.Body.TypeUrl != ValidateRequestMethod {
		return nil, status.Error(http.StatusBadRequest, "request content type must be ValidateRequest")
	}

	body := userproto.ValidateRequest{}
	err := proto.Unmarshal(req.Body.Value, &body)
	if err != nil {
		return nil, status.Error(http.StatusBadRequest, "request body type is not ValidateRequest")
	}

	if req.Header.UID == "" {
		return nil, status.Error(http.StatusBadRequest, "invalid user-id")
	}

	verificationToken := authorization.TokenDetail{
		AccessToken: body.BearerToken,
		UserId:      req.Header.UID,
	}

	_, err = s.userService.Validate(ctx, &verificationToken)

	if err != nil {
		return nil, err
	}

	response := userproto.ResponseMessage{
		Entity: "ValidationResponse",
		Meta:   nil,
		Data:   nil,
	}

	return &response, nil
}

func (s *Server) Refresh(ctx context.Context, req *userproto.RequestMessage) (*userproto.ResponseMessage, error) {
	if req.Body.TypeUrl != RefreshRequestMethod {
		return nil, status.Error(http.StatusBadRequest, "request content type must be RefreshRequest")
	}

	body := userproto.RefreshRequest{}
	err := proto.Unmarshal(req.Body.Value, &body)
	if err != nil {
		return nil, status.Error(http.StatusBadRequest, "request body type is not RefreshRequest")
	}

	if req.Header.UID == "" {
		return nil, status.Error(http.StatusBadRequest, "invalid user-id")
	}

	verificationToken := authorization.TokenDetail{
		AccessToken: body.RefreshToken,
		UserId:      req.Header.UID,
	}

	refreshToken, err := s.userService.RefreshToken(ctx, &verificationToken)

	if err != nil {
		return nil, err
	}

	responseBody := userproto.Response{
		BearerToken:  refreshToken.AccessToken,
		Duration:     refreshToken.ExpireAt,
		RefreshToken: refreshToken.RefreshToke,
	}

	responseBodyByte, err := proto.Marshal(&responseBody)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "got error on generating response")
	}

	response := userproto.ResponseMessage{
		Entity: "Response",
		Meta:   nil,
		Data: &userproto.Any{
			TypeUrl: "Response",
			Value:   responseBodyByte,
		},
	}

	return &response, nil
}

func (s *Server) GetUser(ctx context.Context, req *userproto.RequestMessage) (*userproto.ResponseMessage, error) {
	if req.Body.TypeUrl != GetUserRequestMethod {
		return nil, status.Error(http.StatusBadRequest, "request content type must be GetUserRequest")
	}

	body := userproto.GetUserRequest{}
	err := proto.Unmarshal(req.Body.Value, &body)
	if err != nil {
		return nil, status.Error(http.StatusBadRequest, "request body type is not GetUserRequest")
	}

	contactInfo := entity.ContactInfo{
		Mobile: body.SearchInfo.PhoneNumber,
		Email:  body.SearchInfo.Email,
	}

	user, err := s.userService.GetUser(ctx, contactInfo)
	if err != nil {
		return nil, err
	}

	responseBody := userproto.GetUserResponse{
		Info: &userproto.GetUserResponse_UserInfo{
			UserName:    user.UserName,
			Email:       user.ContactInfo.Email,
			PhoneNumber: user.ContactInfo.Mobile,
			UserId:      user.UserId,
		},
	}

	responseBodyByte, err := proto.Marshal(&responseBody)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, "got error on generating response")
	}

	response := userproto.ResponseMessage{
		Entity: "Response",
		Meta:   nil,
		Data: &userproto.Any{
			TypeUrl: "GetUserResponse",
			Value:   responseBodyByte,
		},
	}
	
	return &response, nil
}
