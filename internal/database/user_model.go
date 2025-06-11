package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserMetadata struct {
	ModelID string `json:"modelId" bson:"modelId"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Metadata  UserMetadata       `json:"metadata" bson:"metadata"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
