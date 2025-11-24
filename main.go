package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Environment represents a single environment configuration
type Environment struct {
	Target      string            `yaml:"target"`
	Environment map[string]string `yaml:"environment"`
}

// Config represents the entire configuration file
type Config struct {
	Environments map[string]Environment `yaml:"environments"`
}

// GetConfigDir returns the user's cs config directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".cs")
	return configDir, nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

// LoadConfig loads the configuration from file or creates a default one
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		if err := createDefaultConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		fmt.Printf("Created default configuration file: %s\n", configPath)
		fmt.Println("Please edit the file to add your environment configurations.")
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.Environments == nil {
		config.Environments = make(map[string]Environment)
	}

	return &config, nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configPath string) error {
	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	defaultConfig := `# CS Switcher Configuration File
# Define your environment configurations here

environments:
  # GLM environment configuration for Claude Code
  glm:
    target: "claude"  # Claude Code command
    environment:
      CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"
      ANTHROPIC_BASE_URL: "https://open.bigmodel.cn/api/anthropic"
      ANTHROPIC_AUTH_TOKEN: "your-glm-api-key"
      ANTHROPIC_MODEL: "glm-4.6"
      ANTHROPIC_SMALL_FAST_MODEL: "glm-4.5-air"
      ANTHROPIC_DEFAULT_SONNET_MODEL: "glm-4.6"
      ANTHROPIC_DEFAULT_OPUS_MODEL: "glm-4.6"
      ANTHROPIC_DEFAULT_HAIKU_MODEL: "glm-4.5-air"
      API_TIMEOUT_MS: "3000000"

# Add more environments as needed
# Example:
#   myenv:
#     target: "node server.js"
#     environment:
#       PORT: "3000"
#       NODE_ENV: "production"
`

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}

// Run executes the target command with the specified environment variables
func Run(envConfig Environment) error {
	if envConfig.Target == "" {
		return fmt.Errorf("target command is empty")
	}

	// Split the target command into command and arguments
	parts := strings.Fields(envConfig.Target)
	if len(parts) == 0 {
		return fmt.Errorf("invalid target command")
	}

	command := parts[0]
	args := parts[1:]

	// Create the command
	cmd := exec.Command(command, args...)

	// Set up environment variables
	if len(envConfig.Environment) > 0 {
		// Start with current environment
		env := os.Environ()

		// Add or override with our environment variables
		for key, value := range envConfig.Environment {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}

		cmd.Env = env
	}

	// Set up standard I/O to connect to the current terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cs <environment>")
		fmt.Println("Available environments:")
		printAvailableEnvironments()
		os.Exit(1)
	}

	environment := os.Args[1]

	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Check if environment exists
	envConfig, exists := cfg.Environments[environment]
	if !exists {
		fmt.Printf("Environment '%s' not found.\n", environment)
		fmt.Println("Available environments:")
		printAvailableEnvironments()
		os.Exit(1)
	}

	// Run the command with environment variables
	if err := Run(envConfig); err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}
}

func printAvailableEnvironments() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("  (Unable to load config: %v)\n", err)
		return
	}

	for name := range cfg.Environments {
		fmt.Printf("  %s\n", name)
	}
}
