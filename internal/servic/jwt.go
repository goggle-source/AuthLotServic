package servic

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string
	jwt.RegisteredClaims
}

func GenerateJWTToken(ctx context.Context, id string, sercretKey *rsa.PrivateKey) (token string, err error) {
	claims := Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err = tokenJWT.SignedString(sercretKey)
	if err != nil {
		return "", fmt.Errorf("token signing error")
	}

	return token, nil
}
