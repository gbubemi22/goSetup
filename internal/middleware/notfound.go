package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleNotFound(c *gin.Context) {
	// Set status code
	c.Status(http.StatusNotFound)

	// Access the environment variable
	serviceName := os.Getenv("SERVICE_NAME")

	// Define response data
	response := map[string]interface{}{
		"success":        false,
		"message":        "Route does not exist",
		"httpStatusCode": http.StatusNotFound,
		"error":          "NOT_FOUND_ERROR",
		"service":        serviceName,
	}

	// Send JSON response
	c.JSON(http.StatusNotFound, response)
}
