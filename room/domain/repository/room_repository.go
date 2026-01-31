package repository

import (
	"context"
	"snake-game/room/domain/entity"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *entity.Room) error
	GetRoom(ctx context.Context, roomID string) (*entity.Room, error)
	UpdateRoom(ctx context.Context, room *entity.Room) error
	DeleteRoom(ctx context.Context, roomID string) error
	AddPlayer(ctx context.Context, roomID, playerID string) error
	RemovePlayer(ctx context.Context, roomID, playerID string) error
	AddMessage(ctx context.Context, roomID string, message *entity.Message) error
	GetMessages(ctx context.Context, roomID string, limit int32) ([]*entity.Message, error)
}