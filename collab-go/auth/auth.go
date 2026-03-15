package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	secret string
}

func NewAuthenticator(secret string) *Authenticator {
	return &Authenticator{
		secret: secret,
	}
}

func (a *Authenticator) ValidateToken(tokenString string) error {
	if tokenString == "" {
		return http.ErrNoCookie
	}

	// 如果配置了本地 JWT 密钥，则使用本地验证
	if a.secret != "" {
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.secret), nil
		}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
		return err
	}

	// 由于 backend 没有专门的 token 验证接口，我们进行基本的 token 格式验证
	// 这里可以根据实际情况扩展，比如检查 token 是否过期等
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	return err
}
