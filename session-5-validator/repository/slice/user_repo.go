package slice

import (
	"time"
	"training-golang/session-5-validator/entity"
)

// langkah kedua membuat interface dan implementasi repository user
// repository user
// IUserRepository mendefinisikan interface untuk repository user
type IUserRepository interface {
	CreateUser(user *entity.User) entity.User
	GetUserByID(id int) (entity.User, bool)
	UpdateUser(id int, user entity.User) (entity.User, bool)
	DeleteUser(id int) bool
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

// create user
func (r *userRepository) CreateUser(user *entity.User) entity.User {
	user.ID = r.nextID //mencari user id selanjutnya
	r.nextID++         //jika user baru dibuat, maka nextID akan bertambah 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.db = append(r.db, *user)
	return *user
}

// get user by id
func (r *userRepository) GetUserByID(id int) (entity.User, bool) {
	for _, user := range r.db {
		if user.ID == id {
			return user, true
		}
	}
	return entity.User{}, false
}

// UpdateUser adalah method untuk mengupdate data user berdasarkan id
func (r *userRepository) UpdateUser(id int, user entity.User) (entity.User, bool) {
	for i, u := range r.db {
		if u.ID == id {
			user.ID = u.ID
			user.CreatedAt = u.CreatedAt
			user.UpdatedAt = time.Now()
			r.db[i] = user
			return user, true
		}
	}
	return entity.User{}, false
}

// DeleteUser adalah method untuk menghapus data user berdasarkan id
func (r *userRepository) DeleteUser(id int) bool {
	for i, user := range r.db {
		if user.ID == id {
			r.db = append(r.db[:i], r.db[i+1:]...) //jika pakai database r.db[i+1 tidak akan terpakai]
			return true
		}
	}
	return false
}

// GetAllUsers adalah method untuk mendapatkan semua data user
func (r *userRepository) GetAllUsers() []entity.User {
	return r.db
}
