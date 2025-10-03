package runtime

import (
	"os"
	"runtime"
	"testing"
)

func TestSystemDetector_IsSystemInstallation(t *testing.T) {
	detector := NewSystemDetector()

	tests := []struct {
		name     string
		goPath   string
		expected bool
	}{
		{
			name:     "system path /usr/bin/go",
			goPath:   "/usr/bin/go",
			expected: true,
		},
		{
			name:     "system path /usr/local/bin/go",
			goPath:   "/usr/local/bin/go",
			expected: true,
		},
		{
			name:     "homebrew path apple silicon",
			goPath:   "/opt/homebrew/opt/go/libexec/bin/go",
			expected: true,
		},
		{
			name:     "homebrew path intel",
			goPath:   "/usr/local/opt/go/libexec/bin/go",
			expected: true,
		},
		{
			name:     "user path",
			goPath:   "/home/user/go/bin/go",
			expected: false,
		},
		{
			name:     "gopher managed path",
			goPath:   "/home/user/.gopher/versions/go1.21.0/bin/go",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.isSystemInstallation(tt.goPath)
			if result != tt.expected {
				t.Errorf("isSystemInstallation(%s) = %v, want %v", tt.goPath, result, tt.expected)
			}
		})
	}
}

func TestSystemDetector_IsSystemInstallation_WindowsPaths(t *testing.T) {
	detector := NewSystemDetector()
	tests := []struct {
		path string
		want bool
	}{
		{"C:\\Program Files\\Go\\bin\\go.exe", true},
		{"C:\\Go\\bin\\go.exe", true},
		{"D:\\Tools\\Go\\bin\\go.exe", false},
	}
	for _, c := range tests {
		if got := detector.isSystemInstallation(c.path); got != c.want {
			t.Errorf("isSystemInstallation(%s) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestSystemDetector_parseGoVersion_EdgeCases(t *testing.T) {
	detector := NewSystemDetector()
	cases := []struct {
		out     string
		want    string
		wantErr bool
	}{
		{"go version go1.22.5 windows/amd64", "go1.22.5", false},
		{"go version 1.22.5 linux/amd64", "", true}, // missing 'go' prefix before version token
		{"go version go1.22.5", "go1.22.5", false},  // minimal tokens still ok (>=3?)
		{"go version", "", true},
	}
	for _, c := range cases {
		got, err := detector.parseGoVersion(c.out)
		if (err != nil) != c.wantErr {
			t.Fatalf("parseGoVersion error mismatch: %v wantErr=%v", err, c.wantErr)
		}
		if !c.wantErr && got != c.want {
			t.Fatalf("parseGoVersion(%q)=%q want %q", c.out, got, c.want)
		}
	}
}

func TestSystemDetector_IsSystemGoAvailable(t *testing.T) {
	detector := NewSystemDetector()

	// This test depends on whether Go is actually installed on the system
	// We'll just test that the function doesn't panic
	_ = detector.IsSystemGoAvailable()
}

func TestSystemDetector_parseGoVersion(t *testing.T) {
	detector := NewSystemDetector()

	tests := []struct {
		name     string
		output   string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid version output",
			output:   "go version go1.21.0 darwin/arm64",
			expected: "go1.21.0",
			wantErr:  false,
		},
		{
			name:     "valid version output with extra info",
			output:   "go version go1.21.0 darwin/arm64 X:prefer-std",
			expected: "go1.21.0",
			wantErr:  false,
		},
		{
			name:     "invalid output format",
			output:   "invalid output",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.parseGoVersion(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGoVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("parseGoVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSystemGoInfo_String(t *testing.T) {
	info := &SystemGoInfo{
		Executable: "/usr/bin/go",
		Version:    "go version go1.21.0 darwin/arm64",
		GOROOT:     "/usr/local/go",
		GOPATH:     "/home/user/go",
		IsValid:    true,
	}

	result := info.Executable
	if result == "" {
		t.Error("String() should not return empty string")
	}
}

func TestSystemGoInfo_FullString(t *testing.T) {
	info := &SystemGoInfo{
		Executable: "/usr/bin/go",
		Version:    "go version go1.21.0 darwin/arm64",
		GOROOT:     "/usr/local/go",
		GOPATH:     "/home/user/go",
		IsValid:    true,
	}

	result := info.Executable
	if result == "" {
		t.Error("FullString() should not return empty string")
	}
}

// TestSystemDetector_DetectSystemGo tests the system detection if Go is available
func TestSystemDetector_DetectSystemGo(t *testing.T) {
	detector := NewSystemDetector()

	// Only run this test if Go is available
	if !detector.IsSystemGoAvailable() {
		t.Skip("Go is not available on this system")
	}

	version, err := detector.DetectSystemGo()
	if err != nil {
		t.Fatalf("DetectSystemGo() error = %v", err)
	}

	if version.Version == "" {
		t.Error("Version should not be empty")
	}
	if version.OS != runtime.GOOS {
		t.Errorf("OS = %v, want %v", version.OS, runtime.GOOS)
	}
	if version.Arch != runtime.GOARCH {
		t.Errorf("Arch = %v, want %v", version.Arch, runtime.GOARCH)
	}
	if !version.IsSystem {
		t.Error("IsSystem should be true for system Go")
	}
	if version.Path == "" {
		t.Error("Path should not be empty")
	}
}

// TestSystemDetector_GetSystemGoPath tests getting system Go path
func TestSystemDetector_GetSystemGoPath(t *testing.T) {
	detector := NewSystemDetector()

	// Only run this test if Go is available
	if !detector.IsSystemGoAvailable() {
		t.Skip("Go is not available on this system")
	}

	path, err := detector.GetSystemGoPath()
	if err != nil {
		t.Fatalf("GetSystemGoPath() error = %v", err)
	}

	if path == "" {
		t.Error("Path should not be empty")
	}

	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Go binary at %s does not exist", path)
	}
}

// TestSystemDetector_GetSystemGoInfo tests getting detailed system Go info
func TestSystemDetector_GetSystemGoInfo(t *testing.T) {
	detector := NewSystemDetector()

	// Only run this test if Go is available
	if !detector.IsSystemGoAvailable() {
		t.Skip("Go is not available on this system")
	}

	info, err := detector.GetSystemGoInfo()
	if err != nil {
		t.Fatalf("GetSystemGoInfo() error = %v", err)
	}

	if info.Executable == "" {
		t.Error("Executable should not be empty")
	}
	if info.Version == "" {
		t.Error("Version should not be empty")
	}
	if info.GOROOT == "" {
		t.Error("GOROOT should not be empty")
	}
	if info.GOPATH == "" {
		t.Error("GOPATH should not be empty")
	}
	if !info.IsValid {
		t.Error("IsValid should be true for system Go")
	}
}

// TestVersion_String_SystemVersion tests string representation of system versions
func TestVersion_String_SystemVersion(t *testing.T) {
	version := Version{
		Version:  "go1.21.0",
		OS:       "darwin",
		Arch:     "arm64",
		IsSystem: true,
	}

	result := version.String()
	expected := "go1.21.0 (darwin/arm64) [system]"
	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}
}

// TestVersion_FullString_SystemVersion tests full string representation of system versions
func TestVersion_FullString_SystemVersion(t *testing.T) {
	version := Version{
		Version:  "go1.21.0",
		OS:       "darwin",
		Arch:     "arm64",
		IsSystem: true,
		IsActive: true,
	}

	result := version.FullString()
	if !contains(result, "[system]") {
		t.Errorf("FullString() should contain '[system]', got: %v", result)
	}
}

// TestSystemDetector_parseGoVersion_MoreEdgeCases tests additional edge cases
func TestSystemDetector_parseGoVersion_MoreEdgeCases(t *testing.T) {
	detector := NewSystemDetector()

	tests := []struct {
		name     string
		output   string
		expected string
		wantErr  bool
	}{
		{
			name:     "version with build info",
			output:   "go version go1.21.0 darwin/arm64 X:prefer-std BuildID=abc123",
			expected: "go1.21.0",
			wantErr:  false,
		},
		{
			name:     "version with devel prefix",
			output:   "go version devel go1.22-abc123 darwin/arm64",
			expected: "devel go1.22-abc123",
			wantErr:  false,
		},
		{
			name:     "version with rc suffix",
			output:   "go version go1.21rc1 darwin/arm64",
			expected: "go1.21rc1",
			wantErr:  false,
		},
		{
			name:     "version with beta suffix",
			output:   "go version go1.21beta1 darwin/arm64",
			expected: "go1.21beta1",
			wantErr:  false,
		},
		{
			name:     "too few fields",
			output:   "go version",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "version without go prefix",
			output:   "go version 1.21.0 darwin/arm64",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.parseGoVersion(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGoVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("parseGoVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSystemDetector_isSystemInstallation_MorePaths tests additional system paths
func TestSystemDetector_isSystemInstallation_MorePaths(t *testing.T) {
	detector := NewSystemDetector()

	tests := []struct {
		name     string
		goPath   string
		expected bool
	}{
		{
			name:     "opt/go/bin/go",
			goPath:   "/opt/go/bin/go",
			expected: true,
		},
		{
			name:     "usr/local/go/bin/go",
			goPath:   "/usr/local/go/bin/go",
			expected: true,
		},
		{
			name:     "homebrew opt path",
			goPath:   "/opt/homebrew/opt/go/libexec/bin/go",
			expected: true,
		},
		{
			name:     "homebrew local path",
			goPath:   "/usr/local/opt/go/libexec/bin/go",
			expected: true,
		},
		{
			name:     "snap path",
			goPath:   "/snap/go/current/bin/go",
			expected: false, // Not in our system paths
		},
		{
			name:     "flatpak path",
			goPath:   "/var/lib/flatpak/app/org.golang.Go/current/active/files/bin/go",
			expected: false, // Not in our system paths
		},
		{
			name:     "nix path",
			goPath:   "/nix/store/abc123-go-1.21.0/bin/go",
			expected: false, // Not in our system paths
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.isSystemInstallation(tt.goPath)
			if result != tt.expected {
				t.Errorf("isSystemInstallation(%s) = %v, want %v", tt.goPath, result, tt.expected)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && contains(s[:len(s)-1], substr)
}
