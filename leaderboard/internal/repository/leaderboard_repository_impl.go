package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"snake-game/leaderboard/domain/entity"
	"snake-game/mongodb"
)

type leaderboardRepositoryImpl struct {
	collection *mongo.Collection
}

func NewLeaderboardRepository() *leaderboardRepositoryImpl {
	return &leaderboardRepositoryImpl{
		collection: mongodb.DB.Collection(mongodb.LeaderboardCollection),
	}
}

func (r *leaderboardRepositoryImpl) CreateEntry(ctx context.Context, entry *entity.LeaderboardEntry) error {
	_, err := r.collection.InsertOne(ctx, entry)
	return err
}

func (r *leaderboardRepositoryImpl) GetEntry(ctx context.Context, userID string) (*entity.LeaderboardEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var entry entity.LeaderboardEntry
	err = r.collection.FindOne(ctx, bson.M{"user_id": objectID}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *leaderboardRepositoryImpl) UpdateEntry(ctx context.Context, entry *entity.LeaderboardEntry) error {
	objectID, err := primitive.ObjectIDFromHex(entry.UserID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"score":        entry.Score,
			"games_won":    entry.GamesWon,
			"games_played": entry.GamesPlayed,
			"updated_at":   entry.UpdatedAt,
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"user_id": objectID}, update)
	return err
}

func (r *leaderboardRepositoryImpl) GetTopEntries(ctx context.Context, limit, offset int32) ([]*entity.LeaderboardEntry, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"score": -1}) // 按分数降序排列
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*entity.LeaderboardEntry
	if err = cursor.All(ctx, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *leaderboardRepositoryImpl) GetUserRank(ctx context.Context, userID string) (int32, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}

	// 获取用户的分数
	var userEntry entity.LeaderboardEntry
	err = r.collection.FindOne(ctx, bson.M{"user_id": objectID}).Decode(&userEntry)
	if err != nil {
		return 0, err
	}

	// 计算比当前用户分数高的用户数量
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"score": bson.M{"$gt": userEntry.Score},
	})
	if err != nil {
		return 0, err
	}

	return int32(count) + 1, nil // 排名 = 比当前用户分数高的人数 + 1
}

func (r *leaderboardRepositoryImpl) GetTotalUsers(ctx context.Context) (int32, error) {
	count, err := r.collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}