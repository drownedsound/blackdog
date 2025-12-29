package app

import (
	"fmt"
	"time"
	"unicode/utf8"
)

type validator struct {
	now                    time.Time
	eighteenYearsAgo       time.Time
	maxNameLength          int
	minContactNumberLength int
}

func newValidator(now time.Time) *validator {
	// TODO: Move magic numbers to a configuration file
	return &validator{
		now:                    now,
		eighteenYearsAgo:       now.AddDate(-18, 0, 0),
		maxNameLength:          30,
		minContactNumberLength: 9,
	}
}

type ApplicationStatus byte

const (
	StatusUnknown ApplicationStatus = iota
	StatusCreated
	StatusInProgress
	StatusApproved
	StatusDeclined
	StatusCancelled
)

type ContactNumberType byte

const (
	TypeUnknown ContactNumberType = iota
	TypeMobile
	TypeHome
	TypeOffice
)

type Application struct {
	CreatedAt       time.Time
	UpdatedAt       time.Time
	OtherApplicants []Applicant
	// Used as the external identifier
	MemberReferenceNumber string
	Applicant
	CreditCard   CreditCard
	PersonalLoan PersonalLoan
	// Used as the internal identifier
	Id int64
	// Shown in centavos
	RequestedAmount int32
	Status          ApplicationStatus
}

func (a *Application) Validate(v *validator) error {
	if a.CreatedAt.After(v.now) {
		return ErrCreatedAtInFuture
	}

	if a.UpdatedAt.After(v.now) {
		return ErrUpdatedAtInFuture
	}

	// TODO: Ensure that service trims the string
	if a.MemberReferenceNumber == "" {
		return ErrMissingMemberReferenceNumber
	}

	if a.RequestedAmount <= 0 {
		return ErrInvalidRequestedAmount
	}

	if a.Status == StatusUnknown || a.Status > StatusCancelled {
		return ErrInvalidStatus
	}

	a.Applicant.Validate(v)
	if err := a.Applicant.Validate(v); err != nil {
		return fmt.Errorf("principal applicant failed validation: %w", err)
	}

	if len(a.OtherApplicants) > 0 {
		for i := range a.OtherApplicants {
			if err := a.OtherApplicants[i].Validate(v); err != nil {
				return fmt.Errorf("applicant %d failed validation: %w", i, err)
			}
		}
	}

	// Application should have either a credit card or personal loan
	if a.CreditCard.ProfileId == 0 && a.PersonalLoan.ProfileId == 0 {
		return ErrMissingProduct
	}

	// Application cannot be for both credit card and personal loan
	if a.CreditCard.ProfileId > 0 && a.PersonalLoan.ProfileId > 0 {
		return ErrTooManyProducts
	}

	if a.CreditCard.ProfileId > 0 {
		if err := a.CreditCard.Validate(); err != nil {
			return fmt.Errorf("credit card failed validation: %w", err)
		}
	}

	if a.PersonalLoan.ProfileId > 0 {
		if err := a.PersonalLoan.Validate(); err != nil {
			return fmt.Errorf("personal loan failed validation: %w", err)
		}
	}

	return nil
}

type CreditCard struct {
	ProfileId  int64
	CurrencyId int64
	// Caps at PHP 20,000,000 when using int32 (4 bytes)
	CreditLimit int32
	// InterestRate represents the rate in basis points
	// Example: 1 bps == 0.01% or 1250 bps == 12.50%
	// 1 bps is 0.01%. Even if InterestRate is 100% that is only 10,000 bps.
	// This fits easily into int16 (2 bytes, max 32,767)
	InterestRate int16
}

func (c *CreditCard) Validate() error {
	if c.ProfileId <= 0 {
		return ErrInvalidCreditCard
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

	return nil
}

type PersonalLoan struct {
	ProfileId  int64
	CurrencyId int64
	// Caps at PHP 20,000,000 when using int32 (4 bytes)
	LoanAmount int32
	// InterestRate represents the rate in basis points
	// Example: 1 bps == 0.01% or 1250 bps == 12.50%
	// 1 bps is 0.01%. Even if InterestRate is 100% that is only 10,000 bps.
	// This fits easily into int16 (2 bytes, max 32,767)
	InterestRate int16
}

func (p *PersonalLoan) Validate() error {
	if p.ProfileId <= 0 {
		return ErrInvalidPersonalLoan
	}

	if p.CurrencyId < 0 {
		return ErrInvalidCurrency
	}

	if p.LoanAmount < 0 {
		return ErrInvalidLoanAmount
	}

	if p.InterestRate < 0 {
		return ErrInvalidInterestRate
	}

	return nil
}

type Applicant struct {
	Birthday       time.Time
	ContactNumbers []ContactNumber
	LastName       string
	FirstName      string
	MiddleName     string
	IsPrincipal    bool
}

func (a *Applicant) Validate(v *validator) error {
	// Checks for zero time.Time (0001-01-01 00:00:00 UTC)
	if a.Birthday.IsZero() {
		return ErrMissingBirthday
	}

	// Use v.now to ensure that time.Now() is only done once
	if a.Birthday.After(v.now) {
		return ErrBirthdayInFuture
	}

	// Use v.eighteenYearsAgo to prevent recalculation of date
	if a.Birthday.After(v.eighteenYearsAgo) {
		return ErrMinimumAgeNotMet
	}

	// TODO: Ensure that service trims the string
	if a.LastName == "" {
		return ErrMissingLastName
	}

	// Fast Path - len(a.LastName) reads the length from the
	// slice header (stack). This is an O(1) operation with minimal cost.
	// If the byte count is within the limit, the rune count is guaranteed
	// to be safe.
	if len(a.LastName) > v.maxNameLength {
		// Slow Path - Handles non-ASCII. Incur the O(N) CPU cost of decoding
		// UTF-8 if the byte count exceeds the limit.
		if utf8.RuneCountInString(a.LastName) > v.maxNameLength {
			return ErrLastNameTooLong
		}
	}

	// TODO: Ensure that service trims the string
	if a.FirstName == "" {
		return ErrMissingFirstName
	}

	// Same approach as LastName above
	if len(a.FirstName) > v.maxNameLength {
		if utf8.RuneCountInString(a.FirstName) > v.maxNameLength {
			return ErrFirstNameTooLong
		}
	}

	// TODO: Ensure that service trims the string
	// Same approach as LastName above
	if len(a.MiddleName) > v.maxNameLength {
		if utf8.RuneCountInString(a.MiddleName) > v.maxNameLength {
			return ErrMiddleNameTooLong
		}
	}

	// Applicant should have at least one contact number
	if len(a.ContactNumbers) == 0 {
		return ErrMissingContactNumber
	}

	for i := range a.ContactNumbers {
		if err := a.ContactNumbers[i].Validate(v); err != nil {
			return fmt.Errorf("contact number %d failed validation: %w", i, err)
		}
	}

	return nil
}

type ContactNumber struct {
	Value string
	Type  ContactNumberType
}

func (c *ContactNumber) Validate(v *validator) error {
	// TODO: Ensure that service trims the string
	if c.Value == "" || len(c.Value) < v.minContactNumberLength {
		return ErrInvalidContactNumber
	}

	if c.Type == TypeUnknown || c.Type > TypeOffice {
		return ErrInvalidContactNumberType
	}

	return nil
}
