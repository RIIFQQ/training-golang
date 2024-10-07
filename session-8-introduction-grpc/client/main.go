package main

import (
	"context"
	"fmt"
	"log"
	pb "training-golang/session-8-introduction-grpc/proto/helloworld/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedGreeterServiceServer
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()
	greeterClient := pb.NewGreeterServiceClient(conn)
	resp, err := greeterClient.SayHello(context.Background(), &pb.SayHelloRequest{
		Name: "Golang",
	})

	if err != nil {
		log.Fatalf("failed to call SayHello: %v", err)
	}
	fmt.Println(resp.Message)

}
