# Google Contacts Exporter (ggl)

This utility exports Google contacts and contact groups to JSON files.

## Description

The Google Contacts Exporter utility allows you to export your Google contacts and contact groups to JSON files. It uses OAuth2 authentication to access your Google account and the Google People API to retrieve your contacts and contact groups.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-state-directory` | Directory to store OAuth2 state | `~/.local/state/ggl` |
| `-output-directory` | Directory to store output files | `.` (current directory) |
| `-print-to-console` | Print data to console instead of writing to files, may be "contacts" or "groups" | (empty) |
| `-verbose` | Enable verbose output | `false` |

## Usage

### Export contacts and groups to files

```bash
ggl -output-directory /path/to/output
```

This will:
1. Authenticate with Google (opening a browser window if needed)
2. Export your contacts to `/path/to/output/contacts.json`
3. Export your contact groups to `/path/to/output/groups.json`

### Print contacts to console

```bash
ggl -print-to-console contacts
```

### Print groups to console

```bash
ggl -print-to-console groups
```

### Enable verbose output

```bash
ggl -verbose -output-directory /path/to/output
```

## Authentication

The first time you run the utility, it will:
1. Open a browser window for you to sign in to your Google account
2. Ask for permission to access your contacts
3. Store the authentication token in the state directory

Subsequent runs will use the stored token, so you won't need to authenticate again unless the token expires or is deleted.

## Data Format

The utility exports the following data:
- Contacts: names, email addresses, phone numbers, addresses, organizations, memberships, birthdays
- Contact groups: names and other group information