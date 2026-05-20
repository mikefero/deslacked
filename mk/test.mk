# Test tooling

.PHONY: test
test: ## Run Go unit tests with race detector
	@echo "Running Go unit tests..."
	@go test -race ./...
	@echo "✓ Go unit tests passed!"

.PHONY: test-no-race
test-no-race: ## Run Go unit tests without race detector
	@echo "Running Go unit tests (no race detector)..."
	@go test ./...
	@echo "✓ Go unit tests passed!"

.PHONY: test-verbose
test-verbose: ## Run Go unit tests with race detector, verbose
	@echo "Running Go unit tests (verbose)..."
	@go test -race ./... -v

.PHONY: test-verbose-no-race
test-verbose-no-race: ## Run Go unit tests verbose, no race detector
	@echo "Running Go unit tests (verbose, no race detector)..."
	@go test ./... -v

.PHONY: test-clean-cache
test-clean-cache: ## Clear the Go test cache
	@echo "Clearing Go test cache..."
	@go clean -testcache
	@echo "✓ Go test cache cleared!"
