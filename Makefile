.PHONY: build test lint fmt check clean install release-dry

BINARY_NAME=claude-code-status-line
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BINARY_NAME) ./cmd/claude-code-status-line

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w cmd internal

check: fmt lint test

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

install: build
	mkdir -p $(HOME)/.local/bin
	cp $(BINARY_NAME) $(HOME)/.local/bin/

release-dry:
	goreleaser release --snapshot --clean
