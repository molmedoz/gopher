package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/molmedoz/gopher/internal/progress"
)

// Downloader handles downloading Go versions
type Downloader struct {
	client  *http.Client
	baseURL string
}

// New creates a new downloader
func New(baseURL string) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 30 * time.Minute, // Long timeout for large downloads
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Follow redirects
				return nil
			},
		},
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// WithClient creates a downloader with a custom http.Client (for testing)
func WithClient(baseURL string, client *http.Client) *Downloader {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Minute}
	}
	return &Downloader{client: client, baseURL: strings.TrimSuffix(baseURL, "/")}
}

// DownloadInfo contains information about a download
type DownloadInfo struct {
	URL      string
	Filename string
	Size     int64
	SHA256   string
}

// GoRelease represents a Go release from the API
type GoRelease struct {
	Version string   `json:"version"`
	Stable  bool     `json:"stable"`
	Date    string   `json:"date"`
	Files   []GoFile `json:"files"`
}

// GoFile represents a file in a Go release
type GoFile struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Size     int64  `json:"size"`
	SHA256   string `json:"sha256"`
	Kind     string `json:"kind"`
}

// VersionInfo represents information about available Go versions
type VersionInfo struct {
	Version     string `json:"version"`
	Stable      bool   `json:"stable"`
	ReleaseDate string `json:"release_date"`
	Files       []File `json:"files"`
}

// File represents a downloadable file for a Go version
type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Size     int64  `json:"size"`
	SHA256   string `json:"sha256"`
}

// GetDownloadInfo returns download information for a Go version
func (d *Downloader) GetDownloadInfo(version string) (*DownloadInfo, error) {
	// Remove 'go' prefix if present
	version = strings.TrimPrefix(version, "go")

	// Determine filename based on OS and architecture
	filename := d.getFilename(version)

	// Construct download URL
	url := fmt.Sprintf("%s/%s", d.baseURL, filename)

	// Get file size and SHA256 from the checksum file
	size, sha256, err := d.getFileInfo(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &DownloadInfo{
		URL:      url,
		Filename: filename,
		Size:     size,
		SHA256:   sha256,
	}, nil
}

// Download downloads a Go version to the specified directory
func (d *Downloader) Download(version string, downloadDir string) (string, error) {
	info, err := d.GetDownloadInfo(version)
	if err != nil {
		return "", fmt.Errorf("failed to get download info: %w", err)
	}

	// Create download directory if it doesn't exist
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	// Construct local file path
	localPath := filepath.Join(downloadDir, info.Filename)

	// Check if file already exists and is valid
	if d.isValidFile(localPath, info.SHA256) {
		return localPath, nil
	}

	// Download the file
	if err := d.downloadFile(info.URL, localPath); err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	// Verify the downloaded file
	if !d.isValidFile(localPath, info.SHA256) {
		os.Remove(localPath) // Clean up invalid file
		return "", fmt.Errorf("downloaded file failed verification (checksum mismatch)")
	}

	return localPath, nil
}

// getFilename returns the appropriate filename for the current platform
func (d *Downloader) getFilename(version string) string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Handle special cases for architecture names
	switch arch {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	case "386":
		arch = "386"
	default:
		arch = "amd64" // Default fallback
	}

	// Handle special cases for OS names
	switch os {
	case "darwin":
		os = "darwin"
	case "linux":
		os = "linux"
	case "windows":
		os = "windows"
		// Don't add .msi to arch here, we'll add it in the filename construction
	default:
		os = "linux" // Default fallback
	}

	if os == "windows" {
		return fmt.Sprintf("go%s.%s-%s.zip", version, os, arch)
	}

	return fmt.Sprintf("go%s.%s-%s.tar.gz", version, os, arch)
}

// getFileInfo retrieves file size and SHA256 from the HTML page
func (d *Downloader) getFileInfo(version string) (int64, string, error) {
	// Download the main downloads page
	pageURL := d.baseURL + "/"

	resp, err := d.client.Get(pageURL)
	if err != nil {
		return 0, "", fmt.Errorf("failed to download downloads page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("failed to download downloads page: HTTP %d (check your internet connection)", resp.StatusCode)
	}

	// Read page content
	pageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("failed to read downloads page: %w", err)
	}

	// Parse HTML to find our file
	filename := d.getFilename(version)
	sha256, size, err := d.parseFileInfoFromHTML(string(pageData), filename)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse file info: %w", err)
	}

	return size, sha256, nil
}

// parseFileInfoFromHTML parses file info from the HTML downloads page
func (d *Downloader) parseFileInfoFromHTML(html, filename string) (string, int64, error) {
	// Look for the specific file in the HTML table
	// The pattern is: <a class="download" href="/dl/filename">filename</a>
	// followed by size and checksum in the same row

	// Find the row containing our filename
	pattern := fmt.Sprintf(`<a class="download" href="/dl/%s">%s</a>`, filename, filename)
	startIndex := strings.Index(html, pattern)
	if startIndex == -1 {
		return "", 0, fmt.Errorf("file not found in downloads page: %s", filename)
	}

	// Find the end of the table row
	rowEnd := strings.Index(html[startIndex:], "</tr>")
	if rowEnd == -1 {
		return "", 0, fmt.Errorf("malformed table row for file: %s", filename)
	}

	rowContent := html[startIndex : startIndex+rowEnd]

	// Extract size and checksum from the row
	// Look for size pattern: <td>XXMB</td> or <td>XXGB</td>
	sizePattern := `<td>(\d+(?:\.\d+)?)(MB|GB)</td>`
	sizeRegex := regexp.MustCompile(sizePattern)
	sizeMatch := sizeRegex.FindStringSubmatch(rowContent)
	if len(sizeMatch) < 3 {
		return "", 0, fmt.Errorf("size not found for file: %s", filename)
	}

	// Convert size to bytes
	sizeStr := sizeMatch[1]
	unit := sizeMatch[2]
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid size format for file %s: %s", filename, sizeStr)
	}

	switch unit {
	case "GB":
		size *= 1024 * 1024 * 1024
	case "MB":
		size *= 1024 * 1024
	}

	// Look for checksum pattern: <tt>checksum</tt>
	checksumPattern := `<tt>([a-f0-9]{64})</tt>`
	checksumRegex := regexp.MustCompile(checksumPattern)
	checksumMatch := checksumRegex.FindStringSubmatch(rowContent)
	if len(checksumMatch) < 2 {
		return "", 0, fmt.Errorf("checksum not found for file: %s", filename)
	}

	return checksumMatch[1], int64(size), nil
}

// getFileSize gets the size of a file by making a HEAD request
func (d *Downloader) getFileSize(filename string) (int64, error) {
	url := fmt.Sprintf("%s/%s", d.baseURL, filename)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get file size: HTTP %d (check your internet connection)", resp.StatusCode)
	}

	return resp.ContentLength, nil
}

// downloadFile downloads a file from URL to local path
func (d *Downloader) downloadFile(url, localPath string) error {
	// Create the file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Make the request
	resp, err := d.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to Umake request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: HTTP %d (check your internet connection)", resp.StatusCode)
	}

	// Get file size for progress tracking
	fileSize := resp.ContentLength
	if fileSize <= 0 {
		// If Content-Length is not available, we can't show progress
		fmt.Printf("Downloading %s...\n", filepath.Base(localPath))
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
		fmt.Println("✓ Download complete")
		return nil
	}

	// Create progress bar
	progressBar := progress.NewProgressBar(fileSize, fmt.Sprintf("Downloading %s", filepath.Base(localPath)))

	// Create progress writer
	progressWriter := progress.NewProgressWriter(file, progressBar)

	// Copy the response body to the file with progress tracking
	_, err = io.Copy(progressWriter, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to Ucopy file: %w", err)
	}

	// Finish progress bar
	progressBar.Finish()

	return nil
}

// isValidFile checks if a file exists and has the correct SHA256
func (d *Downloader) isValidFile(filePath, expectedSHA256 string) bool {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	// Calculate SHA256 of the file
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}

	actualSHA256 := hex.EncodeToString(hash.Sum(nil))
	return actualSHA256 == expectedSHA256
}

// Cleanup removes a downloaded file
func (d *Downloader) Cleanup(filePath string) error {
	return os.Remove(filePath)
}

// ListAvailableVersions fetches all available Go versions from the official page
func (d *Downloader) ListAvailableVersions() ([]VersionInfo, error) {
	// Fetch from the Go downloads page
	pageURL := d.baseURL + "/"

	resp, err := d.client.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases page: HTTP %d (check your internet connection)", resp.StatusCode)
	}

	// Read the HTML content
	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read page content: %w", err)
	}

	// Parse versions from HTML
	versions, err := d.parseVersionsFromHTML(string(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse versions: %w", err)
	}

	// Sort versions by version number (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i].Version, versions[j].Version) > 0
	})

	return versions, nil
}

// isCompatibleFile checks if a file is compatible with the current platform
func (d *Downloader) isCompatibleFile(file GoFile) bool {
	// Check OS compatibility
	osMatch := false
	switch runtime.GOOS {
	case "darwin":
		osMatch = file.OS == "darwin"
	case "linux":
		osMatch = file.OS == "linux"
	case "windows":
		osMatch = file.OS == "windows"
	default:
		osMatch = file.OS == runtime.GOOS
	}

	if !osMatch {
		return false
	}

	// Check architecture compatibility
	archMatch := false
	switch runtime.GOARCH {
	case "amd64":
		archMatch = file.Arch == "amd64" || file.Arch == "x86_64"
	case "arm64":
		archMatch = file.Arch == "arm64" || file.Arch == "aarch64"
	case "386":
		archMatch = file.Arch == "386" || file.Arch == "i386"
	default:
		archMatch = file.Arch == runtime.GOARCH
	}

	// Check if it's a source or binary file (we want binary)
	kindMatch := file.Kind == "archive" || file.Kind == "installer"

	return archMatch && kindMatch
}

// compareVersions compares two version strings
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Remove 'go' prefix for comparison
	v1 = strings.TrimPrefix(v1, "go")
	v2 = strings.TrimPrefix(v2, "go")

	// Parse version parts
	parts1 := parseVersionParts(v1)
	parts2 := parseVersionParts(v2)

	// Compare major version
	if parts1.major != parts2.major {
		if parts1.major > parts2.major {
			return 1
		}
		return -1
	}

	// Compare minor version
	if parts1.minor != parts2.minor {
		if parts1.minor > parts2.minor {
			return 1
		}
		return -1
	}

	// Compare patch version
	if parts1.patch != parts2.patch {
		if parts1.patch > parts2.patch {
			return 1
		}
		return -1
	}

	// Compare prerelease identifiers
	return comparePrerelease(parts1.prerelease, parts2.prerelease)
}

// versionParts represents parsed version components
type versionParts struct {
	major      int
	minor      int
	patch      int
	prerelease string
}

// parseVersionParts parses a version string into components
func parseVersionParts(version string) versionParts {
	parts := versionParts{}

	// Check for prerelease identifiers (rc, beta, alpha)
	prereleaseIndex := -1
	prereleasePatterns := []string{"rc", "beta", "alpha"}

	for _, pattern := range prereleasePatterns {
		if index := strings.Index(version, pattern); index != -1 {
			prereleaseIndex = index
			break
		}
	}

	// Extract prerelease if found
	if prereleaseIndex != -1 {
		parts.prerelease = version[prereleaseIndex:]
		version = version[:prereleaseIndex]
	}

	// Split by dots
	dotParts := strings.Split(version, ".")

	// Parse major version
	if len(dotParts) > 0 {
		parts.major = parseVersionNumber(dotParts[0])
	}

	// Parse minor version
	if len(dotParts) > 1 {
		parts.minor = parseVersionNumber(dotParts[1])
	}

	// Parse patch version
	if len(dotParts) > 2 {
		parts.patch = parseVersionNumber(dotParts[2])
	}

	return parts
}

// parseVersionNumber parses a numeric version part
func parseVersionNumber(s string) int {
	// Extract only numeric part
	numStr := ""
	for _, r := range s {
		if r >= '0' && r <= '9' {
			numStr += string(r)
		} else {
			break
		}
	}

	if numStr == "" {
		return 0
	}

	// Convert to int
	result := 0
	for _, r := range numStr {
		result = result*10 + int(r-'0')
	}
	return result
}

// comparePrerelease compares prerelease identifiers
func comparePrerelease(p1, p2 string) int {
	// Empty prerelease is considered newer than non-empty
	if p1 == "" && p2 == "" {
		return 0
	}
	if p1 == "" {
		return 1
	}
	if p2 == "" {
		return -1
	}

	// Parse prerelease identifiers
	parts1 := parsePrereleaseParts(p1)
	parts2 := parsePrereleaseParts(p2)

	// Compare each part
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var part1, part2 string
		if i < len(parts1) {
			part1 = parts1[i]
		}
		if i < len(parts2) {
			part2 = parts2[i]
		}

		// Compare parts
		cmp := comparePrereleasePart(part1, part2)
		if cmp != 0 {
			return cmp
		}
	}

	return 0
}

// parsePrereleaseParts splits prerelease string into parts
func parsePrereleaseParts(prerelease string) []string {
	// Split by dots
	parts := strings.Split(prerelease, ".")
	return parts
}

// comparePrereleasePart compares two prerelease parts
func comparePrereleasePart(p1, p2 string) int {
	// Handle empty parts
	if p1 == "" && p2 == "" {
		return 0
	}
	if p1 == "" {
		return -1
	}
	if p2 == "" {
		return 1
	}

	// Check if parts are numeric
	num1, isNum1 := parsePrereleaseNumber(p1)
	num2, isNum2 := parsePrereleaseNumber(p2)

	// If both are numeric, compare as numbers
	if isNum1 && isNum2 {
		if num1 < num2 {
			return -1
		}
		if num1 > num2 {
			return 1
		}
		return 0
	}

	// If one is numeric and one is not, numeric comes first
	if isNum1 && !isNum2 {
		return -1
	}
	if !isNum1 && isNum2 {
		return 1
	}

	// Both are non-numeric, compare lexicographically
	if p1 < p2 {
		return -1
	}
	if p1 > p2 {
		return 1
	}
	return 0
}

// parsePrereleaseNumber tries to parse a string as a number
func parsePrereleaseNumber(s string) (int, bool) {
	if s == "" {
		return 0, false
	}

	// Check if all characters are digits
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}
	}

	// Convert to int
	result := 0
	for _, r := range s {
		result = result*10 + int(r-'0')
	}
	return result, true
}

// parseVersionsFromHTML parses Go versions from the HTML page
func (d *Downloader) parseVersionsFromHTML(html string) ([]VersionInfo, error) {
	var versions []VersionInfo
	versionMap := make(map[string]VersionInfo)

	// Find all href patterns that contain Go versions
	// Look for patterns like href="/dl/go1.25.1.windows-amd64.msi"
	start := 0
	for {
		// Find the next href="/dl/go pattern
		hrefStart := strings.Index(html[start:], "href=\"/dl/go")
		if hrefStart == -1 {
			break
		}

		hrefStart += start + 6 // Skip 'href="'
		hrefEnd := strings.Index(html[hrefStart:], "\"")
		if hrefEnd == -1 {
			break
		}

		href := html[hrefStart : hrefStart+hrefEnd]
		version := d.extractVersionFromHref(href)
		if version != "" && !d.isVersionInMap(versionMap, version) {
			d.addVersionToMap(versionMap, version)
		}

		// Move past this match
		start = hrefStart + hrefEnd
	}

	// Convert map to slice
	for _, version := range versionMap {
		versions = append(versions, version)
	}

	return versions, nil
}

// extractVersionFromHref extracts version from href like "/dl/go1.25.1.windows-amd64.msi"
func (d *Downloader) extractVersionFromHref(href string) string {
	// Remove "/dl/" prefix
	if !strings.HasPrefix(href, "/dl/") {
		return ""
	}

	filename := strings.TrimPrefix(href, "/dl/")

	// Extract version from filename like "go1.25.1.windows-amd64.msi"
	// Find the platform part (after the version)
	if strings.HasPrefix(filename, "go") {
		// Look for common platform patterns
		platforms := []string{".windows-", ".darwin-", ".linux-", ".freebsd-", ".openbsd-", ".netbsd-", ".solaris-", ".aix-", ".dragonfly-", ".illumos-", ".plan9-", ".src."}

		for _, platform := range platforms {
			platformIndex := strings.Index(filename, platform)
			if platformIndex != -1 {
				version := filename[:platformIndex]
				if d.isVersionString(version) {
					return version
				}
			}
		}
	}

	return ""
}

// addVersionToMap adds a version to the version map
func (d *Downloader) addVersionToMap(versionMap map[string]VersionInfo, version string) {
	// Determine if it's stable (not beta, rc, etc.)
	stable := !strings.Contains(strings.ToLower(version), "beta") &&
		!strings.Contains(strings.ToLower(version), "rc") &&
		!strings.Contains(strings.ToLower(version), "alpha")

	// Create a compatible file entry for current platform
	compatibleFiles := []File{
		{
			Filename: d.getFilename(version),
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Size:     0,  // Will be filled when actually downloading
			SHA256:   "", // Will be filled when actually downloading
		},
	}

	versionMap[version] = VersionInfo{
		Version:     version,
		Stable:      stable,
		ReleaseDate: "", // We don't have release dates from HTML
		Files:       compatibleFiles,
	}
}

// isVersionString checks if a string looks like a Go version
func (d *Downloader) isVersionString(s string) bool {
	// Must start with "go" and contain dots
	if !strings.HasPrefix(s, "go") {
		return false
	}

	// Must contain at least one dot
	if !strings.Contains(s, ".") {
		return false
	}

	// Must be at least "go1.x" format
	if len(s) < 4 {
		return false
	}

	// Check if it's a valid version format
	versionPart := strings.TrimPrefix(s, "go")
	parts := strings.Split(versionPart, ".")
	if len(parts) < 2 {
		return false
	}

	// Each part should be numeric or contain valid version identifiers
	for _, part := range parts {
		if part == "" {
			return false
		}
		// Allow numeric parts and common version suffixes
		if !d.isValidVersionPart(part) {
			return false
		}
	}

	return true
}

// isValidVersionPart checks if a version part is valid
func (d *Downloader) isValidVersionPart(part string) bool {
	// Allow pure numeric parts
	if d.isNumeric(part) {
		return true
	}

	// Allow common version suffixes
	suffixes := []string{"beta", "rc", "alpha", "pre"}
	for _, suffix := range suffixes {
		if strings.Contains(strings.ToLower(part), suffix) {
			return true
		}
	}

	return false
}

// isNumeric checks if a string is numeric
func (d *Downloader) isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// normalizeVersionString normalizes a version string
func (d *Downloader) normalizeVersionString(s string) string {
	// Remove common HTML artifacts
	s = strings.Trim(s, " \t\n\r\"'<>")

	// Ensure it starts with "go"
	if !strings.HasPrefix(s, "go") {
		return ""
	}

	// Remove any trailing characters that aren't part of the version
	// Find the end of the version string
	end := len(s)
	for i, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '"' || r == '\'' || r == '>' {
			end = i
			break
		}
	}

	return s[:end]
}

// isVersionInMap checks if a version is already in the map
func (d *Downloader) isVersionInMap(m map[string]VersionInfo, version string) bool {
	_, exists := m[version]
	return exists
}
