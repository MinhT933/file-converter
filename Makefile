# Binary output
BIN_DIR := bin
SERVER_BIN := $(BIN_DIR)/server
WORKER_BIN := $(BIN_DIR)/worker

# Entry points
SERVER_ENTRY := cmd/server/main.go
WORKER_ENTRY := cmd/worker/main.go

# Build server
.PHONY: build-server
build-server:
	@echo "🔨 Building server..."
	go build -o $(SERVER_BIN) $(SERVER_ENTRY)

# Build worker
.PHONY: build-worker
build-worker:
	@echo "🔨 Building worker..."
	go build -o $(WORKER_BIN) $(WORKER_ENTRY)

# Run server
.PHONY: run-server
run-server:
	@echo "🚀 Running server..."
	go run $(SERVER_ENTRY)

# Run worker
.PHONY: run-worker
run-worker:
	@echo "🚀 Running worker..."
	go run $(WORKER_ENTRY)

# Clean
.PHONY: clean
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BIN_DIR)
