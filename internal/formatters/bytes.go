package formatters

import "fmt"

// FormatBytes formats bytes into human readable format (B, KB, MB, GB, TB, PB, EB)
//
// Examples:
//   - FormatBytes(0)          -> "0 B"
//   - FormatBytes(1023)       -> "1023 B"
//   - FormatBytes(1024)       -> "1.0 KB"
//   - FormatBytes(1048576)    -> "1.0 MB"
//   - FormatBytes(1073741824) -> "1.0 GB"
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatSpeed formats bytes per second into human readable speed format
//
// Examples:
//   - FormatSpeed(0)           -> "0 B/s"
//   - FormatSpeed(500)         -> "500 B/s"
//   - FormatSpeed(1536)        -> "1.5 KB/s"
//   - FormatSpeed(2097152)     -> "2.0 MB/s"
//   - FormatSpeed(10485760)    -> "10.0 MB/s"
func FormatSpeed(bytesPerSecond float64) string {
	const unit = 1024

	if bytesPerSecond < unit {
		return fmt.Sprintf("%.0f B/s", bytesPerSecond)
	}

	if bytesPerSecond < unit*unit {
		return fmt.Sprintf("%.1f KB/s", bytesPerSecond/unit)
	}

	if bytesPerSecond < unit*unit*unit {
		return fmt.Sprintf("%.1f MB/s", bytesPerSecond/(unit*unit))
	}

	return fmt.Sprintf("%.1f GB/s", bytesPerSecond/(unit*unit*unit))
}

// FormatPercentage formats a value (0.0 to 1.0) as a percentage string
//
// Examples:
//   - FormatPercentage(0.0)    -> "0.0%"
//   - FormatPercentage(0.5)    -> "50.0%"
//   - FormatPercentage(0.753)  -> "75.3%"
//   - FormatPercentage(1.0)    -> "100.0%"
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.1f%%", value*100)
}
