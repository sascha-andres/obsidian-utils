package main

import (
	"errors"
	"log/slog"
	"maps"
	"math"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
)

var (
	folder, dailyFolder, logLevel, reportPath                     string
	startDateFlag, untilFlag, workTimeKey, projects, otherTimeKey string
	printConfig, dryRun                                           bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_FROM_TOGGL")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where the daily notes are stored inside the vault")
	flag.StringVar(&reportPath, "report", "", "path to toggl detailed report (csv)")
	flag.StringVar(&projects, "project", "IQ", "comma separated list of projects to include, empty for all")
	flag.StringVar(&startDateFlag, "start-date", "", "lowest date to import (2006-01-02), empty for no limit")
	flag.StringVar(&untilFlag, "until", "", "highest date to import (2006-01-02), empty for yesterday")
	flag.StringVar(&workTimeKey, "work-time-key", "work time", "frontmatter key for work time in minutes")
	flag.StringVar(&otherTimeKey, "other-time-key", "other time", "frontmatter key for time in minutes of projects not in -project, empty to skip")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&dryRun, "dry-run", false, "do not write files, only log what would change")
}

// main is the entry point of the program.
func main() {
	flag.Parse()
	internal.PrintFlags()

	logger := internal.CreateLogger(logLevel, "OBS_UTIL_FROM_TOGGL")
	if err := run(logger); err != nil {
		logger.Error("error running obs-from-toggl", "err", err)
		os.Exit(1)
	}
}

// run reads the toggl report and writes the work time to daily notes that have no value yet.
func run(logger *slog.Logger) error {
	if folder == "" {
		return errors.New("-folder must be non empty")
	}
	resolvedFolder, err := obsidianutils.ApplyDirectoryPlaceHolder(folder)
	if err != nil {
		return err
	}
	if dailyFolder == "" {
		return errors.New("-daily-folder must be non empty")
	}
	if reportPath == "" {
		return errors.New("-report must be non empty")
	}
	if workTimeKey == "" {
		return errors.New("-work-time-key must be non empty")
	}
	// today's data is not complete, so the last day to import is yesterday
	lastDay := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
	if untilFlag != "" {
		if _, err := time.Parse(time.DateOnly, untilFlag); err != nil {
			return err
		}
		lastDay = min(lastDay, untilFlag)
	}
	if startDateFlag != "" {
		if _, err := time.Parse(time.DateOnly, startDateFlag); err != nil {
			return err
		}
	}
	includedProjects := map[string]bool{}
	for _, project := range strings.Split(projects, ",") {
		if project = strings.TrimSpace(project); project != "" {
			includedProjects[project] = true
		}
	}

	if printConfig {
		logger.Info("folder", "folder", resolvedFolder)
		logger.Info("daily-folder", "daily-folder", dailyFolder)
		logger.Info("report", "report", reportPath)
		logger.Info("projects", "projects", slices.Sorted(maps.Keys(includedProjects)))
		logger.Info("range", "start-date", startDateFlag, "until", lastDay)
		logger.Info("keys", "work-time-key", workTimeKey, "other-time-key", otherTimeKey)
		return nil
	}

	report, err := os.Open(reportPath)
	if err != nil {
		return err
	}
	defer report.Close()
	workTime, otherTime, err := readWorkTime(report, includedProjects)
	if err != nil {
		return err
	}
	days := maps.Clone(workTime)
	if otherTimeKey != "" {
		maps.Copy(days, otherTime)
	}

	var scanned, updated int
	for _, day := range slices.Sorted(maps.Keys(days)) {
		if day < startDateFlag || day > lastDay {
			continue
		}
		values := map[string]any{}
		if work, ok := workTime[day]; ok {
			values[workTimeKey] = int(math.Round(work.Minutes()))
		}
		if otherTimeKey != "" {
			values[otherTimeKey] = int(math.Round(otherTime[day].Minutes()))
		}
		found, changed, err := internal.FillDailyNote(logger, path.Join(resolvedFolder, dailyFolder), day, values, dryRun)
		if err != nil {
			return err
		}
		if found {
			scanned++
		}
		if changed {
			updated++
		}
	}

	logger.Info("done importing toggl work time", "scanned", scanned, "updated", updated, "dry-run", dryRun)
	return nil
}
