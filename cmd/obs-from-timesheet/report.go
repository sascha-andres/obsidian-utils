package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// reportFormat names the columns of a report and the layout of a date column and a time column joined with a space.
// If stopDate is empty, an entry stops on its start date, or on the next day if the stop time is before the start time.
type reportFormat struct {
	client, project, startDate, startTime, stopDate, stopTime, layout string
}

// formats contains the report format per source type.
var formats = map[string]reportFormat{
	"toggl":    {"Client", "Project", "Start date", "Start time", "Stop date", "Stop time", time.DateOnly + " " + time.TimeOnly},
	"timenote": {"client", "project", "date", "start", "", "stop", time.DateOnly + " 15:04"},
}

// readWorkTime reads a csv report in the given format and sums the time per day (2006-01-02).
// work contains the time of entries of a project in projects or a client in clients, other the time of all other entries.
// If projects and clients are empty, all entries count as work. An entry that crosses midnight is split between the days.
// An entry without stop time is still running and is skipped.
func readWorkTime(r io.Reader, format reportFormat, projects, clients map[string]bool) (work, other map[string]time.Duration, err error) {
	// toggl writes a byte order mark, which encoding/csv does not accept before a quote
	buffered := bufio.NewReader(r)
	if bom, err := buffered.Peek(3); err == nil && string(bom) == "\ufeff" {
		_, _ = buffered.Discard(3)
	}
	reader := csv.NewReader(buffered)
	header, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}
	columns := map[string]int{}
	for i, name := range header {
		columns[name] = i
	}
	required := []string{format.project, format.startDate, format.startTime, format.stopTime}
	if format.stopDate != "" {
		required = append(required, format.stopDate)
	}
	if len(clients) > 0 {
		required = append(required, format.client)
	}
	for _, name := range required {
		if _, ok := columns[name]; !ok {
			return nil, nil, fmt.Errorf("column %q not in report", name)
		}
	}

	work, other = map[string]time.Duration{}, map[string]time.Duration{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return work, other, nil
		}
		if err != nil {
			return nil, nil, err
		}
		if record[columns[format.stopTime]] == "" {
			continue
		}
		result := work
		if len(projects)+len(clients) > 0 && !projects[record[columns[format.project]]] &&
			(len(clients) == 0 || !clients[record[columns[format.client]]]) {
			result = other
		}
		// the report has no time zone, the source uses the time zone of the user
		start, err := time.ParseInLocation(format.layout, record[columns[format.startDate]]+" "+record[columns[format.startTime]], time.Local)
		if err != nil {
			return nil, nil, err
		}
		stopDate := record[columns[format.startDate]]
		if format.stopDate != "" {
			stopDate = record[columns[format.stopDate]]
		}
		stop, err := time.ParseInLocation(format.layout, stopDate+" "+record[columns[format.stopTime]], time.Local)
		if err != nil {
			return nil, nil, err
		}
		if format.stopDate == "" && stop.Before(start) {
			stop = stop.AddDate(0, 0, 1)
		}
		if stop.Before(start) {
			return nil, nil, fmt.Errorf("entry stops before it starts: %s", start)
		}
		for start.Before(stop) {
			nextDay := time.Date(start.Year(), start.Month(), start.Day()+1, 0, 0, 0, 0, start.Location())
			end := stop
			if nextDay.Before(stop) {
				end = nextDay
			}
			result[start.Format(time.DateOnly)] += end.Sub(start)
			start = end
		}
	}
}
