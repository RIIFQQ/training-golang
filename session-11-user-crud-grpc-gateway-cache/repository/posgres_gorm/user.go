package posgresgorm

import (
	"context"
	"errors"
	"log"
	"training-golang/session-11-user-crud-grpc-gateway-cache/entity"

	"gorm.io/gorm"
)

type gormDBIface interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan slice untuk menyimpan data user
type userRepository struct {
	db gormDBIface //untuk integrasi ke database
}

// NewUserRepository adalah factory function untuk membuat instance dari userRepository
func NewUserRepository(db gormDBIface) IUserRepository {
	return &userRepository{db: db}
}

// create user
func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	if er := r.db.WithContext(ctx).Create(user).Error; er != nil {
		log.Printf("Error creating user: %v\n", er)
		return entity.User{}, er
	}
	return *user, nil
}

// get user by id
func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	if er := r.db.WithContext(ctx).Select("id, name, email, created_at, updated_at").First(&user, id).Error; er != nil {
		log.Printf("Error getting user: %v\n", er)
		return entity.User{}, er
	}
	return user, nil
}

// UpdateUser adalah method untuk mengupdate data user berdasarkan id
func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	var existingUser entity.User
	if er := r.db.WithContext(ctx).Select("id, name, email, created_at, updated_at").First(&existingUser, id).Error; er != nil {
		log.Printf("Error getting user: %v\n", er)
		return entity.User{}, er
	}

	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Password = user.Password
	if err := r.db.WithContext(ctx).Save(&existingUser).Error; err != nil {
		return entity.User{}, err
	}
	return existingUser, nil

}

// DeleteUser adalah method untuk menghapus data user berdasarkan id
func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil

}

// GetAllUsers adalah method untuk mendapatkan semua data user
func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.WithContext(ctx).Select("id, name, email, created_at, updated_at").Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, nil
		}
		log.Printf("Error getting all users: %v\n", err)
		return nil, err
	}
	return users, nil
}
