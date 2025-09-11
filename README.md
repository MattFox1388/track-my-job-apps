# Track My Job Apps

A desktop application for tracking job applications built with Go, Wails, and React.

## Features

- Track job applications with company, position, location, salary range, and notes
- Full-text search across company names, positions, and notes using SQLite FTS5
- Modern React frontend with Wails desktop integration
- SQLite database with GORM ORM

## Database Schema

### JobApplication Model
```go
type JobApplication struct {
    AppId         uint      `gorm:"primaryKey;autoIncrement" json:"appId"`
    Company       string    `gorm:"type:varchar(255);not null" json:"company"`
    Position      string    `gorm:"type:varchar(255);not null" json:"position"`
    Location      string    `gorm:"type:varchar(255)" json:"location"`
    SalaryRange   string    `gorm:"type:varchar(100)" json:"salaryRange"`
    WorkplaceType string    `gorm:"type:varchar(50)" json:"workplaceType"`
    Status        Status    `gorm:"type:varchar(50);default:SUBMITTED" json:"status"`
    Notes         string    `gorm:"type:text" json:"notes"`
    Website       string    `gorm:"type:varchar(500)" json:"website"`
    DateApplied   time.Time `gorm:"type:date" json:"dateApplied"`
}
```

### Status Types
- `SUBMITTED`
- `REJECTED` 
- `PHONE_SCREEN`
- `REMOTE_INTERVIEW`
- `ON_SITE_INTERVIEW`

### Full-Text Search
The app uses SQLite FTS5 for fast full-text search across:
- Company names
- Position titles  
- Notes

Virtual table: `apps_fts` with `content='apps'` and `content_rowid='appId'`

## Building

**Important**: This project requires SQLite with FTS5 support enabled.

### Build with FTS5 Support
```bash
go build -tags "fts5" .
```

### Wails Development
```bash
wails dev -tags "fts5"
```

### Wails Production Build
```bash
wails build -tags "fts5"
```

## Development Setup

1. Install dependencies:
   ```bash
   go mod download
   cd frontend && npm install
   ```

2. Run in development mode:
   ```bash
   wails dev -tags "fts5"
   ```

## Project Structure

```
├── app.go              # Wails app context and methods
├── database.go         # Database initialization and operations
├── main.go            # Application entry point
├── models.go          # Data models and types
├── frontend/          # React frontend
│   ├── src/
│   └── dist/
└── build/             # Build artifacts
```

## Database Operations

- `CreateApp(app *JobApplication)` - Create new job application
- `GetAllApps()` - Retrieve all applications
- `GetAppByID(id uint)` - Get specific application
- `UpdateApp(app *JobApplication)` - Update existing application
- `DeleteApp(id uint)` - Delete application
- `SearchApps(query string)` - Full-text search with FTS5

## Search Examples

```go
// Search for companies containing "Google"
apps, _ := SearchApps("Google")

// Search for remote positions
apps, _ := SearchApps("remote")

// Search for specific technologies in notes
apps, _ := SearchApps("python OR golang")
```

## Notes

- Database file: `job_apps.db`
- FTS5 virtual table automatically stays in sync with main table
- Build tags are required for FTS5 support in SQLite
