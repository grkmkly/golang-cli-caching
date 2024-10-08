package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkPort struct {
	Id        primitive.ObjectID `bson:"_id"`
	Link      string             `bson:"link"`
	Port      string             `bson:"port"`
	CreatedAt time.Time          `bson:"createdat"`
}
