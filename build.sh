#!/bin/bash

# 构建 Snake Game 微服务系统 - 清晰架构版本

set -e  # 遇到错误时退出

echo "Building Snake Game Microservices with Clean Architecture..."

# 确保 Go modules 被下载
echo "Ensuring dependencies..."
go mod tidy

# 为每个服务创建输出目录
mkdir -p bin/

# 清理旧的二进制文件
rm -f bin/*

# 构建各个服务 (使用静态链接以兼容 Alpine)
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

echo "Building lobby service..."
cd lobby
go build -ldflags="-s -w" -o ../bin/lobby_server server.go
cd ..

echo "Building matching service..."
cd matching
go build -ldflags="-s -w" -o ../bin/matching_server server.go
cd ..

echo "Building room service..."
cd room
go build -ldflags="-s -w" -o ../bin/room_server server.go
cd ..

echo "Building leaderboard service..."
cd leaderboard
go build -ldflags="-s -w" -o ../bin/leaderboard_server server.go
cd ..

echo "Building game service..."
cd game
go build -ldflags="-s -w" -o ../bin/game_server server.go
cd ..

echo "Building friends service..."
cd friends
go build -ldflags="-s -w" -o ../bin/friends_server server.go
cd ..

echo "Building gateway service..."
cd gateway
go build -ldflags="-s -w" -o ../bin/gateway_server main.go
cd ..

echo "Build completed successfully!"
echo "Binaries are located in the 'bin/' directory"

echo ""
echo "To run individual services:"
echo "  ./bin/gateway_server      # API Gateway on port 8080"
echo "  ./bin/lobby_server        # Runs on port 50051"
echo "  ./bin/matching_server     # Runs on port 50052"
echo "  ./bin/room_server         # Runs on port 50053"
echo "  ./bin/leaderboard_server  # Runs on port 50054"
echo "  ./bin/game_server         # Runs on port 50055"
echo "  ./bin/friends_server      # Runs on port 50056"
echo ""
echo "Or use docker-compose to run all services:"
echo "  docker-compose up --build"