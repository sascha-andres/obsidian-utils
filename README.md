# Add a meeting note
Adds a meeting note to specified folder

# Configure

Each flag has an env variable prefixed with OBS_AM_. The env variable
contains uppercase letters and underscores instead of minus.

## Folder

Flag `-folder` specifies the folder path to store generated MD file(s)

## No date prefix

Flag `-no-date-prefix` instructs obsidian am to not add a YYYYMMDD prefix to the filename

## Recurring

To create a recurring event use the `-recurring` flag.

## Interval

For a recurring event `-interval` specifies the recurring interval. Valid values are:

- daily
- weekly
- bi-weekly

## Times

The `-times` flag instructs obsidian-am to create n times the event with respect to interval given.

## Config

To print the config used the `-print-config` flag