package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ConversationId primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	SenderId       primitive.ObjectID `bson:"senderId" json:"senderId"`
	Message        string             `json:"message"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

