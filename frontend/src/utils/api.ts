// API 工具函数，用于与后端微服务通信

import axios from 'axios';

// API 配置 - 通过nginx反向代理访问API网关
// 生产环境使用相对路径 /api，开发环境可以直接访问本地gateway
const API_CONFIG = {
  GATEWAY_URL: process.env.NEXT_PUBLIC_GATEWAY_URL || '/api',
};

// Axios 实例
const createAxiosInstance = (baseURL: string) => {
  const instance = axios.create({
    baseURL,
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // 请求拦截器
  instance.interceptors.request.use(
    (config) => {
      // 可以在这里添加认证 token
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // 响应拦截器
  instance.interceptors.response.use(
    (response) => {
      return response;
    },
    (error) => {
      console.error('API Error:', error);
      return Promise.reject(error);
    }
  );

  return instance;
};

// API网关实例
const gatewayApi = createAxiosInstance(API_CONFIG.GATEWAY_URL);

export const lobbyService = {
  // 用户注册
  register: async (username: string, password: string, email: string) => {
    try {
      const response = await gatewayApi.post('/auth/register', {
        username,
        password,
        email,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 用户登录
  login: async (username: string, password: string) => {
    try {
      const response = await gatewayApi.post('/auth/login', {
        username,
        password,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 获取用户资料
  getUserProfile: async (userId: string) => {
    try {
      const response = await gatewayApi.post('/auth/getUserProfile', {
        userId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 用户登出
  logout: async (userId: string) => {
    try {
      const response = await gatewayApi.post('/auth/logout', {
        userId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },
};

export const matchingService = {
  // 寻找匹配
  findMatch: async (playerId: string, username: string, rating: number) => {
    try {
      const response = await gatewayApi.post('/match/findMatch', {
        playerId,
        username,
        rating,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 取消匹配
  cancelMatch: async (playerId: string) => {
    try {
      const response = await gatewayApi.post('/match/cancelMatch', {
        playerId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 获取等待玩家数量
  getWaitingPlayers: async () => {
    try {
      const response = await gatewayApi.post('/match/getWaitingPlayers', {});
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 获取在线玩家
  getOnlinePlayers: async () => {
    try {
      const response = await gatewayApi.post('/match/getOnlinePlayers', {});
      return response.data;
    } catch (error) {
      throw error;
    }
  },
};

export const roomService = {
  // 创建房间
  createRoom: async (userId: string, roomName: string, maxPlayers: number) => {
    try {
      const response = await gatewayApi.post('/room/createRoom', {
        userId,
        roomName,
        maxPlayers,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 加入房间
  joinRoom: async (roomId: string, userId: string) => {
    try {
      const response = await gatewayApi.post('/room/joinRoom', {
        roomId,
        userId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 离开房间
  leaveRoom: async (roomId: string, userId: string) => {
    try {
      const response = await gatewayApi.post('/room/leaveRoom', {
        roomId,
        userId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 发送消息
  sendMessage: async (roomId: string, senderId: string, content: string, type: string = 'text') => {
    try {
      const response = await gatewayApi.post('/room/sendMessage', {
        roomId,
        senderId,
        content,
        type,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 获取房间消息
  getRoomMessages: async (roomId: string, limit?: number, since?: number) => {
    try {
      const response = await gatewayApi.post('/room/getRoomMessages', {
        roomId,
        limit,
        since,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 开始游戏
  startGame: async (roomId: string) => {
    try {
      const response = await gatewayApi.post('/room/startGame', {
        roomId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },
};

export const leaderboardService = {
  // 获取排行榜
  getLeaderboard: async (limit: number = 10, offset: number = 0) => {
    try {
      const response = await gatewayApi.post('/leaderboard/getLeaderboard', {
        limit,
        offset,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 更新分数
  updateScore: async (userId: string, score: number, gameWon: boolean) => {
    try {
      const response = await gatewayApi.post('/leaderboard/updateScore', {
        userId,
        score,
        gameWon,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 获取用户排名
  getUserRank: async (userId: string) => {
    try {
      const response = await gatewayApi.post('/leaderboard/getUserRank', {
        userId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },
};

export const gameService = {
  // 加入游戏
  joinGame: async (roomId: string, playerId: string) => {
    try {
      const response = await gatewayApi.post('/game/joinGame', {
        roomId,
        playerId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 离开游戏
  leaveGame: async (roomId: string, playerId: string) => {
    try {
      const response = await gatewayApi.post('/game/leaveGame', {
        roomId,
        playerId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 移动指令
  move: async (roomId: string, playerId: string, direction: string) => {
    try {
      const response = await gatewayApi.post('/game/move', {
        roomId,
        playerId,
        direction,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 获取游戏状态
  getGameState: async (roomId: string) => {
    try {
      const response = await gatewayApi.post('/game/getGameState', {
        roomId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },
};

export const friendsService = {
  // 添加好友
  addFriend: async (userId: string, friendUsername: string) => {
    try {
      const response = await gatewayApi.post('/friends/addFriend', {
        userId,
        friendUsername,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 删除好友
  removeFriend: async (userId: string, friendUserId: string) => {
    try {
      const response = await gatewayApi.post('/friends/removeFriend', {
        userId,
        friendUserId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 获取好友列表
  getFriends: async (userId: string) => {
    try {
      const response = await gatewayApi.post('/friends/getFriends', {
        userId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 发送好友请求
  sendFriendRequest: async (userId: string, targetUserId: string) => {
    try {
      const response = await gatewayApi.post('/friends/sendFriendRequest', {
        userId,
        targetUserId,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 回应好友请求
  respondFriendRequest: async (userId: string, requestUserId: string, accepted: boolean) => {
    try {
      const response = await gatewayApi.post('/friends/respondFriendRequest', {
        userId,
        requestUserId,
        accepted,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },
};