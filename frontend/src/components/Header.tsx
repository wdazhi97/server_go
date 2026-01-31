import React from 'react';
import Link from 'next/link';

const Header: React.FC = () => {
  return (
    <header style={{
      backgroundColor: '#111',
      padding: '1rem',
      borderBottom: '1px solid #333'
    }}>
      <nav style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        maxWidth: '1200px',
        margin: '0 auto'
      }}>
        <Link href="/" style={{
          color: '#4CAF50',
          fontSize: '1.5rem',
          fontWeight: 'bold',
          textDecoration: 'none'
        }}>
          🐍 贪吃蛇游戏
        </Link>
        
        <div style={{
          display: 'flex',
          gap: '1rem'
        }}>
          <Link href="/" style={{
            color: '#fff',
            textDecoration: 'none',
            padding: '0.5rem 1rem',
            borderRadius: '4px',
            transition: 'background-color 0.3s'
          }}>
            首页
          </Link>
          <Link href="/leaderboard" style={{
            color: '#fff',
            textDecoration: 'none',
            padding: '0.5rem 1rem',
            borderRadius: '4px',
            transition: 'background-color 0.3s'
          }}>
            排行榜
          </Link>
          <Link href="/single-player" style={{
            color: '#fff',
            textDecoration: 'none',
            padding: '0.5rem 1rem',
            borderRadius: '4px',
            transition: 'background-color 0.3s'
          }}>
            单人游戏
          </Link>
        </div>
      </nav>
    </header>
  );
};

export default Header;