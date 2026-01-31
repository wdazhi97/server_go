import React from 'react';
import Header from './Header';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Header />
      <main style={{ flex: 1, padding: '20px 0' }}>
        {children}
      </main>
      <footer style={{ 
        textAlign: 'center', 
        padding: '20px', 
        backgroundColor: '#111', 
        borderTop: '1px solid #333',
        marginTop: 'auto'
      }}>
        <p>© {new Date().getFullYear()} 贪吃蛇游戏 - 微服务版</p>
      </footer>
    </div>
  );
};

export default Layout;