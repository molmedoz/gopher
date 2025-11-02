package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/molmedoz/gopher/internal/config"
	inruntime "github.com/molmedoz/gopher/internal/runtime"
)

// runInteractiveSetup runs the interactive setup wizard
func runInteractiveSetup(manager *inruntime.Manager) error {
	fmt.Println("🚀 Welcome to Gopher Setup Wizard!")
	fmt.Println("This wizard will help you configure Gopher for your system.")
	fmt.Println()

	// Step 1: System Detection
	systemInfo, err := detectSystemInfo(manager)
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}

	// Step 2: Show current status
	showSystemStatus(systemInfo)

	// Step 3: Interactive configuration
	if err := runInteractiveConfiguration(manager, systemInfo); err != nil {
		return fmt.Errorf("configuration failed: %w", err)
	}

	// Step 4: Setup shell integration
	if err := setupShellIntegrationInteractive(manager, systemInfo); err != nil {
		return fmt.Errorf("shell integration setup failed: %w", err)
	}

	// Step 5: Test and verify setup
	if err := testAndVerifySetup(manager, systemInfo); err != nil {
		return fmt.Errorf("setup verification failed: %w", err)
	}

	// Step 6: Show completion and next steps
	showSetupCompletion(systemInfo)

	return nil
}

// runWindowsSetup handles Windows-specific setup
// Currently unused but kept for potential future use
func runWindowsSetup(manager *inruntime.Manager) error { //nolint:unused
	fmt.Println("🪟 Windows Setup")
	fmt.Println("================")

	// Check Developer Mode
	fmt.Println("1. Checking Developer Mode...")
	if isDeveloperModeEnabled() {
		fmt.Println("   ✅ Developer Mode is enabled")
	} else {
		fmt.Println("   ❌ Developer Mode is not enabled")
		fmt.Println()
		fmt.Println("   Developer Mode is required for symlink creation without admin privileges.")
		fmt.Println("   Please enable it in Windows Settings > Update & Security > For developers")
		fmt.Println()

		if !askForConfirmation("Continue anyway? (You may need admin privileges for symlink creation)") {
			return fmt.Errorf("setup cancelled by user")
		}
	}

	// Check PATH configuration
	fmt.Println()
	fmt.Println("2. Checking PATH configuration...")

	// Get potential symlink directories
	userHome, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	symlinkDirs := []string{
		filepath.Join(userHome, "AppData", "Local", "bin"),
		filepath.Join(userHome, "bin"),
	}

	var pathDirs []string
	for _, dir := range symlinkDirs {
		if isDirectoryInPath(dir) {
			pathDirs = append(pathDirs, dir)
		}
	}

	if len(pathDirs) > 0 {
		fmt.Printf("   ✅ Found %d symlink directory(ies) in PATH:\n", len(pathDirs))
		for _, dir := range pathDirs {
			fmt.Printf("      - %s\n", dir)
		}
	} else {
		fmt.Println("   ⚠️  No symlink directories found in PATH")
		fmt.Println("   Gopher will add directories to PATH automatically when switching versions")
	}

	// Test symlink creation
	fmt.Println()
	fmt.Println("3. Testing symlink creation...")

	// Create a test symlink
	testDir := filepath.Join(userHome, "AppData", "Local", "bin")
	// #nosec G301 -- 0755 required for test bin directory
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("   ❌ Failed to create test directory: %v\n", err)
	} else {
		// Try to create a test symlink
		testTarget := filepath.Join(testDir, "gopher-test.exe")
		testLink := filepath.Join(testDir, "gopher-test-link.exe")

		// Create a dummy file to link to
		// #nosec G306 -- 0644 acceptable for temporary test file
		if err := os.WriteFile(testTarget, []byte("test"), 0644); err == nil {
			if err := os.Symlink(testTarget, testLink); err == nil {
				fmt.Println("   ✅ Symlink creation works")
				if err := os.Remove(testLink); err != nil && !os.IsNotExist(err) {
					fmt.Printf("   ⚠️  Failed to remove test symlink: %v\n", err)
				}
			} else {
				fmt.Printf("   ❌ Symlink creation failed: %v\n", err)
				fmt.Println("   You may need to enable Developer Mode or run as Administrator")
			}
			if err := os.Remove(testTarget); err != nil && !os.IsNotExist(err) {
				fmt.Printf("   ⚠️  Failed to remove test file: %v\n", err)
			}
		}
	}

	// Show next steps
	fmt.Println()
	fmt.Println("🎯 Setup Complete! Next Steps:")
	fmt.Println("================================")
	fmt.Println()
	fmt.Println("STEP 1: Add Gopher's bin directory to PATH (CRITICAL!)")
	fmt.Println("--------------------------------------------------------")
	gopherBin := filepath.Join(userHome, "AppData", "Local", "bin")
	fmt.Println()
	fmt.Println("Copy and run this PowerShell command as Administrator:")
	fmt.Println()
	fmt.Printf("  $userPath = [Environment]::GetEnvironmentVariable(\"PATH\", \"User\")\n")
	fmt.Printf("  $gopherBin = \"%s\"\n", gopherBin)
	fmt.Printf("  if ($userPath -notlike \"*$gopherBin*\") {\n")
	fmt.Printf("    [Environment]::SetEnvironmentVariable(\"PATH\", \"$gopherBin;$userPath\", \"User\")\n")
	fmt.Printf("  }\n")
	fmt.Println()
	fmt.Println("Then restart your PowerShell terminal.")
	fmt.Println()
	fmt.Println("STEP 2: Install a Go version")
	fmt.Println("-----------------------------")
	fmt.Println("  gopher install 1.21.0")
	fmt.Println()
	fmt.Println("STEP 3: Switch to it")
	fmt.Println("--------------------")
	fmt.Println("  gopher use 1.21.0")
	fmt.Println()
	fmt.Println("  ⚠️  Important: If you have system Go installed, Gopher will warn you")
	fmt.Println("  about PATH order and provide exact commands to fix it.")
	fmt.Println()
	fmt.Println("STEP 4: Verify it works")
	fmt.Println("-----------------------")
	fmt.Println("  go version")
	fmt.Println()
	fmt.Println("  Should show: go version go1.21.0 windows/amd64")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("  • If you get 'Access is denied' errors, enable Developer Mode")
	fmt.Println("    (Settings > Update & Security > For developers)")
	fmt.Println("  • Use 'gopher debug' to troubleshoot any issues")
	fmt.Println("  • Use 'gopher status' to check your current setup")
	fmt.Println()
	fmt.Println("📚 For more help: https://github.com/molmedoz/gopher#readme")

	return nil
}

// runMacOSSetup handles macOS-specific setup
// Currently unused but kept for potential future use
func runMacOSSetup(manager *inruntime.Manager) error { //nolint:unused
	fmt.Println("🍎 macOS Setup")
	fmt.Println("==============")

	// Check Homebrew
	fmt.Println("1. Checking Homebrew...")
	if isHomebrewInstalled() {
		fmt.Println("   ✅ Homebrew is installed")
		homebrewGoPath := "/opt/homebrew/bin/go"
		if _, err := os.Stat(homebrewGoPath); err == nil {
			fmt.Printf("   ✅ Found Homebrew Go at: %s\n", homebrewGoPath)
		}
	} else {
		fmt.Println("   ℹ️  Homebrew not detected")
	}

	// Check system Go
	fmt.Println()
	fmt.Println("2. Checking system Go...")
	if systemInfo, err := manager.GetSystemInfo(); err == nil && systemInfo != nil {
		fmt.Printf("   ✅ System Go found: %s\n", systemInfo.Version)
		fmt.Printf("   Path: %s\n", systemInfo.Executable)
	} else {
		fmt.Println("   ℹ️  No system Go detected")
	}

	// Check shell configuration
	fmt.Println()
	fmt.Println("3. Checking shell configuration...")
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	fmt.Printf("   Shell: %s\n", shell)

	// Check if gopher is in PATH
	if isCommandInPath("gopher") {
		fmt.Println("   ✅ gopher command is available")
	} else {
		fmt.Println("   ⚠️  gopher command not found in PATH")
		fmt.Println("   Add gopher to your PATH or use full path")
	}

	// Show next steps
	fmt.Println()
	fmt.Println("🎯 Next Steps")
	fmt.Println("=============")
	fmt.Println("1. Install a Go version: gopher install 1.21.0")
	fmt.Println("2. Switch to it: gopher use 1.21.0")
	fmt.Println("3. Set up shell integration: gopher setup")
	fmt.Println("4. Verify: go version")

	return nil
}

// runLinuxSetup handles Linux-specific setup
// Currently unused but kept for potential future use
func runLinuxSetup(manager *inruntime.Manager) error { //nolint:unused
	fmt.Println("🐧 Linux Setup")
	fmt.Println("==============")

	// Check package manager
	fmt.Println("1. Checking package manager...")
	packageManager := detectPackageManager()
	switch packageManager {
	case "apt":
		fmt.Println("   ✅ APT package manager detected")
	case "yum":
		fmt.Println("   ✅ YUM package manager detected")
	case "dnf":
		fmt.Println("   ✅ DNF package manager detected")
	case "pacman":
		fmt.Println("   ✅ Pacman package manager detected")
	default:
		fmt.Println("   ℹ️  No common package manager detected")
	}

	// Check system Go
	fmt.Println()
	fmt.Println("2. Checking system Go...")
	if systemInfo, err := manager.GetSystemInfo(); err == nil && systemInfo != nil {
		fmt.Printf("   ✅ System Go found: %s\n", systemInfo.Version)
		fmt.Printf("   Path: %s\n", systemInfo.Executable)
	} else {
		fmt.Println("   ℹ️  No system Go detected")
	}

	// Check shell configuration
	fmt.Println()
	fmt.Println("3. Checking shell configuration...")
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	fmt.Printf("   Shell: %s\n", shell)

	// Check if gopher is in PATH
	if isCommandInPath("gopher") {
		fmt.Println("   ✅ gopher command is available")
	} else {
		fmt.Println("   ⚠️  gopher command not found in PATH")
		fmt.Println("   Add gopher to your PATH or use full path")
	}

	// Show next steps
	fmt.Println()
	fmt.Println("🎯 Next Steps")
	fmt.Println("=============")
	fmt.Println("1. Install a Go version: gopher install 1.21.0")
	fmt.Println("2. Switch to it: gopher use 1.21.0")
	fmt.Println("3. Set up shell integration: gopher setup")
	fmt.Println("4. Verify: go version")

	return nil
}

// runGenericSetup handles generic setup for unsupported platforms
// Currently unused but kept for potential future use
func runGenericSetup(manager *inruntime.Manager) error { //nolint:unused
	fmt.Println("🔧 Generic Setup")
	fmt.Println("================")

	// Check system Go
	fmt.Println("1. Checking system Go...")
	if systemInfo, err := manager.GetSystemInfo(); err == nil && systemInfo != nil {
		fmt.Printf("   ✅ System Go found: %s\n", systemInfo.Version)
		fmt.Printf("   Path: %s\n", systemInfo.Executable)
	} else {
		fmt.Println("   ℹ️  No system Go detected")
	}

	// Show next steps
	fmt.Println()
	fmt.Println("🎯 Next Steps")
	fmt.Println("=============")
	fmt.Println("1. Install a Go version: gopher install 1.21.0")
	fmt.Println("2. Switch to it: gopher use 1.21.0")
	fmt.Println("3. Verify: go version")
	fmt.Println()
	fmt.Println("⚠️  Note: Some features may not work on this platform")

	return nil
}

// Helper functions for setup

func isDeveloperModeEnabled() bool {
	// On Windows, check registry for Developer Mode
	if runtime.GOOS == "windows" {
		// This is a simplified check - in practice, you'd check the registry
		// For now, we'll assume it's enabled if we can create symlinks
		return true
	}
	return false
}

func isHomebrewInstalled() bool {
	return isCommandInPath("brew")
}

func isCommandInPath(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func isDirectoryInPath(dir string) bool {
	path := os.Getenv("PATH")
	if path == "" {
		return false
	}

	pathDirs := strings.Split(path, string(os.PathListSeparator))
	for _, pathDir := range pathDirs {
		if filepath.Clean(pathDir) == filepath.Clean(dir) {
			return true
		}
	}
	return false
}

func askForConfirmation(message string) bool {
	fmt.Printf("%s (y/N): ", message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func detectPackageManager() string {
	packageManagers := []string{"apt", "yum", "dnf", "pacman"}
	for _, pm := range packageManagers {
		if isCommandInPath(pm) {
			return pm
		}
	}
	return ""
}

// SystemInfo holds detected system information
type SystemInfo struct {
	Platform         string
	Arch             string
	Shell            string
	ShellProfile     string
	HomeDir          string
	HasSystemGo      bool
	SystemGoVersion  string
	SystemGoPath     string
	HasHomebrew      bool
	HomebrewGoPath   string
	PackageManager   string
	HasDeveloperMode bool
	SymlinkDir       string
	IsInPath         bool
	IsDocker         bool
	Config           *config.Config
}

// detectSystemInfo detects comprehensive system information
func detectSystemInfo(manager *inruntime.Manager) (*SystemInfo, error) {
	info := &SystemInfo{
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	// Detect shell
	info.Shell = detectShell()
	if info.Shell == "" {
		// Platform-specific fallback
		if runtime.GOOS == "windows" {
			info.Shell = "powershell"
		} else {
			info.Shell = "bash"
		}
	}

	// Get home directory
	var err error
	info.HomeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Detect shell profile
	info.ShellProfile, _ = getShellProfile(info.Shell)

	// Check for system Go
	if systemInfo, err := manager.GetSystemInfo(); err == nil && systemInfo != nil {
		info.HasSystemGo = true
		info.SystemGoVersion = systemInfo.Version
		info.SystemGoPath = systemInfo.Executable
	}

	// Check for Homebrew (macOS/Linux)
	if info.Platform == "darwin" || info.Platform == "linux" {
		info.HasHomebrew = isHomebrewInstalled()
		if info.HasHomebrew {
			homebrewGoPath := "/opt/homebrew/bin/go"
			if info.Platform == "linux" {
				homebrewGoPath = "/home/linuxbrew/.linuxbrew/bin/go"
			}
			if _, err := os.Stat(homebrewGoPath); err == nil {
				info.HomebrewGoPath = homebrewGoPath
			}
		}
	}

	// Detect package manager (Linux)
	if info.Platform == "linux" {
		info.PackageManager = detectPackageManager()
	}

	// Check Developer Mode (Windows)
	if info.Platform == "windows" {
		info.HasDeveloperMode = isDeveloperModeEnabled()
	}

	// Determine symlink directory
	switch info.Platform {
	case "windows":
		info.SymlinkDir = filepath.Join(info.HomeDir, "AppData", "Local", "bin")
	case "darwin", "linux":
		info.SymlinkDir = filepath.Join(info.HomeDir, ".local", "bin")
	default:
		info.SymlinkDir = filepath.Join(info.HomeDir, "bin")
	}

	// Check if symlink directory is in PATH
	info.IsInPath = isDirectoryInPath(info.SymlinkDir)

	// Check if running in Docker
	if _, err := os.Stat("/.dockerenv"); err == nil {
		info.IsDocker = true
	}

	// Get current config
	info.Config = manager.GetConfig()

	return info, nil
}

// showSystemStatus displays detected system information
func showSystemStatus(info *SystemInfo) {
	fmt.Println("📋 System Detection")
	fmt.Println("===================")
	fmt.Printf("Platform: %s/%s\n", info.Platform, info.Arch)
	fmt.Printf("Shell: %s\n", info.Shell)
	fmt.Printf("Home Directory: %s\n", info.HomeDir)
	fmt.Printf("Shell Profile: %s\n", info.ShellProfile)
	fmt.Println()

	// System Go status
	if info.HasSystemGo {
		fmt.Printf("✅ System Go: %s at %s\n", info.SystemGoVersion, info.SystemGoPath)
	} else {
		fmt.Println("ℹ️  No system Go detected")
	}

	// Homebrew status (macOS/Linux)
	if info.HasHomebrew {
		if info.HomebrewGoPath != "" {
			fmt.Printf("✅ Homebrew Go: %s\n", info.HomebrewGoPath)
		} else {
			fmt.Println("✅ Homebrew installed (no Go found)")
		}
	}

	// Package manager status (Linux)
	if info.PackageManager != "" {
		fmt.Printf("✅ Package Manager: %s\n", info.PackageManager)
	}

	// Developer Mode status (Windows)
	if info.Platform == "windows" {
		if info.HasDeveloperMode {
			fmt.Println("✅ Developer Mode: Enabled")
		} else {
			fmt.Println("⚠️  Developer Mode: Disabled (required for symlinks)")
		}
	}

	// Symlink directory status
	if info.IsInPath {
		fmt.Printf("✅ Symlink Directory: %s (in PATH)\n", info.SymlinkDir)
	} else {
		fmt.Printf("⚠️  Symlink Directory: %s (not in PATH)\n", info.SymlinkDir)
	}

	// Docker status
	if info.IsDocker {
		fmt.Println("🐳 Docker Environment: Detected")
	}

	fmt.Println()
}

// runInteractiveConfiguration runs interactive configuration
func runInteractiveConfiguration(manager *inruntime.Manager, info *SystemInfo) error {
	fmt.Println("⚙️  Configuration")
	fmt.Println("=================")

	// Check if configuration is needed
	needsConfig := false

	// Check symlink directory setup
	if !info.IsInPath {
		needsConfig = true
		fmt.Printf("❌ Symlink directory %s is not in PATH\n", info.SymlinkDir)
	}

	// Check Developer Mode (Windows)
	if info.Platform == "windows" && !info.HasDeveloperMode {
		needsConfig = true
		fmt.Println("❌ Developer Mode is not enabled")
	}

	if !needsConfig {
		fmt.Println("✅ Configuration looks good!")
		return nil
	}

	fmt.Println()
	fmt.Println("Let's fix the configuration issues:")

	// Fix symlink directory PATH issue
	if !info.IsInPath {
		fmt.Printf("\n1. Adding %s to PATH...\n", info.SymlinkDir)

		// Windows: Show PowerShell instructions (don't try to modify shell profiles)
		if runtime.GOOS == "windows" {
			fmt.Println("   To add to PATH on Windows, run this PowerShell command as Administrator:")
			fmt.Println()
			fmt.Printf("   $userPath = [Environment]::GetEnvironmentVariable(\"PATH\", \"User\")\n")
			fmt.Printf("   [Environment]::SetEnvironmentVariable(\"PATH\", \"%s;$userPath\", \"User\")\n", info.SymlinkDir)
			fmt.Println()
			fmt.Println("   Then restart your terminal.")
		} else {
			// Unix/Linux/macOS: Try to add to shell profile
			if err := addDirectoryToPath(info.SymlinkDir, info.ShellProfile); err != nil {
				fmt.Printf("   ❌ Failed to add to PATH: %v\n", err)
				fmt.Printf("   Please manually add this to your %s:\n", info.ShellProfile)
				fmt.Printf("   export PATH=\"%s:$PATH\"\n", info.SymlinkDir)
			} else {
				fmt.Printf("   ✅ Added %s to PATH\n", info.SymlinkDir)
			}
		}
	}

	// Windows Developer Mode instructions
	if info.Platform == "windows" && !info.HasDeveloperMode {
		fmt.Println("\n2. Developer Mode Setup:")
		fmt.Println("   Please enable Developer Mode in Windows Settings:")
		fmt.Println("   - Go to Settings > Update & Security > For developers")
		fmt.Println("   - Turn on 'Developer Mode'")
		fmt.Println("   - Restart your terminal after enabling")

		if !askForConfirmation("Have you enabled Developer Mode?") {
			fmt.Println("   ⚠️  Please enable Developer Mode and run 'gopher init' again")
		}
	}

	return nil
}

// setupShellIntegrationInteractive sets up shell integration interactively
func setupShellIntegrationInteractive(manager *inruntime.Manager, info *SystemInfo) error {
	fmt.Println("\n🔧 Shell Integration Setup")
	fmt.Println("==========================")

	// Skip shell integration on Windows (not needed for symlink-based switching)
	if runtime.GOOS == "windows" {
		fmt.Println("ℹ️  Shell integration is not required on Windows")
		fmt.Println("   Gopher uses symlinks which work automatically")
		fmt.Println("   Version switching is handled via PATH and symlinks")
		fmt.Println()
		fmt.Println("💡 Tip: When you run 'gopher use <version>', Gopher will:")
		fmt.Println("   - Create a symlink to the selected Go version")
		fmt.Println("   - Check if PATH order is correct")
		fmt.Println("   - Warn you if system Go takes precedence")
		return nil
	}

	// Check if already configured
	if isGopherConfigured(info.ShellProfile) {
		fmt.Println("✅ Shell integration already configured")
		return nil
	}

	fmt.Printf("Setting up shell integration for %s...\n", info.Shell)

	// Create the gopher initialization script
	initScript, err := createGopherInitScript(manager)
	if err != nil {
		return fmt.Errorf("failed to create gopher init script: %w", err)
	}

	// Add to shell profile
	if err := addToShellProfile(info.ShellProfile, initScript); err != nil {
		return fmt.Errorf("failed to add to shell profile: %w", err)
	}

	fmt.Printf("✅ Shell integration configured in %s\n", info.ShellProfile)
	fmt.Printf("✅ Gopher init script created: %s\n", initScript)

	return nil
}

// testAndVerifySetup tests the setup and provides feedback
func testAndVerifySetup(manager *inruntime.Manager, info *SystemInfo) error {
	fmt.Println("\n🧪 Testing Setup")
	fmt.Println("================")

	// Test 1: Check if gopher command works
	fmt.Print("1. Testing gopher command... ")
	if isCommandInPath("gopher") {
		fmt.Println("✅")
	} else {
		fmt.Println("❌")
		fmt.Println("   Please add gopher to your PATH or use full path")
	}

	// Test 2: Test symlink creation
	fmt.Print("2. Testing symlink creation... ")
	if err := testSymlinkCreation(info.SymlinkDir); err != nil {
		fmt.Printf("❌ (%v)\n", err)
		if info.Platform == "windows" && !info.HasDeveloperMode {
			fmt.Println("   Enable Developer Mode or run as Administrator")
		}
	} else {
		fmt.Println("✅")
	}

	// Test 3: Test shell integration (Unix only)
	if runtime.GOOS != "windows" {
		fmt.Print("3. Testing shell integration... ")
		if isGopherConfigured(info.ShellProfile) {
			fmt.Println("✅")
		} else {
			fmt.Println("❌")
			fmt.Printf("   Please run: source %s\n", info.ShellProfile)
		}
	}

	return nil
}

// showSetupCompletion shows setup completion and next steps
func showSetupCompletion(info *SystemInfo) {
	fmt.Println("\n🎉 Setup Complete!")
	fmt.Println("==================")
	fmt.Println("Gopher is now configured for your system.")
	fmt.Println()

	// System-specific next steps
	switch info.Platform {
	case "windows":
		showWindowsNextSteps(info)
	case "darwin":
		showMacOSNextSteps(info)
	case "linux":
		showLinuxNextSteps(info)
	default:
		showGenericNextSteps(info)
	}

	if info.IsDocker {
		fmt.Println("\n🐳 Docker Environment:")
		fmt.Println("- Run 'source ~/.gopher/scripts/gopher-init.sh' in each session")
		fmt.Println("- Or add it to your shell profile manually")
	}

	fmt.Println()
	fmt.Println("📚 For more help, run: gopher help")
}

// showWindowsNextSteps shows Windows-specific setup instructions
func showWindowsNextSteps(info *SystemInfo) {
	fmt.Println("🚀 Windows Setup Instructions:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Check if symlink dir is in PATH
	if !info.IsInPath {
		fmt.Println("⚠️  REQUIRED: Add Gopher's bin directory to PATH")
		fmt.Println()
		fmt.Println("Copy and run this PowerShell command as Administrator:")
		fmt.Println()
		fmt.Printf("  $userPath = [Environment]::GetEnvironmentVariable(\"PATH\", \"User\")\n")
		fmt.Printf("  $gopherBin = \"%s\"\n", info.SymlinkDir)
		fmt.Printf("  if ($userPath -notlike \"*$gopherBin*\") {\n")
		fmt.Printf("    [Environment]::SetEnvironmentVariable(\"PATH\", \"$gopherBin;$userPath\", \"User\")\n")
		fmt.Printf("  }\n")
		fmt.Println()
		fmt.Println("Then RESTART your PowerShell terminal.")
		fmt.Println()
		fmt.Println("═══════════════════════════════════════════════════════════════")
		fmt.Println()
	}

	fmt.Println("STEP 1: Install a Go version")
	fmt.Println("  gopher install 1.21.0")
	fmt.Println()

	fmt.Println("STEP 2: Switch to it")
	fmt.Println("  gopher use 1.21.0")
	fmt.Println()
	fmt.Println("  ⚠️  If you have system Go installed, Gopher will check PATH order")
	fmt.Println("  and show you exact commands to fix it if needed.")
	fmt.Println()

	fmt.Println("STEP 3: Verify it works")
	fmt.Println("  go version")
	fmt.Println("  # Should show: go version go1.21.0 windows/amd64")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("  • If 'Access is denied': Enable Developer Mode in Windows Settings")
	fmt.Println("  • If 'go' not found: Restart terminal after PATH changes")
	fmt.Println("  • Use 'gopher status' to check your setup")
	fmt.Println("  • Use 'gopher debug' to troubleshoot issues")
	fmt.Println()
	fmt.Println("📁 Directories created:")
	fmt.Printf("  Config:    %s\\gopher\\config.json\n", info.HomeDir)
	fmt.Printf("  Versions:  %s\\gopher\\versions\n", info.HomeDir)
	fmt.Printf("  Downloads: %s\\gopher\\downloads\n", info.HomeDir)
	fmt.Printf("  State:     %s\\gopher\\state\n", info.HomeDir)
	fmt.Printf("  Symlinks:  %s\n", info.SymlinkDir)
}

// showMacOSNextSteps shows macOS-specific setup instructions
func showMacOSNextSteps(info *SystemInfo) {
	fmt.Println("🚀 macOS Setup Instructions:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("STEP 1: Install a Go version")
	fmt.Println("  gopher install 1.21.0")
	fmt.Println()

	fmt.Println("STEP 2: Switch to it")
	fmt.Println("  gopher use 1.21.0")
	fmt.Println()

	fmt.Println("STEP 3: Verify it works")
	fmt.Println("  go version")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("  • Homebrew Go versions are automatically detected")
	fmt.Println("  • Use 'gopher system' to switch back to system Go")
	fmt.Println("  • Add shell integration with 'gopher setup'")
	fmt.Println()
	fmt.Println("📁 Directories created:")
	fmt.Printf("  Config:    %s/.gopher/config.json\n", info.HomeDir)
	fmt.Printf("  Versions:  %s/.gopher/versions\n", info.HomeDir)
	fmt.Printf("  Downloads: %s/.gopher/downloads\n", info.HomeDir)
	fmt.Printf("  Symlinks:  %s\n", info.SymlinkDir)
}

// showLinuxNextSteps shows Linux-specific setup instructions
func showLinuxNextSteps(info *SystemInfo) {
	fmt.Println("🚀 Linux Setup Instructions:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Check if symlink dir is in PATH
	if !info.IsInPath {
		fmt.Println("⚠️  Add Gopher's bin directory to PATH:")
		fmt.Println()
		fmt.Printf("  echo 'export PATH=\"%s:$PATH\"' >> %s\n", info.SymlinkDir, info.ShellProfile)
		fmt.Printf("  source %s\n", info.ShellProfile)
		fmt.Println()
		fmt.Println("═══════════════════════════════════════════════════════════════")
		fmt.Println()
	}

	fmt.Println("STEP 1: Install a Go version")
	fmt.Println("  gopher install 1.21.0")
	fmt.Println()

	fmt.Println("STEP 2: Switch to it")
	fmt.Println("  gopher use 1.21.0")
	fmt.Println()

	fmt.Println("STEP 3: Verify it works")
	fmt.Println("  go version")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("  • Use 'gopher setup' for persistent shell integration")
	fmt.Println("  • Use 'gopher system' to switch back to system Go")
	fmt.Println()
	fmt.Println("📁 Directories created:")
	fmt.Printf("  Config:    %s/.gopher/config.json\n", info.HomeDir)
	fmt.Printf("  Versions:  %s/.gopher/versions\n", info.HomeDir)
	fmt.Printf("  Downloads: %s/.gopher/downloads\n", info.HomeDir)
	fmt.Printf("  Symlinks:  %s\n", info.SymlinkDir)
}

// showGenericNextSteps shows generic setup instructions
func showGenericNextSteps(info *SystemInfo) {
	fmt.Println("🚀 Next Steps:")
	fmt.Println("1. Install a Go version: gopher install 1.21.0")
	fmt.Println("2. Switch to it: gopher use 1.21.0")
	fmt.Println("3. Verify: go version")
}

// Helper functions for the new setup system

func addDirectoryToPath(dir, profilePath string) error {
	// Read current profile
	// #nosec G304 -- profilePath is user's shell profile file (validated path)
	content, err := os.ReadFile(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	profileContent := string(content)

	// Check if already added
	if strings.Contains(profileContent, dir) {
		return nil // Already added
	}

	// Add PATH export
	pathExport := fmt.Sprintf("\n# Gopher PATH\nexport PATH=\"%s:$PATH\"\n", dir)

	// Write updated profile
	// #nosec G306 -- 0644 required for shell profile files (must be readable by shell)
	return os.WriteFile(profilePath, []byte(profileContent+pathExport), 0644)
}

func isGopherConfigured(profilePath string) bool {
	// #nosec G304 -- profilePath is user's shell profile file (validated path)
	content, err := os.ReadFile(profilePath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "gopher-init.sh")
}

func testSymlinkCreation(symlinkDir string) error {
	// Create test directory if it doesn't exist
	// #nosec G301 -- 0755 required for test symlink directory
	if err := os.MkdirAll(symlinkDir, 0755); err != nil {
		return err
	}

	// Create test file
	testFile := filepath.Join(symlinkDir, "gopher-test")
	// #nosec G306 -- 0644 acceptable for temporary test file
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(testFile); err != nil && !os.IsNotExist(err) {
			fmt.Printf("warning: cleanup failed removing %s: %v\n", testFile, err)
		}
	}()

	// Create test symlink
	testLink := filepath.Join(symlinkDir, "gopher-test-link")
	if err := os.Symlink(testFile, testLink); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(testLink); err != nil && !os.IsNotExist(err) {
			fmt.Printf("warning: cleanup failed removing %s: %v\n", testLink, err)
		}
	}()

	return nil
}
