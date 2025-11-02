package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Card struct {
	Id          bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	Section     string        `json:"section"`
}
