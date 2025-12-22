package main

import (
	"log/slog"
	"os"

	"github.com/goggle-source/authLotServic/internal/config"
)

/*
Добавить в этот сервис индексы, пул соединений
*/

func main() {
	cfg := config.Load()
	log := SetupLogger(cfg.Env)
	log.Info("start servic")
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
