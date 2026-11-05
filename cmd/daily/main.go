package main

import (
	"bytes"
	_ "embed"
	"errors"
	f "flag"
	"log/slog"
	"os"
	"path"
	"text/template"
	"time"

	"go.livingit.de/reuse/flag"

	"github.com/sascha-andres/obsidian-utils/internal"
)

type (

	// DayData represents the structured data for a specific date including year, month, day, and a composite date-only string.
	DayData struct {
		// Year represents the year portion of a date in a string format.
		Year string

		// Month represents the month portion of a date in a string format.
		Month string

		// Day represents the day portion of a date in a string format.
		Day string

		// DateOnly represents the combination of year, month, and day as a single string formatted as a date.
		DateOnly string
	}

	// NoteData represents the note-related data for a specific date, including links to previous, next, and current day's data.
	// Previous is the DayData for the prior day relative to the current date.
	// Next is the DayData for the following day relative to the current date.
	// Current holds the DayData for the current date.
	// DailyNoteFolder specifies the directory path for storing daily notes.
	NoteData struct {
		// Previous represents the DayData for the prior day relative to the current date.
		Previous DayData

		// Next represents the DayData for the following day relative to the current date.
		Next DayData

		// Current holds the DayData for the current date.
		Current DayData

		// DailyNoteFolder defines the path or location where daily notes are stored as a string.
		// which is basically the -daily-folder parameter
		DailyNoteFolder string

		// WorkLocation specifies the location associated with the current note's context, typically related to the daily note's metadata.
		WorkLocation string
	}
)

var (
	folder, forDate, until, dailyFolder, templateFile string
	defaultWorkLocation                               = "Office"
	logLevel                                          string
	printConfig, overwrite                            bool
)

//go:embed DNote.md
var defaultTemplateFile string

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_DAILY")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&defaultWorkLocation, "default-work-location", defaultWorkLocation, "default work location")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&templateFile, "template-file", "", "path to template file")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing file")
	flag.StringVar(&forDate, "for-date", time.Now().Format(time.DateOnly), "date for which to create the daily note (2006-01-02)")
	flag.StringVar(&until, "until", "", "date until which to create daily notes, starting today (2006-01-02 or +N offset); mutually exclusive with -for-date")
}

// main is the entry point of the program.
func main() {
	flag.Parse()
	if printConfig {
		internal.PrintFlags()
	}

	logger := internal.CreateLogger(logLevel, "OBS_UTIL_DAILY")
	if err := run(logger); err != nil {
		logger.Error("error running daily", "err", err)
		os.Exit(1)
	}
}

// explicitlySet reports whether the named flag was actually passed on the command line.
func explicitlySet(name string) bool {
	set := false
	flag.Visit(func(fl *f.Flag) {
		if fl.Name == name {
			set = true
		}
	})
	return set
}

// run initializes and executes the daily note creation process based on specified folder paths, date, and template configuration.
func run(logger *slog.Logger) error {
	logger.Info("start creating a daily note")

	folder, err := internal.ResolveFolder(folder, dailyFolder, "daily-folder")
	if err != nil {
		return err
	}
	if until != "" && explicitlySet("for-date") {
		return errors.New("-for-date and -until are mutually exclusive")
	}
	if forDate == "" {
		forDate = time.Now().Format(time.DateOnly)
	}
	forDate, err = internal.ResolveDateArg(forDate)
	if err != nil {
		return err
	}

	if until != "" {
		resolvedUntil, err := internal.ResolveDateArg(until)
		if err != nil {
			return err
		}
		untilTime, err := time.Parse("2006-01-02", resolvedUntil)
		if err != nil {
			return err
		}
		today := time.Now().Truncate(24 * time.Hour)
		if untilTime.Before(today) {
			return errors.New("-until must be a date in the future")
		}

		if printConfig {
			logger.Info("daily notes folder", "folder", folder)
			logger.Info("date range", "from", today.Format(time.DateOnly), "until", untilTime.Format(time.DateOnly))
			return nil
		}

		for t := today; !t.After(untilTime); t = t.AddDate(0, 0, 1) {
			if err := createDailyNote(logger, folder, t); err != nil {
				return err
			}
		}
		return nil
	}

	if printConfig {
		logger.Info("daily notes folder", "folder", folder)
		logger.Info("for-date", "for-date", forDate)
		return nil
	}

	t, err := time.Parse("2006-01-02", forDate)
	if err != nil {
		return err
	}

	return createDailyNote(logger, folder, t)
}

// createDailyNote creates the daily note file for the given date inside folder, unless it already exists (or -overwrite is set).
func createDailyNote(logger *slog.Logger, folder string, t time.Time) error {
	resultingFile := internal.DailyNotePath(folder, t)

	_ = os.MkdirAll(path.Dir(resultingFile), 0700)

	if _, err := os.Stat(resultingFile); err == nil {
		if !overwrite {
			logger.Warn("file already exists", "file", resultingFile)
			return nil
		}
		logger.Info("overwriting existing file", "file", resultingFile)
	}

	logger.Info("creating file", "file", resultingFile)

	tpl, err := executeTemplate(logger, t)
	if err != nil {
		return err
	}

	if err := os.WriteFile(resultingFile, tpl.Bytes(), 0600); err != nil {
		return err
	}
	logger.Info("file created", "file", resultingFile)

	return nil
}

// dayData returns the template data of the day t.
func dayData(t time.Time) DayData {
	return DayData{Year: t.Format("2006"), Month: t.Format("01"), Day: t.Format("02"), DateOnly: t.Format(time.DateOnly)}
}

// executeTemplate renders the daily note template for the date t.
func executeTemplate(logger *slog.Logger, t time.Time) (bytes.Buffer, error) {
	templateData := NoteData{
		Previous:        dayData(t.AddDate(0, 0, -1)),
		Next:            dayData(t.AddDate(0, 0, 1)),
		Current:         dayData(t),
		DailyNoteFolder: dailyFolder,
		WorkLocation:    defaultWorkLocation,
	}

	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		templateData.WorkLocation = "Weekend"
	}

	templateContent := ""
	if templateFile != "" {
		data, err := os.ReadFile(templateFile)
		if err != nil {
			logger.Error("could not read template file", "file", templateFile, "err", err)
			return bytes.Buffer{}, err
		}
		templateContent = string(data)
	}
	if templateContent == "" {
		templateContent = defaultTemplateFile
	}
	templateEngine, err := template.New("daily").Parse(templateContent)
	if err != nil {
		logger.Error("could not parse template", "err", err)
		return bytes.Buffer{}, err
	}
	var tpl bytes.Buffer
	err = templateEngine.Execute(&tpl, templateData)
	if err != nil {
		logger.Error("could not execute template", "err", err)
		return bytes.Buffer{}, err
	}
	return tpl, nil
}
