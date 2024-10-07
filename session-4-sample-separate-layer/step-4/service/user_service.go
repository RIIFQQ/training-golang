package service

import (
	"training-golang/session-4-sample-separate-layer/step-4/entity"
	"training-golang/session-4-sample-separate-layer/step-4/repository/slice"
)

// service user
type IUserService interface {
	GetAllUsers() []entity.User
}

type userService struct {
	userRepo slice.IUserRepository
}

func NewUserService(userRepo slice.IUserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUsers() []entity.User {
	return s.userRepo.GetAllUsers()
}
