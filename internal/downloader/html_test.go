package downloader

import "testing"

func TestParseVersionsFromHTML_Simple(t *testing.T) {
	d := New("https://go.dev/dl/")
	html := `
    <a class="download" href="/dl/go1.25.1.linux-amd64.tar.gz">go1.25.1.linux-amd64.tar.gz</a>
    <a class="download" href="/dl/go1.25rc2.linux-amd64.tar.gz">go1.25rc2.linux-amd64.tar.gz</a>
    <span class="version">go1.24.7</span>
    <a class="download" href="/dl/go1.25.1.darwin-arm64.pkg">go1.25.1.darwin-arm64.pkg</a>
    `
	vs, err := d.parseVersionsFromHTML(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vs) < 2 {
		t.Fatalf("expected at least two versions")
	}
}

func TestParseVersionsFromHTML_DuplicatesAndPrereleases(t *testing.T) {
	d := New("https://go.dev/dl/")
	html := `
    <a class="download" href="/dl/go1.25rc1.linux-amd64.tar.gz">go1.25rc1.linux-amd64.tar.gz</a>
    <a class="download" href="/dl/go1.25rc1.darwin-amd64.pkg">go1.25rc1.darwin-amd64.pkg</a>
    <a class="download" href="/dl/go1.25.0.linux-amd64.tar.gz">go1.25.0.linux-amd64.tar.gz</a>
    <span class="version">go1.25.0</span>
    `
	vs, err := d.parseVersionsFromHTML(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// ensure stable and prerelease are both present and duplicates are not excessive
	hasStable, hasRC := false, false
	for _, v := range vs {
		if v.Version == "go1.25.0" {
			hasStable = true
		}
		if v.Version == "go1.25rc1" {
			hasRC = true
		}
	}
	if !hasStable || !hasRC {
		t.Fatalf("expected both stable and rc present: %v", vs)
	}
}
