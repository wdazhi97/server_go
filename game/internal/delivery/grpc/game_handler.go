package grpc

import (
	"context"

	"snake-game/game/domain/entity"
	"snake-game/game/internal/usecase"
	pb "snake-game/proto_new"
)

type GameHandler struct {
	usecase *usecase.GameUsecase
	pb.UnimplementedGameServiceServer
}

func NewGameHandler(usecase *usecase.GameUsecase) *GameHandler {
	return &GameHandler{
		usecase: usecase,
	}
}

// JoinGame 加入游戏
func (h *GameHandler) JoinGame(ctx context.Context, req *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	err := h.usecase.JoinGame(ctx, req.RoomId, req.PlayerId)
	if err != nil {
		return &pb.JoinGameResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.JoinGameResponse{
		Success: true,
		Message: "Joined game successfully",
	}, nil
}

// LeaveGame 离开游戏
func (h *GameHandler) LeaveGame(ctx context.Context, req *pb.LeaveGameRequest) (*pb.LeaveGameResponse, error) {
	err := h.usecase.LeaveGame(ctx, req.RoomId, req.PlayerId)
	if err != nil {
		return &pb.LeaveGameResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LeaveGameResponse{
		Success: true,
		Message: "Left game successfully",
	}, nil
}

// Move 移动指令
func (h *GameHandler) Move(ctx context.Context, req *pb.MoveRequest) (*pb.MoveResponse, error) {
	var direction entity.Direction
	switch req.Direction {
	case pb.Direction_UP:
		direction = entity.Direction_UP
	case pb.Direction_DOWN:
		direction = entity.Direction_DOWN
	case pb.Direction_LEFT:
		direction = entity.Direction_LEFT
	case pb.Direction_RIGHT:
		direction = entity.Direction_RIGHT
	default:
		direction = entity.Direction_NONE
	}

	err := h.usecase.Move(ctx, req.RoomId, req.PlayerId, direction)
	if err != nil {
		return &pb.MoveResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.MoveResponse{
		Success: true,
		Message: "Move processed successfully",
	}, nil
}

// GetGameState 获取游戏状态
func (h *GameHandler) GetGameState(ctx context.Context, req *pb.GetGameStateRequest) (*pb.GetGameStateResponse, error) {
	gameState, err := h.usecase.GetGameState(ctx, req.RoomId)
	if err != nil {
		return &pb.GetGameStateResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	if gameState == nil {
		return &pb.GetGameStateResponse{
			Success: false,
			Message: "Game not found",
		}, nil
	}

	// 转换内部实体到协议缓冲区消息
	pbSnakes := make([]*pb.GameSnake, 0, len(gameState.Snakes))
	for _, snake := range gameState.Snakes {
		pbSegments := make([]*pb.SnakeSegment, len(snake.Segments))
		for i, segment := range snake.Segments {
			pbSegments[i] = &pb.SnakeSegment{
				Position: &pb.Position{
					X: segment.Position.X,
					Y: segment.Position.Y,
				},
			}
		}

		pbSnakes = append(pbSnakes, &pb.GameSnake{
			PlayerId:  snake.PlayerID,
			Segments:  pbSegments,
			Color:     snake.Color,
			Length:    int32(snake.Length),
			Score:     int32(snake.Score),
		})
	}

	pbFoods := make([]*pb.Position, len(gameState.Foods))
	for i, food := range gameState.Foods {
		pbFoods[i] = &pb.Position{
			X: food.X,
			Y: food.Y,
		}
	}

	pbWalls := make([]*pb.Position, len(gameState.Walls))
	for i, wall := range gameState.Walls {
		pbWalls[i] = &pb.Position{
			X: wall.X,
			Y: wall.Y,
		}
	}

	return &pb.GetGameStateResponse{
		Success: true,
		Message: "Game state retrieved successfully",
		Snakes:  pbSnakes,
		Foods:   pbFoods,
		Walls:   pbWalls,
		Status:  gameState.Status,
	}, nil
}

// SubscribeGameUpdates 订阅游戏状态更新
func (h *GameHandler) SubscribeGameUpdates(req *pb.SubscribeGameUpdatesRequest, stream pb.GameService_SubscribeGameUpdatesServer) error {
	// 这里应该实现 WebSocket 或流式更新逻辑
	// 为简化，暂时返回未实现
	return nil
}