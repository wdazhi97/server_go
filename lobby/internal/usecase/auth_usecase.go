package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"snake-game/lobby/domain/entity"
	"snake-game/lobby/domain/repository"
	"snake-game/mongodb"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, username, password, email string) (string, error) {
	// 检查用户名是否已存在
	existingUser, _ := uc.userRepo.FindByUsername(ctx, username)
	if existingUser != nil {
		return "", errors.New("username already exists")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// 创建用户实体
	user := &entity.User{
		Username:  username,
		Password:  string(hashedPassword),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastSeen:  time.Now(),
		Online:    false,
	}

	// 保存用户
	userID, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return "", errors.New("failed to register user")
	}

	// 创建初始排行榜记录
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		// 如果转换失败，记录错误但继续
	} else {
		leaderboardEntry := mongodb.Leaderboard{
			UserID:      objectID,
			Score:       0,
			GamesWon:    0,
			GamesPlayed: 0,
			UpdatedAt:   time.Now(),
		}
		_, err = mongodb.DB.Collection(mongodb.LeaderboardCollection).InsertOne(ctx, leaderboardEntry)
		if err != nil {
			// 不返回错误，因为用户已创建成功
			// 只是记录警告
		}
	}

	return userID, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("invalid username or password")
		}
		return nil, errors.New("internal server error")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 更新用户状态为在线
	err = uc.userRepo.UpdateOnlineStatus(ctx, user.ID, true)
	if err != nil {
		// 记录警告，但不返回错误
	}

	return user, nil
}

func (uc *AuthUsecase) GetUserProfile(ctx context.Context, userID string) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("internal server error")
	}
	return user, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, userID string) error {
	err := uc.userRepo.UpdateOnlineStatus(ctx, userID, false)
	if err != nil {
		return errors.New("failed to update user status")
	}
	return nil
}