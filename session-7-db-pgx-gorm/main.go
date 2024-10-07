package main

import (
	"log"
	"training-golang/session-7-db-pgx-gorm/handler"
	posgresgorm_raw "training-golang/session-7-db-pgx-gorm/repository/posgres_gom_raw"
	"training-golang/session-7-db-pgx-gorm/router"
	"training-golang/session-7-db-pgx-gorm/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	//setup service
	dsn := "postgresql://postgres:admin@localhost:5432/go_db"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		log.Fatalln(err)
	}
	//  pake postgresgorm.
	//userRepo := posgresgorm.NewUserRepository(gormDB)
	userRepo := posgresgorm_raw.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	//setup router
	router.SetupRouter(r, userHandler)

	//setup service
	log.Println("running server on port 8080")
	r.Run(":8080")
}
