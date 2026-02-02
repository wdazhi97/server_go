package grpc

import (
	"context"

	"snake-game/matching/internal/usecase"
	pb "snake-game/proto"
)

type MatchingHandler struct {
	usecase *usecase.MatchingUsecase
	pb.UnimplementedMatchingServiceServer
}

func NewMatchingHandler(usecase *usecase.MatchingUsecase) *MatchingHandler {
	return &MatchingHandler{
		usecase: usecase,
	}
}

// FindMatch 寻找匹配
func (h *MatchingHandler) FindMatch(ctx context.Context, req *pb.FindMatchRequest) (*pb.FindMatchResponse, error) {
	roomID, err := h.usecase.FindMatch(ctx, req.PlayerId, req.Username, req.Rating)
	if err != nil {
		return &pb.FindMatchResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.FindMatchResponse{
		Success: true,
		Message: "Match found successfully",
		RoomId:  roomID,
		Players: []*pb.PlayerInfo{}, // 根据实际需求填充玩家信息
	}, nil
}

// CancelMatch 取消匹配
func (h *MatchingHandler) CancelMatch(ctx context.Context, req *pb.CancelMatchRequest) (*pb.CancelMatchResponse, error) {
	err := h.usecase.CancelMatch(ctx, req.PlayerId)
	if err != nil {
		return &pb.CancelMatchResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.CancelMatchResponse{
		Success: true,
		Message: "Match cancelled successfully",
	}, nil
}

// GetWaitingPlayers 获取等待玩家数量
func (h *MatchingHandler) GetWaitingPlayers(ctx context.Context, req *pb.GetWaitingPlayersRequest) (*pb.GetWaitingPlayersResponse, error) {
	count, err := h.usecase.GetWaitingPlayers(ctx)
	if err != nil {
		return &pb.GetWaitingPlayersResponse{
			Count: 0,
		}, nil
	}

	return &pb.GetWaitingPlayersResponse{
		Count: count,
	}, nil
}

// GetOnlinePlayers 获取在线玩家
func (h *MatchingHandler) GetOnlinePlayers(ctx context.Context, req *pb.GetOnlinePlayersRequest) (*pb.GetOnlinePlayersResponse, error) {
	players, err := h.usecase.GetOnlinePlayers(ctx)
	if err != nil {
		return &pb.GetOnlinePlayersResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 转换域实体到协议缓冲区消息
	pbPlayers := make([]*pb.PlayerInfo, len(players))
	for i, player := range players {
		pbPlayers[i] = &pb.PlayerInfo{
			PlayerId: player.ID,
			Username: player.Username,
			Rating:   player.Rating,
		}
	}

	return &pb.GetOnlinePlayersResponse{
		Success: true,
		Message: "Players retrieved successfully",
		Players: pbPlayers,
	}, nil
}