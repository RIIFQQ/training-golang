package main

import (
	"context"
	"log"
	"net"
	"session-16-crud-user-docker-compose/entity"
	grpcHandler "session-16-crud-user-docker-compose/handler/grpc"
	"session-16-crud-user-docker-compose/middleware"
	pb "session-16-crud-user-docker-compose/proto/user_service/v1"
	posgresgorm_raw "session-16-crud-user-docker-compose/repository/posgres_gom_raw"
	"session-16-crud-user-docker-compose/service"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	//setup service
	dsn := "postgresql://postgres:password@pg-db:5432/go_db_crud"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		log.Fatalln(err)
	}

	//setup redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-db:6379",
		Password: "redispass",
		DB:       0,
	})

	//run migration
	if err := gormDB.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalln("Failed to migrate database:", err)
	}

	//  pake postgresgorm.
	//userRepo := posgresgorm.NewUserRepository(gormDB)
	userRepo := posgresgorm_raw.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo, rdb) // service.NewUserService(userRepo)
	// userHandler := grcp.NewUserHandler(userService)

	userHandler := grpcHandler.NewUserHandler(userService)

	//run the grpc server\
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor()))
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Server running on port :50051")
		grpcServer.Serve(lis)
	}()

	//run grpc gateway
	//mux := runtime.NewServeMux()
	conn, err := grpc.NewClient("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}

	gwmux := runtime.NewServeMux()
	if err = pb.RegisterUserServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("failed to register gateway:", err)
	}

	ginServer := gin.Default()
	ginServer.Group("/v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux))

	log.Println("Server running on port :8080")
	ginServer.Run("0.0.0.0:8080")

}
