package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go_mongoDb/internal/common"
	"go_mongoDb/internal/database"
	"go_mongoDb/internal/model"
	"go_mongoDb/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService() *UserService {
	client, err := database.DBinstance()
	if err != nil {
		panic(err)
	}
	return &UserService{
		collection: client.Database("Gomongodb").Collection("user"),
	}
}

func (s *UserService) validateUserInput(user model.User) error {
	if user.Email == "" || user.Username == "" || user.Password == "" {
		return errors.New("email, username, and password are required")
	}
	return nil
}

func (s *UserService) Create(user model.User) (*model.User, error) {
	if err := s.validateUserInput(user); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user email or username already exists
	filter := bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"username": user.Username},
		},
	}

	var existingUser model.User
	err := s.collection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return nil, utils.NewConflictError("user with given email or username already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, errors.New("internal server error")
	}

	// Validate password
	passwordValidation, err := common.ValidatePasswordString(user.Password)
	if err != nil {
		return nil, err
	}

	if !passwordValidation.IsValid {
		return nil, errors.New("password is not valid")
	}

	hashedPassword := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	user.VerifiedEmail = false

	otpToken, err := utils.GenerateRandomNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP token: %w", err)
	}
	user.OtpToken = otpToken
	user.ExpiredAt = utils.GetOtpExpiryTime()
	// Set createdAt and updatedAt
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	fmt.Println(otpToken)

	to := []string{user.Email}
	subject := "Test Email"
	body := fmt.Sprintf("<h1>Hello from Mailtrap! Here is your OTP: %s</h1>", otpToken)

	if err := utils.SendMail(subject, body, to); err != nil {
		log.Fatalf("Could not send email: %v", err)
	}

	fmt.Println("Email sent successfully!")

	_, err = s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// func (s *UserService) validateVerifyEmailInput(user model.User) error {
// 	if user.Email == "" && user.OtpToken == "" {
// 		return errors.New("email and Otp is required")
// 	}
// 	return nil
// }

func (s *UserService) VerifyEmail(email string, otpToken string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user by email
	var user model.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", err)
	}

	// Check if OTP token matches and is not expired
	if user.OtpToken != otpToken {
		return errors.New("invalid OTP token")
	}
	if user.ExpiredAt.Before(time.Now()) {
		return errors.New("OTP token has expired")
	}

	update := bson.M{
		"$set": bson.M{
			"verifiedEmail": true,
			"otpToken":      nil,
			"expiredAt":     nil,
		},
	}
	_, err = s.collection.UpdateOne(ctx, bson.M{"email": email}, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserService) SendMail(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user with the given email exists
	filter := bson.M{"email": email}
	var existingUser model.User
	err := s.collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("user does not exist")
		}
		return "", fmt.Errorf("error finding user: %w", err)
	}

	// Generate OTP token
	otpToken, err := utils.GenerateRandomNumber()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP token: %w", err)
	}

	// Update user's OTP token and expiry time
	update := bson.M{
		"$set": bson.M{
			"otpToken":  otpToken,
			"expiredAt": utils.GetOtpExpiryTime(),
		},
	}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	// Send email with OTP token
	to := []string{email}
	subject := "OTP for account verification"
	body := fmt.Sprintf("Your OTP is: %s", otpToken)

	if err := utils.SendMail(subject, body, to); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s\n", email)

	// Return success message
	successMessage := fmt.Sprintf("Email sent successfully to %s", email)
	return successMessage, nil
}
func (s *UserService) Login(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("user not found")
		}
		return "", errors.New("internal server error")
	}

	// Check if password is correct
	if !utils.VerifyPassword(password, user.Password) {
		return "", errors.New("invalid password")
	}

	if !user.VerifiedEmail {
		return "", errors.New("please verify your email")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}







func (s *UserService) UploadImage(userId string, imageURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := s.collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"image": imageURL,
		},
	}
	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": userId}, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}


// func (s *UserService) UploadImage(userId string, image string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var user model.User
// 	err := s.collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return errors.New("user not found")
// 		}
// 		return fmt.Errorf("error finding user: %w", err)
// 	}

// 	imageURL, err := utils.UploadImage(image)
// 	if err != nil {
// 		return fmt.Errorf("failed to upload image: %w", err)
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"image": imageURL,
// 		},
// 	}
// 	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": userId}, update)
// 	if err != nil {
// 		return fmt.Errorf("failed to update user: %w", err)
// 	}

// 	return nil
// }
