package obsidianutils

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateFileName(t *testing.T) {
	// Fixed time for consistent test results
	fixedTime := time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC)

	// Get current working directory for testing directory placeholder
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []struct {
		name          string
		folder        string
		title         string
		noDatePrefix  bool
		timeForPrefix time.Time
		want          string
		wantErr       bool
	}{
		{
			name:          "Basic title without date prefix",
			folder:        "test",
			title:         "Test Note",
			noDatePrefix:  true,
			timeForPrefix: fixedTime,
			want:          filepath.Join("test", "Test Note.md"),
			wantErr:       false,
		},
		{
			name:          "Basic title with date prefix",
			folder:        "test",
			title:         "Test Note",
			noDatePrefix:  false,
			timeForPrefix: fixedTime,
			want:          filepath.Join("test", "2023-05-15 Test Note.md"),
			wantErr:       false,
		},
		{
			name:          "Title with special characters",
			folder:        "test",
			title:         "Äpfel und Öl: Übungen für Straße",
			noDatePrefix:  true,
			timeForPrefix: fixedTime,
			want:          filepath.Join("test", "Aepfel und Oel Uebungen fuer Strasse.md"),
			wantErr:       false,
		},
		{
			name:          "Title with special characters and date prefix",
			folder:        "test",
			title:         "Äpfel und Öl: Übungen für Straße",
			noDatePrefix:  false,
			timeForPrefix: fixedTime,
			want:          filepath.Join("test", "2023-05-15 Aepfel und Oel Uebungen fuer Strasse.md"),
			wantErr:       false,
		},
		{
			name:          "Empty title",
			folder:        "test",
			title:         "",
			noDatePrefix:  false,
			timeForPrefix: fixedTime,
			want:          filepath.Join("test", "2023-05-15 .md"),
			wantErr:       false,
		},
		{
			name:          "Empty folder",
			folder:        "",
			title:         "Test Note",
			noDatePrefix:  false,
			timeForPrefix: fixedTime,
			want:          "2023-05-15 Test Note.md",
			wantErr:       false,
		},
		{
			name:          "With directory placeholder",
			folder:        "$$PWD$$/test",
			title:         "Test Note",
			noDatePrefix:  false,
			timeForPrefix: fixedTime,
			want:          filepath.Join(currentDir, "test", "2023-05-15 Test Note.md"),
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateFileName(tt.folder, tt.title, tt.noDatePrefix, tt.timeForPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
