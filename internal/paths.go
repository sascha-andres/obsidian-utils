package internal

import (
	"errors"
	"fmt"
	"path"
	"time"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
)

// ResolveFolder checks that folder and sub are non empty, applies the $$PWD$$ placeholder to folder and joins both.
// subFlag is the flag name of sub, used in the error message.
func ResolveFolder(folder, sub, subFlag string) (string, error) {
	if folder == "" {
		return "", errors.New("-folder must be non empty")
	}
	if sub == "" {
		return "", fmt.Errorf("-%s must be non empty", subFlag)
	}
	resolved, err := obsidianutils.ApplyDirectoryPlaceHolder(folder)
	if err != nil {
		return "", err
	}
	return path.Join(resolved, sub), nil
}

// DailyNotePath returns the path of the daily note of t inside dailyFolder.
func DailyNotePath(dailyFolder string, t time.Time) string {
	return path.Join(dailyFolder, t.Format("2006/01/2006-01-02")+".md")
}
