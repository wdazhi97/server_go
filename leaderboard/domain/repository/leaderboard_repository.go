package repository

import (
	"context"
	"snake-game/leaderboard/domain/entity"
)

type LeaderboardRepository interface {
	CreateEntry(ctx context.Context, entry *entity.LeaderboardEntry) error
	GetEntry(ctx context.Context, userID string) (*entity.LeaderboardEntry, error)
	UpdateEntry(ctx context.Context, entry *entity.LeaderboardEntry) error
	GetTopEntries(ctx context.Context, limit, offset int32) ([]*entity.LeaderboardEntry, error)
	GetUserRank(ctx context.Context, userID string) (int32, error)
	GetTotalUsers(ctx context.Context) (int32, error)
}