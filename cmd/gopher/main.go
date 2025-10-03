// Package main provides the Gopher CLI application for managing Go versions.
//
// Gopher is a simple, fast, and dependency-free Go version manager that enables
// you to install, manage, and switch between multiple Go versions, including
// your system-installed Go.
//
// Features:
//   - Fast: Minimal dependencies, built with Go standard library only
//   - Secure: Cryptographic verification of downloaded Go binaries
//   - Simple: Clean CLI interface with intuitive commands
//   - Cross-platform: Works on Linux, macOS, and Windows
//   - Lightweight: Zero external dependencies beyond Go standard library
//   - Auto-cleanup: Automatically manages old versions to save space
//   - System Integration: Seamlessly manages both system and Gopher-installed Go versions
//   - JSON Support: Full JSON output for scripting and automation
//   - Environment Management: Comprehensive GOPATH and GOROOT management
//   - Structured Logging: Configurable logging with multiple levels
//
// Usage:
//
//	gopher <command> [arguments]
//
// Commands:
//
//	list                    List installed Go versions (including system)
//	list-remote             List available Go versions (with pagination and filtering)
//	install <version>       Install a Go version
//	uninstall <version>     Uninstall a Go version
//	use <version>           Switch to a Go version (use 'system' for system Go)
//	current                 Show current Go version
//	system                  Show system Go information
//	alias                   Manage version aliases (create, list, remove, show)
//	init                    Interactive setup wizard for platform-specific configuration
//	setup                   Set up shell integration for persistent Go version switching
//	status                  Show persistence status and shell integration info
//	debug                   Show debug information for troubleshooting
//	version                 Show gopher version
//	help                    Show detailed help information
//
// Options:
//
//	--json                  Output in JSON format
//	--config <path>         Path to configuration file
//	--help                  Show this help message
//	--verbose, -v           Show detailed output (DEBUG level)
//	--quiet, -q             Only show errors (ERROR level)
//
// Examples:
//
//	# Interactive setup wizard
//	gopher init
//
//	# List all installed versions
//	gopher list
//
//	# Install and switch to Go 1.21.0
//	gopher install 1.21.0
//	gopher use 1.21.0
//
//	# Switch to system Go
//	gopher use system
//
//	# Show current Go version
//	gopher current
//
//	# JSON output for scripting
//	gopher --json list
//
// For more information, visit: https://github.com/molmedoz/gopher
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/config"
	"github.com/molmedoz/gopher/internal/downloader"
	"github.com/molmedoz/gopher/internal/env"
	"github.com/molmedoz/gopher/internal/errors"
	inprogress "github.com/molmedoz/gopher/internal/progress"
	inruntime "github.com/molmedoz/gopher/internal/runtime"
)

// Version information - set via ldflags at build time
// Example: go build -ldflags "-X main.appVersion=v1.0.0 -X main.appCommit=abc123"
var (
	appVersion = "dev"     // Version (from git tag), set via: -X main.appVersion=v1.0.0
	appCommit  = "none"    // Git commit hash, set via: -X main.appCommit=abc123
	appDate    = "unknown" // Build date, set via: -X main.appDate=2025-10-13T10:30:00Z
	appBuiltBy = "source"  // Built by (goreleaser, manual, etc.), set via: -X main.appBuiltBy=goreleaser
)

// getVersionString returns the formatted version string
func getVersionString() string {
	if appVersion == "dev" {
		return "gopher dev (built from source)"
	}
	return fmt.Sprintf("gopher %s", appVersion)
}

// getFullVersionInfo returns detailed version information
func getFullVersionInfo() string {
	return fmt.Sprintf(`gopher %s
  commit: %s
  built: %s
  by: %s
  go: %s
  platform: %s/%s`,
		appVersion, appCommit, appDate, appBuiltBy,
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

const (
	usageString = `gopher - Go version manager

USAGE:
    gopher <command> [arguments]

COMMANDS:
    list                    List installed Go versions (including system)
    list-remote             List available Go versions (with pagination and filtering)
    install <version>       Install a Go version
    uninstall <version>     Uninstall a Go version
    use <version>           Switch to a Go version (use 'system' for system Go)
    current                 Show current Go version
    system                  Show system Go information
    alias                   Manage version aliases (create, list, remove, show)
    init                    Interactive setup wizard for platform-specific configuration
    setup                   Set up shell integration for persistent Go version switching
    status                  Show persistence status and shell integration info
    debug                   Show debug information for troubleshooting
    version                 Show gopher version
    help                    Show detailed help information

EXAMPLES:
    gopher list
    gopher install 1.21.0
    gopher use 1.21.0
    gopher use system
    gopher system
    gopher uninstall 1.20.7
    gopher alias create stable 1.21.0
    gopher alias list
    gopher use stable
    
    # Pagination and filtering (flags must come before command)
    gopher --no-interactive list
    gopher --page-size 5 list-remote
    gopher --page 2 --page-size 10 list-remote
    gopher --filter "1.21" list-remote
    gopher --stable list-remote
    gopher --no-interactive list-remote
    gopher --filter "rc" list-remote
    
    # Verbosity control
    gopher --verbose install 1.21.0
    gopher --quiet list
    gopher -v install 1.21.0
    gopher -q list
    
    # JSON output for scripting
    gopher --json list
    gopher --json current

For more information, visit: https://github.com/molmedoz/gopher
`
)

var (
	jsonOutput = flag.Bool("json", false, "Output in JSON format")
	configPath = flag.String("config", "", "Path to config file")
	helpFlag   = flag.Bool("help", false, "Show help information")

	// Pagination flags
	pageSize      = flag.Int("page-size", 10, "Number of versions to show per page")
	page          = flag.Int("page", 1, "Page number to display")
	filter        = flag.String("filter", "", "Filter versions by text (e.g., '1.21', 'stable', 'rc')")
	stable        = flag.Bool("stable", false, "Show only stable versions")
	noInteractive = flag.Bool("no-interactive", false, "Disable interactive pagination (default: interactive)")

	// Alias flags
	override   = flag.Bool("override", false, "Allow overriding existing aliases without confirmation")
	noOverride = flag.Bool("no-override", false, "Exit with error if alias already exists (no override allowed)")
	force      = flag.Bool("force", false, "Force operation without confirmation (overrides all other flags)")

	// Logging flags
	quiet   = flag.Bool("quiet", false, "Only show errors (sets log level to ERROR)")
	verbose = flag.Bool("verbose", false, "Show detailed output (sets log level to DEBUG)")
	q       = flag.Bool("q", false, "Short form of --quiet")
	v       = flag.Bool("v", false, "Short form of --verbose")
)

// main is the entry point of the Gopher CLI application.
//
// It handles command-line argument parsing and delegates
// command execution to the appropriate handler functions.
//
// The main function:
//  1. Parses command-line flags
//  2. Handles the help command
//  3. Loads configuration
//  4. Creates a version manager instance
//  5. Executes the requested command
//
// Supported flags:
//
//	--json, --config, --help, --verbose, --quiet, and command-specific flags
func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usageString)
	}

	// Parse flags first
	flag.Parse()

	// Check for help flag
	if *helpFlag {
		showHelp()
		return
	}

	// Get non-flag arguments
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	command := args[0]
	commandArgs := args[1:]

	// Check for global flags in command args and re-parse if needed
	hasGlobalFlags := false
	for _, arg := range commandArgs {
		if strings.HasPrefix(arg, "-") {
			hasGlobalFlags = true
			break
		}
	}

	if hasGlobalFlags {
		// Re-parse with all arguments to catch global flags after command
		os.Args = os.Args[1:] // Remove program name
		flag.Parse()
		args = flag.Args()
		if len(args) > 0 {
			command = args[0]
			commandArgs = args[1:]
		}
	}

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Create version manager with default environment provider
	manager := inruntime.NewManager(cfg, &env.DefaultProvider{})

	// Execute command
	if err := executeCommand(manager, command, commandArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig() (*config.Config, error) {
	configPath := *configPath
	if configPath == "" {
		configPath = config.GetConfigPath()
	}

	return config.Load(configPath)
}

func executeCommand(manager *inruntime.Manager, command string, args []string) error {
	switch command {
	case "list":
		return listInstalled(manager)
	case "list-remote":
		return listRemote(manager)
	case "install":
		if len(args) < 1 {
			return errors.NewMissingArgument("install (requires version)")
		}
		return installVersion(manager, args[0])
	case "uninstall":
		if len(args) < 1 {
			return errors.NewMissingArgument("uninstall (requires version)")
		}
		return uninstallVersion(manager, args[0])
	case "use":
		if len(args) < 1 {
			return errors.NewMissingArgument("use (requires version or alias)")
		}
		return useVersion(manager, args[0])
	case "current":
		return showCurrent(manager)
	case "system":
		return showSystem(manager)
	case "version":
		return showVersion()
	case "env":
		if len(args) < 1 {
			return showEnvHelp()
		}
		return handleEnvCommand(args[0], args[1:], manager)
	case "init":
		return runInteractiveSetup(manager)
	case "setup":
		return setupShellIntegrationEnhanced(manager)
	case "status":
		return showPersistenceStatus(manager)
	case "debug":
		return showDebugInfo(manager)
	case "alias":
		return handleAliasCommand(args, manager)
	case "help":
		return showHelp()
	default:
		return errors.Newf(errors.ErrCodeInvalidArgument, "unknown command: %s (use 'gopher help' to see available commands)", command)
	}
}

func listInstalled(manager *inruntime.Manager) error {
	versions, err := manager.ListInstalled()
	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to list installed versions")
	}

	if len(versions) == 0 {
		if *jsonOutput {
			fmt.Println("[]")
		} else {
			fmt.Println("No Go versions installed.")
		}
		return nil
	}

	// Calculate pagination
	totalVersions := len(versions)
	totalPages := (totalVersions + *pageSize - 1) / *pageSize

	// Validate page number
	if *page < 1 {
		*page = 1
	}
	if *page > totalPages && totalPages > 0 {
		*page = totalPages
	}

	// If interactive mode is enabled and not JSON output, start interactive pagination
	if !*noInteractive && !*jsonOutput {
		return listInstalledInteractive(versions, totalPages)
	}

	// Calculate start and end indices
	startIndex := (*page - 1) * *pageSize
	endIndex := startIndex + *pageSize
	if endIndex > totalVersions {
		endIndex = totalVersions
	}

	// Get the page of versions
	pageVersions := versions[startIndex:endIndex]

	if *jsonOutput {
		// For JSON output, include pagination metadata
		result := map[string]any{
			"versions": pageVersions,
			"pagination": map[string]any{
				"current_page": *page,
				"total_pages":  totalPages,
				"page_size":    *pageSize,
				"total_count":  totalVersions,
			},
		}
		return outputJSON(result)
	}

	// Display pagination info
	fmt.Printf("Installed Go versions (page %d of %d, showing %d of %d total):\n",
		*page, totalPages, len(pageVersions), totalVersions)
	fmt.Println()

	// Display versions
	for _, v := range pageVersions {
		fmt.Printf("%s\n", v.ColoredDisplayString())
	}

	// Display pagination controls
	if totalPages > 1 {
		fmt.Println()
		fmt.Printf("Page %d of %d", *page, totalPages)
		if *page > 1 {
			fmt.Printf(" | Use --page %d for previous page", *page-1)
		}
		if *page < totalPages {
			fmt.Printf(" | Use --page %d for next page", *page+1)
		}
		fmt.Println()
		fmt.Printf("Use 'gopher --page-size <number> list' to change page size (current: %d)\n", *pageSize)
		fmt.Println("Use 'gopher --no-interactive list' to disable interactive pagination")
	}

	return nil
}

func listRemote(manager *inruntime.Manager) error {
	versions, err := manager.ListAvailable()
	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to list available versions")
	}

	// Apply filtering if specified
	if *filter != "" {
		versions = filterVersions(versions, *filter)
	}

	// Apply stable filter if specified
	if *stable {
		versions = filterStableVersions(versions)
	}

	// Calculate pagination
	totalVersions := len(versions)
	totalPages := (totalVersions + *pageSize - 1) / *pageSize

	// Validate page number
	if *page < 1 {
		*page = 1
	}
	if *page > totalPages && totalPages > 0 {
		*page = totalPages
	}

	// If interactive mode is enabled and not JSON output, start interactive pagination
	if !*noInteractive && !*jsonOutput {
		return listRemoteInteractive(versions, totalPages)
	}

	// Calculate start and end indices
	startIndex := (*page - 1) * *pageSize
	endIndex := startIndex + *pageSize
	if endIndex > totalVersions {
		endIndex = totalVersions
	}

	// Get the page of versions
	pageVersions := versions[startIndex:endIndex]

	if *jsonOutput {
		// For JSON output, include pagination metadata
		result := map[string]any{
			"versions": pageVersions,
			"pagination": map[string]any{
				"current_page": *page,
				"total_pages":  totalPages,
				"page_size":    *pageSize,
				"total_count":  totalVersions,
				"filter":       *filter,
				"stable_only":  *stable,
			},
		}
		return outputJSON(result)
	}

	// Display pagination info
	fmt.Printf("Available Go versions (page %d of %d, showing %d of %d total):\n",
		*page, totalPages, len(pageVersions), totalVersions)

	if *filter != "" {
		fmt.Printf("Filtered by: '%s'\n", *filter)
	}
	if *stable {
		fmt.Printf("Showing only stable versions\n")
	}
	fmt.Println()

	// Display versions
	for i, v := range pageVersions {
		status := "stable"
		// Check if version is stable (no beta, rc, devel, etc.)
		if strings.Contains(strings.ToLower(v.Version), "beta") ||
			strings.Contains(strings.ToLower(v.Version), "rc") ||
			strings.Contains(strings.ToLower(v.Version), "devel") ||
			strings.Contains(strings.ToLower(v.Version), "alpha") {
			status = "unstable"
		}
		fmt.Printf("  %d. %s (%s)\n", startIndex+i+1, v.Version, status)
	}

	// Display pagination controls
	if totalPages > 1 {
		fmt.Println()
		fmt.Printf("Page %d of %d", *page, totalPages)
		if *page > 1 {
			fmt.Printf(" | Use --page %d for previous page", *page-1)
		}
		if *page < totalPages {
			fmt.Printf(" | Use --page %d for next page", *page+1)
		}
		fmt.Println()
		fmt.Printf("Use --page-size <number> to change page size (current: %d)\n", *pageSize)
	}

	return nil
}

// listRemoteInteractive provides interactive pagination for list-remote command
// listRemoteInteractive provides interactive pagination for VersionInfo lists
func listRemoteInteractive(versions []downloader.VersionInfo, totalPages int) error {
	currentPage := *page
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Calculate start and end indices for current page
		startIndex := (currentPage - 1) * *pageSize
		endIndex := startIndex + *pageSize
		if endIndex > len(versions) {
			endIndex = len(versions)
		}

		// Get the page of versions
		pageVersions := versions[startIndex:endIndex]

		// Clear screen (optional - makes it cleaner)
		fmt.Print("\033[2J\033[H")

		// Display header
		fmt.Printf("Available Go versions (page %d of %d, showing %d of %d total):\n",
			currentPage, totalPages, len(pageVersions), len(versions))

		// Display versions
		for i, v := range pageVersions {
			fmt.Printf("%d. %s\n", startIndex+i+1, v.Version)
		}

		// Display navigation options
		fmt.Printf("\nNavigation options:\n")
		fmt.Printf("  n/next - Next page\n")
		fmt.Printf("  p/prev - Previous page\n")
		fmt.Printf("  g <num> - Go to page number\n")
		fmt.Printf("  q/quit - Exit\n")
		fmt.Printf("  <Enter> - Next page\n")

		// Get user input
		fmt.Printf("\nEnter command: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			// Default to next page
			input = "n"
		}

		// Parse command
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])
		switch command {
		case "q", "quit", "exit":
			return nil
		case "n", "next":
			if currentPage < totalPages {
				currentPage++
			}
		case "p", "prev", "previous":
			if currentPage > 1 {
				currentPage--
			}
		case "g", "goto":
			if len(parts) >= 2 {
				if targetPage, err := strconv.Atoi(parts[1]); err == nil {
					if targetPage >= 1 && targetPage <= totalPages {
						currentPage = targetPage
					}
				}
			}
		default:
			fmt.Printf("Unknown command: %s\n", command)
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

// listRemoteInteractiveString provides interactive pagination for string-based version lists
func listRemoteInteractiveString(versions []string, totalPages int) error {
	currentPage := *page
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Calculate start and end indices for current page
		startIndex := (currentPage - 1) * *pageSize
		endIndex := startIndex + *pageSize
		if endIndex > len(versions) {
			endIndex = len(versions)
		}

		// Get the page of versions
		pageVersions := versions[startIndex:endIndex]

		// Clear screen (optional - makes it cleaner)
		fmt.Print("\033[2J\033[H")

		// Display header
		fmt.Printf("Available Go versions (page %d of %d, showing %d of %d total):\n",
			currentPage, totalPages, len(pageVersions), len(versions))

		// Display versions
		for i, version := range pageVersions {
			fmt.Printf("%d. %s\n", startIndex+i+1, version)
		}

		// Display navigation options
		fmt.Printf("\nNavigation options:\n")
		fmt.Printf("  n/next - Next page\n")
		fmt.Printf("  p/prev - Previous page\n")
		fmt.Printf("  g <num> - Go to page number\n")
		fmt.Printf("  q/quit - Exit\n")
		fmt.Printf("  <Enter> - Next page\n")

		// Get user input
		fmt.Printf("\nEnter command: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			// Default to next page
			input = "n"
		}

		// Parse command
		parts := strings.Fields(input)
		command := strings.ToLower(parts[0])

		switch command {
		case "n", "next":
			if currentPage < totalPages {
				currentPage++
			} else {
				fmt.Println("Already on the last page.")
			}
		case "p", "prev":
			if currentPage > 1 {
				currentPage--
			} else {
				fmt.Println("Already on the first page.")
			}
		case "g":
			if len(parts) > 1 {
				if pageNum, err := strconv.Atoi(parts[1]); err == nil {
					if pageNum >= 1 && pageNum <= totalPages {
						currentPage = pageNum
					} else {
						fmt.Printf("Page number must be between 1 and %d.\n", totalPages)
					}
				} else {
					fmt.Println("Invalid page number.")
				}
			} else {
				fmt.Println("Please specify a page number.")
			}
		case "q", "quit":
			return nil
		default:
			fmt.Println("Invalid command. Use n/next, p/prev, g <num>, or q/quit.")
		}
	}

	return scanner.Err()
}

// listInstalledInteractive provides interactive pagination for list command
func listInstalledInteractive(versions []inruntime.Version, totalPages int) error {
	currentPage := *page
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Calculate start and end indices for current page
		startIndex := (currentPage - 1) * *pageSize
		endIndex := startIndex + *pageSize
		if endIndex > len(versions) {
			endIndex = len(versions)
		}

		// Get the page of versions
		pageVersions := versions[startIndex:endIndex]

		// Clear screen (optional - makes it cleaner)
		fmt.Print("\033[2J\033[H")

		// Display header
		fmt.Printf("Installed Go versions (page %d of %d, showing %d of %d total):\n",
			currentPage, totalPages, len(pageVersions), len(versions))
		fmt.Println()

		// Display versions
		for _, v := range pageVersions {
			fmt.Printf("%s\n", v.ColoredDisplayString())
		}

		// Display navigation options
		fmt.Println()
		fmt.Printf("Page %d of %d\n", currentPage, totalPages)
		fmt.Println("Commands:")
		fmt.Println("  n, next, →     - Next page")
		fmt.Println("  p, prev, ←     - Previous page")
		fmt.Println("  <number>       - Go to specific page")
		fmt.Println("  q, quit, exit  - Quit")
		fmt.Println("  h, help        - Show this help")
		fmt.Print("\nEnter command: ")

		// Read user input
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(strings.ToLower(scanner.Text()))

		// Handle commands
		switch input {
		case "n", "next", "→", "":
			if currentPage < totalPages {
				currentPage++
			} else {
				fmt.Println("Already on the last page!")
				fmt.Print("Press Enter to continue...")
				scanner.Scan()
			}

		case "p", "prev", "←":
			if currentPage > 1 {
				currentPage--
			} else {
				fmt.Println("Already on the first page!")
				fmt.Print("Press Enter to continue...")
				scanner.Scan()
			}

		case "q", "quit", "exit":
			fmt.Println("Goodbye!")
			return nil

		case "h", "help":
			fmt.Println("\nNavigation Help:")
			fmt.Println("  n, next, →     - Go to next page")
			fmt.Println("  p, prev, ←     - Go to previous page")
			fmt.Println("  <number>       - Jump to specific page (1-" + strconv.Itoa(totalPages) + ")")
			fmt.Println("  q, quit, exit  - Exit the program")
			fmt.Println("  h, help        - Show this help")
			fmt.Print("\nPress Enter to continue...")
			scanner.Scan()
			continue

		default:
			// Try to parse as page number
			if pageNum, err := strconv.Atoi(input); err == nil {
				if pageNum >= 1 && pageNum <= totalPages {
					currentPage = pageNum
				} else {
					fmt.Printf("Invalid page number! Please enter a number between 1 and %d.\n", totalPages)
					fmt.Print("Press Enter to continue...")
					scanner.Scan()
				}
			} else {
				fmt.Println("Invalid command! Type 'h' or 'help' for available commands.")
				fmt.Print("Press Enter to continue...")
				scanner.Scan()
			}
		}
	}

	return nil
}

// filterVersions filters versions based on the provided filter text
// filterVersionsString filters a list of version strings based on a filter string
func filterVersionsString(versions []string, filter string) []string {
	if filter == "" {
		return versions
	}

	filter = strings.ToLower(filter)
	var filtered []string

	for _, v := range versions {
		// Check if version number contains the filter
		if strings.Contains(strings.ToLower(v), filter) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

// filterVersions filters a list of VersionInfo based on a filter string
func filterVersions(versions []downloader.VersionInfo, filter string) []downloader.VersionInfo {
	if filter == "" {
		return versions
	}

	filter = strings.ToLower(filter)
	var filtered []downloader.VersionInfo

	for _, v := range versions {
		// Check if version number contains the filter
		if strings.Contains(strings.ToLower(v.Version), filter) {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

// filterStableVersionsString filters out non-stable versions from a list of version strings
func filterStableVersionsString(versions []string) []string {
	var stable []string

	for _, v := range versions {
		// Check if version is stable (no beta, rc, devel, etc.)
		if !strings.Contains(strings.ToLower(v), "beta") &&
			!strings.Contains(strings.ToLower(v), "rc") &&
			!strings.Contains(strings.ToLower(v), "devel") &&
			!strings.Contains(strings.ToLower(v), "alpha") {
			stable = append(stable, v)
		}
	}

	return stable
}

// filterStableVersions filters out non-stable versions from a list of VersionInfo
func filterStableVersions(versions []downloader.VersionInfo) []downloader.VersionInfo {
	var stable []downloader.VersionInfo

	for _, v := range versions {
		// Check if version is stable (no beta, rc, devel, etc.)
		if !strings.Contains(strings.ToLower(v.Version), "beta") &&
			!strings.Contains(strings.ToLower(v.Version), "rc") &&
			!strings.Contains(strings.ToLower(v.Version), "devel") &&
			!strings.Contains(strings.ToLower(v.Version), "alpha") {
			stable = append(stable, v)
		}
	}

	return stable
}

func installVersion(manager *inruntime.Manager, version string) error {
	if err := manager.Install(version); err != nil {
		return errors.Wrapf(err, errors.ErrCodeInstallationFailed, "failed to install version %s", version)
	}
	return nil
}

func uninstallVersion(manager *inruntime.Manager, version string) error {
	spinner := inprogress.NewSpinner(fmt.Sprintf("Uninstalling Go %s", version))
	spinner.Start()

	err := manager.Uninstall(version)

	spinner.Stop()

	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUninstallationFailed, "failed to uninstall version %s", version)
	}

	return nil
}

func useVersion(manager *inruntime.Manager, version string) error {
	fmt.Printf("Switching to Go %s...\n", version)

	if err := manager.Use(version); err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to switch to version %s", version)
	}

	fmt.Printf("Successfully switched to Go %s\n", version)
	return nil
}

func showCurrent(manager *inruntime.Manager) error {
	current, err := manager.GetCurrent()
	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to get current version")
	}

	if *jsonOutput {
		return outputJSON(current)
	}

	fmt.Printf("Current Go version: %s\n", current.String())
	return nil
}

func showSystem(manager *inruntime.Manager) error {
	systemInfo, err := manager.GetSystemInfo()
	if err != nil {
		return errors.Wrapf(err, errors.ErrCodeUnknown, "failed to get system information")
	}

	if *jsonOutput {
		return outputJSON(systemInfo)
	}

	fmt.Println("System Go Information:")
	fmt.Printf("  Version: %s\n", systemInfo.Version)
	fmt.Printf("  GOROOT: %s\n", systemInfo.GOROOT)
	fmt.Printf("  GOPATH: %s\n", systemInfo.GOPATH)
	fmt.Printf("  Executable: %s\n", systemInfo.Executable)
	fmt.Printf("  Valid: %t\n", systemInfo.IsValid)
	return nil
}

func showVersion() error {
	if *jsonOutput {
		versionInfo := map[string]interface{}{
			"version":    appVersion,
			"commit":     appCommit,
			"date":       appDate,
			"built_by":   appBuiltBy,
			"go_version": runtime.Version(),
			"platform":   runtime.GOOS + "/" + runtime.GOARCH,
		}
		return outputJSON(versionInfo)
	}

	if *verbose {
		// Show detailed version info in verbose mode
		fmt.Println(getFullVersionInfo())
	} else {
		// Show simple version
		fmt.Println(getVersionString())
	}

	return nil
}

// showHelp displays comprehensive help information for the Gopher CLI.
//
// It provides detailed information about available commands, options, and usage
// examples. The help output can be formatted as either plain text or JSON
// depending on the --json flag.
//
// The help information includes:
//   - Version information
//   - Description and features
//   - Available commands with descriptions
//   - Command-line options and flags
//   - Usage examples
//   - Configuration information
//
// Returns an error if there's an issue formatting the help output.
func showHelp() error {
	if *jsonOutput {
		helpInfo := map[string]any{
			"version":     appVersion,
			"description": "Go version manager",
			"commands": map[string]string{
				"init":        "Interactive setup wizard for platform-specific configuration",
				"list":        "List installed Go versions (including system)",
				"list-remote": "List available Go versions (with pagination and filtering)",
				"install":     "Install a Go version",
				"uninstall":   "Uninstall a Go version",
				"use":         "Switch to a Go version (use 'system' for system Go)",
				"current":     "Show current Go version",
				"system":      "Show system Go information",
				"alias":       "Manage version aliases (create, list, remove, show)",
				"setup":       "Set up shell integration for persistent Go version switching",
				"status":      "Show persistence status and shell integration info",
				"debug":       "Show debug information for troubleshooting",
				"env":         "Manage environment variables and configuration",
				"version":     "Show gopher version",
				"help":        "Show detailed help information",
			},
			"examples": []string{
				"gopher init",
				"gopher list",
				"gopher install 1.21.0",
				"gopher use 1.21.0",
				"gopher use system",
				"gopher system",
				"gopher uninstall 1.20.7",
				"gopher alias create stable 1.21.0",
				"gopher alias list",
				"gopher use stable",
				"gopher setup",
				"gopher status",
				"gopher debug",
				"gopher env list",
				"gopher list-remote --page-size 5",
				"gopher list-remote --filter '1.21'",
				"gopher list-remote --filter 'stable'",
			},
			"documentation": "https://github.com/molmedoz/gopher",
		}
		return outputJSON(helpInfo)
	}

	// Detailed help output
	fmt.Println("Gopher - Go Version Manager")
	fmt.Println("============================")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  Gopher is a simple, fast, and dependency-free Go version manager.")
	fmt.Println("  It enables you to install, manage, and switch between multiple Go")
	fmt.Println("  versions, including your system-installed Go.")
	fmt.Println()
	fmt.Println("FEATURES:")
	fmt.Println("  • Fast: Minimal dependencies, built with Go standard library only")
	fmt.Println("  • Secure: Cryptographic verification of downloaded Go binaries")
	fmt.Println("  • Simple: Clean CLI interface with intuitive commands")
	fmt.Println("  • Cross-platform: Works on Linux, macOS, and Windows")
	fmt.Println("  • Lightweight: Zero external dependencies beyond Go standard library")
	fmt.Println("  • Auto-cleanup: Automatically manages old versions to save space")
	fmt.Println("  • System Integration: Seamlessly manages both system and Gopher-installed Go versions")
	fmt.Println("  • JSON Support: Full JSON output for scripting and automation")
	fmt.Println("  • Smart Detection: Automatically detects Homebrew, system packages, and manual installations")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  init                    Interactive setup wizard for platform-specific configuration")
	fmt.Println("  list                    List installed Go versions (including system)")
	fmt.Println("  list-remote             List available Go versions (with pagination and filtering)")
	fmt.Println("  install <version>       Install a Go version")
	fmt.Println("  uninstall <version>     Uninstall a Go version")
	fmt.Println("  use <version>           Switch to a Go version (use 'system' for system Go)")
	fmt.Println("  current                 Show current Go version")
	fmt.Println("  system                  Show system Go information")
	fmt.Println("  alias                   Manage version aliases (create, list, remove, show)")
	fmt.Println("  setup                   Set up shell integration for persistent Go version switching")
	fmt.Println("  status                  Show persistence status and shell integration info")
	fmt.Println("  debug                   Show debug information for troubleshooting")
	fmt.Println("  version                 Show gopher version")
	fmt.Println("  help                    Show detailed help information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Interactive setup wizard")
	fmt.Println("  gopher init")
	fmt.Println()
	fmt.Println("  # List all installed versions (interactive by default)")
	fmt.Println("  gopher list")
	fmt.Println("  gopher --no-interactive list")
	fmt.Println()
	fmt.Println("  # Install and use Go 1.21.0")
	fmt.Println("  gopher install 1.21.0")
	fmt.Println("  gopher use 1.21.0")
	fmt.Println()
	fmt.Println("  # Switch to system Go")
	fmt.Println("  gopher use system")
	fmt.Println()
	fmt.Println("  # Show system Go information")
	fmt.Println("  gopher system")
	fmt.Println()
	fmt.Println("  # Set up persistent Go version switching")
	fmt.Println("  gopher setup")
	fmt.Println("  gopher status")
	fmt.Println()
	fmt.Println("  # Debug information")
	fmt.Println("  gopher debug")
	fmt.Println()
	fmt.Println("  # Remove old version")
	fmt.Println("  gopher uninstall 1.20.7")
	fmt.Println()
	fmt.Println("  # Pagination and filtering")
	fmt.Println("  gopher list-remote --page-size 5")
	fmt.Println("  gopher list-remote --page 2 --page-size 10")
	fmt.Println("  gopher list-remote --filter '1.21'")
	fmt.Println("  gopher list-remote --stable")
	fmt.Println("  gopher list-remote --interactive")
	fmt.Println("  gopher list-remote --filter 'rc'")
	fmt.Println()
	fmt.Println("  # Environment management")
	fmt.Println("  gopher env list")
	fmt.Println("  gopher env show go1.21.0")
	fmt.Println("  gopher env set gopath_mode=version-specific")
	fmt.Println("  gopher env set custom_gopath=/path/to/workspace")
	fmt.Println("  gopher env reset")
	fmt.Println()
	fmt.Println("  # JSON output for scripting")
	fmt.Println("  gopher list --json")
	fmt.Println("  gopher current --json")
	fmt.Println("  gopher system --json")
	fmt.Println("  gopher env show go1.21.0 --json")
	fmt.Println()
	fmt.Println("CONFIGURATION:")
	fmt.Println("  Gopher stores its configuration in:")
	fmt.Println("  • Linux/macOS: ~/.gopher/config.json")
	fmt.Printf("  • Windows: %s\\gopher\\config.json\n", "%USERPROFILE%")
	fmt.Println()
	fmt.Println("  Environment variables:")
	fmt.Println("  • GOPHER_CONFIG: Path to custom configuration file")
	fmt.Println("  • GOPHER_INSTALL_DIR: Custom installation directory")
	fmt.Println("  • GOPHER_DOWNLOAD_DIR: Custom download directory")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --json                  Output in JSON format")
	fmt.Println("  --config <path>         Path to configuration file")
	fmt.Println("  --help                  Show this help message")
	fmt.Println("  --verbose, -v           Show detailed output (DEBUG level)")
	fmt.Println("  --quiet, -q             Only show errors (ERROR level)")
	fmt.Println()
	fmt.Println("PAGINATION & FILTERING (for list-remote):")
	fmt.Println("  --page-size <number>    Number of versions per page (default: 10)")
	fmt.Println("  --page <number>         Page number to display (default: 1)")
	fmt.Println("  --filter <text>         Filter versions by text (e.g., '1.21', 'stable', 'rc')")
	fmt.Println("  --stable                Show only stable versions")
	fmt.Println("  --interactive           Enable interactive pagination (wait for user input)")
	fmt.Println()
	fmt.Println("DOCUMENTATION:")
	fmt.Println("  https://github.com/molmedoz/gopher")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/molmedoz/gopher")
	return nil
}

// showEnvHelp shows help for the env command
func showEnvHelp() error {
	fmt.Println("Environment Management Commands:")
	fmt.Println()
	fmt.Println("  gopher env show [version]     - Show environment variables for a version")
	fmt.Println("  gopher env set <key>=<value>  - Set a configuration option")
	fmt.Println("  gopher env list               - List all configuration options")
	fmt.Println("  gopher env reset              - Reset to default configuration")
	fmt.Println()
	fmt.Println("Configuration Options:")
	fmt.Println("  gopath_mode                  - GOPATH management: shared, version-specific, custom")
	fmt.Println("  custom_gopath                - Custom GOPATH when mode is 'custom'")
	fmt.Println("  goproxy                      - Go proxy URL")
	fmt.Println("  gosumdb                      - Go checksum database")
	fmt.Println("  set_environment              - Whether to set environment variables")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gopher env show go1.21.0")
	fmt.Println("  gopher env set gopath_mode=version-specific")
	fmt.Println("  gopher env set custom_gopath=/path/to/go/workspace")
	fmt.Println("  gopher env list")
	return nil
}

// handleEnvCommand handles environment-related commands
func handleEnvCommand(subcommand string, args []string, manager *inruntime.Manager) error {
	switch subcommand {
	case "show":
		if len(args) < 1 {
			return errors.NewMissingArgument("show (requires version)")
		}
		return showEnvForVersion(args[0], manager)
	case "set":
		if len(args) < 1 {
			return errors.NewMissingArgument("set (requires key=value)")
		}
		return setConfigOption(args[0], manager)
	case "list":
		return listConfigOptions(manager)
	case "reset":
		return resetConfig(manager)
	default:
		return errors.Newf(errors.ErrCodeInvalidArgument, "unknown env subcommand: %s", subcommand)
	}
}

// showEnvForVersion shows environment variables for a specific version
func showEnvForVersion(version string, manager *inruntime.Manager) error {
	// Normalize version
	version = strings.TrimPrefix(version, "go")
	version = "go" + version

	// Get environment variables
	envVars := manager.GetConfig().GetEnvironmentVariables(version)

	if *jsonOutput {
		return outputJSON(map[string]any{
			"version":     version,
			"environment": envVars,
		})
	}

	fmt.Printf("Environment variables for Go %s:\n", version)
	fmt.Println()
	for key, value := range envVars {
		fmt.Printf("  %s=%s\n", key, value)
	}

	return nil
}

// setConfigOption sets a configuration option
func setConfigOption(keyValue string, manager *inruntime.Manager) error {
	if err := errors.ValidateKeyValuePair(keyValue); err != nil {
		return err
	}

	parts := strings.SplitN(keyValue, "=", 2)

	key, value := parts[0], parts[1]

	// Update config based on key
	config := manager.GetConfig()
	switch key {
	case "gopath_mode":
		if err := errors.ValidateConfigValue(key, value); err != nil {
			return err
		}
		config.GOPATHMode = value
	case "custom_gopath":
		config.CustomGOPATH = value
	case "goproxy":
		config.GOPROXY = value
	case "gosumdb":
		config.GOSUMDB = value
	case "set_environment":
		if err := errors.ValidateConfigValue(key, value); err != nil {
			return err
		}
		config.SetEnvironment = value == "true"
	default:
		return errors.NewUnknownConfigOption(key)
	}

	// Save config
	configPath := getConfigPath()
	if err := config.Save(configPath); err != nil {
		return errors.NewConfigSaveFailed(configPath, err)
	}

	fmt.Printf("✓ Configuration updated: %s=%s\n", key, value)
	return nil
}

// listConfigOptions lists all configuration options
func listConfigOptions(manager *inruntime.Manager) error {
	config := manager.GetConfig()

	if *jsonOutput {
		return outputJSON(config)
	}

	fmt.Println("Current Configuration:")
	fmt.Println()
	fmt.Printf("  Install Directory: %s\n", config.InstallDir)
	fmt.Printf("  Download Directory: %s\n", config.DownloadDir)
	fmt.Printf("  Mirror URL: %s\n", config.MirrorURL)
	fmt.Printf("  Auto Cleanup: %t\n", config.AutoCleanup)
	fmt.Printf("  Max Versions: %d\n", config.MaxVersions)
	fmt.Printf("  GOPATH Mode: %s\n", config.GOPATHMode)
	fmt.Printf("  Custom GOPATH: %s\n", config.CustomGOPATH)
	fmt.Printf("  GOPROXY: %s\n", config.GOPROXY)
	fmt.Printf("  GOSUMDB: %s\n", config.GOSUMDB)
	fmt.Printf("  Set Environment: %t\n", config.SetEnvironment)

	return nil
}

// resetConfig resets configuration to defaults
func resetConfig(manager *inruntime.Manager) error {
	config := config.DefaultConfig()
	configPath := getConfigPath()

	if err := config.Save(configPath); err != nil {
		return fmt.Errorf("failed to save default configuration: %w", err)
	}

	fmt.Println("✓ Configuration reset to defaults")
	return nil
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	if *configPath != "" {
		return *configPath
	}
	return config.GetConfigPath()
}

func outputJSON(data any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// setupShellIntegration sets up shell integration for persistent Go version switching
func setupShellIntegration(manager *inruntime.Manager) error {
	// Create the gopher initialization script
	initScript, err := createGopherInitScript(manager)
	if err != nil {
		return fmt.Errorf("failed to create gopher init script: %w", err)
	}

	// Detect the user's shell
	shell := detectShell()
	if shell == "" {
		return fmt.Errorf("unable to detect shell")
	}

	// Check if we're in a Docker container
	isDocker := false
	if _, err := os.Stat("/.dockerenv"); err == nil {
		isDocker = true
	}

	profilePath, err := getShellProfile(shell)
	if err != nil {
		// If we can't get a shell profile (e.g., in Docker), provide alternative instructions
		if isDocker {
			if *jsonOutput {
				return outputJSON(map[string]any{
					"status":      "docker_environment",
					"shell":       shell,
					"init_script": initScript,
					"message":     "Docker environment detected. Manual setup required.",
					"instructions": []string{
						"Add the following to your shell profile:",
						fmt.Sprintf("source %s", initScript),
						"Or run: source " + initScript,
					},
				})
			}

			fmt.Printf("✓ Gopher init script created: %s\n", initScript)
			fmt.Println("Docker environment detected - manual setup required:")
			fmt.Println()
			fmt.Printf("1. Add this line to your shell profile:\n   source %s\n", initScript)
			fmt.Println()
			fmt.Printf("2. Or run this command in each session:\n   source %s\n", initScript)
			fmt.Println()
			fmt.Println("3. To make it persistent, add the source command to:")
			fmt.Println("   - ~/.bashrc (for bash)")
			fmt.Println("   - ~/.zshrc (for zsh)")
			fmt.Println("   - ~/.config/fish/config.fish (for fish)")
			fmt.Println("   - ~/.profile (for sh)")
			return nil
		}
		return fmt.Errorf("failed to get shell profile: %w", err)
	}

	// Add gopher initialization to the shell profile
	if err := addToShellProfile(profilePath, initScript); err != nil {
		return fmt.Errorf("failed to add to shell profile: %w", err)
	}

	if *jsonOutput {
		return outputJSON(map[string]any{
			"status":       "success",
			"shell":        shell,
			"profile_path": profilePath,
			"init_script":  initScript,
			"message":      "Shell integration configured successfully",
		})
	}

	fmt.Printf("✓ Shell integration configured for %s\n", shell)
	fmt.Printf("  Profile: %s\n", profilePath)
	fmt.Printf("  Init script: %s\n", initScript)
	fmt.Printf("  Please restart your shell or run: source %s\n", profilePath)

	return nil
}

// showPersistenceStatus shows the current persistence status and shell integration info
func showPersistenceStatus(manager *inruntime.Manager) error {
	// Check if state file exists
	stateFile := filepath.Join(manager.GetConfig().InstallDir, "..", "state", "active-version")
	stateExists := false
	var activeVersion string

	if content, err := os.ReadFile(stateFile); err == nil {
		stateExists = true
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "active_version=") {
				activeVersion = strings.TrimPrefix(line, "active_version=")
				break
			}
		}
	}

	// Check shell integration
	shell := detectShell()
	profilePath, _ := getShellProfile(shell)
	profileExists := false
	integrationExists := false

	if content, err := os.ReadFile(profilePath); err == nil {
		profileExists = true
		integrationExists = strings.Contains(string(content), "gopher-init.sh")
	}

	// Check init script
	initScript := filepath.Join(manager.GetConfig().InstallDir, "..", "scripts", "gopher-init.sh")
	initScriptExists := false
	if _, err := os.Stat(initScript); err == nil {
		initScriptExists = true
	}

	status := map[string]any{
		"persistence": map[string]any{
			"enabled":        stateExists,
			"active_version": activeVersion,
			"state_file":     stateFile,
		},
		"shell_integration": map[string]any{
			"shell":           shell,
			"profile_path":    profilePath,
			"profile_exists":  profileExists,
			"integration_set": integrationExists,
			"init_script":     initScript,
			"script_exists":   initScriptExists,
		},
	}

	if *jsonOutput {
		return outputJSON(status)
	}

	// Human-readable output
	fmt.Println("Gopher Persistence Status")
	fmt.Println("========================")
	fmt.Println()

	// Persistence status
	fmt.Println("Persistence:")
	if stateExists {
		fmt.Printf("  ✓ Enabled (active version: %s)\n", activeVersion)
	} else {
		fmt.Println("  ✗ Disabled")
	}
	fmt.Printf("  State file: %s\n", stateFile)
	fmt.Println()

	// Shell integration status
	fmt.Println("Shell Integration:")
	if shell != "" {
		fmt.Printf("  Shell: %s\n", shell)
		fmt.Printf("  Profile: %s\n", profilePath)
		if profileExists {
			fmt.Printf("  Profile exists: ✓\n")
		} else {
			fmt.Printf("  Profile exists: ✗\n")
		}
		if integrationExists {
			fmt.Printf("  Integration configured: ✓\n")
		} else {
			fmt.Printf("  Integration configured: ✗\n")
		}
	} else {
		fmt.Println("  Shell: Unknown")
	}
	fmt.Printf("  Init script: %s\n", initScript)
	if initScriptExists {
		fmt.Printf("  Script exists: ✓\n")
	} else {
		fmt.Printf("  Script exists: ✗\n")
	}
	fmt.Println()

	// Recommendations
	switch {
	case !stateExists:
		fmt.Println("Recommendations:")
		fmt.Println("  • Run 'gopher use <version>' to enable persistence")
		fmt.Println("  • Run 'gopher setup' to configure shell integration")
	case !integrationExists:
		fmt.Println("Recommendations:")
		fmt.Println("  • Run 'gopher setup' to configure shell integration")
		fmt.Println("  • Restart your shell after setup")
	default:
		fmt.Println("✓ Persistence and shell integration are properly configured!")
	}

	return nil
}

// Helper functions for shell integration (copied from manager.go for CLI access)

func createGopherInitScript(manager *inruntime.Manager) (string, error) {
	scriptDir := filepath.Join(manager.GetConfig().InstallDir, "..", "scripts")
	if err := os.MkdirAll(scriptDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create script directory: %w", err)
	}

	scriptPath := filepath.Join(scriptDir, "gopher-init.sh")

	scriptContent := `#!/bin/bash
# Gopher Go Version Manager - Shell Integration
# This script is automatically generated and should not be edited manually

# Function to get the active Go version
gopher_get_active_version() {
    local state_file="$HOME/.gopher/state/active-version"
    if [[ -f "$state_file" ]]; then
        local version=$(grep "active_version=" "$state_file" | cut -d'=' -f2)
        if [[ -n "$version" ]]; then
            echo "$version"
            return 0
        fi
    fi
    return 1
}

# Function to set up Go environment
gopher_setup_go_env() {
    local version="$1"
    if [[ -z "$version" ]]; then
        return 1
    fi

    # Handle system version
    if [[ "$version" == "system" ]]; then
        # Use system Go - try to find the original system Go installation
        # Check common system Go locations first
        local system_go_locations=(
            "/usr/local/go/bin/go"
            "/usr/lib/go/bin/go"
            "/opt/go/bin/go"
        )
        
        local go_path=""
        for location in "${system_go_locations[@]}"; do
            if [[ -x "$location" ]]; then
                go_path="$location"
                break
            fi
        done
        
        # If not found in common locations, try to find via PATH but avoid symlinks
        if [[ -z "$go_path" ]] && command -v go >/dev/null 2>&1; then
            local path_go=$(command -v go)
            # Only use if it's not a symlink (to avoid gopher-managed versions)
            if [[ ! -L "$path_go" ]]; then
                go_path="$path_go"
            fi
        fi
        
        if [[ -n "$go_path" ]]; then
            local goroot_dir=$(dirname $(dirname "$go_path"))
            export GOROOT="$goroot_dir"
            export PATH="$goroot_dir/bin:$PATH"
        else
            return 1
        fi
    else
        # Set up GOROOT for gopher-managed version
        local goroot="$HOME/.gopher/versions/$version"
        if [[ -d "$goroot" ]]; then
            export GOROOT="$goroot"
            export PATH="$goroot/bin:$PATH"
        else
            return 1
        fi
    fi

    # Set up GOPATH based on configuration
    local gopath_mode="shared"  # Default mode
    local state_file="$HOME/.gopher/state/active-version"
    if [[ -f "$state_file" ]]; then
        local config_file="$HOME/.gopher/config.json"
        if [[ -f "$config_file" ]]; then
            # Try to read GOPATH mode from config (simplified)
            local mode=$(grep -o '"gopath_mode":"[^"]*"' "$config_file" | cut -d'"' -f4)
            if [[ -n "$mode" ]]; then
                gopath_mode="$mode"
            fi
        fi
    fi

    case "$gopath_mode" in
        "shared")
            export GOPATH="$HOME/go"
            ;;
        "version-specific")
            export GOPATH="$HOME/.gopher/gopath/$version"
            ;;
        "custom")
            # For custom mode, we'd need to read from config
            export GOPATH="$HOME/go"
            ;;
    esac

    # Set up other Go environment variables
    export GOPROXY="https://proxy.golang.org,direct"
    export GOSUMDB="sum.golang.org"
    
    # Add GOPATH/bin to PATH for Go tools
    export PATH="$GOPATH/bin:$PATH"
    
    return 0
}

# Auto-initialize gopher when shell starts
if [[ -z "$GOPHER_INITIALIZED" ]]; then
    active_version=$(gopher_get_active_version)
    if [[ -n "$active_version" ]]; then
        gopher_setup_go_env "$active_version"
        export GOPHER_INITIALIZED=1
    fi
fi

# Gopher command aliases
alias gopher-use='gopher use'
alias gopher-list='gopher list'
alias gopher-current='gopher current'
alias gopher-install='gopher install'
alias gopher-uninstall='gopher uninstall'
`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return "", fmt.Errorf("failed to write init script: %w", err)
	}

	return scriptPath, nil
}

func detectShell() string {
	// Windows: Detect PowerShell
	if runtime.GOOS == "windows" {
		// Check if running in PowerShell by checking PSModulePath env var
		if os.Getenv("PSModulePath") != "" {
			// Check if it's PowerShell Core (pwsh) or Windows PowerShell
			if _, err := exec.LookPath("pwsh"); err == nil {
				return "pwsh"
			}
			return "powershell"
		}
		// Default to PowerShell on Windows
		return "powershell"
	}

	// Unix/Linux/macOS: Existing detection logic
	// First try the SHELL environment variable
	shell := os.Getenv("SHELL")
	if shell != "" {
		shellName := filepath.Base(shell)
		switch shellName {
		case "bash", "zsh", "fish", "sh":
			return shellName
		default:
			// Try to detect from path
			return detectShellFromPath(shell)
		}
	}

	// If SHELL is not set, try to detect from process
	// This is common in Docker containers
	if shellPath, err := os.Readlink("/proc/self/exe"); err == nil {
		return detectShellFromPath(shellPath)
	}

	// Try to detect from parent process
	if parentShell := os.Getenv("0"); parentShell != "" {
		return detectShellFromPath(parentShell)
	}

	// Check if we're in a Docker container and default to bash
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "bash"
	}

	// Check for common shell binaries in PATH
	shells := []string{"bash", "zsh", "fish", "sh"}
	for _, shellName := range shells {
		if _, err := exec.LookPath(shellName); err == nil {
			return shellName
		}
	}

	// Last resort: default to bash
	return "bash"
}

func detectShellFromPath(shellPath string) string {
	if shellPath == "" {
		return ""
	}
	shellName := filepath.Base(shellPath)
	switch {
	case strings.Contains(shellName, "bash"):
		return "bash"
	case strings.Contains(shellName, "zsh"):
		return "zsh"
	case strings.Contains(shellName, "fish"):
		return "fish"
	default:
		return ""
	}
}

func getShellProfile(shell string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	switch shell {
	case "powershell", "pwsh":
		// PowerShell profile (both Windows PowerShell 5.x and PowerShell Core 6+)
		// Use the same profile path for both
		return filepath.Join(homeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
	case "bash":
		// Use .bashrc as the primary profile file
		return filepath.Join(homeDir, ".bashrc"), nil
	case "zsh":
		return filepath.Join(homeDir, ".zshrc"), nil
	case "fish":
		return filepath.Join(homeDir, ".config", "fish", "config.fish"), nil
	case "sh":
		// For sh, try .profile first, then .bashrc as fallback
		profile := filepath.Join(homeDir, ".profile")
		if _, err := os.Stat(profile); err == nil {
			return profile, nil
		}
		return filepath.Join(homeDir, ".bashrc"), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

func addToShellProfile(profilePath, initScript string) error {
	// Check if gopher is already in the profile
	content, err := os.ReadFile(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	profileContent := string(content)
	if strings.Contains(profileContent, "gopher-init.sh") {
		// Already configured
		return nil
	}

	// Add gopher initialization
	initLine := fmt.Sprintf("\n# Gopher Go Version Manager\nsource %s\n", initScript)

	// Append to profile
	if err := os.WriteFile(profilePath, []byte(profileContent+initLine), 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	return nil
}

func showDebugInfo(manager *inruntime.Manager) error {
	fmt.Println("=== Gopher Debug Information ===")
	fmt.Println()

	// Show current Go version
	fmt.Println("Current Go version:")
	if current, err := manager.GetCurrent(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  %s\n", current.Version)
	}
	fmt.Println()

	// Show which go executable is being used
	fmt.Println("Go executable path:")
	if goPath, err := exec.LookPath("go"); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  %s\n", goPath)
	}
	fmt.Println()

	// Show PATH
	fmt.Println("PATH environment variable:")
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		fmt.Println("  (empty)")
	} else {
		pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))
		for i, dir := range pathDirs {
			if dir != "" {
				fmt.Printf("  %d. %s\n", i+1, dir)
			}
		}
	}
	fmt.Println()

	// Show gopher symlink locations
	fmt.Println("Gopher symlink locations:")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("  Error getting home directory: %v\n", err)
	} else {
		var symlinkPaths []string
		if runtime.GOOS == "windows" {
			symlinkPaths = []string{
				filepath.Join(homeDir, "AppData", "Local", "bin", "go.exe"),
				filepath.Join(homeDir, "bin", "go.exe"),
			}
		} else {
			symlinkPaths = []string{
				"/usr/local/bin/go",
				"/usr/bin/go",
				filepath.Join(homeDir, ".local", "bin", "go"),
				filepath.Join(homeDir, "bin", "go"),
			}
		}

		for _, symlinkPath := range symlinkPaths {
			if _, err := os.Lstat(symlinkPath); err == nil {
				if target, err := os.Readlink(symlinkPath); err == nil {
					fmt.Printf("  ✓ %s -> %s\n", symlinkPath, target)
				} else {
					fmt.Printf("  ✗ %s (not a symlink)\n", symlinkPath)
				}
			} else {
				fmt.Printf("  - %s (does not exist)\n", symlinkPath)
			}
		}
	}
	fmt.Println()

	// Show installed versions
	fmt.Println("Installed Go versions:")
	installed, err := manager.ListInstalled()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		for _, v := range installed {
			fmt.Printf("  - %s\n", v.Version)
		}
	}
	fmt.Println()

	// Show system Go info
	fmt.Println("System Go information:")
	if systemInfo, err := manager.GetSystemInfo(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Available: %v\n", systemInfo != nil)
		if systemInfo != nil {
			fmt.Printf("  Version: %s\n", systemInfo.Version)
			fmt.Printf("  Executable: %s\n", systemInfo.Executable)
			fmt.Printf("  GOROOT: %s\n", systemInfo.GOROOT)
		}
	}

	return nil
}

// handleAliasCommand handles the alias command and its subcommands
func handleAliasCommand(args []string, manager *inruntime.Manager) error {
	if len(args) == 0 {
		return showAliasHelp()
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "create", "add":
		if len(subArgs) < 2 {
			return fmt.Errorf("alias create requires name and version (e.g., 'gopher alias create stable 1.21.0')")
		}
		return createAlias(manager, subArgs[0], subArgs[1])
	case "list", "ls":
		return listAliases(manager)
	case "show", "get":
		if len(subArgs) < 1 {
			return fmt.Errorf("alias show requires an alias name (e.g., 'gopher alias show stable')")
		}
		return showAlias(manager, subArgs[0])
	case "remove", "rm", "delete":
		if len(subArgs) < 1 {
			return fmt.Errorf("alias remove requires an alias name (e.g., 'gopher alias remove stable')")
		}
		return removeAlias(manager, subArgs[0])
	case "update":
		if len(subArgs) < 2 {
			return fmt.Errorf("alias update requires name and version (e.g., 'gopher alias update stable 1.22.0')")
		}
		return updateAlias(manager, subArgs[0], subArgs[1])
	case "bulk":
		return handleBulkAliasCommand(subArgs, manager)
	case "by-version":
		if len(subArgs) < 1 {
			return fmt.Errorf("version required for 'by-version' subcommand")
		}
		return showAliasesByVersion(manager, subArgs[0])
	case "suggest":
		if len(subArgs) < 1 {
			return fmt.Errorf("version required for 'suggest' subcommand")
		}
		return suggestAliases(manager, subArgs[0])
	case "export":
		return handleAliasExport(subArgs, manager)
	case "import":
		return handleAliasImport(subArgs, manager)
	case "help":
		return showAliasHelp()
	default:
		return fmt.Errorf("unknown alias subcommand: %s (use 'gopher alias help' for available commands)", subcommand)
	}
}

// showAliasesByVersion shows all aliases for a specific version
func showAliasesByVersion(manager *inruntime.Manager, version string) error {
	aliases, err := manager.AliasManager().GetAliasesByVersion(version)
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		fmt.Printf("No aliases found for version %s\n", version)
		return nil
	}

	fmt.Printf("Aliases for version %s:\n", version)
	for _, alias := range aliases {
		fmt.Printf("  %s -> %s (created: %s)\n",
			alias.Name,
			alias.Version,
			alias.Created.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// suggestAliases suggests common alias names for a version
func suggestAliases(manager *inruntime.Manager, version string) error {
	suggestions := manager.AliasManager().SuggestAliases(version)

	if len(suggestions) == 0 {
		fmt.Printf("No suggestions available for version %s\n", version)
		return nil
	}

	fmt.Printf("Suggested aliases for version %s:\n", version)
	for i, suggestion := range suggestions {
		fmt.Printf("  %d. %s\n", i+1, suggestion)
	}

	fmt.Println("\nUsage examples:")
	for _, suggestion := range suggestions[:min(3, len(suggestions))] {
		fmt.Printf("  gopher alias create %s %s\n", suggestion, version)
	}

	return nil
}

// handleAliasExport handles the export command
func handleAliasExport(args []string, manager *inruntime.Manager) error {
	if len(args) < 1 {
		return fmt.Errorf("filename required for export (e.g., 'gopher alias export aliases.json')")
	}

	filename := args[0]
	format := "json"
	tags := []string{}

	// Parse additional flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--tags", "-t":
			if i+1 < len(args) {
				tags = strings.Split(args[i+1], ",")
				i++
			} else {
				return fmt.Errorf("tags value required after --tags")
			}
		}
	}

	if err := manager.AliasManager().ExportAliases(filename, format, tags); err != nil {
		return err
	}

	fmt.Printf("✓ Exported aliases to %s (%s format)\n", filename, format)
	return nil
}

// handleAliasImport handles the import command
func handleAliasImport(args []string, manager *inruntime.Manager) error {
	if len(args) < 1 {
		return fmt.Errorf("filename required for import (e.g., 'gopher alias import aliases.json')")
	}

	filename := args[0]

	// Determine conflict resolution mode
	allowOverride := *override
	noOverride := *noOverride
	force := *force

	// Force overrides all other flags
	if force {
		allowOverride = false
		noOverride = false
	}

	if err := manager.AliasManager().ImportAliases(filename, allowOverride, noOverride, force); err != nil {
		return err
	}

	fmt.Printf("✓ Imported aliases from %s\n", filename)
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// showAliasHelp shows help for the alias command
func showAliasHelp() error {
	fmt.Println(`gopher alias - Manage version aliases

USAGE:
    gopher alias <subcommand> [arguments]

SUBCOMMANDS:
    create <name> <version>    Create a new alias (e.g., 'gopher alias create stable 1.21.0')
    list                      List all aliases
    show <name>               Show details of a specific alias
    by-version <version>      Show all aliases for a specific version
    suggest <version>         Suggest common alias names for a version
    export <file>             Export aliases to JSON file
    import <file>             Import aliases from JSON file
    remove <name>             Remove an alias
    update <name> <version>   Update an existing alias
    bulk                      Bulk alias operations (create multiple aliases)
    help                      Show this help

EXAMPLES:
    gopher alias create stable 1.21.0
    gopher alias create latest 1.22.0
    gopher alias list
    gopher alias show stable
    gopher alias by-version 1.21.0    # Show all aliases for version 1.21.0
    gopher alias suggest 1.21.0       # Suggest common aliases for version 1.21.0
    gopher alias export aliases.json  # Export aliases to JSON file
    gopher alias import aliases.json  # Import aliases from JSON file
    gopher alias remove stable
    gopher alias update stable 1.22.0
    gopher use stable          # Use an alias with the 'use' command
    
    # Interactive conflict resolution (default behavior)
    gopher alias create stable 1.22.0
    gopher alias update stable 1.23.0
    
    # Override mode (no confirmation)
    gopher alias create --override stable 1.22.0
    gopher alias update --override stable 1.23.0
    
    # No override mode (exit on conflicts)
    gopher alias create --no-override stable 1.22.0
    gopher alias update --no-override stable 1.23.0
    
    # Force mode (overrides all other flags)
    gopher alias create --force stable 1.22.0
    gopher alias update --force stable 1.23.0
    
    # Bulk operations
    gopher alias bulk create stable=1.21.0 latest=1.22.0 dev=1.23.0
    gopher alias bulk create --override stable=1.22.0 latest=1.23.0

ALIAS NAMING RULES:
    - Only letters, numbers, dots, hyphens, and underscores allowed
    - 1-50 characters long
    - Cannot use reserved names (commands, 'system', 'sys', etc.)
    - Case-sensitive

RESERVED NAMES:
    list, install, uninstall, use, current, system, version, help, init, setup, status, debug, env, alias, sys, go, git, ls, cd, pwd`)
	return nil
}

// createAlias creates a new alias
func createAlias(manager *inruntime.Manager, name, version string) error {
	// Determine conflict resolution mode
	allowOverride := *override
	noOverride := *noOverride
	force := *force

	// Force overrides all other flags
	if force {
		allowOverride = false
		noOverride = false
	}

	if err := manager.AliasManager().CreateAliasInteractive(name, version, allowOverride, noOverride, force); err != nil {
		return err
	}

	// Show success message for non-interactive modes
	if force || allowOverride || noOverride {
		fmt.Printf("✓ Created alias '%s' -> %s\n", name, version)
	}
	return nil
}

// listAliases lists all aliases
func listAliases(manager *inruntime.Manager) error {
	aliases, err := manager.AliasManager().ListAliases()
	if err != nil {
		return fmt.Errorf("failed to list aliases: %w", err)
	}

	if len(aliases) == 0 {
		fmt.Println("No aliases found.")
		fmt.Println("Create one with: gopher alias create <name> <version>")
		return nil
	}

	fmt.Printf("Found %d alias(es):\n", len(aliases))
	fmt.Println()

	for _, alias := range aliases {
		fmt.Printf("  %-20s -> %s\n", alias.Name, alias.Version)
		fmt.Printf("    Created: %s\n", alias.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Updated: %s\n", alias.Updated.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

// showAlias shows details of a specific alias
func showAlias(manager *inruntime.Manager, name string) error {
	alias, exists := manager.AliasManager().GetAlias(name)
	if !exists {
		return fmt.Errorf("alias '%s' not found", name)
	}

	fmt.Printf("Alias: %s\n", alias.Name)
	fmt.Printf("Version: %s\n", alias.Version)
	fmt.Printf("Created: %s\n", alias.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", alias.Updated.Format("2006-01-02 15:04:05"))

	return nil
}

// removeAlias removes an alias
func removeAlias(manager *inruntime.Manager, name string) error {
	if err := manager.AliasManager().RemoveAlias(name); err != nil {
		return err
	}

	fmt.Printf("✓ Removed alias '%s'\n", name)
	return nil
}

// updateAlias updates an existing alias
func updateAlias(manager *inruntime.Manager, name, version string) error {
	// Determine conflict resolution mode
	allowOverride := *override
	noOverride := *noOverride
	force := *force

	// Force overrides all other flags
	if force {
		allowOverride = false
		noOverride = false
	}

	if err := manager.AliasManager().UpdateAliasInteractive(name, version, allowOverride, noOverride, force); err != nil {
		return err
	}

	// Show success message for non-interactive modes
	if force || allowOverride || noOverride {
		fmt.Printf("✓ Updated alias '%s' -> %s\n", name, version)
	}
	return nil
}

// handleBulkAliasCommand handles the bulk alias command
func handleBulkAliasCommand(args []string, manager *inruntime.Manager) error {
	if len(args) < 1 {
		return showBulkAliasHelp()
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "create":
		return createBulkAliases(manager, subArgs)
	case "help":
		return showBulkAliasHelp()
	default:
		return fmt.Errorf("unknown bulk subcommand: %s (use 'gopher alias bulk help' for available commands)", subcommand)
	}
}

// showBulkAliasHelp shows help for the bulk alias command
func showBulkAliasHelp() error {
	fmt.Println(`gopher alias bulk - Bulk alias operations

USAGE:
    gopher alias bulk create <name1=version1> [name2=version2] ...

EXAMPLES:
    gopher alias bulk create stable=1.21.0 latest=1.22.0 dev=1.23.0
    gopher alias bulk create --override stable=1.21.0 latest=1.22.0
    gopher alias bulk create --no-override stable=1.21.0 latest=1.22.0
    gopher alias bulk create --force stable=1.21.0 latest=1.22.0

OPTIONS:
    --override     Allow overriding existing aliases without confirmation
    --no-override  Exit with error if alias already exists (no override allowed)
    --force        Force operation without confirmation (overrides all other flags)

BULK OPERATION FEATURES:
    • Create multiple aliases in one command
    • Interactive conflict resolution for existing aliases
    • Summary of new aliases and conflicts
    • Batch confirmation for updates`)
	return nil
}

// createBulkAliases creates multiple aliases
func createBulkAliases(manager *inruntime.Manager, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("bulk create requires at least one alias (e.g., 'gopher alias bulk create stable=1.21.0')")
	}

	// Parse aliases from arguments
	aliases := make(map[string]string)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid alias format '%s' (expected 'name=version')", arg)
		}
		aliases[parts[0]] = parts[1]
	}

	// Determine conflict resolution mode
	allowOverride := *override
	noOverride := *noOverride
	force := *force

	// Force overrides all other flags
	if force {
		allowOverride = false
		noOverride = false
	}

	// Create aliases
	if err := manager.AliasManager().CreateAliasesBulk(aliases, allowOverride, noOverride, force); err != nil {
		return err
	}

	return nil
}

// setupShellIntegrationEnhanced provides an enhanced setup experience
func setupShellIntegrationEnhanced(manager *inruntime.Manager) error {
	fmt.Println("🔧 Gopher Environment Setup")
	fmt.Println("===========================")
	fmt.Println()

	// Detect system information
	systemInfo, err := detectSystemInfo(manager)
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}

	// Show current status
	fmt.Println("📋 Current Status")
	fmt.Println("=================")
	fmt.Printf("Platform: %s/%s\n", systemInfo.Platform, systemInfo.Arch)
	fmt.Printf("Shell: %s\n", systemInfo.Shell)
	fmt.Printf("Symlink Directory: %s\n", systemInfo.SymlinkDir)
	fmt.Println()

	// Windows-specific setup instructions
	if runtime.GOOS == "windows" {
		fmt.Println("🪟 Windows Configuration")
		fmt.Println("========================")
		fmt.Println()

		// Check if symlink dir is in PATH
		if !systemInfo.IsInPath {
			fmt.Println("⚠️  STEP 1: Add Gopher's bin directory to PATH (REQUIRED)")
			fmt.Println()
			fmt.Println("Copy and run this PowerShell command as Administrator:")
			fmt.Println()
			fmt.Printf("  $userPath = [Environment]::GetEnvironmentVariable(\"PATH\", \"User\")\n")
			fmt.Printf("  $gopherBin = \"%s\"\n", systemInfo.SymlinkDir)
			fmt.Printf("  if ($userPath -notlike \"*$gopherBin*\") {\n")
			fmt.Printf("    [Environment]::SetEnvironmentVariable(\"PATH\", \"$gopherBin;$userPath\", \"User\")\n")
			fmt.Printf("  }\n")
			fmt.Println()
			fmt.Println("Then RESTART your PowerShell terminal and run 'gopher setup' again.")
			fmt.Println()
		} else {
			fmt.Println("✅ Gopher's bin directory is in PATH")
			fmt.Println()
		}

		fmt.Println("STEP 2: Install and use a Go version")
		fmt.Println("  gopher install 1.21.0")
		fmt.Println("  gopher use 1.21.0")
		fmt.Println()

		fmt.Println("⚠️  STEP 3: Fix PATH order (if you have system Go)")
		fmt.Println("  Gopher will automatically check PATH order and warn you if needed.")
		fmt.Println("  Follow the PowerShell commands it provides.")
		fmt.Println()

		fmt.Println("═══════════════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("ℹ️  Windows Note:")
		fmt.Println("  Shell integration (.bashrc, .sh files) is NOT needed on Windows.")
		fmt.Println("  Gopher uses symlinks which work automatically once PATH is set.")
		fmt.Println()

		return nil
	}

	// Unix/Linux/macOS - show shell integration status
	if systemInfo.IsInPath {
		fmt.Println("✅ Symlink directory is in PATH")
	} else {
		fmt.Println("⚠️  Symlink directory is not in PATH")
	}

	if isGopherConfigured(systemInfo.ShellProfile) {
		fmt.Println("✅ Shell integration already configured")
	} else {
		fmt.Println("❌ Shell integration not configured")
	}
	fmt.Println()

	// Ask if user wants to proceed
	if !askForConfirmation("Do you want to set up shell integration?") {
		fmt.Println("Setup cancelled.")
		return nil
	}

	// Setup shell integration
	fmt.Println("\n🔧 Setting up shell integration...")

	// Create the gopher initialization script
	initScript, err := createGopherInitScript(manager)
	if err != nil {
		return fmt.Errorf("failed to create gopher init script: %w", err)
	}

	// Add to shell profile
	if err := addToShellProfile(systemInfo.ShellProfile, initScript); err != nil {
		return fmt.Errorf("failed to add to shell profile: %w", err)
	}

	// Add symlink directory to PATH if needed
	if !systemInfo.IsInPath {
		fmt.Printf("Adding %s to PATH...\n", systemInfo.SymlinkDir)
		if err := addDirectoryToPath(systemInfo.SymlinkDir, systemInfo.ShellProfile); err != nil {
			fmt.Printf("⚠️  Failed to add to PATH: %v\n", err)
			fmt.Printf("Please manually add this to your %s:\n", systemInfo.ShellProfile)
			fmt.Printf("export PATH=\"%s:$PATH\"\n", systemInfo.SymlinkDir)
		} else {
			fmt.Printf("✅ Added %s to PATH\n", systemInfo.SymlinkDir)
		}
	}

	fmt.Printf("✅ Shell integration configured in %s\n", systemInfo.ShellProfile)
	fmt.Printf("✅ Gopher init script created: %s\n", initScript)

	// Test the setup
	fmt.Println("\n🧪 Testing setup...")
	if err := testSymlinkCreation(systemInfo.SymlinkDir); err != nil {
		fmt.Printf("⚠️  Symlink test failed: %v\n", err)
		if systemInfo.Platform == "windows" && !systemInfo.HasDeveloperMode {
			fmt.Println("Enable Developer Mode or run as Administrator")
		}
	} else {
		fmt.Println("✅ Symlink creation works")
	}

	// Show next steps
	fmt.Println("\n🎉 Setup Complete!")
	fmt.Println("==================")
	fmt.Println("Shell integration is now configured.")
	fmt.Println()
	fmt.Println("🚀 Next Steps:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("STEP 1: Activate shell integration")
	fmt.Printf("  source %s\n", systemInfo.ShellProfile)
	fmt.Println("  (Or restart your terminal)")
	fmt.Println()
	fmt.Println("STEP 2: Install a Go version")
	fmt.Println("  gopher install 1.21.0")
	fmt.Println()
	fmt.Println("STEP 3: Switch to it")
	fmt.Println("  gopher use 1.21.0")
	fmt.Println()
	fmt.Println("STEP 4: Verify it works")
	fmt.Println("  go version")
	fmt.Println("  # Should show: go version go1.21.0")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("📁 Directories created:")
	fmt.Printf("  Config:    %s/.gopher/config.json\n", systemInfo.HomeDir)
	fmt.Printf("  Versions:  %s/.gopher/versions\n", systemInfo.HomeDir)
	fmt.Printf("  Downloads: %s/.gopher/downloads\n", systemInfo.HomeDir)
	fmt.Printf("  Symlinks:  %s\n", systemInfo.SymlinkDir)
	fmt.Println()

	if systemInfo.IsDocker {
		fmt.Println("🐳 Docker Environment:")
		fmt.Println("- Run 'source ~/.gopher/scripts/gopher-init.sh' in each session")
		fmt.Println("- Or add it to your shell profile manually")
		fmt.Println()
	}

	return nil
}
