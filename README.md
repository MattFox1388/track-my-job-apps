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

## Configuration

### Google Drive Backup (Optional)

To enable Google Drive backup functionality, you need to set up OAuth2 credentials:

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Drive API:
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Drive API"
   - Click "Enable"
4. Create OAuth2 credentials:
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Select "Desktop app" as the application type
   - Name it (e.g., "Track My Job Apps")
   - Click "Create"
5. Download the credentials:
   - Click the download button (⬇) next to your newly created OAuth client
   - Save the file as `credentials.json` in the project root directory

On first backup, the app will open your browser for OAuth authorization and save a `token.json` file for future use.

### Turso Database (Optional)

If you want to use Turso as a remote database, set these environment variables:

```bash
# Windows PowerShell
$env:TURSO_DATABASE_URL="libsql://your-database.turso.io"
$env:TURSO_AUTH_TOKEN="your-auth-token-here"

# Windows CMD
set TURSO_DATABASE_URL=libsql://your-database.turso.io
set TURSO_AUTH_TOKEN=your-auth-token-here

# Linux/Mac
export TURSO_DATABASE_URL="libsql://your-database.turso.io"
export TURSO_AUTH_TOKEN="your-auth-token-here"
```

To get your Turso credentials:
1. Sign up at [turso.tech](https://turso.tech)
2. Install the Turso CLI: `curl -sSfL https://get.tur.so/install.sh | bash`
3. Create a database: `turso db create job-apps`
4. Get the URL: `turso db show job-apps --url`
5. Get the token: `turso db tokens create job-apps`

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
