package grpc

import (
	"context"

	"snake-game/leaderboard/internal/usecase"
	pb "snake-game/proto_new"
)

type LeaderboardHandler struct {
	usecase *usecase.LeaderboardUsecase
	pb.UnimplementedLeaderboardServiceServer
}

func NewLeaderboardHandler(usecase *usecase.LeaderboardUsecase) *LeaderboardHandler {
	return &LeaderboardHandler{
		usecase: usecase,
	}
}

// GetLeaderboard 获取排行榜
func (h *LeaderboardHandler) GetLeaderboard(ctx context.Context, req *pb.GetLeaderboardRequest) (*pb.GetLeaderboardResponse, error) {
	entries, _, err := h.usecase.GetLeaderboard(ctx, req.Limit, req.Offset)
	if err != nil {
		return &pb.GetLeaderboardResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 转换条目格式
	pbEntries := make([]*pb.LeaderboardEntry, len(entries))
	for i, entry := range entries {
		pbEntries[i] = &pb.LeaderboardEntry{
			UserId:   entry.UserID,
			Username: "Unknown", // 需要从用户服务获取用户名
			Score:    int32(entry.Score),
			Rank:     int32(i + 1), // 假设 entries 已排序
		}
	}

	return &pb.GetLeaderboardResponse{
		Success: true,
		Message: "Leaderboard retrieved successfully",
		Entries: pbEntries,
	}, nil
}

// UpdateScore 更新分数
func (h *LeaderboardHandler) UpdateScore(ctx context.Context, req *pb.UpdateScoreRequest) (*pb.UpdateScoreResponse, error) {
	err := h.usecase.UpdateScore(ctx, req.UserId, req.Score, req.GameWon)
	if err != nil {
		return &pb.UpdateScoreResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 获取更新后的排名
	rank, _, err := h.usecase.GetUserRank(ctx, req.UserId)
	if err != nil {
		// 如果无法获取排名，仍然返回成功
		return &pb.UpdateScoreResponse{
			Success: true,
			Message: "Score updated successfully",
			NewScore: req.Score,
		}, nil
	}

	return &pb.UpdateScoreResponse{
		Success: true,
		Message: "Score updated successfully",
		NewScore: req.Score,
		Rank:    rank,
	}, nil
}

// GetUserRank 获取用户排名
func (h *LeaderboardHandler) GetUserRank(ctx context.Context, req *pb.GetUserRankRequest) (*pb.GetUserRankResponse, error) {
	rank, total, err := h.usecase.GetUserRank(ctx, req.UserId)
	if err != nil {
		return &pb.GetUserRankResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetUserRankResponse{
		Success: true,
		Message: "Rank retrieved successfully",
		Rank:    rank,
		TotalUsers: total,
	}, nil
}