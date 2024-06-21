package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go_mongoDb/internal/model"
	"go_mongoDb/internal/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	userService := service.NewUserService()
	return &UserController{
		userService: userService,
	}
}

func (controller *UserController) CreateUserHttp(c *gin.Context) {
	var userInput model.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	createdUser, err := controller.userService.Create(userInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}
