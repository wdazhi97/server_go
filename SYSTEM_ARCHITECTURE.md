# 贪吃蛇游戏微服务系统架构

## 概述

这是一个基于微服务架构的贪吃蛇游戏后端系统，采用 Go 语言开发，MongoDB 作为数据库，支持高并发游戏体验。系统包含API网关、服务发现机制和6个核心业务服务。

## 系统架构

### 1. API网关服务 (Gateway Service) - 端口 8080
- **功能**: 统一入口，请求路由，认证与鉴权，限流与熔断
- **API**:
  - `/auth/*`: 认证相关路由（转发到大厅服务）
  - `/match/*`: 匹配相关路由（转发到匹配服务）
  - `/room/*`: 房间相关路由（转发到房间服务）
  - `/leaderboard/*`: 排行榜相关路由（转发到排行榜服务）
  - `/game/*`: 游戏相关路由（转发到游戏服务）
  - `/friends/*`: 好友相关路由（转发到好友服务）
- **特性**:
  - 服务发现与自动注册
  - 负载均衡
  - 认证中间件
  - 限流与熔断机制
  - 日志记录与监控

### 2. 服务注册中心 (Service Registry)
- **功能**: 服务注册与发现，健康检查
- **API**:
  - `RegisterService`: 服务注册
  - `UnregisterService`: 服务注销
  - `DiscoverService`: 服务发现
  - `HealthCheck`: 健康检查

### 3. 大厅服务 (Lobby Service) - 内部服务
- **功能**: 用户注册和登录
- **API**:
  - `Register`: 用户注册
  - `Login`: 用户登录
  - `GetUserProfile`: 获取用户资料
  - `Logout`: 用户登出
- **数据模型**: User, Leaderboard

### 4. 匹配服务 (Matching Service) - 内部服务
- **功能**: 搜索在线玩家，匹配进入同一游戏
- **API**:
  - `FindMatch`: 寻找匹配的玩家
  - `CancelMatch`: 取消匹配
  - `GetWaitingPlayers`: 获取等待匹配的玩家数量
  - `GetOnlinePlayers`: 获取在线玩家列表
- **数据模型**: User

### 5. 房间服务 (Room Service) - 内部服务
- **功能**: 管理游戏房间，提供实时通信(IM功能)
- **API**:
  - `CreateRoom`: 创建房间
  - `JoinRoom`: 加入房间
  - `LeaveRoom`: 离开房间
  - `SendMessage`: 发送消息
  - `GetRoomMessages`: 获取房间消息
  - `StartGame`: 开始游戏
- **数据模型**: GameRoom, Message

### 6. 排行榜服务 (Leaderboard Service) - 内部服务
- **功能**: 分数记录和排名
- **API**:
  - `GetLeaderboard`: 获取排行榜
  - `UpdateScore`: 更新分数
  - `GetUserRank`: 获取用户排名
- **数据模型**: Leaderboard

### 7. 游戏服务 (Game Service) - 内部服务
- **功能**: 游戏状态同步，实时游戏逻辑处理
- **API**:
  - `JoinGame`: 加入游戏
  - `LeaveGame`: 离开游戏
  - `Move`: 移动指令
  - `GetGameState`: 获取游戏状态
  - `SubscribeGameUpdates`: 订阅游戏状态更新(流)
- **数据模型**: GameState

### 8. 好友服务 (Friends Service) - 内部服务
- **功能**: 好友关系管理
- **API**:
  - `AddFriend`: 添加好友
  - `RemoveFriend`: 删除好友
  - `GetFriends`: 获取好友列表
  - `SendFriendRequest`: 发送好友请求
  - `RespondFriendRequest`: 回应好友请求
- **数据模型**: Friend

## 通信方式

- 服务间通信：gRPC + 服务发现
- 前后端通信：REST API（通过API网关统一访问）
- 实时通信：WebSocket（通过网关代理）

## 服务发现机制

使用内置服务注册中心实现服务发现：
- 服务启动时自动向注册中心注册
- API网关从注册中心获取服务列表
- 自动负载均衡和服务健康检查
- 服务故障自动剔除和恢复

## 数据库模型

### User (用户)
- ID, Username, Password, Email, CreatedAt, UpdatedAt, Online, LastSeen

### Friend (好友关系)
- ID, UserID, FriendID, Status(pending/accepted/blocked), CreatedAt

### GameRoom (游戏房间)
- ID, RoomName, CreatorID, Players, MaxPlayers, Status, CreatedAt, GameOptions

### GameOptions (游戏选项)
- FoodCount, WallEnabled, Speed

### GameRecord (游戏记录)
- ID, RoomID, Players, WinnerID, Scores, GameTime, CreatedAt

### PlayerGameResult (玩家游戏结果)
- PlayerID, Score, Rank

### Leaderboard (排行榜)
- ID, UserID, Score, GamesWon, GamesPlayed, UpdatedAt

### Message (房间消息)
- ID, RoomID, SenderID, Content, Type, CreatedAt

### GameSnake (游戏中蛇)
- ID, Points, Color, Length, Score

### Point (坐标点)
- X, Y

### GameState (游戏状态)
- ID, RoomID, Snakes, Foods, Walls, Status, UpdatedAt

## 部署

### 本地部署
1. 确保 MongoDB 已安装并运行
2. 运行 `./build.sh` 构建所有服务
3. 运行 `./start_services.sh` 启动所有服务（包括网关）

### Docker 部署
1. 运行 `docker-compose up --build` 启动所有服务（包括网关）

## 前端对接

前端需要使用 React + Next.js 实现，与后端服务通过API网关进行交互：

1. **注册登录界面** -> `/auth/*` -> 大厅服务
2. **单人游戏/匹配/组队入口** -> `/match/*` -> 匹配服务
3. **房间界面(IM功能)** -> `/room/*` -> 房间服务
4. **游戏界面** -> `/game/*` -> 游戏服务
5. **排行榜界面** -> `/leaderboard/*` -> 排行榜服务
6. **好友管理** -> `/friends/*` -> 好友服务

## 技术栈

- **后端**: Go
- **数据库**: MongoDB
- **协议**: gRPC (服务间), HTTP/REST (网关到前端)
- **前端**: React + Next.js
- **容器化**: Docker + Docker Compose
- **服务发现**: 内置注册中心

## 特性

1. **高可用性**: 微服务架构支持独立部署和扩展
2. **统一入口**: API网关提供统一访问入口
3. **服务发现**: 自动服务注册与发现
4. **负载均衡**: 网关自动负载均衡
5. **认证授权**: 统一认证与授权机制
6. **实时通信**: WebSocket 和 gRPC 流支持实时游戏状态同步
7. **数据持久化**: MongoDB 提供可靠的持久化存储
8. **用户系统**: 完整的用户注册、登录、资料管理功能
9. **社交功能**: 好友系统和房间聊天功能
10. **竞技系统**: 排行榜和积分系统
11. **监控与运维**: 统一日志记录和健康检查