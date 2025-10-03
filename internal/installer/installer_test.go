package installer

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsInstalledFalseWhenMissing(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)
	if inst.IsInstalled("go1.99.99") {
		t.Fatalf("expected not installed")
	}
}

func TestGetGoBinaryPath(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)
	bin := filepath.Join(tdir, "go1.2.3", "bin")
	if err := os.MkdirAll(bin, 0755); err != nil {
		t.Fatal(err)
	}
	goBin := filepath.Join(bin, "go")
	if err := os.WriteFile(goBin, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}
	path, err := inst.GetGoBinaryPath("go1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != goBin {
		t.Fatalf("got %q want %q", path, goBin)
	}
}

func TestCreateVersionMetadata(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)
	ver := "go9.9.9"
	target := filepath.Join(tdir, ver)
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	if err := inst.createVersionMetadata(ver, target); err != nil {
		t.Fatalf("createVersionMetadata error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, ".gopher-metadata")); err != nil {
		t.Fatalf("metadata file not created: %v", err)
	}
}

func TestGetVersionMetadata_ReadBack(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)
	ver := "go7.7.7"
	target := filepath.Join(tdir, ver)
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	if err := inst.createVersionMetadata(ver, target); err != nil {
		t.Fatalf("createVersionMetadata error: %v", err)
	}
	meta, err := inst.GetVersionMetadata(ver)
	if err != nil {
		t.Fatalf("GetVersionMetadata error: %v", err)
	}
	if meta["version"] != ver {
		t.Fatalf("version meta mismatch: %q", meta["version"])
	}
}

func TestUninstall_NotInstalled(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)
	if err := inst.Uninstall("go0.0.1"); err == nil {
		t.Fatalf("expected error when uninstalling non-existent version")
	}
}

func TestUninstall_RemovesDir(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)
	ver := "go3.3.3"
	vdir := filepath.Join(tdir, ver)
	if err := os.MkdirAll(filepath.Join(vdir, "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := inst.Uninstall(ver); err != nil {
		t.Fatalf("unexpected uninstall error: %v", err)
	}
	if _, err := os.Stat(vdir); !os.IsNotExist(err) {
		t.Fatalf("expected version dir removed")
	}
}

func createTarGz(t *testing.T, files map[string][]byte) string {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, data := range files {
		hdr := &tar.Header{Name: name, Mode: 0755, Size: int64(len(data))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	tmp := filepath.Join(t.TempDir(), "go.tar.gz")
	if err := os.WriteFile(tmp, buf.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}
	return tmp
}

func createZip(t *testing.T, files map[string][]byte) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "go.zip")
	zipFile, err := os.Create(tmp)
	if err != nil {
		t.Fatal(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for name, data := range files {
		writer, err := zipWriter.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := writer.Write(data); err != nil {
			t.Fatal(err)
		}
	}

	return tmp
}

func TestInstaller_Install_FromZip(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	// Use the correct binary name for the current platform
	var goBinaryName string
	if runtime.GOOS == "windows" {
		goBinaryName = "go.exe"
	} else {
		goBinaryName = "go"
	}

	// Archive simulates official layout with root "go/" prefix
	zipFile := createZip(t, map[string][]byte{
		"go/bin/" + goBinaryName: []byte("#!/bin/sh\n"),
		"go/VERSION":             []byte("go1.2.3\n"),
	})

	if err := inst.Install("go1.2.3", zipFile); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	// Expect go binary exists under install dir
	bin := filepath.Join(tdir, "go1.2.3", "bin", goBinaryName)
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("installed go binary missing: %v", err)
	}
	// And metadata
	if _, err := os.Stat(filepath.Join(tdir, "go1.2.3", ".gopher-metadata")); err != nil {
		t.Fatalf("metadata missing: %v", err)
	}
}

func TestInstaller_Install_FromTarGz(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)
	// Archive simulates official layout with root "go/" prefix
	tgz := createTarGz(t, map[string][]byte{
		"go/bin/go":  []byte("#!/bin/sh\n"),
		"go/VERSION": []byte("go1.2.3\n"),
	})
	if err := inst.Install("go1.2.3", tgz); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	// Expect go binary exists under install dir
	bin := filepath.Join(tdir, "go1.2.3", "bin", "go")
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("installed go binary missing: %v", err)
	}
	// And metadata
	if _, err := os.Stat(filepath.Join(tdir, "go1.2.3", ".gopher-metadata")); err != nil {
		t.Fatalf("metadata missing: %v", err)
	}
}

func TestInstaller_ListInstalled(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	// Initially should be empty
	versions, err := inst.ListInstalled()
	if err != nil {
		t.Fatalf("ListInstalled error: %v", err)
	}
	if len(versions) != 0 {
		t.Fatalf("expected empty list, got %v", versions)
	}

	// Create some version directories
	versionsToCreate := []string{"go1.20.0", "go1.21.0", "go1.22.0"}
	for _, ver := range versionsToCreate {
		vdir := filepath.Join(tdir, ver)
		if err := os.MkdirAll(filepath.Join(vdir, "bin"), 0755); err != nil {
			t.Fatal(err)
		}
		// Create metadata
		if err := inst.createVersionMetadata(ver, vdir); err != nil {
			t.Fatal(err)
		}
	}

	// Now should list all versions
	versions, err = inst.ListInstalled()
	if err != nil {
		t.Fatalf("ListInstalled error: %v", err)
	}
	if len(versions) != len(versionsToCreate) {
		t.Fatalf("expected %d versions, got %d", len(versionsToCreate), len(versions))
	}

	// Check that all expected versions are present
	found := make(map[string]bool)
	for _, ver := range versions {
		found[ver] = true
	}
	for _, expected := range versionsToCreate {
		if !found[expected] {
			t.Fatalf("expected version %s not found in list", expected)
		}
	}
}

func TestInstaller_GetVersionMetadata_InvalidVersion(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	// Try to get metadata for non-existent version
	_, err := inst.GetVersionMetadata("go999.999.999")
	if err == nil {
		t.Fatalf("expected error for non-existent version")
	}
}

func TestInstaller_GetVersionMetadata_CorruptedMetadata(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	// Create version directory
	ver := "go1.21.0"
	vdir := filepath.Join(tdir, ver)
	if err := os.MkdirAll(vdir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create corrupted metadata file
	metadataFile := filepath.Join(vdir, ".gopher-metadata")
	corruptedContent := "invalid metadata content\nnot in key=value format"
	if err := os.WriteFile(metadataFile, []byte(corruptedContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Should handle corrupted metadata gracefully
	_, err := inst.GetVersionMetadata(ver)
	if err == nil {
		t.Fatalf("expected error for corrupted metadata")
	}
}

func TestInstaller_Install_InvalidArchive(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	// Create invalid archive
	invalidArchive := filepath.Join(t.TempDir(), "invalid.tar.gz")
	invalidContent := "this is not a valid tar.gz file"
	if err := os.WriteFile(invalidArchive, []byte(invalidContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Should return error
	err := inst.Install("go1.21.0", invalidArchive)
	if err == nil {
		t.Fatalf("expected error for invalid archive")
	}
}

func TestInstaller_Install_ArchiveWithoutGoPrefix(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	// Create archive without "go/" prefix
	tgz := createTarGz(t, map[string][]byte{
		"bin/go":      []byte("#!/bin/sh\n"),
		"VERSION":     []byte("go1.2.3\n"),
		"src/main.go": []byte("package main\n"),
	})

	err := inst.Install("go1.2.3", tgz)
	if err == nil {
		t.Fatalf("expected error for archive without go/ prefix")
	}
}

func TestInstaller_Install_ArchiveWithoutGoBinary(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	// Create archive with go/ prefix but no go binary
	tgz := createTarGz(t, map[string][]byte{
		"go/VERSION":     []byte("go1.2.3\n"),
		"go/src/main.go": []byte("package main\n"),
	})

	err := inst.Install("go1.2.3", tgz)
	if err == nil {
		t.Fatalf("expected error for archive without go binary")
	}
}

func TestInstaller_Uninstall_WithMetadata(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	ver := "go1.21.0"
	vdir := filepath.Join(tdir, ver)
	if err := os.MkdirAll(filepath.Join(vdir, "bin"), 0755); err != nil {
		t.Fatal(err)
	}

	// Create metadata
	if err := inst.createVersionMetadata(ver, vdir); err != nil {
		t.Fatal(err)
	}

	// Verify it's installed
	if !inst.IsInstalled(ver) {
		t.Fatalf("version should be installed")
	}

	// Uninstall
	if err := inst.Uninstall(ver); err != nil {
		t.Fatalf("uninstall error: %v", err)
	}

	// Verify it's no longer installed
	if inst.IsInstalled(ver) {
		t.Fatalf("version should not be installed after uninstall")
	}

	// Verify directory is removed
	if _, err := os.Stat(vdir); !os.IsNotExist(err) {
		t.Fatalf("version directory should be removed")
	}
}

func TestInstaller_Install_OverwriteExisting(t *testing.T) {
	tdir := t.TempDir()
	inst := New(tdir)

	ver := "go1.21.0"
	vdir := filepath.Join(tdir, ver)
	if err := os.MkdirAll(filepath.Join(vdir, "bin"), 0755); err != nil {
		t.Fatal(err)
	}

	// Create initial metadata
	if err := inst.createVersionMetadata(ver, vdir); err != nil {
		t.Fatal(err)
	}

	// Create archive to install
	tgz := createTarGz(t, map[string][]byte{
		"go/bin/go":  []byte("#!/bin/sh\n"),
		"go/VERSION": []byte("go1.21.0\n"),
	})

	// Install should overwrite existing
	if err := inst.Install(ver, tgz); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	// Verify it's still installed
	if !inst.IsInstalled(ver) {
		t.Fatalf("version should still be installed")
	}
}
