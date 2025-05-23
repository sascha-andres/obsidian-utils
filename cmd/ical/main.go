package main

import (
	"errors"
	"log"
	"os"

	"github.com/apognu/gocal"
	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
)

var (
	folder, meetingFolder, title, icalFile string
	noDatePrefix, printConfig, dryRun      bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	obsidianutils.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_ICAL")
	obsidianutils.AddCommonFlagPrefixes()
	flag.StringVar(&folder, "folder", "", "base path of obsidian vault")
	flag.StringVar(&meetingFolder, "meeting-folder", "", "where to store the meeting notes")
	flag.BoolVar(&noDatePrefix, "no-date-prefix", false, "pass to not add yyyy-mm-dd prefix to filename")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&dryRun, "dry-run", false, "pass to not create files")
	flag.StringVar(&title, "title", "", "pass title")
	flag.StringVar(&icalFile, "ical-file", "", "pass ical file")
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
	log.SetPrefix("[OBS_UTIL_ICAL] ")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Print("start creating meeting notes from iCal file")
	if icalFile == "" {
		return errors.New("-ical-file must be non empty")
	}
	var (
		f   *os.File
		err error
	)
	if icalFile == "-" {
		f = os.Stdin
	} else {
		f, err = os.OpenFile(icalFile, os.O_RDONLY, 0600)
	}
	if err != nil {
		return err
	}
	c := gocal.NewParser(f)
	err = c.Parse()
	if err != nil {
		return err
	}
	if len(c.Events) == 0 {
		log.Print("no events found")
		return nil
	}
	for _, event := range c.Events {
		log.Print(event.Summary)
	}
	return errors.New("not implemented")
}
