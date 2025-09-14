package tests

import (
	"strings"
	"testing"

	"track-my-job-apps/internal/database"
	"track-my-job-apps/internal/models"
	"track-my-job-apps/internal/parser"
)

func TestIntegrationParseAndSave(t *testing.T) {
	// Initialize test database
	err := database.InitDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Test data
	testHTML := `
	<div class="job__header">
		<div class="job__title">
			<h1>Integration Test Engineer</h1>
			<div class="job__location">
				<div>Seattle, WA</div>
			</div>
		</div>
	</div>
	<img alt="TestCorp Logo" src="test.png">
	mf-URL: https://example.com/job/integration-test
	`

	// Parse the job
	lines := strings.Split(testHTML, "\n")
	jobApp := &models.JobApplication{}

	result, err := parser.ParseGreenhouseJob(lines, jobApp)
	if err != nil {
		t.Fatalf("Failed to parse job: %v", err)
	}

	// Save to database
	err = database.CreateApp(result)
	if err != nil {
		t.Fatalf("Failed to save job: %v", err)
	}

	// Verify it was saved
	if result.AppId == 0 {
		t.Errorf("Expected AppId to be set after saving, got 0")
	}

	// Retrieve it back
	retrieved, err := database.GetAppByID(result.AppId)
	if err != nil {
		t.Fatalf("Failed to retrieve job: %v", err)
	}

	// Verify data integrity
	if retrieved.Position != "Integration Test Engineer" {
		t.Errorf("Expected position 'Integration Test Engineer', got '%s'", retrieved.Position)
	}
	if retrieved.Company != "TestCorp" {
		t.Errorf("Expected company 'TestCorp', got '%s'", retrieved.Company)
	}
	if retrieved.Location != "Seattle, WA" {
		t.Errorf("Expected location 'Seattle, WA', got '%s'", retrieved.Location)
	}

	t.Logf("Integration test successful:")
	t.Logf("  Saved and retrieved job ID: %d", retrieved.AppId)
	t.Logf("  Position: %s", retrieved.Position)
	t.Logf("  Company: %s", retrieved.Company)
	t.Logf("  Location: %s", retrieved.Location)
}
