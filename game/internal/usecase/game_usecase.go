package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"snake-game/game/domain/entity"
	"snake-game/game/domain/repository"
	pb "snake-game/proto"
)

type GameUsecase struct {
	gameRepo     repository.GameRepository
	leaderboardClient pb.LeaderboardServiceClient
}

func NewGameUsecase(gameRepo repository.GameRepository) *GameUsecase {
	// 连接到排行榜服务
	conn, err := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to leaderboard service: %v", err)
		return nil
	}

	leaderboardClient := pb.NewLeaderboardServiceClient(conn)

	return &GameUsecase{
		gameRepo:          gameRepo,
		leaderboardClient: leaderboardClient,
	}
}

func (uc *GameUsecase) JoinGame(ctx context.Context, roomID, playerID string) error {
	game, err := uc.gameRepo.GetGame(ctx, roomID)
	if err != nil {
		return errors.New("failed to get game")
	}

	if game == nil {
		// 如果游戏不存在，创建新游戏
		game = &entity.GameState{
			ID:     fmt.Sprintf("game_%s", roomID),
			RoomID: roomID,
			Snakes: make(map[string]*entity.GameSnake),
			Foods:  []entity.Position{},
			Walls:  []entity.Position{},
			Status: "playing",
		}
		game.Snakes[playerID] = &entity.GameSnake{
			PlayerID:  playerID,
			Segments:  []entity.SnakeSegment{{Position: entity.Position{X: 10, Y: 10}}},
			Color:     fmt.Sprintf("#%06x", rand.Intn(0xffffff)),
			Length:    1,
			Score:     0,
			Alive:     true,
			Direction: entity.Direction_RIGHT,
		}
		// 生成一些食物
		for i := 0; i < 5; i++ {
			game.Foods = append(game.Foods, entity.Position{X: int32(rand.Intn(20)), Y: int32(rand.Intn(20))})
		}
		err = uc.gameRepo.CreateGame(ctx, game)
	} else {
		// 如果游戏存在，添加玩家
		game.Snakes[playerID] = &entity.GameSnake{
			PlayerID:  playerID,
			Segments:  []entity.SnakeSegment{{Position: entity.Position{X: 10, Y: 10}}},
			Color:     fmt.Sprintf("#%06x", rand.Intn(0xffffff)),
			Length:    1,
			Score:     0,
			Alive:     true,
			Direction: entity.Direction_RIGHT,
		}
		err = uc.gameRepo.UpdateGame(ctx, game)
	}

	return err
}

func (uc *GameUsecase) LeaveGame(ctx context.Context, roomID, playerID string) error {
	game, err := uc.gameRepo.GetGame(ctx, roomID)
	if err != nil || game == nil {
		return errors.New("game not found")
	}

	// 从游戏中移除玩家
	delete(game.Snakes, playerID)

	// 检查是否还有其他玩家，如果没有则删除游戏
	if len(game.Snakes) == 0 {
		return uc.gameRepo.DeleteGame(ctx, roomID)
	}

	return uc.gameRepo.UpdateGame(ctx, game)
}

func (uc *GameUsecase) Move(ctx context.Context, roomID, playerID string, direction entity.Direction) error {
	game, err := uc.gameRepo.GetGame(ctx, roomID)
	if err != nil || game == nil {
		return errors.New("game not found")
	}

	snake, exists := game.Snakes[playerID]
	if !exists || !snake.Alive {
		return errors.New("player not in game or dead")
	}

	// 更新蛇的方向
	snake.Direction = direction

	// 这里应该实现实际的游戏移动逻辑
	// 为简化，只是更新蛇的位置
	head := snake.Segments[0].Position
	var newHead entity.Position

	switch direction {
	case entity.Direction_UP:
		newHead = entity.Position{X: head.X, Y: head.Y - 1}
	case entity.Direction_DOWN:
		newHead = entity.Position{X: head.X, Y: head.Y + 1}
	case entity.Direction_LEFT:
		newHead = entity.Position{X: head.X - 1, Y: head.Y}
	case entity.Direction_RIGHT:
		newHead = entity.Position{X: head.X + 1, Y: head.Y}
	default:
		return errors.New("invalid direction")
	}

	// 检查边界
	if newHead.X < 0 || newHead.X >= 20 || newHead.Y < 0 || newHead.Y >= 20 {
		snake.Alive = false
		// 如果蛇死亡，检查是否所有人都死了，如果是则结束游戏
		allDead := true
		for _, s := range game.Snakes {
			if s.Alive {
				allDead = false
				break
			}
		}
		if allDead {
			uc.endGame(ctx, game)
		}
		return uc.gameRepo.UpdateGame(ctx, game)
	}

	// 检查食物
	ateFood := false
	for i, food := range game.Foods {
		if food.X == newHead.X && food.Y == newHead.Y {
			// 吃到了食物
			ateFood = true
			// 移除被吃的食物
			game.Foods = append(game.Foods[:i], game.Foods[i+1:]...)
			// 生成新食物
			game.Foods = append(game.Foods, entity.Position{
				X: int32(rand.Intn(20)),
				Y: int32(rand.Intn(20)),
			})
			snake.Length++
			snake.Score += 10
			break
		}
	}

	// 更新蛇的位置
	newSegments := []entity.SnakeSegment{{Position: newHead}}
	if ateFood {
		// 如果吃了食物，增加身体长度
		newSegments = append(newSegments, snake.Segments...)
	} else {
		// 否则移动身体
		newSegments = append(newSegments, snake.Segments[:len(snake.Segments)-1]...)
	}
	snake.Segments = newSegments

	return uc.gameRepo.UpdateGame(ctx, game)
}

func (uc *GameUsecase) GetGameState(ctx context.Context, roomID string) (*entity.GameState, error) {
	return uc.gameRepo.GetGame(ctx, roomID)
}

// endGame 结束游戏并更新排行榜
func (uc *GameUsecase) endGame(ctx context.Context, game *entity.GameState) {
	game.Status = "finished"
	
	// 更新每个玩家的分数到排行榜服务
	for playerID, snake := range game.Snakes {
		// 调用排行榜服务更新分数
		req := &pb.UpdateScoreRequest{
			UserId:  playerID,
			Score:   int32(snake.Score),
			GameWon: snake.Alive, // 如果蛇还活着，算作获胜
		}
		
		// 在 goroutine 中异步更新排行榜，避免阻塞游戏
		go func(pID string, scoreReq *pb.UpdateScoreRequest) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			_, err := uc.leaderboardClient.UpdateScore(ctx, scoreReq)
			if err != nil {
				log.Printf("Failed to update score for player %s: %v", pID, err)
			}
		}(playerID, req)
	}
}