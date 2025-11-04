package repository

import (
	"context"
	"log"
	"time"

	"github.com/stepanpotapov/Excel-Template-Engine/internal/config"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/utils"
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
	utils.LogMethodInit("ConnectMongoDB")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDBTimeout)
	defer cancel()

	utils.LogMongoTransaction("CONNECT", "Attempting to connect to MongoDB")
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		utils.LogMethodError("ConnectMongoDB", err)
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	utils.LogMongoTransaction("PING", "Verifying MongoDB connection")
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.LogMethodError("ConnectMongoDB", err)
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	utils.LogInfo("Successfully connected to MongoDB at %s", cfg.MongoDBURI)
	utils.LogMethodSuccess("ConnectMongoDB")

	database := client.Database(cfg.MongoDBDatabase)

	return &MongoDBClient{
		Client:   client,
		Database: database,
		Config:   cfg,
	}
}

// Disconnect closes the MongoDB connection
func (m *MongoDBClient) Disconnect() error {
	utils.LogMethodInit("MongoDBClient.Disconnect")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	utils.LogMongoTransaction("DISCONNECT", "Closing MongoDB connection")
	err := m.Client.Disconnect(ctx)
	if err != nil {
		utils.LogMethodError("MongoDBClient.Disconnect", err)
		return err
	}

	utils.LogInfo("Disconnected from MongoDB")
	utils.LogMethodSuccess("MongoDBClient.Disconnect")
	return nil
}

// GetCollection returns a MongoDB collection
func (m *MongoDBClient) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
