SHELL := /bin/bash
GO ?= go

# Color output
BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

BUILD_DIR := $(CURDIR)/build

.PHONY: test-coverage-check
test-coverage-check: MIN_COVERAGE ?= 80
test-coverage-check: ## Check if coverage meets minimum threshold. Usage: make test-coverage-check MIN_COVERAGE=80
	@echo -e "$(BLUE)Checking coverage (minimum: $(MIN_COVERAGE)%)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) test -coverprofile=$(BUILD_DIR)/coverage.out ./...
	@coverage=$$($(GO) tool cover -func=$(BUILD_DIR)/coverage.out | grep total | awk '{print $$3}' | tr -d '%'); \
	echo "Current coverage: $$coverage%"; \
	if awk "BEGIN {exit !($$coverage >= $(MIN_COVERAGE))}"; then \
		echo -e "$(GREEN)✓ Coverage meets minimum threshold$(NC)"; \
	else \
		echo -e "$(RED)✗ Coverage $$coverage% is below minimum $(MIN_COVERAGE)%$(NC)"; \
		exit 1; \
	fi

test-coverage: ## Run tests with coverage report
	@echo -e "$(BLUE)Running tests with coverage...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) test -v -race -coverprofile=$(BUILD_DIR)/coverage.out -covermode=atomic ./...
	@$(GO) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo -e "$(GREEN)✓ Coverage report generated: $(BUILD_DIR)/coverage.html$(NC)"
	@$(GO) tool cover -func=$(BUILD_DIR)/coverage.out | grep total | awk '{print "Total coverage: " $$3}'


.PHONY: audit
audit:
	go fmt ./...
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...
	go test -race -vet=off ./...

.PHONY: code-coverage
code-coverage:
	go test -v -coverprofile /tmp/cover.out ./...
	go tool cover -html /tmp/cover.out -o /tmp/cover.html
	xdg-open /tmp/cover.html

.PHONY: benchmark
benchmark:
	go test -benchmem -bench=. ./...
