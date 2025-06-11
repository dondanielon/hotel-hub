package repositories

import (
	"ais-summoner/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TerrainRepository struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewTerrainRepository(db *mongo.Database) *TerrainRepository {
	return &TerrainRepository{
		collection: db.Collection("terrains"),
		logger:     log.New(log.Writer(), "[TerrainRepository] ", log.LstdFlags),
	}
}

func (tr *TerrainRepository) Insert(ctx context.Context, terrain *models.Terrain) (*models.Terrain, error) {
	terrain.CreatedAt = time.Now()
	terrain.UpdatedAt = time.Now()

	result, err := tr.collection.InsertOne(ctx, terrain)
	if err != nil {
		tr.logger.Printf("Error inserting terrain: %v", err)
		return nil, err
	}

	terrain.ID = result.InsertedID.(primitive.ObjectID)
	return terrain, nil
}

func (tr *TerrainRepository) GetByID(ctx context.Context, id string) (*models.Terrain, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var terrain models.Terrain
	err = tr.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&terrain)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		tr.logger.Printf("Error finding terrain by ID: %v", err)
		return nil, err
	}

	return &terrain, nil
}

func (tr *TerrainRepository) Find(ctx context.Context) ([]*models.Terrain, error) {
	cursor, err := tr.collection.Find(ctx, bson.M{})
	if err != nil {
		tr.logger.Printf("Error finding terrains: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var terrains []*models.Terrain
	for cursor.Next(ctx) {
		var terrain models.Terrain
		if err := cursor.Decode(&terrain); err != nil {
			tr.logger.Printf("Error decoding terrain: %v", err)
			continue
		}
		terrains = append(terrains, &terrain)
	}

	if err := cursor.Err(); err != nil {
		tr.logger.Printf("Cursor error: %v", err)
		return nil, err
	}

	return terrains, nil
}

func (tr *TerrainRepository) Update(ctx context.Context, id string, terrain *models.Terrain) (*models.Terrain, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	terrain.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":      terrain.Name,
			"rotation":  terrain.Rotation,
			"points":    terrain.Points,
			"updatedAt": terrain.UpdatedAt,
		},
	}

	_, err = tr.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		tr.logger.Printf("Error updating terrain: %v", err)
		return nil, err
	}

	return tr.GetByID(ctx, id)
}

func (tr *TerrainRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = tr.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		tr.logger.Printf("Error deleting terrain: %v", err)
		return err
	}

	return nil
}
