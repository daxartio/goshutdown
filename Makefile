DEFAULT_GOAL := all

.PHONY: all
all: fmt check test

.PHONY: deps
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: fmt
fmt:
	go fmt ./...
	golangci-lint run --fix --new ./...

.PHONY: check
check:
	go vet ./...
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: tidy
tidy:
	go mod tidy
