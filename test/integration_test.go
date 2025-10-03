package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/molmedoz/gopher/internal/config"
	"github.com/molmedoz/gopher/internal/env"
	inruntime "github.com/molmedoz/gopher/internal/runtime"
)

// TestIntegration_EndToEndWorkflow tests the complete workflow from installation to usage
func TestIntegration_EndToEndWorkflow(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	installDir := filepath.Join(tmpDir, "install")
	downloadDir := filepath.Join(tmpDir, "download")

	// Create test configuration
	cfg := &config.Config{
		InstallDir:     installDir,
		DownloadDir:    downloadDir,
		MaxVersions:    5,
		SetEnvironment: true,
		AutoCleanup:    true,
	}

	// Create manager
	envProvider := env.NewMockProvider(map[string]string{})
	manager := inruntime.NewManager(cfg, envProvider)

	// Test 1: List installed versions (should be empty initially)
	installed, err := manager.ListInstalled()
	if err != nil {
		t.Fatalf("Failed to list installed versions: %v", err)
	}
	if len(installed) < 1 {
		t.Errorf("Expected at least 1 installed version (system), got %d", len(installed))
	}

	// Test 2: Check if any version is installed
	isInstalled, err := manager.IsInstalled("go1.21.0")
	if err != nil {
		t.Logf("IsInstalled failed: %v", err)
	}
	if isInstalled {
		t.Error("Expected go1.21.0 to not be installed initially")
	}

	// Test 3: Get current version (should be system or unknown)
	_, err = manager.GetCurrent()
	if err != nil {
		t.Logf("GetCurrent failed (expected if no system Go): %v", err)
	}

	// Test 4: Test system version detection
	systemInfo, err := manager.GetSystemInfo()
	if err != nil {
		t.Logf("GetSystemInfo failed (expected if no system Go): %v", err)
	} else {
		t.Logf("System Go info: %+v", systemInfo)
	}

	// Test 5: Test alias manager
	aliasManager := manager.AliasManager()
	if aliasManager == nil {
		t.Fatal("Expected alias manager to be available")
	}

	// Test 6: Test alias operations (without actual installation)
	err = aliasManager.CreateAlias("test", "go1.21.0")
	if err == nil {
		t.Error("Expected error when creating alias for non-installed version")
	}

	// Test 7: Test configuration methods
	installDirResult := manager.GetInstallDir()
	if installDirResult != installDir {
		t.Errorf("Expected install dir %s, got %s", installDir, installDirResult)
	}

	config := manager.GetConfig()
	if config == nil {
		t.Fatal("Expected config to be available")
	}
	if config.InstallDir != installDir {
		t.Errorf("Expected config install dir %s, got %s", installDir, config.InstallDir)
	}
}

// TestIntegration_AliasWorkflow tests the complete alias workflow
func TestIntegration_AliasWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmpDir, "install"),
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := inruntime.NewManager(cfg, envProvider)
	aliasManager := manager.AliasManager()

	// Test alias validation
	tests := []struct {
		name    string
		alias   string
		version string
		wantErr bool
	}{
		{"valid alias", "stable", "go1.21.0", true}, // Will fail because version not installed
		{"empty alias", "", "go1.21.0", true},
		{"reserved name", "system", "go1.21.0", true},
		{"invalid characters", "test@alias", "go1.21.0", true},
		{"too long", "verylongaliasnamethatexceedslimit", "go1.21.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aliasManager.CreateAlias(tt.alias, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Test listing aliases
	aliases, err := aliasManager.ListAliases()
	if err != nil {
		t.Fatalf("Failed to list aliases: %v", err)
	}
	if len(aliases) != 0 {
		t.Errorf("Expected 0 aliases, got %d", len(aliases))
	}
}

// TestIntegration_ConfigurationWorkflow tests configuration management
func TestIntegration_ConfigurationWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	// Test 1: Create default config
	cfg := config.DefaultConfig()
	if cfg.InstallDir == "" {
		t.Error("Expected default install dir to be set")
	}

	// Test 2: Validate config
	err := cfg.Validate()
	if err != nil {
		t.Fatalf("Default config validation failed: %v", err)
	}

	// Test 3: Save and load config
	cfg.InstallDir = filepath.Join(tmpDir, "custom-install")
	cfg.DownloadDir = filepath.Join(tmpDir, "custom-download")
	cfg.MaxVersions = 10

	configPath := filepath.Join(tmpDir, "config.json")
	err = cfg.Save(configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedCfg.InstallDir != cfg.InstallDir {
		t.Errorf("Expected install dir %s, got %s", cfg.InstallDir, loadedCfg.InstallDir)
	}
	if loadedCfg.DownloadDir != cfg.DownloadDir {
		t.Errorf("Expected download dir %s, got %s", cfg.DownloadDir, loadedCfg.DownloadDir)
	}
	if loadedCfg.MaxVersions != cfg.MaxVersions {
		t.Errorf("Expected max versions %d, got %d", cfg.MaxVersions, loadedCfg.MaxVersions)
	}

	// Test 4: Test environment variables
	envProvider := env.NewMockProvider(map[string]string{})
	envVars := cfg.GetEnvironmentVariablesWithEnv("go1.21.0", envProvider)
	if len(envVars) == 0 {
		t.Error("Expected environment variables to be set")
	}

	// Test 5: Test GOPATH generation
	gopath := cfg.GetGOPATH("go1.21.0")
	if gopath == "" {
		t.Error("Expected GOPATH to be set")
	}

	goroot := cfg.GetGOROOT("go1.21.0")
	if goroot == "" {
		t.Error("Expected GOROOT to be set")
	}
}

// TestIntegration_ErrorHandling tests error handling across the system
func TestIntegration_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmpDir, "install"),
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := inruntime.NewManager(cfg, envProvider)

	// Test 1: Install non-existent version
	err := manager.Install("invalid-version")
	if err == nil {
		t.Error("Expected error when installing invalid version")
	}

	// Test 2: Uninstall non-installed version
	err = manager.Uninstall("go1.21.0")
	if err == nil {
		t.Error("Expected error when uninstalling non-installed version")
	}

	// Test 3: Use non-installed version
	err = manager.Use("go1.21.0")
	if err == nil {
		t.Error("Expected error when using non-installed version")
	}

	// Test 4: Test alias operations with invalid data
	aliasManager := manager.AliasManager()

	// Test creating alias for non-installed version
	err = aliasManager.CreateAlias("test", "go1.21.0")
	if err == nil {
		t.Error("Expected error when creating alias for non-installed version")
	}

	// Test removing non-existent alias
	err = aliasManager.RemoveAlias("nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existent alias")
	}

	// Test updating non-existent alias
	err = aliasManager.UpdateAlias("nonexistent", "go1.21.0")
	if err == nil {
		t.Error("Expected error when updating non-existent alias")
	}
}

// TestIntegration_SystemDetection tests system Go detection
func TestIntegration_SystemDetection(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmpDir, "install"),
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := inruntime.NewManager(cfg, envProvider)

	// Test system Go detection
	systemInfo, err := manager.GetSystemInfo()
	if err != nil {
		t.Logf("System Go not available (expected in test environment): %v", err)
		return
	}

	if systemInfo == nil {
		t.Error("Expected system info to be available")
		return
	}

	// Test system Go availability
	systemInfo, err = manager.GetSystemInfo()
	if err != nil {
		t.Logf("System Go not available: %v", err)
	} else {
		t.Logf("System Go available: %+v", systemInfo)
	}

	// Test using system version
	if systemInfo != nil {
		err = manager.Use("system")
		if err != nil {
			t.Logf("Failed to use system version: %v", err)
		}
	}
}

// TestIntegration_FileSystemOperations tests file system operations
func TestIntegration_FileSystemOperations(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
	}

	envProvider := env.NewMockProvider(map[string]string{})
	_ = inruntime.NewManager(cfg, envProvider)

	// Test directory creation
	err := cfg.EnsureDirectories()
	if err != nil {
		t.Fatalf("Failed to ensure directories: %v", err)
	}

	// Check if directories were created
	if _, err := os.Stat(cfg.InstallDir); os.IsNotExist(err) {
		t.Errorf("Install directory was not created: %s", cfg.InstallDir)
	}

	if _, err := os.Stat(cfg.DownloadDir); os.IsNotExist(err) {
		t.Errorf("Download directory was not created: %s", cfg.DownloadDir)
	}
}

// TestIntegration_ConcurrentOperations tests concurrent operations
func TestIntegration_ConcurrentOperations(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmpDir, "install"),
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := inruntime.NewManager(cfg, envProvider)
	aliasManager := manager.AliasManager()

	// Test concurrent alias operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Test alias validation
			aliasManager.ValidateAliasName("test")

			// Test listing aliases
			aliasManager.ListAliases()

			// Test checking if version is installed
			manager.IsInstalled("go1.21.0")

			// Test listing installed versions
			manager.ListInstalled()
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
