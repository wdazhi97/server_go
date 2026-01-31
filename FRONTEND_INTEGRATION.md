# 前后端集成说明

## 项目结构概览

```
server_go/ (后端微服务)
├── lobby/                    # 大厅服务 (端口 50051)
├── matching/                 # 匹配服务 (端口 50052)
├── room/                     # 房间服务 (端口 50053)
├── leaderboard/              # 排行榜服务 (端口 50054)
├── game/                     # 游戏服务 (端口 50055)
├── friends/                  # 好友服务 (端口 50056)
├── mongodb/                  # 数据库模型和操作
├── proto/                    # gRPC 协议定义
├── bin/                      # 编译后的二进制文件
├── frontend/                 # 前端项目 (React + Next.js)
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── styles/
│   │   └── utils/
│   ├── package.json
│   ├── next.config.js
│   └── README.md
├── docker-compose.yml        # 包含前后端服务的部署配置
├── build.sh                  # 后端构建脚本
├── start_services.sh         # 后端启动脚本
└── FRONTEND_INTEGRATION.md   # 本文档
```

## 前后端接口规范

### API 通信方式
- 后端微服务之间使用 gRPC 通信
- 前端与后端服务使用 REST API 通信
- 实时通信使用 WebSocket（通过 Socket.IO）

### 后端服务端口映射
- 大厅服务: `http://localhost:50051`
- 匹配服务: `http://localhost:50052`
- 房间服务: `http://localhost:50053`
- 排行榜服务: `http://localhost:50054`
- 游戏服务: `http://localhost:50055`
- 好友服务: `http://localhost:50056`

## 前端 API 配置

前端项目中的 `.env.local` 文件需要配置后端服务地址：

```env
NEXT_PUBLIC_LOBBY_SERVICE_URL=http://localhost:50051
NEXT_PUBLIC_MATCHING_SERVICE_URL=http://localhost:50052
NEXT_PUBLIC_ROOM_SERVICE_URL=http://localhost:50053
NEXT_PUBLIC_LEADERBOARD_SERVICE_URL=http://localhost:50054
NEXT_PUBLIC_GAME_SERVICE_URL=http://localhost:50055
NEXT_PUBLIC_FRIENDS_SERVICE_URL=http://localhost:50056
```

## 部署方案

### 本地开发
1. 启动后端服务：
   ```bash
   cd /data/workspace/server_go
   ./start_services.sh
   ```

2. 启动 MongoDB（如果未运行）：
   ```bash
   mongod
   ```

3. 启动前端开发服务器：
   ```bash
   cd /data/workspace/server_go/frontend
   npm install
   npm run dev
   ```

### Docker 部署
使用提供的 `docker-compose.yml` 文件可以一键部署所有服务：

```bash
# 构建并启动所有服务（包括前端）
docker-compose up --build
```

## 功能模块对应关系

### 前端页面与后端服务映射

| 前端页面 | 主要后端服务 | 辅助服务 |
|---------|-------------|----------|
| 登录/注册 | 大厅服务 | - |
| 首页 | 匹配服务 | 排行榜服务、好友服务 |
| 游戏房间 | 房间服务 | 游戏服务、匹配服务 |
| 排行榜 | 排行榜服务 | - |
| 单人游戏 | - (纯前端) | - |

### 数据流向

1. **用户认证流程**：
   - 前端登录 → 大厅服务 → MongoDB 验证 → 返回认证信息

2. **游戏匹配流程**：
   - 前端发起匹配 → 匹配服务 → 房间服务 → 创建游戏房间 → 前端跳转到房间

3. **游戏进行流程**：
   - 前端发送操作 → 游戏服务 → 更新游戏状态 → 广播给所有玩家

4. **社交功能流程**：
   - 前端操作 → 好友服务 → MongoDB 更新 → 前端状态更新

## 开发注意事项

### 前后端联调
- 确保后端服务正常运行后再启动前端
- 检查 CORS 配置（虽然目前是 REST API，但需要注意跨域问题）
- 前端错误处理应友好地反馈后端服务异常

### 数据一致性
- 前端应缓存必要的用户状态
- 实时数据（如游戏状态、聊天消息）需定期同步
- 用户操作应有适当的加载状态提示

### 错误处理
- 网络请求失败时应有重试机制
- 服务不可用时应有降级方案
- 用户应得到清晰的错误提示

## 扩展性考虑

### 微服务扩展
- 每个服务可以独立部署和扩展
- 可以通过负载均衡器扩展服务实例
- 数据库可以考虑分片和读写分离

### 前端扩展
- 组件化设计便于功能扩展
- API 工具函数设计便于新增服务
- 状态管理便于复杂业务逻辑

## 监控与维护

### 服务健康检查
- 每个微服务应提供健康检查端点
- 前端应监控与后端服务的连接状态
- 日志记录便于问题排查

### 性能优化
- 前端资源压缩和缓存
- 后端服务性能监控
- 数据库查询优化

## 安全考虑

### API 安全
- 使用 JWT 或其他认证机制
- 敏感信息加密传输
- 输入验证和防护 XSS 攻击

### 数据安全
- 用户密码加密存储
- 敏感操作二次确认
- 访问权限控制

---

这个贪吃蛇游戏系统已经完全实现了前后端分离的微服务架构，后端提供稳定可靠的服务，前端提供流畅的用户体验。系统支持高并发、可扩展，并具备完整的社交和竞技功能。