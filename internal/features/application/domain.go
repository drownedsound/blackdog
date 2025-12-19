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
	ErrInvalidCategoryCode = errors.New("category code is not a valid value")
	ErrInvalidStatusCode   = errors.New("status code is not a valid value")
	ErrInvalidId           = errors.New(
		"id cannot be less than or equal to zero",
	)
	ErrInvalidRequestedAmount = errors.New(
		"requested amount cannot be less than or equal to zero",
	)
)

type (
	ApplicationStatus byte
	ProductCategory   byte
)

const (
	StatusUnknown ApplicationStatus = iota
	StatusCreated
	StatusInProgress
	StatusApproved
	StatusDeclined
	StatusCancelled
)

const (
	CategoryUnknown ProductCategory = iota
	CategoryCard
	CategoryLoan
)

func (s ApplicationStatus) String() string {
	switch s {
	case StatusCreated:
		return "CREATED"
	case StatusInProgress:
		return "IN_PROGRESS"
	case StatusApproved:
		return "APPROVED"
	case StatusDeclined:
		return "DECLINED"
	case StatusCancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

func (c ProductCategory) String() string {
	switch c {
	case CategoryCard:
		return "CREDIT_CARD"
	case CategoryLoan:
		return "PERSONAL_LOAN"
	default:
		return "UNKNOWN"
	}
}

// Application is the aggregate root
type Application struct {
	// TODO: Embed Applicant (principal cardholder)
	// TODO: Attach []Applicant (supplementary cardholders)
	CreditCard
	// TODO: Embed Loan
	CreatedAt         time.Time
	UpdatedAt         time.Time
	MemberReferenceNo string // Used as external identifier
	Id                int64  // Used as internal identifier
	RequestedAmount   int    // Shown in centavos
	CategoryCode      ProductCategory
	StatusCode        ApplicationStatus
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

	if a.CategoryCode == CategoryUnknown || a.CategoryCode > CategoryLoan {
		return ErrInvalidCategoryCode
	}

	if a.StatusCode == StatusUnknown || a.StatusCode > StatusCancelled {
		return ErrInvalidStatusCode
	}

	if a.RequestedAmount <= 0 {
		return ErrInvalidRequestedAmount
	}

	if a.CreditLimit <= 1 {
		return errors.New("invalid credit limit")
	}

	return nil
}

type CreditCard struct {
	CardDesignCode string
	CreditLimit int
}

func (c * CreditCard) Validate() error {
	return nil
}

// TODO: Create Loan value object
// TODO: Create Applicant entity
// TODO: Create Contact value object
// TODO: Create Identification value object
// TODO: Create Address value object
// TODO: Create Education value object
// TODO: Create Employment value object
