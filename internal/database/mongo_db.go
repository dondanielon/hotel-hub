package database

import (
	"ais-summoner/internal/repositories"
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client      *mongo.Client
	db          *mongo.Database
	logger      *log.Logger
	terrainRepo *repositories.TerrainRepository
	userRepo    *repositories.UserRepository
}

func NewMongoDB() *MongoDB {
	logger := log.New(log.Writer(), "[MongoDB] ", log.LstdFlags)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	logger.Println("Connected to MongoDB successfully")
	db := client.Database(os.Getenv("MONGODB_DATABASE"))

	return &MongoDB{
		client:      client,
		db:          db,
		logger:      logger,
		terrainRepo: repositories.NewTerrainRepository(db),
		userRepo:    repositories.NewUserRepository(db),
	}
}

func (m *MongoDB) UserRepository() *repositories.UserRepository {
	return m.userRepo
}

func (m *MongoDB) TerrainRepository() *repositories.TerrainRepository {
	return m.terrainRepo
}

func (m *MongoDB) Close() error {
	m.logger.Println("Disconnecting from MongoDB")
	return m.client.Disconnect(context.Background())
}
