import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';

const SinglePlayer = () => {
  const router = useRouter();
  const [gameBoard, setGameBoard] = useState<Array<Array<string>>>(Array(20).fill(null).map(() => Array(20).fill('empty')));
  const [snake, setSnake] = useState([{ x: 10, y: 10 }]);
  const [food, setFood] = useState({ x: 5, y: 5 });
  const [direction, setDirection] = useState<'UP' | 'DOWN' | 'LEFT' | 'RIGHT'>('RIGHT');
  const [gameRunning, setGameRunning] = useState(false);
  const [score, setScore] = useState(0);
  const [gameOver, setGameOver] = useState(false);
  const directionRef = useRef(direction);
  const gameLoopRef = useRef<NodeJS.Timeout | null>(null);

  // 生成随机食物位置
  const generateFood = useCallback((currentSnake: Array<{x: number, y: number}>) => {
    let newFood: {x: number, y: number};
    do {
      newFood = {
        x: Math.floor(Math.random() * 20),
        y: Math.floor(Math.random() * 20)
      };
      // 确保食物不在蛇身上
    } while (currentSnake.some(segment => segment.x === newFood.x && segment.y === newFood.y));
    
    return newFood;
  }, []);

  // 初始化游戏
  const initGame = useCallback(() => {
    const initialSnake = [{ x: 10, y: 10 }];
    const initialFood = generateFood(initialSnake);
    
    setSnake(initialSnake);
    setFood(initialFood);
    setDirection('RIGHT');
    directionRef.current = 'RIGHT';
    setScore(0);
    setGameOver(false);
    setGameRunning(true);
  }, [generateFood]);

  // 检查碰撞
  const checkCollision = (head: {x: number, y: number}, currentSnake: Array<{x: number, y: number}>) => {
    // 检查边界碰撞
    if (head.x < 0 || head.x >= 20 || head.y < 0 || head.y >= 20) {
      return true;
    }
    
    // 检查自身碰撞（跳过头部）
    for (let i = 1; i < currentSnake.length; i++) {
      if (head.x === currentSnake[i].x && head.y === currentSnake[i].y) {
        return true;
      }
    }
    
    return false;
  };

  // 游戏主循环
  const gameLoop = useCallback(() => {
    if (!gameRunning || gameOver) return;

    setSnake(currentSnake => {
      const head = { ...currentSnake[0] };
      
      // 根据方向移动头部
      switch (directionRef.current) {
        case 'UP':
          head.y -= 1;
          break;
        case 'DOWN':
          head.y += 1;
          break;
        case 'LEFT':
          head.x -= 1;
          break;
        case 'RIGHT':
          head.x += 1;
          break;
      }

      // 检查碰撞
      if (checkCollision(head, currentSnake)) {
        setGameRunning(false);
        setGameOver(true);
        return currentSnake;
      }

      const newSnake = [head, ...currentSnake];
      
      // 检查是否吃到食物
      if (head.x === food.x && head.y === food.y) {
        // 增加分数
        setScore(prev => prev + 10);
        // 生成新食物
        setFood(generateFood(newSnake));
      } else {
        // 没吃到食物则移除尾部
        newSnake.pop();
      }
      
      return newSnake;
    });
  }, [food, gameRunning, gameOver, generateFood]);

  // 设置键盘事件监听
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!gameRunning) return;

      switch (e.key) {
        case 'ArrowUp':
          if (directionRef.current !== 'DOWN') {
            directionRef.current = 'UP';
            setDirection('UP');
          }
          break;
        case 'ArrowDown':
          if (directionRef.current !== 'UP') {
            directionRef.current = 'DOWN';
            setDirection('DOWN');
          }
          break;
        case 'ArrowLeft':
          if (directionRef.current !== 'RIGHT') {
            directionRef.current = 'LEFT';
            setDirection('LEFT');
          }
          break;
        case 'ArrowRight':
          if (directionRef.current !== 'LEFT') {
            directionRef.current = 'RIGHT';
            setDirection('RIGHT');
          }
          break;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [gameRunning]);

  // 游戏循环定时器
  useEffect(() => {
    if (gameRunning && !gameOver) {
      gameLoopRef.current = setInterval(gameLoop, 150);
    }

    return () => {
      if (gameLoopRef.current) {
        clearInterval(gameLoopRef.current);
      }
    };
  }, [gameRunning, gameOver, gameLoop]);

  // 初始化游戏板
  useEffect(() => {
    const board = Array(20).fill(null).map(() => Array(20).fill('empty'));
    
    // 放置蛇
    snake.forEach(segment => {
      if (segment.x >= 0 && segment.x < 20 && segment.y >= 0 && segment.y < 20) {
        board[segment.y][segment.x] = 'snake';
      }
    });
    
    // 放置食物
    if (food.x >= 0 && food.x < 20 && food.y >= 0 && food.y < 20) {
      board[food.y][food.x] = 'food';
    }
    
    setGameBoard(board);
  }, [snake, food]);

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto', textAlign: 'center' }}>
      <Head>
        <title>单人游戏 - 贪吃蛇</title>
      </Head>

      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h1>🎮 单人贪吃蛇</h1>
        <button className="btn" onClick={() => router.push('/')}>返回首页</button>
      </header>

      <div style={{ marginBottom: '20px' }}>
        <h2>得分: {score}</h2>
        {gameOver && (
          <div style={{ color: 'red', fontSize: '24px', fontWeight: 'bold', margin: '20px 0' }}>
            游戏结束!
          </div>
        )}
      </div>

      <div 
        className="game-board" 
        style={{ 
          display: 'grid', 
          gridTemplateColumns: 'repeat(20, 20px)', 
          gridTemplateRows: 'repeat(20, 20px)', 
          gap: '1px', 
          border: '2px solid #333', 
          backgroundColor: '#222',
          margin: '20px auto',
          width: 'fit-content'
        }}
      >
        {gameBoard.map((row, rowIndex) => 
          row.map((cell, colIndex) => (
            <div 
              key={`${rowIndex}-${colIndex}`} 
              className="cell"
              style={{ 
                backgroundColor: cell === 'snake' ? '#4CAF50' : 
                                cell === 'food' ? '#F44336' : '#333',
                borderRadius: cell === 'snake' ? '2px' : cell === 'food' ? '50%' : '0',
              }}
            />
          ))
        )}
      </div>

      <div style={{ marginTop: '20px' }}>
        {!gameRunning ? (
          <button 
            className="btn btn-success" 
            onClick={initGame}
            style={{ fontSize: '18px', padding: '12px 24px' }}
          >
            {gameOver ? '重新开始' : '开始游戏'}
          </button>
        ) : (
          <button 
            className="btn btn-danger" 
            onClick={() => {
              setGameRunning(false);
              if (gameLoopRef.current) clearInterval(gameLoopRef.current);
            }}
          >
            暂停游戏
          </button>
        )}
      </div>

      <div style={{ marginTop: '30px', padding: '15px', backgroundColor: '#1a1a1a', borderRadius: '8px' }}>
        <h3>游戏说明</h3>
        <p>使用方向键 ← → ↑ ↓ 控制蛇的移动</p>
        <p>吃到红色食物可以增长身体并获得分数</p>
        <p>撞到墙壁或自己的身体会导致游戏结束</p>
      </div>
    </div>
  );
};

export default SinglePlayer;