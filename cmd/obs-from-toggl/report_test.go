package main

import (
	"maps"
	"strings"
	"testing"
	"time"
)

const sampleReport = "\ufeff" + `"Client","Project","Duration","Start date","Start time","Stop date","Stop time","Tags","Description","Member","Email"
"Work","IQ","1:00:29","2026-10-20","08:00:00","2026-10-20","09:00:29","","a","m","e"
"Work","IQ","0:30:00","2026-10-20","10:00:00","2026-10-20","10:30:00","","b","m","e"
"","Other","2:00:00","2026-10-20","12:00:00","2026-10-20","14:00:00","","c","m","e"
"Work","IQ","1:00:00","2026-10-20","23:30:00","2026-10-21","00:30:00","","d","m","e"
"Work","Admin","0:15:00","2026-10-21","09:00:00","2026-10-21","09:15:00","","e","m","e"
`

func TestReadWorkTime(t *testing.T) {
	tests := []struct {
		name     string
		projects map[string]bool
		clients  map[string]bool
		want     map[string]time.Duration
		other    map[string]time.Duration
	}{
		{
			name:     "project filter and split at midnight",
			projects: map[string]bool{"IQ": true},
			want:     map[string]time.Duration{"2026-10-20": 2*time.Hour + 29*time.Second, "2026-10-21": 30 * time.Minute},
			other:    map[string]time.Duration{"2026-10-20": 2 * time.Hour, "2026-10-21": 15 * time.Minute},
		},
		{
			name:    "client filter",
			clients: map[string]bool{"Work": true},
			want:    map[string]time.Duration{"2026-10-20": 2*time.Hour + 29*time.Second, "2026-10-21": 45 * time.Minute},
			other:   map[string]time.Duration{"2026-10-20": 2 * time.Hour},
		},
		{
			name:     "client or project",
			projects: map[string]bool{"Other": true},
			clients:  map[string]bool{"Work": true},
			want:     map[string]time.Duration{"2026-10-20": 4*time.Hour + 29*time.Second, "2026-10-21": 45 * time.Minute},
			other:    map[string]time.Duration{},
		},
		{
			name:  "no filter",
			want:  map[string]time.Duration{"2026-10-20": 4*time.Hour + 29*time.Second, "2026-10-21": 45 * time.Minute},
			other: map[string]time.Duration{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOther, err := readWorkTime(strings.NewReader(sampleReport), tt.projects, tt.clients)
			if err != nil {
				t.Fatal(err)
			}
			if !maps.Equal(got, tt.want) {
				t.Fatalf("readWorkTime() work = %v, want %v", got, tt.want)
			}
			if !maps.Equal(gotOther, tt.other) {
				t.Fatalf("readWorkTime() other = %v, want %v", gotOther, tt.other)
			}
		})
	}
}

func TestReadWorkTimeMissingColumn(t *testing.T) {
	if _, _, err := readWorkTime(strings.NewReader("\"Project\",\"Start date\"\n"), nil, nil); err == nil {
		t.Fatal("readWorkTime() without stop columns must fail")
	}
	header := "\"Project\",\"Start date\",\"Start time\",\"Stop date\",\"Stop time\"\n"
	if _, _, err := readWorkTime(strings.NewReader(header), nil, nil); err != nil {
		t.Fatalf("readWorkTime() without client filter must not need the client column: %v", err)
	}
	if _, _, err := readWorkTime(strings.NewReader(header), nil, map[string]bool{"Work": true}); err == nil {
		t.Fatal("readWorkTime() with client filter must fail without client column")
	}
}

func TestReadWorkTimeColumnOrder(t *testing.T) {
	report := `"Stop time","Description","Stop date","Start time","Project","Start date","Client"
"09:00:00","a","2026-10-20","08:00:00","IQ","2026-10-20","Work"
"11:00:00","b","2026-10-20","10:00:00","Other","2026-10-20",""
`
	got, gotOther, err := readWorkTime(strings.NewReader(report), nil, map[string]bool{"Work": true})
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]time.Duration{"2026-10-20": time.Hour}; !maps.Equal(got, want) || !maps.Equal(gotOther, want) {
		t.Fatalf("readWorkTime() work = %v, other = %v, want %v for both", got, gotOther, want)
	}
}
