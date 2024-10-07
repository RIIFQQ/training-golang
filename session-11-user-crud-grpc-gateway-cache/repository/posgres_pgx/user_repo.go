package posgrespgx

import (
	"context"
	"log"
	"training-golang/session-11-user-crud-grpc-gateway-cache/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxPoolIface interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Ping(ctx context.Context) error
}

// langkah kedua membuat interface dan implementasi repository user
// repository user
// IUserRepository mendefinisikan interface untuk repository user
type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan slice untuk menyimpan data user
type userRepository struct {
	db PgxPoolIface //untuk integrasi ke database
}

// NewUserRepository adalah factory function untuk membuat instance dari userRepository
func NewUserRepository(db PgxPoolIface) IUserRepository {
	return &userRepository{db: db}
}

// create user
func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	query := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, NOW(),NOW()) RETURNING id"
	var id int
	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return entity.User{}, err
	}
	user.ID = id
	return *user, nil

}

// get user by id
func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		return entity.User{}, err
	}
	return user, nil
}

// UpdateUser adalah method untuk mengupdate data user berdasarkan id
func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	query := "UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3"
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, id)
	if err != nil {
		log.Printf("Error updating user: %v\n", err)
		return entity.User{}, err
	}
	user.ID = id
	return user, nil
}

// DeleteUser adalah method untuk menghapus data user berdasarkan id
func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil

}

// GetAllUsers adalah method untuk mendapatkan semua data user
func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	query := "SELECT id, name, email, created_at, updated_at FROM users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Printf("Error getting all users: %v\n", err)
		return nil, err
	}
	defer rows.Close() // menutup rows setelah selesai proses
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Printf("Error scan user: %v\n", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}
