package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"snake-game/matching/domain/entity"
	"snake-game/mongodb"
)

type playerRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPlayerRepository() *playerRepositoryImpl {
	return &playerRepositoryImpl{
		collection: mongodb.DB.Collection("players"), // 使用专门的players集合
	}
}

func (r *playerRepositoryImpl) CreatePlayer(ctx context.Context, player *entity.Player) error {
	_, err := r.collection.InsertOne(ctx, player)
	return err
}

func (r *playerRepositoryImpl) UpdatePlayerStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *playerRepositoryImpl) GetPlayer(ctx context.Context, id string) (*entity.Player, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var player entity.Player
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepositoryImpl) DeletePlayer(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}