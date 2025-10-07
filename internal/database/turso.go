package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var tursoDb *sql.DB

// InitTurso initializes connection to Turso database
func InitTurso() error {
	url := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if url == "" {
		return fmt.Errorf("TURSO_DATABASE_URL environment variable is not set")
	}

	if authToken == "" {
		return fmt.Errorf("TURSO_AUTH_TOKEN environment variable is not set")
	}

	// Connection string format for Turso
	connStr := fmt.Sprintf("%s?authToken=%s", url, authToken)

	var err error
	tursoDb, err = sql.Open("libsql", connStr)
	if err != nil {
		return fmt.Errorf("failed to open Turso database: %v", err)
	}

	// Test the connection
	if err = tursoDb.Ping(); err != nil {
		return fmt.Errorf("failed to ping Turso database: %v", err)
	}

	return nil
}

// GetTursoDB returns the Turso database instance
func GetTursoDB() *sql.DB {
	return tursoDb
}

// ExecuteTursoSQL executes a raw SQL query on Turso and returns results
func ExecuteTursoSQL(query string, args ...interface{}) ([]map[string]interface{}, error) {
	if tursoDb == nil {
		return nil, fmt.Errorf("turso database not initialized")
	}

	rows, err := tursoDb.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}

// ExecuteTursoExec executes a SQL command (INSERT, UPDATE, DELETE, etc.) on Turso
func ExecuteTursoExec(query string, args ...interface{}) (int64, error) {
	if tursoDb == nil {
		return 0, fmt.Errorf("turso database not initialized")
	}

	result, err := tursoDb.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %v", err)
	}

	return rowsAffected, nil
}

// CloseTurso closes the Turso database connection
func CloseTurso() error {
	if tursoDb != nil {
		return tursoDb.Close()
	}
	return nil
}
