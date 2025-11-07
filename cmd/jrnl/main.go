package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
)

var (
	folder, forDate, dailyFolder string
	headline                     = "## Other stuff"
	logLevel                     string
	dryRun                       bool
)

func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_JRNL")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&forDate, "for-date", time.Now().Format(time.DateOnly), "date for which to create the daily note (2006-01-02)")
	flag.StringVar(&headline, "headline", headline, fmt.Sprintf("headline under which to place the journal note (default: %s)", headline))
	flag.BoolVar(&dryRun, "dry-run", false, "pass to not edit file but to print added line with some context")
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(_ context.Context) error {
	flag.Parse()

	logger := internal.CreateLogger("OBS_UTIL_DAILY", logLevel)

	logger.Debug("start adding a journal note")
	dailyNoteFolder, err := constructFolder()
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", forDate)
	if err != nil {
		return err
	}

	resultingFile := path.Join(dailyNoteFolder, fmt.Sprintf("%s.md", t.Format("2006/01/2006-01-02")))
	if e, _ := internal.Exists(resultingFile); !e {
		return fmt.Errorf("file %s does not exist, consider to create with daily", resultingFile)
	}

	fileData, err := os.ReadFile(resultingFile)
	if err != nil {
		return err
	}

	newFileData, err := addBulletpoint(fileData, "xxx", headline)
	if err != nil {
		logger.Error("could not add bullet point", "err", err, "file", resultingFile, "headline", headline)
		return err
	}

	if dryRun {
		d := cmp.Diff(string(fileData), string(newFileData))
		if d == "" {
			logger.Error("no changes detected")
			return fmt.Errorf("no changes detected")
		}
		fmt.Println(d)
		return nil
	}

	return nil
}

func addBulletpoint(data []byte, bulletPoint, after string) ([]byte, error) {
	// Convert to lines for easier manipulation
	content := string(data)
	lines := strings.Split(content, "\n")

	trimEq := func(a, b string) bool { return strings.TrimSpace(a) == strings.TrimSpace(b) }
	isHeadline := func(s string) bool {
		// A markdown headline starts with one or more '#'
		s = strings.TrimSpace(s)
		return strings.HasPrefix(s, "#")
	}
	isULItem := func(s string) (bool, string, rune) {
		// Detect unordered list item of the form: optional spaces + ('-', '*', '+') + space
		// Returns (ok, indent, marker)
		if s == "" {
			return false, "", 0
		}
		// Count leading spaces
		i := 0
		for i < len(s) && s[i] == ' ' {
			i++
		}
		indent := s[:i]
		rest := s[i:]
		if len(rest) < 2 {
			return false, "", 0
		}
		switch rest[0] {
		case '-', '*', '+':
			if len(rest) >= 2 && rest[1] == ' ' {
				return true, indent, rune(rest[0])
			}
		}
		return false, "", 0
	}

	// 1) Find the starting point line matching 'after'
	start := -1
	for i := range lines {
		if trimEq(lines[i], after) {
			start = i
			break
		}
	}
	if start == -1 {
		return nil, fmt.Errorf("anchor line not found: %q", after)
	}

	// 2) Determine the scan window: from line after 'after' to before the next headline
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if isHeadline(lines[i]) {
			end = i
			break
		}
	}

	// 3) Search for the first unordered list in [start+1, end)
	listStart := -1
	listIndent := ""
	listMarker := '-'
	for i := start + 1; i < end; i++ {
		if ok, indent, marker := isULItem(lines[i]); ok {
			listStart = i
			listIndent = indent
			listMarker = marker
			break
		}
	}

	if listStart != -1 {
		// 3a) Found a list. Find the last consecutive list item line to append after.
		listEnd := listStart
		for i := listStart; i < end; i++ {
			if ok, _, _ := isULItem(lines[i]); ok {
				listEnd = i
				continue
			}
			break
		}
		// Insert a new list item after listEnd
		insertion := listEnd + 1
		newLine := fmt.Sprintf("%s%c %s", listIndent, listMarker, bulletPoint)
		// Insert while preserving order
		lines = append(lines[:insertion], append([]string{newLine}, lines[insertion:]...)...)
	} else {
		// 3b) No list found before next headline. Create a new list with the bullet point as first item.
		insertion := start + 1
		newLines := []string{}
		// Add a blank line if the immediate next line is not blank and not end boundary
		if insertion < len(lines) && strings.TrimSpace(lines[insertion]) != "" && insertion < end {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, fmt.Sprintf("- %s", bulletPoint))
		// If we inserted right before a headline without a separating blank line, add one for readability
		if insertion < len(lines) && insertion < end && isHeadline(lines[insertion]) {
			newLines = append(newLines, "")
		}
		lines = append(lines[:insertion], append(newLines, lines[insertion:]...)...)
	}

	return []byte(strings.Join(lines, "\n")), nil
}

// constructFolder validates and processes folder paths, applies placeholders, and adjusts dates based on input parameters.
func constructFolder() (string, error) {
	if folder == "" {
		return "", errors.New("-folder must be non empty")
	}
	folder, err := obsidianutils.ApplyDirectoryPlaceHolder(folder)
	if err != nil {
		return "", err
	}
	if dailyFolder == "" {
		return "", errors.New("-daily-folder must be non empty")
	}
	if forDate == "" {
		forDate = time.Now().Format(time.DateOnly)
	}
	if strings.HasPrefix(forDate, "-") {
		relativeString := strings.TrimPrefix(forDate, "-")
		offset, err := strconv.Atoi(relativeString)
		if err != nil {
			return "", err
		}
		forDate = time.Now().AddDate(0, 0, offset*-1).Format(time.DateOnly)
	}
	if strings.HasPrefix(forDate, "+") {
		relativeString := strings.TrimPrefix(forDate, "+")
		offset, err := strconv.Atoi(relativeString)
		if err != nil {
			return "", err
		}
		forDate = time.Now().AddDate(0, 0, offset).Format(time.DateOnly)
	}

	folder = path.Join(folder, dailyFolder)

	return folder, nil
}
