package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"snake-game/lobby/domain/entity"
	"snake-game/lobby/domain/repository"
	"snake-game/mongodb"
)

type userRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository() repository.UserRepository {
	return &userRepositoryImpl{
		collection: mongodb.DB.Collection(mongodb.UserCollection),
	}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	objectID := result.InsertedID.(primitive.ObjectID)
	return objectID.Hex(), nil
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) UpdateOnlineStatus(ctx context.Context, id string, online bool) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"online":     online,
			"last_seen":  time.Now(),
			"updated_at": time.Now(),
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}