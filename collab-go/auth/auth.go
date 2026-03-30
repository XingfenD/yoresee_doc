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

	if a.secret != "" {
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.secret), nil
		}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
		return err
	}

	// TODO: jwt redis support
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	return err
}
