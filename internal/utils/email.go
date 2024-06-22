package utils


import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func SendMail(subject, body string, to []string) error {
	loadEnv()

	host := os.Getenv("MAILTRAP_HOST")
	portStr := os.Getenv("MAILTRAP_PORT")
	user := os.Getenv("MAILTRAP_USER")
	pass := os.Getenv("MAILTRAP_PASS")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	sender := os.Getenv("SENDER")

	m := gomail.NewMessage()
	m.SetHeader("From", sender )
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, user, pass)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %v", err)
	}

	return nil
}
