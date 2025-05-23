package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/apognu/gocal"
	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal/meeting"
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
	flag.SetEnvPrefixForFlag("meeting-folder", "OBS_UTIL_AM")
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
	c.SkipBounds = true
	err = c.Parse()
	if err != nil {
		return err
	}
	if len(c.Events) == 0 {
		log.Print("no events found")
		return nil
	}
	if folder == "" {
		return errors.New("-folder must be non empty")
	}
	if meetingFolder == "" {
		return errors.New("-meeting-folder must be non empty")
	}
	folder, err := obsidianutils.ApplyDirectoryPlaceHolder(folder)
	if err != nil {
		return err
	}
	folder = path.Join(folder, meetingFolder)

	if printConfig {
		log.Println(fmt.Sprintf("meeting notes folder: %q", folder))
		log.Println(fmt.Sprintf("noDatePrefix: %t", noDatePrefix))
		return nil
	}

	for _, event := range c.Events {
		if event.Start.Before(time.Now()) {
			log.Printf("skipping event in the past: %s", event.Summary)
			continue
		}
		fullName, err := obsidianutils.CreateFileName(folder, event.Summary, noDatePrefix, *event.Start)
		if err != nil {
			return err
		}
		if _, err := os.Stat(fullName); err == nil {
			log.Printf("skipping existing file: %s", fullName)
			continue
		}
		if dryRun {
			log.Printf("would create meeting with [%s] on [%s] in [%s]", event.Summary, *event.Start, fullName)
			continue
		}
		m, err := meeting.NewMeeting(meeting.WithTitle(event.Summary))
		if err != nil {
			return err
		}
		c, err := m.CreateContent(event.Summary, *event.Start)
		if err != nil {
			return err
		}

		if err = os.WriteFile(fullName, []byte(c), 0600); err != nil {
			return err
		}
		log.Printf("created meeting with [%s] on [%s] in [%s]", event.Summary, *event.Start, fullName)
	}
	return nil
}
