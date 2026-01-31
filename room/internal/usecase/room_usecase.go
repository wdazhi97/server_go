package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"snake-game/room/domain/entity"
	"snake-game/room/domain/repository"
	pb "snake-game/proto"
)

type RoomUsecase struct {
	roomRepo repository.RoomRepository
	gameClient pb.GameServiceClient
}

func NewRoomUsecase(roomRepo repository.RoomRepository) *RoomUsecase {
	// 连接到游戏服务
	conn, err := grpc.Dial("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}

	gameClient := pb.NewGameServiceClient(conn)

	return &RoomUsecase{
		roomRepo: roomRepo,
		gameClient: gameClient,
	}
}

func (uc *RoomUsecase) CreateRoom(ctx context.Context, userID, roomName string, maxPlayers int32) (string, error) {
	// 生成房间ID
	roomID := fmt.Sprintf("room_%d_%s", time.Now().Unix(), generateRandomString(6))

	room := &entity.Room{
		ID:         roomID,
		Name:       roomName,
		CreatorID:  userID,
		Players:    []string{userID},
		MaxPlayers: int(maxPlayers),
		Status:     "waiting",
		CreatedAt:  time.Now(),
		Messages:   []*entity.Message{},
	}

	err := uc.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		return "", errors.New("failed to create room")
	}

	// 发送系统消息
	systemMsg := &entity.Message{
		ID:         fmt.Sprintf("msg_%d", time.Now().Unix()),
		RoomID:     roomID,
		SenderID:   "system",
		SenderName: "System",
		Content:    fmt.Sprintf("%s 创建了房间", userID),
		Type:       "system",
		CreatedAt:  time.Now(),
	}
	uc.roomRepo.AddMessage(ctx, roomID, systemMsg)

	return roomID, nil
}

func (uc *RoomUsecase) JoinRoom(ctx context.Context, roomID, userID string) error {
	room, err := uc.roomRepo.GetRoom(ctx, roomID)
	if err != nil || room == nil {
		return errors.New("room not found")
	}

	if room.Status != "waiting" {
		return errors.New("cannot join room that is not in waiting state")
	}

	// 检查房间是否已满
	if len(room.Players) >= room.MaxPlayers {
		return errors.New("room is full")
	}

	// 检查玩家是否已在房间中
	for _, pid := range room.Players {
		if pid == userID {
			return errors.New("player already in room")
		}
	}

	// 添加玩家到房间
	err = uc.roomRepo.AddPlayer(ctx, roomID, userID)
	if err != nil {
		return errors.New("failed to add player to room")
	}

	// 发送系统消息
	systemMsg := &entity.Message{
		ID:         fmt.Sprintf("msg_%d", time.Now().Unix()),
		RoomID:     roomID,
		SenderID:   "system",
		SenderName: "System",
		Content:    fmt.Sprintf("%s 加入了房间", userID),
		Type:       "system",
		CreatedAt:  time.Now(),
	}
	uc.roomRepo.AddMessage(ctx, roomID, systemMsg)

	return nil
}

func (uc *RoomUsecase) LeaveRoom(ctx context.Context, roomID, userID string) error {
	room, err := uc.roomRepo.GetRoom(ctx, roomID)
	if err != nil || room == nil {
		return errors.New("room not found")
	}

	// 从房间中移除玩家
	err = uc.roomRepo.RemovePlayer(ctx, roomID, userID)
	if err != nil {
		return errors.New("failed to remove player from room")
	}

	// 发送系统消息
	systemMsg := &entity.Message{
		ID:         fmt.Sprintf("msg_%d", time.Now().Unix()),
		RoomID:     roomID,
		SenderID:   "system",
		SenderName: "System",
		Content:    fmt.Sprintf("%s 离开了房间", userID),
		Type:       "system",
		CreatedAt:  time.Now(),
	}
	uc.roomRepo.AddMessage(ctx, roomID, systemMsg)

	// 如果房主离开且房间还有其他玩家，指定新的房主或解散房间
	if room.CreatorID == userID {
		if len(room.Players) > 0 {
			// 将第一个玩家设为新房主
			room.CreatorID = room.Players[0]
		} else {
			// 没有玩家了，删除房间
			uc.roomRepo.DeleteRoom(ctx, roomID)
		}
	} else if len(room.Players) == 0 {
		// 如果房间里没有玩家了，删除房间
		uc.roomRepo.DeleteRoom(ctx, roomID)
	} else {
		uc.roomRepo.UpdateRoom(ctx, room)
	}

	return nil
}

func (uc *RoomUsecase) SendMessage(ctx context.Context, roomID, senderID, senderName, content, msgType string) error {
	room, err := uc.roomRepo.GetRoom(ctx, roomID)
	if err != nil || room == nil {
		return errors.New("room not found")
	}

	message := &entity.Message{
		ID:         fmt.Sprintf("msg_%d", time.Now().Unix()),
		RoomID:     roomID,
		SenderID:   senderID,
		SenderName: senderName,
		Content:    content,
		Type:       msgType,
		CreatedAt:  time.Now(),
	}

	err = uc.roomRepo.AddMessage(ctx, roomID, message)
	if err != nil {
		return errors.New("failed to send message")
	}

	return nil
}

func (uc *RoomUsecase) GetRoomMessages(ctx context.Context, roomID string, limit int32) ([]*entity.Message, error) {
	messages, err := uc.roomRepo.GetMessages(ctx, roomID, limit)
	if err != nil {
		return nil, errors.New("failed to get messages")
	}

	return messages, nil
}

func (uc *RoomUsecase) StartGame(ctx context.Context, roomID string) error {
	room, err := uc.roomRepo.GetRoom(ctx, roomID)
	if err != nil || room == nil {
		return errors.New("room not found")
	}

	if room.Status != "waiting" {
		return errors.New("game already started or finished")
	}

	if len(room.Players) < 2 {
		return errors.New("not enough players to start game")
	}

	// 更新房间状态
	room.Status = "playing"
	err = uc.roomRepo.UpdateRoom(ctx, room)
	if err != nil {
		return errors.New("failed to update room status")
	}

	// 通知游戏服务开始游戏
	// 注意：这里我们只是更新房间状态，实际的游戏逻辑由游戏服务处理

	return nil
}

// 辅助函数：生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}