package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
			// Handle other types of errors here
			// For example, MongoDB errors
			if mongoErr, ok := err.Err.(*mongo.CommandError); ok {
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
