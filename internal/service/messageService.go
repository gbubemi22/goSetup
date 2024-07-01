package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go_mongoDb/internal/model"
)

type MessageService struct {
	conversationCollection *mongo.Collection
	messageCollection      *mongo.Collection
}

func NewMessageService(db *mongo.Database) *MessageService {
	return &MessageService{
		conversationCollection: db.Collection("conversation"),
		messageCollection:      db.Collection("message"),
	}
}


func (ms *MessageService) validateUserInput(message model.Message) error {
	if message.SenderId == primitive.NilObjectID || message.ConversationId == primitive.NilObjectID {
		return errors.New("ConversationId and SenderId are required")
	}
	return nil
}

func (ms *MessageService) Create(message model.Message) (*model.Message, error) {

	if err := ms.validateUserInput(message); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	_, err := ms.messageCollection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
