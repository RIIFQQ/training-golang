package main

import (
	"context"
	"log"
	"training-golang/session-6-db-pgx-crud/handler"
	posgrespgx "training-golang/session-6-db-pgx-crud/repository/posgres_pgx"
	"training-golang/session-6-db-pgx-crud/router"
	"training-golang/session-6-db-pgx-crud/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	//setup service
	dsn := "postgresql://postgres:admin@localhost:5432/go_db"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalln(err)
	}

	userRepo := posgrespgx.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	//setup router
	router.SetupRouter(r, userHandler)

	//setup service
	log.Println("running server on port 8080")
	r.Run(":8080")
}
