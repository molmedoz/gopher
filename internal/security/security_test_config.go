//go:build security

package security

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSecurityIntegration tests the integration of security features
func TestSecurityIntegration(t *testing.T) {
	t.Run("PathTraversalProtection", func(t *testing.T) {
		// Test various path traversal attempts
		maliciousPaths := []string{
			"../../../etc/passwd",
			"..\\..\\..\\windows\\system32",
			"....//....//....//etc/passwd",
			"..%2F..%2F..%2Fetc%2Fpasswd",
		}

		for _, path := range maliciousPaths {
			if err := ValidatePath(path); err == nil {
				t.Errorf("ValidatePath() should have rejected malicious path: %s", path)
			}
		}
	})

	t.Run("FilePermissionSecurity", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Test creating files with secure permissions
		filePath := filepath.Join(tmpDir, "secure.txt")
		file, err := CreateSecureFile(filePath, 0644)
		if err != nil {
			t.Fatalf("CreateSecureFile() error = %v", err)
		}
		file.Close()

		// Verify file permissions are secure
		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to stat file: %v", err)
		}

		mode := info.Mode()
		if mode&0002 != 0 {
			t.Error("File should not be world writable")
		}
	})

	t.Run("InputSanitization", func(t *testing.T) {
		// Test sanitization of various inputs
		inputs := []struct {
			input    string
			expected string
		}{
			{"normal_input", "normal_input"},
			{"../malicious", "malicious"},
			{"$HOME/secret", "HOME/secret"},
			{"`rm -rf /`", "rm -rf /"},
			{"input;rm -rf /", "inputrm -rf /"},
		}

		for _, test := range inputs {
			result := SanitizePath(test.input)
			if result != test.expected {
				t.Errorf("SanitizePath(%s) = %s, want %s", test.input, result, test.expected)
			}
		}
	})
}

// TestSecurityPerformance tests the performance of security functions
func TestSecurityPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("PathValidationPerformance", func(t *testing.T) {
		// Test performance with many path validations
		paths := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			paths[i] = "test/path/with/many/components"
		}

		for _, path := range paths {
			if err := ValidatePath(path); err != nil {
				t.Errorf("ValidatePath() error = %v", err)
			}
		}
	})
}

// TestSecurityEdgeCases tests edge cases in security functions
func TestSecurityEdgeCases(t *testing.T) {
	t.Run("EmptyAndNilInputs", func(t *testing.T) {
		// Test empty string
		if err := ValidatePath(""); err == nil {
			t.Error("ValidatePath() should reject empty string")
		}

		// Test very long path
		longPath := string(make([]byte, 10000))
		if err := ValidatePath(longPath); err != nil {
			t.Errorf("ValidatePath() should handle long paths: %v", err)
		}
	})

	t.Run("UnicodeAndSpecialCharacters", func(t *testing.T) {
		// Test Unicode characters
		unicodePath := "测试/路径/文件"
		if err := ValidatePath(unicodePath); err != nil {
			t.Errorf("ValidatePath() should handle Unicode: %v", err)
		}

		// Test special characters that might be valid
		specialPath := "file-with-dashes_and_underscores.txt"
		if err := ValidatePath(specialPath); err != nil {
			t.Errorf("ValidatePath() should handle special chars: %v", err)
		}
	})
}
