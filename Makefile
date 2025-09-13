SHELL := /bin/bash
GO    := go

# ---- Project & compose ----
PROJECT        ?= fileconv
COMPOSE_INFRA  ?= compose.infra.yml
COMPOSE_APP    ?= compose.app.yml
NET            ?= fileconv-net          # external network dùng chung infra <-> app
TARGET         ?= dev                   # dev | prod  (match build.target trong compose)

# ---- Go build info ----
BIN_DIR       := bin
SERVER_BIN    := $(BIN_DIR)/server
WORKER_BIN    := $(BIN_DIR)/worker
SERVER_ENTRY  := cmd/server/main.go
WORKER_ENTRY  := cmd/worker/main.go

DATE          := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT        := $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
VERSION       := $(shell git describe --tags --always --dirty 2>/dev/null || echo 'v0.0.0')
LDFLAGS       := -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'

CGO ?= 0
GOOS ?= linux
GOARCH ?= amd64

.PHONY: help
help: ## Hiển thị trợ giúp
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ────────────────────────────────────────────────────────────────────────────────
# Go: build / run / test
# ────────────────────────────────────────────────────────────────────────────────
.PHONY: deps
deps: ## go mod download
	$(GO) mod download

.PHONY: tidy
tidy: ## go mod tidy
	$(GO) mod tidy

.PHONY: fmt
fmt: ## go fmt & go vet
	$(GO) fmt ./...
	$(GO) vet ./...

.PHONY: test
test: ## go test
	$(GO) test ./...

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: build-server
build-server: $(BIN_DIR) ## build server (static, linux/amd64)
	CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	$(GO) build -ldflags "$(LDFLAGS)" -o $(SERVER_BIN) $(SERVER_ENTRY)

.PHONY: build-worker
build-worker: $(BIN_DIR) ## build worker (static, linux/amd64)
	CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	$(GO) build -ldflags "$(LDFLAGS)" -o $(WORKER_BIN) $(WORKER_ENTRY)

.PHONY: build
build: build-server build-worker ## build cả server & worker

.PHONY: run-server
run-server: ## go run server (dành cho local không docker)
	$(GO) run -ldflags "$(LDFLAGS)" $(SERVER_ENTRY)

.PHONY: run-worker
run-worker: ## go run worker (dành cho local không docker)
	$(GO) run -ldflags "$(LDFLAGS)" $(WORKER_ENTRY)

.PHONY: swag
swag: ## gen swagger docs
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/server/main.go --parseInternal --parseDependency

.PHONY: clean
clean: ## xóa bin/
	rm -rf $(BIN_DIR)

# ────────────────────────────────────────────────────────────────────────────────
# Docker Compose: infra / app (dev|prod)
# YÊU CẦU: trong compose.app.yml, phần build nên có: target: ${TARGET:-dev}
# ────────────────────────────────────────────────────────────────────────────────
.PHONY: network
network: ## tạo external network dùng chung
	- docker network create $(NET)

.PHONY: infra-up
infra-up: network ## bật infra (redis/nats/pg/pgadmin)
	docker compose -f $(COMPOSE_INFRA) up -d

.PHONY: infra-down
infra-down: ## tắt infra
	docker compose -f $(COMPOSE_INFRA) down

.PHONY: app-up
app-up: ## bật app+worker (TARGET=dev|prod)
	@echo "▶ TARGET=$(TARGET)"
	TARGET=$(TARGET) docker compose -f $(COMPOSE_APP) up -d --build

.PHONY: app-down
app-down: ## tắt app+worker
	docker compose -f $(COMPOSE_APP) down

.PHONY: up-all
up-all: infra-up app-up ## bật cả infra + app

.PHONY: down-all
down-all: app-down infra-down ## tắt cả infra + app

.PHONY: logs
logs: ## tail logs app & worker
	docker compose -f $(COMPOSE_APP) logs -f app worker

.PHONY: ps
ps: ## xem trạng thái containers
	docker compose -f $(COMPOSE_INFRA) ps
	docker compose -f $(COMPOSE_APP) ps

.PHONY: restart-app
restart-app: ## restart service app
	docker compose -f $(COMPOSE_APP) restart app

.PHONY: restart-worker
restart-worker: ## restart service worker
	docker compose -f $(COMPOSE_APP) restart worker

.PHONY: watch
watch:
	air -c .air.toml 

.PHONY: build-server
	go build -o ./bin/server ./cmd/server/main.go     

.PHONY: scale-worker
scale-worker: ## scale worker, ví dụ: make scale-worker N=3
ifndef N
	$(error "Thiếu N. Ví dụ: make scale-worker N=3")
endif
	TARGET=$(TARGET) docker compose -f $(COMPOSE_APP) up -d --scale worker=$(N) worker
