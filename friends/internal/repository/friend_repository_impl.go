package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"snake-game/friends/domain/entity"
	"snake-game/mongodb"
)

type friendRepositoryImpl struct {
	collection *mongo.Collection
}

func NewFriendRepository() *friendRepositoryImpl {
	return &friendRepositoryImpl{
		collection: mongodb.DB.Collection(mongodb.FriendCollection),
	}
}

func (r *friendRepositoryImpl) CreateFriendship(ctx context.Context, friendship *entity.Friendship) error {
	_, err := r.collection.InsertOne(ctx, friendship)
	return err
}

func (r *friendRepositoryImpl) GetFriendship(ctx context.Context, userID, friendID string) (*entity.Friendship, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	friendObjectID, err := primitive.ObjectIDFromHex(friendID)
	if err != nil {
		return nil, err
	}

	var friendship entity.Friendship
	err = r.collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"user_id": userObjectID, "friend_id": friendObjectID},
			{"user_id": friendObjectID, "friend_id": userObjectID},
		},
	}).Decode(&friendship)
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

func (r *friendRepositoryImpl) UpdateFriendshipStatus(ctx context.Context, userID, friendID, status string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	friendObjectID, err := primitive.ObjectIDFromHex(friendID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{
		"$or": []bson.M{
			{"user_id": userObjectID, "friend_id": friendObjectID},
			{"user_id": friendObjectID, "friend_id": userObjectID},
		},
	}, update)
	return err
}

func (r *friendRepositoryImpl) GetFriends(ctx context.Context, userID string) ([]*entity.Friendship, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"user_id": userObjectID, "status": "accepted"},
			{"friend_id": userObjectID, "status": "accepted"},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendships []*entity.Friendship
	if err = cursor.All(ctx, &friendships); err != nil {
		return nil, err
	}

	return friendships, nil
}

func (r *friendRepositoryImpl) GetPendingRequests(ctx context.Context, userID string) ([]*entity.Friendship, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{
		"friend_id": userObjectID,
		"status":    "pending",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friendships []*entity.Friendship
	if err = cursor.All(ctx, &friendships); err != nil {
		return nil, err
	}

	return friendships, nil
}

func (r *friendRepositoryImpl) DeleteFriendship(ctx context.Context, userID, friendID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	friendObjectID, err := primitive.ObjectIDFromHex(friendID)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteMany(ctx, bson.M{
		"$or": []bson.M{
			{"user_id": userObjectID, "friend_id": friendObjectID},
			{"user_id": friendObjectID, "friend_id": userObjectID},
		},
	})
	return err
}