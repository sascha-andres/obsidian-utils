package main

import (
	"errors"
	"log/slog"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
)

var (
	key, value, dailyFolder, logLevel     string
	notePath, noteType, folder, valueType string
	printConfig                           bool
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_FM")
	flag.StringVar(&logLevel, "log-level", "info", "pass log level (debug/info/warn/error)")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.StringVar(&notePath, "note-path", "", "path to note")
	flag.StringVar(&noteType, "note-type", "", "type of note")
	flag.StringVar(&valueType, "value-type", "string", "type of value (string, bool, int, float)")
	flag.StringVar(&key, "key", "", "key")
	flag.StringVar(&value, "value", "", "value")
}

func main() {
	flag.Parse()
	logger := internal.CreateLogger(logLevel, "OBS_UTIL_FM")
	if err := run(logger); err != nil {
		logger.Error("could not execute utility", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
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
			logger.Error("if note type is daily, -note-path must be non empty and have format 2006-01-02 or be empty")
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
	logger.Debug("working on", "file", completePath)
	processor := obsidianutils.NewSimpleFrontmatterProcessor(completePath)
	switch valueType {
	case "int":
		intValue, err := strconv.Atoi(value)
		if err != nil {
			logger.Error("could not convert to int", "input", value)
			return err
		}
		err = processor.SetValue(key, intValue)
		break
	case "float":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Error("could not convert to float", "input", value)
			return err
		}
		err = processor.SetValue(key, floatValue)
		break
	case "bool":
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			logger.Error("could not convert to bool", "input", value)
			return err
		}
		err = processor.SetValue(key, boolValue)
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
	logger.Info("done working", "file", completePath)
	return os.WriteFile(completePath, doc, 0600)
}
