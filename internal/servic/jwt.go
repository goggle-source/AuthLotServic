package servic

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/goggle-source/authLotServic/domain"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string
	jwt.RegisteredClaims
}

type servicJWT struct {
	secretKey *rsa.PrivateKey
}

func InitServicJWT(key *rsa.PrivateKey) *servicJWT {
	return &servicJWT{
		secretKey: key,
	}
}

func (s *servicJWT) GenerateJWTToken(ctx context.Context, id string) (token string, err error) {
	claims := Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err = tokenJWT.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("token signing error:%s", domain.ErrGenerateJWT)
	}

	return token, nil
}
