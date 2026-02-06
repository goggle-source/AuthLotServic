package interceptor

import (
	"context"
	"crypto/rsa"
	"log/slog"
	"os"
	"sync"

	"github.com/goggle-source/authLotServic/internal/lib/logger"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
)

var (
	privateKey    *rsa.PrivateKey
	secretKeyOnce sync.Once
)

type Claims struct {
	Userid string
	jwt.RegisteredClaims
}

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("request processing starts", slog.String("path", info.FullMethod))

	resp, err = handler(ctx, req)
	if err != nil {
		log.Error("request processing error", logger.Err(err))
	}

	log.Info("end of request processing", slog.Any("result", resp))

	return resp, err
}
