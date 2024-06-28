package utils

import (
	"encoding/json"
	"net/http"
	"os"
)

type CustomError struct {
	Message        string `json:"message"`
	ErrorCode      int    `json:"errorCode,omitempty"`
	HTTPStatusCode int    `json:"httpStatusCode,omitempty"`
	Service        string `json:"service,omitempty"`
	Success        bool   `json:"success,omitempty"`
}

func (e *CustomError) Error() string {
	return e.Message
}

var serviceName = os.Getenv("SERVICE_NAME")

func NewUnauthorizedError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      401,
		HTTPStatusCode: http.StatusUnauthorized,
		Service:        serviceName,
		Success:        false,
	}
}

func NewBadRequestError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      400,
		HTTPStatusCode: http.StatusBadRequest,
		Service:        serviceName,
		Success:        false,
	}
}

func NewConflictError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      409,
		HTTPStatusCode: http.StatusConflict,
		Service:        serviceName,
		Success:        false,
	}
}

func NewInternalServerError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      500,
		HTTPStatusCode: http.StatusInternalServerError,
		Service:        serviceName,
		Success:        false,
	}
}

func NewUnauthenticatedError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      401,
		HTTPStatusCode: http.StatusUnauthorized,
		Service:        serviceName,
		Success:        false,
	}
}

func NewNotFoundError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      404,
		HTTPStatusCode: http.StatusNotFound,
		Service:        serviceName,
		Success:        false,
	}
}

func (e *CustomError) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
