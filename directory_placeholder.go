package obsidianutils

import (
	"os"
	"strings"
)

// ApplyDirectoryPlaceHolder replaces the placeholder "$$PWD$$" in the folder string with the current working directory.
func ApplyDirectoryPlaceHolder(folder string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	folder = strings.Replace(folder, "$$PWD$$", currentDir, -1)
	return folder, nil
}
