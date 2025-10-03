package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.InstallDir == "" {
		t.Error("InstallDir should not be empty")
	}
	if config.DownloadDir == "" {
		t.Error("DownloadDir should not be empty")
	}
	if config.MirrorURL == "" {
		t.Error("MirrorURL should not be empty")
	}
	if config.MaxVersions < 1 {
		t.Error("MaxVersions should be at least 1")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				InstallDir:  "/tmp/install",
				DownloadDir: "/tmp/download",
				MirrorURL:   "https://example.com",
				MaxVersions: 5,
			},
			wantErr: false,
		},
		{
			name: "empty install dir",
			config: &Config{
				InstallDir:  "",
				DownloadDir: "/tmp/download",
				MirrorURL:   "https://example.com",
				MaxVersions: 5,
			},
			wantErr: true,
		},
		{
			name: "empty download dir",
			config: &Config{
				InstallDir:  "/tmp/install",
				DownloadDir: "",
				MirrorURL:   "https://example.com",
				MaxVersions: 5,
			},
			wantErr: true,
		},
		{
			name: "empty mirror URL",
			config: &Config{
				InstallDir:  "/tmp/install",
				DownloadDir: "/tmp/download",
				MirrorURL:   "",
				MaxVersions: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid max versions",
			config: &Config{
				InstallDir:  "/tmp/install",
				DownloadDir: "/tmp/download",
				MirrorURL:   "https://example.com",
				MaxVersions: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a test config
	originalConfig := &Config{
		InstallDir:  "/tmp/install",
		DownloadDir: "/tmp/download",
		MirrorURL:   "https://example.com",
		AutoCleanup: true,
		MaxVersions: 5,
	}

	// Save the config
	if err := originalConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the config
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Compare the configs
	if loadedConfig.InstallDir != originalConfig.InstallDir {
		t.Errorf("InstallDir = %v, want %v", loadedConfig.InstallDir, originalConfig.InstallDir)
	}
	if loadedConfig.DownloadDir != originalConfig.DownloadDir {
		t.Errorf("DownloadDir = %v, want %v", loadedConfig.DownloadDir, originalConfig.DownloadDir)
	}
	if loadedConfig.MirrorURL != originalConfig.MirrorURL {
		t.Errorf("MirrorURL = %v, want %v", loadedConfig.MirrorURL, originalConfig.MirrorURL)
	}
	if loadedConfig.AutoCleanup != originalConfig.AutoCleanup {
		t.Errorf("AutoCleanup = %v, want %v", loadedConfig.AutoCleanup, originalConfig.AutoCleanup)
	}
	if loadedConfig.MaxVersions != originalConfig.MaxVersions {
		t.Errorf("MaxVersions = %v, want %v", loadedConfig.MaxVersions, originalConfig.MaxVersions)
	}
}

func TestConfigEnsureDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	config := &Config{
		InstallDir:  filepath.Join(tempDir, "install"),
		DownloadDir: filepath.Join(tempDir, "download"),
		MirrorURL:   "https://example.com",
		MaxVersions: 5,
	}

	// Ensure directories are created
	if err := config.EnsureDirectories(); err != nil {
		t.Fatalf("Failed to ensure directories: %v", err)
	}

	// Check if directories exist
	if _, err := os.Stat(config.InstallDir); os.IsNotExist(err) {
		t.Errorf("InstallDir %s was not created", config.InstallDir)
	}
	if _, err := os.Stat(config.DownloadDir); os.IsNotExist(err) {
		t.Errorf("DownloadDir %s was not created", config.DownloadDir)
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()

	if path == "" {
		t.Error("GetConfigPath() should not return empty string")
	}

	// Check if it contains expected components
	if !filepath.IsAbs(path) {
		t.Error("GetConfigPath() should return absolute path")
	}
}
