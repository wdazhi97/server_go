import React, { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/router';
import { roomService, gameService } from '../../utils/api';
import GameBoard from '../../components/GameBoard';
import Head from 'next/head';

const RoomPage = () => {
  const router = useRouter();
  const { id: roomId } = router.query;
  const [messages, setMessages] = useState<any[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [players, setPlayers] = useState<any[]>([]);
  const [roomInfo, setRoomInfo] = useState<any>(null);
  const [gameStarted, setGameStarted] = useState(false);
  const [showGame, setShowGame] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const userId = typeof window !== 'undefined' ? localStorage.getItem('userId') : null;
  const username = typeof window !== 'undefined' ? localStorage.getItem('username') : null;

  // 加载房间信息和消息
  useEffect(() => {
    if (!roomId || !userId) return;

    const loadRoomData = async () => {
      try {
        // 获取房间消息
        const msgResponse = await roomService.getRoomMessages(roomId as string);
        if (msgResponse.success) {
          setMessages(msgResponse.messages || []);
        }
      } catch (error) {
        console.error('Error loading room data:', error);
      }
    };

    loadRoomData();

    // 设置定时器刷新消息
    const msgInterval = setInterval(loadRoomData, 2000);

    return () => clearInterval(msgInterval);
  }, [roomId, userId]);

  // 滚动到底部
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  // 发送消息
  const sendMessage = async () => {
    if (!newMessage.trim() || !roomId || !userId) return;

    try {
      const response = await roomService.sendMessage(
        roomId as string,
        userId,
        newMessage,
        'text'
      );

      if (response.success) {
        setNewMessage('');
        // 消息会通过定时器获取到
      }
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  // 开始游戏
  const startGame = async () => {
    if (!roomId) return;

    try {
      const response = await roomService.startGame(roomId as string);
      if (response.success) {
        setGameStarted(true);
        setShowGame(true);
      } else {
        alert(response.message || '开始游戏失败');
      }
    } catch (error) {
      console.error('Error starting game:', error);
      alert('开始游戏失败');
    }
  };

  // 处理游戏结束
  const handleGameOver = (winnerId: string) => {
    setGameStarted(false);
    setShowGame(false);
    alert(`游戏结束！获胜者: ${winnerId}`);
  };

  if (showGame && typeof roomId === 'string' && userId) {
    return (
      <GameBoard 
        roomId={roomId} 
        playerId={userId} 
        onGameOver={handleGameOver} 
      />
    );
  }

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <Head>
        <title>游戏房间 - {roomId}</title>
      </Head>

      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h1>房间: {roomId}</h1>
        <button className="btn" onClick={() => router.back()}>返回首页</button>
      </header>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '20px' }}>
        {/* 游戏区域或聊天区域 */}
        <div>
          {!gameStarted ? (
            <div style={{ textAlign: 'center', padding: '40px 0' }}>
              <h2>等待其他玩家加入...</h2>
              <p>房间 ID: {roomId}</p>
              
              {userId && (
                <div style={{ marginTop: '30px' }}>
                  <button 
                    className="btn btn-success" 
                    onClick={startGame}
                    style={{ fontSize: '18px', padding: '12px 24px' }}
                  >
                    开始游戏
                  </button>
                  <p style={{ marginTop: '15px', color: '#aaa' }}>
                    当所有玩家准备就绪时，房主可以开始游戏
                  </p>
                </div>
              )}
            </div>
          ) : (
            <div style={{ textAlign: 'center', padding: '20px' }}>
              <h2>游戏进行中...</h2>
              <button 
                className="btn btn-success" 
                onClick={() => setShowGame(true)}
                style={{ fontSize: '18px', padding: '12px 24px' }}
              >
                进入游戏
              </button>
            </div>
          )}

          {/* 聊天区域 */}
          <div className="chat-container">
            <h3>房间聊天</h3>
            <div className="chat-messages">
              {messages.length > 0 ? (
                messages.map((msg, index) => (
                  <div 
                    key={index} 
                    className={`chat-message ${msg.type === 'system' ? 'system' : ''}`}
                  >
                    <strong>{msg.sender_username || '系统'}:</strong> {msg.content}
                    <small style={{ float: 'right', opacity: 0.7 }}>
                      {new Date(msg.created_at * 1000).toLocaleTimeString()}
                    </small>
                  </div>
                ))
              ) : (
                <p style={{ color: '#777', fontStyle: 'italic' }}>暂无消息</p>
              )}
              <div ref={messagesEndRef} />
            </div>
            
            <div className="chat-input">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder="输入消息..."
                onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
              />
              <button onClick={sendMessage}>发送</button>
            </div>
          </div>
        </div>

        {/* 玩家列表 */}
        <div style={{ backgroundColor: '#1a1a1a', padding: '20px', borderRadius: '8px', border: '1px solid #333' }}>
          <h3>房间玩家</h3>
          <div>
            {players.length > 0 ? (
              players.map((player, index) => (
                <div key={index} style={{ padding: '10px', borderBottom: '1px solid #333' }}>
                  <span style={{ color: '#4CAF50' }}>🎮 {player.username}</span>
                </div>
              ))
            ) : (
              <p style={{ color: '#777', fontStyle: 'italic' }}>暂无玩家</p>
            )}
          </div>

          <div style={{ marginTop: '30px' }}>
            <h4>操作</h4>
            <button className="btn" style={{ width: '100%', marginBottom: '10px' }}>邀请好友</button>
            <button className="btn btn-danger" style={{ width: '100%' }}>离开房间</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RoomPage;