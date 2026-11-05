package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/apognu/gocal"
	"go.livingit.de/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
	"github.com/sascha-andres/obsidian-utils/internal/meeting"
)

var (
	folder, meetingFolder, icalFile          string
	logLevel                                 string
	noDatePrefix, printConfig, dryRun, force bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_ICAL")
	flag.SetEnvPrefixForFlag("meeting-folder", "OBS_UTIL_AM")
	flag.StringVar(&logLevel, "log-level", "info", "pass log level (debug/info/warn/error)")
	flag.StringVar(&folder, "folder", "", "base path of obsidian vault")
	flag.StringVar(&meetingFolder, "meeting-folder", "", "where to store the meeting notes")
	flag.BoolVar(&noDatePrefix, "no-date-prefix", false, "pass to not add yyyy-mm-dd prefix to filename")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&dryRun, "dry-run", false, "pass to not create files")
	flag.BoolVar(&force, "force", false, "pass to overwrite existing files")
	flag.StringVar(&icalFile, "ical-file", "", "pass ical file")
}

func main() {
	flag.Parse()
	if printConfig {
		internal.PrintFlags()
	}
	logger := internal.CreateLogger(logLevel, "OBS_UTIL_ICAL")
	if err := run(logger); err != nil {
		logger.Error("error running application", "err", err)
		os.Exit(1)
	}
}

// run parses an iCal file, creates meeting notes in a specified folder, and supports options like dry-run and custom formatting.
func run(logger *slog.Logger) error {
	logger.Info("start creating meeting notes from iCal file")
	if icalFile == "" {
		return errors.New("-ical-file must be non empty")
	}
	f := os.Stdin
	if icalFile != "-" {
		file, err := os.Open(icalFile)
		if err != nil {
			return err
		}
		defer file.Close()
		f = file
	}
	c := gocal.NewParser(f)
	c.SkipBounds = true
	if err := c.Parse(); err != nil {
		return err
	}
	if len(c.Events) == 0 {
		logger.Info("no events found")
		return nil
	}
	folder, err := internal.ResolveFolder(folder, meetingFolder, "meeting-folder")
	if err != nil {
		return err
	}

	if printConfig {
		fmt.Printf("meeting notes folder: %q", folder)
		fmt.Printf("noDatePrefix: %t", noDatePrefix)
		return nil
	}

	for _, event := range c.Events {
		if event.Start.Before(time.Now()) {
			logger.Debug("skipping event in the past", "summary", event.Summary, "start", *event.Start)
			continue
		}
		fullName, err := obsidianutils.CreateFileName(folder, event.Summary, noDatePrefix, *event.Start)
		if err != nil {
			return err
		}
		if _, err := os.Stat(fullName); err == nil && !force {
			logger.Debug("skipping existing file", "file", fullName)
			continue
		}
		if dryRun {
			fmt.Printf("would create meeting with [%s] on [%s] in [%s]\n", event.Summary, *event.Start, fullName)
			continue
		}
		c, err := meeting.Content(event.Summary, *event.Start)
		if err != nil {
			return err
		}

		if err = os.WriteFile(fullName, []byte(c), 0600); err != nil {
			return err
		}
		logger.Info("created meeting", "summary", event.Summary, "start", *event.Start, "file", fullName)
	}
	return nil
}
