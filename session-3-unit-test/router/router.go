package router

import (
	"training-golang/session-3-unit-test/handler"
	"training-golang/session-3-unit-test/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/", handler.RootHandler)

	privateEndpoint := r.Group("v1/api")
	privateEndpoint.Use(middleware.AuthMiddleware())
	{
		privateEndpoint.POST("/post", handler.PostHandler)
	}
}
