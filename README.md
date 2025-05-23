# Obsidian Utils

This repository contains a collection of utilities for working with [Obsidian](https://obsidian.md/), a powerful knowledge base that works on top of a local folder of plain text Markdown files.

## Purpose

The purpose of this repository is to provide command-line utilities that enhance the Obsidian experience by automating common tasks and extending Obsidian's functionality. These utilities are designed to be simple, focused, and easy to use.

## Utilities

### [Appointment Manager (am)](cmd/am/README.md)

The Appointment Manager utility creates meeting notes in Obsidian with a specific template. It allows you to create meeting notes with a predefined template and supports recurring meetings with different intervals.

### [Daily Note Creator (daily)](cmd/daily/README.md)

The Daily Note Creator utility creates daily notes in Obsidian with a specific template. It creates markdown files with a predefined template that includes links to previous and next day notes, as well as various sections for tracking activities, health metrics, tasks, and more.

### [Google Contacts Exporter (ggl)](cmd/ggl/README.md)

The Google Contacts Exporter utility exports Google contacts and contact groups to JSON files. It uses OAuth2 authentication to access your Google account and the Google People API to retrieve your contacts and contact groups.

### [Obsidian Frontmatter Editor (obs-fm)](cmd/obs-fm/README.md)

The Obsidian Frontmatter Editor utility modifies frontmatter in Obsidian notes. It can set string, integer, or float values for specified keys in the frontmatter, which is useful for scripting or automating changes to note metadata.

### [iCal Importer (ical)](cmd/ical/README.md)

The iCal Importer utility creates meeting notes in Obsidian from an iCal file. It reads events from the specified iCal file, filters out past events, and creates a markdown file for each future event using a predefined template.

## Installation

Each utility can be installed separately. Please refer to the individual README files for installation instructions.

## Usage

Please refer to the individual README files for usage instructions for each utility.

## License

This project is licensed under the Creative Commons Attribution-NonCommercial 4.0 International License (CC BY-NC 4.0) - see the LICENSE file for details. This license allows you to obtain, run, and modify the code, but prohibits using it for commercial purposes (selling it for money with or without changes).
