package main

import (
	"strings"

	"github.com/molmedoz/gopher/internal/downloader"
)

func filterVersionsHelper(list []downloader.VersionInfo, filter string, stableOnly bool) []downloader.VersionInfo {
	res := make([]downloader.VersionInfo, 0, len(list))
	f := strings.ToLower(filter)
	for _, v := range list {
		if stableOnly && !v.Stable {
			continue
		}
		if f == "" || strings.Contains(strings.ToLower(v.Version), f) {
			res = append(res, v)
		}
	}
	return res
}

func paginateHelper(list []downloader.VersionInfo, page, pageSize int) []downloader.VersionInfo {
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * pageSize
	if start >= len(list) {
		return []downloader.VersionInfo{}
	}
	end := start + pageSize
	if end > len(list) {
		end = len(list)
	}
	return list[start:end]
}
