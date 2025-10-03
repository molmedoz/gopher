package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_Install_AlreadyInstalled(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	// Create a version directory to simulate already installed
	versionDir := filepath.Join(tmp, "go1.21.0")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create metadata file
	writeMetadata(t, tmp, "go1.21.0")

	// Try to install the same version
	err := m.Install("go1.21.0")
	if err == nil {
		t.Fatal("expected error for already installed version")
	}

	expectedError := "VERSION_ALREADY_INSTALLED: version go1.21.0 is already installed"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestManager_Uninstall_NotInstalled(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	// Try to uninstall a version that's not installed
	err := m.Uninstall("go1.21.0")
	if err == nil {
		t.Fatal("expected error for not installed version")
	}

	expectedError := "VERSION_NOT_INSTALLED: version go1.21.0 is not installed"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestManager_Uninstall_Flows(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	// Not installed path
	if err := m.Uninstall("go1.2.3"); err == nil {
		t.Fatal("expected error for uninstalling non-existent version")
	}

	// Create a version to uninstall
	writeMetadata(t, tmp, "go1.21.0")

	// Now it should be installed
	installed, err := m.IsInstalled("go1.21.0")
	if err != nil {
		t.Fatal(err)
	}
	if !installed {
		t.Fatal("expected version to be installed")
	}

	// Uninstall should work
	if err := m.Uninstall("go1.21.0"); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}

	// Should no longer be installed
	installed, err = m.IsInstalled("go1.21.0")
	if err != nil {
		t.Fatal(err)
	}
	if installed {
		t.Fatal("expected version to not be installed after uninstall")
	}
}

func TestManager_IsInstalled(t *testing.T) {
	tmp := t.TempDir()
	m := createTestManager(t, tmp)

	// Test with no versions installed
	installed, err := m.IsInstalled("go1.21.0")
	if err != nil {
		t.Fatal(err)
	}
	if installed {
		t.Fatal("expected version to not be installed")
	}

	// Create a version
	writeMetadata(t, tmp, "go1.21.0")

	// Now it should be installed
	installed, err = m.IsInstalled("go1.21.0")
	if err != nil {
		t.Fatal(err)
	}
	if !installed {
		t.Fatal("expected version to be installed")
	}
}
