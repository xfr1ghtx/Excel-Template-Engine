package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Act represents the main act document
type Act struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BigAct    *BigAct            `json:"bigAct,omitempty" bson:"bigAct,omitempty"`
	Positions []Position         `json:"positions,omitempty" bson:"positions,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
