package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/sascha-andres/reuse/flag"
)

var (
	folder, forDate, dailyFolder, templateFile string
	printConfig                                bool
)

//go:embed DNote.md
var defaultTemplateFile string

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	flag.SetEnvPrefix("OBS_UTIL")
	flag.StringVar(&folder, "folder", "", "base path to obsidian valut")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&templateFile, "template-file", "", "path to template file")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
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

func run() error {
	log.Print("start creating a daily note")

	if folder == "" {
		return errors.New("-folder must be non empty")
	}
	if dailyFolder == "" {
		return errors.New("-daily-folder must be non empty")
	}
	if forDate == "" {
		return errors.New("-for-date must be non empty")
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

	resultingFile := path.Join(folder, fmt.Sprintf("%s.md", t.Format("2006/01/2006-01-02")))

	if _, err := os.Stat(resultingFile); err == nil {
		log.Printf("file %q already exists", resultingFile)
		return nil
	}

	log.Printf("creating file %q", resultingFile)
	if err := os.WriteFile(resultingFile, []byte(""), 0600); err != nil {
		return err
	}

	return nil
}
