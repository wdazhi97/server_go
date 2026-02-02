package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	grpc_handler "snake-game/game/internal/delivery/grpc"
	"snake-game/game/internal/repository"
	"snake-game/game/internal/usecase"
	pb "snake-game/proto"
)

func main() {
	// 初始化仓库层（基于内存的游戏状态管理）
	gameRepo := repository.NewGameMemoryRepository()

	// 初始化业务逻辑层
	gameUsecase := usecase.NewGameUsecase(gameRepo)

	// 初始化通信层
	gameHandler := grpc_handler.NewGameHandler(gameUsecase)

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGameServiceServer(s, gameHandler)

	log.Println("Game service is running on :50055")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}