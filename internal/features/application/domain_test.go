package app

import (
	"testing"
	"time"
)

func TestApplication_Validate_CardProfileCode(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Credit Card With Valid CardProfileCode Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid CardProfileCode Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 0,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
				},
			},
			expectedErr: ErrInvalidCardProfileCode,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     -1_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Credit Card With Invalid InterestRate Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    -1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryLoan,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application With Invalid RequestedAmount Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				// CategoryCode:      CategoryLoan,
				StatusCode:      StatusCreated,
				RequestedAmount: -1,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
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

// func TestProductCategory_String(t *testing.T) {
// 	testCases := []struct {
// 		desc     string
// 		category ProductCategory
// 		expected string
// 	}{
// 		{
// 			desc:     "CategoryCard returns CARD",
// 			category: CategoryCard,
// 			expected: "CREDIT_CARD",
// 		},
// 		{
// 			desc:     "CategoryLoan returns LOAN",
// 			category: CategoryLoan,
// 			expected: "PERSONAL_LOAN",
// 		},
// 		{
// 			desc:     "CategoryUnknown returns UNKNOWN",
// 			category: CategoryUnknown,
// 			expected: "UNKNOWN",
// 		},
// 		{
// 			desc:     "Arbitrary Category returns UNKNOWN",
// 			category: ProductCategory(255),
// 			expected: "UNKNOWN",
// 		},
// 	}
//
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			if got := tC.category.String(); got != tC.expected {
// 				t.Errorf(
// 					"String() mismatch: Category %d, Expected %q, Got %q",
// 					tC.category,
// 					tC.expected,
// 					got,
// 				)
// 			}
// 		})
// 	}
// }

func TestApplication_Validate_StatusCode(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Card Application With Valid StatusCode Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Loan Application With Valid StatusCode Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				// CategoryCode:      CategoryLoan,
				StatusCode:      StatusCreated,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Application With Invalid StatusCode Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				// CategoryCode:      CategoryCard,
				StatusCode:      StatusUnknown,
				RequestedAmount: 100_000_000,
				CreditCard: CreditCard{
					CardProfileCode: 1,
					CreditLimit:     1_000_000,
					InterestRate:    1_000,
				},
			},
			expectedErr: ErrInvalidStatusCode,
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

func TestApplicationStatus_String(t *testing.T) {
	testCases := []struct {
		desc     string
		status   ApplicationStatus
		expected string
	}{
		{
			desc:     "StatusCreated returns CREATED",
			status:   StatusCreated,
			expected: "CREATED",
		},
		{
			desc:     "StatusInProgress returns IN_PROGRESS",
			status:   StatusInProgress,
			expected: "IN_PROGRESS",
		},
		{
			desc:     "StatusApproved returns APPROVED",
			status:   StatusApproved,
			expected: "APPROVED",
		},
		{
			desc:     "StatusDeclined returns DECLINED",
			status:   StatusDeclined,
			expected: "DECLINED",
		},
		{
			desc:     "StatusCancelled returns CANCELLED",
			status:   StatusCancelled,
			expected: "CANCELLED",
		},
		{
			desc:     "StatusUnknown (Zero Value) returns UNKNOWN",
			status:   StatusUnknown,
			expected: "UNKNOWN",
		},
		{
			desc:     "Arbitrary/Invalid Status returns UNKNOWN",
			status:   ApplicationStatus(255),
			expected: "UNKNOWN",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := tC.status.String(); got != tC.expected {
				t.Errorf(
					"String() mismatch: Status %d, Expected %q, Got %q",
					tC.status,
					tC.expected,
					got,
				)
			}
		})
	}
}
