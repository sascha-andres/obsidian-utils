package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/sascha-andres/reuse/flag"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
)

// Scopes required for reading contacts and contact groups
const contactsScope = "https://www.googleapis.com/auth/contacts.readonly"

var (
	stateDirectory  string
	outputDirectory string
	printToConsole  string
	verbose         bool
)

// init initializes the program's environment settings and configuration for Google-related utilities.
// It sets an environment prefix, retrieves the current user, and defines the state directory flag for OAuth2 storage.
func init() {
	flag.SetEnvPrefix("GGL")

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&stateDirectory, "state-directory", path.Join(currentUser.HomeDir, ".local/state/ggl"), "Directory to store OAuth2 state")
	flag.StringVar(&outputDirectory, "output-directory", ".", "Directory to store output files")
	flag.StringVar(&printToConsole, "print-to-console", "", "Print data to console instead of writing to files, may be contacts or groups")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
}

// getClient retrieves a token, saves it, then returns the OAuth2 client
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens
	tokenFile := path.Join(stateDirectory, "token.json")
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		err = saveToken(tokenFile, tok)
		if err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), tok), nil
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("==========================================================")
	fmt.Println("To authorize this application:")
	fmt.Println("1. Copy the following URL")
	fmt.Printf("   %v\n", authURL)
	fmt.Println("2. Open the URL in any browser (can be on another device)")
	fmt.Println("3. Sign in and grant access to your Google account")
	fmt.Println("4. Copy the authorization code provided")
	fmt.Println("5. Paste the authorization code below and press Enter")
	fmt.Println("==========================================================")
	fmt.Print("Enter authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	return config.Exchange(context.TODO(), authCode)
}

// tokenFromFile retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file
func saveToken(path string, token *oauth2.Token) error {
	if verbose {
		fmt.Printf("Saving credential file to: %s\n", path)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// initializeStateDirectory ensures the state directory exists, creating it with the correct permissions if it does not exist.
func initializeStateDirectory() error {
	if _, err := os.Stat(stateDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(stateDirectory, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

// initializeOutputDirectory ensures the output directory exists, creating it with the correct permissions if it does not exist.
func initializeOutputDirectory(writeTo string) error {
	if _, err := os.Stat(writeTo); os.IsNotExist(err) {
		err := os.MkdirAll(writeTo, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

// main initializes the program by parsing flags and executes the run function. It logs a fatal error if the run fails.
func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run orchestrates the execution of initializing directories, authenticating, and exporting contacts and groups.
// It interacts with the Google People API and handles both contacts and contact groups data.
// Returns an error if any of the steps fail.
func run() error {
	ctx := context.Background()

	writeTo, err := obsidianutils.ApplyDirectoryPlaceHolder(outputDirectory)
	if err != nil {
		return err
	}
	err = initializeEnvironment(writeTo)
	if err != nil {
		return err
	}
	srv, err := initializeGoogleApiClient(err, ctx)
	if err != nil {
		return err
	}
	if printToConsole == "" || printToConsole == "contacts" {
		err = handleContacts(srv, writeTo)
		if err != nil {
			return err
		}
	}
	if printToConsole != "" && printToConsole != "groups" {
		return nil
	}
	return handleGroups(srv, writeTo)
}

// handleGroups retrieves Google Contact Groups, exporting data as JSON either by printing to the console or saving to a file.
func handleGroups(srv *people.Service, writeTo string) error {
	// List contact groups
	groups, err := srv.ContactGroups.List().PageSize(1000).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve contact groups: %w", err)
	}

	// Marshal groups data to JSON
	groupsJsonData, err := json.MarshalIndent(groups.ContactGroups, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal groups to JSON: %w", err)
	}

	if printToConsole == "groups" {
		fmt.Println(string(groupsJsonData))
	} else if printToConsole == "" {
		// Write groups data to file
		groupsOutputFile := path.Join(writeTo, "groups.json")
		err = ioutil.WriteFile(groupsOutputFile, groupsJsonData, 0644)
		if err != nil {
			return fmt.Errorf("unable to write groups to file: %w", err)
		}

		groupsAbsPath, _ := filepath.Abs(groupsOutputFile)
		if verbose {
			fmt.Printf("Successfully exported %d groups to %s\n", len(groups.ContactGroups), groupsAbsPath)
		}
	}
	return nil
}

// handleContacts retrieves and processes Google Contacts, exporting data to JSON either by printing or saving to a file.
func handleContacts(srv *people.Service, writeTo string) error {
	// List connections (contacts)
	r, err := srv.People.Connections.List("people/me").
		PersonFields("names,emailAddresses,phoneNumbers,addresses,organizations,memberships,birthdays").
		PageSize(1000).
		Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve contacts: %w", err)
	}

	// Marshal contacts data to JSON
	jsonData, err := json.MarshalIndent(r.Connections, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal contacts to JSON: %w", err)
	}

	if printToConsole == "contacts" {
		fmt.Println(string(jsonData))
	} else if printToConsole == "" {
		// Write contacts data to file
		outputFile := path.Join(writeTo, "contacts.json")
		err = ioutil.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("unable to write contacts to file: %w", err)
		}

		absPath, _ := filepath.Abs(outputFile)
		if verbose {
			fmt.Printf("Successfully exported %d contacts to %s\n", len(r.Connections), absPath)
		}
	}
	return err
}

// initializeGoogleApiClient initializes and returns a Google People Service client using OAuth2.
func initializeGoogleApiClient(err error, ctx context.Context) (*people.Service, error) {
	// Check if credentials.json exists
	credFile := path.Join(stateDirectory, "credentials.json")
	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("missing credentials file: %s\nPlease download it from Google Cloud Console", credFile)
	}

	b, err := os.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	// Configure the OAuth2 client
	config, err := google.ConfigFromJSON(b, contactsScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client, err := getClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token: %w", err)
	}

	// Create the People service
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create People service: %w", err)
	}
	return srv, nil
}

// initializeEnvironment sets up the necessary directories for application state and output, returning an error on failure.
func initializeEnvironment(writeTo string) error {
	err := initializeStateDirectory()
	if err != nil {
		log.Fatal(err)
	}

	err = initializeOutputDirectory(writeTo)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
