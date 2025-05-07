# Appointment Manager (am)

This utility creates meeting notes in Obsidian with a specific template.

## Description

The Appointment Manager utility allows you to create meeting notes in your Obsidian vault. It prompts you for a date, time, and title for the meeting, and then creates a markdown file with a predefined template. The utility can also create recurring meeting notes with different intervals.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-folder` | Base path of Obsidian vault | (required) |
| `-meeting-folder` | Where to store the meeting notes | (required) |
| `-no-date-prefix` | Pass to not add yyyy-mm-dd prefix to filename | `false` |
| `-recurring` | Pass to create recurring meeting notes | `false` |
| `-interval` | Pass interval size (daily/weekly/bi-weekly) | `daily` |
| `-times` | Pass number of times to create meeting notes | `1` |
| `-print-config` | Print configuration | `false` |

## Usage

```bash
am -folder /path/to/vault -meeting-folder "Meetings"
```

This will prompt you for:
1. Date and time of the meeting (format: yyyy-MM-dd HH:mm)
2. Title of the meeting

For recurring meetings:

```bash
am -folder /path/to/vault -meeting-folder "Meetings" -recurring -interval weekly -times 10
```

This will create 10 weekly meeting notes starting from the date you provide.

## Template

The meeting note template includes:
- Frontmatter with date created, date modified, tags, aliases, date, and title
- Link to the daily note for the meeting date
- Sections for attendees and notes