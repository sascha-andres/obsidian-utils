package main

import (
	"errors"
	"log"

	"github.com/sascha-andres/reuse/flag"
)

var (
	key, value, dailyFolder    string
	notePath, noteType, folder string
	printConfig                bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	flag.SetEnvPrefix("OBS_UTIL")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.StringVar(&notePath, "note-path", "", "path to note")
	flag.StringVar(&noteType, "note-type", "", "type of note")
	flag.StringVar(&key, "key", "", "key")
	flag.StringVar(&value, "value", "", "value")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")

	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
	log.SetPrefix("[OBS_UTIL_FM] ")
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("could not execute utility: %s", err)
	}
}

func run() error {
	if folder == "" {
		return errors.New("-folder must be non empty")
	}
	// here
	return errors.New("not implemented")
}
