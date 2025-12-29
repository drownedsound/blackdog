package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
)

type ApplicationRepository struct {
	db    *sql.DB
	idGen IdGenerator

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

	var err error

	queryApp := `
		INSERT INTO APPLICATION (
			id, member_reference_no, status_id, requested_amount, 
			credit_card_id, personal_loan_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if repo.stmtInsertApplication, err = db.Prepare(queryApp); err != nil {
		return nil, fmt.Errorf(
			"preparing insert application stmt failed: %w",
			err,
		)
	}

	queryCC := `
		INSERT INTO CREDIT_CARD (
			id, profile_id, credit_limit
		) VALUES (?, ?, ?)`
	if repo.stmtInsertCreditCard, err = db.Prepare(queryCC); err != nil {
		return nil, fmt.Errorf("preparing insert cc stmt failed: %w", err)
	}

	queryPL := `
		INSERT INTO PERSONAL_LOAN (
			id, profile_id, loan_amount
		) VALUES (?, ?, ?)`
	if repo.stmtInsertPersonalLoan, err = db.Prepare(queryPL); err != nil {
		return nil, fmt.Errorf("preparing insert pl stmt failed: %w", err)
	}

	queryAppl := `
		INSERT INTO APPLICANT (
			id, application_id, is_principal, last_name, first_name, middle_name, birthday
		) VALUES (?, ?, ?, ?, ?, ?, ?)`
	if repo.stmtInsertApplicant, err = db.Prepare(queryAppl); err != nil {
		return nil, fmt.Errorf("preparing insert applicant stmt failed: %w", err)
	}

	queryContact := `
		INSERT INTO CONTACT_NUMBER (
			id, applicant_id, type_id, value
		) VALUES (?, ?, ?, ?)`

	if repo.stmtInsertContactNumber, err = db.Prepare(queryContact); err != nil {
		return nil, fmt.Errorf("preparing insert contact stmt failed: %w", err)
	}

	return repo, nil
}

func (r *ApplicationRepository) Insert(
	ctx context.Context,
	a *app.Application,
) error {
	a.Id = r.idGen.Generate()
	// FIXME: Map to DTO

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	var ccId, plId sql.NullInt64

	if a.CreditCard.ProfileId > 0 {
		ccId = sql.NullInt64{Int64: r.idGen.Generate(), Valid: true}

		stmt := tx.StmtContext(ctx, r.stmtInsertCreditCard)
		defer stmt.Close()

		_, err := stmt.ExecContext(ctx, ccId, a.CreditCard.ProfileId, a.CreditCard.CreditLimit)
		if err != nil {
			return fmt.Errorf("insert credit_card failed: %w", err)
		}
	}

	if a.PersonalLoan.ProfileId > 0 {
		plId = sql.NullInt64{Int64: r.idGen.Generate(), Valid: true}

		stmt := tx.StmtContext(ctx, r.stmtInsertPersonalLoan)
		defer stmt.Close()

		_, err := stmt.ExecContext(ctx, plId, a.PersonalLoan.ProfileId, a.PersonalLoan.LoanAmount)
		if err != nil {
			return fmt.Errorf("insert personal_loan failed: %w", err)
		}
	}

	stmtApp := tx.StmtContext(ctx, r.stmtInsertApplication)
	defer stmtApp.Close()

	_, err = stmtApp.ExecContext(
		ctx,
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
		return fmt.Errorf("insert application failed: %w", err)
	}

	stmtAppl := tx.StmtContext(ctx, r.stmtInsertApplicant)
	defer stmtAppl.Close()

	applicantId := r.idGen.Generate()

	_, err = stmtAppl.ExecContext(
		ctx,
		applicantId,
		a.Id,
		1,
		a.LastName,
		a.FirstName,
		a.MiddleName,
		a.Birthday.Format(time.DateOnly),
	)
	if err != nil {
		return fmt.Errorf("insert applicant failed: %w", err)
	}

	if len(a.ContactNumbers) > 0 {
		stmtContact := tx.StmtContext(ctx, r.stmtInsertContactNumber)
		defer stmtContact.Close()

		for i, contact := range a.ContactNumbers {
			contactId := r.idGen.Generate()

			_, err = stmtContact.ExecContext(
				ctx,
				contactId,
				applicantId,
				contact.Type,
				contact.Value,
			)
			if err != nil {
				return fmt.Errorf("insert contact number %d failed: %w", i, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx failed: %w", err)
	}

	return nil
}

// func (r *ApplicationRepository) Close() error {
// 	var errs []error
//
// 	if r.stmtInsertApplication != nil {
// 		if err := r.stmtInsertApplication.Close(); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}
// 	if r.stmtInsertCreditCard != nil {
// 		if err := r.stmtInsertCreditCard.Close(); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}
// 	if r.stmtInsertPersonalLoan != nil {
// 		if err := r.stmtInsertPersonalLoan.Close(); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}
//
// 	if len(errs) > 0 {
// 		return fmt.Errorf("failed to close one or more statements: %v", errs)
// 	}
//
// 	return nil
// }
