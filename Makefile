
BINARY := gosh
RUN_DIR := .tmp

.PHONY: build
build:
	@go build -o $(BINARY) app/*.go
	@echo "Built $(BINARY)"

.PHONY: run
run:
	@mkdir -p $(RUN_DIR)
	@go build -o $(RUN_DIR)/$(BINARY) app/*.go
	@./$(RUN_DIR)/$(BINARY) "$$@"

.PHONY: fmt
fmt:
	@go fmt ./...
	@echo "Formatted"

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...
# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint the code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build         - Build gosh binary"
	@echo "  make run           - Build and run gosh"
	@echo "  make fmt           - Format code (go fmt)"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make lint          - Lint the code (requires golangci-lint)"
	@echo "  make help          - Show this help message"

