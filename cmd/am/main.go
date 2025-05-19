package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/sascha-andres/reuse/flag"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
)

var (
	folder, interval, meetingFolder, dateTime, title string
	recurring, noDatePrefix, printConfig, dryRun     bool
	times                                            int
)

// init initializes the package by setting up flag options, log flags, and prefix.
func init() {
	obsidianutils.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_AM")
	flag.StringVar(&folder, "folder", "", "base path of obsidian vault")
	flag.StringVar(&meetingFolder, "meeting-folder", "", "where to store the meeting notes")
	flag.BoolVar(&noDatePrefix, "no-date-prefix", false, "pass to not add yyyy-mm-dd prefix to filename")
	flag.BoolVar(&recurring, "recurring", false, "pass to create recurring meeting notes")
	flag.StringVar(&interval, "interval", "daily", "pass interval size (daily/weekly/bi-weekly)")
	flag.IntVar(&times, "times", 1, "pass number of times to create meeting notes")
	flag.BoolVar(&printConfig, "print-config", false, "print configuration")
	flag.BoolVar(&dryRun, "dry-run", false, "pass to not create files")
	flag.StringVar(&dateTime, "date-time", "", "pass date and time in format yyyy-mm-dd hh:mm")
	flag.StringVar(&title, "title", "", "pass title")
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
	log.SetPrefix("[OBS_UTIL_AM] ")
}

// main is the entry point of the program.
func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run executes the main logic of the program, which involves creating a meeting note file.
// It prompts the user for input, validates the input, and generates the content of the note.
// The generated note is then saved to a file.
//
// If the "-folder" flag is not specified or is empty, it returns an error.
// The format of the date and time input is expected to be "yyyy-MM-dd hh:mm".
//
// The "title" input is obtained from the user through a prompt.
//
// The "createFileName" function is used to determine the filename for the meeting note.
// It applies specified replacements to the title and prefixes the file with the appointment date if the "noDatePrefix"
// flag is not set.
//
// The "createContent" function generates the content for the meeting note using a template.
// The template data includes the current time, the appointment date, and the title.
//
// The generated content is then saved to the determined filename.
// An error is returned if any of the steps fail.
func run() error {
	log.Print("start creating a meeting note")
	if folder == "" {
		return errors.New("-folder must be non empty")
	}
	folder, err := obsidianutils.ApplyDirectoryPlaceHolder(folder)
	if err != nil {
		return err
	}
	if meetingFolder == "" {
		return errors.New("-meeting-folder must be non empty")
	}

	folder = path.Join(folder, meetingFolder)

	if recurring {
		if interval != "daily" && interval != "weekly" && interval != "bi-weekly" {
			return errors.New("invalid interval")
		}
		if times < 1 {
			return errors.New("invalid times")
		}
		if times == 1 {
			recurring = false
		}
	} else {
		times = 1
		interval = "daily"
	}

	if printConfig {
		log.Println(fmt.Sprintf("meeting notes folder: %q", folder))
		log.Println(fmt.Sprintf("interval: %q", interval))
		log.Println(fmt.Sprintf("recurring: %t", recurring))
		log.Println(fmt.Sprintf("noDatePrefix: %t", noDatePrefix))
		log.Println(fmt.Sprintf("times: %d", times))
		return nil
	}

	if !dryRun {
		err = os.MkdirAll(folder, 0700)
		if err != nil {
			if !os.IsExist(err) {
				return err
			}
		}
	}

	var ts string
	if dateTime != "" {
		_, err = time.Parse("2006-01-02 15:04", dateTime)
		if err != nil {
			return err
		}
		ts = dateTime
	} else {
		ts, err = promptText("provide date and time (2006-01-02 15:04)", time.Now().Format("2006-01-02 15:04"), func(i string) error {
			_, err := time.Parse("2006-01-02 15:04", i)
			return err
		})
		if err != nil {
			return err
		}
	}

	var localTitle string
	if title != "" {
		localTitle = title
	} else {
		localTitle, err = promptText("get title", "", func(s string) error {
			if strings.TrimSpace(s) == "" {
				return errors.New("empty title")
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	t, err := time.Parse("2006-01-02 15:04", ts)
	if err != nil {
		return err
	}

	for i := 0; i < times; i++ {
		if i > 0 {
			if interval == "daily" {
				t = t.Add(time.Hour * 24)
			}
			if interval == "weekly" {
				t = t.Add(time.Hour * 168)
			}
			if interval == "bi-weekly" {
				t = t.Add(time.Hour * 336)
			}
		}
		log.Printf("trying to create meeting with [%s] on [%s]", localTitle, t)

		fullName, err := createFileName(title, t)
		if err != nil {
			return err
		}
		if dryRun {
			log.Printf("would create meeting with [%s] on [%s] in [%s]", localTitle, t, fullName)
		} else {
			c, err := createContent(title, t)
			if err != nil {
				return err
			}

			if err = os.WriteFile(fullName, []byte(c), 0600); err != nil {
				return err
			}
		}
	}

	return nil
}

// createContent generates the content for a meeting note using a template.
// It takes a title and an appointment time as parameters.
// The template data includes the current time, the appointment date,
// and the title, which are combined and executed using a template engine.
// The resulting content is returned as a string.
// If an error occurs during the parsing or execution of the template,
// an empty string and the error are returned.
func createContent(title string, appointment time.Time) (string, error) {
	tmpl, err := template.New("m").Parse(meetingTemplate)
	if err != nil {
		return "", err
	}
	td := TemplateData{time.Now().Format(time.RFC850), appointment.Format(time.RFC3339), title, appointment.Format("2006-01-02")}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, td)
	return tpl.String(), err
}

// TemplateData represents the data necessary for rendering a template. It contains
// fields for the current time, the appointment time, and the title. This data is used
// to generate the content of a meeting note by executing a template.
type TemplateData struct {

	// Now represents the current time as a string. It is a field in the TemplateData struct, which is used to generate the content of a meeting note by executing a template.
	Now string

	// Appointment represents the appointment time in a TemplateData struct. It is used to generate the content of a meeting note by executing a template.
	Appointment string

	// Title represents the title of a meeting. It is a field in the TemplateData struct,
	// which is used to generate the content of a meeting note by executing a template.
	Title string

	// DayNote is a field of struct type TemplateData. It represents a string used to link to the daily note
	DayNote string
}

// The `meetingTemplate` constant is a string that represents a template for generating meeting notes.
// It uses the Go template syntax and incorporates placeholders for various fields such as the current date,
// the appointment date, and the meeting title.
// This template can be used with the `createContent` function to generate the content of a meeting note by
// providing the necessary data such as the title and appointment time.
// The resulting content is returned as a string.
const meetingTemplate = `---
date created: {{ .Now }}
date modified: {{ .Now }}
tags:
  - meeting
aliases: 
date: {{ .Appointment }}
title: {{ .Title }}
---

[[{{ .DayNote }}]]

# Meeting

## Attendees

## Notes`

// replacements is a map containing character replacements for German umlauts and other special characters.
var replacements = map[string]string{
	"ä": "ae",
	"ö": "oe",
	"ü": "ue",
	"Ä": "Ae",
	"Ö": "Oe",
	"Ü": "Ue",
	"ß": "ss",
	":": "",
}

// createFileName generates a file name for a meeting note based on the provided title and appointment time. It applies specified character replacements to the title and prefixes the file with the appointment date if the noDatePrefix flag is not set. The generated file name is returned as a string.
func createFileName(title string, appointment time.Time) (string, error) {
	fixed := title
	for k, v := range replacements {
		fixed = strings.ReplaceAll(fixed, k, v)
	}
	fName := fmt.Sprintf("%s.md", fixed)
	if !noDatePrefix {
		fName = fmt.Sprintf("%s %s", appointment.Format("2006-01-02"), fName)
	}
	return obsidianutils.ApplyDirectoryPlaceHolder(path.Join(folder, fName))
}

// promptText runs a textual prompt
func promptText(label, defaultValue string, val func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	if nil != val {
		prompt.Validate = val
	}
	return prompt.Run()
}
