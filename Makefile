# ==============================
# Makefile - Go Project CI/CD
# ==============================

.PHONY: all tidy build test coverage lint security clean

# Go files excluding vendor
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")
COVERAGE_FILE := coverage.out
LINT_REPORT := golangci-lint-report.txt
SECURITY_REPORT := gosec-report.json
MIN_COVERAGE := 80

# ==============================
# Default task
# ==============================
all: tidy build lint test coverage security

# ==============================
# Dependency / mod tidy check
# ==============================
tidy:
	@echo ">> Running go mod tidy and checking differences..."
	go mod tidy
	@git diff --exit-code go.mod go.sum || (echo "go.mod/go.sum is not tidy" && exit 1)

# ==============================
# Build
# ==============================
build:
	@echo ">> Building project..."
	go build -v ./...

# ==============================
# Tests
# ==============================
test:
	@echo ">> Running tests with coverage..."
	go test ./... -coverprofile=$(COVERAGE_FILE) -covermode=atomic

# ==============================
# Coverage check
# ==============================
coverage: test
	@echo ">> Checking coverage threshold ($(MIN_COVERAGE)%)..."
	@TOTAL=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $${TOTAL%.*} -lt $(MIN_COVERAGE) ]; then \
		echo "Coverage below threshold: $${TOTAL}%"; exit 1; \
	else \
		echo "Coverage OK: $${TOTAL}%"; \
	fi

# ==============================
# Lint
# ==============================
lint:
	@echo ">> Running golangci-lint..."
	@golangci-lint run ./... --out-format=colored-line-number > $(LINT_REPORT) || (cat $(LINT_REPORT); exit 1)

# ==============================
# Security checks
# ==============================
security:
	@echo ">> Running gosec security checks..."
	@gosec -fmt=json -out $(SECURITY_REPORT) ./...
	@echo "Security report saved to $(SECURITY_REPORT)"

# ==============================
# Clean
# ==============================
clean:
	@echo ">> Cleaning build artifacts..."
	@rm -f $(COVERAGE_FILE) $(LINT_REPORT) $(SECURITY_REPORT)
	@echo "Done."
