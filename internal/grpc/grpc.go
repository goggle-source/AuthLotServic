package grpc

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"

	"github.com/goggle-source/authLotProto/gen/go/auth"
	"github.com/goggle-source/authLotServic/internal/lib/logger"
	"github.com/goggle-source/authLotServic/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer interface {
	Register(ctx context.Context, userRequest models.UserRegister) (token string, err error)
	Login(ctx context.Context, userLogin models.UserLogin) (name string, token string, err error)
	HealthyCheack(ctx context.Context) (map[string]string, error)
}

type ServerAPI struct {
	auth.UnimplementedAuthServer
	auth     AuthServer
	log      *slog.Logger
	validate *validator.Validate
}

func Register(grpc *grpc.Server, authServ AuthServer, log *slog.Logger, val *validator.Validate) {
	auth.RegisterAuthServer(grpc, &ServerAPI{auth: authServ, log: log, validate: val})
}

func (s *ServerAPI) Login(ctx context.Context, in *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	const op = "grpc.Login"

	log := s.log.With(slog.String("op", op), slog.String("email", in.GetEmail()))

	userLogin := models.UserLogin{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	if err := s.validate.Struct(userLogin); err != nil {
		log.Error("error validate date", logger.Err(err))
		return nil, ValidationErrValidator(err)
	}

	name, token, err := s.auth.Login(ctx, userLogin)

	if err != nil {
		log.Error("error login user", logger.Err(err))
		return nil, ValidationError(err)
	}

	return &auth.LoginUserResponse{
		Token: token,
		Name:  name,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, in *auth.RegisterUserRequest) (*auth.RegisterUserResponse, error) {
	const op = "grpc.Register"

	log := s.log.With(slog.String("op", op), slog.String("email", in.GetEmail()))

	log.Info("start register")

	userRegister := models.UserRegister{
		Name:     in.GetName(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	if err := s.validate.Struct(userRegister); err != nil {
		log.Error("error validate date", logger.Err(err))
		return nil, ValidationErrValidator(err)
	}

	token, err := s.auth.Register(ctx, userRegister)
	if err != nil {
		log.Error("error register", logger.Err(err))
		return nil, ValidationError(err)
	}

	return &auth.RegisterUserResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Health(ctx context.Context, in *auth.HealthCheckRequest) (*auth.HealthCheckResponse, error) {
	const op = "grpc.Health"

	log := s.log.With(slog.String("op", op))

	details, err := s.auth.HealthyCheack(ctx)
	if err != nil {
		log.Error("error healthyCheack in service layer", logger.Err(err))
		return &auth.HealthCheckResponse{
			Details: map[string]string{
				"err": "an error has occurred",
			},
		}, status.Error(codes.Internal, "server error")
	}

	return &auth.HealthCheckResponse{
		Details: details,
	}, nil
}
