package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	//setup service
	var mockUserDBInSlice []User
	userRepo := NewUserRepository(mockUserDBInSlice)
	userService := NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	//setup router
	SetupRouter(r, userHandler)

	//setup service
	log.Println("running server on port 8080")
	r.Run(":8080")
}

// langkah pertama membuat struct user

// repository user
// IUserRepository mendefinisikan interface untuk repository user
type IUserRepository interface {
	GetAllUsers() []User
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan slice untuk menyimpan data user
type userRepository struct {
	db     []User //slice untuk menyimpan data user
	nextID int    //id selanjutnya yg akan digunakan untuk user baru
}

// NewUserRepository adalah factory function untuk membuat instance dari userRepository
func NewUserRepository(db []User) IUserRepository {
	return &userRepository{
		db:     db,
		nextID: 1,
	}
}

// GetAllUsers adalah method untuk mendapatkan semua data user
func (r *userRepository) GetAllUsers() []User {
	return r.db
}

// langkah kedua membuat interface dan implementasi repository user

// service user
type IUserService interface {
	GetAllUsers() []User
}

type userService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUsers() []User {
	return s.userRepo.GetAllUsers()
}

// langkah ketiga membuat interface dan implementasi service user
// handler user

type IUserHandler interface {
	GetAllUsers(c *gin.Context)
}

type UserHandler struct {
	userService IUserService
}

func NewUserHandler(userService IUserService) IUserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users := h.userService.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

// langkah keempat membuat interface dan implementasi handler user
//router user

func SetupRouter(r *gin.Engine, userHandler IUserHandler) {
	userPublicEndpoint := r.Group("/users")
	userPublicEndpoint.GET("/", userHandler.GetAllUsers)
}
