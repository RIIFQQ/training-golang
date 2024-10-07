package handler

import (
	"strconv"
	"time"
	"training-golang/session-2-latihan-crud-user-gin/entity"

	"github.com/gin-gonic/gin"
)

var (
	users  []entity.User
	nextID int = 1
)

// create User
func CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user.ID = nextID
	nextID++

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	users = append(users, user)

	c.JSON(201, user)

}

// get user	by id
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(200, user)
			return
		}
	}

	c.JSON(404, gin.H{"error": "User not found"})

}

//update user

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var user entity.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for i, u := range users {
		updateUser := entity.User{
			ID:        id,
			Name:      user.Name,
			Email:     user.Email,
			Password:  u.Password,
			CreatedAt: u.CreatedAt,
			UpdatedAt: time.Now(),
		}
		users[i] = updateUser
		c.JSON(200, updateUser)
		return
	}

	c.JSON(404, gin.H{"error": "User not found"})

}

// delete user
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	//users := []int{0, 1, 2, 3, 4, 5}
	// i := 2
	// users = append(users[:i], user[i+1:]...)
	// users[:i] will be [0, 1]
	// users[i+1:] will be [3, 4, 5]
	// users = append([]int{0,1}, []int{3,4,5}...)
	// users slice will be [0, 1, 3, 4, 5]
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(200, gin.H{"message": "User deleted"})
			return
		}
	}

	c.JSON(404, gin.H{"error": "User not found"})

}

// get all user
func GetAllUsers(c *gin.Context) {
	c.JSON(200, users)
}
