package utils

import (
	"context"
	"fmt"
	"log"
	

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"mime/multipart"
)

// UploadImage uploads a file to an AWS S3 bucket.
func UploadImage(file *multipart.FileHeader) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Load the default AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error loading AWS config: %v", err)
		return "", fmt.Errorf("error loading AWS config: %w", err)
	}

	// Create a new S3 client
	client := s3.NewFromConfig(cfg)

	// Open the file
	f, err := file.Open()
	if err != nil {
		log.Printf("error opening file: %v", err)
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Create an uploader with the S3 client
	uploader := manager.NewUploader(client)

	// Upload the file to S3
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("gbubemi"),
		Key:    aws.String(file.Filename),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		log.Printf("error uploading file to S3: %v", err)
		return "", fmt.Errorf("error uploading file to S3: %w", err)
	}

	// Return the URL of the uploaded file
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "my-bucket", file.Filename)
	return url, nil
}
