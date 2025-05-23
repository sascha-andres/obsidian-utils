package meeting

import (
	"strings"
	"testing"
	"time"
)

func TestNewMeeting(t *testing.T) {
	tests := []struct {
		name    string
		opts    []OptionFunc
		wantErr bool
	}{
		{
			name:    "Valid meeting with title",
			opts:    []OptionFunc{WithTitle("Test Meeting")},
			wantErr: false,
		},
		{
			name:    "Invalid meeting without title",
			opts:    []OptionFunc{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMeeting(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMeeting() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithTitle(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		wantTitle string
	}{
		{
			name:     "Set title",
			title:    "Test Meeting",
			wantTitle: "Test Meeting",
		},
		{
			name:     "Set empty title",
			title:    "",
			wantTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Meeting{}
			opt := WithTitle(tt.title)
			err := opt(m)
			if err != nil {
				t.Errorf("WithTitle() error = %v", err)
			}
			if m.title != tt.wantTitle {
				t.Errorf("WithTitle() title = %v, want %v", m.title, tt.wantTitle)
			}
		})
	}
}

func TestCreateContent(t *testing.T) {
	// Fixed time for consistent test results
	fixedTime := time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name        string
		meetingTitle string
		contentTitle string
		appointment time.Time
		wantContains []string
		wantErr     bool
	}{
		{
			name:        "Basic meeting content",
			meetingTitle: "Test Meeting",
			contentTitle: "Test Meeting",
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
			wantErr:     false,
		},
		{
			name:        "Meeting with different title",
			meetingTitle: "Meeting Title",
			contentTitle: "Content Title",
			appointment: fixedTime,
			wantContains: []string{
				"title: Meeting Title",
				"[[2023-05-15]]",
			},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMeeting(WithTitle(tt.meetingTitle))
			if err != nil {
				t.Fatalf("Failed to create meeting: %v", err)
			}
			
			content, err := m.CreateContent(tt.contentTitle, tt.appointment)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			for _, want := range tt.wantContains {
				if !strings.Contains(content, want) {
					t.Errorf("CreateContent() content does not contain %q\nContent: %s", want, content)
				}
			}
		})
	}
}