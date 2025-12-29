package infra

import (
	"context"
	"database/sql"
	"fmt"

	app "github.com/drownedsound/blackdog/internal/features/application"
)

type ApplicationRepository struct {
	db    *sql.DB
	idGen IdGenerator

	stmtInsertApplication  *sql.Stmt
	stmtInsertCreditCard   *sql.Stmt
	stmtInsertPersonalLoan *sql.Stmt
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

	var ccId, plId sql.NullInt64

	if a.CreditCard.ProfileId > 0 {
		id := r.idGen.Generate()
		ccId = sql.NullInt64{Int64: id, Valid: true}

		stmt := tx.StmtContext(ctx, r.stmtInsertCreditCard)
		defer stmt.Close()

		_, err := stmt.ExecContext(ctx, ccId, a.CreditCard.ProfileId, a.CreditCard.CreditLimit)
		if err != nil {
			return fmt.Errorf("insert credit_card failed: %w", err)
		}
	}

	if a.PersonalLoan.ProfileId > 0 {
		id := r.idGen.Generate()
		plId = sql.NullInt64{Int64: id, Valid: true}

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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx failed: %w", err)
	}

	return nil
}

func (r *ApplicationRepository) Close() error {
	var errs []error

	if r.stmtInsertApplication != nil {
		if err := r.stmtInsertApplication.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if r.stmtInsertCreditCard != nil {
		if err := r.stmtInsertCreditCard.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if r.stmtInsertPersonalLoan != nil {
		if err := r.stmtInsertPersonalLoan.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close one or more statements: %v", errs)
	}

	return nil
}
