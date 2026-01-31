package entity

import (
	"sync"
	"time"
)

type Room struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	CreatorID   string     `json:"creator_id"`
	Players     []string   `json:"players"`  // 玩家ID列表
	MaxPlayers  int        `json:"max_players"`
	Status      string     `json:"status"`   // waiting, playing, finished
	CreatedAt   time.Time  `json:"created_at"`
	Messages    []*Message `json:"messages"`
	mutex       sync.RWMutex
}

type Message struct {
	ID           string    `json:"id"`
	RoomID       string    `json:"room_id"`
	SenderID     string    `json:"sender_id"`
	SenderName   string    `json:"sender_name"`
	Content      string    `json:"content"`
	Type         string    `json:"type"`     // text, system, etc.
	CreatedAt    time.Time `json:"created_at"`
}