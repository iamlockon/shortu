package db

import (
	"context"
	"fmt"
	"time"

	"github.com/iamlockon/shortu/internal/error"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// New creates a mongo client or nil if error presents
func New(config *MongoConfig) (*MongoClient, *error.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetConnStr()))
	if err != nil {
		return nil, error.New(error.InvalidConfigError, fmt.Sprintf("failed to connect mongo: %v", err))
	}
	return &MongoClient{
		client,
	}, nil
}
