# Docker 镜像导入 containerd（k3s）方案

k3s 使用 containerd 作为运行时，Docker 构建的镜像需要导入后才能被 k3s 使用。可选方案如下。

## 方案一：脚本导入（当前推荐）

项目根目录提供 `import-to-containerd.sh`，在本地用 Docker 构建好镜像后，一键导入到 k3s 的 containerd。

```bash
# 先构建镜像（示例）
./build.sh
docker build -f Dockerfile.gateway -t snake-game-gateway:latest .
docker build -f Dockerfile.lobby -t snake-game-lobby:latest .
# ... 其他服务

# 导入全部默认镜像
./import-to-containerd.sh

# 只导入指定镜像
./import-to-containerd.sh snake-game-gateway:latest snake-game-lobby:latest

# 查看默认会导入哪些镜像
./import-to-containerd.sh --list
```

**要求**：本机已安装 Docker 和 k3s（或 containerd），且能执行 `sudo k3s ctr`。  
**适用**：本地/内网开发、一次性或少量镜像更新。

---

## 方案二：本地 Registry

在集群内或本机起一个 Docker Registry，Docker push 到该 registry，k3s 从 registry 拉镜像，**不再需要** `docker save | ctr import`。

1. 部署 registry（例如在 k3s 里或本机）：
   ```bash
   docker run -d -p 5000:5000 --restart=always --name registry registry:2
   ```

2. 构建并推送到本地 registry：
   ```bash
   docker build -f Dockerfile.gateway -t localhost:5000/snake-game-gateway:latest .
   docker push localhost:5000/snake-game-gateway:latest
   ```

3. 配置 k3s 信任该 registry（insecure 或 TLS），并在 deployment 里把 `image` 改为 `localhost:5000/snake-game-gateway:latest`（或你的 registry 地址）。

4. k3s 拉取镜像：
   ```bash
   kubectl get pods -n snake-game-debug  # 触发拉取，或 rollout restart
   ```

**适用**：多台机器共用、CI 构建后推送到同一 registry、不想每次在本机 import。

→ 本地 Registry 的完整步骤见 **[docs/local-registry.md](local-registry.md)**（含 k3s 内部署 Registry、registries.yaml、推送脚本、overlay-registry）。

---

## 方案三：nerdctl 直接构建到 containerd

[nerdctl](https://github.com/containerd/nerdctl) 是 containerd 的 Docker CLI 兼容工具，可以在 k3s 使用的 containerd 里直接构建/拉取镜像，无需经过 Docker。

- 安装 nerdctl 并配置使用 k3s 的 containerd（例如 `containerd-root: /var/lib/rancher/k3s/agent/containerd` 或通过 `CONTAINERD_ADDRESS`）。
- 在项目里用 `nerdctl build -f Dockerfile.gateway -t snake-game-gateway:latest .` 构建，镜像会出现在 containerd 中。
- 若 k3s 使用独立 namespace（如 `k8s.io`），需在构建时指定该 namespace 或通过 `ctr -n k8s.io images import` 把 nerdctl 默认 namespace 的镜像导入到 `k8s.io`。

**适用**：不想装 Docker、希望构建和运行都在 containerd 的场景。

---

## 方案四：远程 Registry（生产常用）

镜像推送到 Docker Hub、Harbor、ECR 等远程 registry，k3s 从该 registry 拉取。

- 在 deployment 里使用完整镜像名，例如 `wdazhi97hub/snake_game:gateway-latest`。
- 若为私有仓库，在 k8s 里配置 `imagePullSecrets`。
- 构建与部署流程：CI 用 Docker 构建 → push → k3s 自动拉取（或 ArgoCD 同步）。

**适用**：生产、多环境、团队协作。

---

## 小结

| 方案           | 操作复杂度 | 适用场景           |
|----------------|------------|--------------------|
| 脚本导入       | 低         | 本地开发、快速验证 |
| 本地 Registry  | 中         | 内网多机、CI 推送  |
| nerdctl        | 中         | 纯 containerd 环境 |
| 远程 Registry  | 中         | 生产、多环境       |

当前仓库已提供 **方案一** 的 `import-to-containerd.sh`，其余方案可按需自行落地。
