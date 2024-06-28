// middleware/error_handler.go
package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go_mongoDb/internal/utils"
)

func ErrorHandlerMiddleware(c *gin.Context) {
	c.Next() // Ensure the rest of the middleware pipeline and handlers are executed

	serviceName := os.Getenv("SERVICE_NAME")

	// Check if there's an error to handle
	err := c.Errors.Last()
	if err != nil {
		// Default error response
		statusCode := http.StatusInternalServerError
		msg := "Internal Server Error"
		errorCode := "INTERNAL_SERVER_ERROR"

		// Determine the specific error type
		switch err.Type {
		case gin.ErrorTypeBind:
			statusCode = http.StatusBadRequest
			msg = "Invalid request data"
			errorCode = "BAD_REQUEST"
		case gin.ErrorTypeRender:
			statusCode = http.StatusInternalServerError
			msg = "Error rendering response"
			errorCode = "INTERNAL_SERVER_ERROR"
		default:
			// Handle custom errors
			if customErr, ok := err.Err.(*utils.CustomError); ok {
				statusCode = customErr.HTTPStatusCode
				msg = customErr.Message
				errorCode = "CUSTOM_ERROR" // You may adjust this logic based on your needs
			} else if mongoErr, ok := err.Err.(*mongo.CommandError); ok {
				switch mongoErr.Code {
				case 11000:
					statusCode = http.StatusConflict
					msg = "Duplicate value entered"
					errorCode = "DUPLICATE_ENTRY"
					// Add more cases as needed for specific MongoDB errors
				}
			}
		}

		// Send JSON response
		c.JSON(statusCode, gin.H{
			"success":        false,
			"message":        msg,
			"httpStatusCode": c.Writer.Status(),
			"error":          errorCode,
			"service":        serviceName,
		})
	}
}
