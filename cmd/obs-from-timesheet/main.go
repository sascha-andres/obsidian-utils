package main

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"math"
	"os"
	"slices"
	"strings"

	"go.livingit.de/reuse/flag"

	"github.com/sascha-andres/obsidian-utils/internal"
)

var (
	folder, dailyFolder, logLevel, reportPath                              string
	startDateFlag, untilFlag, workTimeKey, projects, clients, otherTimeKey string
	sourceType                                                             string
	printConfig, dryRun                                                    bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_FROM_TIMESHEET")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where the daily notes are stored inside the vault")
	flag.StringVar(&reportPath, "report", "", "path to time report (csv)")
	flag.StringVar(&sourceType, "source-type", "", "type of the report (toggl, timenote), empty to ask")
	flag.StringVar(&projects, "project", "", "comma separated list of projects to include")
	flag.StringVar(&clients, "client", "", "comma separated list of clients to include")
	flag.StringVar(&startDateFlag, "start-date", "", "lowest date to import (2006-01-02), empty for no limit")
	flag.StringVar(&untilFlag, "until", "", "highest date to import (2006-01-02), empty for yesterday")
	flag.StringVar(&workTimeKey, "work-time-key", "work time", "frontmatter key for work time in minutes")
	flag.StringVar(&otherTimeKey, "other-time-key", "other time", "frontmatter key for time in minutes of entries that do not count as work, empty to skip")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&dryRun, "dry-run", false, "do not write files, only log what would change")
}

// main is the entry point of the program.
func main() {
	flag.Parse()
	if printConfig {
		internal.PrintFlags()
	}

	logger := internal.CreateLogger(logLevel, "OBS_UTIL_FROM_TIMESHEET")
	if err := run(logger); err != nil {
		logger.Error("error running obs-from-timesheet", "err", err)
		os.Exit(1)
	}
}

// run reads the time report and writes the work time to daily notes that have no value yet.
func run(logger *slog.Logger) error {
	dailyNoteFolder, err := internal.ResolveFolder(folder, dailyFolder, "daily-folder")
	if err != nil {
		return err
	}
	if reportPath == "" {
		return errors.New("-report must be non empty")
	}
	if workTimeKey == "" {
		return errors.New("-work-time-key must be non empty")
	}
	lastDay, err := internal.LastImportDay(startDateFlag, untilFlag)
	if err != nil {
		return err
	}
	if sourceType == "" {
		if sourceType, err = internal.PromptSelect("source type", slices.Sorted(maps.Keys(formats))); err != nil {
			return err
		}
	}
	format, ok := formats[sourceType]
	if !ok {
		return fmt.Errorf("-source-type %q not supported, use one of %v", sourceType, slices.Sorted(maps.Keys(formats)))
	}
	includedProjects, includedClients := splitList(projects), splitList(clients)

	if printConfig {
		logger.Info("daily-folder", "daily-folder", dailyNoteFolder)
		logger.Info("report", "report", reportPath, "source-type", sourceType)
		logger.Info("projects", "projects", slices.Sorted(maps.Keys(includedProjects)))
		logger.Info("clients", "clients", slices.Sorted(maps.Keys(includedClients)))
		logger.Info("range", "start-date", startDateFlag, "until", lastDay)
		logger.Info("keys", "work-time-key", workTimeKey, "other-time-key", otherTimeKey)
		return nil
	}

	report, err := os.Open(reportPath)
	if err != nil {
		return err
	}
	defer report.Close()
	workTime, otherTime, err := readWorkTime(report, format, includedProjects, includedClients)
	if err != nil {
		return err
	}
	values := map[string]map[string]any{}
	for day, work := range workTime {
		values[day] = map[string]any{workTimeKey: int(math.Round(work.Minutes()))}
	}
	if otherTimeKey != "" {
		for day := range otherTime {
			if values[day] == nil {
				values[day] = map[string]any{}
			}
		}
		for day, dayValues := range values {
			dayValues[otherTimeKey] = int(math.Round(otherTime[day].Minutes()))
		}
	}

	scanned, updated, err := internal.FillDailyNotes(logger, dailyNoteFolder, values, startDateFlag, lastDay, dryRun)
	if err != nil {
		return err
	}

	logger.Info("done importing work time", "scanned", scanned, "updated", updated, "dry-run", dryRun)
	return nil
}

// splitList returns the non empty, trimmed values of a comma separated list as a set.
func splitList(list string) map[string]bool {
	result := map[string]bool{}
	for _, value := range strings.Split(list, ",") {
		if value = strings.TrimSpace(value); value != "" {
			result[value] = true
		}
	}
	return result
}
