package grpc

import (
	"context"

	"snake-game/friends/internal/usecase"
	pb "snake-game/proto_new"
)

type FriendsHandler struct {
	usecase *usecase.FriendsUsecase
	pb.UnimplementedFriendsServiceServer
}

func NewFriendsHandler(usecase *usecase.FriendsUsecase) *FriendsHandler {
	return &FriendsHandler{
		usecase: usecase,
	}
}

// AddFriend 添加好友
func (h *FriendsHandler) AddFriend(ctx context.Context, req *pb.AddFriendRequest) (*pb.AddFriendResponse, error) {
	err := h.usecase.AddFriend(ctx, req.UserId, req.FriendUsername)
	if err != nil {
		return &pb.AddFriendResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.AddFriendResponse{
		Success: true,
		Message: "Friend request sent successfully",
	}, nil
}

// RemoveFriend 删除好友
func (h *FriendsHandler) RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (*pb.RemoveFriendResponse, error) {
	err := h.usecase.RemoveFriend(ctx, req.UserId, req.FriendUserId)
	if err != nil {
		return &pb.RemoveFriendResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RemoveFriendResponse{
		Success: true,
		Message: "Friend removed successfully",
	}, nil
}

// GetFriends 获取好友列表
func (h *FriendsHandler) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	friendships, err := h.usecase.GetFriends(ctx, req.UserId)
	if err != nil {
		return &pb.GetFriendsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 转换为协议缓冲区格式
	pbFriends := make([]*pb.FriendInfo, len(friendships))
	for i, friendship := range friendships {
		pbFriends[i] = &pb.FriendInfo{
			UserId:   friendship.FriendID, // 需要根据实际关系调整
			Username: "Unknown", // 需要从用户服务获取用户名
			Online:   false,     // 需要从用户服务获取在线状态
			Status:   friendship.Status,
		}
	}

	return &pb.GetFriendsResponse{
		Success: true,
		Message: "Friends retrieved successfully",
		Friends: pbFriends,
	}, nil
}

// SendFriendRequest 发送好友请求
func (h *FriendsHandler) SendFriendRequest(ctx context.Context, req *pb.SendFriendRequestRequest) (*pb.SendFriendRequestResponse, error) {
	err := h.usecase.SendFriendRequest(ctx, req.UserId, req.TargetUserId)
	if err != nil {
		return &pb.SendFriendRequestResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SendFriendRequestResponse{
		Success: true,
		Message: "Friend request sent successfully",
	}, nil
}

// RespondFriendRequest 回应好友请求
func (h *FriendsHandler) RespondFriendRequest(ctx context.Context, req *pb.RespondFriendRequestRequest) (*pb.RespondFriendRequestResponse, error) {
	err := h.usecase.RespondFriendRequest(ctx, req.UserId, req.RequestUserId, req.Accepted)
	if err != nil {
		return &pb.RespondFriendRequestResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	status := "rejected"
	if req.Accepted {
		status = "accepted"
	}

	return &pb.RespondFriendRequestResponse{
		Success: true,
		Message: "Friend request " + status,
	}, nil
}