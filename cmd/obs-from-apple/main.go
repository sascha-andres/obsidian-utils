package main

import (
	"errors"
	"log/slog"
	"maps"
	"math"
	"os"
	"path"
	"slices"
	"time"

	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
)

var (
	folder, dailyFolder, logLevel, exportPath    string
	startDateFlag, untilFlag, stepsKey, sleepKey string
	printConfig, dryRun                          bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_FROM_APPLE")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where the daily notes are stored inside the vault")
	flag.StringVar(&exportPath, "export", "", "path to apple health export (export.zip or export.xml)")
	flag.StringVar(&startDateFlag, "start-date", "", "lowest date to import (2006-01-02), empty for no limit")
	flag.StringVar(&untilFlag, "until", "", "highest date to import (2006-01-02), empty for yesterday")
	flag.StringVar(&stepsKey, "steps-key", "steps", "frontmatter key for steps")
	flag.StringVar(&sleepKey, "sleep-key", "sleep", "frontmatter key for sleep minutes")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&dryRun, "dry-run", false, "do not write files, only log what would change")
}

// main is the entry point of the program.
func main() {
	flag.Parse()
	internal.PrintFlags()

	logger := internal.CreateLogger(logLevel, "OBS_UTIL_FROM_APPLE")
	if err := run(logger); err != nil {
		logger.Error("error running obs-from-apple", "err", err)
		os.Exit(1)
	}
}

// run reads the apple health export and writes steps and sleep minutes to daily notes that have no value yet.
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
	if exportPath == "" {
		return errors.New("-export must be non empty")
	}
	if stepsKey == "" || sleepKey == "" {
		return errors.New("-steps-key and -sleep-key must be non empty")
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

	if printConfig {
		logger.Info("folder", "folder", resolvedFolder)
		logger.Info("daily-folder", "daily-folder", dailyFolder)
		logger.Info("export", "export", exportPath)
		logger.Info("range", "start-date", startDateFlag, "until", lastDay)
		logger.Info("keys", "steps-key", stepsKey, "sleep-key", sleepKey)
		return nil
	}

	export, err := openExport(exportPath)
	if err != nil {
		return err
	}
	defer export.Close()
	logger.Info("reading export", "export", exportPath)
	data, err := readHealthData(export)
	if err != nil {
		return err
	}

	days := map[string]bool{}
	for day := range data.steps {
		days[day] = true
	}
	for day := range data.sleep {
		days[day] = true
	}

	var scanned, updated int
	for _, day := range slices.Sorted(maps.Keys(days)) {
		if day < startDateFlag || day > lastDay {
			continue
		}
		values := map[string]any{}
		if steps, ok := data.steps.max(day); ok {
			values[stepsKey] = int(math.Round(steps))
		}
		if sleep, ok := data.sleep.max(day); ok {
			values[sleepKey] = int(math.Round(sleep))
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

	logger.Info("done importing apple health data", "scanned", scanned, "updated", updated, "dry-run", dryRun)
	return nil
}
