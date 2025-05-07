# Obsidian Frontmatter Editor (obs-fm)

This utility modifies frontmatter in Obsidian notes.

## Description

The Obsidian Frontmatter Editor utility allows you to modify the frontmatter of Obsidian notes. It can set string, integer, or float values for specified keys in the frontmatter. This is particularly useful for scripting or automating changes to note metadata.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-folder` | Base path to Obsidian vault | (required) |
| `-daily-folder` | Where to store the daily note inside the vault | (required for daily notes) |
| `-print-config` | Print configuration | `false` |
| `-note-path` | Path to note | (required) |
| `-note-type` | Type of note (e.g., "daily") | (empty) |
| `-value-type` | Type of value ("string", "int", "float") | `string` |
| `-key` | Key in the frontmatter to modify | (required) |
| `-value` | Value to set for the key | (empty) |

## Usage

### Modify a regular note

```bash
obs-fm -folder /path/to/vault -note-path "path/to/note.md" -key "tags" -value "meeting"
```

### Modify a daily note

```bash
obs-fm -folder /path/to/vault -daily-folder "Daily Notes" -note-type daily -note-path "2023-09-15" -key "mood" -value "happy"
```

### Set an integer value

```bash
obs-fm -folder /path/to/vault -daily-folder "Daily Notes" -note-type daily -note-path "2023-09-15" -key "steps" -value "10000" -value-type int
```

### Set a float value

```bash
obs-fm -folder /path/to/vault -daily-folder "Daily Notes" -note-type daily -note-path "2023-09-15" -key "weight" -value "75.5" -value-type float
```

## Notes

- For daily notes, if `-note-path` is not specified, the current date will be used.
- The `-note-type` flag is used to determine how to process the note path. Currently, only "daily" is supported as a special type.
- When `-note-type` is set to "daily", the note path is expected to be in the format "YYYY-MM-DD".