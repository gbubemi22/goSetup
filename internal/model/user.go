package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username      string             `bson:"username" json:"username"`
	Password      string             `bson:"password" json:"password"`
	Email         string             `bson:"email" json:"email"`
	VerifiedEmail bool               `bson:"verifiedEmail" json:"verifiedEmail"`
	OtpToken      string             `bson:"otpToken,omitempty" json:"otpToken,omitempty"`
	ExpiredAt     time.Time          `bson:"expiredAt,omitempty" json:"expiredAt,omitempty"`
	Image         string             `bson:"image" json:"image"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
