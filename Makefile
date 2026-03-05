# Makefile for get_jobs
# AI 驱动的求职自动化工具

# 项目名称和版本
PROJECT_NAME := get_jobs
MODULE_NAME := github.com/loks666/get_jobs
VERSION := 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GO := go

# 构建目标
BINARY := $(PROJECT_NAME)
BUILD_DIR := .
MAIN_FILE := cmd/main.go

# Go 构建参数
GO_BUILD_FLAGS := -ldflags "-s -w"
GO_TEST_FLAGS := -v -race -coverprofile=coverage.out

# 颜色输出
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help build run clean test deps fmt lint install dev check run-Prod

# 默认目标：显示帮助
help:
	@echo ""
	@echo -e "$(BLUE)========================================$(NC)"
	@echo -e "$(BLUE)  $(PROJECT_NAME) - Makefile 帮助$(NC)"
	@echo -e "$(BLUE)========================================$(NC)"
	@echo ""
	@echo -e "$(GREEN)可用目标:$(NC)"
	@echo ""
	@echo -e "  $(YELLOW)make build$(NC)       - 构建项目二进制文件"
	@echo -e "  $(YELLOW)make run$(NC)         - 运行项目（开发模式）"
	@echo -e "  $(YELLOW)make run-prod$(NC)    - 运行项目（生产模式）"
	@echo -e "  $(YELLOW)make clean$(NC)       - 清理构建产物和缓存"
	@echo -e "  $(YELLOW)make test$(NC)       - 运行测试"
	@echo -e "  $(YELLOW)make deps$(NC)       - 下载和更新依赖"
	@echo -e "  $(YELLOW)make fmt$(NC)        - 格式化代码"
	@echo -e "  $(YELLOW)make lint$(NC)       - 代码检查"
	@echo -e "  $(YELLOW)make install$(NC)    - 安装依赖并构建"
	@echo -e "  $(YELLOW)make dev$(NC)        - 开发模式（带调试）"
	@echo -e "  $(YELLOW)make check$(NC)      - 检查环境配置"
	@echo ""
	@echo -e "$(BLUE)========================================$(NC)"
	@echo ""

# 构建项目
build:
	@echo -e "$(BLUE)Building $(PROJECT_NAME)...$(NC)"
	@$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY) $(MAIN_FILE)
	@echo -e "$(GREEN)Build successful!$(NC) Binary: $(BINARY)"

# 构建特定平台
build-darwin:
	@echo -e "$(BLUE)Building for darwin...$(NC)"
	@GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY)-darwin-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY)-darwin-arm64 $(MAIN_FILE)
	@echo -e "$(GREEN)Darwin build complete!$(NC)"

build-linux:
	@echo -e "$(BLUE)Building for linux...$(NC)"
	@GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY)-linux-amd64 $(MAIN_FILE)
	@echo -e "$(GREEN)Linux build complete!$(NC)"

build-windows:
	@echo -e "$(BLUE)Building for windows...$(NC)"
	@GOOS=windows $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY).exe $(MAIN_FILE)
	@echo -e "$(GREEN)Windows build complete!$(NC)"

# 运行项目
run: build
	@echo -e "$(BLUE)Running $(PROJECT_NAME)...$(NC)"
	@./$(BINARY)

# 生产模式运行
run-prod: build
	@echo -e "$(BLUE)Running $(PROJECT_NAME) in production mode...$(NC)"
	@./$(BINARY) --production

# 清理构建产物
clean:
	@echo -e "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -f $(BINARY)
	@rm -f $(BINARY)-*
	@rm -f $(BINARY).exe
	@rm -f coverage.out
	@rm -rf dist/
	@echo -e "$(GREEN)Clean complete!$(NC)"

# 运行测试
test:
	@echo -e "$(BLUE)Running tests...$(NC)"
	@$(GO) test $(GO_TEST_FLAGS) ./...
	@echo -e "$(GREEN)Tests complete!$(NC)"

# 运行特定包的测试
test-cover:
	@echo -e "$(BLUE)Running tests with coverage...$(NC)"
	@$(GO) test $(GO_TEST_FLAGS) -covermode=atomic ./...
	@if [ -f coverage.out ]; then \
		echo -e "$(GREEN)Coverage report generated: coverage.out$(NC)"; \
	fi

# 下载和更新依赖
deps:
	@echo -e "$(BLUE)Downloading dependencies...$(NC)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo -e "$(GREEN)Dependencies updated!$(NC)"

# 格式化代码
fmt:
	@echo -e "$(BLUE)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@$(GO) vet ./...
	@echo -e "$(GREEN)Code formatted!$(NC)"

# 代码检查
lint:
	@echo -e "$(BLUE)Running linter...$(NC)"
	@$(GO) vet ./...
	@echo -e "$(GREEN)Linting complete!$(NC)"

# 安装依赖并构建
install: deps build
	@echo -e "$(GREEN)Installation complete!$(NC)"
	@echo "Binary installed: ./$(BINARY)"

# 开发模式（带调试信息）
dev: build
	@echo -e "$(BLUE)Running in development mode...$(NC)"
	@DELVE_PORT=2345 dlv exec $(BINARY) --accept-multiclient --headless --api-version=2

# 检查环境配置
check:
	@echo -e "$(BLUE)Checking environment...$(NC)"
	@echo ""
	@echo -e "$(YELLOW)Go version:$(NC) $$($(GO) version)"
	@echo -e "$(YELLOW)Go env GOPATH:$(NC) $$($(GO) env GOPATH)"
	@echo -e "$(YELLOW)Go env GOROOT:$(NC) $$($(GO) env GOROOT)"
	@echo ""
	@echo -e "$(BLUE)Checking required files...$(NC)"
	@if [ -f $(MAIN_FILE) ]; then \
		echo -e "$(GREEN)✓$(NC) Main file exists: $(MAIN_FILE)"; \
	else \
		echo -e "$(RED)✗$(NC) Main file missing: $(MAIN_FILE)"; \
	fi
	@if [ -f config.yaml ]; then \
		echo -e "$(GREEN)✓$(NC) Config file exists: config.yaml"; \
	else \
		echo -e "$(RED)✗$(NC) Config file missing: config.yaml"; \
	fi
	@if [ -f go.mod ]; then \
		echo -e "$(GREEN)✓$(NC) Go module file exists: go.mod"; \
	else \
		echo -e "$(RED)✗$(NC) Go module file missing: go.mod"; \
	fi
	@echo ""
	@echo -e "$(BLUE)Environment check complete!$(NC)"

# 创建必要目录
init:
	@echo -e "$(BLUE)Initializing project directories...$(NC)"
	@mkdir -p data logs resources
	@echo -e "$(GREEN)Directories created!$(NC)"

# 显示版本信息
version:
	@echo "$(PROJECT_NAME) version $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@$(GO) version

# 完整构建流程（依赖 -> 格式化 -> 检查 -> 构建）
all: deps fmt check build
	@echo -e "$(GREEN)All complete!$(NC)"
