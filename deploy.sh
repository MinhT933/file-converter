#!/usr/bin/env bash
set -Eeuo pipefail

STACK_DIR=/home/ubuntu/app/file-convert

cd "$STACK_DIR" || { echo "❌ STACK_DIR=$STACK_DIR not found"; exit 1; }

echo "==> Stop old containers..."
docker compose -f docker-compose.prod.yml down || true

echo "==> Remove old images safely (keep :$TAG & :latest)..."
for repo in "${REPOS[@]}"; do
  echo "Repo: $repo"

  # liệt kê đúng repo (không ảnh hưởng repo khác), bỏ tag đang giữ & latest
  mapfile -t CANDIDATES < <(
    docker image ls --filter "reference=${repo}:*" --format '{{.Repository}}:{{.Tag}}' \
    | grep -v -E ":(${TAG}|latest)$"
  )

  for img in "${CANDIDATES[@]}"; do
    # skip nếu image đang được container dùng (kể cả stopped)
    if docker ps -a -q --filter "ancestor=${img}" | grep -q .; then
      echo "  ⚠️  Skip (in use): $img"
      continue
    fi

    if [ "${DRY_RUN:-0}" = "1" ]; then
      echo "  (dry-run) would delete: $img"
    else
      echo "  🗑️  Deleting: $img"
      docker rmi -f "$img" || true
    fi
  done
done

# dọn layer rác (không ảnh hưởng image đang dùng)
docker image prune -f >/dev/null 2>&1 || true

echo "==> Pull latest image..."
docker compose -f docker-compose.prod.yml pull

echo "==> Start new containers..."
docker compose -f docker-compose.prod.yml up -d --remove-orphans


echo "==> Current status..."
docker ps
