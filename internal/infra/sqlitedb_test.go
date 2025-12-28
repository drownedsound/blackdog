package infra

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDb creates an in-memory SQLite database, applies the schema,
// and populates reference data.
func setupTestDb(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// 1. Open In-Memory Database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	// 2. Enable Foreign Keys (Critical for this schema)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// 3. Load Schema and Reference Data
	// Assuming the test runs from 'internal/infra', step back to 'data/scripts'
	basePath := "../../data/scripts"
	scripts := []string{
		"01_create_schema.sql",
		"02_populate_ref_data.sql",
	}

	for _, script := range scripts {
		path := filepath.Join(basePath, script)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read script %s: %v", path, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			t.Fatalf("failed to execute script %s: %v", script, err)
		}
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestSQLiteDb_Insert(t *testing.T) {
	t.Parallel()

	db, cleanup := setupTestDb(t)
	defer cleanup()

	// 1. Initialize Dependencies
	// Use Node 1 for tests.
	idGen, err := NewSnowflakeIDGenerator(1)
	if err != nil {
		t.Fatalf("failed to create id generator: %v", err)
	}

	repo, err := NewSQLiteRepository(db, idGen)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	ctx := context.Background()

	t.Run("Successfully Insert Credit Card Application", func(t *testing.T) {
		now := time.Now().UTC()
		application := &app.Application{
			MemberReferenceNo: "TEST-MRN-001",
			Status:            app.StatusCreated,
			RequestedAmount:   5000000, // 50,000.00
			CreatedAt:         now,
			UpdatedAt:         now,
			CreditCard: app.CreditCard{
				ProfileId:   1, // Pasada King
				CreditLimit: 5000000,
			},
			Applicant: app.Applicant{
				IsPrincipal: true,
				FirstName:   "Juan",
				LastName:    "Dela Cruz",
				MiddleName:  "Santos",
				Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				ContactNumbers: []app.ContactNumber{
					{Type: app.TypeMobile, Value: "09171234567"},
					{Type: app.TypeHome, Value: "0281234567"},
				},
			},
		}

		err := repo.Insert(ctx, application)

		if err != nil {
			t.Fatalf("Insert failed: %v", err)
		}

		if application.Id == 0 {
			t.Error("expected application ID to be generated, got 0")
		}

		var dbStatus int
		var dbReqAmount int
		var dbCCId int64
		err = db.QueryRowContext(ctx, `
			SELECT status_id, requested_amount, credit_card_id 
			FROM APPLICATION WHERE id = ?`, 
			application.Id,
			).Scan(&dbStatus, &dbReqAmount, &dbCCId)

		if err != nil {
			t.Fatalf("failed to query inserted application: %v", err)
		}

		if dbStatus != int(app.StatusCreated) {
			t.Errorf("expected status %d, got %d", app.StatusCreated, dbStatus)
		}
		if dbCCId == 0 {
			t.Error("expected credit_card_id to be set, got 0/NULL")
		}

		var applicantName string
		err = db.QueryRowContext(ctx, `
			SELECT last_name FROM APPLICANT WHERE application_id = ?`, 
			application.Id,
			).Scan(&applicantName)
		if err != nil {
			t.Fatalf("failed to query applicant: %v", err)
		}
		if applicantName != "Dela Cruz" {
			t.Errorf("expected applicant last name 'Dela Cruz', got '%s'", applicantName)
		}
	})

	t.Run("Successfully Insert Personal Loan Application", func(t *testing.T) {
		now := time.Now().UTC()
		application := &app.Application{
			MemberReferenceNo: "TEST-MRN-PL-001",
			Status:            app.StatusCreated,
			RequestedAmount:   10000000,
			CreatedAt:         now,
			UpdatedAt:         now,
			PersonalLoan: app.PersonalLoan{
				ProfileId:  1, 
				LoanAmount: 10000000,
			},
			Applicant: app.Applicant{
				IsPrincipal: true,
				FirstName:   "Maria",
				LastName:    "Clara",
				Birthday:    time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC),
				ContactNumbers: []app.ContactNumber{
					{Type: app.TypeMobile, Value: "09181234567"},
				},
			},
		}

		err := repo.Insert(ctx, application)

		if err != nil {
			t.Fatalf("Insert failed: %v", err)
		}
		if application.Id == 0 {
			t.Error("expected application ID to be generated")
		}

		var dbPLId sql.NullInt64
		var dbCCId sql.NullInt64
		err = db.QueryRowContext(ctx, `
			SELECT personal_loan_id, credit_card_id 
			FROM APPLICATION WHERE id = ?`, 
			application.Id,
			).Scan(&dbPLId, &dbCCId)

		if err != nil {
			t.Fatalf("failed to query inserted application: %v", err)
		}

		if !dbPLId.Valid {
			t.Error("expected personal_loan_id to be valid")
		}
		if dbCCId.Valid {
			t.Error("expected credit_card_id to be NULL for PL application")
		}
	})

	t.Run("Fail on Missing Constraint Data", func(t *testing.T) {
		application := &app.Application{
			MemberReferenceNo: "TEST-FAIL-001",
			Status:            app.StatusCreated,
			RequestedAmount:   5000,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Applicant: app.Applicant{
				IsPrincipal: true,
				LastName:    "Doe",
				FirstName:   "John",
				Birthday:    time.Now(),
				ContactNumbers: []app.ContactNumber{
					{Type: app.TypeMobile, Value: "0000000"},
				},
			},
		}

		err := repo.Insert(ctx, application)

		if err == nil {
			t.Error("expected error due to missing product, got nil")
		}

		if err != nil && !strings.Contains(err.Error(), "constraint") {
			t.Logf("Got error as expected, but check message: %v", err)
		}
	})
}
