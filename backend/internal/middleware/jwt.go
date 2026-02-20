package middleware

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/auth"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "未登录或非法访问"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{"error": "请求头中auth格式有误"})
			c.Abort()
			return
		}

		jwtValidator := &auth.JWTValidator{}
		claims, err := jwtValidator.Validate(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的Token"})
			c.Abort()
			return
		}

		if jwtValidator.IsExpired(claims) {
			c.JSON(401, gin.H{"error": "Token已过期"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role_id", claims.RoleID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
