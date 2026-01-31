package entity

import "time"

type LeaderboardEntry struct {
	UserID    string    `bson:"user_id" json:"user_id"`
	Score     int       `bson:"score" json:"score"`
	GamesWon  int       `bson:"games_won" json:"games_won"`
	GamesPlayed int     `bson:"games_played" json:"games_played"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}