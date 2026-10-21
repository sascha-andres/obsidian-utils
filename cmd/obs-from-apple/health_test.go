package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleExport = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE HealthData [
<!ELEMENT HealthData (Record*)>
]>
<HealthData locale="de_DE">
 <Record type="HKQuantityTypeIdentifierStepCount" sourceName="Watch" startDate="2026-10-20 08:00:00 +0200" endDate="2026-10-20 08:10:00 +0200" value="1000"/>
 <Record type="HKQuantityTypeIdentifierStepCount" sourceName="Watch" startDate="2026-10-20 18:00:00 +0200" endDate="2026-10-20 18:10:00 +0200" value="500"/>
 <Record type="HKQuantityTypeIdentifierStepCount" sourceName="Phone" startDate="2026-10-20 08:00:00 +0200" endDate="2026-10-20 08:10:00 +0200" value="1200">
  <MetadataEntry key="HKMetadataKeySyncIdentifier" value="x"/>
 </Record>
 <Record type="HKQuantityTypeIdentifierStepCount" sourceName="Phone" startDate="2026-10-21 08:00:00 +0200" endDate="2026-10-21 08:10:00 +0200" value="42"/>
 <Record type="HKQuantityTypeIdentifierHeartRate" sourceName="Watch" startDate="2026-10-20 08:00:00 +0200" endDate="2026-10-20 08:00:00 +0200" value="60"/>
 <Record type="HKCategoryTypeIdentifierSleepAnalysis" sourceName="Watch" startDate="2026-10-20 22:30:00 +0200" endDate="2026-10-20 23:45:00 +0200" value="HKCategoryValueSleepAnalysisAsleepCore"/>
 <Record type="HKCategoryTypeIdentifierSleepAnalysis" sourceName="Watch" startDate="2026-10-20 23:45:00 +0200" endDate="2026-10-21 00:15:00 +0200" value="HKCategoryValueSleepAnalysisAwake"/>
 <Record type="HKCategoryTypeIdentifierSleepAnalysis" sourceName="Watch" startDate="2026-10-21 00:15:00 +0200" endDate="2026-10-21 06:30:00 +0200" value="HKCategoryValueSleepAnalysisAsleepDeep"/>
 <Record type="HKCategoryTypeIdentifierSleepAnalysis" sourceName="Watch" startDate="2026-10-20 22:00:00 +0200" endDate="2026-10-21 07:00:00 +0200" value="HKCategoryValueSleepAnalysisInBed"/>
 <Record type="HKCategoryTypeIdentifierSleepAnalysis" sourceName="Band" startDate="2026-10-20 23:00:00 +0200" endDate="2026-10-21 06:00:00 +0200" value="HKCategoryValueSleepAnalysisAsleepCore"/>
</HealthData>
`

func TestReadHealthData(t *testing.T) {
	data, err := readHealthData(strings.NewReader(sampleExport))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		values dailyValues
		day    string
		want   float64
		wantOk bool
	}{
		{"steps max of sources", data.steps, "2026-10-20", 1500, true},
		{"steps next day", data.steps, "2026-10-21", 42, true},
		{"steps missing day", data.steps, "2026-10-19", 0, false},
		{"sleep on wake-up day without awake and in bed", data.sleep, "2026-10-21", 7.5, true},
		{"no sleep on bed day", data.sleep, "2026-10-20", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.values.max(tt.day)
			if ok != tt.wantOk || got != tt.want {
				t.Fatalf("max(%s) = %v, %v, want %v, %v", tt.day, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestOpenExportZip(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "export.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	for name, content := range map[string]string{
		"apple_health_export/export_cda.xml": "wrong",
		"apple_health_export/export.xml":     sampleExport,
	} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(w, content); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	rc, err := openExport(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != sampleExport {
		t.Fatalf("openExport read %q, want export.xml content", got)
	}
}
