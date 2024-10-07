package main

import (
	"log"
	"net"

	grpcHandler "training-golang/assignment-2-simple-ewallet-system/user-service/handler/grpc"
	pb "training-golang/assignment-2-simple-ewallet-system/user-service/proto/user_service/v1"
	"training-golang/assignment-2-simple-ewallet-system/user-service/repository/postgres_gorm"
	"training-golang/assignment-2-simple-ewallet-system/user-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	listen, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "assignment2_user.", // schema name
			SingularTable: false,
		}})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	repo := postgres_gorm.NewUserRepository(db) // Initialize your repository implementation
	userService := service.NewUserService(repo)
	userHandler := grpcHandler.NewUserHandler(userService)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Printf("gRPC server started at %s", listen.Addr().String())
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
