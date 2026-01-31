package repository

import (
	"context"
	"snake-game/friends/domain/entity"
)

type FriendRepository interface {
	CreateFriendship(ctx context.Context, friendship *entity.Friendship) error
	GetFriendship(ctx context.Context, userID, friendID string) (*entity.Friendship, error)
	UpdateFriendshipStatus(ctx context.Context, userID, friendID, status string) error
	GetFriends(ctx context.Context, userID string) ([]*entity.Friendship, error)
	GetPendingRequests(ctx context.Context, userID string) ([]*entity.Friendship, error)
	DeleteFriendship(ctx context.Context, userID, friendID string) error
}