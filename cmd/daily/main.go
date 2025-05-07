package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
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
	}
)

var (
	folder, forDate, dailyFolder, templateFile string
	printConfig, overwrite                     bool
)

//go:embed DNote.md
var defaultTemplateFile string

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	obsidianutils.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_DAILY")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&templateFile, "template-file", "", "path to template file")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing file")
	flag.StringVar(&forDate, "for-date", time.Now().Format(time.DateOnly), "date for which to create the daily note (2006-01-02)")
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
	log.SetPrefix("[OBS_UTIL_DAILY] ")
}

// main is the entry point of the program.
func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run initializes and executes the daily note creation process based on specified folder paths, date, and template configuration.
func run() error {
	log.Print("start creating a daily note")

	if folder == "" {
		return errors.New("-folder must be non empty")
	}
	folder, err := obsidianutils.ApplyDirectoryPlaceHolder(folder)
	if err != nil {
		return err
	}
	if dailyFolder == "" {
		return errors.New("-daily-folder must be non empty")
	}
	if forDate == "" {
		forDate = time.Now().Format(time.DateOnly)
	}

	folder = path.Join(folder, dailyFolder)

	if printConfig {
		log.Println(fmt.Sprintf("daily notes folder: %q", folder))
		log.Println(fmt.Sprintf("for-date: %q", forDate))
		return nil
	}

	t, err := time.Parse("2006-01-02", forDate)
	if err != nil {
		return err
	}

	resultingDirectory := path.Join(folder, t.Format("2006/01"))
	resultingFile := path.Join(folder, fmt.Sprintf("%s.md", t.Format("2006/01/2006-01-02")))

	_ = os.MkdirAll(resultingDirectory, 0700)

	if overwrite {
		log.Printf("overwriting existing file %q", resultingFile)
		_ = os.Remove(resultingFile)
	}

	if _, err := os.Stat(resultingFile); err == nil {
		log.Printf("file %q already exists", resultingFile)
		return nil
	}

	log.Printf("creating file %q", resultingFile)

	tpl, err := executeTemplate(t, err)
	if err != nil {
		return err
	}

	if err := os.WriteFile(resultingFile, tpl.Bytes(), 0600); err != nil {
		return err
	}
	log.Printf("file %q created", resultingFile)

	return nil
}

// executeTemplate generates a template using the provided date and error, returning the rendered output or an error.
func executeTemplate(t time.Time, err error) (bytes.Buffer, error) {
	templateData := NoteData{}
	templateData.DailyNoteFolder = dailyFolder
	templateData.Current.Year = t.Format("2006")
	templateData.Current.Month = t.Format("01")
	templateData.Current.Day = t.Format("02")
	templateData.Current.DateOnly = t.Format("2006-01-02")
	templateData.Previous.Year = t.AddDate(0, 0, -1).Format("2006")
	templateData.Previous.Month = t.AddDate(0, 0, -1).Format("01")
	templateData.Previous.Day = t.AddDate(0, 0, -1).Format("02")
	templateData.Previous.DateOnly = t.AddDate(0, 0, -1).Format("2006-01-02")
	templateData.Next.Year = t.AddDate(0, 0, 1).Format("2006")
	templateData.Next.Month = t.AddDate(0, 0, 1).Format("01")
	templateData.Next.Day = t.AddDate(0, 0, 1).Format("02")
	templateData.Next.DateOnly = t.AddDate(0, 0, 1).Format("2006-01-02")
	templateContent := ""
	if templateFile != "" {
		data, err := os.ReadFile(templateFile)
		if err != nil {
			log.Printf("could not read template file %q: %v", templateFile, err)
			return bytes.Buffer{}, err
		}
		templateContent = string(data)
	}
	if templateContent == "" {
		templateContent = defaultTemplateFile
	}
	templateEngine, err := template.New("daily").Parse(templateContent)
	if err != nil {
		log.Printf("could not parse template: %v", err)
		return bytes.Buffer{}, err
	}
	var tpl bytes.Buffer
	err = templateEngine.Execute(&tpl, templateData)
	if err != nil {
		log.Printf("could not execute template: %v", err)
		return bytes.Buffer{}, err
	}
	return tpl, nil
}
