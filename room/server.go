package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	grpc_handler "snake-game/room/internal/delivery/grpc"
	"snake-game/room/internal/repository"
	"snake-game/room/internal/usecase"
	pb "snake-game/proto"
)

func main() {
	// 初始化仓库层（基于内存的房间管理）
	roomRepo := repository.NewRoomMemoryRepository()

	// 初始化业务逻辑层
	roomUsecase := usecase.NewRoomUsecase(roomRepo)

	// 初始化通信层
	roomHandler := grpc_handler.NewRoomHandler(roomUsecase)

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRoomServiceServer(s, roomHandler)

	log.Println("Room service is running on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}