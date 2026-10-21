package main

import (
	"strings"
	"testing"
	"time"
)

const sampleReport = "\ufeff" + `"Description","Duration","Member","Email","Project","Tags","Start date","Start time","Stop date","Stop time"
"a","1:00:29","m","e","IQ","","2026-10-20","08:00:00","2026-10-20","09:00:29"
"b","0:30:00","m","e","IQ","","2026-10-20","10:00:00","2026-10-20","10:30:00"
"c","2:00:00","m","e","Other","","2026-10-20","12:00:00","2026-10-20","14:00:00"
"d","1:00:00","m","e","IQ","","2026-10-20","23:30:00","2026-10-21","00:30:00"
`

func TestReadWorkTime(t *testing.T) {
	tests := []struct {
		name     string
		projects map[string]bool
		want     map[string]time.Duration
	}{
		{
			name:     "project filter and split at midnight",
			projects: map[string]bool{"IQ": true},
			want:     map[string]time.Duration{"2026-10-20": 2*time.Hour + 29*time.Second, "2026-10-21": 30 * time.Minute},
		},
		{
			name: "no filter",
			want: map[string]time.Duration{"2026-10-20": 4*time.Hour + 29*time.Second, "2026-10-21": 30 * time.Minute},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readWorkTime(strings.NewReader(sampleReport), tt.projects)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("readWorkTime() = %v, want %v", got, tt.want)
			}
			for day, want := range tt.want {
				if got[day] != want {
					t.Fatalf("readWorkTime()[%s] = %v, want %v", day, got[day], want)
				}
			}
		})
	}
}

func TestReadWorkTimeMissingColumn(t *testing.T) {
	if _, err := readWorkTime(strings.NewReader("\"Project\",\"Start date\"\n"), nil); err == nil {
		t.Fatal("readWorkTime() without stop columns must fail")
	}
}
