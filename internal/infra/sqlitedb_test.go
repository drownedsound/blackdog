package infra

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
	_ "github.com/mattn/go-sqlite3"
)

// minimalSchema mirrors the production schema required for these tests.
// Keeping this local ensures tests are atomic and don't break
// if external SQL files move.
const minimalSchema = `
	PRAGMA foreign_keys = ON;
	PRAGMA journal_mode = WAL;

	CREATE TABLE REF_PRODUCT_CATEGORY (
	id INTEGER PRIMARY KEY, 
	name TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL
	) STRICT;

	CREATE TABLE REF_APP_STATUS (
	id INTEGER PRIMARY KEY, 
	name TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL,
	is_terminal INTEGER NOT NULL DEFAULT 0 CHECK (is_terminal IN (0, 1))
	) STRICT;

	CREATE TABLE APPLICATION (
	id INTEGER PRIMARY KEY, 
	member_reference_no TEXT NOT NULL UNIQUE,
	category_id INTEGER NOT NULL,
	status_id INTEGER NOT NULL,
	requested_amount INTEGER NOT NULL,
	created_at INTEGER NOT NULL DEFAULT (unixepoch()),
	updated_at INTEGER NOT NULL DEFAULT (unixepoch()),

	FOREIGN KEY (category_id) REFERENCES REF_PRODUCT_CATEGORY(id),
	FOREIGN KEY (status_id) REFERENCES REF_APP_STATUS(id)
	CHECK (requested_amount > 0)
	) STRICT;

	-- Seed Reference Data
	INSERT INTO REF_PRODUCT_CATEGORY (id, name, description) 
	VALUES (1, 'Credit Card', 'CC');
	
	INSERT INTO REF_APP_STATUS (id, name, description) 
	VALUES (1, 'CREATED', 'New');
	`

// setupTestDb initializes a fresh SQLite DB in a temp directory
// for valid isolation.
func setupTestDb(t *testing.T) (*SqliteDb, func()) {
	t.Helper()

	// Use t.TempDir for isolation. On Linux, this is often tmpfs (RAM),
	// providing high performance while still testing file I/O mechanics.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_blackdog.db")

	db, err := NewSQLiteConnection(dbPath)
	if err != nil {
		t.Fatalf("Failed to open connection: %v", err)
	}

	// Initialize Schema
	if _, err := db.Exec(minimalSchema); err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	repo := NewSqliteDb(db)

	cleanup := func() {
		db.Close()
	}

	return repo, cleanup
}

func TestSqliteDb_Save(t *testing.T) {
	// Parallel execution to detect any WAL-mode locking issues
	t.Parallel()

	repo, cleanup := setupTestDb(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("Successfully Persists Application", func(t *testing.T) {
		newApp := &app.Application{
			MemberReferenceNo: "REF-001",
			CategoryCode:      app.CategoryCard,
			StatusCode:        app.StatusCreated,
			RequestedAmount:   500000,
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		}

		err := repo.Save(ctx, newApp)
		if err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		if newApp.Id == 0 {
			t.Error("Expected ID to be set after Save(), got 0")
		}
	})

	t.Run("Fails On Duplicate MemberReferenceNo", func(t *testing.T) {
		// Attempt to insert the same Reference No again
		dupApp := &app.Application{
			MemberReferenceNo: "REF-001", // Already exists from previous test
			CategoryCode:      app.CategoryCard,
			StatusCode:        app.StatusCreated,
			RequestedAmount:   1000,
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		}

		err := repo.Save(ctx, dupApp)
		if err == nil {
			t.Error("Expected error on duplicate insert, got nil")
		}

		// Verify we are wrapping the error correctly
		expectedMsg := "duplicate application"
		if err != nil && len(err.Error()) < len(expectedMsg) {
			t.Errorf(
				"Expected error containing %q, got %q",
				expectedMsg,
				err.Error(),
			)
		}
	})

	t.Run("Fails When Context Canceled", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		a := &app.Application{
			MemberReferenceNo: "REF-CANCEL",
			CategoryCode:      app.CategoryCard,
			StatusCode:        app.StatusCreated,
			RequestedAmount:   100,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err := repo.Save(canceledCtx, a)
		if err == nil {
			t.Error("Expected error on canceled context, got nil")
		}
	})
}

func TestSqliteDb_GetById(t *testing.T) {
	t.Parallel()
	repo, cleanup := setupTestDb(t)
	defer cleanup()

	ctx := context.Background()

	// Seed a record
	savedApp := &app.Application{
		MemberReferenceNo: "GET-001",
		CategoryCode:      app.CategoryCard,
		StatusCode:        app.StatusCreated,
		RequestedAmount:   75000,
		CreatedAt:         time.Now().UTC().Truncate(time.Second),
		UpdatedAt:         time.Now().UTC().Truncate(time.Second),
	}
	if err := repo.Save(ctx, savedApp); err != nil {
		t.Fatalf("Setup failed, could not save app: %v", err)
	}

	t.Run("Successfully Retrieves Application", func(t *testing.T) {
		fetched, err := repo.GetById(ctx, savedApp.Id)
		if err != nil {
			t.Fatalf("GetById() failed: %v", err)
		}

		if fetched.Id != savedApp.Id {
			t.Errorf("Expected ID %d, got %d", savedApp.Id, fetched.Id)
		}
		if fetched.MemberReferenceNo != savedApp.MemberReferenceNo {
			t.Errorf(
				"Expected Ref %s, got %s",
				savedApp.MemberReferenceNo,
				fetched.MemberReferenceNo)
		}
		// Verify Enum Mapping
		if fetched.CategoryCode != app.CategoryCard {
			t.Errorf("Expected CategoryCard, got %v", fetched.CategoryCode)
		}
		// Verify Timestamp Rehydration
		if !fetched.CreatedAt.Equal(savedApp.CreatedAt) {
			t.Errorf(
				"Expected CreatedAt %v, got %v",
				savedApp.CreatedAt,
				fetched.CreatedAt)
		}
	})

	t.Run("Returns ErrNotFound For NonExistent ID", func(t *testing.T) {
		_, err := repo.GetById(ctx, 999999)
		if err != app.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}

func TestNewSQLiteConnection_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	conflictPath := filepath.Join(tmpDir, "bad_db")

	// Create a directory at the path where we want the DB file to be
	if err := os.Mkdir(conflictPath, 0o755); err != nil {
		t.Fatal(err)
	}

	// This should fail because 'bad_db' is a directory, not a file
	_, err := NewSQLiteConnection(conflictPath)
	if err == nil {
		t.Error("Expected error opening DB on top of directory, got nil")
	}
}

func TestSqliteDb_Save_EdgeCases(t *testing.T) {
	t.Parallel()
	repo, cleanup := setupTestDb(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("Fails On Check Constraint Violation", func(t *testing.T) {
		// The schema has CHECK (requested_amount > 0)
		invalidApp := &app.Application{
			MemberReferenceNo: "NEG-AMOUNT",
			CategoryCode:      app.CategoryCard,
			StatusCode:        app.StatusCreated,
			RequestedAmount:   -500, // Invalid
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		}

		err := repo.Save(ctx, invalidApp)
		if err == nil {
			t.Error("Expected error due to CHECK constraint, got nil")
		}

		if !errors.Is(err, app.ErrInsertFailed) {
			t.Errorf("Expected ErrInsertFailed, got %v", err)
		}
	})
}

func TestSqliteDb_GetById_EdgeCases(t *testing.T) {
	t.Parallel()
	repo, cleanup := setupTestDb(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("Fails When Context Canceled", func(t *testing.T) {
		// Cancel the context immediately
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := repo.GetById(canceledCtx, 1)
		if err == nil {
			t.Error("Expected error due to canceled context, got nil")
		}

		// Verify it hits the generic error path, not ErrNotFound
		if err == app.ErrNotFound {
			t.Error("Expected generic query error, got ErrNotFound")
		}
	})
}
