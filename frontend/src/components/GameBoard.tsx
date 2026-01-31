import React, { useEffect, useState } from 'react';
import { gameService } from '../utils/api';

interface Position {
  x: number;
  y: number;
}

interface GameSnake {
  playerId: string;
  segments: { position: Position }[];
  color: string;
  length: number;
  score: number;
}

interface GameState {
  snakes: GameSnake[];
  foods: Position[];
  walls: Position[];
  status: string;
}

interface GameBoardProps {
  roomId: string;
  playerId: string;
  onGameOver: (winnerId: string) => void;
}

const GameBoard: React.FC<GameBoardProps> = ({ roomId, playerId, onGameOver }) => {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [direction, setDirection] = useState<string>('RIGHT');
  const [connected, setConnected] = useState(false);

  // 初始化游戏
  useEffect(() => {
    const initializeGame = async () => {
      try {
        const response = await gameService.joinGame(roomId, playerId);
        if (response.success) {
          setGameState(response);
          setConnected(true);
          
          // 开始监听游戏状态更新
          startGameLoop();
        } else {
          console.error('Failed to join game:', response.message);
        }
      } catch (error) {
        console.error('Error joining game:', error);
      }
    };

    initializeGame();

    // 清理函数
    return () => {
      // 离开游戏
      gameService.leaveGame(roomId, playerId).catch(console.error);
    };
  }, [roomId, playerId]);

  // 游戏循环
  const startGameLoop = () => {
    const gameLoop = setInterval(async () => {
      try {
        const response = await gameService.getGameState(roomId);
        if (response.success) {
          setGameState(response);
          
          // 检查游戏是否结束
          if (response.status === 'finished') {
            clearInterval(gameLoop);
            onGameOver(response.winner_player_id || '');
          }
        }
      } catch (error) {
        console.error('Error getting game state:', error);
      }
    }, 100); // 每100ms更新一次
  };

  // 处理键盘事件
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!connected) return;

      let newDirection = direction;
      switch (e.key) {
        case 'ArrowUp':
          if (direction !== 'DOWN') newDirection = 'UP';
          break;
        case 'ArrowDown':
          if (direction !== 'UP') newDirection = 'DOWN';
          break;
        case 'ArrowLeft':
          if (direction !== 'RIGHT') newDirection = 'LEFT';
          break;
        case 'ArrowRight':
          if (direction !== 'LEFT') newDirection = 'RIGHT';
          break;
      }

      if (newDirection !== direction) {
        setDirection(newDirection);
        // 发送移动指令到服务器
        gameService.move(roomId, playerId, newDirection).catch(console.error);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [direction, connected, roomId, playerId]);

  // 渲染游戏单元格
  const renderCell = (x: number, y: number) => {
    if (!gameState) return <div className="cell" key={`${x}-${y}`} />;

    // 检查是否是墙
    const isWall = gameState.walls.some(wall => wall.x === x && wall.y === y);
    if (isWall) {
      return <div className="cell wall" key={`${x}-${y}`} />;
    }

    // 检查是否是食物
    const isFood = gameState.foods.some(food => food.x === x && food.y === y);
    if (isFood) {
      return <div className="cell food" key={`${x}-${y}`} />;
    }

    // 检查是否是蛇的身体
    let snakeColor = '';
    for (const snake of gameState.snakes) {
      for (const segment of snake.segments) {
        if (segment.position.x === x && segment.position.y === y) {
          snakeColor = snake.color;
          break;
        }
      }
      if (snakeColor) break;
    }

    if (snakeColor) {
      return <div className={`cell snake`} style={{ backgroundColor: snakeColor }} key={`${x}-${y}`} />;
    }

    return <div className="cell" key={`${x}-${y}`} />;
  };

  // 创建游戏板
  const createGameBoard = () => {
    const board = [];
    for (let y = 0; y < 30; y++) {
      for (let x = 0; x < 30; x++) {
        board.push(renderCell(x, y));
      }
    }
    return board;
  };

  return (
    <div className="game-container">
      <h2>贪吃蛇游戏</h2>
      <div className="game-status">
        <p>房间: {roomId}</p>
        <p>状态: {gameState?.status || '连接中...'}</p>
        {gameState && (
          <div>
            <p>你的分数: {
              gameState.snakes.find(s => s.playerId === playerId)?.score || 0
            }</p>
          </div>
        )}
      </div>
      
      <div className="game-board">
        {createGameBoard()}
      </div>
      
      <div className="controls">
        <p>使用方向键控制蛇的移动</p>
      </div>
    </div>
  );
};

export default GameBoard;