package app

import (
	"testing"
	"time"
)

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
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Application Created Yesterday Fails Validation",
			app: Application{
				CreatedAt:         time.Now().AddDate(0, 0, -1),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Application Created Tomorrow Fails Validation",
			app: Application{
				CreatedAt:         time.Now().AddDate(0, 0, 1),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
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
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Application Updated Yesterday Passes Validation",
			app: Application{
				CreatedAt:         time.Now().AddDate(0, 0, -1),
				UpdatedAt:         time.Now().AddDate(0, 0, -1),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Application Updated Tomorrow Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now().AddDate(0, 0, 1),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
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

func TestApplication_Validate_MemberReferenceNo(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Application With Valid MemberReferenceNo Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Application Missing MemberReferenceNo Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "",
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
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

func TestApplication_Validate_CategoryCode(t *testing.T) {
	testCases := []struct {
		desc        string
		app         Application
		expectedErr error
	}{
		{
			desc: "Card Application With Valid CategoryCode Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Loan Application With Valid CategoryCode Passes Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "LOAN",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Application Missing CategoryCode Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: ErrMissingCategoryCode,
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
				CategoryCode:      "CARD",
				StatusCode:        "CREATED",
				RequestedAmount:   100_000_000,
			},
			expectedErr: nil,
		},
		{
			desc: "Application With Invalid RequestedAmount Fails Validation",
			app: Application{
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				MemberReferenceNo: "ABCDE12345",
				CategoryCode:      "LOAN",
				StatusCode:        "CREATED",
				RequestedAmount:   -1,
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
