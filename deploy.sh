#!/usr/bin/env bash
set -Eeuo pipefail

STACK_DIR=/home/ubuntu/app/file-convert

cd "$STACK_DIR" || { echo "❌ STACK_DIR=$STACK_DIR not found"; exit 1; }

echo "==> Stop old containers..."
docker compose -f docker-compose.prod.yml down || true

: "${TAG:?Missing TAG}"
: "${IMAGE_NAME_SERVER:?Missing IMAGE_NAME_SERVER}"
: "${IMAGE_NAME_WORKER:?Missing IMAGE_NAME_WORKER}"
REGISTRY_HOST="${REGISTRY_HOST:-192.168.1.100:5001}"

echo "==> Remove old images (keep :$TAG & :latest, only these 2 repos)..."
for repo in "$REGISTRY_HOST/$IMAGE_NAME_SERVER" "$REGISTRY_HOST/$IMAGE_NAME_WORKER"; do
  echo "Repo: $repo"
  list="$(docker image ls "$repo" --format '{{.Repository}}:{{.Tag}}' \
          | grep -v -E ":(${TAG}|latest)$" || true)"
  if [ -n "$list" ]; then
    echo "$list" | xargs -r -n1 docker rmi -f || true
  else
    echo "Nothing to delete."
  fi
done

docker image prune -f >/dev/null 2>&1 || true

echo "==> Pull latest image..."
docker compose -f docker-compose.prod.yml pull

echo "==> Start new containers..."
docker compose -f docker-compose.prod.yml up -d --remove-orphans


echo "==> Current status..."
docker ps
