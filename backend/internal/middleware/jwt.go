package middleware

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/auth"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

var JWTAuth = &JWTAuthMiddleware{}

type JWTAuthMiddleware struct {
}

func (m *JWTAuthMiddleware) ValidateAuthorizationHeader(authHeader string) (*auth.Claims, error) {
	return m.handle(authHeader)
}

func (m *JWTAuthMiddleware) handle(authHeader string) (*auth.Claims, error) {
	if authHeader == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusTokenInvalid, "unlogin or illegal access")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, status.GenErrWithCustomMsg(status.StatusTokenInvalid, "invalid token format")
	}
	token := parts[1]
	claims, err := auth.ParseToken(token)
	if err != nil {
		return nil, status.StatusTokenInvalid
	}

	if (&auth.JWTValidator{}).IsExpired(claims) {
		return nil, status.StatusTokenExpired
	}
	if err := auth.ValidateJWTTokenInRedis(claims.ExternalID, token); err != nil {
		return nil, err
	}
	return claims, nil
}
