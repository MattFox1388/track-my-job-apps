package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"track-my-job-apps/internal/backup"
	"track-my-job-apps/internal/database"
	"track-my-job-apps/internal/models"
	"track-my-job-apps/internal/parser"
)

// App struct
type App struct {
	ctx    context.Context
	backup *backup.BackupService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts up and can be used to
// initialize the application.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize backup service (don't fail if backup setup is incomplete)
	backupService, err := backup.NewBackupService()
	if err != nil {
		log.Printf("Warning: Failed to initialize backup service: %v", err)
		log.Printf("Backup will be skipped. Make sure you have credentials.json and completed OAuth setup.")
	} else {
		log.Println("Backup service initialized successfully")
		a.backup = backupService
	}
}

func (a *App) TrackJobApp(jobAppData string, platform string) (*models.JobApplication, error) {
	fmt.Printf("Received job app data from %s: %s\n", platform, jobAppData)

	lines := strings.Split(jobAppData, "\n")
	jobApp := &models.JobApplication{}

	// Platform-specific parsing logic
	switch platform {
	case "linkedin":
		return parser.ParseLinkedInJob(lines, jobApp)
	case "greenhouse":
		return parser.ParseGreenhouseJob(lines, jobApp)
	default:
		return parser.ParseLinkedInJob(lines, jobApp) // Default to LinkedIn parsing
	}
}

func (a *App) SaveJobApp(jobApp *models.JobApplication) error {
	if err := database.CreateApp(jobApp); err != nil {
		fmt.Printf("Error saving job app: %v\n", err)
		return err
	}

	fmt.Printf("Saved job app: %s at %s (ID: %d)\n", jobApp.Position, jobApp.Company, jobApp.AppId)
	return nil
}

// GetAllJobApps returns all job applications from the database
func (a *App) GetAllJobApps() ([]models.JobApplication, error) {
	apps, err := database.GetAllApps()
	if err != nil {
		fmt.Printf("Error getting job apps: %v\n", err)
		return nil, err
	}

	fmt.Printf("Retrieved %d job applications\n", len(apps))
	return apps, nil
}

func (a *App) SearchByCompany(companyName string) ([]models.JobApplication, error) {
	apps, err := database.SearchByCompany(companyName)
	if err != nil {
		fmt.Printf("Error searching by company: %v\n", err)
		return nil, err
	}
	return apps, nil
}

// Helper function to safely get string from map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (a *App) BeforeClose(ctx context.Context) bool {
	if a.backup != nil {
		log.Println("Backing up database before closing...")
		if err := a.backup.BackupDatabase("job_apps.db"); err != nil {
			log.Printf("Error backing up database: %v", err)
		}
	}
	return false
}

// TestBackup manually triggers a backup (for testing)
func (a *App) TestBackup() error {
	if a.backup == nil {
		return fmt.Errorf("backup service not initialized")
	}

	log.Println("Testing backup...")
	err := a.backup.BackupDatabase("job_apps.db")
	if err != nil {
		log.Printf("Backup test failed: %v", err)
		return err
	}
	log.Println("Backup test successful!")
	return nil
}
