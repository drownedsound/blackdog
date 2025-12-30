package infra

import (
	"context"
	"database/sql"
	"testing"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
	_ "github.com/mattn/go-sqlite3"
)

// -- Mocks --

type mockIdGenerator struct {
	nextId int64
}

func (m *mockIdGenerator) Generate() int64 {
	m.nextId++
	return m.nextId
}

// -- Helpers --

func setupTestDB(t *testing.T) (*sql.DB, *ApplicationRepository) {
	// shared cache allows multiple connections to the same in-memory DB
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared&_foreign_keys=on")
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

	// Inline Schema (Minimal DDL required for Repo to work)
	schema := `
	CREATE TABLE REF_APPLICATION_STATUS (id integer PRIMARY KEY, name text, description text, is_terminal integer);
	CREATE TABLE REF_CURRENCY (id integer PRIMARY KEY, code text);
	CREATE TABLE REF_CREDIT_CARD (id integer PRIMARY KEY, name text, description text, currency_id integer, product_ceiling integer, interest_rate integer);
	CREATE TABLE REF_PERSONAL_LOAN (id integer PRIMARY KEY, name text, description text, currency_id integer, product_ceiling integer, interest_rate integer);
	CREATE TABLE REF_CONTACT_TYPE (id integer PRIMARY KEY, name text, description text);

	CREATE TABLE CREDIT_CARD (id integer PRIMARY KEY, profile_id integer, credit_limit integer);
	CREATE TABLE PERSONAL_LOAN (id integer PRIMARY KEY, profile_id integer, loan_amount integer);

	CREATE TABLE APPLICATION (
	id integer PRIMARY KEY, member_reference_no text, status_id integer, 
	credit_card_id integer, personal_loan_id integer, requested_amount integer, 
	created_at integer, updated_at integer
	);

	CREATE TABLE APPLICANT (
	id integer PRIMARY KEY, application_id integer, is_principal integer, 
	last_name text, first_name text, middle_name text, birthday text
	);

	CREATE TABLE CONTACT_NUMBER (id integer PRIMARY KEY, applicant_id integer, type_id integer, value text);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Schema init failed: %v", err)
	}

	// Populate Reference Data
	refData := `
	INSERT INTO REF_APPLICATION_STATUS VALUES (1, 'CREATED', 'Desc', 0);
	INSERT INTO REF_CURRENCY VALUES (1, 'PHP');
	INSERT INTO REF_CREDIT_CARD VALUES (1, 'Pasada King', 'Desc', 1, 75000, 325);
	INSERT INTO REF_PERSONAL_LOAN VALUES (1, 'Home Loan', 'Desc', 1, 2000000, 1450);
	INSERT INTO REF_CONTACT_TYPE VALUES (1, 'Mobile', 'Desc');
	`
	if _, err := db.Exec(refData); err != nil {
		t.Fatalf("Ref data init failed: %v", err)
	}

	repo, err := NewApplicationRepository(db, &mockIdGenerator{nextId: 100})
	if err != nil {
		t.Fatalf("Repo init failed: %v", err)
	}

	return db, repo
}

// -- Tests --

func TestApplicationRepository_FullLifecycle(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second) // Align with DB integer precision

	// 1. Prepare Complex Entity (Credit Card + Multiple Applicants)
	input := &app.Application{
		MemberReferenceNumber: "REF-CC-001",
		Status:                app.StatusCreated,
		RequestedAmount:       50000,
		CreatedAt:             now,
		UpdatedAt:             now,
		CreditCard: app.CreditCard{
			ProfileId:   1, // Pasada King
			CreditLimit: 10000,
		},
		Applicant: app.Applicant{
			FirstName: "Principal", LastName: "User", IsPrincipal: true,
			Birthday: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			ContactNumbers: []app.ContactNumber{
				{Type: app.TypeMobile, Value: "09170000001"},
			},
		},
		OtherApplicants: []app.Applicant{
			{
				FirstName: "Co", LastName: "Borrower", IsPrincipal: false,
				Birthday: time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC),
				ContactNumbers: []app.ContactNumber{
					{Type: app.TypeMobile, Value: "09170000002"},
				},
			},
		},
	}

	// 2. Insert
	if err := repo.Insert(ctx, input); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if input.Id == 0 {
		t.Error("Expected ID to be generated")
	}

	// 3. Get By ID
	saved, err := repo.GetByInternalId(ctx, input.Id)
	if err != nil {
		t.Fatalf("GetByInternalId failed: %v", err)
	}

	// 4. Verification
	if saved.MemberReferenceNumber != input.MemberReferenceNumber {
		t.Errorf("Want MRN %s, got %s", input.MemberReferenceNumber, saved.MemberReferenceNumber)
	}
	if len(saved.OtherApplicants) != 1 {
		t.Errorf("Want 1 OtherApplicant, got %d", len(saved.OtherApplicants))
	}
	if saved.Applicant.ContactNumbers[0].Value != "09170000001" {
		t.Errorf("Want Principal Contact 09170000001, got %s", saved.Applicant.ContactNumbers[0].Value)
	}
}

func TestApplicationRepository_PersonalLoan(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	input := &app.Application{
		MemberReferenceNumber: "REF-PL-001",
		Status:                app.StatusCreated,
		RequestedAmount:       100000,
		CreatedAt:             time.Now().UTC().Truncate(time.Second),
		UpdatedAt:             time.Now().UTC().Truncate(time.Second),
		PersonalLoan: app.PersonalLoan{
			ProfileId:  1,
			LoanAmount: 100000,
		},
		Applicant: app.Applicant{
			FirstName: "Solo", LastName: "User", IsPrincipal: true,
			Birthday:       time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC),
			ContactNumbers: []app.ContactNumber{{Type: app.TypeMobile, Value: "0999"}},
		},
	}

	if err := repo.Insert(context.Background(), input); err != nil {
		t.Fatalf("Insert PL failed: %v", err)
	}

	saved, err := repo.GetByInternalId(context.Background(), input.Id)
	if err != nil {
		t.Fatalf("Get PL failed: %v", err)
	}

	if saved.PersonalLoan.ProfileId != 1 {
		t.Error("Expected PersonalLoan ProfileId 1")
	}
	if saved.CreditCard.ProfileId != 0 {
		t.Error("Expected no CreditCard")
	}
}

func TestApplicationRepository_Errors(t *testing.T) {
	db, repo := setupTestDB(t)
	defer db.Close()

	// 1. Transaction Failure (Context Cancelled)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	input := &app.Application{MemberReferenceNumber: "FAIL", Id: 500}
	err := repo.Insert(ctx, input)
	if err == nil {
		t.Error("Expected error on cancelled context, got nil")
	}

	// 2. Prepare Failure (Closed DB)
	db2, _ := sql.Open("sqlite3", ":memory:")
	db2.Close()
	_, err = NewApplicationRepository(db2, &mockIdGenerator{})
	if err == nil {
		t.Error("Expected error on NewApplicationRepository with closed DB")
	}

	// 3. Get Not Found
	_, err = repo.GetByInternalId(context.Background(), 999999)
	if err != app.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
