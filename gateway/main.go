package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"snake-game/gateway/internal/delivery/http"
	"snake-game/gateway/internal/repository"
	"snake-game/gateway/internal/usecase"
	"snake-game/gateway/domain/entity"
	"snake-game/otel"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 初始化 OpenTelemetry
	ctx := context.Background()
	shutdown, err := otel.InitTracer(ctx, "api-gateway")
	if err != nil {
		log.Printf("Warning: Failed to initialize OpenTelemetry: %v", err)
	} else {
		defer shutdown(ctx)
	}

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 添加 OpenTelemetry 中间件
	r.Use(otel.GinMiddleware("api-gateway"))

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Content-Type", "Authorization")
	r.Use(cors.New(config))

	// 初始化仓库层
	serviceRegistry := repository.NewServiceRegistry()

	// 从环境变量读取服务地址（支持k8s服务发现）
	lobbyService := &entity.ServiceInfo{
		Name:    "lobby",
		Address: getEnv("LOBBY_SERVICE_ADDR", "localhost:50051"),
		Health:  true,
	}
	matchingService := &entity.ServiceInfo{
		Name:    "matching",
		Address: getEnv("MATCHING_SERVICE_ADDR", "localhost:50052"),
		Health:  true,
	}
	roomService := &entity.ServiceInfo{
		Name:    "room",
		Address: getEnv("ROOM_SERVICE_ADDR", "localhost:50053"),
		Health:  true,
	}
	leaderboardService := &entity.ServiceInfo{
		Name:    "leaderboard",
		Address: getEnv("LEADERBOARD_SERVICE_ADDR", "localhost:50054"),
		Health:  true,
	}
	gameService := &entity.ServiceInfo{
		Name:    "game",
		Address: getEnv("GAME_SERVICE_ADDR", "localhost:50055"),
		Health:  true,
	}
	friendsService := &entity.ServiceInfo{
		Name:    "friends",
		Address: getEnv("FRIENDS_SERVICE_ADDR", "localhost:50056"),
		Health:  true,
	}

	// 注册服务
	serviceRegistry.RegisterService(nil, lobbyService)
	serviceRegistry.RegisterService(nil, matchingService)
	serviceRegistry.RegisterService(nil, roomService)
	serviceRegistry.RegisterService(nil, leaderboardService)
	serviceRegistry.RegisterService(nil, gameService)
	serviceRegistry.RegisterService(nil, friendsService)

	// 初始化业务逻辑层
	apiGatewayUsecase := usecase.NewAPIGatewayUsecase(serviceRegistry)

	// 初始化通信层
	gatewayHandler := http.NewGatewayHandler(apiGatewayUsecase)

	// 设置路由
	gatewayHandler.SetupRoutes(r)

	// 启动服务器
	port := ":8080"
	log.Printf("API Gateway starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}