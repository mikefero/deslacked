# Build tooling

REPO_PATH         := github.com/mikefero/deslacked
VERSION_PKG       := $(REPO_PATH)/internal/version
BINARY_NAME       ?= deslacked
GO_BUILD_OUTPUT   ?= $(BIN_DIR)/$(BINARY_NAME)
DOCKER_IMAGE_NAME ?= deslacked

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

.PHONY: build-docker
build-docker: check-docker ## Build Docker image for the local platform
	@echo "Building deslacked Docker image..."
	@docker buildx build \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATETIME) \
		-t $(DOCKER_IMAGE_NAME):$(APP_VERSION) \
		-t $(DOCKER_IMAGE_NAME):latest \
		--load \
		.
	@echo "✓ Image: $(DOCKER_IMAGE_NAME):$(APP_VERSION)"
