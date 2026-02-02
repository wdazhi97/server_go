#!/bin/bash
# 将本地构建好的镜像打上 Registry 标签并推送
# 用法:
#   ./push-to-registry.sh                    # 推送到默认 localhost:30500
#   REGISTRY=192.168.1.100:30500 ./push-to-registry.sh

set -e

REGISTRY="${REGISTRY:-localhost:30500}"

IMAGES=(
  snake-game-gateway:latest
  snake-game-lobby:latest
  snake-game-matching:latest
  snake-game-room:latest
  snake-game-leaderboard:latest
  snake-game-game:latest
  snake-game-friends:latest
  snake-game-frontend:latest
)

for img in "${IMAGES[@]}"; do
  if ! docker image inspect "$img" &>/dev/null; then
    echo "Skip: $img (本地不存在)" >&2
    continue
  fi
  echo "Pushing $img -> $REGISTRY/$img"
  docker tag "$img" "$REGISTRY/$img"
  docker push "$REGISTRY/$img"
  echo "Done: $REGISTRY/$img"
done

echo "All done. Registry: $REGISTRY"
