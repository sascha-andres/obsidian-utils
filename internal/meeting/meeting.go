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

// meetingTemplate is the template of a meeting note.
var meetingTemplate = template.Must(template.New("m").Parse(`---
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

## Notes`))

// Content generates the content of a meeting note with title on appointment.
// The title is trimmed, line breaks become "->" and colons become "-". An empty title is an error.
func Content(title string, appointment time.Time) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", errors.New("title is required")
	}
	title = strings.ReplaceAll(title, "\n", "->")
	title = strings.ReplaceAll(title, ":", "-")
	td := TemplateData{time.Now().Format(time.RFC850), appointment.Format(time.RFC3339), title, appointment.Format("2006-01-02")}
	var tpl bytes.Buffer
	err := meetingTemplate.Execute(&tpl, td)
	return tpl.String(), err
}
