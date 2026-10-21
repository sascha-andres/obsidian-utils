package main

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	recordTypeSteps = "HKQuantityTypeIdentifierStepCount"
	recordTypeSleep = "HKCategoryTypeIdentifierSleepAnalysis"

	// healthDateLayout is the layout of startDate and endDate in export.xml.
	healthDateLayout = "2006-01-02 15:04:05 -0700"

	// sleepDayShift moves the end of a sleep segment so that a night (18:00 to 18:00) belongs to the wake-up day.
	sleepDayShift = 6 * time.Hour
)

// asleepValues are the sleep analysis values that count as sleep. InBed and Awake do not count.
var asleepValues = map[string]bool{
	"HKCategoryValueSleepAnalysisAsleep":            true,
	"HKCategoryValueSleepAnalysisAsleepUnspecified": true,
	"HKCategoryValueSleepAnalysisAsleepCore":        true,
	"HKCategoryValueSleepAnalysisAsleepDeep":        true,
	"HKCategoryValueSleepAnalysisAsleepREM":         true,
}

// dailyValues holds a value per day (2006-01-02) per source.
type dailyValues map[string]map[string]float64

// add adds value to the total of source on day.
func (d dailyValues) add(day, source string, value float64) {
	if d[day] == nil {
		d[day] = map[string]float64{}
	}
	d[day][source] += value
}

// max returns the highest total of all sources on day. Sources record the same activity, so a sum counts it more than one time.
func (d dailyValues) max(day string) (float64, bool) {
	sources, ok := d[day]
	if !ok {
		return 0, false
	}
	result := 0.0
	for _, v := range sources {
		result = math.Max(result, v)
	}
	return result, true
}

// interval is a time span of a sleep segment.
type interval struct {
	start, end time.Time
}

// sleepSegments holds the sleep segments per day (2006-01-02) per source.
type sleepSegments map[string]map[string][]interval

// add adds a segment of source to day.
func (s sleepSegments) add(day, source string, segment interval) {
	if s[day] == nil {
		s[day] = map[string][]interval{}
	}
	s[day][source] = append(s[day][source], segment)
}

// hours sums the hours per day per source. Time that segments of one source share counts one time only, because a source can contain copies of segments of other sources.
func (s sleepSegments) hours() dailyValues {
	result := dailyValues{}
	for day, sources := range s {
		for source, segments := range sources {
			slices.SortFunc(segments, func(a, b interval) int { return a.start.Compare(b.start) })
			var covered time.Time
			for _, segment := range segments {
				start := segment.start
				if start.Before(covered) {
					start = covered
				}
				if segment.end.After(start) {
					result.add(day, source, segment.end.Sub(start).Hours())
					covered = segment.end
				}
			}
		}
	}
	return result
}

// healthData contains the aggregated daily values read from an export.
type healthData struct {
	steps dailyValues
	sleep dailyValues
}

// openExport opens export.xml, either directly or from inside an export zip.
func openExport(exportPath string) (io.ReadCloser, error) {
	if !strings.EqualFold(path.Ext(exportPath), ".zip") {
		return os.Open(exportPath)
	}
	zr, err := zip.OpenReader(exportPath)
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		if path.Base(f.Name) == "export.xml" {
			rc, err := f.Open()
			if err != nil {
				_ = zr.Close()
				return nil, err
			}
			return zipEntry{ReadCloser: rc, archive: zr}, nil
		}
	}
	_ = zr.Close()
	return nil, errors.New("no export.xml in zip file")
}

// zipEntry closes the zip archive together with the entry.
type zipEntry struct {
	io.ReadCloser
	archive *zip.ReadCloser
}

// Close closes the entry and the archive.
func (z zipEntry) Close() error {
	return errors.Join(z.ReadCloser.Close(), z.archive.Close())
}

// readHealthData streams export.xml and aggregates steps and sleep hours per day.
func readHealthData(r io.Reader) (healthData, error) {
	data := healthData{steps: dailyValues{}}
	sleep := sleepSegments{}
	decoder := xml.NewDecoder(r)
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			data.sleep = sleep.hours()
			return data, nil
		}
		if err != nil {
			return data, err
		}
		element, ok := token.(xml.StartElement)
		if !ok || element.Name.Local != "Record" {
			continue
		}
		attributes := map[string]string{}
		for _, a := range element.Attr {
			attributes[a.Name.Local] = a.Value
		}
		switch attributes["type"] {
		case recordTypeSteps:
			if err := addSteps(data.steps, attributes); err != nil {
				return data, err
			}
		case recordTypeSleep:
			if err := addSleep(sleep, attributes); err != nil {
				return data, err
			}
		}
	}
}

// addSteps adds a step count record to the day it started on.
func addSteps(steps dailyValues, attributes map[string]string) error {
	start, err := time.Parse(healthDateLayout, attributes["startDate"])
	if err != nil {
		return err
	}
	value, err := strconv.ParseFloat(attributes["value"], 64)
	if err != nil {
		return fmt.Errorf("invalid step count %q: %w", attributes["value"], err)
	}
	steps.add(start.Format(time.DateOnly), attributes["sourceName"], value)
	return nil
}

// addSleep adds an asleep record to its wake-up day.
func addSleep(sleep sleepSegments, attributes map[string]string) error {
	if !asleepValues[attributes["value"]] {
		return nil
	}
	start, err := time.Parse(healthDateLayout, attributes["startDate"])
	if err != nil {
		return err
	}
	end, err := time.Parse(healthDateLayout, attributes["endDate"])
	if err != nil {
		return err
	}
	sleep.add(end.Add(sleepDayShift).Format(time.DateOnly), attributes["sourceName"], interval{start: start, end: end})
	return nil
}
