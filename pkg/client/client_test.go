package client

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := &Config{
		ServerURL: "https://test-server.com",
		ProjectID: "test-project",
		Platform:  "linux",
		Arch:      "amd64",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	// 检查默认值
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", config.Timeout)
	}

	if config.BackupCount != 3 {
		t.Errorf("Expected default backup count to be 3, got %d", config.BackupCount)
	}

	if len(config.PreserveFiles) == 0 {
		t.Error("Expected default preserve files to be set")
	}
}

func TestNewClientInvalidConfig(t *testing.T) {
	testCases := []struct {
		name   string
		config *Config
	}{
		{
			name:   "nil config",
			config: nil,
		},
		{
			name: "missing ServerURL",
			config: &Config{
				ProjectID: "test-project",
				Platform:  "linux",
				Arch:      "amd64",
			},
		},
		{
			name: "missing ProjectID",
			config: &Config{
				ServerURL: "https://test-server.com",
				Platform:  "linux",
				Arch:      "amd64",
			},
		},
		{
			name: "invalid platform",
			config: &Config{
				ServerURL: "https://test-server.com",
				ProjectID: "test-project",
				Platform:  "invalid",
				Arch:      "amd64",
			},
		},
		{
			name: "invalid arch",
			config: &Config{
				ServerURL: "https://test-server.com",
				ProjectID: "test-project",
				Platform:  "linux",
				Arch:      "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewClient(tc.config)
			if err == nil {
				t.Error("Expected error for invalid config")
			}
		})
	}
}

func TestCheckForUpdates(t *testing.T) {
	// 这里需要mock HTTP客户端，暂时跳过实际的网络请求测试
	t.Skip("Need to implement HTTP client mocking")
}

func TestValidateConfig(t *testing.T) {
	validConfig := &Config{
		ServerURL: "https://test-server.com", 
		ProjectID: "test-project",
		Platform:  "linux",
		Arch:      "amd64",
	}

	err := validateConfig(validConfig)
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got %v", err)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}
	
	if !contains(slice, "apple") {
		t.Error("Expected to find 'apple' in slice")
	}
	
	if contains(slice, "grape") {
		t.Error("Expected not to find 'grape' in slice")
	}
}