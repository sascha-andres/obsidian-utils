package main

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/absidian-utils"
)

var (
	key, value, dailyFolder    string
	notePath, noteType, folder string
	printConfig                bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	flag.SetEnvPrefix("OBS_UTIL")
	flag.SetEnvPrefixForFlag("note-path", "OBS_UTIL_FM")
	flag.SetEnvPrefixForFlag("note-type", "OBS_UTIL_FM")
	flag.SetEnvPrefixForFlag("key", "OBS_UTIL_FM")
	flag.SetEnvPrefixForFlag("value", "OBS_UTIL_FM")
	// TODO commonly used flags in utilities should be declared in a shared func
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.StringVar(&notePath, "note-path", "", "path to note")
	flag.StringVar(&noteType, "note-type", "", "type of note")
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
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	folder = strings.Replace(folder, "$$PWD$$", currentDir, -1)
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
	err = processor.SetValue(key, value)
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
