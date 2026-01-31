import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import { 
  matchingService, 
  leaderboardService, 
  friendsService, 
  lobbyService 
} from '../utils/api';
import Head from 'next/head';
import Header from '../components/Header';

const Home = () => {
  const [user, setUser] = useState<any>(null);
  const [onlinePlayers, setOnlinePlayers] = useState<any[]>([]);
  const [waitingPlayers, setWaitingPlayers] = useState<number>(0);
  const [leaderboard, setLeaderboard] = useState<any[]>([]);
  const [friends, setFriends] = useState<any[]>([]);
  const [searchUsername, setSearchUsername] = useState('');
  const [matchingStatus, setMatchingStatus] = useState<string>('');
  const router = useRouter();

  // 检查用户登录状态
  useEffect(() => {
    const userId = localStorage.getItem('userId');
    const username = localStorage.getItem('username');
    
    if (!userId || !username) {
      router.push('/login');
      return;
    }

    setUser({ id: userId, username });

    // 加载数据
    loadUserData(userId);
    loadLeaderboard();
    loadFriends(userId);
  }, []);

  // 加载用户数据
  const loadUserData = async (userId: string) => {
    try {
      // 获取在线玩家
      const onlineResp = await matchingService.getOnlinePlayers();
      if (onlineResp.success) {
        setOnlinePlayers(onlineResp.players || []);
      }

      // 获取等待玩家数
      const waitingResp = await matchingService.getWaitingPlayers();
      setWaitingPlayers(waitingResp.count || 0);

      // 获取用户资料
      const profileResp = await lobbyService.getUserProfile(userId);
      if (profileResp.success) {
        setUser((prev: any) => ({ ...prev, ...profileResp.user, score: profileResp.score }));
      }
    } catch (error) {
      console.error('Error loading user data:', error);
    }
  };

  // 加载排行榜
  const loadLeaderboard = async () => {
    try {
      const response = await leaderboardService.getLeaderboard(10, 0);
      if (response.success) {
        setLeaderboard(response.entries || []);
      }
    } catch (error) {
      console.error('Error loading leaderboard:', error);
    }
  };

  // 加载好友列表
  const loadFriends = async (userId: string) => {
    try {
      const response = await friendsService.getFriends(userId);
      if (response.success) {
        setFriends(response.friends || []);
      }
    } catch (error) {
      console.error('Error loading friends:', error);
    }
  };

  // 开始匹配
  const startMatching = async () => {
    if (!user) return;

    try {
      setMatchingStatus('正在寻找对手...');
      const response = await matchingService.findMatch(
        user.id, 
        user.username, 
        user.score || 0
      );
      
      if (response.success) {
        if (response.room_id) {
          // 匹配成功，跳转到游戏房间
          router.push(`/room/${response.room_id}`);
        } else {
          setMatchingStatus(response.message);
        }
      } else {
        setMatchingStatus(response.message || '匹配失败');
      }
    } catch (error) {
      setMatchingStatus('匹配失败，请重试');
      console.error('Error matching:', error);
    }
  };

  // 搜索并添加好友
  const searchAndAddFriend = async () => {
    if (!searchUsername.trim()) return;

    try {
      // 这里需要根据用户名查找用户ID
      // 为简化实现，我们假设有这样的API
      // 在实际实现中，可能需要一个专门的查找用户API
      alert(`功能开发中：查找并添加好友 ${searchUsername}`);
    } catch (error) {
      console.error('Error adding friend:', error);
    }
  };

  // 退出登录
  const handleLogout = () => {
    localStorage.removeItem('userId');
    localStorage.removeItem('username');
    localStorage.removeItem('token');
    router.push('/login');
  };

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <Head>
        <title>贪吃蛇游戏 - 首页</title>
      </Head>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', margin: '20px auto', maxWidth: '1200px' }}>
        <h1>🐍 贪吃蛇游戏</h1>
        <div>
          <span style={{ marginRight: '20px' }}>欢迎, {user?.username}</span>
          <button className="btn" onClick={handleLogout}>退出</button>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '20px', marginBottom: '30px' }}>
        {/* 用户信息 */}
        <div className="user-info">
          <h2>个人信息</h2>
          <p><strong>用户名:</strong> {user?.username}</p>
          <p><strong>积分:</strong> {user?.score || 0}</p>
          <p><strong>胜场:</strong> {user?.games_won || 0}</p>
          <p><strong>总局数:</strong> {user?.games_played || 0}</p>
        </div>

        {/* 匹配区域 */}
        <div style={{ backgroundColor: '#1a1a1a', padding: '20px', borderRadius: '8px', border: '1px solid #333' }}>
          <h2>匹配游戏</h2>
          <p>当前等待玩家: {waitingPlayers}</p>
          <div style={{ marginTop: '15px' }}>
            <button 
              className="btn btn-success" 
              onClick={startMatching}
              style={{ width: '100%', marginBottom: '10px' }}
            >
              开始匹配
            </button>
            {matchingStatus && (
              <p style={{ color: 'orange', fontSize: '14px' }}>{matchingStatus}</p>
            )}
          </div>
          <button 
            className="btn" 
            onClick={() => router.push('/single-player')}
            style={{ width: '100%' }}
          >
            单人游戏
          </button>
        </div>

        {/* 好友列表 */}
        <div style={{ backgroundColor: '#1a1a1a', padding: '20px', borderRadius: '8px', border: '1px solid #333' }}>
          <h2>好友</h2>
          <div style={{ marginBottom: '15px' }}>
            <input
              type="text"
              placeholder="输入用户名添加好友"
              value={searchUsername}
              onChange={(e) => setSearchUsername(e.target.value)}
              style={{
                padding: '8px',
                marginRight: '10px',
                width: 'calc(100% - 80px)',
                backgroundColor: '#222',
                border: '1px solid #444',
                borderRadius: '4px',
                color: 'white'
              }}
            />
            <button className="btn" onClick={searchAndAddFriend}>添加</button>
          </div>
          <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
            {friends.length > 0 ? (
              friends.map((friend, index) => (
                <div key={index} style={{ padding: '8px', borderBottom: '1px solid #333' }}>
                  <span style={{ color: friend.online ? '#4CAF50' : '#aaa' }}>
                    {friend.username} {friend.online ? '🟢' : '🔴'}
                  </span>
                </div>
              ))
            ) : (
              <p style={{ color: '#777', fontStyle: 'italic' }}>暂无好友</p>
            )}
          </div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
        {/* 在线玩家 */}
        <div style={{ backgroundColor: '#1a1a1a', padding: '20px', borderRadius: '8px', border: '1px solid #333' }}>
          <h2>在线玩家 ({onlinePlayers.length})</h2>
          <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
            {onlinePlayers.length > 0 ? (
              onlinePlayers.map((player, index) => (
                <div key={index} style={{ padding: '8px', borderBottom: '1px solid #333' }}>
                  <span>{player.username} (积分: {player.rating || 0})</span>
                </div>
              ))
            ) : (
              <p style={{ color: '#777', fontStyle: 'italic' }}>暂无在线玩家</p>
            )}
          </div>
        </div>

        {/* 排行榜 */}
        <div style={{ backgroundColor: '#1a1a1a', padding: '20px', borderRadius: '8px', border: '1px solid #333' }}>
          <h2>排行榜</h2>
          <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
            {leaderboard.length > 0 ? (
              leaderboard.map((entry, index) => (
                <div 
                  key={index} 
                  style={{ 
                    padding: '8px', 
                    borderBottom: '1px solid #333',
                    backgroundColor: entry.user_id === user?.id ? '#333' : 'transparent'
                  }}
                >
                  <span>
                    #{entry.rank}. {entry.username} - {entry.score} 分
                  </span>
                </div>
              ))
            ) : (
              <p style={{ color: '#777', fontStyle: 'italic' }}>排行榜为空</p>
            )}
          </div>
        </div>
      </div>

      <div style={{ marginTop: '30px', textAlign: 'center' }}>
        <button 
          className="btn btn-success" 
          onClick={() => router.push('/leaderboard')}
          style={{ marginRight: '10px' }}
        >
          查看完整排行榜
        </button>
        <button 
          className="btn" 
          onClick={() => router.push('/profile')}
        >
          个人资料
        </button>
      </div>
    </div>
  );
};

export default Home;