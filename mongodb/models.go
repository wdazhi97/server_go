package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"` // 应该存储哈希值
	Email     string             `bson:"email" json:"email"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Online    bool               `bson:"online" json:"online"`
	LastSeen  time.Time          `bson:"last_seen" json:"last_seen"`
}

// Friend 好友关系模型
type Friend struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	FriendID  primitive.ObjectID `bson:"friend_id" json:"friend_id"`
	Status    string             `bson:"status" json:"status"` // pending, accepted, blocked
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// GameRoom 游戏房间模型
type GameRoom struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	RoomName    string               `bson:"room_name" json:"room_name"`
	CreatorID   primitive.ObjectID   `bson:"creator_id" json:"creator_id"`
	Players     []primitive.ObjectID `bson:"players" json:"players"`
	MaxPlayers  int                  `bson:"max_players" json:"max_players"`
	Status      string               `bson:"status" json:"status"` // waiting, playing, finished
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	GameOptions GameOptions          `bson:"game_options" json:"game_options"`
}

// GameOptions 游戏选项
type GameOptions struct {
	FoodCount    int `bson:"food_count" json:"food_count"`
	WallEnabled  bool `bson:"wall_enabled" json:"wall_enabled"`
	Speed        int `bson:"speed" json:"speed"`
}

// GameRecord 游戏记录模型
type GameRecord struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	RoomID      primitive.ObjectID   `bson:"room_id" json:"room_id"`
	Players     []PlayerGameResult   `bson:"players" json:"players"`
	WinnerID    primitive.ObjectID   `bson:"winner_id" json:"winner_id"`
	Scores      map[string]int       `bson:"scores" json:"scores"`
	GameTime    time.Duration        `bson:"game_time" json:"game_time"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
}

// PlayerGameResult 玩家游戏结果
type PlayerGameResult struct {
	PlayerID primitive.ObjectID `bson:"player_id" json:"player_id"`
	Score    int                `bson:"score" json:"score"`
	Rank     int                `bson:"rank" json:"rank"`
}

// Leaderboard 排行榜模型
type Leaderboard struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Score     int                `bson:"score" json:"score"`
	GamesWon  int                `bson:"games_won" json:"games_won"`
	GamesPlayed int              `bson:"games_played" json:"games_played"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Message 房间消息模型
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    primitive.ObjectID `bson:"room_id" json:"room_id"`
	SenderID  primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	Content   string             `bson:"content" json:"content"`
	Type      string             `bson:"type" json:"type"` // text, system, etc.
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// GameSnake 游戏中的蛇模型
type GameSnake struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Points []Point            `bson:"points" json:"points"`
	Color  string             `bson:"color" json:"color"`
	Length int                `bson:"length" json:"length"`
	Score  int                `bson:"score" json:"score"`
}

// Point 坐标点
type Point struct {
	X int `bson:"x" json:"x"`
	Y int `bson:"y" json:"y"`
}

// GameState 游戏状态
type GameState struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    primitive.ObjectID `bson:"room_id" json:"room_id"`
	Snakes    []GameSnake        `bson:"snakes" json:"snakes"`
	Foods     []Point            `bson:"foods" json:"foods"`
	Walls     []Point            `bson:"walls" json:"walls"`
	Status    string             `bson:"status" json:"status"` // waiting, playing, paused, finished
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}