package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// togglDateTimeLayout is the layout of a date column and a time column joined with a space.
const togglDateTimeLayout = time.DateOnly + " " + time.TimeOnly

// readWorkTime reads a toggl detailed csv report and sums the time per day (2006-01-02).
// work contains the time of entries of a project in projects or a client in clients, other the time of all other entries.
// If projects and clients are empty, all entries count as work. An entry that crosses midnight is split between the days.
func readWorkTime(r io.Reader, projects, clients map[string]bool) (work, other map[string]time.Duration, err error) {
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
	required := []string{"Project", "Start date", "Start time", "Stop date", "Stop time"}
	if len(clients) > 0 {
		required = append(required, "Client")
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
		result := work
		if len(projects)+len(clients) > 0 && !projects[record[columns["Project"]]] &&
			(len(clients) == 0 || !clients[record[columns["Client"]]]) {
			result = other
		}
		// the report has no time zone, toggl uses the time zone of the user
		start, err := time.ParseInLocation(togglDateTimeLayout, record[columns["Start date"]]+" "+record[columns["Start time"]], time.Local)
		if err != nil {
			return nil, nil, err
		}
		stop, err := time.ParseInLocation(togglDateTimeLayout, record[columns["Stop date"]]+" "+record[columns["Stop time"]], time.Local)
		if err != nil {
			return nil, nil, err
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
