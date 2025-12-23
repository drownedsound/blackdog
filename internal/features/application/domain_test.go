package app

import (
	"testing"
	"time"
	"errors"
)

func Test_CreditCard_Validate(t *testing.T) {
	testCases := []struct {
		desc        string
		card         CreditCard
		expectedErr error
	}{
		{
			desc: "Credit Card With Valid Details Passes Validation",
			card: CreditCard{
				ProfileId: 1,
				CurrencyId: 1,
				CreditLimit: 1_000_000,
				InterestRate: 300,
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid Profile Id Fails Validation",
			card: CreditCard{
				ProfileId: 0,
				CurrencyId: 1,
				CreditLimit: 1_000_000,
				InterestRate: 300,
			},
			expectedErr: ErrInvalidCreditCard,
		},
		{
			desc: "Credit Card With Invalid Credit Limit Fails Validation",
			card: CreditCard{
				ProfileId: 1,
				CurrencyId: 1,
				CreditLimit: -1_000_000,
				InterestRate: 300,
			},
			expectedErr: ErrInvalidCreditLimit,
		},
		{
			desc: "Credit Card With Invalid Interest Rate Fails Validation",
			card: CreditCard{
				ProfileId: 1,
				CurrencyId: 1,
				CreditLimit: 1_000_000,
				InterestRate: -300,
			},
			expectedErr: ErrInvalidInterestRate,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.card.Validate()
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
	testCases := []struct {
		desc        string
		loan         PersonalLoan
		expectedErr error
	}{
		{
			desc: "Personal Loan With Valid Details Passes Validation",
			loan: PersonalLoan{
				ProfileId: 1,
				CurrencyId: 1,
				LoanAmount: 1_000_000,
				InterestRate: 300,
			},
			expectedErr: nil,
		},
		{
			desc: "Personal Loan With Invalid Profile Id Fails Validation",
			loan: PersonalLoan{
				ProfileId: 0,
				CurrencyId: 1,
				LoanAmount: 1_000_000,
				InterestRate: 300,
			},
			expectedErr: ErrInvalidPersonalLoan,
		},
		{
			desc: "Personal Loan With Invalid Loan Amount Fails Validation",
			loan: PersonalLoan {
				ProfileId: 1,
				CurrencyId: 1,
				LoanAmount: -1_000_000,
				InterestRate: 300,
			},
			expectedErr: ErrInvalidLoanAmount,
		},
		{
			desc: "Personal Loan With Invalid Interest Rate Fails Validation",
			loan: PersonalLoan{
				ProfileId: 1,
				CurrencyId: 1,
				LoanAmount: 1_000_000,
				InterestRate: -300,
			},
			expectedErr: ErrInvalidInterestRate,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.loan.Validate()
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
	testCases := []struct {
		desc        string
		contact         ContactNumber
		expectedErr error
	}{
		{
			desc: "Contact Number With Valid Details Passes Validation",
			contact: ContactNumber{
				Value: "9171234567",
				Type: TypeMobile,
			},
			expectedErr: nil,
		},
		{
			desc: "Contact Number With No Phone Number Fails Validation",
			contact: ContactNumber{
				Value: "",
				Type: TypeHome,
			},
			expectedErr: ErrInvalidContactNumber,
		},
		{
			desc: "Contact Number With Invalid Phone Number Fails Validation",
			contact: ContactNumber{
				Value: "21234567",
				Type: TypeHome,
			},
			expectedErr: ErrInvalidContactNumber,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.contact.Validate()
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
	testCases := []struct {
		desc        string
		app         Applicant
		expectedErr error
	}{
		{
			desc: "Applicant With Valid Details Passes Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: nil,
		},
		{
			desc: "Applicant With Missing Birthday Fails Validation",
			app: Applicant {
				Birthday: time.Time{},
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrMissingBirthday,
		},
		{
			desc: "Applicant With Birthday in Future Fails Validation",
			app: Applicant {
				Birthday:  time.Now().AddDate(0, 0, 1).UTC(),
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrBirthdayInFuture,
		},
		{
			desc: "Applicant Not Meeting Age Requirements Fails Validation",
			app: Applicant {
				Birthday:  time.Now().AddDate(-16, 0, 0).UTC(),
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrMinimumAgeNotMet,
		},
		{
			desc: "Applicant Missing Last Name Fails Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrMissingLastName,
		},
		{
			desc: "Applicant With Lengthy Last Name Fails Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith Cornell Vedder Staley Cobain Hetfield Hatfield",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrLastNameTooLong,
		},
		{
			desc: "Applicant Missing First Name Fails Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrMissingFirstName,
		},
		{
			desc: "Applicant Lengthy First Name Fails Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "John Paul George Simon Layne Eddie Kurt Chris Julian",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrFirstNameTooLong,
		},
		{
			desc: "Applicant Lengthy Middle Name Fails Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {
					{ Value: "9171234567", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "John",
				MiddleName: "Smith Cornell Vedder Staley Cobain Hetfield Hatfield",
				IsPrincipal: true,
			},
			expectedErr: ErrMiddleNameTooLong,
		},
		{
			desc: "Applicant With No Contact Number Fails Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {},
				LastName: "Smith",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrMissingContactNumber,
		},
		{
			desc: "Applicant With Invalid Contact Number Fails Validation",
			app: Applicant {
				Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC), 
				ContactNumbers: []ContactNumber {
					{ Value: "12345", Type: TypeMobile },
				},
				LastName: "Smith",
				FirstName: "John",
				MiddleName: "Doe",
				IsPrincipal: true,
			},
			expectedErr: ErrInvalidContactNumber,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if !errors.Is(err, tC.expectedErr) {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
					)
			}
		})
	}
}

// func TestApplication_Validate_CardProfile(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Credit Card With Valid CardProfile Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Credit Card With Invalid CardProfile Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  0,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrInvalidCardProfile,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestApplication_Validate_CreditLimit(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Credit Card With Valid CreditLimit Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Credit Card With Zero CreditLimit Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  0,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Credit Card With Invalid CreditLimit Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  -1_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrInvalidCreditLimit,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestApplication_Validate_InterestRate(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Credit Card With Valid InterestRate Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Credit Card With Zero InterestRate Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 0,
// 				},
// 			},
// 			expectedErr: ErrInvalidInterestRate,
// 		},
// 		{
// 			desc: "Credit Card With Invalid InterestRate Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: -1_000,
// 				},
// 			},
// 			expectedErr: ErrInvalidInterestRate,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestApplication_Validate_CreatedAt(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Application Created Today Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Application Created Yesterday Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now().AddDate(0, 0, -1),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Application Created Tomorrow Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now().AddDate(0, 0, 1),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrCreatedAtInFuture,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestApplication_Validate_UpdatedAt(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Application Updated Today Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Application Updated Yesterday Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now().AddDate(0, 0, -1),
// 				UpdatedAt:         time.Now().AddDate(0, 0, -1),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Application Upadated Tomorrow Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now().AddDate(0, 0, 1),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrUpdatedAtInFuture,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestApplication_Validate_MemberReferenceNumber(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Application With Valid MemberReferenceNumber" +
// 				"Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Application With Invalid MemberReferenceNumber" +
// 				"Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrMissingMemberRefNo,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestApplication_Validate_RequestedAmount(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Application With Correct RequestedAmount Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Application With Zero RequestedAmount Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   0,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrInvalidRequestedAmount,
// 		},
// 		{
// 			desc: "Application With Invalid RequestedAmount Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   -1,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrInvalidRequestedAmount,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
//
// func TestApplication_Validate_Status(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		app         Application
// 		expectedErr error
// 	}{
// 		{
// 			desc: "Card Application With Valid Status Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Loan Application With Valid Status Passes Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCreated,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			desc: "Application With Invalid Status (Unknown)" +
// 				"Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusUnknown,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrInvalidStatus,
// 		},
// 		{
// 			desc: "Application With Invalid Status (Out of Range)" +
// 				"Fails Validation",
// 			app: Application{
// 				CreatedAt:         time.Now(),
// 				UpdatedAt:         time.Now(),
// 				MemberReferenceNo: "ABCDE12345",
// 				Status:            StatusCancelled + 1,
// 				RequestedAmount:   100_000_000,
// 				CreditCard: CreditCard{
// 					CardProfile:  1,
// 					CreditLimit:  1_000_000,
// 					InterestRate: 1_000,
// 				},
// 			},
// 			expectedErr: ErrInvalidStatus,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := tC.app.Validate()
// 			if err != tC.expectedErr {
// 				t.Errorf(
// 					"Validate Error: Expected = %v, Actual = %v",
// 					tC.expectedErr, err,
// 				)
// 			}
// 		})
// 	}
// }
