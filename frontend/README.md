# 贪吃蛇游戏前端

这是一个基于 React + Next.js 的贪吃蛇游戏前端项目，与后端微服务系统配合使用。

## 项目概述

前端项目包含以下功能模块：

### 页面
- **首页** - 游戏大厅，包含用户信息、匹配、好友列表等
- **登录/注册** - 用户身份验证
- **游戏房间** - 多人游戏房间及聊天功能
- **排行榜** - 显示玩家积分排名
- **单人游戏** - 离线单人贪吃蛇游戏

### 技术栈
- **框架**: Next.js 14
- **语言**: TypeScript
- **状态管理**: Redux (可选)
- **网络请求**: Axios
- **实时通信**: Socket.IO Client (预留)

## 快速开始

### 环境要求
- Node.js 16+
- npm 或 yarn

### 安装依赖
```bash
npm install
# 或
yarn install
```

### 环境配置
复制 `.env.example` 文件为 `.env.local` 并根据实际情况修改API网关地址：

```bash
cp .env.example .env.local
```

注意：前端现在通过API网关（默认端口8080）统一访问所有后端服务，而不是直接访问各个微服务。

### 开发模式
```bash
npm run dev
# 或
yarn dev
```

应用将在 `http://localhost:3000` 上运行。

### 构建生产版本
```bash
npm run build
# 然后
npm start
```

## 项目结构

```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── components/         # 可复用组件
│   │   ├── GameBoard.tsx   # 游戏面板组件
│   │   ├── Header.tsx      # 页面头部组件
│   │   └── Layout.tsx      # 页面布局组件
│   ├── pages/              # 页面组件
│   │   ├── index.tsx       # 首页
│   │   ├── login.tsx       # 登录/注册页
│   │   ├── leaderboard.tsx # 排行榜页
│   │   ├── single-player.tsx # 单人游戏页
│   │   └── room/[id].tsx   # 游戏房间页
│   ├── styles/             # 样式文件
│   │   └── globals.css     # 全局样式
│   └── utils/              # 工具函数
│       └── api.ts          # API 调用工具
├── .env.example            # 环境变量示例
├── next.config.js          # Next.js 配置
├── tsconfig.json           # TypeScript 配置
└── package.json            # 项目配置
```

## API 集成

前端通过 `src/utils/api.ts` 文件与后端微服务通信，所有请求都通过API网关统一处理，封装了以下服务的 API 调用：

- **大厅服务** - 用户注册、登录、资料管理 (/auth/*)
- **匹配服务** - 玩家匹配、在线玩家查询 (/match/*)
- **房间服务** - 房间管理、消息发送 (/room/*)
- **排行榜服务** - 排行榜查询、分数更新 (/leaderboard/*)
- **游戏服务** - 游戏状态同步 (/game/*)
- **好友服务** - 好友管理 (/friends/*)

## 主要功能

### 用户系统
- 用户注册和登录
- 个人资料查看和管理

### 游戏匹配
- 快速匹配其他玩家
- 查看等待匹配的玩家数

### 多人游戏
- 创建和加入游戏房间
- 房间内实时聊天功能
- 多人同场竞技

### 单人游戏
- 离线单人贪吃蛇游戏
- 方向键控制
- 得分系统

### 社交功能
- 好友添加和管理
- 好友在线状态显示

### 排行榜
- 实时积分排名
- 个人排名查询

## 部署

### 静态导出
```bash
npm run build
npm run export  # 生成静态文件到 out/ 目录
```

### 服务端渲染部署
构建后可部署到任何支持 Node.js 的服务器或云平台（Vercel、Netlify 等）。

## 开发规范

### 代码风格
- 使用 TypeScript 编写
- 遵循 ESLint 和 Prettier 规范
- 组件命名使用帕斯卡命名法
- 函数和变量使用驼峰命名法

### 文件组织
- 页面组件放在 `pages/` 目录
- 可复用组件放在 `components/` 目录
- 工具函数放在 `utils/` 目录
- 样式文件放在 `styles/` 目录

## API 通信

所有后端服务调用都通过 `api.ts` 中的工具函数进行，这些函数会自动处理：
- 通过API网关路由到相应后端服务
- 错误处理
- 认证信息附加
- 请求/响应拦截

## 未来扩展

- 添加游戏回放功能
- 实现观战模式
- 增加更多游戏模式
- 添加音效和视觉效果
- 实现成就系统
- 添加聊天表情和互动功能