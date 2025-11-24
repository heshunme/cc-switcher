# CC Switcher

[![English](https://img.shields.io/badge/Language-English-blue.svg)](README.md) 简体中文

一个用于快速使用特定环境变量启动指定命令的命令行工具，尤其适用于快速使用不同模型配置快速启动claude code。

## 特性

- 🚀 **零依赖**: 单个可执行文件，无需运行时依赖
- 🔄 **多环境支持**: 支持在多套环境配置间快速切换
- ⚙️ **YAML配置**: 使用人性化的YAML格式配置环境
- 📁 **自动配置管理**: 配置文件自动存储在用户目录下
- 🌐 **跨平台**: 支持 Windows、macOS 和 Linux

## 快速开始

1. 下载对应平台的可执行文件：
   - Windows: `cs.exe`
   - Linux/macOS: `cs-*`

2. 将可执行文件放到系统PATH中

3. 首次运行时自动创建配置文件：

   ```bash
   cs glm
   ```

## 配置文件格式

配置文件自动创建在用户目录下：

- Windows: `%USERPROFILE%\.cs\config.yaml`
- Linux/macOS: `~/.cs/config.yaml`

### 默认GLM配置（Claude Code）

```yaml
environments:
  # GLM environment configuration for Claude Code
  glm:
    target: "claude"  # Claude Code command
    environment:
      CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"
      ANTHROPIC_BASE_URL: "https://open.bigmodel.cn/api/anthropic"
      ANTHROPIC_AUTH_TOKEN: "your-glm-api-key"  # 替换为实际的API密钥
      ANTHROPIC_MODEL: "glm-4.6"
      ANTHROPIC_SMALL_FAST_MODEL: "glm-4.5-air"
      ANTHROPIC_DEFAULT_SONNET_MODEL: "glm-4.6"
      ANTHROPIC_DEFAULT_OPUS_MODEL: "glm-4.6"
      ANTHROPIC_DEFAULT_HAIKU_MODEL: "glm-4.5-air"
      API_TIMEOUT_MS: "3000000"
```

### 添加更多环境配置

```yaml
  # 示例：Node.js 开发环境
  node-dev:
    target: "node server.js"
    environment:
      PORT: "3000"
      NODE_ENV: "development"
      DEBUG: "true"

  # 示例：Python 虚拟环境
  python-env:
    target: "python app.py"
    environment:
      PYTHONPATH: "/path/to/project"
      DJANGO_SETTINGS_MODULE: "myproject.settings"
```

## 使用方法

```bash
# 使用 glm 环境启动命令
cs glm

# 查看可用环境
cs

```

## 从源码编译

如果需要从源码构建：

### 环境要求

- Go 1.24.1 或更高版本

### 快速编译

```bash
# 构建当前平台版本
go build -o cs main.go

# Windows 版本
go build -o cs.exe main.go
```

### 跨平台编译

```bash
# Windows 64位
GOOS=windows GOARCH=amd64 go build -o cs-windows.exe main.go

# Linux 64位
GOOS=linux GOARCH=amd64 go build -o cs-linux main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o cs-macos-intel main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o cs-macos-arm64 main.go
```

## 发布版本

预编译的二进制文件可在 [GitHub Releases](https://github.com/yourusername/cc-switcher/releases) 页面下载：

- Windows 64位 (`cs-windows.exe`)
- Linux 64位 (`cs-linux`)
- Linux ARM64 (`cs-linux-arm64`)
- macOS Intel (`cs-macos-intel`)
- macOS Apple Silicon (`cs-macos-arm64`)

## 贡献指南

开发和贡献相关指南请参阅 [MAINTAINER.md](MAINTAINER.md)。

## 工作原理

1. `cs <environment>` 读取指定环境的配置
2. 注入配置的环境变量到当前环境
3. 启动配置的目标命令
4. 继承当前终端的标准输入/输出

## 许可证

MIT License
