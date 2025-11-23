package runtime

import (
	"path/filepath"

	"github.com/molmedoz/gopher/internal/config"
	"github.com/molmedoz/gopher/internal/downloader"
	"github.com/molmedoz/gopher/internal/env"
	"github.com/molmedoz/gopher/internal/installer"
)

// NewManager creates a new version manager with all required dependencies.
//
// It initializes a Manager instance with concrete implementations of all
// required interfaces, using the adapter pattern to maintain clean separation
// between the core logic and external dependencies.
//
// The manager is configured with:
//   - Downloader: For downloading Go releases from official mirrors
//   - Installer: For installing and uninstalling Go versions
//   - SystemDetector: For detecting system-installed Go versions
//   - SymlinkManager: For managing version switching via symlinks
//   - FileSystem: For file system operations
//   - AliasManager: For managing version aliases
//
// Parameters:
//   - cfg: Configuration containing paths, URLs, and settings
//   - envProvider: Environment variable provider for configuration
//
// Returns a fully initialized Manager ready for use.
//
// Example:
//
//	cfg := config.Load("/path/to/config.json")
//	manager := NewManager(cfg, envProvider)
//	err := manager.Install("1.21.0")
func NewManager(cfg *config.Config, envProvider env.Provider) *Manager {
	manager := &Manager{
		config:       cfg,
		downloader:   downloader.New(cfg.MirrorURL),
		installer:    installer.New(cfg.InstallDir),
		aliasManager: nil, // Will be set below
		envProvider:  envProvider,
	}

	// Create alias manager with manager reference
	manager.aliasManager = NewAliasManagerWithManager(cfg, manager)

	return manager
}

// NewAliasManagerWithManager creates a new alias manager with manager reference
func NewAliasManagerWithManager(config *config.Config, manager *Manager) *AliasManager {
	// Get safe root directory (parent of InstallDir, e.g., ~/.gopher or ~/gopher)
	// This avoids path traversal via parent directory access
	installDirAbs, err := filepath.Abs(config.InstallDir)
	if err != nil {
		// Fallback to simple path join if abs resolution fails
	aliasesFile := filepath.Join(filepath.Dir(config.InstallDir), "aliases.json")
		return &AliasManager{
			config:      config,
			aliases:     make(map[string]*Alias),
			aliasesFile: aliasesFile,
			manager:     manager,
		}
	}
	safeRoot := filepath.Dir(installDirAbs) // Parent of versions directory (e.g., ~/.gopher)
	aliasesFile := filepath.Join(safeRoot, "aliases.json")
	return &AliasManager{
		config:      config,
		aliases:     make(map[string]*Alias),
		aliasesFile: aliasesFile,
		manager:     manager,
	}
}

// NewAliasManager creates a new alias manager (legacy constructor for backward compatibility)
func NewAliasManager(cfg *config.Config) *AliasManager {
	// Get safe root directory (parent of InstallDir, e.g., ~/.gopher or ~/gopher)
	// This avoids path traversal via parent directory access
	installDirAbs, err := filepath.Abs(cfg.InstallDir)
	if err != nil {
		// Fallback to simple path join if abs resolution fails
	aliasesFile := filepath.Join(filepath.Dir(cfg.InstallDir), "aliases.json")
		return &AliasManager{
			config:      cfg,
			aliases:     make(map[string]*Alias),
			aliasesFile: aliasesFile,
			manager:     nil, // Will be set when used with a manager
		}
	}
	safeRoot := filepath.Dir(installDirAbs) // Parent of versions directory (e.g., ~/.gopher)
	aliasesFile := filepath.Join(safeRoot, "aliases.json")
	return &AliasManager{
		config:      cfg,
		aliases:     make(map[string]*Alias),
		aliasesFile: aliasesFile,
		manager:     nil, // Will be set when used with a manager
	}
}
