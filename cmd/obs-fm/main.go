package main

import (
	"errors"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/absidian-utils"
)

var (
	key, value, dailyFolder               string
	notePath, noteType, folder, valueType string
	printConfig                           bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	obsidianutils.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_FM")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.StringVar(&notePath, "note-path", "", "path to note")
	flag.StringVar(&noteType, "note-type", "", "type of note")
	flag.StringVar(&valueType, "value-type", "string", "type of value")
	flag.StringVar(&key, "key", "", "key")
	flag.StringVar(&value, "value", "", "value")

	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
	log.SetPrefix("[OBS_UTIL_FM] ")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("could not execute utility: %s", err)
	}
}

func run() error {
	var (
		dailyTimestamp = time.Now()
		err            error
	)

	if folder == "" {
		return errors.New("-folder must be non empty")
	}
	folder, err = obsidianutils.ApplyDirectoryPlaceHolder(folder)
	if err != nil {
		return err
	}
	if noteType == "daily" {
		if dailyFolder == "" {
			return errors.New("-daily-folder must be non empty if note-type is daily")
		}
		if notePath == "" {
			notePath = dailyTimestamp.Format("2006-01-02")
		}
		dailyTimestamp, err = time.Parse("2006-01-02", notePath)
		if err != nil {
			log.Printf("if note type is daily, -note-path must be non empty and have format 2006-01-02 or be empty")
			return errors.New("if note type is daily, -note-path must be non empty and have format 2006-01-02 or be empty")
		}
	}
	if notePath == "" {
		return errors.New("-note-path must be non empty")
	}
	if key == "" {
		return errors.New("-key must be non empty")
	}

	completePath := path.Join(folder, dailyFolder, dailyTimestamp.Format("2006/01"), dailyTimestamp.Format("2006-01-02")+".md")
	log.Printf("working on: %q", completePath)
	processor := obsidianutils.NewSimpleFrontmatterProcessor(completePath)
	switch valueType {
	case "int":
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("could not convert %q to int", value)
			return err
		}
		err = processor.SetValue(key, intValue)
		break
	case "float":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("could not convert %q to float", value)
			return err
		}
		err = processor.SetValue(key, floatValue)
		break
	default:
		err = processor.SetValue(key, value)
	}
	if err != nil {
		return err
	}
	doc, err := processor.GenerateMarkDownDocument()
	if err != nil {
		return err
	}
	log.Printf("done working on %q", completePath)
	return os.WriteFile(completePath, doc, 0600)
}
