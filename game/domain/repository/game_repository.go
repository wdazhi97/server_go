package repository

import (
	"context"
	"snake-game/game/domain/entity"
)

type GameRepository interface {
	CreateGame(ctx context.Context, gameState *entity.GameState) error
	GetGame(ctx context.Context, roomID string) (*entity.GameState, error)
	UpdateGame(ctx context.Context, gameState *entity.GameState) error
	DeleteGame(ctx context.Context, roomID string) error
}