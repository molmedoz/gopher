package runtime

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/molmedoz/gopher/internal/config"
)

func TestAliasManager_CreateAlias(t *testing.T) {
	tests := []struct {
		name      string
		aliasName string
		version   string
		wantErr   bool
		errorMsg  string
		setup     func(*AliasManager) // Setup function to run before test
	}{
		{"valid alias", "stable", "go1.21.0", false, "", nil},
		{"duplicate alias", "stable", "go1.22.0", true, "ALIAS_ALREADY_EXISTS: alias 'stable' already exists (use 'gopher alias remove stable' first)", func(am *AliasManager) {
			// Create the first alias to test duplicate detection
			_ = am.CreateAlias("stable", "go1.21.0")
		}},
		{"invalid version", "test", "invalid-version", true, "INVALID_VERSION: invalid version format: invalid-version", nil},
		{"empty alias name", "", "go1.21.0", true, "INVALID_ALIAS_NAME: alias name cannot be empty", nil},
		{"reserved name", "system", "go1.21.0", true, "RESERVED_NAME: name 'system' is reserved", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new AliasManager for each test to reset sync.Once
			tmp := t.TempDir()
			cfg := &config.Config{
				InstallDir: filepath.Join(tmp, "install"),
			}
			am := NewAliasManager(cfg)

			// Run setup if provided
			if tt.setup != nil {
				tt.setup(am)
			}

			err := am.CreateAlias(tt.aliasName, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errorMsg != "" {
				if err.Error() != tt.errorMsg {
					t.Errorf("CreateAlias() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

func TestAliasManager_GetAlias(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		t.Fatal(err)
	}

	// Add a test alias
	am.aliases["test"] = &Alias{
		Name:    "test",
		Version: "go1.21.0",
		Created: time.Now(),
		Updated: time.Now(),
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		t.Fatal(err)
	}

	// Test getting existing alias
	alias, exists := am.GetAlias("test")
	if !exists {
		t.Fatal("expected alias to exist")
	}
	if alias.Name != "test" {
		t.Errorf("expected alias name 'test', got %s", alias.Name)
	}

	// Test getting non-existent alias
	_, exists = am.GetAlias("nonexistent")
	if exists {
		t.Fatal("expected alias to not exist")
	}
}

func TestAliasManager_ListAliases(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		t.Fatal(err)
	}

	// Add test aliases
	am.aliases["alias1"] = &Alias{
		Name:    "alias1",
		Version: "go1.21.0",
		Created: time.Now(),
		Updated: time.Now(),
	}
	am.aliases["alias2"] = &Alias{
		Name:    "alias2",
		Version: "go1.22.0",
		Created: time.Now(),
		Updated: time.Now(),
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		t.Fatal(err)
	}

	// Test listing aliases
	aliases, err := am.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases error: %v", err)
	}

	if len(aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(aliases))
	}
}

func TestAliasManager_RemoveAlias(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		t.Fatal(err)
	}

	// Add a test alias
	am.aliases["test"] = &Alias{
		Name:    "test",
		Version: "go1.21.0",
		Created: time.Now(),
		Updated: time.Now(),
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		t.Fatal(err)
	}

	// Test removing existing alias
	err := am.RemoveAlias("test")
	if err != nil {
		t.Fatalf("RemoveAlias error: %v", err)
	}

	// Verify alias was removed
	if _, exists := am.aliases["test"]; exists {
		t.Fatal("expected alias to be removed")
	}

	// Test removing non-existent alias
	err = am.RemoveAlias("nonexistent")
	if err == nil {
		t.Fatal("expected error for removing non-existent alias")
	}
}

func TestAliasManager_UpdateAlias(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		t.Fatal(err)
	}

	// Add a test alias
	am.aliases["test"] = &Alias{
		Name:    "test",
		Version: "go1.21.0",
		Created: time.Now(),
		Updated: time.Now(),
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		t.Fatal(err)
	}

	// Test updating existing alias
	err := am.UpdateAlias("test", "go1.22.0")
	if err != nil {
		t.Fatalf("UpdateAlias error: %v", err)
	}

	// Verify alias was updated
	alias := am.aliases["test"]
	if alias.Version != "go1.22.0" {
		t.Errorf("expected version 'go1.22.0', got %s", alias.Version)
	}

	// Test updating non-existent alias
	err = am.UpdateAlias("nonexistent", "go1.22.0")
	if err == nil {
		t.Fatal("expected error for updating non-existent alias")
	}
}

func TestAliasManager_GetAliasesByVersion(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Load aliases first
	if err := am.LoadAliases(); err != nil {
		t.Fatal(err)
	}

	// Add test aliases
	am.aliases["alias1"] = &Alias{
		Name:    "alias1",
		Version: "go1.21.0",
		Created: time.Now(),
		Updated: time.Now(),
	}
	am.aliases["alias2"] = &Alias{
		Name:    "alias2",
		Version: "go1.21.0",
		Created: time.Now(),
		Updated: time.Now(),
	}
	am.aliases["alias3"] = &Alias{
		Name:    "alias3",
		Version: "go1.22.0",
		Created: time.Now(),
		Updated: time.Now(),
	}

	// Save aliases
	if err := am.SaveAliases(); err != nil {
		t.Fatal(err)
	}

	// Test getting aliases by version
	aliases, err := am.GetAliasesByVersion("go1.21.0")
	if err != nil {
		t.Fatalf("GetAliasesByVersion error: %v", err)
	}

	if len(aliases) != 2 {
		t.Errorf("expected 2 aliases for go1.21.0, got %d", len(aliases))
	}
}

func TestAliasManager_LoadAliases(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Test loading from non-existent file
	err := am.LoadAliases()
	if err != nil {
		t.Fatalf("LoadAliases error: %v", err)
	}

	// Should have empty aliases map
	if len(am.aliases) != 0 {
		t.Errorf("expected empty aliases map, got %d aliases", len(am.aliases))
	}
}

func TestAliasManager_SaveAliases(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Add a test alias
	am.aliases["test"] = &Alias{
		Name:    "test",
		Version: "go1.21.0",
		Created: time.Now(),
		Updated: time.Now(),
	}

	// Save aliases
	err := am.SaveAliases()
	if err != nil {
		t.Fatalf("SaveAliases error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(am.aliasesFile); os.IsNotExist(err) {
		t.Fatal("aliases file was not created")
	}
}

func TestAliasManager_ValidateAliasName(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	tests := []struct {
		name      string
		aliasName string
		wantErr   bool
		errorMsg  string
	}{
		{"valid name", "stable", false, ""},
		{"valid name with numbers", "go1.21", false, ""},
		{"valid name with hyphens", "my-alias", false, ""},
		{"valid name with underscores", "my_alias", false, ""},
		{"valid name with dots", "my.alias", false, ""},
		{"empty name", "", true, "alias name cannot be empty"},
		{"reserved name", "system", true, "'system' is a reserved name and cannot be used as an alias"},
		{"reserved name", "list", true, "'list' is a reserved name and cannot be used as an alias"},
		{"invalid characters", "my@alias", true, "alias name contains invalid characters"},
		{"too long", "this_is_a_very_long_alias_name_that_exceeds_fifty_characters", true, "alias name is too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := am.ValidateAliasName(tt.aliasName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAliasName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errorMsg != "" {
				if err.Error() != tt.errorMsg {
					t.Errorf("ValidateAliasName() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

func TestAliasManager_IsVersionInstalled(t *testing.T) {
	tmp := t.TempDir()
	cfg := &config.Config{
		InstallDir: filepath.Join(tmp, "install"),
	}
	am := NewAliasManager(cfg)

	// Test with nil manager (should return true as fallback)
	result := am.isVersionInstalled("go1.21.0")
	if !result {
		t.Error("expected isVersionInstalled to return true with nil manager")
	}
}
