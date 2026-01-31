# 贪吃蛇游戏微服务系统

这是一个基于微服务架构的贪吃蛇游戏全栈系统，采用 Go 语言开发后端服务，React + Next.js 开发前端界面，MongoDB 作为数据库，支持高并发游戏体验。系统包含API网关、服务发现机制和6个核心业务服务。

## 系统架构

本系统由以下组件组成：

### 1. API网关服务 (Gateway Service) - 端口 8080
- 统一入口，请求路由
- 服务发现与负载均衡
- 认证与鉴权
- 限流与熔断

### 2. 大厅服务 (Lobby Service) - 内部服务
- 用户注册和登录
- 身份验证和会话管理
- 基础用户信息管理

### 3. 匹配服务 (Matching Service) - 内部服务
- 搜索在线玩家
- 匹配玩家进入同一游戏
- 匹配策略管理

### 4. 房间服务 (Room Service) - 内部服务
- 管理游戏房间
- 房间内实时通信（IM功能）
- 房间状态管理

### 5. 排行榜服务 (Leaderboard Service) - 内部服务
- 分数记录和排名
- 排行榜查询
- 历史战绩管理

### 6. 游戏服务 (Game Service) - 内部服务
- 游戏状态同步
- 实时游戏逻辑处理
- 蛇的位置和状态管理

### 7. 好友服务 (Friends Service) - 内部服务
- 好友关系管理
- 好友请求处理
- 好友状态查询

## 技术栈

- **后端**: Go
- **前端**: React + Next.js
- **数据库**: MongoDB
- **协议**: gRPC (服务间), HTTP/REST (网关到前端)
- **服务发现**: 内置注册中心

## 快速开始

### 1. 环境准备

确保系统已安装以下软件：
- Go 1.25+
- Node.js 16+
- MongoDB
- Docker (可选，用于容器化部署)
- Docker Compose (可选)

### 2. 后端服务构建

```bash
# 构建所有后端服务（包括网关）
./build.sh
```

### 3. 前端项目构建

```bash
cd frontend
npm install
npm run dev  # 开发模式
# 或
npm run build && npm start  # 生产模式
```

### 4. 启动服务

#### 本地启动

```bash
# 启动 MongoDB (如果尚未运行)
mongod

# 启动后端所有服务（包括网关）
./start_services.sh

# 在另一个终端启动前端开发服务器
cd frontend
npm run dev
```

#### Docker 启动

```bash
# 使用 Docker Compose 启动所有服务（包括网关）
docker-compose up --build
```

### 5. 服务端口

- 前端: `:3000`
- API网关: `:8080`
- 内部服务端口（仅容器内访问）:
  - 大厅服务: `:50051`
  - 匹配服务: `:50052`
  - 房间服务: `:50053`
  - 排行榜服务: `:50054`
  - 游戏服务: `:50055`
  - 好友服务: `:50056`

## 项目结构

```
snake-game/
├── gateway/                  # API网关服务
│   └── main.go               # 网关主程序
├── lobby/                    # 大厅服务
├── matching/                 # 匹配服务
├── room/                     # 房间服务
├── leaderboard/              # 排行榜服务
├── game/                     # 游戏服务
├── friends/                  # 好友服务
├── mongodb/                  # 数据库模型和操作
│   ├── models.go             # 数据模型定义
│   └── connection.go         # 数据库连接
├── proto/                    # gRPC 协议定义
│   ├── snake_game.proto      # 协议定义文件
│   ├── snake_game.pb.go      # 生成的 Go 代码
│   └── snake_game_grpc.pb.go # 生成的 gRPC 代码
├── frontend/                 # 前端项目 (React + Next.js)
│   ├── public/               # 静态资源
│   ├── src/
│   │   ├── components/       # 可复用组件
│   │   ├── pages/            # 页面组件
│   │   ├── styles/           # 样式文件
│   │   └── utils/            # 工具函数
│   ├── package.json          # 前端依赖配置
│   ├── next.config.js        # Next.js 配置
│   └── README.md             # 前端项目文档
├── bin/                      # 编译后的二进制文件
├── docker-compose.yml        # Docker Compose 配置
├── Dockerfile                # Docker 镜像配置
├── build.sh                  # 构建脚本
├── start_services.sh         # 启动服务脚本
├── SYSTEM_ARCHITECTURE.md    # 系统架构文档
├── GATEWAY_SERVICE.md        # API网关服务文档
├── FRONTEND_INTEGRATION.md   # 前后端集成说明
├── PROJECT_SUMMARY.md        # 项目总结
└── README.md
```

## 功能特性

### 前端功能
- **登录/注册界面**: 用户身份验证 -> `/auth/*` -> 大厅服务
- **游戏大厅**: 用户信息、匹配、好友列表 -> `/match/*` -> 匹配服务
- **房间界面**: 房间管理、IM功能 -> `/room/*` -> 房间服务
- **游戏界面**: 实时状态同步 -> `/game/*` -> 游戏服务
- **排行榜界面**: 排行榜查询 -> `/leaderboard/*` -> 排行榜服务
- **好友管理**: 好友系统 -> `/friends/*` -> 好友服务

### 后端功能
- 完整的微服务架构
- API网关统一入口
- 服务发现与负载均衡
- 数据持久化存储
- 实时通信支持
- 用户认证和授权
- 游戏逻辑处理
- 社交功能支持

## API 文档

- 服务间通信通过 gRPC 实现，具体接口定义见 `proto/snake_game.proto` 文件
- 前后端通信通过 REST API 实现，经由API网关路由到对应服务，具体接口定义见 `frontend/src/utils/api.ts`
- API网关路由规则见 `GATEWAY_SERVICE.md`

## 部署

使用 Docker Compose 进行一体化部署:

```bash
docker-compose up -d
```

## 开发指南

### 后端开发

1. 在 `proto/snake_game.proto` 中定义新的 gRPC 服务或消息
2. 重新生成 Go 代码: `protoc --go_out=. --go-grpc_out=. proto/snake_game.proto`
3. 在相应的服务中实现新功能
4. 如需添加新的API路由，在 `gateway/main.go` 中添加对应的处理函数
5. 重新构建: `./build.sh`

### 前端开发

1. 创建新的页面组件在 `frontend/src/pages/`
2. 创建可复用组件在 `frontend/src/components/`
3. 如需新增 API 调用，在 `frontend/src/utils/api.ts` 中添加相应函数
4. 启动开发服务器: `cd frontend && npm run dev`

### 数据库模型

所有数据库模型定义在 `mongodb/models.go` 中，如需添加新模型，请遵循现有模式。

## 架构特点

1. **高可用性**: 微服务架构支持独立部署和扩展
2. **统一入口**: API网关提供统一访问入口
3. **服务发现**: 自动服务注册与发现
4. **负载均衡**: 网关自动负载均衡
5. **认证授权**: 统一认证与授权机制
6. **前后端分离**: 前端和后端完全解耦，便于独立开发和维护
7. **实时通信**: 支持实时游戏状态同步和聊天功能
8. **数据持久化**: MongoDB 提供可靠的持久化存储
9. **用户系统**: 完整的用户注册、登录、资料管理功能
10. **社交功能**: 好友系统和房间聊天功能
11. **竞技系统**: 排行榜和积分系统

## 项目完成度

- ✅ API网关服务完成（统一入口、路由、服务发现）
- ✅ 后端6个微服务完成
- ✅ 前端React + Next.js界面完成
- ✅ 数据库集成完成
- ✅ 前后端API集成完成（通过网关）
- ✅ 容器化部署方案完成
- ✅ 完整的游戏功能实现

## 未来扩展

- 游戏回放功能
- 观战模式
- 更多游戏模式
- 装备/道具系统
- 社交功能增强
- 移动端适配
- 监控与告警系统

## 维护

- 日志记录和监控
- 性能优化
- 安全加固
- 容量规划