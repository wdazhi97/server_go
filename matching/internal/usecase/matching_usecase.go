package usecase

import (
	"context"
	"errors"
	"fmt"

	"snake-game/matching/domain/entity"
	"snake-game/matching/domain/repository"
)

type MatchingUsecase struct {
	playerRepo repository.PlayerRepository
}

func NewMatchingUsecase(playerRepo repository.PlayerRepository) *MatchingUsecase {
	return &MatchingUsecase{
		playerRepo: playerRepo,
	}
}

func (uc *MatchingUsecase) FindMatch(ctx context.Context, playerID, username string, rating int32) (string, error) {
	// 创建玩家并设置为等待匹配状态
	player := &entity.Player{
		ID:       playerID,
		Username: username,
		Rating:   rating,
		Status:   "waiting",
	}

	err := uc.playerRepo.CreatePlayer(ctx, player)
	if err != nil {
		return "", errors.New("failed to add player to matching queue")
	}

	// 这里应该实现实际的匹配逻辑
	// 简化实现：直接创建一个房间ID并返回
	roomID := fmt.Sprintf("room_%s", playerID[:8])

	// 更新玩家状态为已匹配
	err = uc.playerRepo.UpdatePlayerStatus(ctx, playerID, "matched")
	if err != nil {
		return "", errors.New("failed to update player status")
	}

	return roomID, nil
}

func (uc *MatchingUsecase) CancelMatch(ctx context.Context, playerID string) error {
	err := uc.playerRepo.UpdatePlayerStatus(ctx, playerID, "idle")
	if err != nil {
		return errors.New("failed to cancel match")
	}
	return nil
}

func (uc *MatchingUsecase) GetWaitingPlayers(ctx context.Context) (int32, error) {
	// 这里应该查询数据库中状态为"waiting"的玩家数量
	// 为了简化，暂时返回一个示例值
	return 0, nil
}

func (uc *MatchingUsecase) GetOnlinePlayers(ctx context.Context) ([]*entity.Player, error) {
	// 这里应该查询数据库中在线的玩家
	// 为了简化，暂时返回空数组
	return []*entity.Player{}, nil
}