# 根 Earthfile：公共 base + 各服务构建，上下文为仓库根
# 执行: earthly +all  或 earthly +gateway-docker 等
VERSION 0.8

# ---------- 各服务 build 阶段（共享 go.mod，上下文为根） ----------

gateway-build:
  FROM golang:1.23-bookworm
  WORKDIR /app
  COPY go.mod go.sum ./
  ENV GOTOOLCHAIN=auto
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gateway_server ./gateway/main.go
  SAVE ARTIFACT gateway_server

lobby-build:
  FROM golang:1.23-bookworm
  WORKDIR /app
  COPY go.mod go.sum ./
  ENV GOTOOLCHAIN=auto
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o lobby_server ./lobby/server.go
  SAVE ARTIFACT lobby_server

matching-build:
  FROM golang:1.23-bookworm
  WORKDIR /app
  COPY go.mod go.sum ./
  ENV GOTOOLCHAIN=auto
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o matching_server ./matching/server.go
  SAVE ARTIFACT matching_server

room-build:
  FROM golang:1.23-bookworm
  WORKDIR /app
  COPY go.mod go.sum ./
  ENV GOTOOLCHAIN=auto
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o room_server ./room/server.go
  SAVE ARTIFACT room_server

leaderboard-build:
  FROM golang:1.23-bookworm
  WORKDIR /app
  COPY go.mod go.sum ./
  ENV GOTOOLCHAIN=auto
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/leaderboard_server ./leaderboard/server.go
  SAVE ARTIFACT bin/leaderboard_server

game-build:
  FROM golang:1.23-bookworm
  WORKDIR /app
  COPY go.mod go.sum ./
  ENV GOTOOLCHAIN=auto
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o game_server ./game/server.go
  SAVE ARTIFACT game_server

friends-build:
  FROM golang:1.23-bookworm
  WORKDIR /app
  COPY go.mod go.sum ./
  ENV GOTOOLCHAIN=auto
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o friends_server ./friends/server.go
  SAVE ARTIFACT friends_server

# ---------- 各服务 docker 阶段（依赖 buildtool + 对应 build） ----------

gateway-docker:
  BUILD +gateway-build
  BUILD ./buildtool+runtime
  FROM ./buildtool+runtime
  COPY +gateway-build/gateway_server ./bin/gateway_server
  RUN chmod +x ./bin/*
  SAVE IMAGE snake-game-gateway:latest

lobby-docker:
  BUILD +lobby-build
  BUILD ./buildtool+runtime
  FROM ./buildtool+runtime
  COPY +lobby-build/bin/lobby_server ./bin/
  RUN chmod +x ./bin/*
  SAVE IMAGE snake-game-lobby:latest

matching-docker:
  BUILD +matching-build
  BUILD ./buildtool+runtime
  FROM ./buildtool+runtime
  COPY +matching-build/matching_server ./bin/matching_server
  RUN chmod +x ./bin/*
  SAVE IMAGE snake-game-matching:latest

room-docker:
  BUILD +room-build
  BUILD ./buildtool+runtime
  FROM ./buildtool+runtime
  COPY +room-build/room_server ./bin/room_server
  RUN chmod +x ./bin/*
  SAVE IMAGE snake-game-room:latest

leaderboard-docker:
  BUILD +leaderboard-build
  BUILD ./buildtool+runtime
  FROM ./buildtool+runtime
  COPY +leaderboard-build/leaderboard_server ./bin/leaderboard_server
  RUN chmod +x ./bin/*
  SAVE IMAGE snake-game-leaderboard:latest

game-docker:
  BUILD +game-build
  BUILD ./buildtool+runtime
  FROM ./buildtool+runtime
  COPY +game-build/game_server ./bin/game_server
  RUN chmod +x ./bin/*
  SAVE IMAGE snake-game-game:latest

friends-docker:
  BUILD +friends-build
  BUILD ./buildtool+runtime
  FROM ./buildtool+runtime
  COPY +friends-build/friends_server ./bin/friends_server
  RUN chmod +x ./bin/*
  SAVE IMAGE snake-game-friends:latest

# ---------- 前端 ----------

frontend-build:
  FROM ubuntu:24.04
  RUN apt-get update && apt-get install -y --no-install-recommends curl && \
      curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
      apt-get install -y nodejs && rm -rf /var/lib/apt/lists/*
  WORKDIR /app
  COPY frontend/package*.json ./
  RUN npm ci --only=production=false
  COPY frontend/ ./
  RUN npm run build
  SAVE ARTIFACT out

frontend-docker:
  FROM ubuntu:24.04
  RUN apt-get update && apt-get install -y --no-install-recommends nginx && rm -rf /var/lib/apt/lists/*
  COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
  COPY +frontend-build/out /usr/share/nginx/html
  CMD ["nginx", "-g", "daemon off;"]
  SAVE IMAGE snake-game-frontend:latest

# ---------- 一键全部 ----------

all:
  BUILD +gateway-docker
  BUILD +lobby-docker
  BUILD +matching-docker
  BUILD +room-docker
  BUILD +leaderboard-docker
  BUILD +game-docker
  BUILD +friends-docker
  BUILD +frontend-docker
