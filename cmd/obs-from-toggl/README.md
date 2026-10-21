# Toggl Importer (obs-from-toggl)

This utility writes the work time from a Toggl Track detailed report into existing daily notes.

## Description

The utility reads the CSV file of a Toggl Track detailed report. It sums the time of the entries per day and writes the result in minutes to the frontmatter of the daily note of that day.

The utility only writes a value if the key is missing or `0`. It does not change other values. It does not create daily notes and it does not change today's note, because the data of today is not complete.

- Only entries of the projects in `-project` count. If `-project` is empty, all entries count.
- If `-other-time-key` is set, the time of all other entries goes to that key. The same rule applies: only a missing or `0` value is written.
- An entry that crosses midnight is split between the two days.
- The time is calculated from the start and stop columns in the local time zone. Use the same time zone as in Toggl.
- The value is rounded to full minutes per day.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-folder` | Base path to Obsidian vault | (required) |
| `-daily-folder` | Where the daily notes are stored inside the vault | (required) |
| `-report` | Path to the Toggl detailed report (CSV) | (required) |
| `-project` | Comma separated list of projects to include, empty for all | `IQ` |
| `-start-date` | Lowest date to import (2006-01-02) | (no limit) |
| `-until` | Highest date to import (2006-01-02) | yesterday |
| `-work-time-key` | Frontmatter key for work time in minutes | `work time` |
| `-other-time-key` | Frontmatter key for the time in minutes of projects not in `-project`. Each day in the report gets the key, also with `0`. | `other time` (empty: not written) |
| `-dry-run` | Do not write files, only log what would change | `false` |
| `-print-config` | Print configuration | `false` |
| `-log-level` | Log level (debug, info, warn) | `info` |

Each flag can also be set with an environment variable with the prefix `OBS_UTIL_FROM_TOGGL_`, for example `OBS_UTIL_FROM_TOGGL_REPORT`. `-folder`, `-daily-folder`, `-print-config` and `-log-level` use the common prefix `OBS_UTIL_`.

## Usage

```bash
obs-from-toggl -folder /path/to/vault -daily-folder "Daily Notes" -report ~/Downloads/TogglTrack_Report_Detailed_report.csv -dry-run -log-level debug
obs-from-toggl -folder /path/to/vault -daily-folder "Daily Notes" -report ~/Downloads/TogglTrack_Report_Detailed_report.csv -project "IQ,Other"
```
