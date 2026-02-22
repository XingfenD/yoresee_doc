package router

import (
	"github.com/XingfenD/yoresee_doc/internal/api"
	"github.com/XingfenD/yoresee_doc/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	root := r.Group("/")

	root.GET("/health", api.HealthHandlerImpl.GinHandle())

	root.GET("/system-info", api.SystemInfoHandlerImpl.GinHandle())
	root.POST("/test/post", api.TestPostHandlerImpl.GinHandle())
	root.POST("/login", api.AuthLoginHandlerImpl.GinHandle())
	root.POST("/register", api.AuthRegisterHandlerImpl.GinHandle())

	// protected routes
	protected := root.Group("/")
	protected.Use(middleware.JWTAuth.GinHandle())
	{
		protected.GET("/document/:documentExternalID/content", api.GetDocumentContentHandlerImpl.GinHandle())
		protected.GET("/test/protected", api.TestProtectedHandlerImpl.GinHandle())
	}
}
