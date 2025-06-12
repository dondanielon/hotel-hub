package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Terrain struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Rotation  float64            `json:"rotation" bson:"rotation"`
	Points    []Vector2          `json:"points" bson:"points"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
