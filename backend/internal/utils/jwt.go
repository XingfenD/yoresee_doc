package utils

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	RoleID   int64  `json:"role_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, roleID int64, username string) (string, error) {
	cfg := config.GlobalConfig.Backend.Jwt
	now := time.Now()
	expireTime := now.Add(time.Duration(cfg.Expire) * time.Second)

	claims := Claims{
		UserID:   userID,
		RoleID:   roleID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "doc_manager",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(cfg.Secret))
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
