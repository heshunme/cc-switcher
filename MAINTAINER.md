# Maintainer Guide

This guide contains information for maintainers and developers working on the CC Switcher project.

## Development Setup

### Prerequisites
- Go 1.24.1 or later
- Make (for build automation)

### Initial Setup
```bash
# Install development dependencies (linters, security tools)
make dev-setup
```

## Build System

### Makefile Commands
```bash
# Development commands
make dev-setup          # Install development dependencies
make build              # Build for current platform
make build-all          # Build for all platforms
make clean              # Remove build artifacts
make test               # Run tests
make help               # Show all available commands

# Code quality
make fmt                # Format code
make lint               # Run golangci-lint
make security           # Run security scan with gosec

# Platform-specific builds
make build-windows      # Windows 64-bit
make build-linux        # Linux 64-bit
make build-linux-arm64  # Linux ARM64
make build-macos-intel  # macOS Intel
make build-macos-arm64  # macOS Apple Silicon

# Release preparation
make release            # Full build for all platforms
```

### Build Configuration
- Uses static linking (`CGO_ENABLED=0`) for cross-platform compatibility
- Version and build time embedded via ldflags
- Optimized binaries with `-s -w` flags to remove debug information

### Cross-Compilation Parameters
- `GOOS`: Target operating system (windows, linux, darwin)
- `GOARCH`: Target architecture (amd64, arm64, 386)
- `CGO_ENABLED=0`: Disable CGO for static linking

## CI/CD Pipeline

### GitHub Actions Workflows

#### 1. CI Workflow (`.github/workflows/ci.yml`)
- **Triggers**: Push/PR to main and develop branches
- **Steps**:
  - Code testing with race detection and coverage
  - Linting with golangci-lint
  - Multi-platform build validation
  - Security scanning with gosec
  - Upload coverage to Codecov

#### 2. Release Workflow (`.github/workflows/release.yml`)
- **Triggers**:
  - Push of `v*` tags
  - Manual workflow dispatch
- **Steps**:
  - Cross-platform builds for all supported architectures
  - Artifact upload
  - Release creation with changelog
  - Binary testing
- **Outputs**: GitHub Release with platform-specific binaries

#### 3. Auto-Tag Workflow (`.github/workflows/auto-tag.yml`)
- **Triggers**: Push to main branch
- **Features**:
  - Automatic version number generation
  - Tag creation and push
  - Release workflow triggering

### Release Process

#### Automated Release (Tag-based)
```bash
# Create and push version tag
git tag v1.0.0
git push origin v1.0.0
```

#### Manual Release (GitHub Actions)
1. Go to GitHub Actions page
2. Select "Release" workflow
3. Click "Run workflow"
4. Enter version number
5. Trigger release

### Release Artifacts
Each release automatically generates:
- `cs-windows.exe` - Windows 64-bit
- `cs-linux` - Linux 64-bit
- `cs-linux-arm64` - Linux ARM64
- `cs-macos-intel` - macOS Intel
- `cs-macos-arm64` - macOS Apple Silicon

### Version Management
- Tags follow semantic versioning (`v1.0.0`, `v1.0.1`, etc.)
- Pre-release versions use hyphens (`v1.0.0-beta`, `v1.0.0-rc1`)
- Auto-tag workflow generates versions based on commit history

## Testing Strategy

### Current State
- No unit test files currently exist
- CI includes functional testing of compiled binaries
- Manual testing described in README

### Recommended Test Additions
1. **Unit Tests**: For configuration parsing and environment management
2. **Integration Tests**: For command execution with different environment setups
3. **Cross-Platform Tests**: Validate behavior on all supported platforms
4. **Configuration Validation**: Test YAML parsing and error handling

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Run specific test
go test -v ./path/to/package
```

## Code Quality

### Linting
```bash
# Run linter
make lint

# Run with specific options
golangci-lint run --timeout=5m
```

### Security Scanning
```bash
# Run security scan
make security

# Run with gosec directly
gosec ./...
```

### Code Formatting
```bash
# Format code
make fmt

# Or use go fmt directly
go fmt ./...
```

## Project Structure

### Key Files
```
cc-switcher/
├── main.go              # Core application logic
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── Makefile            # Build automation
├── CLAUDE.md           # Claude Code guidance
├── MAINTAINER.md       # This file
├── README.md           # User documentation (English)
├── README_CN.md        # User documentation (Chinese)
├── .github/
│   └── workflows/      # CI/CD configurations
└── cs.exe              # Pre-built binary
```

### Architecture Principles
- Single-file Go application with clear function separation
- Minimal external dependencies (only `gopkg.in/yaml.v3`)
- Cross-platform compatibility
- User-friendly configuration system

## Contributing Guidelines

### Development Workflow
1. Fork the repository
2. Create feature branch
3. Make changes with appropriate testing
4. Ensure code quality checks pass
5. Submit pull request

### Code Style
- Follow Go conventions
- Use clear function and variable names
- Add appropriate comments for complex logic
- Maintain cross-platform compatibility

### Testing Requirements
- All new features should include tests
- Ensure builds work on all supported platforms
- Validate configuration handling
- Test error conditions

## Release Checklist

### Before Release
- [ ] All tests pass
- [ ] Code quality checks pass
- [ ] Documentation updated
- [ ] Version number updated in go.mod if needed
- [ ] CHANGELOG.md updated (if exists)

### Release Process
- [ ] Create and push version tag
- [ ] Verify automated release workflow
- [ ] Test released binaries
- [ ] Update GitHub Release notes if needed

### Post-Release
- [ ] Verify downloads work correctly
- [ ] Update documentation if needed
- [ ] Announce release (if applicable)

## Common Issues and Solutions

### Build Issues
- **CGO errors**: Ensure `CGO_ENABLED=0` for cross-platform builds
- **Missing dependencies**: Run `make dev-setup` to install linters
- **Platform-specific builds**: Use correct `GOOS` and `GOARCH` values

### CI/CD Issues
- **Permission errors**: Ensure workflow has `contents: write` permission
- **Tag conflicts**: Check if tag already exists before creating
- **Build failures**: Check GitHub Actions logs for specific errors

### Development Issues
- **Linting failures**: Run `make lint` locally to fix before committing
- **Test failures**: Ensure all dependencies are available
- **Cross-platform issues**: Test on multiple platforms when possible