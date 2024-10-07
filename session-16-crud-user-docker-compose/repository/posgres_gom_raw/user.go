package posgresgormraw

import (
	"context"
	"errors"
	"log"
	"session-16-crud-user-docker-compose/entity"

	"github.com/redis/go-redis/v9"
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
	db    gormDBIface //untuk integrasi ke database
	redis *redis.Client
}

// NewUserRepository adalah factory function untuk membuat instance dari userRepository
func NewUserRepository(db gormDBIface) IUserRepository {
	return &userRepository{db: db}
}

// create user
func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	query := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, NOW(),NOW()) RETURNING id"
	var createdID int
	if err := r.db.WithContext(ctx).Raw(query, user.Name, user.Email, user.Password).Scan(&createdID).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		return entity.User{}, err
	}
	user.ID = createdID
	return *user, nil
}

// get user by id
func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1"
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&user).Error; err != nil {
		log.Printf("Error getting user: %v\n", err)
		return entity.User{}, err
	}

	return user, nil
}

// UpdateUser adalah method untuk mengupdate data user berdasarkan id
func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	query := "Update users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3"
	if er := r.db.WithContext(ctx).Exec(query, user.Name, user.Email, id).Error; er != nil {
		log.Printf("Error getting user: %v\n", er)
		return entity.User{}, er
	}
	return user, nil

}

// DeleteUser adalah method untuk menghapus data user berdasarkan id
func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	if err := r.db.WithContext(ctx).Exec(query, id).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil

}

// GetAllUsers adalah method untuk mendapatkan semua data user
func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	query := "SELECT id, name, email, created_at, updated_at FROM users"
	if err := r.db.WithContext(ctx).Raw(query).Scan(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, nil
		}
		log.Printf("Error getting all users: %v\n", err)
		return nil, err
	}
	return users, nil
}
