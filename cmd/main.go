package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/goggle-source/authLotServic/internal/app"
	"github.com/goggle-source/authLotServic/internal/config"
)

/*
Добавить в этот сервис индексы, пул соединений
так же добавить проверку контекста
добавить статус для check, чтобы вместе с метрики отдавать статус сервера(работает, нагружен, не работает)
*/

func main() {
	cfg := config.Load()
	log := SetupLogger(cfg.Env)
	log.Info("start servic")

	application := app.New(log, cfg)

	application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "Local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case "Prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
