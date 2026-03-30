package auth

import (
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTValidator struct{}

func (v *JWTValidator) Validate(tokenString string) (*Claims, error) {
	return ParseToken(tokenString)
}

func (v *JWTValidator) IsExpired(claims *Claims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

type Claims struct {
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(externalID string, username string) (string, error) {
	cfg := config.GlobalConfig.Backend.Jwt
	now := time.Now()
	expireTime := now.Add(time.Duration(cfg.Expire) * time.Second)

	claims := Claims{
		ExternalID: externalID,
		Username:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "doc_manager",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	if err := StoreJWTToken(externalID, token, time.Until(expireTime)); err != nil {
		return "", fmt.Errorf("store jwt token failed: %w", err)
	}

	return token, nil
}

func ParseToken(token string) (*Claims, error) {
	cfg := config.GlobalConfig.Backend.Jwt
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
