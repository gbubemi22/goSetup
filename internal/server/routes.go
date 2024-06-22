package server

import (
	"github.com/gin-gonic/gin"
	"go_mongoDb/internal/controller"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	userController := controller.NewUserController()

	r.POST("/v1/auth/users/create", userController.CreateUserHttp)
	r.POST("/v1/auth/users/verify-email", userController.VerifyEmailHandler)

	r.POST("/v1/auth/users/send-email", userController.SendEmailHandler)

	r.POST("/v1/auth/users/login", userController.LoginHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
