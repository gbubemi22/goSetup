package service

import (
	"context"
	"errors"
	"go_mongoDb/internal/model"
	"go_mongoDb/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationService struct {
	conversationCollection *mongo.Collection
	userCollection         *mongo.Collection
}

func NewConversationService(db *mongo.Database) *ConversationService {
	return &ConversationService{
		conversationCollection: db.Collection("conversation"),
		userCollection:         db.Collection("user"),
	}
}

func (cs *ConversationService) validateUserInput(conversation model.Conversation) error {
	if conversation.SenderId == primitive.NilObjectID || conversation.ReceiverId == primitive.NilObjectID {
		return errors.New("senderId and receiverId are required")
	}
	return nil
}

func (cs *ConversationService) Create(conversation model.Conversation) (*model.Conversation, error) {
	if err := cs.validateUserInput(conversation); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if both sender and receiver exist
	userFilter := bson.M{
		"$or": []bson.M{
			{"_id": conversation.SenderId},
			{"_id": conversation.ReceiverId},
		},
	}

	count, err := cs.userCollection.CountDocuments(ctx, userFilter)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	if count < 2 {
		return nil, utils.NewBadRequestError("One or both users do not exist")
	}

	// Check if conversation exists
	conversationFilter := bson.M{
		"$and": []bson.M{
			{"senderId": conversation.SenderId},
			{"receiverId": conversation.ReceiverId},
		},
	}

	var existingConv model.Conversation
	err = cs.conversationCollection.FindOne(ctx, conversationFilter).Decode(&existingConv)
	if err == nil {
		return nil, utils.NewConflictError("Conversation already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, errors.New("internal server error")
	}

	conversation.ID = primitive.NewObjectID()
	conversation.CreatedAt = time.Now()
	conversation.UpdatedAt = time.Now()

	_, err = cs.conversationCollection.InsertOne(ctx, conversation)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}



func (cs *ConversationService) GetConversationWithUsers(ctx context.Context, convID primitive.ObjectID) (*model.Conversation, *model.User, *model.User, error) {
	var conversation model.Conversation
	err := cs.conversationCollection.FindOne(ctx, bson.M{"_id": convID}).Decode(&conversation)
	if err != nil {
		return nil, nil, nil, err
	}

	var sender model.User
	err = cs.userCollection.FindOne(ctx, bson.M{"_id": conversation.SenderId}).Decode(&sender)
	if err != nil {
		return &conversation, nil, nil, err
	}

	var receiver model.User
	err = cs.userCollection.FindOne(ctx, bson.M{"_id": conversation.ReceiverId}).Decode(&receiver)
	if err != nil {
		return &conversation, &sender, nil, err
	}

	return &conversation, &sender, &receiver, nil
}
