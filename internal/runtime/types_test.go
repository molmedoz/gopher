package runtime

import (
	"runtime"
	"testing"
	"time"
)

func TestVersionString(t *testing.T) {
	v := &Version{
		Version:     "1.21.0",
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstalledAt: time.Now(),
		IsActive:    false,
		IsSystem:    false,
	}

	expected := "1.21.0 (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
	if v.String() != expected {
		t.Errorf("Version.String() = %s, want %s", v.String(), expected)
	}
}

func TestVersionString_SystemVersion(t *testing.T) {
	v := &Version{
		Version:     "1.21.0",
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstalledAt: time.Now(),
		IsActive:    false,
		IsSystem:    true,
	}

	expected := "1.21.0 (" + runtime.GOOS + "/" + runtime.GOARCH + ") [system]"
	if v.String() != expected {
		t.Errorf("Version.String() = %s, want %s", v.String(), expected)
	}
}

func TestVersionFullString(t *testing.T) {
	v := &Version{
		Version:     "1.21.0",
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstalledAt: time.Now(),
		IsActive:    false,
		IsSystem:    false,
	}

	expected := "1.21.0 (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
	if v.FullString() != expected {
		t.Errorf("Version.FullString() = %s, want %s", v.FullString(), expected)
	}
}

func TestVersionFullString_SystemVersion(t *testing.T) {
	v := &Version{
		Version:     "1.21.0",
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		InstalledAt: time.Now(),
		IsActive:    false,
		IsSystem:    true,
	}

	expected := "1.21.0 (" + runtime.GOOS + "/" + runtime.GOARCH + ") [system]"
	if v.FullString() != expected {
		t.Errorf("Version.FullString() = %s, want %s", v.FullString(), expected)
	}
}

func TestVersionIsCompatible(t *testing.T) {
	tests := []struct {
		name     string
		version  *Version
		expected bool
	}{
		{
			name: "compatible version",
			version: &Version{
				OS:       runtime.GOOS,
				Arch:     runtime.GOARCH,
				IsSystem: false,
			},
			expected: true,
		},
		{
			name: "incompatible OS",
			version: &Version{
				OS:       "windows",
				Arch:     runtime.GOARCH,
				IsSystem: false,
			},
			expected: false,
		},
		{
			name: "incompatible arch",
			version: &Version{
				OS: runtime.GOOS,
				Arch: func() string {
					// Use an architecture that's different from current
					if runtime.GOARCH == "amd64" {
						return "arm64"
					}
					return "amd64"
				}(),
				IsSystem: false,
			},
			expected: false,
		},
		{
			name: "system version",
			version: &Version{
				OS:       "any",
				Arch:     "any",
				IsSystem: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.IsCompatible()
			if result != tt.expected {
				t.Errorf("Version.IsCompatible() = %v, want %v", result, tt.expected)
			}
		})
	}
}
