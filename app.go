package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"

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

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	// Check if database exists
	// dbExists := database.DatabaseExists()

	// Initialize backup service first (needed for restore)
	backupService, err := backup.NewBackupService()
	if err != nil {
		log.Printf("Warning: Failed to initialize backup service: %v", err)
		log.Printf("Backup will be skipped. Make sure you have credentials.json and completed OAuth setup.")
	} else {
		log.Println("Backup service initialized successfully")
		a.backup = backupService

		// Try to restore from backup
		log.Println("Database not found, attempting to restore from backup...")
		if err := a.backup.RestoreDatabase("job_apps.db"); err != nil {
			log.Printf("Warning: Failed to restore database from backup: %v", err)
			log.Println("Continuing with fresh database...")
		} else {
			log.Println("Database restored successfully from backup")
		}
	}

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Turso database connection
	if err := database.InitTurso(); err != nil {
		log.Printf("Warning: Failed to initialize Turso database: %v", err)
		log.Println("Turso features will be unavailable. Set TURSO_DATABASE_URL and TURSO_AUTH_TOKEN env vars.")
	} else {
		log.Println("Turso database initialized successfully")
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
	case "workday":
		return parser.ParseWorkdayJob(lines, jobApp)
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

func (a *App) UpdateJobApp(jobApp *models.JobApplication) error {
	if err := database.UpdateApp(jobApp); err != nil {
		fmt.Printf("Error updating job app: %v\n", err)
		return err
	}
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
	fmt.Printf("Searching by company: %s\n", companyName)
	apps, err := database.SearchByCompany(strings.ToLower(companyName))
	if err != nil {
		fmt.Printf("Error searching by company: %v\n", err)
		return nil, err
	}
	fmt.Printf("Found %d job applications for company: %s\n", len(apps), companyName)
	return apps, nil
}

func (a *App) SearchByFullText(query string) ([]models.JobApplication, error) {
	fmt.Printf("Searching by full text: %s\n", query)
	apps, err := database.SearchByFullText(strings.ToLower(query))
	if err != nil {
		fmt.Printf("Error searching by full text: %v\n", err)
		return nil, err
	}
	return apps, nil
}

// GetJobAppCount returns the total number of job applications in the database
func (a *App) GetJobAppCount() (int64, error) {
	count, err := database.GetAppCount()
	if err != nil {
		fmt.Printf("Error getting job app count: %v\n", err)
		return 0, err
	}
	fmt.Printf("Total job applications: %d\n", count)
	return count, nil
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

// GetRecentPostingsWithQualifications fetches recent job postings with qualifications from Turso
func (a *App) GetRecentPostingsWithQualifications() ([]map[string]interface{}, error) {
	query := `
		SELECT
		p.id,
		p.link,
		p.descrip,
		p.postedDate,
		p.companyName,
		p.reviewOutput,
		json_group_array(q.qualification_text) AS qualifications
		FROM postings p
		LEFT JOIN qualifications q
		ON p.id = q.postings_id
		GROUP BY p.id, p.descrip
		ORDER BY p.ingestDate DESC
		LIMIT 100;
	`

	log.Println("Fetching recent postings with qualifications from Turso...")
	results, err := database.ExecuteTursoSQL(query)
	if err != nil {
		log.Printf("Failed to fetch postings: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d postings", len(results))

	// Print first few results for debugging
	for i, row := range results {
		if i >= 10 { // Only print first 3 rows to avoid spam
			break
		}
		log.Printf("Row %d:", i)
		for key, value := range row {
			log.Printf("  %s: %v", key, value)
		}
	}

	return results, nil
}

// DeletePosting deletes a job posting and its qualifications from Turso by ID
func (a *App) DeletePosting(postingID int64) error {
	log.Printf("Deleting posting with ID: %d", postingID)

	// First delete qualifications
	deleteQualQuery := "DELETE FROM qualifications WHERE postings_id = ?"
	_, err := database.ExecuteTursoExec(deleteQualQuery, postingID)
	if err != nil {
		log.Printf("Error deleting qualifications: %v", err)
		return fmt.Errorf("failed to delete qualifications: %v", err)
	}

	// Then delete the posting
	deletePostingQuery := "DELETE FROM postings WHERE id = ?"
	rowsAffected, err := database.ExecuteTursoExec(deletePostingQuery, postingID)
	if err != nil {
		log.Printf("Error deleting posting: %v", err)
		return fmt.Errorf("failed to delete posting: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no posting found with ID: %d", postingID)
	}

	log.Printf("Successfully deleted posting ID %d", postingID)
	return nil
}
