package usecase

import (
	"context"
	"errors"
	"time"

	"snake-game/friends/domain/entity"
	"snake-game/friends/domain/repository"
)

type FriendsUsecase struct {
	repo repository.FriendRepository
}

func NewFriendsUsecase(repo repository.FriendRepository) *FriendsUsecase {
	return &FriendsUsecase{
		repo: repo,
	}
}

func (uc *FriendsUsecase) AddFriend(ctx context.Context, userID, friendUsername string) error {
	// 这里需要通过用户服务查找用户名对应的ID
	// 为简化，假设我们知道目标用户的ID
	friendID := friendUsername // 实际应用中需要通过用户名查找ID

	// 检查是否已经是好友
	friendship, err := uc.repo.GetFriendship(ctx, userID, friendID)
	if err == nil && friendship != nil {
		if friendship.Status == "accepted" {
			return errors.New("already friends")
		} else if friendship.Status == "pending" {
			return errors.New("friend request already sent")
		}
	}

	// 创建好友请求
	friendship = &entity.Friendship{
		UserID:    userID,
		FriendID:  friendID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	return uc.repo.CreateFriendship(ctx, friendship)
}

func (uc *FriendsUsecase) RemoveFriend(ctx context.Context, userID, friendUserID string) error {
	return uc.repo.DeleteFriendship(ctx, userID, friendUserID)
}

func (uc *FriendsUsecase) GetFriends(ctx context.Context, userID string) ([]*entity.Friendship, error) {
	return uc.repo.GetFriends(ctx, userID)
}

func (uc *FriendsUsecase) SendFriendRequest(ctx context.Context, userID, targetUserID string) error {
	// 检查是否已经是好友或已有请求
	friendship, err := uc.repo.GetFriendship(ctx, userID, targetUserID)
	if err == nil && friendship != nil {
		if friendship.Status == "accepted" {
			return errors.New("already friends")
		} else if friendship.Status == "pending" {
			return errors.New("friend request already exists")
		}
	}

	// 创建好友请求
	friendship = &entity.Friendship{
		UserID:    userID,
		FriendID:  targetUserID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	return uc.repo.CreateFriendship(ctx, friendship)
}

func (uc *FriendsUsecase) RespondFriendRequest(ctx context.Context, userID, requestUserID string, accepted bool) error {
	friendship, err := uc.repo.GetFriendship(ctx, requestUserID, userID)
	if err != nil || friendship == nil {
		return errors.New("friend request not found")
	}

	if friendship.Status != "pending" {
		return errors.New("friend request already processed")
	}

	status := "rejected"
	if accepted {
		status = "accepted"
	}

	return uc.repo.UpdateFriendshipStatus(ctx, requestUserID, userID, status)
}