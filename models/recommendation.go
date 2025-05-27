package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recommendation struct {
	PropertyID    primitive.ObjectID `bson:"propertyId" json:"propertyId"`
	RecommendedBy primitive.ObjectID `bson:"recommendedBy" json:"recommendedBy"`
	RecommendedAt time.Time          `bson:"recommendedAt" json:"recommendedAt"`
}
