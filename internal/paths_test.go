package internal

import "testing"

func TestResolveFolder(t *testing.T) {
	tests := []struct {
		folder, sub, want string
		wantErr           bool
	}{
		{"/vault", "daily", "/vault/daily", false},
		{"", "daily", "", true},
		{"/vault", "", "", true},
	}
	for _, tt := range tests {
		got, err := ResolveFolder(tt.folder, tt.sub, "daily-folder")
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("ResolveFolder(%q, %q) = %q, %v, want %q, error %t", tt.folder, tt.sub, got, err, tt.want, tt.wantErr)
		}
	}
}
