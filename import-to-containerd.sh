#!/bin/bash
# 将 Docker 镜像导入到 containerd（供 k3s 使用）
# 用法:
#   ./import-to-containerd.sh                          # 导入默认全部镜像
#   ./import-to-containerd.sh snake-game-gateway:latest  # 只导入指定镜像
#   ./import-to-containerd.sh --list                   # 只打印默认镜像列表，不导入

set -e

# k3s 的 ctr：优先用 k3s 自带的，否则用系统 ctr
CTR="${CTR:-}"
if [[ -z "$CTR" ]]; then
  if command -v k3s &>/dev/null; then
    CTR="sudo k3s ctr"
  elif command -v ctr &>/dev/null; then
    CTR="ctr"
  else
    echo "Error: 未找到 ctr。请安装 k3s 或 containerd，或设置 CTR=your_ctr_command" >&2
    exit 1
  fi
fi

# 默认要导入的镜像（与 deployment 里使用的名称一致）
DEFAULT_IMAGES=(
  snake-game-gateway:latest
  snake-game-lobby:latest
  snake-game-matching:latest
  snake-game-room:latest
  snake-game-leaderboard:latest
  snake-game-game:latest
  snake-game-friends:latest
  snake-game-frontend:latest
)

import_one() {
  local img="$1"
  if ! docker image inspect "$img" &>/dev/null; then
    echo "Skip: $img (本地不存在)" >&2
    return 0
  fi
  echo "Importing $img ..."
  docker save "$img" | $CTR images import -
  echo "Done: $img"
}

if [[ "${1:-}" == "--list" ]]; then
  printf '%s\n' "${DEFAULT_IMAGES[@]}"
  exit 0
fi

if [[ $# -gt 0 ]]; then
  for img in "$@"; do
    import_one "$img"
  done
else
  for img in "${DEFAULT_IMAGES[@]}"; do
    import_one "$img"
  done
fi

echo "All done."
