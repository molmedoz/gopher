package downloader

import (
	"io"
	"net/http"
)

// HTTPClient interface for HTTP operations
type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

// FileSystem interface for file operations
type FileSystem interface {
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
	MkdirAll(path string, perm uint32) error
	Remove(name string) error
	Stat(name string) (interface{}, error)
}

// ProgressWriter interface for progress tracking
type ProgressWriter interface {
	Write(p []byte) (n int, err error)
}
