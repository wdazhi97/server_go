package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	grpc_handler "snake-game/matching/internal/delivery/grpc"
	"snake-game/matching/internal/repository"
	"snake-game/matching/internal/usecase"
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
	playerRepo := repository.NewPlayerRepository()

	// 初始化业务逻辑层
	matchingUsecase := usecase.NewMatchingUsecase(playerRepo)

	// 初始化通信层
	matchingHandler := grpc_handler.NewMatchingHandler(matchingUsecase)

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMatchingServiceServer(s, matchingHandler)

	log.Println("Matching service is running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}