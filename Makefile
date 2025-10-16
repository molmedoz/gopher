# Gopher - Go Version Manager Makefile

# Variables
PROJECT_NAME := gopher
BUILD_DIR := build
BINARY_NAME := gopher

# Get version from git tag or default to dev
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_BY := manual

# ldflags for version injection
LDFLAGS := -ldflags "\
	-X main.appVersion=$(VERSION) \
	-X main.appCommit=$(COMMIT) \
	-X main.appDate=$(BUILD_DATE) \
	-X main.appBuiltBy=$(BUILD_BY) \
	-s -w"

# Go variables
GO := go
GOMOD := $(GO) mod
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOCLEAN := $(GO) clean
GOINSTALL := $(GO) install
GOFMT := $(GO) fmt
GOPATH := $(shell $(GO) env GOPATH)
GOBIN := $(GOPATH)/bin

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "$(CYAN)Gopher - Go Version Manager$(NC)"
	@echo "$(YELLOW)============================$(NC)"
	@echo ""
	@echo "$(GREEN)Build Commands:$(NC)"
	@echo "  build          - Build the binary"
	@echo "  build-all      - Build for all platforms"
	@echo "  install        - Install to system"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "$(GREEN)Development Commands:$(NC)"
	@echo "  test           - Run all tests with race detection and coverage"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  fmt            - Format Go code with go fmt"
	@echo "  imports        - Format imports with goimports"
	@echo "  format         - Format code and imports (comprehensive)"
	@echo "  check-format   - Check code format (CI mode, no modifications)"
	@echo "  check-imports  - Check imports format (CI mode, no modifications)"
	@echo "  lint           - Run linter"
	@echo "  vet            - Run go vet"
	@echo "  check          - Run all checks (fmt + vet + test)"
	@echo "  ci             - Run CI checks locally (matches GitHub Actions)"
	@echo "  install-tools  - Install all development tools"
	@echo ""
	@echo "$(GREEN)Release Commands:$(NC)"
	@echo "  release              - Create release build"
	@echo "  dist                 - Create distribution packages"
	@echo "  goreleaser-check     - Validate GoReleaser config"
	@echo "  goreleaser-snapshot  - Build snapshot (doesn't publish)"
	@echo "  prepare-release      - Prepare for release (requires VERSION=X.Y.Z)"
	@echo "                         Use GitHub Actions 'Create Release' workflow to actually release"
	@echo ""
	@echo "$(GREEN)Utility Commands:$(NC)"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy dependencies"
	@echo "  check          - Run all checks"
	@echo "  version        - Show version information"
	@echo ""
	@echo "$(GREEN)Security Commands:$(NC)"
	@echo "  security-tools - Install security tools"
	@echo "  security-scan  - Run security scan"
	@echo "  vuln-check     - Check for vulnerabilities"
	@echo "  security-test  - Run security tests"
	@echo "  security-all   - Run all security checks"
	@echo ""
	@echo "$(GREEN)Docker Commands:$(NC)"
	@echo "  docker-build   - Build all Docker images"
	@echo "  docker-test    - Run Docker tests"
	@echo "  docker-clean   - Clean Docker images"
	@echo "  docker-test-unix-no-go      - Test Unix no Go scenario"
	@echo "  docker-test-unix-with-go    - Test Unix with Go scenario"
	@echo "  docker-test-windows-no-go   - Test Windows no Go scenario"
	@echo "  docker-test-windows-with-go - Test Windows with Go scenario"
	@echo "  docker-test-windows-simulated - Test Windows simulation container"
	@echo "  docker-test-macos-no-go     - Test macOS no Go scenario"
	@echo "  docker-test-macos-with-go   - Test macOS with Go scenario"
	@echo ""
	@echo "$(YELLOW)Usage: make <command>$(NC)"
	@echo "$(YELLOW)Example: make build$(NC)"

# Build Commands
.PHONY: build
build: ## Build the binary
	@echo "$(BLUE)Building $(BINARY_NAME) $(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gopher
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@echo "$(CYAN)Version: $(VERSION) (commit: $(COMMIT))$(NC)"

.PHONY: build-all
build-all: ## Build for all platforms
	@echo "$(BLUE)Building for all platforms ($(VERSION))...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@echo "$(YELLOW)Building for Linux amd64...$(NC)"
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/gopher
	@echo "$(YELLOW)Building for Linux arm64...$(NC)"
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/gopher
	@echo "$(YELLOW)Building for macOS amd64...$(NC)"
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/gopher
	@echo "$(YELLOW)Building for macOS arm64...$(NC)"
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/gopher
	@echo "$(YELLOW)Building for Windows amd64...$(NC)"
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/gopher
	@echo "$(GREEN)✅ All platform builds complete ($(VERSION))$(NC)"

.PHONY: install
install: build ## Install to system
	@echo "$(BLUE)Installing $(BINARY_NAME) to system...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)✅ $(BINARY_NAME) installed to /usr/local/bin/$(NC)"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✅ Build artifacts cleaned$(NC)"

# Development Commands
.PHONY: test
test: ## Run all tests with race detection and coverage
	@echo "$(BLUE)Running tests with race detection and coverage...$(NC)"
	@$(GOTEST) -race -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...
	@$(GO) tool cover -func=coverage.out | tail -n 1
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	@$(GOTEST) -v ./...
	@echo "$(GREEN)✅ Tests completed$(NC)"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage (HTML + summary)
	@echo "$(BLUE)Running tests with coverage (detailed)...$(NC)"
	@$(GOTEST) -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...
	@$(GO) tool cover -func=coverage.out | tail -n 1
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

.PHONY: fmt
fmt: ## Format Go code with go fmt
	@echo "$(BLUE)Formatting Go code...$(NC)"
	@$(GOFMT) ./...
	@echo "$(GREEN)✅ Go code formatted$(NC)"

.PHONY: imports
imports: ## Format imports with goimports
	@echo "$(BLUE)Formatting imports with goimports...$(NC)"
	@if [ -f "$(GOBIN)/goimports" ]; then \
		$(GOBIN)/goimports -w .; \
		echo "$(GREEN)✅ Imports formatted$(NC)"; \
	else \
		echo "$(YELLOW)Warning: goimports not installed$(NC)"; \
		echo "$(CYAN)Run 'make install-tools' to install goimports$(NC)"; \
		exit 1; \
	fi

.PHONY: format
format: fmt imports ## Format code and imports (comprehensive)
	@echo "$(GREEN)✅ All formatting complete$(NC)"

.PHONY: check-format
check-format: ## Check if code is formatted (CI mode, doesn't modify files)
	@echo "$(BLUE)Checking code format...$(NC)"
	@UNFORMATTED=$$(gofmt -s -l . 2>&1); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "$(RED)The following files are not formatted:$(NC)"; \
		echo "$$UNFORMATTED"; \
		echo "$(CYAN)Run 'make fmt' to fix$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Code format check passed$(NC)"

.PHONY: check-imports
check-imports: ## Check if imports are formatted (CI mode, doesn't modify files)
	@echo "$(BLUE)Checking imports...$(NC)"
	@if [ -f "$(GOBIN)/goimports" ]; then \
		UNFORMATTED=$$($(GOBIN)/goimports -l . 2>&1); \
		if [ -n "$$UNFORMATTED" ]; then \
			echo "$(RED)The following files have import issues:$(NC)"; \
			echo "$$UNFORMATTED"; \
			echo "$(CYAN)Run 'make imports' to fix$(NC)"; \
			exit 1; \
		fi; \
		echo "$(GREEN)✅ Import check passed$(NC)"; \
	else \
		echo "$(YELLOW)Warning: goimports not installed$(NC)"; \
		echo "$(CYAN)Run 'make install-tools' to install goimports$(NC)"; \
		exit 1; \
	fi

.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@if [ -f "$(GOBIN)/golangci-lint" ]; then \
		$(GOBIN)/golangci-lint run; \
	else \
		echo "$(YELLOW)Warning: golangci-lint not installed, skipping linting$(NC)"; \
		echo "$(CYAN)Run 'make install-tools' to install golangci-lint$(NC)"; \
	fi
	@echo "$(GREEN)✅ Linting completed$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)✅ Go vet completed$(NC)"

# Release Commands
.PHONY: release
release: test build-all ## Create release build
	@echo "$(BLUE)Creating release build...$(NC)"
	@mkdir -p $(BUILD_DIR)/release
	@cp $(BUILD_DIR)/$(BINARY_NAME)-* $(BUILD_DIR)/release/
	@echo "$(GREEN)✅ Release build complete in $(BUILD_DIR)/release/$(NC)"

.PHONY: dist
dist: release ## Create distribution packages
	@echo "$(BLUE)Creating distribution packages...$(NC)"
	@cd $(BUILD_DIR)/release && \
		tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
		tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && \
		tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
		tar -czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
		zip $(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "$(GREEN)✅ Distribution packages created$(NC)"

# GoReleaser Commands
.PHONY: goreleaser-check
goreleaser-check: ## Validate GoReleaser configuration
	@echo "$(BLUE)Validating GoReleaser config...$(NC)"
	@goreleaser check
	@echo "$(GREEN)✅ GoReleaser config valid$(NC)"

.PHONY: goreleaser-snapshot
goreleaser-snapshot: ## Build snapshot release (doesn't publish)
	@echo "$(BLUE)Building snapshot release with GoReleaser...$(NC)"
	@goreleaser release --snapshot --clean
	@echo "$(GREEN)✅ Snapshot build complete in dist/$(NC)"

.PHONY: prepare-release
prepare-release: ## Prepare for release (run before creating GitHub release)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)❌ Error: VERSION not specified$(NC)"; \
		echo "Usage: make prepare-release VERSION=1.0.0"; \
		exit 1; \
	fi
	@echo "$(BLUE)Preparing release $(VERSION)...$(NC)"
	@echo "$(YELLOW)Running pre-release checks...$(NC)"
	@$(MAKE) ci
	@echo "$(GREEN)✅ Pre-release checks passed$(NC)"
	@echo ""
	@echo "$(CYAN)Next steps:$(NC)"
	@echo "1. Commit all changes"
	@echo "2. Push to your branch"
	@echo "3. Create and merge PR to master"
	@echo "4. Go to: https://github.com/molmedoz/gopher/actions/workflows/create-release.yml"
	@echo "5. Click 'Run workflow'"
	@echo "6. Enter version: $(VERSION)"
	@echo "7. Click 'Run workflow'"
	@echo ""
	@echo "$(GREEN)The workflow will validate, test, tag, build, and release automatically!$(NC)"

# Utility Commands
.PHONY: deps
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@$(GOMOD) download
	@echo "$(GREEN)✅ Dependencies downloaded$(NC)"

.PHONY: tidy
tidy: ## Tidy dependencies
	@echo "$(BLUE)Tidying dependencies...$(NC)"
	@$(GOMOD) tidy
	@echo "$(GREEN)✅ Dependencies tidied$(NC)"

.PHONY: check
check: fmt vet test ## Run all checks (development mode - modifies files)
	@echo "$(GREEN)✅ All checks completed$(NC)"

.PHONY: ci
ci: check-format check-imports vet lint test ## Run CI checks locally (matches GitHub Actions, no modifications)
	@echo "$(GREEN)✅ All CI checks passed - ready for push!$(NC)"

.PHONY: version
version: ## Show version information
	@echo "$(CYAN)Gopher Version Information$(NC)"
	@echo "$(YELLOW)=========================$(NC)"
	@echo "$(GREEN)Version:$(NC) $(VERSION)"
	@echo "$(GREEN)Go Version:$(NC) $(shell go version)"
	@echo "$(GREEN)Build Time:$(NC) $(shell date)"
	@echo "$(GREEN)Git Commit:$(NC) $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# Development setup
.PHONY: setup
setup: deps ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@$(MAKE) deps
	@echo "$(GREEN)✅ Development environment ready$(NC)"

# Quick development cycle
.PHONY: dev
dev: fmt vet test build ## Run development cycle
	@echo "$(GREEN)✅ Development cycle complete$(NC)"

# Docker Commands
.PHONY: test-docker
test-docker: ## Run unit tests in a container (volume-mounted, safe HOME)
	@echo "$(BLUE)Running tests in Docker (volume-mounted)...$(NC)"
	@docker run --rm -i \
		-v "$(PWD)":/workspace \
		-w /workspace \
		-e HOME=/tmp/gopher-home \
		-e GOPHER_CONFIG=/tmp/gopher-home/.gopher/config.json \
		golang:1.22-bullseye /bin/bash -lc "mkdir -p /tmp/gopher-home && /usr/local/go/bin/go version && /usr/local/go/bin/go test ./..."
	@echo "$(GREEN)✅ Docker tests completed$(NC)"

.PHONY: test-docker-coverage
test-docker-coverage: ## Run tests with coverage in a container (volume-mounted)
	@echo "$(BLUE)Running coverage in Docker (volume-mounted)...$(NC)"
	@docker run --rm -i \
		-v "$(PWD)":/workspace \
		-w /workspace \
		-e HOME=/tmp/gopher-home \
		-e GOPHER_CONFIG=/tmp/gopher-home/.gopher/config.json \
		golang:1.22-bullseye /bin/bash -lc "mkdir -p /tmp/gopher-home && /usr/local/go/bin/go test -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./... && /usr/local/go/bin/go tool cover -func=coverage.out | tail -n 1 && /usr/local/go/bin/go tool cover -html=coverage.out -o coverage.html"
	@echo "$(GREEN)✅ Docker coverage report generated: coverage.html$(NC)"

.PHONY: docker-build
docker-build: ## Build all Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	@./docker/build.sh build
	@echo "$(GREEN)✅ Docker images built$(NC)"

.PHONY: docker-test
docker-test: ## Run Docker tests
	@echo "$(BLUE)Running Docker tests...$(NC)"
	@./docker/build.sh test
	@echo "$(GREEN)✅ Docker tests completed$(NC)"

.PHONY: docker-clean
docker-clean: ## Clean Docker images
	@echo "$(BLUE)Cleaning Docker images...$(NC)"
	@./docker/build.sh clean
	@echo "$(GREEN)✅ Docker images cleaned$(NC)"

.PHONY: docker-test-unix-no-go
docker-test-unix-no-go: ## Test Unix no Go scenario
	@echo "$(BLUE)Testing Unix no Go scenario...$(NC)"
	@./docker/build.sh test-unix-no-go
	@echo "$(GREEN)✅ Unix no Go test completed$(NC)"

.PHONY: docker-test-unix-with-go
docker-test-unix-with-go: ## Test Unix with Go scenario
	@echo "$(BLUE)Testing Unix with Go scenario...$(NC)"
	@./docker/build.sh test-unix-with-go
	@echo "$(GREEN)✅ Unix with Go test completed$(NC)"

.PHONY: docker-test-windows-no-go
docker-test-windows-no-go: ## Test Windows no Go scenario
	@echo "$(BLUE)Testing Windows no Go scenario...$(NC)"
	@./docker/build.sh test-windows-no-go
	@echo "$(GREEN)✅ Windows no Go test completed$(NC)"

.PHONY: docker-test-windows-with-go
docker-test-windows-with-go: ## Test Windows with Go scenario
	@echo "$(BLUE)Testing Windows with Go scenario...$(NC)"
	@./docker/build.sh test-windows-with-go
	@echo "$(GREEN)✅ Windows with Go test completed$(NC)"

.PHONY: docker-test-windows-simulated
docker-test-windows-simulated: ## Test Windows simulation container
	@echo "$(BLUE)Testing Windows simulation container...$(NC)"
	@./docker/build.sh test-windows-simulated
	@echo "$(GREEN)✅ Windows simulation test completed$(NC)"

.PHONY: docker-test-macos-no-go
docker-test-macos-no-go: ## Test macOS no Go scenario
	@echo "$(BLUE)Testing macOS no Go scenario...$(NC)"
	@./docker/build.sh test-macos-no-go
	@echo "$(GREEN)✅ macOS no Go test completed$(NC)"

.PHONY: docker-test-macos-with-go
docker-test-macos-with-go: ## Test macOS with Go scenario
	@echo "$(BLUE)Testing macOS with Go scenario...$(NC)"
	@./docker/build.sh test-macos-with-go
	@echo "$(GREEN)✅ macOS with Go test completed$(NC)"

# Install development tools
.PHONY: install-tools
install-tools: ## Install all development tools (goimports, golangci-lint, etc.)
	@echo "$(BLUE)Installing development tools...$(NC)"
	@echo "$(YELLOW)Installing goimports...$(NC)"
	@$(GOINSTALL) golang.org/x/tools/cmd/goimports@latest
	@echo "$(YELLOW)Installing golangci-lint...$(NC)"
	@$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(YELLOW)Installing gofumpt...$(NC)"
	@$(GOINSTALL) mvdan.cc/gofumpt@latest
	@echo "$(GREEN)✅ Development tools installed$(NC)"
	@echo "$(CYAN)Installed tools:$(NC)"
	@echo "  - goimports (import formatting)"
	@echo "  - golangci-lint (comprehensive linting)"
	@echo "  - gofumpt (strict formatting)"

.PHONY: install-dev-tools
install-dev-tools: install-tools ## Alias for install-tools

# Security Commands
.PHONY: security-tools
security-tools: ## Install security tools
	@echo "$(BLUE)Installing security tools...$(NC)"
	@$(GOINSTALL) golang.org/x/vuln/cmd/govulncheck@latest
	@$(GOINSTALL) honnef.co/go/tools/cmd/staticcheck@latest
	@echo "$(YELLOW)Note: gosec installation skipped due to repository access issues$(NC)"
	@echo "$(GREEN)✅ Security tools installed$(NC)"

.PHONY: security-scan
security-scan: ## Run security scan
	@echo "$(BLUE)Running security scan...$(NC)"
	@if [ -f "$(GOBIN)/staticcheck" ]; then \
		echo "$(YELLOW)Running staticcheck...$(NC)"; \
		$(GOBIN)/staticcheck ./...; \
	else \
		echo "$(YELLOW)Warning: staticcheck not installed, skipping$(NC)"; \
	fi
	@echo "$(YELLOW)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Security scan completed$(NC)"

.PHONY: vuln-check
vuln-check: ## Check for vulnerabilities
	@echo "$(BLUE)Checking for vulnerabilities...$(NC)"
	@if [ -f "$(GOBIN)/govulncheck" ]; then \
		$(GOBIN)/govulncheck ./...; \
		echo "$(GREEN)✅ Vulnerability check completed$(NC)"; \
	else \
		echo "$(YELLOW)Warning: govulncheck not installed$(NC)"; \
		echo "$(CYAN)Run 'make security-tools' to install govulncheck$(NC)"; \
		exit 1; \
	fi

.PHONY: security-test
security-test: ## Run security tests
	@echo "$(BLUE)Running security tests...$(NC)"
	@$(GOTEST) -v -tags=security ./internal/security/...
	@echo "$(GREEN)✅ Security tests completed$(NC)"

.PHONY: security-all
security-all: security-tools security-scan vuln-check security-test ## Run all security checks
	@echo "$(GREEN)✅ All security checks completed$(NC)"
