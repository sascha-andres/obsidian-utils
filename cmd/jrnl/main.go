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
	"github.com/sascha-andres/obsidian-utils/internal/jrnl"
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
	defer logger.Debug("journal note added")

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

	defaultEntry := strings.Join(flag.GetVerbs(), " ")
	bulletPoint, err := internal.PromptText("Journal entry", defaultEntry, func(s string) error {
		if len(s) == 0 {
			return errors.New("journal entry cannot be empty")
		}
		return nil
	})

	if err != nil {
		return err
	}
	line := ""
	if strings.HasPrefix(bulletPoint, "[ ]") {
		line = fmt.Sprintf("%s", bulletPoint)
	} else {
		line = fmt.Sprintf("(%s) %s", time.Now().Format("15:04"), bulletPoint)
	}

	newFileData, err := jrnl.AddBulletpoint(fileData, line, headline)
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

	return os.WriteFile(resultingFile, newFileData, 0640)
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
