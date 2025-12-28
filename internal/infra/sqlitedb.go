package infra

import (
	"context"
	"database/sql"
	"fmt"

	app "github.com/drownedsound/blackdog/internal/features/application"
)

type SQLiteDB struct {
	db     *sql.DB
	idGen  IdGenerator 

	stmtInsertApplication      *sql.Stmt
	stmtInsertCreditCard       *sql.Stmt
	stmtInsertPersonalLoan        *sql.Stmt
	stmtInsertApplicant *sql.Stmt
	stmtInsertContact   *sql.Stmt
}

func NewSQLiteDB(db *sql.DB, idGen IdGenerator) (
	repo *SQLiteDB, 
	err error,
) {
	repo = &SQLiteDB{
		db:    db,
		idGen: idGen,
	}

	if repo.stmtInsertApplication, err = db.Prepare(`
		INSERT INTO APPLICATION (
			id, member_reference_no, status_id, requested_amount, 
			credit_card_id, personal_loan_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`); err != nil {
		return nil, fmt.Errorf("prepare insert app: %w", err)
	}

	if repo.stmtInsertCreditCard, err = db.Prepare(`
		INSERT INTO CREDIT_CARD (id, profile_id, credit_limit) 
		VALUES (?, ?, ?)`); err != nil {
		return nil, fmt.Errorf("prepare insert cc: %w", err)
	}

	if repo.stmtInsertPersonalLoan, err = db.Prepare(`
		INSERT INTO PERSONAL_LOAN (id, profile_id, loan_amount) 
		VALUES (?, ?, ?)`); err != nil {
		return nil, fmt.Errorf("prepare insert pl: %w", err)
	}

	if repo.stmtInsertApplicant, err = db.Prepare(`
		INSERT INTO APPLICANT (
			id, application_id, is_principal, 
			last_name, first_name, middle_name, birthday
		) VALUES (?, ?, ?, ?, ?, ?, ?)`); err != nil {
		return nil, fmt.Errorf("prepare insert applicant: %w", err)
	}

	if repo.stmtInsertContact, err = db.Prepare(`
		INSERT INTO CONTACT_NUMBER (id, applicant_id, type_id, value) 
		VALUES (?, ?, ?, ?)`); err != nil {
		return nil, fmt.Errorf("prepare insert contact: %w", err)
	}

	return repo nil
}

func (r *SQLiteDB) Insert(ctx context.Context, a *app.Application) (err error) {
	// 0. Pre-Generate the Aggregate Root Id
	a.Id = r.idGen.Generate()

	// 1. Start Transaction
	tx, err = r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// 2. Bind Prepared Statements
	txStmtApp := tx.StmtContext(ctx, r.stmtInsertApplication)
	txStmtCC := tx.StmtContext(ctx, r.stmtInsertCreditCard)
	txStmtPL := tx.StmtContext(ctx, r.stmtInsertPersonalLoan)
	txStmtApplicant := tx.StmtContext(ctx, r.stmtInsertApplicant)
	txStmtContact := tx.StmtContext(ctx, r.stmtInsertContact)

	defer func() {
		txStmtApp.Close()
		txStmtCC.Close()
		txStmtPL.Close()
		txStmtApplicant.Close()
		txStmtContact.Close()
	}()

	// 3. Insert Product
	var ccId, plId sql.NullInt64

	if a.CreditCard.ProfileId > 0 {
		newCCId := r.idGen.Generate()
		_, err := txStmtCC.ExecContext(
			ctx, 
			newCCId, 
			a.CreditCard.ProfileId, 
			a.CreditCard.CreditLimit,
			)
		if err != nil {
			return fmt.Errorf("insert credit card: %w", err)
		}
		ccId.Int64 = newCCId
		ccId.Valid = true
	} else if a.PersonalLoan.ProfileId > 0 {
		// Generate Id for the Personal Loan entity
		newPLId := r.idGen.Generate()
		_, err := txStmtPL.ExecContext(
			ctx, 
			newPLId, 
			a.PersonalLoan.ProfileId, 
			a.PersonalLoan.LoanAmount,
			)
		if err != nil {
			return fmt.Errorf("insert personal loan: %w", err)
		}
		plId.Int64 = newPLId
		plId.Valid = true
	}

	// 4. Insert Application (Using the pre-generated a.Id)
	_, err = txStmtApp.ExecContext(ctx,
		a.Id, 
		a.MemberReferenceNo,
		a.Status,
		a.RequestedAmount,
		ccId,
		plId,
		a.CreatedAt.Unix(),
		a.UpdatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert application: %w", err)
	}

	// 5. Insert Principal Applicant
	// We pass a.Id as the foreign key
	if err := r.insertApplicantWithContacts(
		ctx, 
		txStmtApplicant, 
		txStmtContact, 
		a.Id, 
		&a.Applicant, 
		true,
		); err != nil {
		return fmt.Errorf("insert principal: %w", err)
	}

	// 6. Insert Other Applicants
	for i := range a.OtherApplicants {
		if err := r.insertApplicantWithContacts(
			ctx, 
			txStmtApplicant, 
			txStmtContact, 
			a.Id, 
			&a.OtherApplicants[i], false,
			); err != nil {
			return fmt.Errorf("insert other applicant %d: %w", i, err)
		}
	}

	// 7. Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return 
}

// insertApplicantWithContacts helper
func (r *SQLiteRepository) insertApplicantWithContacts(
	ctx context.Context,
	stmtApplicant *sql.Stmt,
	stmtContact *sql.Stmt,
	appId int64,
	applicant *app.Applicant,
	isPrincipal bool,
) (err error) {
	// Generate Id for Applicant
	applicantId := r.idGen.Generate()
	
	isPrincipalInt := 0
	if isPrincipal {
		isPrincipalInt = 1
	}

	_, err = stmtApplicant.ExecContext(ctx,
		applicantId, // <--- Explicit Id
		appId,       // Foreign Key
		isPrincipalInt,
		applicant.LastName,
		applicant.FirstName,
		applicant.MiddleName,
		applicant.Birthday.Format("2006-01-02"),
	)
	if err != nil {
		return err
	}

	// Insert Contact Numbers
	for _, contact := range applicant.ContactNumbers {
		// Generate Id for Contact
		contactId := r.idGen.Generate()
		
		_, err := stmtContact.ExecContext(ctx,
			contactId, // <--- Explicit Id
			applicantId,
			contact.Type,
			contact.Value,
		)
		if err != nil {
			return fmt.Errorf("insert contact %s: %w", contact.Value, err)
		}
	}

	return 
}
