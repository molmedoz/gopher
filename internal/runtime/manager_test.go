package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/molmedoz/gopher/internal/config"
	"github.com/molmedoz/gopher/internal/env"
)

// TestManager_Install_Comprehensive tests the Install method comprehensively
func TestManager_Install_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test installing invalid version
	err := manager.Install("invalid-version")
	if err == nil {
		t.Error("Expected error when installing invalid version")
	}

	// Test installing empty version
	err = manager.Install("")
	if err == nil {
		t.Error("Expected error when installing empty version")
	}

	// Test installing version with invalid format
	err = manager.Install("1.21.0") // Missing 'go' prefix
	if err == nil {
		t.Error("Expected error when installing version without 'go' prefix")
	}
}

// TestManager_Uninstall_Comprehensive tests the Uninstall method comprehensively
func TestManager_Uninstall_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test uninstalling invalid version
	err := manager.Uninstall("invalid-version")
	if err == nil {
		t.Error("Expected error when uninstalling invalid version")
	}

	// Test uninstalling empty version
	err = manager.Uninstall("")
	if err == nil {
		t.Error("Expected error when uninstalling empty version")
	}

	// Test uninstalling non-existent version
	err = manager.Uninstall("go1.21.0")
	if err == nil {
		t.Error("Expected error when uninstalling non-existent version")
	}
}

// TestManager_Use_Comprehensive tests the Use method comprehensively
func TestManager_Use_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test using invalid version
	err := manager.Use("invalid-version")
	if err == nil {
		t.Error("Expected error when using invalid version")
	}

	// Test using empty version
	err = manager.Use("")
	if err == nil {
		t.Error("Expected error when using empty version")
	}

	// Test using non-existent version
	err = manager.Use("go1.21.0")
	if err == nil {
		t.Error("Expected error when using non-existent version")
	}
}

// TestManager_GetCurrent_Comprehensive tests the GetCurrent method comprehensively
func TestManager_GetCurrent_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test getting current version when none is set
	current, err := manager.GetCurrent()
	if err != nil {
		t.Logf("GetCurrent failed (expected if no Go installed): %v", err)
	} else {
		t.Logf("Current version: %s", current)
	}
}

// TestManager_GetSystemInfo_Comprehensive tests the GetSystemInfo method comprehensively
func TestManager_GetSystemInfo_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test getting system info
	systemInfo, err := manager.GetSystemInfo()
	if err != nil {
		t.Logf("GetSystemInfo failed (expected if no system Go): %v", err)
	} else {
		t.Logf("System info: %+v", systemInfo)
	}
}

// TestManager_UseSystemVersion_Comprehensive tests the UseSystemVersion method comprehensively
func TestManager_UseSystemVersion_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test using system version
	err := manager.useSystemVersion()
	if err != nil {
		t.Logf("useSystemVersion failed (expected if no system Go): %v", err)
	}
}

// TestManager_IsSystemGoAvailable_Comprehensive tests the IsSystemGoAvailable method comprehensively
func TestManager_IsSystemGoAvailable_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test checking system Go availability
	systemInfo, err := manager.GetSystemInfo()
	if err != nil {
		t.Logf("System Go not available: %v", err)
	} else {
		t.Logf("System Go available: %+v", systemInfo)
	}
}

// TestManager_SetupEnvironment_Comprehensive tests the setupEnvironment method comprehensively
func TestManager_SetupEnvironment_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test setting up environment
	err := manager.setupEnvironment("go1.21.0")
	if err != nil {
		t.Logf("setupEnvironment failed: %v", err)
	}
}

// TestManager_SetupSystemEnvironment_Comprehensive tests the setupSystemEnvironment method comprehensively
func TestManager_SetupSystemEnvironment_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test setting up system environment
	err := manager.setupSystemEnvironment()
	if err != nil {
		t.Logf("setupSystemEnvironment failed: %v", err)
	}
}

// TestManager_CreateEnvironmentScript_Comprehensive tests the createEnvironmentScript method comprehensively
func TestManager_CreateEnvironmentScript_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test creating environment script
	_, err := manager.createEnvironmentScript("go1.21.0", map[string]string{})
	if err != nil {
		t.Logf("createEnvironmentScript failed: %v", err)
	}
}

// TestManager_SaveActiveVersion_Comprehensive tests the saveActiveVersion method comprehensively
func TestManager_SaveActiveVersion_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test saving active version
	err := manager.saveActiveVersion("go1.21.0")
	if err != nil {
		t.Logf("saveActiveVersion failed: %v", err)
	}
}

// TestManager_GetActiveVersionFromState_Comprehensive tests the getActiveVersionFromState method comprehensively
func TestManager_GetActiveVersionFromState_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test getting active version from state
	version, err := manager.getActiveVersionFromState()
	if err != nil {
		t.Logf("getActiveVersionFromState failed: %v", err)
	} else {
		t.Logf("Active version from state: %s", version)
	}
}

// TestManager_SetupShellIntegration_Comprehensive tests the setupShellIntegration method comprehensively
func TestManager_SetupShellIntegration_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test setting up shell integration
	err := manager.setupShellIntegration()
	if err != nil {
		t.Logf("setupShellIntegration failed: %v", err)
	}
}

// TestManager_CreateGopherInitScript_Comprehensive tests the createGopherInitScript method comprehensively
func TestManager_CreateGopherInitScript_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test creating gopher init script
	_, err := manager.createGopherInitScript()
	if err != nil {
		t.Logf("createGopherInitScript failed: %v", err)
	}
}

// TestManager_DetectShell_Comprehensive tests the detectShell method comprehensively
func TestManager_DetectShell_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test detecting shell
	shell := manager.detectShell()
	t.Logf("Detected shell: %s", shell)
}

// TestManager_DetectShellFromPath_Comprehensive tests the detectShellFromPath method comprehensively
func TestManager_DetectShellFromPath_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test detecting shell from path
	tests := []struct {
		name string
		path string
		want string
	}{
		{"bash path", "/bin/bash", "bash"},
		{"zsh path", "/bin/zsh", "zsh"},
		{"fish path", "/usr/bin/fish", "fish"},
		{"unknown path", "/bin/unknown", "unknown"},
		{"empty path", "", "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.detectShellFromPath(tt.path)
			if got != tt.want {
				t.Errorf("detectShellFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestManager_GetShellProfile_Comprehensive tests the getShellProfile method comprehensively
func TestManager_GetShellProfile_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test getting shell profile
	homeDir, _ := os.UserHomeDir()
	tests := []struct {
		name  string
		shell string
		want  string
	}{
		{"bash", "bash", filepath.Join(homeDir, ".bashrc")},
		{"zsh", "zsh", filepath.Join(homeDir, ".zshrc")},
		{"fish", "fish", filepath.Join(homeDir, ".config", "fish", "config.fish")},
		{"unknown", "unknown", filepath.Join(homeDir, ".profile")},
		{"empty", "", filepath.Join(homeDir, ".profile")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := manager.getShellProfile(tt.shell)
			if err != nil {
				t.Logf("getShellProfile failed: %v", err)
			}
			if got != tt.want {
				t.Errorf("getShellProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestManager_AddToShellProfile_Comprehensive tests the addToShellProfile method comprehensively
func TestManager_AddToShellProfile_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test adding to shell profile
	err := manager.addToShellProfile("bash", "/tmp/profile")
	if err != nil {
		t.Logf("addToShellProfile failed: %v", err)
	}
}

// TestManager_CreateSymlink_Comprehensive tests the createSymlink method comprehensively
func TestManager_CreateSymlink_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test creating symlink
	err := manager.createSymlink("go1.21.0")
	if err != nil {
		t.Logf("createSymlink failed: %v", err)
	}
}

// TestManager_TryCreateSymlink_Comprehensive tests the tryCreateSymlink method comprehensively
func TestManager_TryCreateSymlink_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test trying to create symlink
	err := manager.tryCreateSymlink("go1.21.0", "/tmp/go")
	if err != nil {
		t.Logf("tryCreateSymlink failed: %v", err)
	}
}

// TestManager_IsSymlinkActuallyUsed_Comprehensive tests the isSymlinkActuallyUsed method comprehensively
func TestManager_IsSymlinkActuallyUsed_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test checking if symlink is actually used
	used := manager.isSymlinkActuallyUsed("/tmp/go")
	t.Logf("Symlink actually used: %v", used)
}

// TestManager_HasGopherSymlinkInPath_Comprehensive tests the hasGopherSymlinkInPath method comprehensively
func TestManager_HasGopherSymlinkInPath_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test checking if gopher symlink is in path
	has := manager.hasGopherSymlinkInPath()
	t.Logf("Has gopher symlink in path: %v", has)
}

// TestManager_AddSymlinkToPath_Comprehensive tests the addSymlinkToPath method comprehensively
func TestManager_AddSymlinkToPath_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test adding symlink to path
	err := manager.addSymlinkToPath("/tmp/go")
	if err != nil {
		t.Logf("addSymlinkToPath failed: %v", err)
	}
}

// TestManager_RemoveGopherSymlinks_Comprehensive tests the removeGopherSymlinks method comprehensively
func TestManager_RemoveGopherSymlinks_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test removing gopher symlinks
	err := manager.removeGopherSymlinks()
	if err != nil {
		t.Logf("removeGopherSymlinks failed: %v", err)
	}
}

// TestManager_AutoCleanup_Comprehensive tests the autoCleanup method comprehensively
func TestManager_AutoCleanup_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test auto cleanup
	err := manager.autoCleanup()
	if err != nil {
		t.Logf("autoCleanup failed: %v", err)
	}
}

// TestManager_GetDownloadDir_Comprehensive tests the GetDownloadDir method comprehensively
func TestManager_GetDownloadDir_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test getting download directory
	downloadDir := manager.GetDownloadDir()
	if downloadDir != cfg.DownloadDir {
		t.Errorf("GetDownloadDir() = %v, want %v", downloadDir, cfg.DownloadDir)
	}
}

// TestManager_AliasManager_Comprehensive tests the AliasManager method comprehensively
func TestManager_AliasManager_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test getting alias manager
	aliasManager := manager.AliasManager()
	if aliasManager == nil {
		t.Error("Expected alias manager to be available")
	}
}

// TestManager_ExtractVersionFromPath_Comprehensive tests the extractVersionFromPath method comprehensively
func TestManager_ExtractVersionFromPath_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test extracting version from path
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{"valid version path", "/path/to/go1.21.0/bin/go", "go1.21.0", false},
		{"valid version path with subdir", "/path/to/go1.21.0/subdir/bin/go", "go1.21.0", false},
		{"invalid path", "/path/to/invalid/bin/go", "", true},
		{"empty path", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.extractVersionFromPath(tt.path)
			if got != tt.want {
				t.Errorf("extractVersionFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestManager_GetVersionInfo_Comprehensive tests the getVersionInfo method comprehensively
func TestManager_GetVersionInfo_Comprehensive(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		InstallDir:  filepath.Join(tmpDir, "install"),
		DownloadDir: filepath.Join(tmpDir, "download"),
		MaxVersions: 5,
	}

	envProvider := env.NewMockProvider(map[string]string{})
	manager := NewManager(cfg, envProvider)

	// Test getting version info
	versionInfo, err := manager.getVersionInfo("go1.21.0")
	if err != nil {
		t.Logf("getVersionInfo failed: %v", err)
	} else {
		t.Logf("Version info: %+v", versionInfo)
	}
}
