package app

import (
	"testing"
	"time"
)

func TestApplication_Validate_CardProfile(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Credit Card With Valid CardProfile Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid CardProfile Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  0,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrInvalidCardProfile,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func TestApplication_Validate_CreditLimit(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Credit Card With Valid CreditLimit Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Zero CreditLimit Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  0,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid CreditLimit Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  -1_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrInvalidCreditLimit,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func TestApplication_Validate_InterestRate(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Credit Card With Valid InterestRate Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Zero InterestRate Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 0,
				},
			},
			expectedErr: ErrInvalidInterestRate,
		},
		{
			desc: "Credit Card With Invalid InterestRate Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: -1_000,
				},
			},
			expectedErr: ErrInvalidInterestRate,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func TestApplication_Validate_CreatedAt(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Application Created Today Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application Created Yesterday Fails Validation",
			app: Application{
				CreatedAt:         time.Now().AddDate(0, 0, -1),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application Created Tomorrow Fails Validation",
			app: Application{
				CreatedAt:         time.Now().AddDate(0, 0, 1),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrCreatedAtInFuture,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func TestApplication_Validate_UpdatedAt(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Application Updated Today Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application Updated Yesterday Passes Validation",
			app: Application{
				CreatedAt:         time.Now().AddDate(0, 0, -1),
				UpdatedAt:         time.Now().AddDate(0, 0, -1),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application Upadated Tomorrow Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now().AddDate(0, 0, 1),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrUpdatedAtInFuture,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func TestApplication_Validate_MemberReferenceNumber(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Application With Valid MemberReferenceNumber" +
				"Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application With Invalid MemberReferenceNumber" +
				"Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrMissingMemberRefNo,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func TestApplication_Validate_RequestedAmount(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Application With Correct RequestedAmount Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application With Zero RequestedAmount Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   0,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrInvalidRequestedAmount,
		},
		{
			desc: "Application With Invalid RequestedAmount Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   -1,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrInvalidRequestedAmount,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func TestApplication_Validate_Status(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Card Application With Valid Status Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Loan Application With Valid Status Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCreated,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application With Invalid Status (Unknown)" +
				"Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusUnknown,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrInvalidStatus,
		},
		{
			desc: "Application With Invalid Status (Out of Range)" +
				"Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				Status:            StatusCancelled + 1,
				RequestedAmount:   100_000_000,
				CreditCard: CreditCard{
					CardProfile:  1,
					CreditLimit:  1_000_000,
					InterestRate: 1_000,
				},
			},
			expectedErr: ErrInvalidStatus,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.app.Validate()
			if err != tC.expectedErr {
				t.Errorf(
					"Validate Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}
