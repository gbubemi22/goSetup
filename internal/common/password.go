package common

import (
	"errors"
	"regexp"
)

type PasswordValidationResponse struct {
	IsValid bool
}

func ValidatePasswordString(password string) (*PasswordValidationResponse, error) {
	if len(password) < 8 {
		return &PasswordValidationResponse{IsValid: false}, errors.New("password must be at least 8 characters long")
	}
	if len(password) > 20 {
		return &PasswordValidationResponse{IsValid: false}, errors.New("password must be no more than 20 characters long")
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)

	if !hasDigit {
		return &PasswordValidationResponse{IsValid: false}, errors.New("password must contain at least one digit")
	}
	if !hasLower {
		return &PasswordValidationResponse{IsValid: false}, errors.New("password must contain at least one lowercase letter")
	}
	if !hasUpper {
		return &PasswordValidationResponse{IsValid: false}, errors.New("password must contain at least one uppercase letter")
	}

	return &PasswordValidationResponse{IsValid: true}, nil
}
