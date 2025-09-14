package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type BackupService struct {
	service *drive.Service
	ctx     context.Context
}

// NewBackupService creates a new backup service
func NewBackupService() (*BackupService, error) {
	ctx := context.Background()

	// Load credentials from file (you'll need to set this up)
	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	return &BackupService{service: srv, ctx: ctx}, nil
}

// BackupDatabase uploads the SQLite database to Google Drive
func (bs *BackupService) BackupDatabase(dbPath string) error {
	file, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("unable to open database file: %v", err)
	}
	defer file.Close()

	// Check if backup already exists
	existingFileID, err := bs.findExistingBackup()
	if err != nil {
		log.Printf("Warning: Could not check for existing backup: %v", err)
	}

	// Create file metadata - remove appDataFolder
	driveFile := &drive.File{
		Name: "job_apps_backup.db",
		// Remove Parents line - will go to root Drive folder
	}

	// Upload or update the file
	if existingFileID != "" {
		// Update existing file
		_, err = bs.service.Files.Update(existingFileID, driveFile).Media(file).Do()
		if err != nil {
			return fmt.Errorf("unable to update backup: %v", err)
		}
		log.Println("Database backup updated successfully")
	} else {
		// Create new file
		_, err = bs.service.Files.Create(driveFile).Media(file).Do()
		if err != nil {
			return fmt.Errorf("unable to create backup: %v", err)
		}
		log.Println("Database backup created successfully")
	}

	return nil
}

// findExistingBackup finds existing backup file ID
func (bs *BackupService) findExistingBackup() (string, error) {
	files, err := bs.service.Files.List().
		Q("name='job_apps_backup.db'").
		Do()
	if err != nil {
		return "", err
	}

	if len(files.Files) > 0 {
		return files.Files[0].Id, nil
	}
	return "", nil
}

// getClient retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		log.Println("Token retrieved from web")
		log.Println(tok)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Create a random state
	state := "random-state-string"
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Start local server to catch callback
	server := &http.Server{Addr: ":80"}
	authCode := make(chan string, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received callback: %s", r.URL.String())

		code := r.URL.Query().Get("code")
		receivedState := r.URL.Query().Get("state")

		log.Printf("Code: %s, State: %s", code, receivedState)

		if receivedState == state && code != "" {
			authCode <- code
			w.Write([]byte("Authorization successful! You can close this window."))
		} else {
			w.Write([]byte("Authorization failed!"))
		}
	})

	go func() {
		log.Printf("Starting server on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Open browser
	err := openBrowser(authURL)
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}

	// Wait for auth code
	log.Printf("Waiting for auth code...")
	code := <-authCode
	log.Printf("Received auth code: %s", code)
	server.Close()

	log.Printf("Exchanging code for token...")
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	log.Printf("Token exchange successful!")
	return tok
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

// tokenFromFile retrieves a token from a local file.
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

// saveToken saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
