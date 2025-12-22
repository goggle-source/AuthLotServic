package grpc

import (
	"context"

	auth "github.com/goggle-source/authLotProto/gen/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer interface {
	Register(ctx context.Context, email string, password string, name string) (token string, err error)
	Login(ctx context.Context, email string, password string) (name string, token string, err error)
	ValidateToken(ctx context.Context, token string) (isValid bool, err error)
	Check(ctx context.Context) (details map[string]string, status int, err error)
}

type ServerAPI struct {
	auth.UnimplementedAuthServer
	auth AuthServer
}

func (s *ServerAPI) Login(ctx context.Context, in *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	name, token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())

	if err != nil {
		//Добавить валидацию ошибок
		return nil, status.Error(codes.Internal, "server error")
	}

	return &auth.LoginUserResponse{
		Token: token,
		Name:  name,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, in *auth.RegisterUserRequest) (*auth.RegisterUserResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if len(in.GetName()) < 3 {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	token, err := s.auth.Register(ctx, in.GetEmail(), in.GetPassword(), in.GetName())
	if err != nil {
		//Добавить валидацию ошибок
		return nil, status.Error(codes.Internal, "server error")
	}

	return &auth.RegisterUserResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) ValidateToken(ctx context.Context, in *auth.ValidTokenRequest) (*auth.ValidBoolResponse, error) {
	if in.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	ok, err := s.auth.ValidateToken(ctx, in.GetToken())
	if err != nil {
		//Добавить валидацию ошибок
		return nil, status.Error(codes.Internal, "server error")
	}

	return &auth.ValidBoolResponse{
		IsValid: ok,
	}, nil
}

func (s *ServerAPI) Health(ctx context.Context, in *auth.HealthCheckRequest) (*auth.HealthCheckResponse, error) {

	details, statusServic, err := s.auth.Check(ctx)
	if err != nil {
		//Добавить валидацию ошибок
		return &auth.HealthCheckResponse{
			Details: map[string]string{
				"err": "an error has occurred",
			},
		}, status.Error(codes.Internal, "server error")
	}
	/*
		добавить проверку на статус сервиса, где 1 это FREE, 2 SUCCESS и 3 это LIMIT и менять статус в ответе в соответствие с возвращаемым значением
	*/
	_ = statusServic

	return &auth.HealthCheckResponse{
		Details: details,
		Status:  auth.HealthCheckResponse_FREE,
	}, nil
}
