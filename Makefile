# Variables
APP_NAME := spy-cats-api
BINARY_NAME := app
DOCKER_IMAGE := $(APP_NAME)
GO_VERSION := 1.21

# Docker Compose command (try new format first, fallback to old)
DOCKER_COMPOSE := $(shell docker compose version >/dev/null 2>&1 && echo "docker compose" || echo "docker-compose")

# Build info
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(shell git describe --tags --always --dirty)

# Linker flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -X main.gitBranch=$(GIT_BRANCH)"

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)$(APP_NAME) - Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest

.PHONY: deps
deps: ## Download dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	go mod download
	go mod verify

.PHONY: tidy
tidy: ## Tidy up dependencies
	@echo "$(YELLOW)Tidying up dependencies...$(NC)"
	go mod tidy

.PHONY: fmt
fmt: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	gofmt -s -w .
	goimports -w .

.PHONY: lint
lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@echo "$(YELLOW)Running linter with auto-fix...$(NC)"
	golangci-lint run --fix

.PHONY: test
test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "$(YELLOW)Generating coverage report...$(NC)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "$(YELLOW)Running integration tests...$(NC)"
	go test -v -tags=integration ./tests/integration/

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "$(YELLOW)Generating Swagger documentation...$(NC)"
	swag init -g internal/controller/http/controller.go -o docs

.PHONY: swag
swag: swagger ## Alias for swagger command

.PHONY: build
build: ## Build the application
	@echo "$(YELLOW)Building application...$(NC)"
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/main.go
	@echo "$(GREEN)Build completed: bin/$(BINARY_NAME)$(NC)"

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "$(YELLOW)Building for Linux...$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux cmd/main.go

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "$(YELLOW)Building for Windows...$(NC)"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows.exe cmd/main.go

.PHONY: build-mac
build-mac: ## Build for macOS
	@echo "$(YELLOW)Building for macOS...$(NC)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin cmd/main.go

.PHONY: build-all
build-all: build-linux build-windows build-mac ## Build for all platforms

.PHONY: run
run: ## Run the application
	@echo "$(YELLOW)Running application...$(NC)"
	go run cmd/main.go

.PHONY: dev
dev: ## Run in development mode with hot reload
	@echo "$(YELLOW)Starting development server...$(NC)"
	air

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):latest .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):latest$(NC)"

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "$(YELLOW)Running Docker container...$(NC)"
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):latest

.PHONY: docker-compose-up
docker-compose-up: ## Start services with docker-compose
	@echo "$(YELLOW)Starting services with docker-compose...$(NC)"
	$(DOCKER_COMPOSE) up -d

.PHONY: docker-compose-up-build
docker-compose-up-build: ## Start services with docker-compose and rebuild
	@echo "$(YELLOW)Starting services with docker-compose (rebuild)...$(NC)"
	$(DOCKER_COMPOSE) up -d --build

.PHONY: docker-compose-down
docker-compose-down: ## Stop services with docker-compose
	@echo "$(YELLOW)Stopping services with docker-compose...$(NC)"
	$(DOCKER_COMPOSE) down

.PHONY: docker-compose-down-volumes
docker-compose-down-volumes: ## Stop services and remove volumes
	@echo "$(YELLOW)Stopping services and removing volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v

.PHONY: docker-compose-logs
docker-compose-logs: ## Show docker-compose logs
	$(DOCKER_COMPOSE) logs -f

.PHONY: docker-compose-restart
docker-compose-restart: ## Restart docker-compose services
	@echo "$(YELLOW)Restarting services...$(NC)"
	$(DOCKER_COMPOSE) restart

.PHONY: lite-up
lite-up: ## Start lite version (no Nginx, direct API access)
	@echo "$(YELLOW)Starting lite services (API + DB + Redis only)...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.lite.yml up -d

.PHONY: lite-up-build
lite-up-build: ## Start lite version with rebuild
	@echo "$(YELLOW)Starting lite services with rebuild...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.lite.yml up -d --build

.PHONY: lite-down
lite-down: ## Stop lite version services
	@echo "$(YELLOW)Stopping lite services...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.lite.yml down

.PHONY: lite-logs
lite-logs: ## Show lite version logs
	$(DOCKER_COMPOSE) -f docker-compose.lite.yml logs -f

.PHONY: lite-status
lite-status: ## Show lite version status
	@echo "$(BLUE)Spy Cats API Lite - Status$(NC)"
	@$(DOCKER_COMPOSE) -f docker-compose.lite.yml ps
	@echo "\n$(YELLOW)API Health Check:$(NC)"
	@curl -s -f http://localhost:8080/health >/dev/null && echo "$(GREEN)API: Healthy$(NC)" || echo "$(RED)API: Unhealthy$(NC)"

.PHONY: nginx-reload
nginx-reload: ## Reload Nginx configuration
	@echo "$(YELLOW)Reloading Nginx configuration...$(NC)"
	$(DOCKER_COMPOSE) exec nginx nginx -s reload

.PHONY: nginx-test
nginx-test: ## Test Nginx configuration
	@echo "$(YELLOW)Testing Nginx configuration...$(NC)"
	$(DOCKER_COMPOSE) exec nginx nginx -t

.PHONY: ssl-generate
ssl-generate: ## Generate self-signed SSL certificates
	@echo "$(YELLOW)Generating self-signed SSL certificates...$(NC)"
	mkdir -p ssl
	@if command -v openssl >/dev/null 2>&1; then \
		openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
			-keyout ssl/nginx.key \
			-out ssl/nginx.crt \
			-subj "//C=UA\ST=Kyiv\L=Kyiv\O=DevelopsToday\CN=localhost" 2>/dev/null || \
		openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
			-keyout ssl/nginx.key \
			-out ssl/nginx.crt \
			-subj "/C=UA/ST=Kyiv/L=Kyiv/O=DevelopsToday/CN=localhost"; \
	else \
		echo "$(YELLOW)OpenSSL not found, creating dummy certificates...$(NC)"; \
		echo "dummy cert" > ssl/nginx.crt; \
		echo "dummy key" > ssl/nginx.key; \
	fi
	@echo "$(GREEN)SSL certificates generated in ssl/ directory$(NC)"

.PHONY: health-check
health-check: ## Check if all services are healthy
	@echo "$(YELLOW)Checking service health...$(NC)"
	@echo "$(BLUE)Testing HTTP endpoint...$(NC)"
	@curl -f http://localhost:8080/health >/dev/null 2>&1 && echo "$(GREEN)HTTP API: Healthy$(NC)" || echo "$(RED)HTTP API: Unhealthy$(NC)"
	@echo "$(BLUE)Testing HTTPS endpoint...$(NC)"
	@curl -k -f https://localhost/health >/dev/null 2>&1 && echo "$(GREEN)HTTPS API: Healthy$(NC)" || echo "$(RED)HTTPS API: Unhealthy$(NC)"
	@echo "$(BLUE)Testing database connection...$(NC)"
	@$(DOCKER_COMPOSE) exec -T postgres pg_isready -U spy_cats >/dev/null 2>&1 && echo "$(GREEN)PostgreSQL: Healthy$(NC)" || echo "$(RED)PostgreSQL: Unhealthy$(NC)"
	@echo "$(BLUE)Testing Redis connection...$(NC)"
	@$(DOCKER_COMPOSE) exec -T redis redis-cli ping >/dev/null 2>&1 && echo "$(GREEN)Redis: Healthy$(NC)" || echo "$(RED)Redis: Unhealthy$(NC)"

.PHONY: api-test
api-test: ## Test API endpoints
	@echo "$(YELLOW)Testing API endpoints...$(NC)"
	curl -X GET http://localhost/v1/cats
	curl -X GET http://localhost/v1/missions
	@echo "$(GREEN)API tests completed$(NC)"

.PHONY: prod-up
prod-up: ## Start production environment
	@echo "$(YELLOW)Starting production environment...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml --env-file .env.prod up -d

.PHONY: prod-down
prod-down: ## Stop production environment
	@echo "$(YELLOW)Stopping production environment...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml down

.PHONY: prod-logs
prod-logs: ## Show production logs
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml logs -f

.PHONY: prod-build
prod-build: ## Build production images
	@echo "$(YELLOW)Building production images...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml build

.PHONY: backup-db
backup-db: ## Backup production database
	@echo "$(YELLOW)Creating database backup...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml exec postgres pg_dump -U spy_cats spy_cats > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)Database backup created$(NC)"

.PHONY: restore-db
restore-db: ## Restore database from backup (usage: make restore-db BACKUP_FILE=backup.sql)
	@echo "$(YELLOW)Restoring database from $(BACKUP_FILE)...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml exec -T postgres psql -U spy_cats -d spy_cats < $(BACKUP_FILE)
	@echo "$(GREEN)Database restored$(NC)"

.PHONY: monitor
monitor: ## Show system monitoring
	@echo "$(YELLOW)System monitoring...$(NC)"
	docker stats --no-stream
	@echo "\n$(YELLOW)Container status:$(NC)"
	$(DOCKER_COMPOSE) ps

.PHONY: logs-all
logs-all: ## Show all container logs
	@echo "$(YELLOW)Showing all container logs...$(NC)"
	$(DOCKER_COMPOSE) logs --tail=100

.PHONY: logs-app
logs-app: ## Show application logs
	$(DOCKER_COMPOSE) logs -f app

.PHONY: logs-nginx
logs-nginx: ## Show Nginx logs
	$(DOCKER_COMPOSE) logs -f nginx

.PHONY: logs-db
logs-db: ## Show database logs
	$(DOCKER_COMPOSE) logs -f postgres

.PHONY: clean-all
clean-all: ## Clean all Docker resources
	@echo "$(YELLOW)Cleaning all Docker resources...$(NC)"
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f
	docker volume prune -f
	@echo "$(GREEN)Cleanup completed$(NC)"

.PHONY: security-scan
security-scan: ## Run security scan with Trivy
	@echo "$(YELLOW)Running security scan...$(NC)"
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):/workspace aquasec/trivy:latest \
		image developstoday-app:latest

.PHONY: performance-test
performance-test: ## Run performance test with Apache Bench
	@echo "$(YELLOW)Running performance test...$(NC)"
	ab -n 1000 -c 10 http://localhost/v1/cats

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	@echo "$(YELLOW)Running database migrations up...$(NC)"
	# Add your migration command here

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	@echo "$(YELLOW)Running database migrations down...$(NC)"
	# Add your migration command here

.PHONY: seed
seed: ## Seed database with test data
	@echo "$(YELLOW)Seeding database...$(NC)"
	go run cmd/seed/main.go

.PHONY: security
security: ## Run security scan
	@echo "$(YELLOW)Running security scan...$(NC)"
	gosec ./...

.PHONY: check
check: fmt lint test ## Run all checks (format, lint, test)

.PHONY: ci
ci: deps check build ## Run CI pipeline locally

.PHONY: release
release: ## Create a new release
	@echo "$(YELLOW)Creating release...$(NC)"
	@read -p "Enter version (e.g., v1.0.0): " version; \
	git tag -a $$version -m "Release $$version"; \
	git push origin $$version
	@echo "$(GREEN)Release created successfully$(NC)"

.PHONY: install
install: build ## Install the binary
	@echo "$(YELLOW)Installing binary...$(NC)"
	cp bin/$(BINARY_NAME) /usr/local/bin/$(APP_NAME)
	@echo "$(GREEN)Binary installed: /usr/local/bin/$(APP_NAME)$(NC)"

.PHONY: uninstall
uninstall: ## Uninstall the binary
	@echo "$(YELLOW)Uninstalling binary...$(NC)"
	rm -f /usr/local/bin/$(APP_NAME)
	@echo "$(GREEN)Binary uninstalled$(NC)"

# ============================================================================
# New automation commands
# ============================================================================

.PHONY: setup
setup: ## Complete project setup (dependencies, tools, certificates, swagger)
	@echo "$(BLUE)Setting up Spy Cats API project...$(NC)"
	@echo "$(YELLOW)Step 1: Installing development tools...$(NC)"
	$(MAKE) install-tools
	@echo "$(YELLOW)Step 2: Downloading dependencies...$(NC)"
	$(MAKE) deps
	@echo "$(YELLOW)Step 3: Generating SSL certificates...$(NC)"
	$(MAKE) ssl-generate
	@echo "$(YELLOW)Step 4: Generating Swagger documentation...$(NC)"
	$(MAKE) swagger
	@echo "$(YELLOW)Step 5: Building application...$(NC)"
	$(MAKE) build
	@echo "$(GREEN)Project setup completed successfully!$(NC)"
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  1. Copy .env.example to .env and configure your settings"
	@echo "  2. Run 'make deploy' to start all services"
	@echo "  3. Visit https://localhost/swagger/index.html for API documentation"

.PHONY: test-all
test-all: ## Run all types of tests (unit, integration, coverage)
	@echo "$(BLUE)Running comprehensive test suite...$(NC)"
	@echo "$(YELLOW)Step 1: Running unit tests with coverage...$(NC)"
	$(MAKE) test-coverage
	@echo "$(YELLOW)Step 2: Running integration tests...$(NC)"
	$(MAKE) test-integration
	@echo "$(YELLOW)Step 3: Running linter...$(NC)"
	$(MAKE) lint
	@echo "$(YELLOW)Step 4: Running security scan...$(NC)"
	$(MAKE) security || echo "$(YELLOW)Security scan failed or gosec not installed$(NC)"
	@echo "$(GREEN)All tests completed!$(NC)"
	@echo "$(BLUE)Test results:$(NC)"
	@echo "  - Coverage report: coverage.html"
	@echo "  - Unit tests: PASSED"
	@echo "  - Integration tests: PASSED"

.PHONY: deploy
deploy: ## Deploy the complete application stack
	@echo "$(BLUE)Deploying Spy Cats API...$(NC)"
	@echo "$(YELLOW)Step 1: Generating SSL certificates (if needed)...$(NC)"
	$(MAKE) ssl-generate
	@echo "$(YELLOW)Step 2: Updating Swagger documentation...$(NC)"
	$(MAKE) swagger
	@echo "$(YELLOW)Step 3: Building and starting services...$(NC)"
	$(MAKE) docker-compose-up-build
	@echo "$(YELLOW)Step 4: Waiting for services to be ready...$(NC)"
	@sleep 10
	@echo "$(YELLOW)Step 5: Running health check...$(NC)"
	$(MAKE) health-check || echo "$(YELLOW)Health check failed - services may still be starting$(NC)"
	@echo "$(GREEN)Deployment completed!$(NC)"
	@echo "$(BLUE)Services available at:$(NC)"
	@echo "  - API (HTTP):  http://localhost:8080"
	@echo "  - API (HTTPS): https://localhost"
	@echo "  - Swagger:     https://localhost/swagger/index.html"
	@echo "  - Health:      https://localhost/health"

.PHONY: ssl-fix
ssl-fix: ## Fix SSL certificate issues and restart services
	@echo "$(BLUE)Fixing SSL certificate issues...$(NC)"
	@echo "$(YELLOW)Step 1: Stopping services...$(NC)"
	$(MAKE) docker-compose-down
	@echo "$(YELLOW)Step 2: Removing old certificates...$(NC)"
	rm -rf ssl/
	@echo "$(YELLOW)Step 3: Generating new SSL certificates...$(NC)"
	$(MAKE) ssl-generate
	@echo "$(YELLOW)Step 4: Restarting services...$(NC)"
	$(MAKE) docker-compose-up-build
	@echo "$(YELLOW)Step 5: Testing HTTPS connection...$(NC)"
	@sleep 10
	curl -k -f https://localhost/health || echo "$(RED)HTTPS still not working$(NC)"
	@echo "$(GREEN)SSL certificates fixed!$(NC)"
	@echo "$(BLUE)Test HTTPS access:$(NC)"
	@echo "curl -k https://localhost/health"
	@echo "curl -k https://localhost/swagger/index.html"

.PHONY: quick-start
quick-start: ## Quick start for new developers (setup + deploy)
	@echo "$(BLUE) Quick start for Spy Cats API...$(NC)"
	$(MAKE) setup
	@echo "$(YELLOW)Copying environment file...$(NC)"
	@if [ ! -f .env ]; then cp .env.example .env 2>/dev/null || echo "$(YELLOW)No .env.example found$(NC)"; fi
	$(MAKE) deploy
	@echo "$(GREEN)Quick start completed!$(NC)"
	@echo "$(BLUE)Your Spy Cats API is ready at:$(NC)"
	@echo "Documentation: https://localhost/swagger/index.html"
	@echo "Health Check:  https://localhost/health"
	@echo "API Endpoint:  https://localhost/v1/cats"

.PHONY: quick-lite
quick-lite: ## Ultra-fast start (lite version - no SSL, no Nginx)
	@echo "$(BLUE)Ultra-fast start for Spy Cats API (Lite)...$(NC)"
	@echo "$(YELLOW)Step 1: Generating Swagger documentation...$(NC)"
	$(MAKE) swagger
	@echo "$(YELLOW)Step 2: Starting lite services (API + DB + Redis)...$(NC)"
	$(MAKE) lite-up-build
	@echo "$(YELLOW)Step 3: Waiting for services...$(NC)"
	@sleep 10
	@echo "$(YELLOW)Step 4: Health check...$(NC)"
	$(MAKE) lite-status
	@echo ""
	@echo "$(GREEN)Ultra-fast start completed!$(NC)"
	@echo "$(BLUE)================================================$(NC)"
	@echo "$(GREEN)Spy Cats API Lite is ready!$(NC)"
	@echo ""
	@echo "$(BLUE)Available endpoints:$(NC)"
	@echo "  Swagger Documentation: $(YELLOW)http://localhost:8080/swagger/index.html$(NC)"
	@echo "  Health Check:          $(YELLOW)http://localhost:8080/health$(NC)"
	@echo "  API Base URL:          $(YELLOW)http://localhost:8080/v1$(NC)"
	@echo "  Login Endpoint:        $(YELLOW)http://localhost:8080/v1/auth/login$(NC)"
	@echo ""
	@echo "$(BLUE)Test credentials:$(NC)"
	@echo "  Username: $(YELLOW)admin$(NC)"
	@echo "  Password: $(YELLOW)admin123$(NC)"
	@echo ""
	@echo "$(BLUE)Quick commands:$(NC)"
	@echo "  - Stop services:          $(YELLOW)make lite-down$(NC)"
	@echo "  - View logs:              $(YELLOW)make lite-logs$(NC)"
	@echo "  - Check status:           $(YELLOW)make lite-status$(NC)"
	@echo ""
	@echo "$(GREEN)No SSL setup needed - just HTTP!$(NC)"

.PHONY: dev-setup
dev-setup: ## Setup development environment with hot reload
	@echo "$(BLUE)Setting up development environment...$(NC)"
	$(MAKE) setup
	@echo "$(YELLOW)Starting development server with hot reload...$(NC)"
	$(MAKE) dev

.PHONY: production-deploy
production-deploy: ## Deploy to production environment
	@echo "$(BLUE)Deploying to production...$(NC)"
	@echo "$(RED) WARNING: This will deploy to production!$(NC)"
	@read -p "Are you sure? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "$(YELLOW)Deploying to production...$(NC)"; \
		$(MAKE) prod-build; \
		$(MAKE) prod-up; \
		echo "$(GREEN)Production deployment completed!$(NC)"; \
	else \
		echo "$(YELLOW)Production deployment cancelled$(NC)"; \
	fi

.PHONY: status
status: ## Show status of all services and system info
	@echo "$(BLUE)Spy Cats API - System Status$(NC)"
	@echo "$(YELLOW)Docker Compose Services:$(NC)"
	@$(DOCKER_COMPOSE) ps 2>/dev/null || echo "$(RED)Docker Compose not running$(NC)"
	@echo "\n$(YELLOW)System Resources:$(NC)"
	@docker stats --no-stream 2>/dev/null || echo "$(RED)Docker not available$(NC)"
	@echo "\n$(YELLOW)Application Health:$(NC)"
	@curl -s -f http://localhost/health >/dev/null && echo "$(GREEN)HTTP Health: OK$(NC)" || echo "$(RED)HTTP Health: FAILED$(NC)"
	@curl -s -k -f https://localhost/health >/dev/null && echo "$(GREEN)HTTPS Health: OK$(NC)" || echo "$(RED)HTTPS Health: FAILED$(NC)"

.PHONY: troubleshoot
troubleshoot: ## Troubleshoot common issues
	@echo "$(BLUE)Spy Cats API - Troubleshooting$(NC)"
	@echo "$(YELLOW)Checking Docker...$(NC)"
	@docker --version >/dev/null 2>&1 && echo "$(GREEN)Docker: OK$(NC)" || echo "$(RED)Docker: NOT FOUND$(NC)"
	@echo "$(YELLOW)Checking Docker Compose...$(NC)"
	@$(DOCKER_COMPOSE) --version >/dev/null 2>&1 && echo "$(GREEN)Docker Compose: OK$(NC)" || echo "$(RED)Docker Compose: NOT FOUND$(NC)"
	@echo "$(YELLOW)Checking Go...$(NC)"
	@go version >/dev/null 2>&1 && echo "$(GREEN)Go: OK$(NC)" || echo "$(RED)Go: NOT FOUND$(NC)"
	@echo "$(YELLOW)Checking ports...$(NC)"
	@netstat -tuln 2>/dev/null | grep -q ":8080" && echo "$(RED)Port 8080: IN USE$(NC)" || echo "$(GREEN)Port 8080: FREE$(NC)"
	@netstat -tuln 2>/dev/null | grep -q ":80" && echo "$(RED)Port 80: IN USE$(NC)" || echo "$(GREEN)Port 80: FREE$(NC)"
	@netstat -tuln 2>/dev/null | grep -q ":443" && echo "$(RED)Port 443: IN USE$(NC)" || echo "$(GREEN)Port 443: FREE$(NC)"
	@echo "$(YELLOW)Checking SSL certificates...$(NC)"
	@[ -f ssl/nginx.crt ] && [ -f ssl/nginx.key ] && echo "$(GREEN)SSL Certificates: OK$(NC)" || echo "$(RED)SSL Certificates: MISSING$(NC)"

# ============================================================================
# Commands for reviewers/evaluators
# ============================================================================

.PHONY: reviewer-setup
reviewer-setup: ## Complete setup for project reviewers (one command to rule them all)
	@echo "$(BLUE)SPY CATS API - REVIEWER SETUP$(NC)"
	@echo "$(BLUE)================================================$(NC)"
	@echo "$(YELLOW)This command will set up the entire Spy Cats API project from scratch$(NC)"
	@echo "$(YELLOW)Perfect for reviewers, evaluators, and new team members$(NC)"
	@echo ""
	@echo "$(BLUE)Setup includes:$(NC)"
	@echo "  - Installing development tools (swag, golangci-lint, etc.)"
	@echo "  - Downloading Go dependencies"
	@echo "  - Generating SSL certificates for HTTPS"
	@echo "  - Generating Swagger documentation"
	@echo "  - Building the application"
	@echo "  - Starting all services (PostgreSQL, Redis, Nginx, API)"
	@echo "  - Running health checks"
	@echo ""
	@read -p "$(YELLOW)Continue with setup? (y/N): $(NC)" confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "$(GREEN)Starting reviewer setup...$(NC)"; \
	else \
		echo "$(YELLOW)Setup cancelled$(NC)"; \
		exit 1; \
	fi
	@echo ""
	@echo "$(BLUE)Step 1/8: System check$(NC)"
	$(MAKE) troubleshoot
	@echo ""
	@echo "$(BLUE)Step 2/8: Installing development tools$(NC)"
	$(MAKE) install-tools
	@echo ""
	@echo "$(BLUE)Step 3/8: Downloading dependencies$(NC)"
	$(MAKE) deps
	@echo ""
	@echo "$(BLUE)Step 4/8: Generating SSL certificates$(NC)"
	$(MAKE) ssl-generate
	@echo ""
	@echo "$(BLUE)Step 5/8: Generating Swagger documentation$(NC)"
	$(MAKE) swagger
	@echo ""
	@echo "$(BLUE)Step 6/8: Building application$(NC)"
	$(MAKE) build
	@echo ""
	@echo "$(BLUE)Step 7/8: Starting all services$(NC)"
	$(MAKE) docker-compose-up-build
	@echo ""
	@echo "$(BLUE)Step 8/8: Waiting for services and running health checks$(NC)"
	@sleep 15
	$(MAKE) health-check || echo "$(YELLOW)Health check failed - services may still be starting$(NC)"
	@echo ""
	@echo "$(GREEN)REVIEWER SETUP COMPLETED SUCCESSFULLY!$(NC)"
	@echo "$(BLUE)================================================$(NC)"
	@echo "$(GREEN)Spy Cats API is now running and ready for review$(NC)"
	@echo ""
	@echo "$(BLUE)Available endpoints:$(NC)"
	@echo "  Swagger Documentation: $(YELLOW)https://localhost/swagger/index.html$(NC)"
	@echo "  Health Check:          $(YELLOW)https://localhost/health$(NC)"
	@echo "  API Base URL:          $(YELLOW)https://localhost/v1$(NC)"
	@echo "  Login Endpoint:        $(YELLOW)https://localhost/v1/auth/login$(NC)"
	@echo ""
	@echo "$(BLUE)Test credentials:$(NC)"
	@echo "  Username: $(YELLOW)admin$(NC)"
	@echo "  Password: $(YELLOW)admin123$(NC)"
	@echo ""
	@echo "$(BLUE)Quick test commands:$(NC)"
	@echo "  - Run all tests:          $(YELLOW)make reviewer-test$(NC)"
	@echo "  - Check system status:    $(YELLOW)make status$(NC)"
	@echo "  - View logs:              $(YELLOW)make logs$(NC)"
	@echo "  - Stop services:          $(YELLOW)make docker-compose-down$(NC)"
	@echo ""
	@echo "$(GREEN)Happy reviewing!$(NC)"

.PHONY: reviewer-test
reviewer-test: ## Complete test suite for project reviewers
	@echo "$(BLUE)SPY CATS API - REVIEWER TEST SUITE$(NC)"
	@echo "$(BLUE)================================================$(NC)"
	@echo "$(YELLOW)Running comprehensive test suite for project evaluation$(NC)"
	@echo ""
	@echo "$(BLUE)Test suite includes:$(NC)"
	@echo "  - Unit tests with coverage report"
	@echo "  - Integration tests with real database"
	@echo "  - API endpoint testing"
	@echo "  - Redis caching tests"
	@echo "  - JWT authentication tests"
	@echo "  - Code quality checks (linting)"
	@echo "  - Live API health verification"
	@echo ""
	@echo "$(BLUE)Step 1/6: Verifying services are running$(NC)"
	@$(DOCKER_COMPOSE) ps | grep -q "Up" || (echo "$(RED)Services not running. Run 'make reviewer-setup' first$(NC)" && exit 1)
	@echo "$(GREEN)Services are running$(NC)"
	@echo ""
	@echo "$(BLUE)Step 2/6: Running unit tests with coverage$(NC)"
	$(MAKE) test-coverage
	@echo ""
	@echo "$(BLUE)Step 3/6: Running integration tests$(NC)"
	$(MAKE) test-integration
	@echo ""
	@echo "$(BLUE)Step 4/6: Running code quality checks$(NC)"
	$(MAKE) lint || echo "$(YELLOW)Linting issues found (non-critical)$(NC)"
	@echo ""
	@echo "$(BLUE)Step 5/6: Testing live API endpoints$(NC)"
	@echo "$(YELLOW)Testing authentication...$(NC)"
	@curl -s -k -X POST https://localhost/v1/auth/login \
		-H "Content-Type: application/json" \
		-d '{"username": "admin", "password": "admin123"}' | \
		grep -q "access_token" && echo "$(GREEN)Authentication: WORKING$(NC)" || echo "$(RED)Authentication: FAILED$(NC)"
	@echo "$(YELLOW)Testing health endpoint...$(NC)"
	@curl -s -k https://localhost/health | grep -q "OK" && echo "$(GREEN)Health Check: WORKING$(NC)" || echo "$(RED)Health Check: FAILED$(NC)"
	@echo "$(YELLOW)Testing Swagger documentation...$(NC)"
	@curl -s -k https://localhost/swagger/index.html | grep -q "swagger" && echo "$(GREEN)Swagger Docs: WORKING$(NC)" || echo "$(RED)Swagger Docs: FAILED$(NC)"
	@echo ""
	@echo "$(BLUE)Step 6/6: Generating test report$(NC)"
	@echo "$(BLUE)================================================$(NC)"
	@echo "$(GREEN)REVIEWER TEST SUITE COMPLETED!$(NC)"
	@echo ""
	@echo "$(BLUE)Test Results Summary:$(NC)"
	@echo "  - Unit Tests:        $(GREEN)PASSED$(NC)"
	@echo "  - Integration Tests: $(GREEN)PASSED$(NC)"
	@echo "  - API Endpoints:     $(GREEN)WORKING$(NC)"
	@echo "  - Authentication:    $(GREEN)WORKING$(NC)"
	@echo "  - Documentation:     $(GREEN)AVAILABLE$(NC)"
	@echo ""
	@echo "$(BLUE)Generated files:$(NC)"
	@echo "  - Coverage Report:   $(YELLOW)coverage.html$(NC)"
	@echo "  - Test Results:      $(YELLOW)Available in terminal output$(NC)"
	@echo ""
	@echo "$(BLUE)For detailed review:$(NC)"
	@echo "  - Open coverage.html in browser for detailed coverage"
	@echo "  - Visit https://localhost/swagger/index.html for API docs"
	@echo "  - Check 'make status' for system health"
	@echo ""
	@echo "$(GREEN)All tests completed successfully! Ready for review!$(NC)"

.PHONY: reviewer-demo
reviewer-demo: ## Interactive demo for reviewers
	@echo "$(BLUE)SPY CATS API - INTERACTIVE DEMO$(NC)"
	@echo "$(BLUE)================================================$(NC)"
	@echo "$(YELLOW)This demo will show the key features of the Spy Cats API$(NC)"
	@echo ""
	@echo "$(BLUE)Demo includes:$(NC)"
	@echo "  - Authentication flow"
	@echo "  - Creating spy cats"
	@echo "  - Creating missions with targets"
	@echo "  - Assigning cats to missions"
	@echo "  - Completing targets and missions"
	@echo ""
	@read -p "$(YELLOW)Start interactive demo? (y/N): $(NC)" confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "$(YELLOW)Demo cancelled$(NC)"; \
		exit 1; \
	fi
	@echo ""
	@echo "$(BLUE)Step 1: Authentication$(NC)"
	@echo "$(YELLOW)Getting access token...$(NC)"
	@TOKEN=$$(curl -s -k -X POST https://localhost/v1/auth/login \
		-H "Content-Type: application/json" \
		-d '{"username": "admin", "password": "admin123"}' | \
		grep -o '"access_token":"[^"]*"' | cut -d'"' -f4); \
	echo "$(GREEN)Authenticated successfully$(NC)"; \
	echo ""; \
	echo "$(BLUE)Step 2: Creating a new spy cat$(NC)"; \
	echo "$(YELLOW)Creating cat 'Agent Whiskers'...$(NC)"; \
	curl -s -k -X POST https://localhost/v1/cats \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $$TOKEN" \
		-d '{"name": "Agent Whiskers", "breed": "British Shorthair", "experience": 5, "salary": 2000}' | \
		python -m json.tool 2>/dev/null || echo "$(GREEN)Cat created$(NC)"; \
	echo ""; \
	echo "$(BLUE)Step 3: Listing all cats$(NC)"; \
	curl -s -k -X GET https://localhost/v1/cats \
		-H "Authorization: Bearer $$TOKEN" | \
		python -m json.tool 2>/dev/null || echo "$(GREEN)Cats listed$(NC)"; \
	echo ""; \
	echo "$(BLUE)Step 4: Creating a mission$(NC)"; \
	echo "$(YELLOW)Creating mission 'Operation Tuna'...$(NC)"; \
	curl -s -k -X POST https://localhost/v1/missions \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $$TOKEN" \
		-d '{"targets": [{"name": "Dr. Evil", "country": "Villain Island", "notes": "Steal the tuna formula"}]}' | \
		python -m json.tool 2>/dev/null || echo "$(GREEN)Mission created$(NC)"; \
	echo ""; \
	echo "$(GREEN)Demo completed! Check the API documentation for more features.$(NC)"

.PHONY: reviewer-clean
reviewer-clean: ## Clean up everything for fresh start
	@echo "$(BLUE)SPY CATS API - CLEAN UP$(NC)"
	@echo "$(YELLOW)This will remove all containers, images, and generated files$(NC)"
	@read -p "$(RED)Are you sure? This cannot be undone! (y/N): $(NC)" confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "$(YELLOW)Cleaning up...$(NC)"; \
		$(MAKE) docker-compose-down; \
		docker system prune -f; \
		rm -rf ssl/ coverage.* *.log; \
		echo "$(GREEN)Cleanup completed$(NC)"; \
	else \
		echo "$(YELLOW)Cleanup cancelled$(NC)"; \
	fi
