package grpc

import (
	"context"

	"snake-game/lobby/internal/usecase"
	pb "snake-game/proto"
)

type LobbyHandler struct {
	usecase *usecase.AuthUsecase
	pb.UnimplementedLobbyServiceServer
}

func NewLobbyHandler(usecase *usecase.AuthUsecase) *LobbyHandler {
	return &LobbyHandler{
		usecase: usecase,
	}
}

// Register 用户注册
func (h *LobbyHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userID, err := h.usecase.Register(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  userID,
	}, nil
}

// Login 用户登录
func (h *LobbyHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.usecase.Login(ctx, req.Username, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Success:  true,
		Message:  "Login successful",
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

// GetUserProfile 获取用户资料
func (h *LobbyHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := h.usecase.GetUserProfile(ctx, req.UserId)
	if err != nil {
		return &pb.GetUserProfileResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 获取用户排行榜信息
	// 这里需要根据实际情况实现排行榜获取逻辑

	return &pb.GetUserProfileResponse{
		Success: true,
		User: &pb.User{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Online:    user.Online,
			CreatedAt: user.CreatedAt.Unix(),
		},
		Score:       0, // Placeholder - 实际实现时需要从排行榜服务获取
		GamesWon:    0, // Placeholder
		GamesPlayed: 0, // Placeholder
	}, nil
}

// Logout 用户登出
func (h *LobbyHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.usecase.Logout(ctx, req.UserId)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}