#!/bin/bash

# 获取版本信息
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}

echo "构建版本信息:"
echo "  Version: $VERSION"
echo ""

# 构建Docker镜像
docker build \
    --build-arg VERSION="$VERSION" \
    -t kotorimiku/komiko:$VERSION \
    -t kotorimiku/komiko:latest \
    .

echo ""
echo "构建完成！"
echo "镜像标签: komiko:$VERSION, komiko:latest"