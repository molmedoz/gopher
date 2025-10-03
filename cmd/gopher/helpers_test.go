package main

import (
	"testing"

	"github.com/molmedoz/gopher/internal/downloader"
)

func sampleList() []downloader.VersionInfo {
	return []downloader.VersionInfo{
		{Version: "go1.25.1", Stable: true},
		{Version: "go1.25rc1", Stable: false},
		{Version: "go1.24.7", Stable: true},
	}
}

func TestFilterVersionsHelper(t *testing.T) {
	list := sampleList()
	got := filterVersionsHelper(list, "1.25", false)
	if len(got) != 2 {
		t.Fatalf("want 2 got %d", len(got))
	}
	got = filterVersionsHelper(list, "1.25", true)
	if len(got) != 1 || got[0].Version != "go1.25.1" {
		t.Fatalf("stable filter mismatch")
	}
}

func TestPaginateHelper(t *testing.T) {
	list := sampleList()
	p := paginateHelper(list, 1, 2)
	if len(p) != 2 {
		t.Fatalf("want 2 got %d", len(p))
	}
	p2 := paginateHelper(list, 2, 2)
	if len(p2) != 1 {
		t.Fatalf("want 1 got %d", len(p2))
	}
	p3 := paginateHelper(list, 3, 2)
	if len(p3) != 0 {
		t.Fatalf("want 0 got %d", len(p3))
	}
}

func TestPaginateHelper_Edges(t *testing.T) {
	list := sampleList()
	// pageSize <= 0 defaults to 10
	p := paginateHelper(list, 1, 0)
	if len(p) != len(list) {
		t.Fatalf("default pageSize should return all when small list")
	}
	// page <= 0 defaults to 1
	p2 := paginateHelper(list, 0, 1)
	if len(p2) != 1 {
		t.Fatalf("page<=0 should default to page 1")
	}
}

func TestFilterVersionsHelper_Edges(t *testing.T) {
	list := sampleList()
	// empty filter returns all (subject to stableOnly)
	got := filterVersionsHelper(list, "", false)
	if len(got) != len(list) {
		t.Fatalf("empty filter should return all")
	}
	// stableOnly with empty filter returns only stable
	got2 := filterVersionsHelper(list, "", true)
	if len(got2) != 2 {
		t.Fatalf("stableOnly should filter out prereleases")
	}
}

func TestFilterVersionsHelper_NoMatches(t *testing.T) {
	list := sampleList()
	got := filterVersionsHelper(list, "9.99", false)
	if len(got) != 0 {
		t.Fatalf("expected no matches, got %d", len(got))
	}
}
