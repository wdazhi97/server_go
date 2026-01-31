import React, { useState } from 'react';
import { useRouter } from 'next/router';
import { lobbyService } from '../utils/api';
import Head from 'next/head';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState('');
  const [isRegister, setIsRegister] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      if (isRegister) {
        // 注册
        const response = await lobbyService.register(username, password, email);
        if (response.success) {
          alert('注册成功！请登录。');
          setIsRegister(false);
        } else {
          setError(response.message || '注册失败');
        }
      } else {
        // 登录
        const response = await lobbyService.login(username, password);
        if (response.success) {
          // 保存用户信息到本地存储
          localStorage.setItem('userId', response.userId);
          localStorage.setItem('username', response.username);
          localStorage.setItem('token', btoa(`${username}:${password}`)); // 简单的认证令牌

          // 跳转到主页
          router.push('/');
        } else {
          setError(response.message || '登录失败');
        }
      }
    } catch (err) {
      setError('网络错误，请稍后再试');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-container">
      <Head>
        <title>{isRegister ? '注册' : '登录'} - 贪吃蛇游戏</title>
      </Head>
      
      <h2 style={{ textAlign: 'center', marginBottom: '20px' }}>
        {isRegister ? '注册新账户' : '登录到游戏'}
      </h2>
      
      {error && <div style={{ color: 'red', marginBottom: '15px' }}>{error}</div>}
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">用户名:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        
        {!isRegister && (
          <div className="form-group">
            <label htmlFor="email">邮箱:</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required={!isRegister}
            />
          </div>
        )}
        
        <div className="form-group">
          <label htmlFor="password">密码:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        
        <button 
          type="submit" 
          className="btn btn-success" 
          disabled={loading}
          style={{ width: '100%', marginTop: '10px' }}
        >
          {loading ? '处理中...' : (isRegister ? '注册' : '登录')}
        </button>
      </form>
      
      <div style={{ textAlign: 'center', marginTop: '15px' }}>
        <button 
          className="btn" 
          onClick={() => setIsRegister(!isRegister)}
          style={{ marginTop: '10px' }}
        >
          {isRegister ? '已有账户？去登录' : '没有账户？去注册'}
        </button>
      </div>
    </div>
  );
};

export default Login;