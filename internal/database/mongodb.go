package database

import (
	"auction-web/internal/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient connects the application with mongo database and creates new mongo client
func NewMongoClient(uri string, timeout time.Duration) (client *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, logger.WrapError(err, "failed to connect MongoDB")
	}

	// Ping the database to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, logger.WrapError(err, "failed to ping MongoDB")
	}

	return client, nil
}

// DisconnectMongoClient disconnects the application with mongo database
func DisconnectMongoClient(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		logger.WrapError(err, "failed to disconnect mongo client")
	}
}
