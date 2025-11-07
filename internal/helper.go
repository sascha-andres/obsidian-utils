package internal

import (
	"errors"
	"os"
)

// Exists checks if the file or directory at the specified path exists and is accessible. Returns true if the file or directory exists and is accessible, false otherwise.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	// Some other error (e.g., permission issues)
	return false, err
}
