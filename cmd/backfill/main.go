package main

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"go.livingit.de/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
)

var (
	folder, dailyFolder, startDateFlag, logLevel, frontmatterKey string
	printConfig, dryRun                                          bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_BACKFILL")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where the daily notes are stored inside the vault")
	flag.StringVar(&startDateFlag, "start-date", "", "lowest date to start scanning from (2006-01-02)")
	flag.StringVar(&frontmatterKey, "frontmatter-key", "", "frontmatter key to backfill")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&dryRun, "dry-run", false, "do not write files, only log what would change")
}

// main is the entry point of the program.
func main() {
	flag.Parse()
	if printConfig {
		internal.PrintFlags()
	}

	logger := internal.CreateLogger(logLevel, "OBS_UTIL_BACKFILL")
	if err := run(logger); err != nil {
		logger.Error("error running backfill", "err", err)
		os.Exit(1)
	}
}

// run scans daily notes from -start-date up to today and forward-fills zero/missing entries for -frontmatter-key.
func run(logger *slog.Logger) error {
	dailyNoteFolder, err := internal.ResolveFolder(folder, dailyFolder, "daily-folder")
	if err != nil {
		return err
	}
	if startDateFlag == "" {
		return errors.New("-start-date must be non empty")
	}
	if frontmatterKey == "" {
		return errors.New("-frontmatter-key must be non empty")
	}

	startDate, err := time.Parse(time.DateOnly, startDateFlag)
	if err != nil {
		return err
	}
	today, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	if err != nil {
		return err
	}

	if printConfig {
		logger.Info("daily-folder", "daily-folder", dailyNoteFolder)
		logger.Info("start-date", "start-date", startDate.Format(time.DateOnly))
		logger.Info("frontmatter-key", "frontmatter-key", frontmatterKey)
		return nil
	}

	if startDate.After(today) {
		return errors.New("-start-date must not be in the future")
	}

	var (
		memoized     any
		haveMemoized bool
		scanned      int
		updated      int
	)

	for t := startDate; !t.After(today); t = t.AddDate(0, 0, 1) {
		notePath := internal.DailyNotePath(dailyNoteFolder, t)
		exists, err := internal.Exists(notePath)
		if err != nil {
			return err
		}
		if !exists {
			logger.Debug("daily note does not exist, skipping", "file", notePath)
			continue
		}

		scanned++
		processor := obsidianutils.NewSimpleFrontmatterProcessor(notePath)
		value, getErr := processor.GetValue(frontmatterKey)
		present := getErr == nil

		isZero, ok := internal.IsZeroValue(value, present)
		if !ok {
			logger.Warn("could not interpret value, skipping", "file", notePath, "value", value)
			continue
		}

		if !isZero {
			memoized = value
			haveMemoized = true
			logger.Debug("memorized value", "file", notePath, "value", value)
			continue
		}

		if !haveMemoized {
			logger.Debug("zero value and nothing memorized yet, leaving as is", "file", notePath)
			continue
		}

		if !t.Before(today) {
			logger.Debug("skipping today's or a future daily note", "file", notePath)
			continue
		}

		logger.Debug("backfilling value", "file", notePath, "old", value, "new", memoized)
		if dryRun {
			updated++
			continue
		}

		if err := processor.SetValue(frontmatterKey, memoized); err != nil {
			return err
		}
		doc, err := processor.GenerateMarkDownDocument()
		if err != nil {
			return err
		}
		if err := os.WriteFile(notePath, doc, 0600); err != nil {
			return err
		}
		updated++
	}

	logger.Info("done scanning daily notes", "scanned", scanned, "updated", updated, "dry-run", dryRun)
	return nil
}
