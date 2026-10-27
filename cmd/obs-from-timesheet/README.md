# Timesheet Importer (obs-from-timesheet)

This utility writes the work time from a time report into existing daily notes. Supported reports:

- `toggl`: Toggl Track detailed report.
- `timenote`: timenote report.

## Description

The utility reads the CSV file of the report. It sums the time of the entries per day and writes the result in minutes to the frontmatter of the daily note of that day.

The utility only writes a value if the key is missing or `0`. It does not change other values. It does not create daily notes and it does not change today's note, because the data of today is not complete.

- An entry counts if its project is in `-project` or its client is in `-client`. If both are empty, all entries count.
- The report must have a `Client` column if `-client` is set.
- If `-other-time-key` is set, the time of all other entries goes to that key. The same rule applies: only a missing or `0` value is written.
- An entry that crosses midnight is split between the two days. A timenote report has no stop date: if the stop time is before the start time, the entry stops on the next day.
- An entry without stop time is still running and is skipped.
- The time is calculated from the start and stop columns in the local time zone. Use the same time zone as in the source app.
- The value is rounded to full minutes per day.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-folder` | Base path to Obsidian vault | (required) |
| `-daily-folder` | Where the daily notes are stored inside the vault | (required) |
| `-report` | Path to the report (CSV) | (required) |
| `-source-type` | Type of the report: `toggl` or `timenote` | (ask) |
| `-project` | Comma separated list of projects to include | (empty) |
| `-client` | Comma separated list of clients to include | (empty) |
| `-start-date` | Lowest date to import (2006-01-02) | (no limit) |
| `-until` | Highest date to import (2006-01-02) | yesterday |
| `-work-time-key` | Frontmatter key for work time in minutes | `work time` |
| `-other-time-key` | Frontmatter key for the time in minutes of entries that do not count as work. Each day in the report gets the key, also with `0`. | `other time` (empty: not written) |
| `-dry-run` | Do not write files, only log what would change | `false` |
| `-print-config` | Print configuration | `false` |
| `-log-level` | Log level (debug, info, warn) | `info` |

Each flag can also be set with an environment variable with the prefix `OBS_UTIL_FROM_TIMESHEET_`, for example `OBS_UTIL_FROM_TIMESHEET_REPORT`. `-folder`, `-daily-folder`, `-print-config` and `-log-level` use the common prefix `OBS_UTIL_`.

## Usage

```bash
obs-from-timesheet -folder /path/to/vault -daily-folder "Daily Notes" -report ~/Downloads/TogglTrack_Report_Detailed_report.csv -source-type toggl -dry-run -log-level debug
obs-from-timesheet -folder /path/to/vault -daily-folder "Daily Notes" -report ~/Downloads/TogglTrack_Report_Detailed_report.csv -source-type toggl -project "IQ,Other"
obs-from-timesheet -folder /path/to/vault -daily-folder "Daily Notes" -report ~/Downloads/timenote-2026-11-02-2026-11-08.csv -source-type timenote -client "Infront Quant AG"
```
