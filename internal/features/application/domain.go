package app

import (
	"fmt"
	// "errors"
	// "strings"
	"time"
	"unicode/utf8"
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

type ContactNumberType byte

const (
	TypeUnknown ContactNumberType = iota
	TypeMobile
	TypeHome
	TypeOffice
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
	// if a.ProfileId <= 0 {
	// 	return ErrInvalidProduct
	// }
	//
	// // Handles CreditLimit zero value passed by service 
	// if a.CreditLimit != 0 && a.CreditLimit < 1 {
	// 	return ErrInvalidCreditLimit
	// }
	//
	// if a.InterestRate <= 0 {
	// 	return ErrInvalidInterestRate
	// }
	//
	// now := time.Now()
	//
	// if a.CreatedAt.After(now) {
	// 	return ErrCreatedAtInFuture
	// }
	//
	// if a.UpdatedAt.After(now) {
	// 	return ErrUpdatedAtInFuture
	// }
	//
	// if strings.TrimSpace(a.MemberReferenceNo) == "" {
	// 	return ErrMissingMemberRefNo
	// }
	//
	// if a.RequestedAmount <= 0 {
	// 	return ErrInvalidRequestedAmount
	// }
	//
	// if a.Status == StatusUnknown || a.Status > StatusCancelled {
	// 	return ErrInvalidStatus
	// }

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
	ProfileId int64
	CurrencyId int64
	LoanAmount int
	// InterestRate represents the rate in basis points
	// Example: 1 bps == 0.01% or 1250 bps == 12.50%
	InterestRate int
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
	Birthday time.Time
	ContactNumbers []ContactNumber
	LastName string
	FirstName string
	MiddleName string
	IsPrincipal bool
}

func (a *Applicant) Validate() error {
	// Checks for zero time.Time (0001-01-01 00:00:00 UTC)
	if a.Birthday.IsZero() {
		return ErrMissingBirthday
	}	

	now := time.Now().UTC()
	if a.Birthday.After(now) {
		return ErrBirthdayInFuture	
	}

	if a.Birthday.AddDate(18, 0, 0).After(now) {
		return ErrMinimumAgeNotMet	
	}

	// TODO: Ensure that service trims the string
	if a.LastName == "" {
		return ErrMissingLastName
	}
	
	// Does not use len() in case non-ASCII characters are used
	if utf8.RuneCountInString(a.LastName) > 30 {
		return ErrLastNameTooLong
	}
	
	// TODO: Ensure that service trims the string
	if a.FirstName == "" {
		return ErrMissingFirstName
	}
	
	// Does not use len() in case non-ASCII characters are used
	if utf8.RuneCountInString(a.FirstName) > 30 {
		return ErrFirstNameTooLong
	}

	// TODO: Ensure that service trims the string
	// Does not use len() in case non-ASCII characters are used
	if utf8.RuneCountInString(a.MiddleName) > 30 {
		return ErrMiddleNameTooLong
	}

	if len(a.ContactNumbers) == 0 {
		return ErrMissingContactNumber
	}

	for i := range a.ContactNumbers {
		if err := a.ContactNumbers[i].Validate(); err != nil {
			return fmt.Errorf("contact number %d failed validation: %w", i, err)
		}
	}

	return nil
}

type ContactNumber struct {
	Value string
	Type ContactNumberType
}

func (c *ContactNumber) Validate() error {
	// TODO: Ensure that service trims the string
	if c.Value == "" || len(c.Value) < 9 {
		return ErrInvalidContactNumber
	}

	if c.Type == TypeUnknown || c.Type > TypeOffice {
		return ErrInvalidContactNumberType
	}

	return nil
}
