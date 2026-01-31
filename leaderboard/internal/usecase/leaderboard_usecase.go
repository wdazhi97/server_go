package usecase

import (
	"context"
	"errors"
	"time"

	"snake-game/leaderboard/domain/entity"
	"snake-game/leaderboard/domain/repository"
)

type LeaderboardUsecase struct {
	repo repository.LeaderboardRepository
}

func NewLeaderboardUsecase(repo repository.LeaderboardRepository) *LeaderboardUsecase {
	return &LeaderboardUsecase{
		repo: repo,
	}
}

func (uc *LeaderboardUsecase) GetLeaderboard(ctx context.Context, limit, offset int32) ([]*entity.LeaderboardEntry, int32, error) {
	entries, err := uc.repo.GetTopEntries(ctx, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to get leaderboard")
	}

	total, err := uc.repo.GetTotalUsers(ctx)
	if err != nil {
		return nil, 0, errors.New("failed to get total users")
	}

	return entries, total, nil
}

func (uc *LeaderboardUsecase) UpdateScore(ctx context.Context, userID string, score int32, gameWon bool) error {
	entry, err := uc.repo.GetEntry(ctx, userID)
	if err != nil {
		// 如果用户不存在，创建新的排行榜条目
		newEntry := &entity.LeaderboardEntry{
			UserID:      userID,
			Score:       int(score),
			GamesWon:    0,
			GamesPlayed: 1,
			UpdatedAt:   time.Now(),
		}
		if gameWon {
			newEntry.GamesWon = 1
		}
		return uc.repo.CreateEntry(ctx, newEntry)
	}

	// 更新现有条目
	if score > int32(entry.Score) {
		entry.Score = int(score)
	}
	entry.GamesPlayed++
	if gameWon {
		entry.GamesWon++
	}
	entry.UpdatedAt = time.Now()

	return uc.repo.UpdateEntry(ctx, entry)
}

func (uc *LeaderboardUsecase) GetUserRank(ctx context.Context, userID string) (int32, int32, error) {
	rank, err := uc.repo.GetUserRank(ctx, userID)
	if err != nil {
		return 0, 0, errors.New("failed to get user rank")
	}

	total, err := uc.repo.GetTotalUsers(ctx)
	if err != nil {
		return 0, 0, errors.New("failed to get total users")
	}

	return rank, total, nil
}