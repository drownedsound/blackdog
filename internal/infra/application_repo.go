package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
)

type ApplicationRepository struct {
	db                      *sql.DB
	idGen                   IdGenerator
	stmtInsertApplication   *sql.Stmt
	stmtInsertCreditCard    *sql.Stmt
	stmtInsertPersonalLoan  *sql.Stmt
	stmtInsertApplicant     *sql.Stmt
	stmtInsertContactNumber *sql.Stmt
}

func NewApplicationRepository(db *sql.DB, idGen IdGenerator) (
	*ApplicationRepository,
	error,
) {
	repo := &ApplicationRepository{
		db:    db,
		idGen: idGen,
	}

	prepare := func(query string) (*sql.Stmt, error) {
		return db.Prepare(query)
	}

	var err error

	repo.stmtInsertApplication, err = prepare(`
	INSERT INTO APPLICATION (
		id, member_reference_no, status_id, requested_amount, 
		credit_card_id, personal_loan_id, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf(
			"preparing insert application stmt failed: %w",
			err,
		)
	}

	repo.stmtInsertCreditCard, err = prepare(`
	INSERT INTO CREDIT_CARD (
		id, profile_id, credit_limit
	) VALUES (?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("preparing insert cc stmt failed: %w", err)
	}

	repo.stmtInsertPersonalLoan, err = prepare(`
	INSERT INTO PERSONAL_LOAN (
		id, profile_id, loan_amount
	) VALUES (?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("preparing insert pl stmt failed: %w", err)
	}

	repo.stmtInsertApplicant, err = prepare(`
	INSERT INTO APPLICANT (
		id, application_id, is_principal, 
		last_name, first_name, middle_name, birthday
	) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf(
			"preparing insert applicant stmt failed: %w", err,
		)
	}

	repo.stmtInsertContactNumber, err = prepare(`
	INSERT INTO CONTACT_NUMBER (
		id, applicant_id, type_id, value
	) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("preparing insert contact stmt failed: %w", err)
	}

	return repo, nil
}

func (r *ApplicationRepository) Insert(
	ctx context.Context,
	a *app.Application,
) error {
	a.Id = r.idGen.Generate()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	txStmtApp := tx.StmtContext(ctx, r.stmtInsertApplication)
	defer txStmtApp.Close()

	var txStmtCC *sql.Stmt
	if a.CreditCard.ProfileId > 0 {
		txStmtCC = tx.StmtContext(ctx, r.stmtInsertCreditCard)
		defer txStmtCC.Close()
	}

	var txStmtPL *sql.Stmt
	if a.PersonalLoan.ProfileId > 0 {
		txStmtPL = tx.StmtContext(ctx, r.stmtInsertPersonalLoan)
		defer txStmtPL.Close()
	}

	txStmtApplicant := tx.StmtContext(ctx, r.stmtInsertApplicant)
	defer txStmtApplicant.Close()

	txStmtContact := tx.StmtContext(ctx, r.stmtInsertContactNumber)
	defer txStmtContact.Close()

	var ccId, plId sql.NullInt64

	if a.CreditCard.ProfileId > 0 {
		ccId = sql.NullInt64{Int64: r.idGen.Generate(), Valid: true}

		if _, err := txStmtCC.ExecContext(
			ctx, ccId, a.CreditCard.ProfileId, a.CreditCard.CreditLimit,
		); err != nil {
			return fmt.Errorf("insert credit card failed: %w", err)
		}
	}

	if a.PersonalLoan.ProfileId > 0 {
		plId = sql.NullInt64{Int64: r.idGen.Generate(), Valid: true}

		if _, err := txStmtPL.ExecContext(
			ctx, plId, a.PersonalLoan.ProfileId, a.PersonalLoan.LoanAmount,
		); err != nil {
			return fmt.Errorf("insert personal_loan failed: %w", err)
		}
	}

	if _, err := txStmtApp.ExecContext(
		ctx,
		a.Id,
		a.MemberReferenceNumber,
		a.Status,
		a.RequestedAmount,
		ccId,
		plId,
		a.CreatedAt.Unix(),
		a.UpdatedAt.Unix(),
	); err != nil {
		return fmt.Errorf("insert application failed: %w", err)
	}

	if err := r.insertApplicant(
		ctx,
		txStmtApplicant,
		txStmtContact,
		&a.Applicant,
		a.Id,
	); err != nil {
		return fmt.Errorf("insert principal failed: %w", err)
	}

	for i := range a.OtherApplicants {
		if err := r.insertApplicant(
			ctx,
			txStmtApplicant,
			txStmtContact,
			&a.OtherApplicants[i],
			a.Id,
		); err != nil {
			return fmt.Errorf("insert other applicant %d failed: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx failed: %w", err)
	}

	return nil
}

func (r *ApplicationRepository) insertApplicant(
	ctx context.Context,
	stmtAppl *sql.Stmt,
	stmtContact *sql.Stmt,
	appl *app.Applicant,
	appId int64,
) error {
	applicantId := r.idGen.Generate()

	isPrincipal := 0
	if appl.IsPrincipal {
		isPrincipal = 1
	}

	if _, err := stmtAppl.ExecContext(
		ctx,
		applicantId,
		appId,
		isPrincipal,
		appl.LastName,
		appl.FirstName,
		appl.MiddleName,
		appl.Birthday.Format(time.DateOnly),
	); err != nil {
		return err
	}

	for i, contact := range appl.ContactNumbers {
		contactId := r.idGen.Generate()
		if _, err := stmtContact.ExecContext(
			ctx,
			contactId,
			applicantId,
			contact.Type,
			contact.Value,
		); err != nil {
			return fmt.Errorf("contact %d: %w", i, err)
		}
	}

	return nil
}
