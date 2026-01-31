package entity

import (
	"sync"
	"time"
)

type GameState struct {
	ID           string                 `json:"id"`
	RoomID       string                 `json:"room_id"`
	Snakes       map[string]*GameSnake  `json:"snakes"`  // 玩家ID -> 蛇对象
	Foods        []Position             `json:"foods"`
	Walls        []Position             `json:"walls"`
	Status       string                 `json:"status"`  // waiting, playing, paused, finished
	UpdatedAt    time.Time              `json:"updated_at"`
	mutex        sync.RWMutex           // 内部同步锁
}

type GameSnake struct {
	PlayerID string        `json:"player_id"`
	Segments []SnakeSegment `json:"segments"`
	Color    string        `json:"color"`
	Length   int           `json:"length"`
	Score    int           `json:"score"`
	Alive    bool          `json:"alive"`
	Direction Direction    `json:"direction"`
}

type SnakeSegment struct {
	Position Position `json:"position"`
}

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type Direction int32

const (
	Direction_NONE  Direction = 0
	Direction_UP    Direction = 1
	Direction_DOWN  Direction = 2
	Direction_LEFT  Direction = 3
	Direction_RIGHT Direction = 4
)