package router

import (
	"github.com/XingfenD/yoresee_doc/internal/api"
	"github.com/XingfenD/yoresee_doc/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	root := r.Group("/")

	root.GET("/health", api.HealthHandlerImpl.GinHandle())

	api_ := root.Group("/api")

	api_.POST("/test/post", api.TestPostHandlerImpl.GinHandle())

	// protected routes
	protected := api_.Group("/")
	protected.Use(middleware.JWTAuth.GinHandle())
	{
		protected.GET("/test/protected", api.TestProtectedHandlerImpl.GinHandle())
	}
}
