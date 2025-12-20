package app

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidCardProfileCode = errors.New(
		"card profile code is not a valid value",
	)
	ErrInvalidCreditLimit = errors.New(
		"credit limit cannot be less than or equal to zero",
	)
	ErrInvalidInterestRate = errors.New(
		"interest rate cannot be less than or equal to zero",
	)
	ErrCreatedAtInFuture  = errors.New("created at cannot be in the future")
	ErrUpdatedAtInFuture  = errors.New("updated at cannot be in the future")
	ErrMissingMemberRefNo = errors.New(
		"member reference number cannot be empty",
	)
	// ErrInvalidCategoryCode = errors.New("category code is not a valid value")
	// ErrInvalidId           = errors.New(
	// 	"id cannot be less than or equal to zero",
	// )
	ErrInvalidRequestedAmount = errors.New(
		"requested amount cannot be less than or equal to zero",
	)
	ErrInvalidStatusCode = errors.New("status code is not a valid value")
)

type (
	ApplicationStatus byte

// ProductCategory   byte
)

const (
	StatusUnknown ApplicationStatus = iota
	StatusCreated
	StatusInProgress
	StatusApproved
	StatusDeclined
	StatusCancelled
)

// const (
// 	CategoryUnknown ProductCategory = iota
// 	CategoryCard
// 	CategoryLoan
// )

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

// func (c ProductCategory) String() string {
// 	switch c {
// 	case CategoryCard:
// 		return "CREDIT_CARD"
// 	case CategoryLoan:
// 		return "PERSONAL_LOAN"
// 	default:
// 		return "UNKNOWN"
// 	}
// }

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
	// CategoryCode      ProductCategory
	StatusCode ApplicationStatus
}

func (a *Application) Validate() error {
	if strings.TrimSpace(a.CardProfileCode) == "" {
		return ErrInvalidCardProfileCode
	}

	if a.CreditLimit <= 1 {
		return ErrInvalidCreditLimit
	}

	if a.InterestRate <= 0 {
		return ErrInvalidInterestRate
	}

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

	// if a.CategoryCode == CategoryUnknown || a.CategoryCode > CategoryLoan {
	// 	return ErrInvalidCategoryCode
	// }

	if a.RequestedAmount <= 0 {
		return ErrInvalidRequestedAmount
	}

	if a.StatusCode == StatusUnknown || a.StatusCode > StatusCancelled {
		return ErrInvalidStatusCode
	}

	return nil
}

type CreditCard struct {
	CardProfileCode string
	CreditLimit     int
	// InterestRate represents the rate in basis points
	// Example: 1 bps == 0.01% or 1250 bps == 12.50%
	InterestRate int
}

// TODO: Create Loan value object
// TODO: Create Applicant entity
// TODO: Create Contact value object
// TODO: Create Identification value object
// TODO: Create Address value object
// TODO: Create Education value object
// TODO: Create Employment value object
