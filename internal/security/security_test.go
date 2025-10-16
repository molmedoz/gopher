package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errCode string
	}{
		{"valid relative path", "go1.21.0", false, ""},
		{"valid nested path", "versions/go1.21.0", false, ""},
		{"valid absolute path", "/tmp/go1.21.0", false, ""},
		{"empty path", "", true, ErrCodeInvalidPath},
		{"path traversal with dots", "../go1.21.0", true, ErrCodePathTraversal},
		{"path traversal with slashes", "../../go1.21.0", true, ErrCodePathTraversal},
		{"path traversal with backslashes", "..\\go1.21.0", true, ErrCodePathTraversal},
		{"suspicious tilde", "~/go1.21.0", true, ErrCodeUnsafePath},
		{"suspicious dollar", "$HOME/go1.21.0", true, ErrCodeUnsafePath},
		{"suspicious backtick", "`go1.21.0`", true, ErrCodeUnsafePath},
		{"suspicious pipe", "go1.21.0|rm", true, ErrCodeUnsafePath},
		{"suspicious ampersand", "go1.21.0&rm", true, ErrCodeUnsafePath},
		{"suspicious semicolon", "go1.21.0;rm", true, ErrCodeUnsafePath},
		{"suspicious parentheses", "go1.21.0(rm)", true, ErrCodeUnsafePath},
		{"suspicious angle brackets", "go1.21.0<rm", true, ErrCodeUnsafePath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if se, ok := err.(*SecurityError); ok {
					if se.Code != tt.errCode {
						t.Errorf("ValidatePath() error code = %v, want %v", se.Code, tt.errCode)
					}
				} else {
					t.Errorf("ValidatePath() error type = %T, want *SecurityError", err)
				}
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"clean path", "go1.21.0", "go1.21.0"},
		{"path with dots", "../go1.21.0", "go1.21.0"},
		{"path with multiple dots", "../../go1.21.0", "go1.21.0"},
		{"path with backslashes", "..\\go1.21.0", "go1.21.0"},
		{"path with tilde", "~/go1.21.0", "go1.21.0"},
		{"path with dollar", "$HOME/go1.21.0", "HOME/go1.21.0"},
		{"path with backtick", "`go1.21.0`", "go1.21.0"},
		{"path with pipe", "go1.21.0|rm", "go1.21.0rm"},
		{"path with ampersand", "go1.21.0&rm", "go1.21.0rm"},
		{"path with semicolon", "go1.21.0;rm", "go1.21.0rm"},
		{"path with parentheses", "go1.21.0(rm)", "go1.21.0rm"},
		{"path with angle brackets", "go1.21.0<rm", "go1.21.0rm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizePath(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizePath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateDirectoryPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid directory", "versions", false},
		{"valid nested directory", "versions/go1.21.0", false},
		{"valid absolute directory", "/tmp/versions", false},
		{"directory with trailing slash", "versions/", false},
		{"directory with trailing backslash", "versions\\", false},
		{"empty directory", "", true},
		{"path traversal", "../versions", true},
		{"suspicious characters", "versions$", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDirectoryPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDirectoryPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsSafePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"safe path", "go1.21.0", true},
		{"safe nested path", "versions/go1.21.0", true},
		{"unsafe path traversal", "../go1.21.0", false},
		{"unsafe suspicious chars", "go1.21.0$", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSafePath(tt.path)
			if got != tt.want {
				t.Errorf("IsSafePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSafePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{"valid path", "go1.21.0", "go1.21.0", false},
		{"path with dots", "../go1.21.0", "", true},
		{"clean path", "./go1.21.0", "go1.21.0", false},
		{"nested path", "versions/go1.21.0", "versions/go1.21.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSafePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSafePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetSafePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurityError(t *testing.T) {
	tests := []struct {
		name     string
		err      *SecurityError
		expected string
	}{
		{
			"error with details",
			&SecurityError{Code: "TEST", Message: "test error", Details: "details"},
			"test error: details",
		},
		{
			"error without details",
			&SecurityError{Code: "TEST", Message: "test error"},
			"test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("SecurityError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		setup        func() string
		expectedMode os.FileMode
		wantErr      bool
	}{
		{
			"valid file permissions",
			func() string {
				filePath := filepath.Join(tmpDir, "valid.txt")
				file, _ := os.Create(filePath)
				file.Close()
				_ = os.Chmod(filePath, 0644)
				return filePath
			},
			0644,
			false,
		},
		{
			"world writable file",
			func() string {
				filePath := filepath.Join(tmpDir, "world_writable.txt")
				file, _ := os.Create(filePath)
				file.Close()
				_ = os.Chmod(filePath, 0666)
				return filePath
			},
			0644,
			true,
		},
		{
			"group writable but not readable",
			func() string {
				filePath := filepath.Join(tmpDir, "group_writable.txt")
				file, _ := os.Create(filePath)
				file.Close()
				_ = os.Chmod(filePath, 0620)
				return filePath
			},
			0644,
			true,
		},
		{
			"nonexistent file",
			func() string {
				return filepath.Join(tmpDir, "nonexistent.txt")
			},
			0644,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup()
			err := ValidateFilePermissions(filePath, tt.expectedMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetSecureFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Create a file
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Set secure permissions
	err = SetSecureFilePermissions(filePath, 0644)
	if err != nil {
		t.Errorf("SetSecureFilePermissions() error = %v", err)
	}

	// Check that permissions are secure
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	mode := info.Mode()
	if mode&0002 != 0 {
		t.Error("File is world writable after SetSecureFilePermissions")
	}
}

func TestCreateSecureFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "secure.txt")

	// Create secure file
	file, err := CreateSecureFile(filePath, 0644)
	if err != nil {
		t.Errorf("CreateSecureFile() error = %v", err)
	}
	file.Close()

	// Check that file exists and has secure permissions
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	mode := info.Mode()
	if mode&0002 != 0 {
		t.Error("Created file is world writable")
	}
}

func TestCreateSecureDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "secure_dir")

	// Create secure directory
	err := CreateSecureDirectory(dirPath, 0755)
	if err != nil {
		t.Errorf("CreateSecureDirectory() error = %v", err)
	}

	// Check that directory exists and has secure permissions
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}

	mode := info.Mode()
	if mode&0002 != 0 {
		t.Error("Created directory is world writable")
	}
}
