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

type UserRepository struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
		logger:     log.New(log.Writer(), "[UserRepository] ", log.LstdFlags),
	}
}

func (ur *UserRepository) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.Metadata.ModelID == "" {
		user.Metadata.ModelID = "019590ed-2942-7503-b8db-0a185f81a1de"
	}

	result, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		ur.logger.Printf("Error inserting user: %v", err)
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = ur.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		ur.logger.Printf("Error finding user by ID: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := ur.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		ur.logger.Printf("Error finding user by username: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := ur.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		ur.logger.Printf("Error finding user by email: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Find(ctx context.Context) ([]*models.User, error) {
	cursor, err := ur.collection.Find(ctx, bson.M{})
	if err != nil {
		ur.logger.Printf("Error finding users: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			ur.logger.Printf("Error decoding user: %v", err)
			continue
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		ur.logger.Printf("Cursor error: %v", err)
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) Update(ctx context.Context, id string, user *models.User) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"username":  user.Username,
			"email":     user.Email,
			"metadata":  user.Metadata,
			"updatedAt": user.UpdatedAt,
		},
	}

	_, err = ur.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		ur.logger.Printf("Error updating user: %v", err)
		return nil, err
	}

	return ur.GetByID(ctx, id)
}

func (ur *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = ur.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		ur.logger.Printf("Error deleting user: %v", err)
		return err
	}

	return nil
}
