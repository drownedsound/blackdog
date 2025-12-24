package app

import (
	"testing"
	"time"
	"errors"
)

func Test_CreditCard_Validate(t *testing.T) {
	base := CreditCard{
		ProfileId:    1,
		CurrencyId:   1, 
		CreditLimit:  1_000_000,
		InterestRate: 300, 
	}

	testCases := []struct {
		desc        string
		mutate func(*CreditCard)
		expectedErr error
	}{
		{
			desc: "Credit Card With Valid Details Passes Validation",
			mutate: func(c *CreditCard) {}, 
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid Profile Id Fails Validation",
			mutate: func (c *CreditCard)  {
				c.ProfileId = 0
			},
			expectedErr: ErrInvalidCreditCard,
		},
		{
			desc: "Credit Card With Invalid Currency Id Fails Validation",
			mutate: func(c * CreditCard) {
				c.CurrencyId = -1
			},
			expectedErr: ErrInvalidCurrency,
		},
		{
			desc: "Credit Card With Unset Credit Limit Passes Validation",
			mutate: func(c * CreditCard) {
				c.CreditLimit = 0
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid Credit Limit Fails Validation",
			mutate: func(c * CreditCard) {
				c.CreditLimit = -1_000_000
			},
			expectedErr: ErrInvalidCreditLimit,
		},
		{
			desc: "Credit Card With Unset Interest Rate Passes Validation",
			mutate: func(c * CreditCard) {
				c.InterestRate = 0
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid Interest Rate Fails Validation",
			mutate: func(c * CreditCard) {
				c.InterestRate = -300
			},
			expectedErr: ErrInvalidInterestRate,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			card := base
			tC.mutate(&card)
			err := card.Validate()

			if !errors.Is(err, tC.expectedErr) {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
					)
			}
		})
	}
}

func Test_PersonaLoan_Validate(t *testing.T) {
	base := PersonalLoan{
		ProfileId: 1,
		CurrencyId: 1,
		LoanAmount: 1_000_000,
		InterestRate: 300,
	}

	testCases := []struct {
		desc        string
		mutate func(*PersonalLoan)
		expectedErr error
	}{
		{
			desc: "Personal Loan With Valid Details Passes Validation",
			mutate: func (p *PersonalLoan) {}, 
			expectedErr: nil,
		},
		{
			desc: "Personal Loan With Invalid Profile Id Fails Validation",
			mutate: func(p *PersonalLoan) {
				p.ProfileId = 0
			},
			expectedErr: ErrInvalidPersonalLoan,
		},
		{
			desc: "Personal Loan With Invalid Currency Id Fails Validation",
			mutate: func (p *PersonalLoan) {
				p.CurrencyId = -1
			},
			expectedErr: ErrInvalidCurrency,
		},
		{
			desc: "Personal Loan With Unset Loan Amount Passes Validation",
			mutate: func(p *PersonalLoan) {
				p.LoanAmount = 0
			},
			expectedErr: nil,
		},
		{
			desc: "Personal Loan With Invalid Loan Amount Fails Validation",
			mutate: func(p *PersonalLoan) {
				p.LoanAmount = -1_000_000
			},
			expectedErr: ErrInvalidLoanAmount,
		},
		{
			desc: "Personal Loan With Unset Interest Rate Fails Validation",
			mutate: func (p *PersonalLoan) {
				p.InterestRate = 0
			},
			expectedErr: nil,
		},
		{
			desc: "Personal Loan With Invalid Interest Rate Fails Validation",
			mutate: func (p *PersonalLoan) {
				p.InterestRate = -300
			},
			expectedErr: ErrInvalidInterestRate,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			loan := base
			tC.mutate(&loan)
			err := loan.Validate()

			if !errors.Is(err, tC.expectedErr) {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
					)
			}
		})
	}
}

func Test_ContactNumber_Validate(t *testing.T) {
	base := ContactNumber {
		Value: "9171234567",
		Type: TypeMobile,
	}

	testCases := []struct {
		desc        string
		mutate func (*ContactNumber)
		expectedErr error
	}{
		{
			desc: "Contact Number With Valid Details Passes Validation",
			mutate: func(c *ContactNumber) {},
			expectedErr: nil,
		},
		{
			desc: "Contact Number With No Phone Number Fails Validation",
			mutate: func(c *ContactNumber) {
				c.Value = ""
			},
			expectedErr: ErrInvalidContactNumber,
		},
		{
			desc: "Contact Number With Invalid Phone Number Fails Validation",
			mutate: func (c *ContactNumber) {
				c.Value ="21234567" 
			},
			expectedErr: ErrInvalidContactNumber,
		},
		{
			desc: "Contact Number With Invalid Type Fails Validation",
			mutate: func(c *ContactNumber) {
				c.Type = TypeUnknown
			},
			expectedErr: ErrInvalidContactNumberType,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			contact := base
			tC.mutate(&contact)
			err := contact.Validate()

			if !errors.Is(err, tC.expectedErr) {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
					)
			}
		})
	}
}

func Test_Applicant_Validate(t *testing.T) {
	base := Applicant{
		LastName:    "Doe",
		FirstName:   "John",
		MiddleName:  "Smith",
		Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
		IsPrincipal: true,
		ContactNumbers: []ContactNumber{
			{Value: "09171234567", Type: TypeMobile},
		},
	}

	testCases := []struct {
		desc        string
		mutate func(*Applicant)
		expectedErr error
	}{
		{
			desc: "Applicant With Valid Details Passes Validation",
			mutate: func(a *Applicant) {},
			expectedErr: nil,
		},
		{
			desc: "Applicant With Missing Birthday Fails Validation",
			mutate: func(a *Applicant) {
				a.Birthday = time.Time {}
			},
			expectedErr: ErrMissingBirthday,
		},
		{
			desc: "Applicant With Birthday in Future Fails Validation",
			mutate: func(a *Applicant) {
				a.Birthday = time.Now().AddDate(0, 0, 1).UTC()
			},
			expectedErr: ErrBirthdayInFuture,
		},
		{
			desc: "Applicant Not Meeting Age Requirements Fails Validation",
			mutate: func(a *Applicant) {
				a.Birthday = time.Now().AddDate(-16, 0, 1).UTC()
			},
			expectedErr: ErrMinimumAgeNotMet,
		},
		{
			desc: "Applicant Missing Last Name Fails Validation",
			mutate: func(a *Applicant) {
				a.LastName = ""
			},
			expectedErr: ErrMissingLastName,
		},
		{
			desc: "Applicant With Lengthy Last Name Fails Validation",
			mutate: func(a *Applicant) {
				a.LastName = "Smith Cornell Vedder Staley Cobain Hetfield Hatfield"
			},
			expectedErr: ErrLastNameTooLong,
		},
		{
			desc: "Valid Name With Multi-Byte Characters Exceeding Byte Limit",
			mutate: func(a *Applicant) {
				// 16 characters * 2 bytes = 32 bytes.
				// 32 > 30 (Bytes) -> Enters Slow Path
				// 16 <= 30 (Runes) -> Validates OK
				a.LastName = "ÑÑÑÑÑÑÑÑÑÑÑÑÑÑÑÑ" 
			},
			expectedErr: nil,
		},
		{
			desc: "Applicant Missing First Name Fails Validation",
			mutate: func(a *Applicant) {
				a.FirstName = ""
			},
			expectedErr: ErrMissingFirstName,
		},
		{
			desc: "Applicant Lengthy First Name Fails Validation",
			mutate: func(a *Applicant) {
				a.FirstName = "John Paul George Simon Layne Eddie Kurt Chris Julian"
			},
			expectedErr: ErrFirstNameTooLong,
		},
		{
			desc: "Applicant Lengthy Middle Name Fails Validation",
			mutate: func(a *Applicant) {
				a.MiddleName = "Smith Cornell Vedder Staley Cobain Hetfield Hatfield"
			},
			expectedErr: ErrMiddleNameTooLong,
		},
		{
			desc: "Applicant With No Contact Number Fails Validation",
			mutate: func(a *Applicant) {
				a.ContactNumbers = []ContactNumber{}
			},
			expectedErr: ErrMissingContactNumber,
		},
		{
			desc: "Applicant With Invalid Contact Number Fails Validation",
			mutate: func(a *Applicant) {
				a.ContactNumbers = []ContactNumber{
					{ Value: "9171234567", Type: TypeMobile },
					{ Value: "12345", Type: TypeMobile },
				}
			},
			expectedErr: ErrInvalidContactNumber,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			app := base
			tC.mutate(&app)
			err := app.Validate()

			if !errors.Is(err, tC.expectedErr) {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
					)
			}
		})
	}
}

func Test_Application_Validate(t *testing.T) {
	base := Application{
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		OtherApplicants: []Applicant{},
		MemberReferenceNo: "ABCDE12345",
		Applicant: Applicant {
			Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
			ContactNumbers: []ContactNumber {
				{ Value: "9171234567", Type: TypeMobile },
			},
			LastName: "Smith",
			FirstName: "John",
			MiddleName: "Doe",
			IsPrincipal: true,
		},
		CreditCard: CreditCard {
			ProfileId: 1,
			CurrencyId: 1,
			CreditLimit: 1_000_000,
			InterestRate: 300,
		},
		Status:            StatusCreated,
		RequestedAmount:   100_000_000,
	}

	testCases := []struct {
		desc        string
		mutate func(*Application)
		expectedErr error
	}{
		{
			desc: "Application With Valid Details Validation",
			mutate: func(a *Application) {},
			expectedErr: nil,
		},
		{
			desc: "Application Created Yesterday Passes Validation",
			mutate: func(a *Application) {
				a.CreatedAt = time.Now().AddDate(0, 0, -1).UTC()
			},
			expectedErr: nil,
		},
		{
			desc: "Application Created Tomorrow Fails Validation",
			mutate: func(a *Application) {
				a.CreatedAt = time.Now().AddDate(0, 0, 1).UTC()
			},
			expectedErr: ErrCreatedAtInFuture,
		},
		{
			desc: "Application Updated Yesterday Passes Validation",
			mutate: func(a *Application) {
				a.UpdatedAt = time.Now().AddDate(0, 0, -1).UTC()
			},
			expectedErr: nil,
		},
		{
			desc: "Application Updated Tomorrow Fails Validation",
			mutate: func(a *Application) {
				a.UpdatedAt = time.Now().AddDate(0, 0, 1).UTC()
			},
			expectedErr: ErrUpdatedAtInFuture,
		},
		{
			desc: "Application With Missing MemberReferenceNumber" +
			"Fails Validation",
			mutate: func(a *Application) {
				a.MemberReferenceNo = ""
			},
			expectedErr: ErrMissingMemberRefNo,
		},
		{
			desc: "Application With Unset RequestedAmount Fails Validation",
			mutate: func(a *Application) {
				a.RequestedAmount = 0
			},
			expectedErr: ErrInvalidRequestedAmount,
		},
		{
			desc: "Application With Invalid RequestedAmount Fails Validation",
			mutate: func(a *Application) {
				a.RequestedAmount = -1_000_000
			},
			expectedErr: ErrInvalidRequestedAmount,
		},
		{
			desc: "Application With Invalid Status (Unknown) Fails Validation",
			mutate: func(a *Application) {
				a.Status = StatusUnknown
			},
			expectedErr: ErrInvalidStatus,
		},
		{
			desc: "Application With Invalid Status (Out of Range)" +
			"Fails Validation",
			mutate: func(a *Application) {
				a.Status = StatusCancelled + 1
			},
			expectedErr: ErrInvalidStatus,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			app := base
			tC.mutate(&app)
			err := app.Validate()

			if !errors.Is(err, tC.expectedErr) {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
					)
			}
		})
	}
}
