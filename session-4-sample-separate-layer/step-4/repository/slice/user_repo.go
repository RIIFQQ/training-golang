package slice

import "training-golang/session-4-sample-separate-layer/step-4/entity"

// langkah kedua membuat interface dan implementasi repository user
// repository user
// IUserRepository mendefinisikan interface untuk repository user
type IUserRepository interface {
	GetAllUsers() []entity.User
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan slice untuk menyimpan data user
type userRepository struct {
	db     []entity.User //slice untuk menyimpan data user
	nextID int           //id selanjutnya yg akan digunakan untuk user baru
}

// NewUserRepository adalah factory function untuk membuat instance dari userRepository
func NewUserRepository(db []entity.User) IUserRepository {
	return &userRepository{
		db:     db,
		nextID: 1,
	}
}

// GetAllUsers adalah method untuk mendapatkan semua data user
func (r *userRepository) GetAllUsers() []entity.User {
	return r.db
}
