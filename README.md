# CC Switcher

English [![简体中文](https://img.shields.io/badge/Language-简体中文-red.svg)](README_CN.md)

A command-line tool for quickly launching specified commands with specific environment variables, especially suitable for rapidly starting Claude Code with different model configurations.

## Features

- 🚀 **Zero Dependencies**: Single executable file, no runtime dependencies required
- 🔄 **Multi-Environment Support**: Quickly switch between multiple environment configurations
- ⚙️ **YAML Configuration**: User-friendly YAML format for environment configuration
- 📁 **Automatic Configuration Management**: Configuration files automatically stored in user directory
- 🌐 **Cross-Platform**: Supports Windows, macOS, and Linux

## Quick Start

1. Download the executable file for your platform:
   - Windows: `cs.exe`
   - Linux/macOS: `cs-*`

2. Place the executable file in your system PATH

3. Configuration files are automatically created on first run:

   ```bash
   cs glm
   ```

## Configuration File Format

Configuration files are automatically created in the user directory:

- Windows: `%USERPROFILE%\.cs\config.yaml`
- Linux/macOS: `~/.cs/config.yaml`

### Default GLM Configuration (Claude Code)

```yaml
environments:
  # GLM environment configuration for Claude Code
  glm:
    target: "claude"  # Claude Code command
    environment:
      CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"
      ANTHROPIC_BASE_URL: "https://open.bigmodel.cn/api/anthropic"
      ANTHROPIC_AUTH_TOKEN: "your-glm-api-key"  # Replace with actual API key
      ANTHROPIC_MODEL: "glm-4.6"
      ANTHROPIC_SMALL_FAST_MODEL: "glm-4.5-air"
      ANTHROPIC_DEFAULT_SONNET_MODEL: "glm-4.6"
      ANTHROPIC_DEFAULT_OPUS_MODEL: "glm-4.6"
      ANTHROPIC_DEFAULT_HAIKU_MODEL: "glm-4.5-air"
      API_TIMEOUT_MS: "3000000"
```

### Adding More Environment Configurations

```yaml
  # Example: Node.js development environment
  node-dev:
    target: "node server.js"
    environment:
      PORT: "3000"
      NODE_ENV: "development"
      DEBUG: "true"

  # Example: Python virtual environment
  python-env:
    target: "python app.py"
    environment:
      PYTHONPATH: "/path/to/project"
      DJANGO_SETTINGS_MODULE: "myproject.settings"
```

## Usage

```bash
# Launch command using glm environment
cs glm

# List available environments
cs
```

## Building from Source

For users who want to build from source:

### Prerequisites

- Go 1.24.1 or later

### Quick Build

```bash
# Build for current platform
go build -o cs main.go

# For Windows
go build -o cs.exe main.go
```

### Cross-Platform Build

```bash
# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o cs-windows.exe main.go

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o cs-linux main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o cs-macos-intel main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o cs-macos-arm64 main.go
```

## Releases

Pre-compiled binaries are available in the [GitHub Releases](https://github.com/yourusername/cc-switcher/releases) section for:

- Windows 64-bit (`cs-windows.exe`)
- Linux 64-bit (`cs-linux`)
- Linux ARM64 (`cs-linux-arm64`)
- macOS Intel (`cs-macos-intel`)
- macOS Apple Silicon (`cs-macos-arm64`)

## Contributing

For development and contribution guidelines, please see [MAINTAINER.md](MAINTAINER.md).

## How It Works

1. `cs <environment>` reads the configuration for the specified environment
2. Injects configured environment variables into the current environment
3. Launches the configured target command
4. Inherits standard input/output from the current terminal

## License

MIT License
