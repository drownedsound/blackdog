package app

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrCreatedAtInFuture  = errors.New("created at cannot be in the future")
	ErrUpdatedAtInFuture  = errors.New("updated at cannot be in the future")
	ErrMissingMemberRefNo = errors.New(
		"member reference number cannot be empty",
	)
	ErrMissingCategoryCode = errors.New("category code cannot be empty")
	ErrInvalidId           = errors.New(
		"id cannot be less than or equal to zero",
	)
	ErrInvalidRequestedAmount = errors.New(
		"requested amount cannot be less than or equal to zero",
	)
)

// Application is the aggregate root
type Application struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	MemberReferenceNo string // external identifier
	CategoryCode      string
	StatusCode        string

	Id              int64 // internal identifier
	RequestedAmount int64 // Shown in centavos
}

func (a *Application) Validate() error {
	now := time.Now()

	if a.CreatedAt.After(now) {
		return ErrCreatedAtInFuture
	}

	if a.UpdatedAt.After(now) {
		return ErrUpdatedAtInFuture
	}

	if strings.TrimSpace(a.MemberReferenceNo) == "" {
		return ErrMissingMemberRefNo
	}

	if strings.TrimSpace(a.CategoryCode) == "" {
		return ErrMissingCategoryCode
	}

	if a.RequestedAmount <= 0 {
		return ErrInvalidRequestedAmount
	}

	return nil
}

// TODO: Create Applicant entity
// TODO: Create Loan value object
// TODO: Create CreditCard value object
// TODO: Create Contact value object
// TODO: Create Address value object
// TODO: Create Identification value object
// TODO: Create Education value object
// TODO: Create Employment value object
