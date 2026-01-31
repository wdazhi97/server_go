import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import { leaderboardService } from '../utils/api';
import Head from 'next/head';

const Leaderboard = () => {
  const [leaderboard, setLeaderboard] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    loadLeaderboard();
  }, []);

  const loadLeaderboard = async () => {
    try {
      setLoading(true);
      const response = await leaderboardService.getLeaderboard(50, 0); // 获取前50名
      if (response.success) {
        setLeaderboard(response.entries || []);
      }
    } catch (error) {
      console.error('Error loading leaderboard:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
      <Head>
        <title>排行榜 - 贪吃蛇游戏</title>
      </Head>

      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
        <h1>🏆 排行榜</h1>
        <button className="btn" onClick={() => router.push('/')}>返回首页</button>
      </header>

      {loading ? (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <p>加载中...</p>
        </div>
      ) : (
        <div>
          <table style={{ width: '100%', borderCollapse: 'collapse', backgroundColor: '#1a1a1a', borderRadius: '8px', overflow: 'hidden' }}>
            <thead>
              <tr style={{ backgroundColor: '#222' }}>
                <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid #333' }}>排名</th>
                <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid #333' }}>玩家</th>
                <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid #333' }}>积分</th>
                <th style={{ padding: '12px', textAlign: 'left', borderBottom: '1px solid #333' }}>胜场</th>
              </tr>
            </thead>
            <tbody>
              {leaderboard.length > 0 ? (
                leaderboard.map((entry, index) => (
                  <tr 
                    key={index} 
                    style={{ 
                      backgroundColor: index < 3 ? '#222' : 'transparent',
                      borderLeft: index < 3 ? (index === 0 ? '4px solid gold' : index === 1 ? '4px solid silver' : '4px solid #cd7f32') : 'none'
                    }}
                  >
                    <td style={{ padding: '12px', borderBottom: '1px solid #333' }}>
                      <span style={{ 
                        display: 'inline-block', 
                        width: '24px', 
                        height: '24px', 
                        lineHeight: '24px', 
                        textAlign: 'center', 
                        borderRadius: '50%', 
                        backgroundColor: index === 0 ? '#FFD700' : index === 1 ? '#C0C0C0' : index === 2 ? '#CD7F32' : '#333',
                        color: index < 3 ? '#000' : '#fff'
                      }}>
                        {entry.rank}
                      </span>
                    </td>
                    <td style={{ padding: '12px', borderBottom: '1px solid #333' }}>
                      {entry.username}
                      {index === 0 && <span style={{ marginLeft: '8px', color: '#FFD700' }}>👑</span>}
                    </td>
                    <td style={{ padding: '12px', borderBottom: '1px solid #333' }}>{entry.score}</td>
                    <td style={{ padding: '12px', borderBottom: '1px solid #333' }}>
                      {entry.wins || '-'}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4} style={{ padding: '20px', textAlign: 'center', color: '#777' }}>
                    暂无排行榜数据
                  </td>
                </tr>
              )}
            </tbody>
          </table>

          <div style={{ marginTop: '30px', textAlign: 'center' }}>
            <button 
              className="btn" 
              onClick={loadLeaderboard}
              disabled={loading}
            >
              {loading ? '刷新中...' : '刷新排行榜'}
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Leaderboard;