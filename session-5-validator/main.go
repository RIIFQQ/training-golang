package main

import (
	"log"
	"training-golang/session-5-validator/entity"
	"training-golang/session-5-validator/handler"
	"training-golang/session-5-validator/repository/slice"
	"training-golang/session-5-validator/router"
	"training-golang/session-5-validator/service"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	//setup service
	var mockUserDBInSlice []entity.User
	userRepo := slice.NewUserRepository(mockUserDBInSlice)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	//setup router
	router.SetupRouter(r, userHandler)

	//setup service
	log.Println("running server on port 8080")
	r.Run(":8080")
}
