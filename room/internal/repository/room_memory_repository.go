package repository

import (
	"context"
	"sync"

	"snake-game/room/domain/entity"
)

type roomMemoryRepository struct {
	rooms map[string]*entity.Room
	mutex sync.RWMutex
}

func NewRoomMemoryRepository() *roomMemoryRepository {
	return &roomMemoryRepository{
		rooms: make(map[string]*entity.Room),
	}
}

func (r *roomMemoryRepository) CreateRoom(ctx context.Context, room *entity.Room) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.rooms[room.ID] = room
	return nil
}

func (r *roomMemoryRepository) GetRoom(ctx context.Context, roomID string) (*entity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	room, exists := r.rooms[roomID]
	if !exists {
		return nil, nil
	}
	
	return room, nil
}

func (r *roomMemoryRepository) UpdateRoom(ctx context.Context, room *entity.Room) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.rooms[room.ID] = room
	return nil
}

func (r *roomMemoryRepository) DeleteRoom(ctx context.Context, roomID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.rooms, roomID)
	return nil
}

func (r *roomMemoryRepository) AddPlayer(ctx context.Context, roomID, playerID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	room, exists := r.rooms[roomID]
	if !exists {
		return nil
	}
	
	// 检查玩家是否已经在房间中
	for _, pid := range room.Players {
		if pid == playerID {
			return nil // 玩家已在房间中
		}
	}
	
	// 检查房间是否已满
	if len(room.Players) >= room.MaxPlayers {
		return nil // 房间已满
	}
	
	room.Players = append(room.Players, playerID)
	r.rooms[roomID] = room
	
	return nil
}

func (r *roomMemoryRepository) RemovePlayer(ctx context.Context, roomID, playerID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	room, exists := r.rooms[roomID]
	if !exists {
		return nil
	}
	
	// 找到玩家并从列表中移除
	for i, pid := range room.Players {
		if pid == playerID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			break
		}
	}
	
	r.rooms[roomID] = room
	return nil
}

func (r *roomMemoryRepository) AddMessage(ctx context.Context, roomID string, message *entity.Message) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	room, exists := r.rooms[roomID]
	if !exists {
		return nil
	}
	
	room.Messages = append(room.Messages, message)
	r.rooms[roomID] = room
	
	return nil
}

func (r *roomMemoryRepository) GetMessages(ctx context.Context, roomID string, limit int32) ([]*entity.Message, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	room, exists := r.rooms[roomID]
	if !exists {
		return nil, nil
	}
	
	// 返回最新的消息
	startIdx := 0
	if len(room.Messages) > int(limit) {
		startIdx = len(room.Messages) - int(limit)
	}
	
	messages := make([]*entity.Message, len(room.Messages[startIdx:]))
	copy(messages, room.Messages[startIdx:])
	
	return messages, nil
}