package middleware

import (
	"context"
	"encoding/base64"
	"session-16-crud-user-docker-compose/config"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		//check for public method that do not require authentication
		publicMethods := []string{
			"/proto.user_service.v1.UserService/GetUsers",
			"/proto.user_service.v1.UserService/GetUserByID",
		}
		for _, method := range publicMethods {
			if method == info.FullMethod {
				return handler(ctx, req)
			}
		}
		//extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		//get authorization token header
		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		//check for basic auth
		if !strings.HasPrefix(authHeader[0], "Basic ") {
			return nil, status.Errorf(codes.Unauthenticated, "Invailid authorization schema")
		}

		//decode base64 auth header
		decoded, err := base64.StdEncoding.DecodeString(authHeader[0][6:])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization token")
		}

		//split the credentials into username and password
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization token")
		}
		username, password := creds[0], creds[1]

		//validate the credentials
		if username != config.AuthBasicUsername || password != config.AuthBasicpassword {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
		}

		return handler(ctx, req)
	}

}
