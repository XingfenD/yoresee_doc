package router

import (
	"github.com/XingfenD/yoresee_doc/internal/api"
	"github.com/XingfenD/yoresee_doc/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	root := r.Group("/")

	// RegisterController(root, "GET", "/health", controller.HealthControllerImpl.Name(), controller.HealthControllerImpl.Handle)
	root.GET("/health", api.HealthHandlerImpl.GinHander())

	api := root.Group("/api")

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// 示例：注册受保护的路由
		// RegisterController(protected, "GET", "/user/profile", authCtrl.Name(), authCtrl.GetProfile)
		// RegisterController(protected, "GET", "/document/list", docCtrl.Name(), docCtrl.List)
	}
}
