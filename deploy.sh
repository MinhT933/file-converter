#!/usr/bin/env bash
set -Eeuo pipefail

# Hardcode STACK_DIR
STACK_DIR=/home/ubuntu/app/file-convert

cd "$STACK_DIR" || { echo "❌ STACK_DIR=$STACK_DIR not found"; exit 1; }

echo "==> Stop old containers..."
docker compose -f docker-compose.prod.yml down || true

# echo "==> Remove old image..."
# for repo in "192.168.1.100:5001/$IMAGE_NAME_SERVER" "192.168.1.100:5001/$IMAGE_NAME_WORKER"; do
#   docker images --format '{{.Repository}}:{{.Tag}}' \
#   | awk -v repo="$repo" -v keep="$TAG" \
#       'index($0, repo ":")==1 && $0 != (repo ":" keep) && $0 != (repo ":latest") {print}' \
#   | xargs -r -n1 docker rmi -f || true
# done

echo "==> Pull latest image..."
docker compose -f docker-compose.prod.yml pull

echo "==> Start new containers..."
docker compose -f docker-compose.prod.yml up -d --remove-orphans


echo "==> Current status..."
docker ps
