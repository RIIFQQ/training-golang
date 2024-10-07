package grcp

import (
	"context"
	"fmt"
	"log"
	"strings"
	"training-golang/session-11-user-crud-grpc-gateway-cache/entity"
	pb "training-golang/session-11-user-crud-grpc-gateway-cache/proto/user_service/v1"
	"training-golang/session-11-user-crud-grpc-gateway-cache/service"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	// Add this import
)

// langkah ketiga membuat interface dan implementasi service user
// handler user

type IUserHandler interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
}

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUserResponse, error) {
	users, err := h.userService.GetAllUsers(ctx) // get all users

	if err != nil {
		return nil, err
	}
	var userProto []*pb.User
	for _, user := range users {
		userProto = append(userProto, &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}
	return &pb.GetUserResponse{
		Users: userProto,
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := h.userService.GetUserByID(ctx, int(req.Id))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
	return res, nil
}

// create user
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.MutationResponse, error) {
	createdUser, err := h.userService.CreateUser(ctx, &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success create user with id %d", createdUser.ID),
	}, nil

}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.MutationResponse, error) {
	updatedUser, err := h.userService.UpdateUser(ctx, int(req.Id), entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success update user with id %d", updatedUser.ID),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.MutationResponse, error) {
	err := h.userService.DeleteUser(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success delete user with id %d ", req.Id),
	}, nil
}

// custom untuk handler validation
func fieldErrorMessage(errorMessage string) string {
	switch {
	case strings.Contains(errorMessage, "'Name' failed on the 'required' tag"):
		return "Name is mandatory"
	case strings.Contains(errorMessage, "'Name'failed on the 'min' tag"):
		return "Name must be at least 3 characters"
	case strings.Contains(errorMessage, "'Email'failed on the 'required' tag"):
		return "Email is mandatory"
	case strings.Contains(errorMessage, "'Email'failed on the 'email' tag"):
		return "Email is not valid"
	}
	return errorMessage
}
