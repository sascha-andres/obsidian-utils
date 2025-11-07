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

	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
)

var (
	folder, forDate, dailyFolder string
	headline                     = "## Other"
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
	flag.BoolVar(&dryRun, "dry-run", false, "pass to not edit file but to print added line with +- 5 lines")
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

	newFileData, err := addBulletpoint(fileData, headline)
	if err != nil {
		logger.Error("could not add bullet point", "err", err, "file", resultingFile, "headline", headline)
		return err
	}
	_ = newFileData

	return nil
}

func addBulletpoint(data []byte, bulletPoint, after string) ([]byte, error) {

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
