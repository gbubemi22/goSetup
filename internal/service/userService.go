package service

import (
	"context"
	"errors"
	"time"
	
	"go_mongoDb/internal/database"
	"go_mongoDb/internal/utils"
	"go_mongoDb/internal/model"
	"go_mongoDb/internal/common"


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
	    return nil, errors.New("user with given email or username already exists")
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
	// Set createdAt and updatedAt
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
 
	_, err = s.collection.InsertOne(ctx, user)
	if err != nil {
	    return nil, err
	}
 
	return &user, nil	
}