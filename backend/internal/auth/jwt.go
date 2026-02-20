package auth

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/utils"
)

type JWTValidator struct{}

func (v *JWTValidator) Validate(tokenString string) (*utils.Claims, error) {
	return utils.ParseToken(tokenString)
}

func (v *JWTValidator) IsExpired(claims *utils.Claims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}
