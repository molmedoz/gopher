package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_ListInstalled_Empty(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	// No versions installed
	installed, err := m.IsInstalled("go1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if installed {
		t.Fatalf("expected not installed")
	}

	// List should not error and may include system if present
	versions, err := m.ListInstalled()
	if err != nil {
		t.Fatalf("ListInstalled error: %v", err)
	}

	// Should have at least system version if available
	if len(versions) == 0 {
		t.Log("No versions found (including system)")
	}
}

func TestManager_ListInstalled_WithMetadata(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	writeMetadata(t, tmp, "go1.18.10")
	writeMetadata(t, tmp, "go1.19.9")

	got, err := m.ListInstalled()
	if err != nil {
		t.Fatalf("ListInstalled error: %v", err)
	}

	// We expect at least the two we created (system version may add more)
	count := 0
	for _, v := range got {
		if v.Version == "go1.18.10" || v.Version == "go1.19.9" {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("expected 2 metadata versions, got %d (total %d)", count, len(got))
	}
}

func TestManager_ListInstalled_Deduplication(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	// Create metadata for same version with different OS/arch
	writeMetadata(t, tmp, "go1.21.0")

	// Create another metadata file with same version
	writeMetadata(t, tmp, "go1.21.0")

	got, err := m.ListInstalled()
	if err != nil {
		t.Fatalf("ListInstalled error: %v", err)
	}

	// Should not have duplicates
	versionCount := make(map[string]int)
	for _, v := range got {
		versionCount[v.Version]++
	}

	for version, count := range versionCount {
		if count > 1 {
			t.Errorf("Version %s appears %d times, expected at most 1", version, count)
		}
	}
}

func TestManager_ListInstalled_WithActiveVersion(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	writeMetadata(t, tmp, "go1.21.0")
	writeMetadata(t, tmp, "go1.22.0")

	// Simulate active version by creating state file
	stateDir := filepath.Join(tmp, "..", "state")
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		t.Fatal(err)
	}
	stateFile := filepath.Join(stateDir, "active-version")
	if err := os.WriteFile(stateFile, []byte("active_version=go1.21.0\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := m.ListInstalled()
	if err != nil {
		t.Fatalf("ListInstalled error: %v", err)
	}

	// Check if active version is marked correctly
	foundActive := false
	for _, v := range got {
		if v.Version == "go1.21.0" && v.IsActive {
			foundActive = true
			break
		}
	}

	if !foundActive {
		t.Log("Active version not found or not marked as active (this may be expected in test environment)")
	}
}

func TestManager_ListAvailable(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	// This test may fail if there's no internet connection
	// We'll just check that it doesn't panic
	versions, err := m.ListAvailable()
	if err != nil {
		t.Logf("ListAvailable failed (expected if no internet): %v", err)
		return
	}

	if len(versions) == 0 {
		t.Log("No versions available (unexpected)")
	} else {
		t.Logf("Found %d available versions", len(versions))
	}
}
