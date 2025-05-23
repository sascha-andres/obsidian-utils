# iCal Importer (ical)

This utility creates meeting notes in Obsidian from an iCal file.

## Description

The iCal Importer utility allows you to import calendar events from an iCal file into your Obsidian vault as meeting notes. It reads the events from the specified iCal file, filters out past events, and creates a markdown file for each future event using a predefined template. This is useful for automatically creating meeting notes for upcoming calendar events.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-folder` | Base path of Obsidian vault | (required) |
| `-meeting-folder` | Where to store the meeting notes | (required) |
| `-no-date-prefix` | Pass to not add yyyy-mm-dd prefix to filename | `false` |
| `-ical-file` | Path to the iCal file or "-" for stdin | (required) |
| `-dry-run` | Pass to not create files (preview only) | `false` |
| `-print-config` | Print configuration | `false` |

## Usage

```bash
ical -folder /path/to/vault -meeting-folder "Meetings" -ical-file calendar.ics
```

This will:
1. Read events from the calendar.ics file
2. Skip any events that are in the past
3. Create a meeting note for each future event
4. Skip creation for events that already have a corresponding file

To preview what would be created without actually creating files:

```bash
ical -folder /path/to/vault -meeting-folder "Meetings" -ical-file calendar.ics -dry-run
```

To read from stdin:

```bash
cat calendar.ics | ical -folder /path/to/vault -meeting-folder "Meetings" -ical-file -
```

## Template

The meeting note template includes:
- Frontmatter with date created, date modified, tags, aliases, date, and title
- Link to the daily note for the meeting date
- Sections for attendees and notes