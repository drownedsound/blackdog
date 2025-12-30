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
	stmtGetApplicationById  *sql.Stmt
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

	repo.stmtGetApplicationById, err = prepare(`
	SELECT 
		app.id, app.member_reference_no, app.status_id, app.requested_amount, 
		app.created_at, app.updated_at, cc.profile_id, cc.credit_limit, 
		rcc.interest_rate, rcc.currency_id, pl.profile_id, pl.loan_amount, 
		rpl.interest_rate, rpl.currency_id, a.id, a.is_principal, 
		a.first_name, a.middle_name, a.last_name, a.birthday, c.value, c.type_id
	FROM APPLICATION app
	LEFT JOIN CREDIT_CARD cc ON app.credit_card_id = cc.id
	LEFT JOIN REF_CREDIT_CARD rcc ON cc.profile_id = rcc.id
	LEFT JOIN PERSONAL_LOAN pl ON app.personal_loan_id = pl.id
	LEFT JOIN REF_PERSONAL_LOAN rpl ON pl.profile_id = rpl.id
	JOIN APPLICANT a ON app.id = a.application_id
	LEFT JOIN CONTACT_NUMBER c ON a.id = c.applicant_id
	WHERE app.id = ?
	ORDER BY a.is_principal DESC, a.id
	`)
	if err != nil {
		return nil, fmt.Errorf("preparing get application stmt failed: %w", err)
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

func (r *ApplicationRepository) GetByInternalId(
	ctx context.Context,
	id int64,
) (*app.Application, error) {
	rows, err := r.stmtGetApplicationById.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get application failed: %w", err)
	}
	defer rows.Close()

	var (
		application        *app.Application
		currentApplicant   *app.Applicant
		currentApplicantId int64
	)

	for rows.Next() {
		var (
			// Application fields
			appId, statusId, reqAmount int64
			createdAt, updatedAt       int64
			mrn                        string

			// Credit Card fields (Nullable)
			ccProfileId, ccLimit, ccRate, ccCurrency sql.NullInt64

			// Personal Loan fields (Nullable)
			plProfileId, plAmount, plRate, plCurrency sql.NullInt64

			// Applicant fields
			applId, isPrincipal          int64
			applFirst, applLast, applDob string
			applMiddle                   sql.NullString

			// Contact fields (Nullable)
			contactValue sql.NullString
			contactType  sql.NullInt64
		)

		if err := rows.Scan(
			&appId, &mrn, &statusId, &reqAmount, &createdAt, &updatedAt,
			&ccProfileId, &ccLimit, &ccRate, &ccCurrency,
			&plProfileId, &plAmount, &plRate, &plCurrency,
			&applId, &isPrincipal, &applFirst, &applMiddle, &applLast, &applDob,
			&contactValue, &contactType,
		); err != nil {
			return nil, fmt.Errorf("database row scan failed: %w", err)
		}

		// Initialize Application on first row
		if application == nil {
			application = &app.Application{
				Id:                    appId,
				MemberReferenceNumber: mrn,
				Status:                app.ApplicationStatus(statusId),
				RequestedAmount:       int32(reqAmount),
				CreatedAt:             time.Unix(createdAt, 0).UTC(),
				UpdatedAt:             time.Unix(updatedAt, 0).UTC(),
				// Initialize slice to avoid nil slice issues
				OtherApplicants: make([]app.Applicant, 0, 1),
			}

			if ccProfileId.Valid {
				application.CreditCard = app.CreditCard{
					ProfileId:    ccProfileId.Int64,
					CreditLimit:  int32(ccLimit.Int64),
					InterestRate: int16(ccRate.Int64),
					CurrencyId:   ccCurrency.Int64,
				}
			}

			if plProfileId.Valid {
				application.PersonalLoan = app.PersonalLoan{
					ProfileId:    plProfileId.Int64,
					LoanAmount:   int32(plAmount.Int64),
					InterestRate: int16(plRate.Int64),
					CurrencyId:   plCurrency.Int64,
				}
			}
		}

		// Handle Applicant switching
		if applId != currentApplicantId {
			// Save the finished applicant to the application struct
			if currentApplicant != nil {
				if currentApplicant.IsPrincipal {
					application.Applicant = *currentApplicant
				} else {
					application.OtherApplicants = append(
						application.OtherApplicants, *currentApplicant,
					)
				}
			}

			// Start new applicant
			dob, err := time.Parse(time.DateOnly, applDob)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid dob format %s: %w", applDob, err,
				)
			}

			currentApplicant = &app.Applicant{
				FirstName:      applFirst,
				LastName:       applLast,
				Birthday:       dob,
				IsPrincipal:    isPrincipal == 1,
				ContactNumbers: make([]app.ContactNumber, 0, 2),
			}
			if applMiddle.Valid {
				currentApplicant.MiddleName = applMiddle.String
			}
			currentApplicantId = applId
		}

		// Add Contact Number
		if contactValue.Valid {
			currentApplicant.ContactNumbers = append(
				currentApplicant.ContactNumbers,
				app.ContactNumber{
					Value: contactValue.String,
					Type:  app.ContactNumberType(contactType.Int64),
				},
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	// No rows found
	if application == nil {
		return nil, app.ErrNotFound
	}

	// Save the very last applicant (since loop finishes before saving it)
	if currentApplicant != nil {
		if currentApplicant.IsPrincipal {
			application.Applicant = *currentApplicant
		} else {
			application.OtherApplicants = append(
				application.OtherApplicants, *currentApplicant,
			)
		}
	}

	return application, nil
}
