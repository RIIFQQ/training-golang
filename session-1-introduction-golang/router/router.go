package router

import (
	"training-golang/session-1-introduction-golang/handler"
	"training-golang/session-1-introduction-golang/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/", handler.RootHandler)

	privateRoute := r.Group("/api/v1")
	privateRoute.Use(middleware.AuthMiddleware())
	{
		privateRoute.POST("/post", handler.PostHandler)
	}
}
