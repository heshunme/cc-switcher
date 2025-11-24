# Makefile for cc-switcher

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=cs
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build parameters
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Targets
.PHONY: all build clean test deps help run build-all release

# Default target
all: clean deps test build

# Build for current platform
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) -v

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux -v

# Build for Linux ARM64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 -v

# Build for macOS Intel
build-macos-intel:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-macos-intel -v

# Build for macOS ARM64
build-macos-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-macos-arm64 -v

# Build for all platforms
build-all: build-windows build-linux build-linux-arm64 build-macos-intel build-macos-arm64

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)*
	rm -f coverage.out

# Run tests
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run linter
lint:
	golangci-lint run


# Run the binary
run: build
	./$(BINARY_NAME)

# Install locally
install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/

# Create release artifacts
release: clean deps test build-all

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  build            - Build binary for current platform"
	@echo "  build-windows    - Build binary for Windows"
	@echo "  build-linux      - Build binary for Linux"
	@echo "  build-linux-arm64- Build binary for Linux ARM64"
	@echo "  build-macos-intel- Build binary for macOS Intel"
	@echo "  build-macos-arm64- Build binary for macOS ARM64"
	@echo "  build-all        - Build binaries for all platforms"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run tests"
	@echo "  deps             - Download dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  run              - Build and run binary"
	@echo "  install          - Install binary to GOPATH/bin"
	@echo "  release          - Create release artifacts"
	@echo "  dev-setup        - Setup development environment"
	@echo "  help             - Show this help message"