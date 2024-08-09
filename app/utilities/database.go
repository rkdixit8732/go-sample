package utilities

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDB *mongo.Client
)

// Initialize MongoDB client and connect
func InitMongoDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	mongoDB = client

}

// Check MongoDB connection status
func CheckMongoDBConnection() error {
	if mongoDB == nil {
		return errors.New("MongoDB client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := mongoDB.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

// Get MongoDB connection client
func GetMongoClient() *mongo.Client {
	return mongoDB
}
