.PHONY: all tidy build test lint coverage security

GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")

all: tidy build lint test coverage security

tidy:
	@echo ">> Running go mod tidy"
	go mod tidy
	git diff --exit-code go.mod go.sum

build:
	@echo ">> Building project"
	go build -v ./...

test:
	@echo ">> Running tests with coverage"
	go test ./... -coverprofile=coverage.out -covermode=atomic

coverage: test
	@echo ">> Checking coverage threshold (min 80%)"
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//' | awk '{if ($$1+0 < 80) {print "Coverage below threshold: "$$1"%"; exit 1}}'

lint:
	@echo ">> Running golangci-lint"
	golangci-lint run ./...

security:
	@echo ">> Running gosec security checks"
	gosec ./...
