#!/usr/bin/env bash
set -Eeuo pipefail

# Hardcode STACK_DIR
STACK_DIR=/home/ubuntu/app/file-convert

cd "$STACK_DIR" || { echo "❌ STACK_DIR=$STACK_DIR not found"; exit 1; }

echo "==> Stop old containers..."
docker compose -f docker-compose.prod.yml down

echo "==> Remove old image..."
docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' \
| awk -v keep="$TAG" '$1 ~ ("^"repo":") && $1 != (repo":"keep) {print $2}' \
| xargs -r docker rmi -f
done

echo "==> Pull latest image..."
docker compose -f docker-compose.prod.yml pull

echo "==> Start new containers..."
docker compose -f docker-compose.prod.yml up -d --remove-orphans


echo "==> Current status..."
docker ps
