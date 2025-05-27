package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                      primitive.ObjectID   `bson:"_id,omitempty"`
	Name                    string               `bson:"name" json:"name"`
	Email                   string               `bson:"email" json:"email"`
	Password                string               `bson:"password" json:"password"`
	CreatedAt               time.Time            `bson:"createdAt" json:"createdAt"`
	Favorites               []primitive.ObjectID `bson:"favorites" json:"favorites"`
	RecommendationsReceived []Recommendation     `bson:"recommendationsReceived" json:"recommendationsReceived"`
}
