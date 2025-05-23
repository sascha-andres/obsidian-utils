package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"

	"github.com/sascha-andres/reuse/flag"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"

	obsidianutils "github.com/sascha-andres/obsidian-utils"
	"github.com/sascha-andres/obsidian-utils/internal"
)

// Scopes required for reading contacts and contact groups
const contactsScope = "https://www.googleapis.com/auth/contacts.readonly"

var (
	stateDirectory, outputDirectory, printToConsole, logLevel string
	verbose                                                   bool
)

// init initializes the program's environment settings and configuration for Google-related utilities.
// It sets an environment prefix, retrieves the current user, and defines the state directory flag for OAuth2 storage.
func init() {
	obsidianutils.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("GGL")

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&stateDirectory, "state-directory", path.Join(currentUser.HomeDir, ".local/state/ggl"), "Directory to store OAuth2 state")
	flag.StringVar(&outputDirectory, "output-directory", ".", "Directory to store output files")
	flag.StringVar(&printToConsole, "print-to-console", "", "Print data to console instead of writing to files, may be contacts or groups")
	flag.StringVar(&logLevel, "log-level", "info", "Log level, one of: debug, info, warn, error, fatal")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
}

// getClient retrieves a token, saves it, then returns the OAuth2 client
func getClient(logger *slog.Logger, config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens
	tokenFile := path.Join(stateDirectory, "token.json")
	tok, err := tokenFromFile(logger, tokenFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		err = saveToken(logger, tokenFile, tok)
		if err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), tok), nil
}

// openBrowser opens a browser window with the specified URL
func openBrowser(url string) {
	var err error

	switch {
	case len(os.Getenv("BROWSER")) > 0:
		err = exec.Command(os.Getenv("BROWSER"), url).Start()
	case os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "":
		err = exec.Command("xdg-open", url).Start()
	case os.Getenv("XDG_SESSION_TYPE") == "wayland":
		err = exec.Command("xdg-open", url).Start()
	case runtime.GOOS == "darwin":
		err = exec.Command("open", url).Start()
	case runtime.GOOS == "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

// getTokenFromWeb requests a token from the web by starting a local web server to handle the OAuth callback
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// Set up a local web server to receive the callback
	codeChan := make(chan string)
	errChan := make(chan error)

	// Create a random port between 8000 and 9000
	const callbackPath = "/oauth2callback"
	redirectURL := "http://localhost:8080" + callbackPath
	config.RedirectURL = redirectURL

	// Set up the server
	server := &http.Server{Addr: ":8080"}

	// Set up the handler
	http.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no code in callback")
			http.Error(w, "No code provided", http.StatusBadRequest)
			return
		}

		// Display a success message to the user
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, "<html><body><h1>Authentication Successful</h1><p>You can close this window now.</p></body></html>")

		// Send the code to the channel
		codeChan <- code
	})

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// Generate the auth URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// Prompt the user to open the URL
	fmt.Println("==========================================================")
	fmt.Println("To authorize this application:")
	fmt.Println("1. A browser window should open automatically.")
	fmt.Println("   If it doesn't, please open the following URL:")
	fmt.Printf("   %v\n", authURL)
	fmt.Println("2. Sign in and grant access to your Google account")
	fmt.Println("==========================================================")

	// Try to open the browser automatically
	openBrowser(authURL)

	// Wait for the code or an error
	var code string
	select {
	case code = <-codeChan:
		// Got the code, continue
	case err := <-errChan:
		// Shutdown the server
		_ = server.Shutdown(context.Background())
		return nil, err
	}

	// Shutdown the server
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	// Exchange the code for a token
	return config.Exchange(context.TODO(), code)
}

// tokenFromFile retrieves a token from a local file
func tokenFromFile(logger *slog.Logger, file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			logger.Error("error closing token file", "err", err)
		}
	}()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file
func saveToken(logger *slog.Logger, path string, token *oauth2.Token) error {
	if verbose {
		fmt.Printf("Saving credential file to: %s\n", path)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			logger.Error("error closing token file", "err", err)
		}
	}()
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
	logger := internal.CreateLogger("OBS_UTIL_DAILY", logLevel)
	if err := run(logger); err != nil {
		logger.Error("error running", "err", err)
	}
}

// run orchestrates the execution of initializing directories, authenticating, and exporting contacts and groups.
// It interacts with the Google People API and handles both contacts and contact groups data.
// Returns an error if any of the steps fail.
func run(logger *slog.Logger) error {
	ctx := context.Background()

	writeTo, err := obsidianutils.ApplyDirectoryPlaceHolder(outputDirectory)
	if err != nil {
		return err
	}
	err = initializeEnvironment(logger, writeTo)
	if err != nil {
		return err
	}
	srv, err := initializeGoogleApiClient(logger, ctx)
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
		err = os.WriteFile(groupsOutputFile, groupsJsonData, 0644)
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
		// Write contacts data to a file
		outputFile := path.Join(writeTo, "contacts.json")
		err = os.WriteFile(outputFile, jsonData, 0644)
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
func initializeGoogleApiClient(logger *slog.Logger, ctx context.Context) (*people.Service, error) {
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
	client, err := getClient(logger, config)
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
func initializeEnvironment(logger *slog.Logger, writeTo string) error {
	err := initializeStateDirectory()
	if err != nil {
		logger.Error("error state directory", "err", err)
		return err
	}

	err = initializeOutputDirectory(writeTo)
	if err != nil {
		logger.Error("error initializing output directory", "err", err)
	}
	return err
}
