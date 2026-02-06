package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/go-playground/validator/v10"
	grpcServ "github.com/goggle-source/authLotServic/internal/grpc"
	"github.com/goggle-source/authLotServic/internal/grpc/interceptor"
	"github.com/goggle-source/authLotServic/internal/lib/logger"
	"github.com/goggle-source/authLotServic/internal/servic"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log  *slog.Logger
	GRPC *grpc.Server
	Port int
}

func Init(log *slog.Logger, authServic servic.ServicApp, port int, validate *validator.Validate) *App {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.LoggingInterceptor,
		recovery.UnaryServerInterceptor(recoveryOpts...),
	))

	grpcServ.Register(gRPCServer, &authServic, log, validate)

	return &App{
		log:  log,
		GRPC: gRPCServer,
		Port: port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("start run server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))
	if err != nil {
		log.Error("error listen TCP port", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.GRPC.Serve(l); err != nil {
		log.Error("error start grpc server", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.Port))

	a.GRPC.GracefulStop()
}
