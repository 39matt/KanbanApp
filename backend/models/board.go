package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Board struct {
	ID        bson.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name      string          `bson:"name" json:"name"`
	CardIDs   []bson.ObjectID `bson:"cardIds,omitempty" json:"cardIds"`
	Alias     string          `bson:"alias" json:"alias"`
	CreatedAt time.Time       `bson:"createdAt" json:"createdAt"`
}
