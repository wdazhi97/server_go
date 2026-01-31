package repository

import (
	"context"
	"sync"

	"snake-game/game/domain/entity"
)

type gameMemoryRepository struct {
	games map[string]*entity.GameState
	mutex sync.RWMutex
}

func NewGameMemoryRepository() *gameMemoryRepository {
	return &gameMemoryRepository{
		games: make(map[string]*entity.GameState),
	}
}

func (r *gameMemoryRepository) CreateGame(ctx context.Context, gameState *entity.GameState) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.games[gameState.RoomID] = gameState
	return nil
}

func (r *gameMemoryRepository) GetGame(ctx context.Context, roomID string) (*entity.GameState, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	gameState, exists := r.games[roomID]
	if !exists {
		return nil, nil
	}
	
	return gameState, nil
}

func (r *gameMemoryRepository) UpdateGame(ctx context.Context, gameState *entity.GameState) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.games[gameState.RoomID] = gameState
	return nil
}

func (r *gameMemoryRepository) DeleteGame(ctx context.Context, roomID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.games, roomID)
	return nil
}