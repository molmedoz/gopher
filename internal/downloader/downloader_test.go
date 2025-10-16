package downloader

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestParseVersionNumber(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"1", 1},
		{"21", 21},
		{"25rc1", 25},
		{"0beta2", 0},
		{"", 0},
	}

	for _, c := range cases {
		if got := parseVersionNumber(c.in); got != c.want {
			t.Fatalf("parseVersionNumber(%q)=%d, want %d", c.in, got, c.want)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int // -1 if a<b, 0 if eq, 1 if a>b
	}{
		{"go1.25.1", "go1.25.0", 1},
		{"go1.25.1", "go1.24.9", 1},
		{"go1.25rc2", "go1.25rc1", 1},
		{"go1.25", "go1.25rc3", 1}, // stable > prerelease
		{"go1.21.3", "go1.21.3", 0},
		{"go1.19.9", "go1.20.0", -1},
	}

	for _, c := range cases {
		if got := compareVersions(c.a, c.b); got != c.want {
			t.Fatalf("compareVersions(%q,%q)=%d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestParseVersionParts(t *testing.T) {
	p := parseVersionParts("1.25.3rc2")
	if p.major != 1 || p.minor != 25 || p.patch != 3 {
		t.Fatalf("unexpected parts: %+v", p)
	}
	if p.prerelease != "rc2" {
		t.Fatalf("unexpected prerelease: %q", p.prerelease)
	}
	p2 := parseVersionParts("1.20")
	if p2.patch != 0 || p2.prerelease != "" {
		t.Fatalf("unexpected parts: %+v", p2)
	}
}

func TestParsePrereleaseParts(t *testing.T) {
	parts := parsePrereleaseParts("rc1.alpha2")
	if len(parts) != 2 || parts[0] != "rc1" || parts[1] != "alpha2" {
		t.Fatalf("bad prerelease parts: %v", parts)
	}
}

func TestComparePrerelease(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"", "rc1", 1},
		{"rc2", "rc1", 1},
		{"beta1", "alpha1", 1},
		{"rc1", "rc1", 0},
	}

	for _, c := range cases {
		if got := comparePrerelease(c.a, c.b); got != c.want {
			t.Fatalf("comparePrerelease(%q,%q)=%d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestComparePrereleasePart(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"rc1", "rc2", -1},
		{"rc2", "rc1", 1},
		{"rc10", "rc2", -1},
		{"alpha1", "alpha1", 0},
		{"", "rc1", -1},
		{"rc1", "", 1},
		{"beta", "beta", 0},
	}
	for _, c := range cases {
		if got := comparePrereleasePart(c.a, c.b); got != c.want {
			t.Fatalf("comparePrereleasePart(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestIsCompatibleFile(t *testing.T) {
	d := New("https://go.dev/dl/")
	// matching
	f := GoFile{OS: runtime.GOOS, Arch: runtime.GOARCH, Kind: "archive"}
	if !d.isCompatibleFile(f) {
		t.Fatalf("expected compatible for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	// mismatch OS
	f2 := GoFile{OS: "zzz", Arch: runtime.GOARCH, Kind: "archive"}
	if d.isCompatibleFile(f2) {
		t.Fatalf("expected not compatible for mismatched OS")
	}
	// mismatch kind
	f3 := GoFile{OS: runtime.GOOS, Arch: runtime.GOARCH, Kind: "source"}
	if d.isCompatibleFile(f3) {
		t.Fatalf("expected not compatible for non-binary kind")
	}
}

func TestNew(t *testing.T) {
	d := New("https://go.dev/dl/")
	if d.baseURL != "https://go.dev/dl" {
		t.Errorf("Expected baseURL 'https://go.dev/dl', got '%s'", d.baseURL)
	}
	if d.client == nil {
		t.Error("Expected client to be initialized")
	}
	if d.client.Timeout != 30*time.Minute {
		t.Errorf("Expected timeout 30m, got %v", d.client.Timeout)
	}
}

func TestWithClient(t *testing.T) {
	customClient := &http.Client{Timeout: 5 * time.Minute}
	d := WithClient("https://example.com", customClient)
	if d.baseURL != "https://example.com" {
		t.Errorf("Expected baseURL 'https://example.com', got '%s'", d.baseURL)
	}
	if d.client != customClient {
		t.Error("Expected client to be the custom client")
	}
}

func TestWithClientNil(t *testing.T) {
	d := WithClient("https://example.com", nil)
	if d.baseURL != "https://example.com" {
		t.Errorf("Expected baseURL 'https://example.com', got '%s'", d.baseURL)
	}
	if d.client == nil {
		t.Error("Expected client to be initialized")
	}
	if d.client.Timeout != 30*time.Minute {
		t.Errorf("Expected timeout 30m, got %v", d.client.Timeout)
	}
}

func TestGetFilename(t *testing.T) {
	d := New("https://go.dev/dl/")

	// Test with current platform
	result := d.getFilename("1.21.0")
	if result == "" {
		t.Error("Expected non-empty filename")
	}

	// Test that it contains expected elements
	if !strings.Contains(result, "go1.21.0") {
		t.Errorf("Expected filename to contain 'go1.21.0', got '%s'", result)
	}
	if !strings.Contains(result, runtime.GOOS) {
		t.Errorf("Expected filename to contain OS '%s', got '%s'", runtime.GOOS, result)
	}
	if !strings.Contains(result, runtime.GOARCH) {
		t.Errorf("Expected filename to contain arch '%s', got '%s'", runtime.GOARCH, result)
	}
}

func TestDownloadInfo(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Mock downloads page with multiple platforms
			html := `
			<html>
			<body>
			<table>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.linux-amd64.tar.gz">go1.21.0.linux-amd64.tar.gz</a></td>
				<td>62.0MB</td>
				<td><tt>d0398903a16ba2232b389fb31032ddf57cac34efda306a0eebac34684945965e</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.linux-arm64.tar.gz">go1.21.0.linux-arm64.tar.gz</a></td>
				<td>60.5MB</td>
				<td><tt>818d46ede85682dd551ad378ef37a4d247006f12ec59b5b755601d2ce114369a</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.darwin-amd64.tar.gz">go1.21.0.darwin-amd64.tar.gz</a></td>
				<td>63.0MB</td>
				<td><tt>b314de9f704ab122c077d2ec8e67e3670affe8865479d1f01991e7ac55d65e70</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.darwin-arm64.tar.gz">go1.21.0.darwin-arm64.tar.gz</a></td>
				<td>62.0MB</td>
				<td><tt>5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.windows-amd64.zip">go1.21.0.windows-amd64.zip</a></td>
				<td>68.5MB</td>
				<td><tt>4a0f5b2f2e0d8b6f3c964e97b0e6e6d3e4b5f8e2a0f5b2f2e0d8b6f3c964e97b</tt></td>
			</tr>
			</table>
			</body>
			</html>
			`
			w.Write([]byte(html))
		} else if r.URL.Path == "/go1.21.0.darwin-arm64.tar.gz" {
			// Mock file download
			content := "mock file content"
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(content))
		}
	}))
	defer server.Close()

	d := New(server.URL)

	// Test GetDownloadInfo
	info, err := d.GetDownloadInfo("1.21.0")
	if err != nil {
		t.Fatalf("GetDownloadInfo failed: %v", err)
	}

	// Check that we got info for the current platform
	expectedFilename := fmt.Sprintf("go1.21.0.%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		expectedFilename = fmt.Sprintf("go1.21.0.%s-%s.zip", runtime.GOOS, runtime.GOARCH)
	}

	if info.Filename != expectedFilename {
		t.Errorf("Expected filename '%s', got '%s'", expectedFilename, info.Filename)
	}

	// Verify size is parsed correctly (any positive number is fine for mock)
	if info.Size <= 0 {
		t.Errorf("Expected positive size, got %d", info.Size)
	}

	// Verify SHA256 is present
	if info.SHA256 == "" {
		t.Errorf("Expected SHA256 to be present")
	}
}

func TestDownload(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Mock downloads page with multiple platforms
			html := `
			<html>
			<body>
			<table>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.linux-amd64.tar.gz">go1.21.0.linux-amd64.tar.gz</a></td>
				<td>62.0MB</td>
				<td><tt>d0398903a16ba2232b389fb31032ddf57cac34efda306a0eebac34684945965e</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.linux-arm64.tar.gz">go1.21.0.linux-arm64.tar.gz</a></td>
				<td>60.5MB</td>
				<td><tt>818d46ede85682dd551ad378ef37a4d247006f12ec59b5b755601d2ce114369a</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.darwin-amd64.tar.gz">go1.21.0.darwin-amd64.tar.gz</a></td>
				<td>63.0MB</td>
				<td><tt>b314de9f704ab122c077d2ec8e67e3670affe8865479d1f01991e7ac55d65e70</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.darwin-arm64.tar.gz">go1.21.0.darwin-arm64.tar.gz</a></td>
				<td>62.0MB</td>
				<td><tt>5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871</tt></td>
			</tr>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.windows-amd64.zip">go1.21.0.windows-amd64.zip</a></td>
				<td>68.5MB</td>
				<td><tt>4a0f5b2f2e0d8b6f3c964e97b0e6e6d3e4b5f8e2a0f5b2f2e0d8b6f3c964e97b</tt></td>
			</tr>
			</table>
			</body>
			</html>
			`
			w.Write([]byte(html))
		} else {
			// Handle any platform-specific file download
			// Extract filename from path (e.g., "/go1.21.0.linux-amd64.tar.gz")
			content := "mock file content"
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(content))
		}
	}))
	defer server.Close()

	d := New(server.URL)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Test Download
	filePath, err := d.Download("1.21.0", tmpDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Downloaded file does not exist")
	}

	// Check file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if string(content) != "mock file content" {
		t.Errorf("Expected content 'mock file content', got '%s'", string(content))
	}
}

func TestDownloadExistingFile(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Mock downloads page
			html := `
			<html>
			<body>
			<table>
			<tr>
				<td><a class="download" href="/dl/go1.21.0.darwin-arm64.tar.gz">go1.21.0.darwin-arm64.tar.gz</a></td>
				<td>62.0MB</td>
				<td><tt>5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871</tt></td>
			</tr>
			</table>
			</body>
			</html>
			`
			w.Write([]byte(html))
		}
	}))
	defer server.Close()

	d := New(server.URL)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Create existing file with correct hash
	filePath := filepath.Join(tmpDir, "go1.21.0.darwin-arm64.tar.gz")
	err := os.WriteFile(filePath, []byte("mock file content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// Test Download with existing file
	resultPath, err := d.Download("1.21.0", tmpDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	if resultPath != filePath {
		t.Errorf("Expected path %s, got %s", filePath, resultPath)
	}

	// Check that file content is unchanged
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "mock file content" {
		t.Errorf("Expected content 'mock file content', got '%s'", string(content))
	}
}

func TestIsValidFile(t *testing.T) {
	d := New("https://go.dev/dl/")

	// Create temporary file
	tmpFile := t.TempDir() + "/test.txt"
	err := os.WriteFile(tmpFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with correct hash (SHA256 of "test content")
	hash := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"
	if !d.isValidFile(tmpFile, hash) {
		t.Error("Expected file to be valid with correct hash")
	}

	// Test with incorrect hash
	if d.isValidFile(tmpFile, "wronghash") {
		t.Error("Expected file to be invalid with incorrect hash")
	}

	// Test with non-existent file
	if d.isValidFile("/nonexistent/file", hash) {
		t.Error("Expected non-existent file to be invalid")
	}
}

func TestCleanup(t *testing.T) {
	d := New("https://go.dev/dl/")

	// Create temporary file
	tmpFile := t.TempDir() + "/test.txt"
	err := os.WriteFile(tmpFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test cleanup
	err = d.Cleanup(tmpFile)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Check file is deleted
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("Expected file to be deleted after cleanup")
	}
}

func TestListAvailableVersions(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<html>
		<body>
		<table>
		<tr>
			<td><a class="download" href="/dl/go1.21.0.darwin-arm64.tar.gz">go1.21.0.darwin-arm64.tar.gz</a></td>
			<td>62.0 MB</td>
			<td><tt>5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871</tt></td>
		</tr>
		<tr>
			<td><a class="download" href="/dl/go1.20.0.darwin-arm64.tar.gz">go1.20.0.darwin-arm64.tar.gz</a></td>
			<td>62.0 MB</td>
			<td><tt>def456ghi789</tt></td>
		</tr>
		</table>
		</body>
		</html>
		`
		w.Write([]byte(html))
	}))
	defer server.Close()

	d := New(server.URL)

	// Test ListAvailableVersions
	versions, err := d.ListAvailableVersions()
	if err != nil {
		t.Fatalf("ListAvailableVersions failed: %v", err)
	}

	if len(versions) == 0 {
		t.Error("Expected at least one version")
	}

	// Check that we have the expected versions
	found21 := false
	found20 := false
	for _, v := range versions {
		if v.Version == "go1.21.0" {
			found21 = true
		}
		if v.Version == "go1.20.0" {
			found20 = true
		}
	}

	if !found21 {
		t.Error("Expected to find go1.21.0")
	}
	if !found20 {
		t.Error("Expected to find go1.20.0")
	}
}

func TestParseFileInfoFromHTML(t *testing.T) {
	d := New("https://go.dev/dl/")

	html := `
	<html>
	<body>
	<table>
	<tr>
		<td><a class="download" href="/dl/go1.21.0.darwin-arm64.tar.gz">go1.21.0.darwin-arm64.tar.gz</a></td>
		<td>62.0MB</td>
		<td><tt>5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871</tt></td>
	</tr>
	</table>
	</body>
	</html>
	`

	sha256, size, err := d.parseFileInfoFromHTML(html, "go1.21.0.darwin-arm64.tar.gz")
	if err != nil {
		t.Fatalf("parseFileInfoFromHTML failed: %v", err)
	}

	if sha256 != "5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871" {
		t.Errorf("Expected SHA256 '5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871', got '%s'", sha256)
	}
	if size != 65011712 { // 62.0 MB
		t.Errorf("Expected size 65011712, got %d", size)
	}
}

func TestParseFileInfoFromHTMLNotFound(t *testing.T) {
	d := New("https://go.dev/dl/")

	html := `<html><body><table></table></body></html>`

	_, _, err := d.parseFileInfoFromHTML(html, "nonexistent.tar.gz")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestGetFileSize(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Set("Content-Length", "52428800") // 50MB
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	d := New(server.URL)

	// Test getFileSize
	size, err := d.getFileSize("go1.21.0.linux-amd64.tar.gz")
	if err != nil {
		t.Fatalf("getFileSize failed: %v", err)
	}

	if size != 52428800 {
		t.Errorf("Expected size 52428800, got %d", size)
	}
}

func TestGetFileSizeError(t *testing.T) {
	// Create a test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	d := New(server.URL)

	// Test getFileSize with error
	_, err := d.getFileSize("nonexistent.tar.gz")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestIsVersionString(t *testing.T) {
	d := New("https://go.dev/dl/")

	tests := []struct {
		input    string
		expected bool
	}{
		{"go1.21.0", true},
		{"go1.20.0", true},
		{"go1.21.0beta1", true},
		{"go1.21.0rc1", true},
		{"1.21.0", false},
		{"go", false},
		{"", false},
		{"not-a-version", false},
	}

	for _, test := range tests {
		result := d.isVersionString(test.input)
		if result != test.expected {
			t.Errorf("isVersionString(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsValidVersionPart(t *testing.T) {
	d := New("https://go.dev/dl/")

	tests := []struct {
		input    string
		expected bool
	}{
		{"1", true},
		{"21", true},
		{"0", true},
		{"beta", true},
		{"rc1", true},
		{"alpha2", true},
		{"pre1", true},
		{"", false},
		{"invalid", false},
		{"1beta", true},
	}

	for _, test := range tests {
		result := d.isValidVersionPart(test.input)
		if result != test.expected {
			t.Errorf("isValidVersionPart(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	d := New("https://go.dev/dl/")

	tests := []struct {
		input    string
		expected bool
	}{
		{"1", true},
		{"21", true},
		{"0", true},
		{"123", true},
		{"", false},
		{"1a", false},
		{"a1", false},
		{"beta", false},
	}

	for _, test := range tests {
		result := d.isNumeric(test.input)
		if result != test.expected {
			t.Errorf("isNumeric(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestNormalizeVersionString(t *testing.T) {
	d := New("https://go.dev/dl/")

	tests := []struct {
		input    string
		expected string
	}{
		{"go1.21.0", "go1.21.0"},
		{"go1.21.0beta1", "go1.21.0beta1"},
		{"go1.21.0rc1", "go1.21.0rc1"},
		{"1.21.0", ""},
		{"go", "go"},
		{"", ""},
		{"go1.21.0 ", "go1.21.0"},
		{" go1.21.0", "go1.21.0"},
		{"<go1.21.0>", "go1.21.0"},
		{"\"go1.21.0\"", "go1.21.0"},
	}

	for _, test := range tests {
		result := d.normalizeVersionString(test.input)
		if result != test.expected {
			t.Errorf("normalizeVersionString(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestIsVersionInMap(t *testing.T) {
	d := New("https://go.dev/dl/")

	versionMap := map[string]VersionInfo{
		"go1.21.0": {Version: "go1.21.0"},
		"go1.20.0": {Version: "go1.20.0"},
	}

	tests := []struct {
		version  string
		expected bool
	}{
		{"go1.21.0", true},
		{"go1.20.0", true},
		{"go1.19.0", false},
		{"", false},
	}

	for _, test := range tests {
		result := d.isVersionInMap(versionMap, test.version)
		if result != test.expected {
			t.Errorf("isVersionInMap(%s) = %v, expected %v", test.version, result, test.expected)
		}
	}
}
