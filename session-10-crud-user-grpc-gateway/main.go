package main

import (
	"context"
	"log"
	"net"
	grpcHandler "training-golang/session-10-crud-user-grpc-gateway/handler/grpc"
	"training-golang/session-10-crud-user-grpc-gateway/middleware"
	pb "training-golang/session-10-crud-user-grpc-gateway/proto/user_service/v1"
	posgresgorm_raw "training-golang/session-10-crud-user-grpc-gateway/repository/posgres_gom_raw"
	"training-golang/session-10-crud-user-grpc-gateway/service"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	//setup service
	dsn := "postgresql://postgres:admin@localhost:5432/go_db"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		log.Fatalln(err)
	}
	//  pake postgresgorm.
	//userRepo := posgresgorm.NewUserRepository(gormDB)
	userRepo := posgresgorm_raw.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo) // service.NewUserService(userRepo)
	// userHandler := grcp.NewUserHandler(userService)

	userHandler := grpcHandler.NewUserHandler(userService)

	//run the grpc server\
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor()))
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Server running on port :50051")
		grpcServer.Serve(lis)
	}()

	//run grpc gateway
	//mux := runtime.NewServeMux()
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}

	gwmux := runtime.NewServeMux()
	if err = pb.RegisterUserServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("failed to register gateway: %v", err)
	}

	ginServer := gin.Default()
	ginServer.Group("/v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux))

	log.Println("Server running on port :8080")
	ginServer.Run(":8080")

}
