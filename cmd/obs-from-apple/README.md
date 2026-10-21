# Apple Health Importer (obs-from-apple)

This utility writes steps and sleep minutes from an Apple Health export into existing daily notes.

## Description

The utility reads the "Export All Health Data" file of the Apple Health app. It accepts the `export.zip` file or the extracted `export.xml` file. It aggregates the values per day and writes them to the frontmatter of the daily note of that day.

The utility only writes a value if the key is missing or `0`. It does not change other values. It does not create daily notes and it does not change today's note, because the data of today is not complete.

- **Steps**: `HKQuantityTypeIdentifierStepCount`. The steps are summed per source (for example iPhone, Apple Watch). The source with the highest total is used, because the sources record the same steps. Stored as an integer.
- **Sleep minutes**: `HKCategoryTypeIdentifierSleepAnalysis` with the stages Core, Deep, REM and Unspecified. InBed and Awake do not count. The minutes are summed per source, and time that segments of one source share counts one time only. The source with the highest total is used. A night belongs to the wake-up day (sleep that ends between 18:00 of the previous day and 18:00 counts for that day). Stored as an integer (minutes).

The Apple Health export does not contain a sleep score.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-folder` | Base path to Obsidian vault | (required) |
| `-daily-folder` | Where the daily notes are stored inside the vault | (required) |
| `-export` | Path to the Apple Health export (`export.zip` or `export.xml`) | (required) |
| `-start-date` | Lowest date to import (2006-01-02) | (no limit) |
| `-until` | Highest date to import (2006-01-02) | yesterday |
| `-steps-key` | Frontmatter key for steps | `steps` |
| `-sleep-key` | Frontmatter key for sleep minutes | `sleep` |
| `-dry-run` | Do not write files, only log what would change | `false` |
| `-print-config` | Print configuration | `false` |
| `-log-level` | Log level (debug, info, warn) | `info` |

Each flag can also be set with an environment variable with the prefix `OBS_UTIL_FROM_APPLE_`, for example `OBS_UTIL_FROM_APPLE_EXPORT`. `-folder`, `-daily-folder`, `-print-config` and `-log-level` use the common prefix `OBS_UTIL_`.

## Usage

```bash
obs-from-apple -folder /path/to/vault -daily-folder "Daily Notes" -export ~/Downloads/export.zip -dry-run -log-level debug
obs-from-apple -folder /path/to/vault -daily-folder "Daily Notes" -export ~/Downloads/export.zip -start-date 2026-10-01
```
