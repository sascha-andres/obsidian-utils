package internal

import (
	"log/slog"
	"maps"
	"slices"
	"time"
)

// LastImportDay validates startDate and until (2006-01-02 or empty) and returns the last day to import:
// yesterday, because today's data is not complete, or until if it is earlier.
func LastImportDay(startDate, until string) (string, error) {
	lastDay := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
	if until != "" {
		if _, err := time.Parse(time.DateOnly, until); err != nil {
			return "", err
		}
		lastDay = min(lastDay, until)
	}
	if startDate != "" {
		if _, err := time.Parse(time.DateOnly, startDate); err != nil {
			return "", err
		}
	}
	return lastDay, nil
}

// FillDailyNotes calls FillDailyNote in date order for each day (2006-01-02) of values from startDate (empty for no limit) to lastDay.
// scanned counts the daily notes that exist, updated the notes that changed.
func FillDailyNotes(logger *slog.Logger, dailyFolder string, values map[string]map[string]any, startDate, lastDay string, dryRun bool) (scanned, updated int, err error) {
	for _, day := range slices.Sorted(maps.Keys(values)) {
		if day < startDate || day > lastDay {
			continue
		}
		found, changed, err := FillDailyNote(logger, dailyFolder, day, values[day], dryRun)
		if err != nil {
			return scanned, updated, err
		}
		if found {
			scanned++
		}
		if changed {
			updated++
		}
	}
	return scanned, updated, nil
}
