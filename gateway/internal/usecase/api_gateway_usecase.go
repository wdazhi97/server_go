package usecase

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"snake-game/gateway/domain/repository"
	pb "snake-game/proto_new"
)

type APIGatewayUsecase struct {
	serviceRegistry repository.ServiceRegistry
}

func NewAPIGatewayUsecase(serviceRegistry repository.ServiceRegistry) *APIGatewayUsecase {
	return &APIGatewayUsecase{
		serviceRegistry: serviceRegistry,
	}
}

// ForwardRequest 转发请求到后端服务
func (uc *APIGatewayUsecase) ForwardRequest(c *gin.Context, serviceName string) {
	service, err := uc.serviceRegistry.GetService(context.Background(), serviceName)
	if err != nil || service == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
		return
	}

	if !service.Health {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service is not healthy"})
		return
	}

	// 解析请求体
	var reqBody map[string]interface{}
	if c.Request.Body != nil {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if len(body) > 0 {
			if err := json.Unmarshal(body, &reqBody); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
				return
			}
		}
	}

	// 根据服务类型决定如何处理请求
	switch serviceName {
	case "lobby":
		uc.handleLobbyRequest(c, service.Address, reqBody)
	case "matching":
		uc.handleMatchingRequest(c, service.Address, reqBody)
	case "room":
		uc.handleRoomRequest(c, service.Address, reqBody)
	case "leaderboard":
		uc.handleLeaderboardRequest(c, service.Address, reqBody)
	case "game":
		uc.handleGameRequest(c, service.Address, reqBody)
	case "friends":
		uc.handleFriendsRequest(c, service.Address, reqBody)
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
	}
}

// handleLobbyRequest 处理大厅服务请求
func (uc *APIGatewayUsecase) handleLobbyRequest(c *gin.Context, address string, reqBody map[string]interface{}) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to lobby service"})
		return
	}
	defer conn.Close()

	clientLobby := pb.NewLobbyServiceClient(conn)

	action := strings.TrimPrefix(c.Request.URL.Path, "/auth/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch action {
	case "register":
		username, ok1 := reqBody["username"].(string)
		password, ok2 := reqBody["password"].(string)
		email, ok3 := reqBody["email"].(string)
		
		if !ok1 || !ok2 || !ok3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientLobby.Register(ctx, &pb.RegisterRequest{
			Username: username,
			Password: password,
			Email:    email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
			"userId":  resp.UserId,
		})

	case "login":
		username, ok1 := reqBody["username"].(string)
		password, ok2 := reqBody["password"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientLobby.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  resp.Success,
			"message":  resp.Message,
			"userId":   resp.UserId,
			"username": resp.Username,
		})

	case "getUserProfile":
		userId, ok := reqBody["userId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId"})
			return
		}

		resp, err := clientLobby.GetUserProfile(ctx, &pb.GetUserProfileRequest{
			UserId: userId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"user":    resp.User,
			"score":   resp.Score,
		})

	case "logout":
		userId, ok := reqBody["userId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId"})
			return
		}

		resp, err := clientLobby.Logout(ctx, &pb.LogoutRequest{
			UserId: userId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
	}
}

// handleMatchingRequest 处理匹配服务请求
func (uc *APIGatewayUsecase) handleMatchingRequest(c *gin.Context, address string, reqBody map[string]interface{}) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to matching service"})
		return
	}
	defer conn.Close()

	clientMatching := pb.NewMatchingServiceClient(conn)

	action := strings.TrimPrefix(c.Request.URL.Path, "/match/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch action {
	case "findMatch":
		playerId, ok1 := reqBody["playerId"].(string)
		username, ok2 := reqBody["username"].(string)
		rating, ok3 := reqBody["rating"].(float64)
		
		if !ok1 || !ok2 || !ok3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientMatching.FindMatch(ctx, &pb.FindMatchRequest{
			PlayerId: playerId,
			Username: username,
			Rating:   int32(rating),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
			"roomId":  resp.RoomId,
			"players": resp.Players,
		})

	case "cancelMatch":
		playerId, ok := reqBody["playerId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing playerId"})
			return
		}

		resp, err := clientMatching.CancelMatch(ctx, &pb.CancelMatchRequest{
			PlayerId: playerId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "getWaitingPlayers":
		resp, err := clientMatching.GetWaitingPlayers(ctx, &pb.GetWaitingPlayersRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count": resp.Count,
		})

	case "getOnlinePlayers":
		resp, err := clientMatching.GetOnlinePlayers(ctx, &pb.GetOnlinePlayersRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"players": resp.Players,
		})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
	}
}

// handleRoomRequest 处理房间服务请求
func (uc *APIGatewayUsecase) handleRoomRequest(c *gin.Context, address string, reqBody map[string]interface{}) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to room service"})
		return
	}
	defer conn.Close()

	clientRoom := pb.NewRoomServiceClient(conn)

	action := strings.TrimPrefix(c.Request.URL.Path, "/room/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch action {
	case "createRoom":
		userId, ok1 := reqBody["userId"].(string)
		roomName, ok2 := reqBody["roomName"].(string)
		maxPlayersFloat, ok3 := reqBody["maxPlayers"].(float64)
		
		if !ok1 || !ok2 || !ok3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}
		maxPlayers := int32(maxPlayersFloat)

		resp, err := clientRoom.CreateRoom(ctx, &pb.CreateRoomRequest{
			UserId:     userId,
			RoomName:   roomName,
			MaxPlayers: maxPlayers,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"roomId":  resp.RoomId,
			"message": resp.Message,
		})

	case "joinRoom":
		roomId, ok1 := reqBody["roomId"].(string)
		userId, ok2 := reqBody["userId"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientRoom.JoinRoom(ctx, &pb.JoinRoomRequest{
			RoomId: roomId,
			UserId: userId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "leaveRoom":
		roomId, ok1 := reqBody["roomId"].(string)
		userId, ok2 := reqBody["userId"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientRoom.LeaveRoom(ctx, &pb.LeaveRoomRequest{
			RoomId: roomId,
			UserId: userId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "sendMessage":
		roomId, ok1 := reqBody["roomId"].(string)
		senderId, ok2 := reqBody["senderId"].(string)
		content, ok3 := reqBody["content"].(string)
		msgType, ok4 := reqBody["type"].(string)
		
		if !ok1 || !ok2 || !ok3 || !ok4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientRoom.SendMessage(ctx, &pb.SendMessageRequest{
			RoomId:    roomId,
			SenderId:  senderId,
			Content:   content,
			Type:      msgType,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "getRoomMessages":
		roomId, ok1 := reqBody["roomId"].(string)
		limitFloat, hasLimit := reqBody["limit"].(float64)
		sinceFloat, hasSince := reqBody["since"].(float64)
		
		if !ok1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing roomId"})
			return
		}
		
		var limit int32
		var since int64
		if hasLimit {
			limit = int32(limitFloat)
		} else {
			limit = 50 // 默认值
		}
		
		if hasSince {
			since = int64(sinceFloat)
		}

		resp, err := clientRoom.GetRoomMessages(ctx, &pb.GetRoomMessagesRequest{
			RoomId: roomId,
			Limit:  limit,
			Since:  since,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  resp.Success,
			"messages": resp.Messages,
		})

	case "startGame":
		roomId, ok := reqBody["roomId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing roomId"})
			return
		}

		resp, err := clientRoom.StartGame(ctx, &pb.StartGameRequest{
			RoomId: roomId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
	}
}

// handleLeaderboardRequest 处理排行榜服务请求
func (uc *APIGatewayUsecase) handleLeaderboardRequest(c *gin.Context, address string, reqBody map[string]interface{}) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to leaderboard service"})
		return
	}
	defer conn.Close()

	clientLeaderboard := pb.NewLeaderboardServiceClient(conn)

	action := strings.TrimPrefix(c.Request.URL.Path, "/leaderboard/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch action {
	case "getLeaderboard":
		limitFloat, hasLimit := reqBody["limit"].(float64)
		offsetFloat, hasOffset := reqBody["offset"].(float64)
		
		var limit, offset int32
		if hasLimit {
			limit = int32(limitFloat)
		} else {
			limit = 10 // 默认值
		}
		
		if hasOffset {
			offset = int32(offsetFloat)
		}

		resp, err := clientLeaderboard.GetLeaderboard(ctx, &pb.GetLeaderboardRequest{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"entries": resp.Entries,
		})

	case "updateScore":
		userId, ok1 := reqBody["userId"].(string)
		scoreFloat, ok2 := reqBody["score"].(float64)
		gameWon, ok3 := reqBody["gameWon"].(bool)
		
		if !ok1 || !ok2 || !ok3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}
		score := int32(scoreFloat)

		resp, err := clientLeaderboard.UpdateScore(ctx, &pb.UpdateScoreRequest{
			UserId:  userId,
			Score:   score,
			GameWon: gameWon,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "getUserRank":
		userId, ok := reqBody["userId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId"})
			return
		}

		resp, err := clientLeaderboard.GetUserRank(ctx, &pb.GetUserRankRequest{
			UserId: userId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":    resp.Success,
			"rank":       resp.Rank,
			"totalUsers": resp.TotalUsers,
		})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
	}
}

// handleGameRequest 处理游戏服务请求
func (uc *APIGatewayUsecase) handleGameRequest(c *gin.Context, address string, reqBody map[string]interface{}) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to game service"})
		return
	}
	defer conn.Close()

	clientGame := pb.NewGameServiceClient(conn)

	action := strings.TrimPrefix(c.Request.URL.Path, "/game/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch action {
	case "joinGame":
		roomId, ok1 := reqBody["roomId"].(string)
		playerId, ok2 := reqBody["playerId"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientGame.JoinGame(ctx, &pb.JoinGameRequest{
			RoomId:   roomId,
			PlayerId: playerId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "leaveGame":
		roomId, ok1 := reqBody["roomId"].(string)
		playerId, ok2 := reqBody["playerId"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientGame.LeaveGame(ctx, &pb.LeaveGameRequest{
			RoomId:   roomId,
			PlayerId: playerId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "move":
		roomId, ok1 := reqBody["roomId"].(string)
		playerId, ok2 := reqBody["playerId"].(string)
		directionStr, ok3 := reqBody["direction"].(string)
		
		if !ok1 || !ok2 || !ok3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// 将字符串方向转换为 protobuf 枚举
		var direction pb.Direction
		switch directionStr {
		case "UP":
			direction = pb.Direction_UP
		case "DOWN":
			direction = pb.Direction_DOWN
		case "LEFT":
			direction = pb.Direction_LEFT
		case "RIGHT":
			direction = pb.Direction_RIGHT
		default:
			direction = pb.Direction_NONE
		}

		resp, err := clientGame.Move(ctx, &pb.MoveRequest{
			RoomId:      roomId,
			PlayerId:    playerId,
			Direction:   direction,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "getGameState":
		roomId, ok := reqBody["roomId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing roomId"})
			return
		}

		resp, err := clientGame.GetGameState(ctx, &pb.GetGameStateRequest{
			RoomId: roomId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"snakes":  resp.Snakes,
			"foods":   resp.Foods,
			"walls":   resp.Walls,
			"status":  resp.Status,
		})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
	}
}

// handleFriendsRequest 处理好友服务请求
func (uc *APIGatewayUsecase) handleFriendsRequest(c *gin.Context, address string, reqBody map[string]interface{}) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to friends service"})
		return
	}
	defer conn.Close()

	clientFriends := pb.NewFriendsServiceClient(conn)

	action := strings.TrimPrefix(c.Request.URL.Path, "/friends/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch action {
	case "addFriend":
		userId, ok1 := reqBody["userId"].(string)
		friendUsername, ok2 := reqBody["friendUsername"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientFriends.AddFriend(ctx, &pb.AddFriendRequest{
			UserId:         userId,
			FriendUsername: friendUsername,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "removeFriend":
		userId, ok1 := reqBody["userId"].(string)
		friendUserId, ok2 := reqBody["friendUserId"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientFriends.RemoveFriend(ctx, &pb.RemoveFriendRequest{
			UserId:      userId,
			FriendUserId: friendUserId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "getFriends":
		userId, ok := reqBody["userId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId"})
			return
		}

		resp, err := clientFriends.GetFriends(ctx, &pb.GetFriendsRequest{
			UserId: userId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"friends": resp.Friends,
		})

	case "sendFriendRequest":
		userId, ok1 := reqBody["userId"].(string)
		targetUserId, ok2 := reqBody["targetUserId"].(string)
		
		if !ok1 || !ok2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientFriends.SendFriendRequest(ctx, &pb.SendFriendRequestRequest{
			UserId:       userId,
			TargetUserId: targetUserId,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	case "respondFriendRequest":
		userId, ok1 := reqBody["userId"].(string)
		requestUserId, ok2 := reqBody["requestUserId"].(string)
		accepted, ok3 := reqBody["accepted"].(bool)
		
		if !ok1 || !ok2 || !ok3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		resp, err := clientFriends.RespondFriendRequest(ctx, &pb.RespondFriendRequestRequest{
			UserId:        userId,
			RequestUserId: requestUserId,
			Accepted:      accepted,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
			"message": resp.Message,
		})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
	}
}

// HealthCheck 检查服务健康状况
func (uc *APIGatewayUsecase) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"timestamp":   time.Now().Unix(),
		"service":     "api-gateway",
		"version":     "1.0.0",
	})
}

// GetServiceDiscovery 获取服务发现信息
func (uc *APIGatewayUsecase) GetServiceDiscovery(c *gin.Context) {
	services, err := uc.serviceRegistry.GetAllServices(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get services"})
		return
	}

	serviceList := make([]map[string]interface{}, len(services))
	for i, service := range services {
		serviceList[i] = map[string]interface{}{
			"name":    service.Name,
			"address": service.Address,
			"health":  service.Health,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"services": serviceList,
	})
}