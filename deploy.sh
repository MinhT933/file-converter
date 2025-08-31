#!/usr/bin/env bash
set -Eeuo pipefail

# Hardcode STACK_DIR
STACK_DIR=/home/ubuntu/app/file-convert


echo "==> Current dir before cd:"
pwd

cd "$STACK_DIR" || { echo "❌ STACK_DIR=$STACK_DIR not found"; exit 1; }

echo "==> Current dir after cd:"
pwd

echo "==> Pull latest image..."
docker compose -f docker-compose.prod.yml pull 

echo "==> Restart service..."
docker compose -f docker-compose.prod.yml up -d 

echo "==> Current status..."
docker ps
