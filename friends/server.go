package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	grpc_handler "snake-game/friends/internal/delivery/grpc"
	"snake-game/friends/internal/repository"
	"snake-game/friends/internal/usecase"
	"snake-game/mongodb"
	pb "snake-game/proto_new"
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
	friendRepo := repository.NewFriendRepository()

	// 初始化业务逻辑层
	friendsUsecase := usecase.NewFriendsUsecase(friendRepo)

	// 初始化通信层
	friendsHandler := grpc_handler.NewFriendsHandler(friendsUsecase)

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":50056")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFriendsServiceServer(s, friendsHandler)

	log.Println("Friends service is running on :50056")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}