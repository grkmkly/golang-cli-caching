package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Username     string
	Server       string
	Database     string
	Collection   *mongo.Collection
	Client       *mongo.Client
	Ctx          context.Context
	ClientOption *options.ClientOptions
}
