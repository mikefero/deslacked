# Build tooling

REPO_PATH       := github.com/mikefero/deslacked
VERSION_PKG     := $(REPO_PATH)/internal/version
BINARY_NAME     ?= deslacked
GO_BUILD_OUTPUT ?= $(BIN_DIR)/$(BINARY_NAME)

.PHONY: build
build: ## Build the deslacked binary
	@echo "Building deslacked binary..."
	@mkdir -p "$(BIN_DIR)"
	@CGO_ENABLED=0 go build \
		-ldflags="-X $(VERSION_PKG).AppVersion=$(APP_VERSION) \
			-X $(VERSION_PKG).AppCommit=$(GIT_COMMIT) \
			-X $(VERSION_PKG).BuildDate=$(BUILD_DATETIME)" \
		-o $(GO_BUILD_OUTPUT) ./main.go
	@echo "✓ Binary: $(GO_BUILD_OUTPUT)"
