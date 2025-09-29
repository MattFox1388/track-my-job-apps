package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"track-my-job-apps/internal/models"
)

// ParseGreenhouseJob extracts job data from Greenhouse HTML content
func ParseGreenhouseJob(lines []string, jobApp *models.JobApplication) (*models.JobApplication, error) {
	// Join all lines into a single HTML string
	htmlContent := strings.Join(lines, "\n")

	// 1. Extract job title
	extractJobTitle(htmlContent, jobApp)

	// 2. Extract location from <div class="job__location"><div>
	locationPattern := regexp.MustCompile(`(?s)<div[^>]*class="[^"]*job__location[^"]*"[^>]*>.*?<div[^>]*>(.*?)</div>`)
	if locationMatch := locationPattern.FindStringSubmatch(htmlContent); len(locationMatch) > 1 {
		jobApp.Location = strings.TrimSpace(CleanHTMLTags(locationMatch[1]))
	}

	// 3. Find salary patterns like $100,000 - $150,000
	salaryPattern := regexp.MustCompile(`\$(\d+(?:,\d{3})*)\s*-\s*\$(\d+(?:,\d{3})*)`)
	salaryMatches := salaryPattern.FindAllStringSubmatch(htmlContent, -1)
	var salaryRanges []string
	for _, match := range salaryMatches {
		if len(match) >= 3 {
			salaryRanges = append(salaryRanges, fmt.Sprintf("$%s - $%s", match[1], match[2]))
		}
	}

	salaryPattern2 := regexp.MustCompile(`\$(\d+(?:,\d{3})*)\s+to\s+\$(\d+(?:,\d{3})*)`)
	salaryMatches2 := salaryPattern2.FindAllStringSubmatch(htmlContent, -1)
	for _, match := range salaryMatches2 {
		if len(match) >= 3 {
			salaryRanges = append(salaryRanges, fmt.Sprintf("$%s - $%s", match[1], match[2]))
		}
	}

	if len(salaryRanges) > 0 {
		jobApp.SalaryRange = strings.Join(salaryRanges, ", ")
	}

	// 4. Extract URL from "mf-URL:" pattern and extract company name
	urlPattern := regexp.MustCompile(`mf-URL:\s*(.+)`)
	if urlMatch := urlPattern.FindStringSubmatch(htmlContent); len(urlMatch) > 1 {
		url := strings.TrimSpace(urlMatch[1])

		// Extract company name from Greenhouse URL pattern: https://job-boards.greenhouse.io/company/jobs/
		companyPattern := regexp.MustCompile(`https?://job-boards\.greenhouse\.io/([^/]+)/`)
		if companyMatch := companyPattern.FindStringSubmatch(url); len(companyMatch) > 1 {
			jobApp.Company = strings.TrimSpace(companyMatch[1])
		}

		// Store URL in notes
		if jobApp.Notes == "" {
			jobApp.Notes = fmt.Sprintf("Source URL: %s", url)
		} else {
			jobApp.Notes += fmt.Sprintf("\nSource URL: %s", url)
		}
	}

	// Extract company name from the page (look for Cloudflare in this case)
	companyPattern := regexp.MustCompile(`<img[^>]*alt="([^"]*Logo)"`)
	if companyMatch := companyPattern.FindStringSubmatch(htmlContent); len(companyMatch) > 1 {
		companyName := strings.Replace(companyMatch[1], " Logo", "", 1)
		jobApp.Company = strings.TrimSpace(companyName)
	}

	jobApp.Status = models.SUBMITTED
	jobApp.DateApplied = models.DateOnly{Time: time.Now()}

	return jobApp, nil
}

// extractJobTitle extracts job title from Greenhouse HTML content
func extractJobTitle(htmlContent string, jobApp *models.JobApplication) {
	// Try multiple possible CSS classes: job__title, job_description_title, job-description__title
	titlePattern := regexp.MustCompile(`(?s)<div[^>]*class="[^"]*(?:job__title|job_description_title|job-description__title)[^"]*"[^>]*>.*?<h1[^>]*>(.*?)</h1>`)
	if titleMatch := titlePattern.FindStringSubmatch(htmlContent); len(titleMatch) > 1 {
		jobApp.Position = strings.TrimSpace(CleanHTMLTags(titleMatch[1]))
	}
}

// CleanHTMLTags removes HTML tags and decodes entities
func CleanHTMLTags(html string) string {
	// Remove HTML tags
	tagPattern := regexp.MustCompile(`<[^>]*>`)
	cleaned := tagPattern.ReplaceAllString(html, "")

	// Decode common HTML entities
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&lt;", "<")
	cleaned = strings.ReplaceAll(cleaned, "&gt;", ">")
	cleaned = strings.ReplaceAll(cleaned, "&quot;", "\"")
	cleaned = strings.ReplaceAll(cleaned, "&#39;", "'")
	cleaned = strings.ReplaceAll(cleaned, "&nbsp;", " ")

	return strings.TrimSpace(cleaned)
}

// ParseLinkedInJob extracts job data from LinkedIn HTML content
func ParseLinkedInJob(lines []string, jobApp *models.JobApplication) (*models.JobApplication, error) {
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "Share" {
			continue
		}
		if strings.TrimSpace(line) == "Show more options" {
			continue
		}
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	lines = filteredLines

	for i, line := range lines {
		if strings.Contains(line, "$") {
			jobApp.SalaryRange = strings.TrimSpace(line)
		}

		workspaceTypeStr := "Matches your job preferences, workplace type is"
		foundWorkspaceType := strings.Index(line, "Matches your job preferences, workplace type is")
		if foundWorkspaceType != -1 {
			jobApp.WorkplaceType = strings.TrimSpace(line[foundWorkspaceType+len(workspaceTypeStr):])
		}
		switch i {
		case 0:
			jobApp.Company = strings.TrimSpace(line)
		case 1:
			jobApp.Position = strings.TrimSpace(line)
		case 2:
			parts := strings.Split(line, "·")
			if len(parts) > 0 {
				jobApp.Location = strings.TrimSpace(parts[0])
			}
		}
	}
	jobApp.Status = models.SUBMITTED
	jobApp.DateApplied = models.DateOnly{Time: time.Now()}

	return jobApp, nil
}

// ParseWorkdayJob extracts job data from Workday HTML content
func ParseWorkdayJob(lines []string, jobApp *models.JobApplication) (*models.JobApplication, error) {
	// Join all lines into a single HTML string
	htmlContent := strings.Join(lines, "\n")

	// 1. Extract job title from <h2 data-automation-id="jobPostingHeader">
	titlePattern := regexp.MustCompile(`(?s)<h2[^>]*data-automation-id="jobPostingHeader"[^>]*>(.*?)</h2>`)
	if titleMatch := titlePattern.FindStringSubmatch(htmlContent); len(titleMatch) > 1 {
		jobApp.Position = strings.TrimSpace(CleanHTMLTags(titleMatch[1]))
	}

	// 2. Extract locations from <div data-automation-id="locations"> and get all <dd> elements
	locationSectionPattern := regexp.MustCompile(`(?s)<div[^>]*data-automation-id="locations"[^>]*>.*?<dl[^>]*>(.*?)</dl>`)
	if locationMatch := locationSectionPattern.FindStringSubmatch(htmlContent); len(locationMatch) > 1 {
		ddPattern := regexp.MustCompile(`(?s)<dd[^>]*>(.*?)</dd>`)
		ddMatches := ddPattern.FindAllStringSubmatch(locationMatch[1], -1)
		var locations []string
		for _, match := range ddMatches {
			if len(match) > 1 {
				location := strings.TrimSpace(CleanHTMLTags(match[1]))
				if location != "" && location != "Remote" && len(locations) < 10 { // Take first 10 locations
					locations = append(locations, location)
				}
			}
		}
		if len(locations) > 0 {
			jobApp.Location = strings.Join(locations, "\n")
		}
	}

	// 3. Extract job type from <div data-automation-id="time"> and get all <dd> elements
	timePattern := regexp.MustCompile(`(?s)<div[^>]*data-automation-id="time"[^>]*>.*?<dl[^>]*>.*?<dd[^>]*>(.*?)</dd>`)
	if timeMatch := timePattern.FindStringSubmatch(htmlContent); len(timeMatch) > 1 {
		jobApp.WorkplaceType = strings.TrimSpace(CleanHTMLTags(timeMatch[1]))
	}

	// 4. Extract salary patterns like $100,000 - $150,000 or $100,000.00 - $150,000.00
	salaryPattern := regexp.MustCompile(`\$(\d+(?:,\d{3})*)(?:\.\d{2})?\s*-\s*\$(\d+(?:,\d{3})*)(?:\.\d{2})?`)
	salaryMatches := salaryPattern.FindAllStringSubmatch(htmlContent, -1)
	var salaryRanges []string
	for _, match := range salaryMatches {
		if len(match) >= 3 {
			salaryRanges = append(salaryRanges, fmt.Sprintf("$%s - $%s", match[1], match[2]))
		}
	}

	salaryPattern2 := regexp.MustCompile(`\$(\d+(?:,\d{3})*)(?:\.\d{2})?\s+to\s+\$(\d+(?:,\d{3})*)(?:\.\d{2})?`)
	salaryMatches2 := salaryPattern2.FindAllStringSubmatch(htmlContent, -1)
	for _, match := range salaryMatches2 {
		if len(match) >= 3 {
			salaryRanges = append(salaryRanges, fmt.Sprintf("$%s - $%s", match[1], match[2]))
		}
	}

	if len(salaryRanges) > 0 {
		jobApp.SalaryRange = strings.Join(salaryRanges, ", ")
	}

	// 5. Extract URL from "mf-URL:" pattern and extract company name
	urlPattern := regexp.MustCompile(`mf-URL:\s*(.+)`)
	if urlMatch := urlPattern.FindStringSubmatch(htmlContent); len(urlMatch) > 1 {
		url := strings.TrimSpace(urlMatch[1])

		// Extract company name from Workday URL pattern: https://company.wd#.myworkdayjobs.com/
		companyPattern := regexp.MustCompile(`https?://([^.]+)\.wd\d+\.myworkdayjobs\.com`)
		if companyMatch := companyPattern.FindStringSubmatch(url); len(companyMatch) > 1 {
			jobApp.Company = strings.TrimSpace(companyMatch[1])
		}

		// Store URL in notes
		if jobApp.Notes == "" {
			jobApp.Notes = fmt.Sprintf("Source URL: %s", url)
		} else {
			jobApp.Notes += fmt.Sprintf("\nSource URL: %s", url)
		}
	}

	jobApp.Status = models.SUBMITTED
	jobApp.DateApplied = models.DateOnly{Time: time.Now()}

	return jobApp, nil
}
