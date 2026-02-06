package app

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/goggle-source/authLotServic/cmd/migrate"
	grpcapp "github.com/goggle-source/authLotServic/internal/app/grpc/app"
	"github.com/goggle-source/authLotServic/internal/config"
	secretkey "github.com/goggle-source/authLotServic/internal/lib/SecretKey"
	"github.com/goggle-source/authLotServic/internal/repository"
	"github.com/goggle-source/authLotServic/internal/servic"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Cfg) *App {
	db := repository.Init(cfg, log)

	err := migrate.RunMigrations(cfg)
	if err != nil {
		panic(err)
	}

	validate := validator.New()

	privateKey, err := secretkey.LoadPrivateKey(cfg.Path)
	if err != nil {
		panic(err)
	}

	bisnesServic := servic.Init(log, db, privateKey)

	grpcServer := grpcapp.Init(log, *bisnesServic, cfg.GRPC.Port, validate)

	return &App{
		GRPCServer: grpcServer,
	}
}
