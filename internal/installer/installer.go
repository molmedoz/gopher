package installer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/progress"
	"github.com/molmedoz/gopher/internal/security"
)

// Installer handles installing Go versions
type Installer struct {
	installDir string
}

// New creates a new installer
func New(installDir string) *Installer {
	return &Installer{
		installDir: installDir,
	}
}

// Install installs a Go version from a downloaded file
func (i *Installer) Install(version, filePath string) error {
	// Print installation start message
	fmt.Printf("Installing Go %s\n", version)

	// Validate input paths for security
	if err := security.ValidatePath(version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}
	if err := security.ValidatePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	if err := security.ValidateDirectoryPath(i.installDir); err != nil {
		return fmt.Errorf("invalid install directory: %w", err)
	}

	// Ensure install directory exists
	if err := os.MkdirAll(i.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Determine target directory
	targetDir := filepath.Join(i.installDir, version)

	// Remove existing installation if it exists
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to remove existing installation: %w", err)
	}

	// Extract the archive with progress
	if err := i.extractArchive(filePath, targetDir); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Create version metadata with spinner
	metadataSpinner := progress.NewSpinner("Creating version metadata")
	metadataSpinner.Start()
	err := i.createVersionMetadata(version, targetDir)
	metadataSpinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to create version metadata: %w", err)
	}

	fmt.Printf("✓ Successfully installed Go %s\n", version)
	return nil
}

// Uninstall removes a Go version
func (i *Installer) Uninstall(version string) error {
	// Validate input paths for security
	if err := security.ValidatePath(version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}
	if err := security.ValidateDirectoryPath(i.installDir); err != nil {
		return fmt.Errorf("invalid install directory: %w", err)
	}

	targetDir := filepath.Join(i.installDir, version)

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed (use 'gopher list' to see installed versions)", version)
	}

	return os.RemoveAll(targetDir)
}

// IsInstalled checks if a version is installed
func (i *Installer) IsInstalled(version string) bool {
	targetDir := filepath.Join(i.installDir, version)
	_, err := os.Stat(targetDir)
	return !os.IsNotExist(err)
}

// ListInstalled returns a list of installed versions
func (i *Installer) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(i.installDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read install directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	return versions, nil
}

// extractArchive extracts a Go archive to the target directory
func (i *Installer) extractArchive(filePath, targetDir string) error {
	spinner := progress.NewSpinner("Extracting archive")
	spinner.Start()
	defer spinner.Stop()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	// Determine archive type and extract accordingly
	ext := filepath.Ext(filePath)
	switch ext {
	case ".gz":
		if strings.HasSuffix(filePath, ".tar.gz") {
			err = i.extractTarGz(file, targetDir)
		} else {
			err = fmt.Errorf("unsupported .gz format: %s", filePath)
		}
	case ".zip":
		err = i.extractZip(filePath, targetDir)
	case ".msi":
		err = i.extractMSI(filePath, targetDir)
	default:
		err = fmt.Errorf("unsupported archive format: %s", ext)
	}

	if err != nil {
		return err
	}

	return nil
}

// extractTarGz extracts a tar.gz archive
func (i *Installer) extractTarGz(file *os.File, targetDir string) error {
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	var hasGoPrefix bool
	var hasGoBinary bool
	var goBinaryName string
	if runtime.GOOS == "windows" {
		goBinaryName = "go.exe"
	} else {
		goBinaryName = "go"
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Check if this archive has the required "go/" prefix
		if strings.HasPrefix(header.Name, "go/") {
			hasGoPrefix = true
		}

		// Check if this archive contains the go binary
		if strings.HasSuffix(header.Name, "/bin/"+goBinaryName) {
			hasGoBinary = true
		}

		// Skip the root "go" directory and extract contents directly
		path := strings.TrimPrefix(header.Name, "go/")

		targetPath := filepath.Join(targetDir, path)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to copy file content: %w", err)
			}

			outFile.Close()

			// Set file permissions
			if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to set file permissions: %w", err)
			}
		}
	}

	// Validate archive structure
	if !hasGoPrefix {
		return fmt.Errorf("archive does not have required 'go/' prefix")
	}
	if !hasGoBinary {
		return fmt.Errorf("archive does not contain go binary")
	}

	return nil
}

// extractZip extracts a ZIP archive
func (i *Installer) extractZip(filePath, targetDir string) error {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	var hasGoPrefix bool
	var hasGoBinary bool
	var goBinaryName string
	if runtime.GOOS == "windows" {
		goBinaryName = "go.exe"
	} else {
		goBinaryName = "go"
	}

	for _, file := range reader.File {
		// Check if this archive has the required "go/" prefix
		if strings.HasPrefix(file.Name, "go/") {
			hasGoPrefix = true
		}

		// Skip the root "go" directory and extract contents directly
		path := strings.TrimPrefix(file.Name, "go/")
		targetPath := filepath.Join(targetDir, path)

		// Check if this archive contains the go binary (after trimming go/ prefix)
		if strings.HasSuffix(path, "/bin/"+goBinaryName) || path == "bin/"+goBinaryName {
			hasGoBinary = true
		}

		// Skip empty directories
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, file.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Extract file
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		outFile, err := os.Create(targetPath)
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create file: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}

		// Set file permissions
		if err := os.Chmod(targetPath, file.FileInfo().Mode()); err != nil {
			return fmt.Errorf("failed to set file permissions: %w", err)
		}
	}

	// Validate archive structure
	if !hasGoPrefix {
		return fmt.Errorf("archive does not have required 'go/' prefix")
	}
	if !hasGoBinary {
		return fmt.Errorf("archive does not contain go binary")
	}

	return nil
}

// extractMSI extracts a Windows MSI file
// Note: This is a simplified implementation. In a real implementation,
// you might want to use a proper MSI extraction library or call
// Windows APIs through CGO.
func (i *Installer) extractMSI(filePath, targetDir string) error {
	// For now, we'll return an error for MSI files on non-Windows systems
	if runtime.GOOS != "windows" {
		return fmt.Errorf("MSI extraction is only supported on Windows")
	}

	// This is a placeholder. In a real implementation, you would:
	// 1. Use Windows APIs to extract the MSI
	// 2. Or use a third-party library like github.com/akavel/rsrc
	// 3. Or call msiexec.exe programmatically

	return fmt.Errorf("MSI extraction not implemented yet")
}

// createVersionMetadata creates metadata for the installed version
func (i *Installer) createVersionMetadata(version, targetDir string) error {
	metadata := map[string]any{
		"version":      version,
		"os":           runtime.GOOS,
		"arch":         runtime.GOARCH,
		"installed_at": time.Now().Format(time.RFC3339),
		"install_dir":  targetDir,
	}

	// Write metadata to a file
	metadataPath := filepath.Join(targetDir, ".gopher-metadata")
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	// Simple metadata format (could be JSON, YAML, etc.)
	for key, value := range metadata {
		if _, err := fmt.Fprintf(file, "%s=%v\n", key, value); err != nil {
			return fmt.Errorf("failed to write metadata: %w", err)
		}
	}

	return nil
}

// GetVersionMetadata reads metadata for an installed version
func (i *Installer) GetVersionMetadata(version string) (map[string]string, error) {
	targetDir := filepath.Join(i.installDir, version)
	metadataPath := filepath.Join(targetDir, ".gopher-metadata")

	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	metadata := make(map[string]string)

	// Read metadata line by line
	var line string
	for {
		_, err := fmt.Fscanf(file, "%s\n", &line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read metadata: %w", err)
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			metadata[parts[0]] = parts[1]
		}
	}

	return metadata, nil
}

// GetGoBinaryPath returns the path to the go binary for a version
func (i *Installer) GetGoBinaryPath(version string) (string, error) {
	targetDir := filepath.Join(i.installDir, version)

	// Look for go binary in the expected location
	binaryName := "go"
	if runtime.GOOS == "windows" {
		binaryName = "go.exe"
	}

	binaryPath := filepath.Join(targetDir, "bin", binaryName)

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return "", fmt.Errorf("go binary not found for version %s (installation may be corrupted)", version)
	}

	return binaryPath, nil
}
