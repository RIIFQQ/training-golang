package main

import (
	"context"
	"log"
	"time"

	userPb "training-golang/ewallet-system/user-service/proto/user_service/v1"
	walletPb "training-golang/ewallet-system/wallet-service/proto/wallet_service/v1"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := userPb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts); err != nil {
		log.Fatalf("did not connect user service grpc: %v", err)
	}

	if err := walletPb.RegisterWalletServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("did not connect wallet service grpc: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Any("*any", gin.WrapH(mux))

	log.Println("gateway running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
