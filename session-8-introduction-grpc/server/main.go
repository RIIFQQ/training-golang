package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "training-golang/session-8-introduction-grpc/proto/helloworld/v1"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServiceServer
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *server) SayHello(ctx context.Context, in *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {
	return &pb.SayHelloResponse{
		Message: fmt.Sprintf("Hello %s", in.Name),
	}, nil
}
