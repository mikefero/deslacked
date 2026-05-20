# Tools

GOLANGCI_LINT_VERSION ?= v2.12.2
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

.PHONY: install-tools
install-tools: ## Install required development tools
	@mkdir -p "$(BIN_DIR)"
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@curl -sSfL \
		"https://golangci-lint.run/install.sh" \
		| sh -s -- -b "$(BIN_DIR)" "$(GOLANGCI_LINT_VERSION)"
	@echo "✓ Tools installed to $(BIN_DIR)/"

.PHONY: check-golangci-lint
check-golangci-lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "'golangci-lint' not found — run 'make install-tools'"; \
		exit 1; \
	fi

.PHONY: lint-go
lint-go: check-golangci-lint ## Lint Go source files with golangci-lint
	@echo "Linting Go source files..."
	@"$(GOLANGCI_LINT)" run ./...
	@echo "✓ Go lint passed!"

.PHONY: check-docker
check-docker:
	@if ! command -v docker &> /dev/null; then \
		echo "'docker' not found — install Docker to run lint-web"; \
		exit 1; \
	fi

.PHONY: lint-web
lint-web: check-docker ## Lint JavaScript, HTML, and CSS files via Docker
	@echo "Linting web files..."
	@docker run --rm \
		-v "$(PROJECT_ROOT):/app" \
		-w /app \
		node:lts-alpine \
		sh -c "npm install --silent && \
			npx eslint internal/web/static/ && \
			npx html-validate internal/web/static/*.html && \
			npx stylelint internal/web/static/*.css"
	@echo "✓ Web lint passed!"

.PHONY: lint
lint: lint-go lint-web ## Lint all source files (Go and JavaScript)

.PHONY: format
format: check-golangci-lint ## Apply gofmt/goimports formatting to Go source files
	@echo "Formatting Go source files..."
	@"$(GOLANGCI_LINT)" fmt
	@echo "✓ Go formatting applied!"

.PHONY: go-mod-upgrade
go-mod-upgrade: ## Upgrade go modules
	@go get -u ./...
	@go mod tidy
	@go mod verify > /dev/null
