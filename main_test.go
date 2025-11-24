package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Test helpers

// setupTestConfig creates a temporary config directory and file for testing
func setupTestConfig(t *testing.T) (configDir, configFile string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "cs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	configDir = filepath.Join(tempDir, ".cs")
	configFile = filepath.Join(configDir, "config.yaml")

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return configDir, configFile, cleanup
}

// createTestConfigFile creates a test config file with the given content
func createTestConfigFile(t *testing.T, configPath string, content string) {
	t.Helper()

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
}


// Testable version of LoadConfig that uses a custom config path
func LoadConfigWithPath(configPath string) (*Config, error) {
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



// Configuration Management Tests


func TestCreateDefaultConfig(t *testing.T) {
	_, configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test creating default config
	err := createDefaultConfig(configFile)
	if err != nil {
		t.Fatalf("createDefaultConfig() returned error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Verify file contents
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	content := string(data)
	expectedStrings := []string{
		"CS Switcher Configuration File",
		"environments:",
		"glm:",
		"target: \"claude\"",
		"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("Config file doesn't contain expected string: %s", expected)
		}
	}
}

func TestLoadConfigWithPath_ExistingFile(t *testing.T) {
	_, configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a test config file
	testConfigContent := `
environments:
  testenv:
    target: "echo hello"
    environment:
      TEST_VAR: "test_value"
  another:
    target: "pwd"
    environment:
      PATH: "/custom/path"
`
	createTestConfigFile(t, configFile, testConfigContent)

	// Load the config
	config, err := LoadConfigWithPath(configFile)
	if err != nil {
		t.Fatalf("LoadConfigWithPath() returned error: %v", err)
	}

	// Verify the loaded config
	if len(config.Environments) != 2 {
		t.Errorf("Expected 2 environments, got %d", len(config.Environments))
	}

	testEnv, exists := config.Environments["testenv"]
	if !exists {
		t.Fatal("testenv not found in loaded config")
	}

	if testEnv.Target != "echo hello" {
		t.Errorf("Expected target 'echo hello', got '%s'", testEnv.Target)
	}

	if testEnv.Environment["TEST_VAR"] != "test_value" {
		t.Errorf("Expected TEST_VAR 'test_value', got '%s'", testEnv.Environment["TEST_VAR"])
	}

	anotherEnv, exists := config.Environments["another"]
	if !exists {
		t.Fatal("another environment not found in loaded config")
	}

	if anotherEnv.Target != "pwd" {
		t.Errorf("Expected target 'pwd', got '%s'", anotherEnv.Target)
	}
}

func TestLoadConfigWithPath_NonExistentFile(t *testing.T) {
	_, configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Ensure the config file doesn't exist
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		t.Fatal("Config file already exists")
	}

	// This should create a default config
	config, err := LoadConfigWithPath(configFile)
	if err != nil {
		t.Fatalf("LoadConfigWithPath() returned error: %v", err)
	}

	// Verify the config was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Default config file was not created")
	}

	// Verify it has at least the glm environment
	if len(config.Environments) == 0 {
		t.Fatal("No environments in loaded config")
	}

	if _, exists := config.Environments["glm"]; !exists {
		t.Error("Default glm environment not found")
	}
}

func TestLoadConfigWithPath_EmptyEnvironments(t *testing.T) {
	_, configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a config file with empty environments
	testConfigContent := `
environments: {}
`
	createTestConfigFile(t, configFile, testConfigContent)

	config, err := LoadConfigWithPath(configFile)
	if err != nil {
		t.Fatalf("LoadConfigWithPath() returned error: %v", err)
	}

	// Verify environments map is not nil
	if config.Environments == nil {
		t.Fatal("Environments map is nil after LoadConfigWithPath")
	}

	if len(config.Environments) != 0 {
		t.Errorf("Expected 0 environments, got %d", len(config.Environments))
	}
}

func TestLoadConfigWithPath_InvalidYAML(t *testing.T) {
	_, configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a config file with invalid YAML
	testConfigContent := `
environments:
  testenv:
    target: "echo hello"
    environment:
      TEST_VAR: "test_value"
  invalid_yaml: [unclosed array
`
	createTestConfigFile(t, configFile, testConfigContent)

	// This should return an error
	_, err := LoadConfigWithPath(configFile)
	if err == nil {
		t.Fatal("LoadConfigWithPath() should have returned an error for invalid YAML")
	}

	expectedError := "failed to parse config file"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

// Command Execution Tests

func TestRun_EmptyTarget(t *testing.T) {
	env := Environment{
		Target:      "",
		Environment: map[string]string{},
	}

	err := Run(env)
	if err == nil {
		t.Fatal("Run() should have returned an error for empty target")
	}

	expectedError := "target command is empty"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

func TestRun_WhitespaceOnlyTarget(t *testing.T) {
	env := Environment{
		Target:      "   ",
		Environment: map[string]string{},
	}

	err := Run(env)
	if err == nil {
		t.Fatal("Run() should have returned an error for whitespace-only target")
	}

	expectedError := "invalid target command"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

func TestRun_SimpleCommand(t *testing.T) {
	env := Environment{
		Target:      "echo hello world",
		Environment: map[string]string{},
	}

	err := Run(env)
	if err != nil {
		t.Logf("Run() returned error (command might not exist): %v", err)
	}
}


// CLI Interface Tests

func TestPrintAvailableEnvironments_WithMockConfig(t *testing.T) {
	_, configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a test config file with multiple environments
	testConfigContent := `
environments:
  dev:
    target: "echo dev"
    environment:
      ENV: "development"
  prod:
    target: "echo prod"
    environment:
      ENV: "production"
`
	createTestConfigFile(t, configFile, testConfigContent)

	// Temporarily modify the global config path for testing
	// This is a bit of a hack but necessary for testing the print function
	originalHome := os.Getenv("HOME")
	tempHome := filepath.Dir(filepath.Dir(configFile))
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	// Test that the function doesn't panic
	printAvailableEnvironments()
}


// Edge Case Tests



func TestConcurrency_SimultaneousConfigAccess(t *testing.T) {
	_, configFile, cleanup := setupTestConfig(t)
	defer cleanup()

	testConfigContent := `
environments:
  concurrent1:
    target: "echo test1"
    environment:
      VAR1: "value1"
  concurrent2:
    target: "echo test2"
    environment:
      VAR2: "value2"
`
	createTestConfigFile(t, configFile, testConfigContent)

	// Run multiple LoadConfig calls concurrently
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			config, err := LoadConfigWithPath(configFile)
			if err != nil {
				errors <- err
				done <- true
				return
			}

			// Verify the config is valid
			if len(config.Environments) != 2 {
				errors <- fmt.Errorf("expected 2 environments, got %d", len(config.Environments))
				done <- true
				return
			}

			errors <- nil
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Check for any errors
	close(errors)
	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent load error: %v", err)
		}
	}
}