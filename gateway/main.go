package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"snake-game/gateway/internal/delivery/http"
	"snake-game/gateway/internal/repository"
	"snake-game/gateway/internal/usecase"
	"snake-game/gateway/domain/entity"
)

func main() {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Content-Type", "Authorization")
	r.Use(cors.New(config))

	// 初始化仓库层
	serviceRegistry := repository.NewServiceRegistry()

	// 注册内部服务地址
	lobbyService := &entity.ServiceInfo{
		Name:    "lobby",
		Address: "localhost:50051",
		Health:  true,
	}
	matchingService := &entity.ServiceInfo{
		Name:    "matching",
		Address: "localhost:50052",
		Health:  true,
	}
	roomService := &entity.ServiceInfo{
		Name:    "room",
		Address: "localhost:50053",
		Health:  true,
	}
	leaderboardService := &entity.ServiceInfo{
		Name:    "leaderboard",
		Address: "localhost:50054",
		Health:  true,
	}
	gameService := &entity.ServiceInfo{
		Name:    "game",
		Address: "localhost:50055",
		Health:  true,
	}
	friendsService := &entity.ServiceInfo{
		Name:    "friends",
		Address: "localhost:50056",
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