package middleware

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/api"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var JWTAuth = &JWTAuthMiddleware{}

type JWTAuthMiddleware struct {
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

	if jwtValidator.IsExpired(claims) {
		return nil, status.StatusTokenExpired
	}
	return claims, nil
}

func (m *JWTAuthMiddleware) GinHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		logrus.Infof("authHeader: %s", authHeader)
		claims, err := m.handle(authHeader)
		if err != nil {
			c.JSON(401, api.GenBaseRespWithErr(err))
			c.Abort()
			return
		}

		c.Set("user_external_id", claims.ExternalID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
