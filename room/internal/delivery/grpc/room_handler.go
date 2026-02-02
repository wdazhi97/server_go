package grpc

import (
	"context"

	"snake-game/room/internal/usecase"
	pb "snake-game/proto"
)

type RoomHandler struct {
	usecase *usecase.RoomUsecase
	pb.UnimplementedRoomServiceServer
}

func NewRoomHandler(usecase *usecase.RoomUsecase) *RoomHandler {
	return &RoomHandler{
		usecase: usecase,
	}
}

// CreateRoom 创建房间
func (h *RoomHandler) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	roomID, err := h.usecase.CreateRoom(ctx, req.UserId, req.RoomName, req.MaxPlayers)
	if err != nil {
		return &pb.CreateRoomResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.CreateRoomResponse{
		Success: true,
		Message: "Room created successfully",
		RoomId:  roomID,
	}, nil
}

// JoinRoom 加入房间
func (h *RoomHandler) JoinRoom(ctx context.Context, req *pb.JoinRoomRequest) (*pb.JoinRoomResponse, error) {
	err := h.usecase.JoinRoom(ctx, req.RoomId, req.UserId)
	if err != nil {
		return &pb.JoinRoomResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.JoinRoomResponse{
		Success: true,
		Message: "Joined room successfully",
	}, nil
}

// LeaveRoom 离开房间
func (h *RoomHandler) LeaveRoom(ctx context.Context, req *pb.LeaveRoomRequest) (*pb.LeaveRoomResponse, error) {
	err := h.usecase.LeaveRoom(ctx, req.RoomId, req.UserId)
	if err != nil {
		return &pb.LeaveRoomResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LeaveRoomResponse{
		Success: true,
		Message: "Left room successfully",
	}, nil
}

// SendMessage 发送消息
func (h *RoomHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	err := h.usecase.SendMessage(ctx, req.RoomId, req.SenderId, req.SenderId, req.Content, req.Type)
	if err != nil {
		return &pb.SendMessageResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SendMessageResponse{
		Success: true,
		Message: "Message sent successfully",
	}, nil
}

// GetRoomMessages 获取房间消息
func (h *RoomHandler) GetRoomMessages(ctx context.Context, req *pb.GetRoomMessagesRequest) (*pb.GetRoomMessagesResponse, error) {
	messages, err := h.usecase.GetRoomMessages(ctx, req.RoomId, req.Limit)
	if err != nil {
		return &pb.GetRoomMessagesResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 转换消息格式
	pbMessages := make([]*pb.Message, len(messages))
	for i, msg := range messages {
		pbMessages[i] = &pb.Message{
			Id:             msg.ID,
			RoomId:         msg.RoomID,
			SenderId:       msg.SenderID,
			SenderUsername: msg.SenderName,
			Content:        msg.Content,
			Type:           msg.Type,
			CreatedAt:      msg.CreatedAt.Unix(),
		}
	}

	return &pb.GetRoomMessagesResponse{
		Success:  true,
		Message:  "Messages retrieved successfully",
		Messages: pbMessages,
	}, nil
}

// StartGame 开始游戏
func (h *RoomHandler) StartGame(ctx context.Context, req *pb.StartGameRequest) (*pb.StartGameResponse, error) {
	err := h.usecase.StartGame(ctx, req.RoomId)
	if err != nil {
		return &pb.StartGameResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.StartGameResponse{
		Success: true,
		Message: "Game started successfully",
	}, nil
}