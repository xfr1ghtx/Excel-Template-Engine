package repository

import (
	"context"
	"errors"

	"github.com/stepanpotapov/Excel-Template-Engine/internal/models"
	"github.com/stepanpotapov/Excel-Template-Engine/internal/utils"
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
	utils.LogMethodInit("ActRepository.Create")

	utils.LogMongoTransaction("INSERT", "Inserting new act into database")
	result, err := r.collection.InsertOne(ctx, act)
	if err != nil {
		utils.LogMethodError("ActRepository.Create", err)
		return "", err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err := errors.New("failed to convert inserted ID to ObjectID")
		utils.LogMethodError("ActRepository.Create", err)
		return "", err
	}

	utils.LogInfo("Successfully created act with ID: %s", insertedID.Hex())
	utils.LogMethodSuccess("ActRepository.Create")
	return insertedID.Hex(), nil
}

// FindByID retrieves an act by its ID
func (r *actRepository) FindByID(ctx context.Context, id string) (*models.Act, error) {
	utils.LogMethodInit("ActRepository.FindByID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.LogError("Invalid ObjectID format: %v", err)
		utils.LogMethodError("ActRepository.FindByID", err)
		return nil, errors.New("invalid ID format")
	}

	utils.LogMongoTransaction("SELECT", "Finding act by ID: "+id)
	var act models.Act
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&act)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.LogError("Act not found with ID: %s", id)
			utils.LogMethodError("ActRepository.FindByID", err)
			return nil, errors.New("act not found")
		}
		utils.LogMethodError("ActRepository.FindByID", err)
		return nil, err
	}

	utils.LogInfo("Successfully found act with ID: %s", id)
	utils.LogMethodSuccess("ActRepository.FindByID")
	return &act, nil
}

// Update updates an existing act in the database
func (r *actRepository) Update(ctx context.Context, id string, act *models.Act) error {
	utils.LogMethodInit("ActRepository.Update")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.LogError("Invalid ObjectID format: %v", err)
		utils.LogMethodError("ActRepository.Update", err)
		return errors.New("invalid ID format")
	}

	update := bson.M{
		"$set": act,
	}

	utils.LogMongoTransaction("UPDATE", "Updating act with ID: "+id)
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		utils.LogMethodError("ActRepository.Update", err)
		return err
	}

	if result.MatchedCount == 0 {
		err := errors.New("act not found")
		utils.LogError("Act not found with ID: %s", id)
		utils.LogMethodError("ActRepository.Update", err)
		return err
	}

	utils.LogInfo("Successfully updated act with ID: %s", id)
	utils.LogMethodSuccess("ActRepository.Update")
	return nil
}
