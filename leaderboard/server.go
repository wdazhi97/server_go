package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	grpc_handler "snake-game/leaderboard/internal/delivery/grpc"
	"snake-game/leaderboard/internal/repository"
	"snake-game/leaderboard/internal/usecase"
	"snake-game/mongodb"
	pb "snake-game/proto"
)

func main() {
	// 从环境变量获取 MongoDB URI
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// 连接数据库
	err := mongodb.Connect(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Disconnect()

	// 初始化仓库层
	leaderboardRepo := repository.NewLeaderboardRepository()

	// 初始化业务逻辑层
	leaderboardUsecase := usecase.NewLeaderboardUsecase(leaderboardRepo)

	// 初始化通信层
	leaderboardHandler := grpc_handler.NewLeaderboardHandler(leaderboardUsecase)

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLeaderboardServiceServer(s, leaderboardHandler)

	log.Println("Leaderboard service is running on :50054")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}