SHELL := /bin/bash
.DEFAULT_GOAL := help

# Paths
PROJECT_ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BIN_DIR      := $(PROJECT_ROOT)bin

# Image metadata (populated from git/date for local builds; overridden by GHA)
APP_VERSION    ?= $(shell cat version 2>/dev/null || echo v0.0.0-alpha.1.dev)
BUILD_DATETIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT     ?= $(shell git rev-parse --verify HEAD 2>/dev/null || echo unknown)

# Use https://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=deslacked to generate
# ANSI Shadow title for application name
.PHONY: help
help: ## Display this help screen
	@echo ''
	@echo '██████╗ ███████╗███████╗██╗      █████╗  ██████╗██╗  ██╗███████╗██████╗ ';
	@echo '██╔══██╗██╔════╝██╔════╝██║     ██╔══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗';
	@echo '██║  ██║█████╗  ███████╗██║     ███████║██║     █████╔╝ █████╗  ██║  ██║';
	@echo '██║  ██║██╔══╝  ╚════██║██║     ██╔══██║██║     ██╔═██╗ ██╔══╝  ██║  ██║';
	@echo '██████╔╝███████╗███████║███████╗██║  ██║╚██████╗██║  ██╗███████╗██████╔╝';
	@echo '╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═════╝ ';
	@echo ''
	@echo "deslacked — turning your coworkers into ghosts, one profile pic at a time"
	@echo ""
	@echo "Usage: make [target] [VARIABLE=value]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) mk/*.mk | sed 's/^[^:]*://' | sort -u | awk 'BEGIN {FS = ":.*?## "}; !seen[$$1]++ {printf "  \033[36m%-26s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Variables:"
	@echo "  APP_VERSION              Application version (default: $(APP_VERSION))"
	@echo "  BINARY_NAME              Output binary name (default: $(BINARY_NAME))"

# Include build, test, and tools targets
include mk/build.mk
include mk/test.mk
include mk/tools.mk
