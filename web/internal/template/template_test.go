package template

import (
	"testing"
)

func TestNewFormatTimeFunc(t *testing.T) {
	formatTime := NewFormatTimeFunc("Asia/Shanghai")
	want := "2026-01-01 00:00"

	tests := []struct {
		name string
		in   any
	}{
		{
			name: "rfc3339 string",
			in:   "2025-12-31T16:00:00Z",
		},
		{
			name: "rfc3339 string with nanos",
			in:   "2025-12-31T16:00:00.000000000Z",
		},
		{
			name: "formatted string",
			in:   "2026-01-01 00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTime(tt.in); got != want {
				t.Fatalf("formatTime() = %q, want %q", got, want)
			}
		})
	}
}
