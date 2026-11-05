package internal

import (
	"testing"
	"time"
)

func TestResolveDateArg(t *testing.T) {
	tests := []struct {
		raw, want string
		wantErr   bool
	}{
		{"2026-01-02", "2026-01-02", false},
		{"+2", time.Now().AddDate(0, 0, 2).Format(time.DateOnly), false},
		{"-3", time.Now().AddDate(0, 0, -3).Format(time.DateOnly), false},
		{"+x", "", true},
	}
	for _, tt := range tests {
		got, err := ResolveDateArg(tt.raw)
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("ResolveDateArg(%q) = %q, %v, want %q, error %t", tt.raw, got, err, tt.want, tt.wantErr)
		}
	}
}
