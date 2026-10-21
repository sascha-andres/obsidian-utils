package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
)

// FillDailyNote sets each value whose key is missing or zero in the frontmatter of the daily note of day (2006-01-02) in dailyFolder.
// found is false if the note does not exist. changed reports whether a value was set. With dryRun the note is not written.
func FillDailyNote(logger *slog.Logger, dailyFolder, day string, values map[string]any, dryRun bool) (found, changed bool, err error) {
	t, err := time.Parse(time.DateOnly, day)
	if err != nil {
		return false, false, err
	}
	notePath := path.Join(dailyFolder, fmt.Sprintf("%s.md", t.Format("2006/01/2006-01-02")))
	found, err = Exists(notePath)
	if err != nil || !found {
		logger.Debug("daily note does not exist, skipping", "file", notePath)
		return found, false, err
	}

	processor := obsidianutils.NewSimpleFrontmatterProcessor(notePath)
	for key, value := range values {
		old, getErr := processor.GetValue(key)
		isZero, ok := IsZeroValue(old, getErr == nil)
		if !ok {
			logger.Warn("could not interpret value, skipping", "file", notePath, "key", key, "value", old)
			continue
		}
		if !isZero {
			continue
		}
		logger.Debug("setting value", "file", notePath, "key", key, "old", old, "new", value)
		if err := processor.SetValue(key, value); err != nil {
			return true, false, err
		}
		changed = true
	}
	if !changed || dryRun {
		return true, changed, nil
	}
	doc, err := processor.GenerateMarkDownDocument()
	if err != nil {
		return true, false, err
	}
	return true, true, os.WriteFile(notePath, doc, 0600)
}
