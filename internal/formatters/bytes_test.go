package formatters

import "testing"

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"zero bytes", 0, "0 B"},
		{"one byte", 1, "1 B"},
		{"small bytes", 512, "512 B"},
		{"max bytes before KB", 1023, "1023 B"},
		{"exactly 1 KB", 1024, "1.0 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"max KB before MB", 1048575, "1024.0 KB"},
		{"exactly 1 MB", 1048576, "1.0 MB"},
		{"2.5 MB", 2621440, "2.5 MB"},
		{"exactly 1 GB", 1073741824, "1.0 GB"},
		{"1.5 GB", 1610612736, "1.5 GB"},
		{"exactly 1 TB", 1099511627776, "1.0 TB"},
		{"2 TB", 2199023255552, "2.0 TB"},
		{"exactly 1 PB", 1125899906842624, "1.0 PB"},
		{"exactly 1 EB", 1152921504606846976, "1.0 EB"},
		{"large EB value", 2305843009213693952, "2.0 EB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.input)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"zero speed", 0, "0 B/s"},
		{"small speed", 100, "100 B/s"},
		{"500 B/s", 500, "500 B/s"},
		{"max bytes per second", 1023, "1023 B/s"},
		{"exactly 1 KB/s", 1024, "1.0 KB/s"},
		{"1.5 KB/s", 1536, "1.5 KB/s"},
		{"10 KB/s", 10240, "10.0 KB/s"},
		{"exactly 1 MB/s", 1048576, "1.0 MB/s"},
		{"2.5 MB/s", 2621440, "2.5 MB/s"},
		{"10 MB/s", 10485760, "10.0 MB/s"},
		{"exactly 1 GB/s", 1073741824, "1.0 GB/s"},
		{"2 GB/s", 2147483648, "2.0 GB/s"},
		{"very fast", 5368709120, "5.0 GB/s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSpeed(tt.input)
			if result != tt.expected {
				t.Errorf("FormatSpeed(%.0f) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatPercentage(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"zero percent", 0.0, "0.0%"},
		{"one percent", 0.01, "1.0%"},
		{"ten percent", 0.1, "10.0%"},
		{"twenty five percent", 0.25, "25.0%"},
		{"fifty percent", 0.5, "50.0%"},
		{"seventy five point three", 0.753, "75.3%"},
		{"ninety nine point nine", 0.999, "99.9%"},
		{"one hundred percent", 1.0, "100.0%"},
		{"over hundred percent", 1.5, "150.0%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPercentage(tt.input)
			if result != tt.expected {
				t.Errorf("FormatPercentage(%.3f) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkFormatBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatBytes(1073741824) // 1 GB
	}
}

func BenchmarkFormatSpeed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatSpeed(10485760) // 10 MB/s
	}
}

func BenchmarkFormatPercentage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatPercentage(0.753)
	}
}
