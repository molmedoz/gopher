package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/molmedoz/gopher/internal/config"
	"github.com/molmedoz/gopher/internal/env"
)

// Helper function to write metadata files
func writeMetadata(t *testing.T, dir, version string) {
	t.Helper()
	vdir := filepath.Join(dir, version)
	// #nosec G301 -- 0755 acceptable for test directory
	if err := os.MkdirAll(vdir, 0755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(vdir, ".gopher-metadata")
	content := "version=" + version + "\n" +
		"os=linux\n" +
		"arch=amd64\n" +
		"installed_at=2023-01-01T00:00:00Z\n" +
		"install_dir=" + vdir + "\n"
	// #nosec G306 -- 0644 acceptable for test files
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// Helper function to create test manager
func createTestManager(t *testing.T, tmp string) *Manager {
	cfg := &config.Config{
		InstallDir:  tmp,
		DownloadDir: filepath.Join(tmp, "dl"),
		MirrorURL:   "https://go.dev/dl/",
	}
	// Use mock environment provider for testing
	mockEnv := env.NewMockProvider(map[string]string{
		"PATH":  "/usr/local/bin:/usr/bin:/bin",
		"SHELL": "/bin/bash",
	})
	// Use NewManager to properly initialize all components
	manager := NewManager(cfg, mockEnv)
	return manager
}
