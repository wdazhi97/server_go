package repository

import (
	"context"
	"snake-game/lobby/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (string, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	UpdateOnlineStatus(ctx context.Context, id string, online bool) error
}