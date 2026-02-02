#!/bin/bash
# 使用 Earthly 构建全部镜像并导入到 containerd（k3s）
# 依赖：earthly、Docker、k3s（或 containerd）
set -e

echo "Building all images with Earthly..."
earthly +all

echo "Importing images to containerd..."
./import-to-containerd.sh

echo "Done. Images are in containerd."
