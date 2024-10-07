package handler

import (
	"net/http"
	"training-golang/session-4-sample-separate-layer/step-4/service"

	"github.com/gin-gonic/gin"
)

// langkah ketiga membuat interface dan implementasi service user
// handler user

type IUserHandler interface {
	GetAllUsers(c *gin.Context)
}

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) IUserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users := h.userService.GetAllUsers()
	c.JSON(http.StatusOK, users)
}
