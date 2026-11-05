package internal

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFillDailyNotesRespectsRange(t *testing.T) {
	dir := t.TempDir()
	for _, day := range []string{"2026-01-01", "2026-01-02", "2026-01-03"} {
		d, _ := time.Parse(time.DateOnly, day)
		notePath := DailyNotePath(dir, d)
		if err := os.MkdirAll(filepath.Dir(notePath), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(notePath, []byte("---\nsteps: 0\n---\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	values := map[string]map[string]any{
		"2026-01-01": {"steps": 1},
		"2026-01-02": {"steps": 2},
		"2026-01-03": {"steps": 3},
		"2026-01-04": {"steps": 4},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	scanned, updated, err := FillDailyNotes(logger, dir, values, "2026-01-02", "2026-01-04", false)
	if err != nil || scanned != 2 || updated != 2 {
		t.Fatalf("FillDailyNotes = %d, %d, %v, want 2, 2, nil", scanned, updated, err)
	}
}

func TestLastImportDay(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
	if got, err := LastImportDay("", ""); err != nil || got != yesterday {
		t.Errorf("LastImportDay(\"\", \"\") = %q, %v, want %q", got, err, yesterday)
	}
	if got, err := LastImportDay("", "2000-01-01"); err != nil || got != "2000-01-01" {
		t.Errorf("LastImportDay with earlier until = %q, %v, want 2000-01-01", got, err)
	}
	if _, err := LastImportDay("bad", ""); err == nil {
		t.Error("LastImportDay with invalid start date: want error")
	}
}
