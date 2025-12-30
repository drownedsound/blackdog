package infra

import (
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Register SQLite driver
)

func TestNewSQLiteConnection(t *testing.T) {
	// 1. Happy Path: Valid File (WAL mode requires a file-backed DB to persist journaling)
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := NewSQLiteConnection(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteConnection failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	// Verify PRAGMA settings (Mechanical Sympathy: WAL for concurrency)
	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		t.Errorf("Query journal_mode failed: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("Expected journal_mode=wal, got %s", journalMode)
	}

	var busyTimeout int
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&busyTimeout); err != nil {
		t.Errorf("Query busy_timeout failed: %v", err)
	}
	if busyTimeout != 10000 {
		t.Errorf("Expected busy_timeout=10000, got %d", busyTimeout)
	}
}

func TestNewSQLiteConnection_Errors(t *testing.T) {
	// 2. Error Path: Invalid Path (OS permission denied or invalid dir)
	_, err := NewSQLiteConnection("/sys/invalid/path/db.sqlite")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}
