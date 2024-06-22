package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)
func GenerateRandomNumber() (string, error) {
	// Define the maximum value for a 6-digit number (exclusive)
	max := big.NewInt(900000) // 999999 - 100000 + 1

	// Generate a random number between 0 and 899999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("error generating random number: %v", err)
	}

	// Add 100000 to ensure the number is 6 digits
	n = n.Add(n, big.NewInt(100000))

	return n.String(), nil
}


func GetOtpExpiryTime() time.Time {
	expiryDuration := 10 * time.Minute
	expiredAt := time.Now().Add(expiryDuration)
	return expiredAt
}