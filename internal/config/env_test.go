package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetGOPATH_Shared_DefaultsToHomeGo(t *testing.T) {
	t.Setenv("GOPATH", "")
	cfg := &Config{InstallDir: "/tmp/install", DownloadDir: "/tmp/dl", GOPATHMode: "shared"}
	gp := cfg.GetGOPATH("go1.25.1")
	home := getUserHomeDir()
	want := filepath.Join(home, "go")
	if gp != want {
		t.Fatalf("GetGOPATH(shared)=%q want %q", gp, want)
	}
}

func TestGetGOPATH_Shared_RespectsEnv(t *testing.T) {
	t.Setenv("GOPATH", "/tmp/mygopath")
	cfg := &Config{InstallDir: "/tmp/install", DownloadDir: "/tmp/dl", GOPATHMode: "shared"}
	if got := cfg.GetGOPATH("go1.25.1"); got != "/tmp/mygopath" {
		t.Fatalf("GetGOPATH(shared with env)=%q", got)
	}
}

func TestGetGOPATH_VersionSpecific(t *testing.T) {
	t.Setenv("GOPATH", "")
	tmp := t.TempDir()
	cfg := &Config{InstallDir: tmp, DownloadDir: filepath.Join(tmp, "dl"), GOPATHMode: "version-specific"}
	got := cfg.GetGOPATH("go1.21.0")
	want := filepath.Join(tmp, "go1.21.0", "gopath")
	if got != want {
		t.Fatalf("GetGOPATH(version-specific)=%q want %q", got, want)
	}
}

func TestGetEnvironmentVariables_ComposesPATH(t *testing.T) {
	tmp := t.TempDir()
	cfg := &Config{InstallDir: tmp, DownloadDir: filepath.Join(tmp, "dl"), GOPATHMode: "version-specific", SetEnvironment: true}
	ver := "go1.20.5"
	env := cfg.GetEnvironmentVariables(ver)
	if env == nil {
		t.Fatalf("expected env map")
	}
	goroot := filepath.Join(tmp, ver)
	gopath := filepath.Join(tmp, ver, "gopath")
	if env["GOROOT"] != goroot {
		t.Fatalf("GOROOT=%q want %q", env["GOROOT"], goroot)
	}
	if env["GOPATH"] != gopath {
		t.Fatalf("GOPATH=%q want %q", env["GOPATH"], gopath)
	}
	bin := filepath.Join(goroot, "bin")
	if !strings.HasPrefix(env["PATH"], bin+string(os.PathListSeparator)) && env["PATH"] != bin {
		t.Fatalf("PATH does not start with goroot bin: %q", env["PATH"])
	}
}

func TestGetGOROOT_PathIsUnderInstallDir(t *testing.T) {
	tmp := t.TempDir()
	cfg := &Config{InstallDir: tmp, DownloadDir: filepath.Join(tmp, "dl")}
	got := cfg.GetGOROOT("go1.18.10")
	if !strings.HasPrefix(got, tmp) || !strings.Contains(got, "go1.18.10") {
		t.Fatalf("GOROOT=%q not under install dir %q", got, tmp)
	}
}

func TestGetGOPATH_CustomMode(t *testing.T) {
	tmp := t.TempDir()
	cfg := &Config{InstallDir: tmp, DownloadDir: filepath.Join(tmp, "dl"), GOPATHMode: "custom", CustomGOPATH: filepath.Join(tmp, "ws")}
	if got := cfg.GetGOPATH("go1.2.3"); got != cfg.CustomGOPATH {
		t.Fatalf("GetGOPATH(custom)=%q want %q", got, cfg.CustomGOPATH)
	}
}

func TestGetEnvironmentVariables_Disabled(t *testing.T) {
	tmp := t.TempDir()
	cfg := &Config{InstallDir: tmp, DownloadDir: filepath.Join(tmp, "dl"), SetEnvironment: false}
	if env := cfg.GetEnvironmentVariables("go1.2.3"); env != nil {
		t.Fatalf("expected nil env when SetEnvironment=false, got %v", env)
	}
}

func TestEnsureDirectories_Creates(t *testing.T) {
	tmp := t.TempDir()
	cfg := &Config{InstallDir: filepath.Join(tmp, "inst"), DownloadDir: filepath.Join(tmp, "dl")}
	if err := cfg.EnsureDirectories(); err != nil {
		t.Fatalf("EnsureDirectories error: %v", err)
	}
	for _, p := range []string{cfg.InstallDir, cfg.DownloadDir} {
		if st, err := os.Stat(p); err != nil || !st.IsDir() {
			t.Fatalf("dir %s not created", p)
		}
	}
}

// Sanity: GetConfigPath returns absolute path appropriate for OS
func TestGetConfigPath_IsAbsolute(t *testing.T) {
	p := GetConfigPath()
	if !filepath.IsAbs(p) {
		t.Fatalf("GetConfigPath not absolute: %q (%s)", p, runtime.GOOS)
	}
}
