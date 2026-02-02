# 构建与导入 containerd

镜像由 **Earthly** 构建，构建后通过脚本导入到 **containerd**（供 k3s 使用）。不再使用 Dockerfile。

## 结构

| 位置 | 作用 |
|------|------|
| **buildtool/** | 公共运行时 base（Ubuntu 24.04 + ca-certificates），各服务 docker 阶段依赖它 |
| **Earthfile**（根目录） | 各服务的 **build** target（在根上下文编译，因 go.mod 在根目录）+ **all**；不拷贝父目录 |
| **gateway/Earthfile、lobby/Earthfile、…** | 各服务目录**单独构建**：通过 `BUILD ../+gateway-build` 使用根目录的编译结果，本目录只做 docker 阶段，不 COPY 父目录 |
| **build-and-import.sh** | 先 `earthly +all` 构建全部镜像，再 `import-to-containerd.sh` 导入到 containerd |
| **import-to-containerd.sh** | 将本地 Docker 中的镜像导入到 k3s/containerd |

## 依赖

- [Earthly](https://earthly.dev/)（已安装则跳过）
- Docker（Earthly 的 buildkit 需要）
- k3s 或 containerd（导入目标）

## 用法

**在仓库根目录执行：**

```bash
# 构建全部镜像并导入到 containerd（推荐）
./build-and-import.sh
```

或分步：

```bash
# 只构建全部镜像（写入本地 Docker）
earthly +all

# 按目录单独构建（在仓库根目录执行，不拷贝父目录）
earthly ./gateway+docker
earthly ./lobby+docker
earthly ./frontend+docker
# ...

# 将已构建的镜像导入到 containerd
./import-to-containerd.sh

# 只导入指定镜像
./import-to-containerd.sh snake-game-gateway:latest snake-game-lobby:latest

# 查看默认会导入的镜像列表
./import-to-containerd.sh --list
```

## 镜像名称（与 deployment 一致）

- snake-game-gateway:latest
- snake-game-lobby:latest
- snake-game-matching:latest
- snake-game-room:latest
- snake-game-leaderboard:latest
- snake-game-game:latest
- snake-game-friends:latest
- snake-game-frontend:latest

## 保留的旧文件（可选）

- **Dockerfile.backup** / **Dockerfile.simple**：多阶段 Docker 构建，仅作备份，日常请用 Earthly。
