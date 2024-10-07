package main

import (
	"log"
	"net"
	grpcHandler "training-golang/session-9-crud-user-grpc/handler/grpc"
	pb "training-golang/session-9-crud-user-grpc/proto/user_service/v1"
	posgresgorm_raw "training-golang/session-9-crud-user-grpc/repository/posgres_gom_raw"
	"training-golang/session-9-crud-user-grpc/service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Server running on port :8080")
	grpcServer.Serve(lis)
}

// 	//setup router
// 	grcpServer := grpc.NewServer()
// 	pb.RegisterUserServiceServer(grcpServer, userHandler)

// 	lis, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	log.Println("Server running on port :8080")
// 	grcpServer.Serve(lis)

// }
