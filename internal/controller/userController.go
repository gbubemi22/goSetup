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

type VerifyEmailRequest struct {
	Email    string `json:"email" binding:"required"`
	OtpToken string `json:"otpToken" binding:"required"`
}

// VerifyEmailHandler handles the verification of a user's email
func (controller *UserController) VerifyEmailHandler(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and OTP token are required"})
		return
	}
	controller.userService.VerifyEmail(req.Email, req.OtpToken)

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

type SendEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func (controller *UserController) SendEmailHandler(c *gin.Context) {
	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	controller.userService.SendMail(req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})

}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (controller *UserController) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	   }
	 
	   token, err := controller.userService.Login(req.Email, req.Password)
	   if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // Use specific error message
		return
	   }
	 
	   c.JSON(http.StatusOK, gin.H{"token": token})

}


// if err != nil {
	// 	if err.Error() == "user not found" || err.Error() == "invalid password" {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	} else {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	// 	}
	// 	return
	// }