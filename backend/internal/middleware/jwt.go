package middleware

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
)

var JWTAuth = &JWTAuthMiddleware{}

type JWTAuthMiddleware struct {
}

func (m *JWTAuthMiddleware) ValidateAuthorizationHeader(authHeader string) (*utils.Claims, error) {
	return m.handle(authHeader)
}

func (m *JWTAuthMiddleware) handle(authHeader string) (*utils.Claims, error) {
	if authHeader == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusTokenInvalid, "unlogin or illegal access")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, status.GenErrWithCustomMsg(status.StatusTokenInvalid, "invalid token format")
	}
	token := parts[1]
	jwtValidator := &utils.JWTValidator{}
	claims, err := jwtValidator.Validate(token)
	if err != nil {
		return nil, status.StatusTokenInvalid
	}

	// TODO: jwt + redis
	if jwtValidator.IsExpired(claims) {
		return nil, status.StatusTokenExpired
	}
	return claims, nil
}
