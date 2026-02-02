# 本地 Registry 方案（k3s 从 Registry 拉镜像）

不再使用 `docker save | ctr import`，改为：Docker 构建 → push 到本地 Registry → k3s 从 Registry 拉取。

## 一、在 k3s 里部署 Registry

已提供 `deployment/env/debug/base/registry/deployment.yaml`，随 debug 环境一起部署即可（ArgoCD 同步或手动 apply）。

- 集群内访问：`registry.snake-game-debug.svc.cluster.local:5000`
- 宿主机推送：通过 NodePort **30500**，即 `localhost:30500` 或 `<节点IP>:30500`

## 二、配置 k3s 信任本地 Registry

让 k3s 把 `registry.local:5000` 解析到集群内的 Registry 服务。

在 **k3s 节点**上创建或编辑 `/etc/rancher/k3s/registries.yaml`：

```yaml
mirrors:
  "registry.local:5000":
    endpoint:
      - "http://registry.snake-game-debug.svc.cluster.local:5000"
```

然后重启 k3s：

```bash
sudo systemctl restart k3s
# 或
sudo systemctl restart k3s-agent
```

> 若 Registry 部署在别的 namespace，把上面的 `snake-game-debug` 改成对应 namespace。

## 三、构建、打标签、推送

推送时使用 **NodePort 地址**（宿主机能访问），例如本机：`localhost:30500`。

```bash
# 在 server_go 根目录
./build.sh

# 构建各服务镜像（示例）
docker build -f Dockerfile.gateway -t snake-game-gateway:latest .
docker build -f Dockerfile.lobby -t snake-game-lobby:latest .
# ... 其他服务、前端

# 使用脚本一次性打标签并推送到本地 Registry（REGISTRY 默认 localhost:30500）
./push-to-registry.sh

# 或指定 Registry 地址（多节点时用节点 IP）
REGISTRY=192.168.1.100:30500 ./push-to-registry.sh
```

## 四、让 k3s 使用本地 Registry 的镜像

两种方式任选其一。

### 方式 A：使用 overlay-registry（推荐）

用带「镜像替换」的 overlay 部署，这样会从 `registry.local:5000` 拉取镜像。

```bash
# 本地 apply 示例（路径按你仓库为准）
kubectl apply -k deployment/env/debug/overlay-registry

# 若用 ArgoCD：把 Application 的 path 从 overlay 改为 overlay-registry
# 例如 path: deployment/env/debug/overlay-registry
```

### 方式 B：不改 overlay，只改 ArgoCD 的 kustomize path

在 ArgoCD 里把 debug 应用的 Kustomize 路径从 `deployment/env/debug/overlay` 改为 `deployment/env/debug/overlay-registry`，并同步。

## 五、验证

```bash
# 查看 Registry 是否在跑
kubectl get pods -n snake-game-debug -l app=registry

# 查看本地 Registry 里的镜像（需在集群内或 port-forward）
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
  curl -s http://registry.snake-game-debug.svc.cluster.local:5000/v2/_catalog

# 触发一次拉取（例如重启 gateway）
kubectl rollout restart deployment/gateway -n snake-game-debug
kubectl get pods -n snake-game-debug -w
```

## 六、小结

| 步骤 | 操作 |
|------|------|
| 1 | 部署 Registry（随 debug base 或单独 apply） |
| 2 | 在 k3s 节点配置 `/etc/rancher/k3s/registries.yaml` 并重启 k3s |
| 3 | 本地 `docker build` 后执行 `./push-to-registry.sh`（或手动 tag + push 到 `localhost:30500`） |
| 4 | 使用 `overlay-registry` 部署，或把 ArgoCD 指向 `overlay-registry` |

之后只需：改代码 → 构建 → `./push-to-registry.sh`，k3s 会按 `IfNotPresent` 从本地 Registry 拉取新镜像。
