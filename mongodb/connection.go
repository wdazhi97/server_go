package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

const (
	DatabaseName      = "snake_game_db"
	UserCollection    = "users"
	FriendCollection  = "friends"
	RoomCollection    = "rooms"
	GameRecordCollection = "game_records"
	LeaderboardCollection = "leaderboards"
	MessageCollection = "messages"
	GameStateCollection = "game_states"
)

// Connect 连接到 MongoDB
func Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	Client = client
	DB = client.Database(DatabaseName)

	log.Println("Connected to MongoDB successfully")
	return nil
}

// Disconnect 断开 MongoDB 连接
func Disconnect() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Disconnect(ctx)
}