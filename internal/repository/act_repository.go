package repository

import (
	"context"
	"errors"
	"log"

	"github.com/stepanpotapov/Excel-Template-Engine/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ActRepository defines the interface for act data operations
type ActRepository interface {
	Create(ctx context.Context, act *models.Act) (string, error)
	FindByID(ctx context.Context, id string) (*models.Act, error)
	Update(ctx context.Context, id string, act *models.Act) error
}

// actRepository implements ActRepository
type actRepository struct {
	collection *mongo.Collection
}

// NewActRepository creates a new ActRepository
func NewActRepository(mongoClient *MongoDBClient) ActRepository {
	collection := mongoClient.GetCollection(mongoClient.Config.MongoDBCollection)
	return &actRepository{
		collection: collection,
	}
}

// Create inserts a new act into the database
func (r *actRepository) Create(ctx context.Context, act *models.Act) (string, error) {
	result, err := r.collection.InsertOne(ctx, act)
	if err != nil {
		log.Printf("Error creating act: %v", err)
		return "", err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to convert inserted ID to ObjectID")
	}

	log.Printf("Successfully created act with ID: %s", insertedID.Hex())
	return insertedID.Hex(), nil
}

// FindByID retrieves an act by its ID
func (r *actRepository) FindByID(ctx context.Context, id string) (*models.Act, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return nil, errors.New("invalid ID format")
	}

	var act models.Act
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&act)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Act not found with ID: %s", id)
			return nil, errors.New("act not found")
		}
		log.Printf("Error finding act: %v", err)
		return nil, err
	}

	log.Printf("Successfully found act with ID: %s", id)
	return &act, nil
}

// Update updates an existing act in the database
func (r *actRepository) Update(ctx context.Context, id string, act *models.Act) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return errors.New("invalid ID format")
	}

	update := bson.M{
		"$set": act,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Printf("Error updating act: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Printf("Act not found with ID: %s", id)
		return errors.New("act not found")
	}

	log.Printf("Successfully updated act with ID: %s", id)
	return nil
}

