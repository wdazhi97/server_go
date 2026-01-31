package http

import (
	"github.com/gin-gonic/gin"
	"snake-game/gateway/internal/usecase"
)

type GatewayHandler struct {
	usecase *usecase.APIGatewayUsecase
}

func NewGatewayHandler(usecase *usecase.APIGatewayUsecase) *GatewayHandler {
	return &GatewayHandler{
		usecase: usecase,
	}
}

// SetupRoutes 设置路由
func (h *GatewayHandler) SetupRoutes(r *gin.Engine) {
	// 健康检查端点
	r.GET("/health", h.usecase.HealthCheck)

	// 服务发现端点
	r.GET("/discovery", h.usecase.GetServiceDiscovery)

	// 认证相关路由
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "lobby")
		})
		authGroup.POST("/login", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "lobby")
		})
		authGroup.POST("/getUserProfile", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "lobby")
		})
		authGroup.POST("/logout", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "lobby")
		})
	}

	// 匹配相关路由
	matchGroup := r.Group("/match")
	{
		matchGroup.POST("/findMatch", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "matching")
		})
		matchGroup.POST("/cancelMatch", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "matching")
		})
		matchGroup.POST("/getWaitingPlayers", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "matching")
		})
		matchGroup.POST("/getOnlinePlayers", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "matching")
		})
	}

	// 房间相关路由
	roomGroup := r.Group("/room")
	{
		roomGroup.POST("/createRoom", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "room")
		})
		roomGroup.POST("/joinRoom", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "room")
		})
		roomGroup.POST("/leaveRoom", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "room")
		})
		roomGroup.POST("/sendMessage", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "room")
		})
		roomGroup.POST("/getRoomMessages", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "room")
		})
		roomGroup.POST("/startGame", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "room")
		})
	}

	// 排行榜相关路由
	leaderboardGroup := r.Group("/leaderboard")
	{
		leaderboardGroup.POST("/getLeaderboard", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "leaderboard")
		})
		leaderboardGroup.POST("/updateScore", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "leaderboard")
		})
		leaderboardGroup.POST("/getUserRank", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "leaderboard")
		})
	}

	// 游戏相关路由
	gameGroup := r.Group("/game")
	{
		gameGroup.POST("/joinGame", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "game")
		})
		gameGroup.POST("/leaveGame", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "game")
		})
		gameGroup.POST("/move", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "game")
		})
		gameGroup.POST("/getGameState", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "game")
		})
	}

	// 好友相关路由
	friendsGroup := r.Group("/friends")
	{
		friendsGroup.POST("/addFriend", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "friends")
		})
		friendsGroup.POST("/removeFriend", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "friends")
		})
		friendsGroup.POST("/getFriends", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "friends")
		})
		friendsGroup.POST("/sendFriendRequest", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "friends")
		})
		friendsGroup.POST("/respondFriendRequest", func(c *gin.Context) {
			h.usecase.ForwardRequest(c, "friends")
		})
	}
}