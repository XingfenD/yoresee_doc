package utils

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost := bcrypt.DefaultCost
	if config.GlobalConfig != nil {
		c := config.GlobalConfig.Backend.Security.PasswordHashCost
		if c >= bcrypt.MinCost && c <= bcrypt.MaxCost {
			cost = c
		}
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
