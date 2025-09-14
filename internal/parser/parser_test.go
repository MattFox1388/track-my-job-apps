package parser

import (
	"os"
	"strings"
	"testing"
	"time"

	"track-my-job-apps/internal/models"
)

func TestParseGreenhouseJob(t *testing.T) {
	// Read the test HTML data from scratch.txt (in project root)
	testData, err := os.ReadFile("../../scratch.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Convert to lines as the function expects
	lines := strings.Split(string(testData), "\n")

	// Create a new JobApplication instance
	jobApp := &models.JobApplication{}

	// Test the parsing function
	result, err := ParseGreenhouseJob(lines, jobApp)
	if err != nil {
		t.Fatalf("ParseGreenhouseJob failed: %v", err)
	}

	// Test 1: Job title extraction
	expectedPosition := "Software Engineer, Enterprise Readiness"
	if result.Position != expectedPosition {
		t.Errorf("Expected position '%s', got '%s'", expectedPosition, result.Position)
	}

	// Test 2: Location extraction
	expectedLocation := "Hybrid"
	if result.Location != expectedLocation {
		t.Errorf("Expected location '%s', got '%s'", expectedLocation, result.Location)
	}

	// Test 3: Company extraction
	expectedCompany := "Cloudflare"
	if result.Company != expectedCompany {
		t.Errorf("Expected company '%s', got '%s'", expectedCompany, result.Company)
	}

	// Test 4: URL extraction in notes
	if !strings.Contains(result.Notes, "https://job-boards.greenhouse.io/cloudflare/jobs/6467350?gh_jid=6467350&gh_src=5ylsd31") {
		t.Errorf("Expected URL to be extracted in notes, got: %s", result.Notes)
	}

	// Test 5: Status should be set to SUBMITTED
	if result.Status != models.SUBMITTED {
		t.Errorf("Expected status to be SUBMITTED, got %s", result.Status)
	}

	// Test 6: DateApplied should be set to today
	today := time.Now().Format("2006-01-02")
	appliedDate := result.DateApplied.Time.Format("2006-01-02")
	if appliedDate != today {
		t.Errorf("Expected date applied to be today (%s), got %s", today, appliedDate)
	}

	// Print all extracted data for manual verification
	t.Logf("Extracted Job Application:")
	t.Logf("  Company: %s", result.Company)
	t.Logf("  Position: %s", result.Position)
	t.Logf("  Location: %s", result.Location)
	t.Logf("  Salary Range: %s", result.SalaryRange)
	t.Logf("  Status: %s", result.Status)
	t.Logf("  Date Applied: %s", result.DateApplied.Time.Format("2006-01-02"))
	t.Logf("  Notes: %s", result.Notes)
}

func TestParseGreenhouseJob2(t *testing.T) {
	// Read the test HTML data from scratch.txt (in project root)
	testData, err := os.ReadFile("../../scratch2.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Convert to lines as the function expects
	lines := strings.Split(string(testData), "\n")

	// Create a new JobApplication instance
	jobApp := &models.JobApplication{}

	// Test the parsing function
	result, err := ParseGreenhouseJob(lines, jobApp)
	if err != nil {
		t.Fatalf("ParseGreenhouseJob failed: %v", err)
	}

	// Test 1: Job title extraction
	expectedPosition := "Senior Software Engineer, Product Foundations"
	if result.Position != expectedPosition {
		t.Errorf("Expected position '%s', got '%s'", expectedPosition, result.Position)
	}

	// Test 2: Location extraction
	expectedLocation := "Denver, Colorado"
	if !strings.Contains(result.Location, expectedLocation) {
		t.Errorf("Expected location '%s', got '%s'", expectedLocation, result.Location)
	}

	// Test 3: Company extraction
	expectedCompany := "checkr"
	if strings.ToLower(result.Company) != expectedCompany {
		t.Errorf("Expected company '%s', got '%s'", expectedCompany, result.Company)
	}

	// Test 4: URL extraction in notes
	if !strings.Contains(result.Notes, "https://job-boards.greenhouse.io/checkr/jobs/7078460") {
		t.Errorf("Expected URL to be extracted in notes, got: %s", result.Notes)
	}

	// Test 5: Status should be set to SUBMITTED
	if result.Status != models.SUBMITTED {
		t.Errorf("Expected status to be SUBMITTED, got %s", result.Status)
	}

	// Test 6: DateApplied should be set to today
	today := time.Now().Format("2006-01-02")
	appliedDate := result.DateApplied.Time.Format("2006-01-02")
	if appliedDate != today {
		t.Errorf("Expected date applied to be today (%s), got %s", today, appliedDate)
	}

	// Print all extracted data for manual verification
	t.Logf("Extracted Job Application:")
	t.Logf("  Company: %s", result.Company)
	t.Logf("  Position: %s", result.Position)
	t.Logf("  Location: %s", result.Location)
	t.Logf("  Salary Range: %s", result.SalaryRange)
	t.Logf("  Status: %s", result.Status)
	t.Logf("  Date Applied: %s", result.DateApplied.Time.Format("2006-01-02"))
	t.Logf("  Notes: %s", result.Notes)
}

func TestParseGreenhouseJob3(t *testing.T) {
	// Read the test HTML data from scratch.txt (in project root)
	testData, err := os.ReadFile("../../scratch3.txt")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Convert to lines as the function expects
	lines := strings.Split(string(testData), "\n")

	// Create a new JobApplication instance
	jobApp := &models.JobApplication{}

	// Test the parsing function
	result, err := ParseGreenhouseJob(lines, jobApp)
	if err != nil {
		t.Fatalf("ParseGreenhouseJob failed: %v", err)
	}

	// Test 1: Job title extraction
	expectedPosition := "Software Engineer, Product Engineering"
	if result.Position != expectedPosition {
		t.Errorf("Expected position '%s', got '%s'", expectedPosition, result.Position)
	}

	// Test 2: Location extraction
	expectedLocation := "New York, NY"
	if !strings.Contains(result.Location, expectedLocation) {
		t.Errorf("Expected location '%s', got '%s'", expectedLocation, result.Location)
	}

	// Test 3: Company extraction
	expectedCompany := "figma"
	if strings.ToLower(result.Company) != expectedCompany {
		t.Errorf("Expected company '%s', got '%s'", expectedCompany, result.Company)
	}

	// Test 4: URL extraction in notes
	if !strings.Contains(result.Notes, "https://job-boards.greenhouse.io/figma/jobs/5552540004?gh_jid=5552540004&gh_src=28109e334us&source=LinkedIn") {
		t.Errorf("Expected URL to be extracted in notes, got: %s", result.Notes)
	}

	// Test 5: Status should be set to SUBMITTED
	if result.Status != models.SUBMITTED {
		t.Errorf("Expected status to be SUBMITTED, got %s", result.Status)
	}

	// Test 6: DateApplied should be set to today
	today := time.Now().Format("2006-01-02")
	appliedDate := result.DateApplied.Time.Format("2006-01-02")
	if appliedDate != today {
		t.Errorf("Expected date applied to be today (%s), got %s", today, appliedDate)
	}

	// Print all extracted data for manual verification
	t.Logf("Extracted Job Application:")
	t.Logf("  Company: %s", result.Company)
	t.Logf("  Position: %s", result.Position)
	t.Logf("  Location: %s", result.Location)
	t.Logf("  Salary Range: %s", result.SalaryRange)
	t.Logf("  Status: %s", result.Status)
	t.Logf("  Date Applied: %s", result.DateApplied.Time.Format("2006-01-02"))
	t.Logf("  Notes: %s", result.Notes)
}

func TestParseGreenhouseJobWithSalary(t *testing.T) {
	// Test with sample HTML that includes salary information using correct CSS classes
	testHTML := `
	<div class="job__header">
		<div class="job__title">
			<h1 class="section-header section-header--large font-primary">Senior Software Engineer</h1>
			<div class="job__location">
				<svg class="svg-icon"></svg>
				<div>Remote</div>
			</div>
		</div>
	</div>
	<div>
		<p>Salary: $120,000 - $180,000 per year</p>
		<p>Additional compensation: $150,000-$200,000 equity</p>
	</div>
	<img alt="TestCompany Logo" src="test.png">
	mf-URL: https://example.com/job/123
	`

	lines := strings.Split(testHTML, "\n")
	jobApp := &models.JobApplication{}

	result, err := ParseGreenhouseJob(lines, jobApp)
	if err != nil {
		t.Fatalf("ParseGreenhouseJob failed: %v", err)
	}

	// Test salary extraction
	expectedSalaryCount := 2 // Should find two salary ranges
	salaryRanges := strings.Split(result.SalaryRange, ", ")
	if len(salaryRanges) != expectedSalaryCount {
		t.Errorf("Expected %d salary ranges, got %d: %s", expectedSalaryCount, len(salaryRanges), result.SalaryRange)
	}

	// Check if both salary ranges are captured
	if !strings.Contains(result.SalaryRange, "$120,000 - $180,000") {
		t.Errorf("Expected to find '$120,000 - $180,000' in salary range, got: %s", result.SalaryRange)
	}
	if !strings.Contains(result.SalaryRange, "$150,000 - $200,000") {
		t.Errorf("Expected to find '$150,000 - $200,000' in salary range, got: %s", result.SalaryRange)
	}

	// Test other fields
	if result.Position != "Senior Software Engineer" {
		t.Errorf("Expected position 'Senior Software Engineer', got '%s'", result.Position)
	}
	if result.Location != "Remote" {
		t.Errorf("Expected location 'Remote', got '%s'", result.Location)
	}
	if result.Company != "TestCompany" {
		t.Errorf("Expected company 'TestCompany', got '%s'", result.Company)
	}
	if !strings.Contains(result.Notes, "https://example.com/job/123") {
		t.Errorf("Expected URL in notes, got: %s", result.Notes)
	}

	// Log the results for verification
	t.Logf("Salary Test Results:")
	t.Logf("  Position: %s", result.Position)
	t.Logf("  Location: %s", result.Location)
	t.Logf("  Company: %s", result.Company)
	t.Logf("  Salary Range: %s", result.SalaryRange)
	t.Logf("  Notes: %s", result.Notes)
}

func TestCleanHTMLTags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<h1>Software Engineer</h1>",
			expected: "Software Engineer",
		},
		{
			input:    "<div class=\"test\">Hello &amp; Welcome</div>",
			expected: "Hello & Welcome",
		},
		{
			input:    "Plain text with no tags",
			expected: "Plain text with no tags",
		},
		{
			input:    "<p>Multiple&nbsp;<strong>tags</strong>&nbsp;here</p>",
			expected: "Multiple tags here",
		},
		{
			input:    "  <span>  Whitespace  </span>  ",
			expected: "Whitespace",
		},
	}

	for _, test := range tests {
		result := CleanHTMLTags(test.input)
		if result != test.expected {
			t.Errorf("CleanHTMLTags('%s') = '%s', expected '%s'", test.input, result, test.expected)
		}
	}
}

func TestParseLinkedInJob(t *testing.T) {
	// Test LinkedIn parsing with sample data
	testLines := []string{
		"TestCompany Inc",
		"Senior Developer",
		"San Francisco, CA · Remote",
		"$100,000 - $150,000",
		"Full-time",
		"Share",
		"Show more options",
		"Matches your job preferences, workplace type is Remote",
	}

	jobApp := &models.JobApplication{}
	result, err := ParseLinkedInJob(testLines, jobApp)
	if err != nil {
		t.Fatalf("ParseLinkedInJob failed: %v", err)
	}

	// Test extraction
	if result.Company != "TestCompany Inc" {
		t.Errorf("Expected company 'TestCompany Inc', got '%s'", result.Company)
	}
	if result.Position != "Senior Developer" {
		t.Errorf("Expected position 'Senior Developer', got '%s'", result.Position)
	}
	if result.Location != "San Francisco, CA" {
		t.Errorf("Expected location 'San Francisco, CA', got '%s'", result.Location)
	}
	if result.WorkplaceType != "Remote" {
		t.Errorf("Expected workplace type 'Remote', got '%s'", result.WorkplaceType)
	}

	t.Logf("LinkedIn Test Results:")
	t.Logf("  Company: %s", result.Company)
	t.Logf("  Position: %s", result.Position)
	t.Logf("  Location: %s", result.Location)
	t.Logf("  Workplace Type: %s", result.WorkplaceType)
}
