package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	grpc_handler "snake-game/lobby/internal/delivery/grpc"
	"snake-game/lobby/internal/repository"
	"snake-game/lobby/internal/usecase"
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
	userRepo := repository.NewUserRepository()

	// 初始化业务逻辑层
	authUsecase := usecase.NewAuthUsecase(userRepo)

	// 初始化通信层
	lobbyHandler := grpc_handler.NewLobbyHandler(authUsecase)

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLobbyServiceServer(s, lobbyHandler)

	log.Println("Lobby service is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}