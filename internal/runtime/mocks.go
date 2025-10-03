package runtime

import (
	"github.com/molmedoz/gopher/internal/downloader"
)

// MockDownloader implements minimal downloader interface for testing
type MockDownloader struct {
	DownloadFunc              func(version, downloadDir string) (string, error)
	CleanupFunc               func(version string) error
	ListAvailableVersionsFunc func() ([]downloader.VersionInfo, error)
}

func (m *MockDownloader) Download(version, downloadDir string) (string, error) {
	if m.DownloadFunc != nil {
		return m.DownloadFunc(version, downloadDir)
	}
	return "/tmp/go" + version + ".tar.gz", nil
}

func (m *MockDownloader) Cleanup(version string) error {
	if m.CleanupFunc != nil {
		return m.CleanupFunc(version)
	}
	return nil
}

func (m *MockDownloader) ListAvailableVersions() ([]downloader.VersionInfo, error) {
	if m.ListAvailableVersionsFunc != nil {
		return m.ListAvailableVersionsFunc()
	}
	return []downloader.VersionInfo{
		{Version: "go1.21.0", Stable: true, ReleaseDate: "2023-08-08", Files: []downloader.File{}},
		{Version: "go1.20.0", Stable: true, ReleaseDate: "2023-02-01", Files: []downloader.File{}},
	}, nil
}

// Note: MockInstaller, MockSystemDetector, MockSymlinkManager, and MockFileSystem
// have been removed after adapter refactoring. Tests now use real implementations
// or need to be updated to work without these mocks.
//
// For tests that need mocking, use test doubles for installer.Installer and
// downloader.Downloader directly, or use real instances with temp directories.
