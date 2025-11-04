package repository

import (
	"context"
	"log"
	"time"

	"github.com/stepanpotapov/Excel-Template-Engine/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient wraps the MongoDB client
type MongoDBClient struct {
	Client   *mongo.Client
	Database *mongo.Database
	Config   *config.Config
}

// ConnectMongoDB establishes a connection to MongoDB
func ConnectMongoDB(cfg *config.Config) *MongoDBClient {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDBTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Printf("Successfully connected to MongoDB at %s", cfg.MongoDBURI)

	database := client.Database(cfg.MongoDBDatabase)

	return &MongoDBClient{
		Client:   client,
		Database: database,
		Config:   cfg,
	}
}

// Disconnect closes the MongoDB connection
func (m *MongoDBClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.Client.Disconnect(ctx)
	if err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
		return err
	}

	log.Println("Disconnected from MongoDB")
	return nil
}

// GetCollection returns a MongoDB collection
func (m *MongoDBClient) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

