#!/bin/bash
set -e  # Chỉ dừng khi có lỗi, không ghi log mỗi lệnh (-x)

cd /app

# Tự động reload khi code thay đổi
if [ "${DEV_MODE:-false}" = "true" ]; then
  echo "🔄 Starting in development mode with auto-reload..."
  
  while true; do
    echo "Generating Swagger docs..."
    swag init -g cmd/server/main.go --parseInternal --parseDependency
    
    # Compile lại server từ source
    echo "🔨 Compiling server..."
    go build -o myapp-server-dev ./cmd/server
    
    echo "Starting server..."
    ./myapp-server-dev &
    SERVER_PID=$!
    
    # Tính hash ban đầu
    LAST_HASH=$(find /app/cmd /app/internal -type f -name "*.go" -exec md5sum {} \; | sort | md5sum | cut -d' ' -f1)
    echo "👀 Watching for code changes (press Ctrl+C to stop)..."
    
    # Chỉ hiển thị thời gian kiểm tra cuối cùng mỗi 30 giây
    CHECK_COUNT=0
    LAST_LOG_TIME=$(date +%s)
    
    while true; do
      sleep 2
      CURRENT_HASH=$(find /app/cmd /app/internal -type f -name "*.go" -exec md5sum {} \; | sort | md5sum | cut -d' ' -f1)
      
      # Chỉ log mỗi 30 giây
      CURRENT_TIME=$(date +%s)
      CHECK_COUNT=$((CHECK_COUNT + 1))
      
      if [ $((CURRENT_TIME - LAST_LOG_TIME)) -gt 30 ]; then
        echo "⏱️  Still watching... (checked $CHECK_COUNT times)"
        LAST_LOG_TIME=$CURRENT_TIME
        CHECK_COUNT=0
      fi
      
      # Nếu phát hiện thay đổi
      if [ "$CURRENT_HASH" != "$LAST_HASH" ]; then
        echo "🔄 Code changed, reloading..."
        kill $SERVER_PID || echo "Failed to kill $SERVER_PID"
        break
      fi
    done
  done
else
  # Production mode
  echo "⚙️ Re-generating Swagger docs..."
  swag init -g cmd/server/main.go --parseInternal --parseDependency
  echo "✅ Swagger docs generated"
  echo "🚀 Starting server..."
  exec ./myapp-server
fi