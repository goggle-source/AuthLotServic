package grpc

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	auth "github.com/goggle-source/authLotProto/gen/go/auth"
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
	ValidateUser(ctx context.Context, userID int) (bool, error)
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

	log := s.log.With(slog.String("op", op))

	log.Info("start Login")

	userLogin := models.UserLogin{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	if err := s.validate.Struct(userLogin); err != nil {
		return nil, ValidationErrValidator(err)
	}

	name, token, err := s.auth.Login(ctx, userLogin)

	if err != nil {
		log.Error("error login user", logger.Err(err))
		return nil, ValidationError(err)
	}
	log.Info("success login")

	return &auth.LoginUserResponse{
		Token: token,
		Name:  name,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, in *auth.RegisterUserRequest) (*auth.RegisterUserResponse, error) {
	const op = "grpc.Register"

	log := s.log.With(slog.String("op", op))

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
	log.Info("success register")

	return &auth.RegisterUserResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Health(ctx context.Context, in *auth.HealthCheckRequest) (*auth.HealthCheckResponse, error) {
	const op = "grpc.Health"

	log := s.log.With(slog.String("op", op))

	log.Info("start health")

	details, err := s.auth.HealthyCheack(ctx)
	if err != nil {
		log.Error("error healthyCheack in service layer", logger.Err(err))
		return &auth.HealthCheckResponse{
			Details: map[string]string{
				"err": "an error has occurred",
			},
		}, status.Error(codes.Internal, "server error")
	}
	log.Info("success health")

	return &auth.HealthCheckResponse{
		Details: details,
	}, nil
}

func (s *ServerAPI) ValidateUserId(ctx context.Context, in *auth.UserIdRequest) (*auth.ValidIsIdResponse, error) {
	const op = "grpc.ValidateUserId"

	log := s.log.With(slog.String("op", op))

	log.Info("start validateUserId")

	if in.Id == 0 {
		log.Error("is not id")
		return &auth.ValidIsIdResponse{}, status.Error(codes.InvalidArgument, "id is required")
	}

	isValid, err := s.auth.ValidateUser(ctx, int(in.Id))
	if err != nil {
		log.Error("error validateUser", logger.Err(err))
		return &auth.ValidIsIdResponse{}, ValidationError(err)
	}
	log.Info("success validateUserId")

	return &auth.ValidIsIdResponse{
		IsValidId: isValid,
	}, nil
}
