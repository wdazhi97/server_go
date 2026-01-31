package repository

import (
	"context"
	"snake-game/matching/domain/entity"
)

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, player *entity.Player) error
	UpdatePlayerStatus(ctx context.Context, id string, status string) error
	GetPlayer(ctx context.Context, id string) (*entity.Player, error)
	DeletePlayer(ctx context.Context, id string) error
}