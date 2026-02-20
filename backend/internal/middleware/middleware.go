package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	GinHandle() gin.HandlerFunc
}
