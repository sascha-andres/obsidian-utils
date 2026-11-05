package meeting

import (
	"strings"
	"testing"
	"time"
)

func TestContent(t *testing.T) {
	// Fixed time for consistent test results
	fixedTime := time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		title        string
		appointment  time.Time
		wantContains []string
		wantErr      bool
	}{
		{
			name:        "Basic meeting content",
			title:       "Test Meeting",
			appointment: fixedTime,
			wantContains: []string{
				"tags:",
				"- meeting",
				"date: 2023-05-15T10:00:00Z",
				"title: Test Meeting",
				"[[2023-05-15]]",
				"# Meeting",
				"## Attendees",
				"## Notes",
			},
		},
		{
			name:         "Title is cleaned",
			title:        " a: b\nc ",
			appointment:  fixedTime,
			wantContains: []string{"title: a- b->c\n"},
		},
		{
			name:        "Empty title",
			title:       "  ",
			appointment: fixedTime,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := Content(tt.title, tt.appointment)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Content() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(content, want) {
					t.Errorf("Content() content does not contain %q\nContent: %s", want, content)
				}
			}
		})
	}
}
