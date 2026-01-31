#!/bin/bash

# 启动贪吃蛇游戏微服务系统
echo "Starting Snake Game Microservices..."

# 检查二进制文件是否存在
if [ ! -f "bin/gateway_server" ]; then
    echo "Building gateway service..."
    go build -o bin/gateway_server gateway/main.go
fi

if [ ! -f "bin/lobby_server" ]; then
    echo "Building lobby service..."
    go build -o bin/lobby_server lobby/server.go
fi

if [ ! -f "bin/matching_server" ]; then
    echo "Building matching service..."
    go build -o bin/matching_server matching/server.go
fi

if [ ! -f "bin/room_server" ]; then
    echo "Building room service..."
    go build -o bin/room_server room/server.go
fi

if [ ! -f "bin/leaderboard_server" ]; then
    echo "Building leaderboard service..."
    go build -o bin/leaderboard_server leaderboard/server.go
fi

if [ ! -f "bin/game_server" ]; then
    echo "Building game service..."
    go build -o bin/game_server game/server.go
fi

if [ ! -f "bin/friends_server" ]; then
    echo "Building friends service..."
    go build -o bin/friends_server friends/server.go
fi

# 设置 MongoDB 连接字符串
export MONGODB_URI="mongodb://localhost:27017"

# 启动服务
echo "Starting Gateway Service on port 8080..."
./bin/gateway_server &
GATEWAY_PID=$!

echo "Starting Lobby Service on port 50051..."
./bin/lobby_server &
LOBBY_PID=$!

echo "Starting Matching Service on port 50052..."
./bin/matching_server &
MATCHING_PID=$!

echo "Starting Room Service on port 50053..."
./bin/room_server &
ROOM_PID=$!

echo "Starting Leaderboard Service on port 50054..."
./bin/leaderboard_server &
LEADERBOARD_PID=$!

echo "Starting Game Service on port 50055..."
./bin/game_server &
GAME_PID=$!

echo "Starting Friends Service on port 50056..."
./bin/friends_server &
FRIENDS_PID=$!

# 保存 PID 到文件以便后续停止服务
echo "$GATEWAY_PID $LOBBY_PID $MATCHING_PID $ROOM_PID $LEADERBOARD_PID $GAME_PID $FRIENDS_PID" > service_pids.txt

echo "All services started successfully!"
echo "Gateway: http://localhost:8080"
echo "Lobby: localhost:50051"
echo "Matching: localhost:50052"
echo "Room: localhost:50053"
echo "Leaderboard: localhost:50054"
echo "Game: localhost:50055"
echo "Friends: localhost:50056"

# 等待所有后台进程
wait