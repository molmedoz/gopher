package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// MockHTTPClient implements HTTPClient interface for testing
type MockHTTPClient struct {
	responses map[string]*http.Response
}

func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		responses: make(map[string]*http.Response),
	}
}

func (m *MockHTTPClient) AddResponse(url string, statusCode int, headers map[string]string, body string) {
	response := &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	// Set headers
	for key, value := range headers {
		response.Header.Set(key, value)
	}

	m.responses[url] = response
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}
	return nil, fmt.Errorf("no mock response for URL: %s", url)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Get(req.URL.String())
}

func (m *MockHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Do(req)
}

// TestDownloader creates a downloader with mocked dependencies for testing
func TestDownloader(t *testing.T, baseURL string) (*Downloader, *MockHTTPClient) {
	mockClient := NewMockHTTPClient()
	downloader := &Downloader{
		client:  &http.Client{Transport: mockClient},
		baseURL: baseURL,
	}
	return downloader, mockClient
}

// SetupMockDownloadScenario sets up a complete mock download scenario
func SetupMockDownloadScenario(t *testing.T, version, filename string) (*Downloader, *MockHTTPClient, string) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Mock downloads page
			html := fmt.Sprintf(`
			<html>
			<body>
			<table>
			<tr>
				<td><a class="download" href="/dl/%s">%s</a></td>
				<td>62.0MB</td>
				<td><tt>5633d479dfae75ba7a78914ee380fa202bd6126e7c6b7c22e3ebc9e1a6ddc871</tt></td>
			</tr>
			</table>
			</body>
			</html>
			`, filename, filename)
			w.Write([]byte(html))
		} else if r.URL.Path == "/dl/"+filename {
			// Mock file download
			w.Header().Set("Content-Length", "65011712") // 62MB
			w.Write([]byte("mock file content"))
		}
	}))

	downloader := New(server.URL)
	return downloader, nil, server.URL
}

// AssertFileContent checks if a file contains expected content
func AssertFileContent(t *testing.T, filePath, expectedContent string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filePath, err)
	}
	if string(content) != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, string(content))
	}
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", filePath)
	}
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, filePath string) {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("Expected file %s to not exist", filePath)
	}
}

// CreateTestFile creates a test file with specified content
func CreateTestFile(t *testing.T, filePath, content string) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file %s: %v", filePath, err)
	}
}
