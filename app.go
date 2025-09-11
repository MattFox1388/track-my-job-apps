package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx context.Context
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
	if err := InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

func (a *App) TrackJobApp(jobAppData string, platform string) (*JobApplication, error) {
	fmt.Printf("Received job app data from %s: %s\n", platform, jobAppData)

	lines := strings.Split(jobAppData, "\n")	
	jobApp := &JobApplication{}
	
	// Platform-specific parsing logic
	switch platform {
	case "linkedin":
		return parseLinkedInJob(lines, jobApp)
	case "greenhouse":
		return parseGreenhouseJob(lines, jobApp)
	default:
		return parseLinkedInJob(lines, jobApp) // Default to LinkedIn parsing
	}
}

func parseLinkedInJob(lines []string, jobApp *JobApplication) (*JobApplication, error) {

	for i, line := range lines {
		if (strings.TrimSpace(line) == "") {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	for i, line := range lines {
		if strings.Contains(line, "$") {
			jobApp.SalaryRange = strings.TrimSpace(line)
		}

		workspaceTypeStr := "Matches your job preferences, workplace type is"
		foundWorkspaceType := strings.Index(line, "Matches your job preferences, workplace type is")
		if foundWorkspaceType != -1 {
			jobApp.WorkplaceType = strings.TrimSpace(line[foundWorkspaceType+len(workspaceTypeStr):])
		}
		if i == 1 {
			jobApp.Company = strings.TrimSpace(line)
		} else if i == 4 {
			jobApp.Position = strings.TrimSpace(line)
		} else if i== 5 {
			parts := strings.Split(line, "·")
			if len(parts) > 0 {
				jobApp.Location = strings.TrimSpace(parts[0])
			}
		}
	}
	jobApp.Status = SUBMITTED
	jobApp.DateApplied = DateOnly{time.Now()}
	
	return jobApp, nil
}

func parseGreenhouseJob(lines []string, jobApp *JobApplication) (*JobApplication, error) {
	
}

func (a *App) SaveJobApp(jobApp *JobApplication) error {
	if err := CreateApp(jobApp); err != nil {
        fmt.Printf("Error saving job app: %v\n", err)
        return err
    }
    
    fmt.Printf("Saved job app: %s at %s (ID: %d)\n", jobApp.Position, jobApp.Company, jobApp.AppId)
    return nil
}

// GetAllJobApps returns all job applications from the database
func (a *App) GetAllJobApps() ([]JobApplication, error) {
	apps, err := GetAllApps()
	if err != nil {
		fmt.Printf("Error getting job apps: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("Retrieved %d job applications\n", len(apps))
	return apps, nil
}

func (a *App) SearchByCompany(companyName string) ([]JobApplication, error) {
	apps, err := SearchByCompany(companyName)
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
