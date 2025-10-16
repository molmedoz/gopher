package downloader

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadWithMock(t *testing.T) {
	// Create test server
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
			_, _ = w.Write([]byte(html))
		} else if r.URL.Path == "/go1.21.0.darwin-arm64.tar.gz" {
			// Mock file download
			content := "mock file content"
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(content))
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
	AssertFileExists(t, filePath)

	// Check file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	t.Logf("File content length: %d", len(content))
	t.Logf("File content: %q", string(content))
	AssertFileContent(t, filePath, "mock file content")
}

func TestGetDownloadInfo(t *testing.T) {
	// Create test server
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
			_, _ = w.Write([]byte(html))
		} else if r.URL.Path == "/go1.21.0.darwin-arm64.tar.gz" {
			// Mock file download
			content := "mock file content"
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(content))
		}
	}))
	defer server.Close()

	d := New(server.URL)

	// Test GetDownloadInfo
	info, err := d.GetDownloadInfo("1.21.0")
	if err != nil {
		t.Fatalf("GetDownloadInfo failed: %v", err)
	}

	if info.Filename != "go1.21.0.darwin-arm64.tar.gz" {
		t.Errorf("Expected filename 'go1.21.0.darwin-arm64.tar.gz', got '%s'", info.Filename)
	}
	if info.Size != 65011712 {
		t.Errorf("Expected size 65011712, got %d", info.Size)
	}
	if info.SHA256 != "5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871" {
		t.Errorf("Expected SHA256 '5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871', got '%s'", info.SHA256)
	}
}
