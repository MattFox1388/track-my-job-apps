package database

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"track-my-job-apps/internal/models"
)

var db *gorm.DB

// InitDatabase initializes the SQLite database connection and creates tables
func InitDatabase() error {
	var err error

	// Open SQLite database with pure Go driver
	sqlDB, err := sql.Open("sqlite", "job_apps.db?_pragma=foreign_keys(1)")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	db, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.JobApplication{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// Create FTS5 virtual table for search
	err = createFTSTable()
	if err != nil {
		return fmt.Errorf("failed to create FTS table: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// CreateApp creates a new job application in the database
func CreateApp(app *models.JobApplication) error {
	result := db.Create(app)
	if result.Error != nil {
		return fmt.Errorf("failed to create app: %v", result.Error)
	}
	return nil
}

// GetAllApps retrieves all job applications from the database
func GetAllApps() ([]models.JobApplication, error) {
	var apps []models.JobApplication
	result := db.Order("date_applied DESC").Find(&apps).Limit(20)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get apps: %v", result.Error)
	}
	return apps, nil
}

// GetAppByID retrieves a job application by ID
func GetAppByID(id uint) (*models.JobApplication, error) {
	var app models.JobApplication
	result := db.First(&app, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get app: %v", result.Error)
	}
	return &app, nil
}

// UpdateApp updates a job application
func UpdateApp(app *models.JobApplication) error {
	result := db.Save(app)
	if result.Error != nil {
		return fmt.Errorf("failed to update app: %v", result.Error)
	}
	return nil
}

// DeleteApp deletes a job application
func DeleteApp(id uint) error {
	result := db.Delete(&models.JobApplication{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete app: %v", result.Error)
	}
	return nil
}

// createFTSTable creates the FTS5 virtual table for full-text search
func createFTSTable() error {
	// Create FTS5 virtual table
	sql := `CREATE VIRTUAL TABLE IF NOT EXISTS apps_fts USING fts5(
		company, 
		position, 
		notes, 
		content='apps', 
		content_rowid='app_id'
	)`

	result := db.Exec(sql)
	return result.Error
}

// SearchApps performs full-text search on job applications
func SearchApps(query string) ([]models.JobApplication, error) {
	var apps []models.JobApplication

	// Search using FTS5 and join with main table
	sql := `SELECT a.* FROM apps a 
			JOIN apps_fts fts ON a.appId = fts.rowid 
			WHERE apps_fts MATCH ? 
			ORDER BY bm25(apps_fts)`

	result := db.Raw(sql, query).Scan(&apps)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to search apps: %v", result.Error)
	}

	return apps, nil
}

func SearchByCompany(companyName string) ([]models.JobApplication, error) {
	var apps []models.JobApplication
	result := db.Where("LOWER(company) LIKE LOWER(?)", "%"+companyName+"%").Find(&apps).Limit(10)
	if result.Error != nil {
		return nil, result.Error
	}
	return apps, nil
}
