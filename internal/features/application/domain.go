package app

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidCardProfile = errors.New(
		"card profile is not a valid value",
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
	ErrInvalidId = errors.New(
		"id cannot be less than or equal to zero",
	)
	ErrInvalidRequestedAmount = errors.New(
		"requested amount cannot be less than or equal to zero",
	)
	ErrInvalidStatus = errors.New("status is not a valid value")
)

type ApplicationStatus byte

const (
	StatusUnknown ApplicationStatus = iota
	StatusCreated
	StatusInProgress
	StatusApproved
	StatusDeclined
	StatusCancelled
)

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
	Status            ApplicationStatus
	// TODO: Add ApprovedAmount
}

func (a *Application) Validate() error {
	if a.CardProfile <= 0 {
		return ErrInvalidCardProfile
	}

	// Handles CreditLimit zero value passed by service via DTO
	if a.CreditLimit != 0 && a.CreditLimit < 1 {
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

	if a.RequestedAmount <= 0 {
		return ErrInvalidRequestedAmount
	}

	if a.Status == StatusUnknown || a.Status > StatusCancelled {
		return ErrInvalidStatus
	}

	return nil
}

type CreditCard struct {
	CardProfile int64
	CreditLimit int
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
