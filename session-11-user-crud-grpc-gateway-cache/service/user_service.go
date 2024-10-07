package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"training-golang/session-11-user-crud-grpc-gateway-cache/entity"
	posgresgorm_raw "training-golang/session-11-user-crud-grpc-gateway-cache/repository/posgres_gom_raw"

	"github.com/redis/go-redis/v9"
)

// service user
type IUserService interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

type userService struct {
	userRepo posgresgorm_raw.IUserRepository
	rdb      *redis.Client
}

const redisUserByIDkey = "user:%d"

func NewUserService(userRepo posgresgorm_raw.IUserRepository, rdb *redis.Client) IUserService {
	return &userService{userRepo: userRepo, rdb: rdb}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("error creating user: %v", err)
	}
	createdUserJson, err := json.Marshal(createdUser)
	if err != nil {
		log.Println("failed to marshall user to JSON", err)
		return createdUser, err
	}

	//set the redis key with a 1-minute expiration

	redisKey := fmt.Sprintf(redisUserByIDkey, createdUser.ID)

	if err = s.rdb.Set(ctx, redisKey, createdUserJson, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set redis key with expitation", err)
		return createdUser, err
	}

	//invalidate the all_users cache to ensure it will be refreshed on the next call
	if err = s.rdb.Del(ctx, "all_users").Err(); err != nil {
		log.Println("failed to validate all_users cache", err)
	}

	return createdUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	redisKey := fmt.Sprintf(redisUserByIDkey, id)

	//try to get user from redis cache
	val, err := s.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		//unmarshal the user from redis cache
		if err = json.Unmarshal([]byte(val), &user); err == nil {
			return user, nil
		}
		log.Println("failed to unmarshal user from redis cache")
	}

	//if cache miss or unmarshal failed, get user from the database
	user, err = s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("user not found: %v", err)
	}

	//cache the user in redis
	userJSON, _ := json.Marshal(user)
	if err = s.rdb.Set(ctx, redisKey, userJSON, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set user data to redis:", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	updatedUser, err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed update user: %v", err)
	}

	redisKey := fmt.Sprintf(redisUserByIDkey, updatedUser.ID)
	updatedUserJSON, _ := json.Marshal(updatedUser)
	if err = s.rdb.Set(ctx, redisKey, updatedUserJSON, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set user data to redis:", err)
	}

	//invalidate the all_users cache to ensure it will be refreshed on the next call
	if err = s.rdb.Del(ctx, "all_users").Err(); err != nil {
		log.Println("failed to validate all_users cache", err)
	}
	return updatedUser, nil
}
func (s *userService) DeleteUser(ctx context.Context, id int) error {
	redisKey := fmt.Sprintf(redisUserByIDkey, id)

	//delete user cache from redis
	if err := s.rdb.Del(ctx, redisKey).Err(); err != nil {
		log.Println("failed to delete user cache:", err)
	}

	//invalidate the all_users cache to ensure it will be refreshed on the next call
	if err := s.rdb.Del(ctx, "all_users").Err(); err != nil {
		log.Println("failed to validate all_users cache", err)
	}

	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed delete user: %v", err)
	}
	return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	redisKey := "all_users"
	var users []entity.User

	//try to get all users from redis cache
	val, err := s.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		//unmarshal the users from redis cache
		if err = json.Unmarshal([]byte(val), &users); err == nil {
			return users, nil
		}
		log.Println("failed to unmarshal users from redis cache")
	}

	//if cache miss, get all users from the database
	users, err = s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get all users: %v", err)
	}
	//cache the users in redis
	userJSON, _ := json.Marshal(users)
	if err = s.rdb.Set(ctx, redisKey, userJSON, 1*time.Minute).Err(); err != nil {
		log.Println("failed to set user data to redis:", err)
	}
	return users, nil
}
