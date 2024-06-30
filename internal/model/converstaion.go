package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Conversation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SenderId   primitive.ObjectID `bson:"senderId" json: "senderId"`
	ReceiverId primitive.ObjectID `bson:"receiverId" json: "receiverId"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
