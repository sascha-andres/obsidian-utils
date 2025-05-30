package meeting

import (
	"bytes"
	"errors"
	"strings"
	"text/template"
	"time"
)

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

// Meeting represents a structure for managing and handling meeting-related data using a TemplateData instance.
type Meeting struct {
	title string
}

// OptionFunc defines a function type that modifies a Meeting instance or returns an error.
type OptionFunc func(m *Meeting) error

// WithTitle sets the title of the meeting in the TemplateData and returns an OptionFunc for configuring a Meeting.
func WithTitle(title string) OptionFunc {
	return func(m *Meeting) error {
		m.title = title
		return nil
	}
}

// NewMeeting initializes and returns a new Meeting instance with the provided options or an error if an option fails.
func NewMeeting(opts ...OptionFunc) (*Meeting, error) {
	m := &Meeting{}
	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}
	if m.title == "" {
		return nil, errors.New("title is required")
	}
	return m, nil
}

// CreateContent generates the content for a meeting note using a template.
// It takes a title and an appointment time as parameters.
// The template data includes the current time, the appointment date,
// and the title, which are combined and executed using a template engine.
// The resulting content is returned as a string.
// If an error occurs during the parsing or execution of the template,
// an empty string and the error are returned.
func (m *Meeting) CreateContent(title string, appointment time.Time) (string, error) {
	tmpl, err := template.New("m").Parse(meetingTemplate)
	if err != nil {
		return "", err
	}
	err = m.cleanTitle()
	if err != nil {
		return "", err
	}
	td := TemplateData{time.Now().Format(time.RFC850), appointment.Format(time.RFC3339), m.title, appointment.Format("2006-01-02")}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, td)
	return tpl.String(), err
}

func (m *Meeting) cleanTitle() error {
	m.title = strings.TrimSpace(m.title)
	m.title = strings.ReplaceAll(m.title, "\n", "->")
	m.title = strings.ReplaceAll(m.title, ":", "-")
	return nil
}
