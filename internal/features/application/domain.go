package app

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidProduct = errors.New(
		"product applied for is not valid",
		)
	ErrInvalidCurrency = errors.New(
		"currency used cannot be less than zero"	
		)
	ErrInvalidCreditLimit = errors.New(
		"credit limit cannot be less than zero",
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
	CreatedAt         time.Time
	UpdatedAt         time.Time
	OtherApplicants []Applicant
	MemberReferenceNo string // Used as external identifier
	Applicant
	CreditCard CreditCard
	PersonalLoan PersonalLoan
	Id                int64  // Used as internal identifier
	RequestedAmount   int    // Shown in centavos
	Status            ApplicationStatus
}

func (a *Application) Validate() error {
	if a.ProfileId <= 0 {
		return ErrInvalidProduct
	}

	// Handles CreditLimit zero value passed by service 
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
	ProfileId int64
	CurrencyId int64
	CreditLimit int
	// InterestRate represents the rate in basis points
	// Example: 1 bps == 0.01% or 1250 bps == 12.50%
	InterestRate int
}

func (c *CreditCard) Validate() error {
	if c.ProfileId <= 0 {
		return ErrInvalidProduct
	}

	if c.CurrencyId < 0 {
		return ErrInvalidCurrency	
	}

	if c.CreditLimit < 0 {
		return ErrInvalidCreditLimit
	}

	if c.InterestRate < 0 {
		return ErrInvalidInterestRate
	}
}

type PersonalLoan struct {
	ProfileId int64
	CurrencyId int64
	LoanAmount int
	// InterestRate represents the rate in basis points
	// Example: 1 bps == 0.01% or 1250 bps == 12.50%
	InterestRate int
}

type Applicant struct {
	Birthday time.Time
	ContactNumbers []ContactNumer
	LastName string
	FirstName string
	MiddleName string
	IsPrincipal bool
}

type ContactNumber struct {
	Value string
	TypeId int64
}

