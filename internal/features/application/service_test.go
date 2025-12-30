package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"
)

func TestService_CreateApplication(t *testing.T) {
	// Base request with all required fields populated.
	baseReq := CreateApplicationRequest{
		MemberReferenceNumber: "APP-001",
		CardProfile:           1,
		RequestedAmount:       500_000,
		Birthday:              "1990-01-01",
		FirstName:             "John",
		LastName:              "Doe",
		ContactNumbers: []ContactNumberRequest{
			{Value: "09171234567", Type: 1}, // Type 1 = Mobile
		},
	}

	testCases := []struct {
		desc          string
		mutate        func(req *CreateApplicationRequest)
		mockInsert    func(ctx context.Context, a *Application) error
		expectedErr   error
		expectedId    int64
		expectedState ApplicationStatus
	}{
		{
			desc:   "Repository Save Succeeds",
			mutate: func(req *CreateApplicationRequest) {}, // No changes
			mockInsert: func(ctx context.Context, a *Application) error {
				a.Id = 101
				return nil
			},
			expectedErr:   nil,
			expectedId:    101,
			expectedState: StatusCreated,
		},
		{
			desc:   "Repository Save Fails",
			mutate: func(req *CreateApplicationRequest) {}, // No changes
			mockInsert: func(ctx context.Context, a *Application) error {
				return ErrInsertFailed
			},
			expectedErr: fmt.Errorf(
				"application.service failed to save entity: %w",
				ErrInsertFailed,
			),
		},
		{
			desc: "Domain Validation Failure (Invalid CardProfile)",
			mutate: func(req *CreateApplicationRequest) {
				req.CardProfile = -1
			},
			mockInsert: func(ctx context.Context, a *Application) error {
				t.Error("Repository Save should not be called on validation failure")
				return nil
			},
			expectedErr:   ErrInvalidCreditCard,
			expectedId:    0,
			expectedState: StatusUnknown,
		},
		{
			desc: "Validation Failure (Missing Birthday)",
			mutate: func(req *CreateApplicationRequest) {
				req.Birthday = ""
			},
			mockInsert: func(ctx context.Context, a *Application) error {
				return nil
			},
			expectedErr: errors.New("principal mapping failed: invalid birthday format"),
		},
		{
			desc: "Create Application With Other Applicants Succeeds",
			mutate: func(req *CreateApplicationRequest) {
				req.OtherApplicants = []ApplicantRequest{
					{
						FirstName: "Jane",
						LastName:  "Doe",
						Birthday:  "1992-02-02",
						ContactNumbers: []ContactNumberRequest{
							{Value: "09181234567", Type: 1},
						},
					},
				}
			},
			mockInsert: func(ctx context.Context, a *Application) error {
				a.Id = 102
				return nil
			},
			expectedErr:   nil,
			expectedId:    102,
			expectedState: StatusCreated,
		},
		{
			desc: "Create Application With Invalid Other Applicant Fails",
			mutate: func(req *CreateApplicationRequest) {
				req.OtherApplicants = []ApplicantRequest{
					{
						FirstName: "Jane",
						LastName:  "Doe",
						Birthday:  "invalid-date", // Triggers parsing error
						ContactNumbers: []ContactNumberRequest{
							{Value: "09181234567", Type: 1},
						},
					},
				}
			},
			mockInsert: func(ctx context.Context, a *Application) error {
				return nil
			},
			expectedErr: errors.New("other applicant 0 mapping failed: invalid birthday format"),
		},
		{
			desc: "Create Personal Loan Application Succeeds",
			mutate: func(req *CreateApplicationRequest) {
				req.CardProfile = 0
				req.LoanProfile = 1
			},
			mockInsert: func(ctx context.Context, a *Application) error {
				a.Id = 103
				return nil
			},
			expectedErr:   nil,
			expectedId:    103,
			expectedState: StatusCreated,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := &mockRepo{
				insertFunc: tC.mockInsert,
			}

			// Create a copy of the base request for mutation
			req := baseReq
			if tC.mutate != nil {
				tC.mutate(&req)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			resp, err := svc.CreateApplication(context.Background(), req)

			if tC.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tC.expectedErr)
					return
				}

				if err.Error() != tC.expectedErr.Error() && !errors.Is(err, tC.expectedErr) {
					// Handle wrapped error strings
					if len(tC.expectedErr.Error()) > 0 &&
						(len(err.Error()) < len(tC.expectedErr.Error()) ||
							err.Error()[:len(tC.expectedErr.Error())] != tC.expectedErr.Error()) {
						t.Errorf(
							"CreateApplication Error:\nExpected: %q\nActual:   %q",
							tC.expectedErr.Error(),
							err.Error(),
						)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected Error: %v", err)
			}

			if resp.Id != tC.expectedId {
				t.Errorf(
					"CreateApplication Error: Id Expected %d, Actual %d",
					tC.expectedId,
					resp.Id,
				)
			}
			if resp.Status != tC.expectedState {
				t.Errorf(
					"CreateApplication Error: Status Expected %d, Actual %d",
					tC.expectedState,
					resp.Status,
				)
			}
		})
	}
}

func TestService_GetApplicationById(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		desc                string
		req                 GetApplicationRequest
		mockGetByInternalId func(ctx context.Context, id int64) (*Application, error)
		expectedErr         error
		expectedState       ApplicationStatus
	}{
		{
			desc: "Get Existing Application Succeeds",
			req:  GetApplicationRequest{Id: 99},
			mockGetByInternalId: func(ctx context.Context, id int64) (*Application, error) {
				if id == 99 {
					return &Application{
						CreatedAt:             now,
						UpdatedAt:             now,
						MemberReferenceNumber: "APP-99",
						CreditCard: CreditCard{
							ProfileId: 1,
						},
						Status:          StatusApproved,
						Id:              99,
						RequestedAmount: 1_000_000,
					}, nil
				}
				return nil, ErrNotFound
			},
			expectedErr:   nil,
			expectedState: StatusApproved,
		},
		{
			desc: "Get Non-Existent Application Fails",
			req:  GetApplicationRequest{Id: 999},
			mockGetByInternalId: func(ctx context.Context, id int64) (*Application, error) {
				return nil, ErrNotFound
			},
			expectedErr: fmt.Errorf(
				"application.service failed to get application: %w",
				ErrNotFound,
			),
		},
		{
			desc: "Repository Failure Returns Error",
			req:  GetApplicationRequest{Id: 50},
			mockGetByInternalId: func(ctx context.Context, id int64) (*Application, error) {
				return nil, errors.New("connection refused")
			},
			expectedErr: fmt.Errorf(
				"application.service failed to get application: %w",
				errors.New("connection refused")),
		},
		{
			desc: "Get Personal Loan Application Succeeds",
			req:  GetApplicationRequest{Id: 100},
			mockGetByInternalId: func(ctx context.Context, id int64) (*Application, error) {
				return &Application{
					CreatedAt:             now,
					UpdatedAt:             now,
					MemberReferenceNumber: "APP-100",
					PersonalLoan: PersonalLoan{
						ProfileId:    1,
						InterestRate: 200,
					},
					Status:          StatusApproved,
					Id:              100,
					RequestedAmount: 100_000,
				}, nil
			},
			expectedErr:   nil,
			expectedState: StatusApproved,
		},
		{
			desc: "Get Application With Other Applicants Succeeds",
			req:  GetApplicationRequest{Id: 101},
			mockGetByInternalId: func(ctx context.Context, id int64) (*Application, error) {
				return &Application{
					CreatedAt:             now,
					UpdatedAt:             now,
					MemberReferenceNumber: "APP-101",
					CreditCard: CreditCard{
						ProfileId: 1,
					},
					OtherApplicants: []Applicant{
						{
							FirstName: "Jane",
							LastName:  "Doe",
							Birthday:  now.AddDate(-25, 0, 0),
							ContactNumbers: []ContactNumber{
								{Value: "09181234567", Type: TypeMobile},
							},
						},
					},
					Status:          StatusApproved,
					Id:              101,
					RequestedAmount: 100_000,
				}, nil
			},
			expectedErr:   nil,
			expectedState: StatusApproved,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := &mockRepo{
				getByInternalIdFunc: tC.mockGetByInternalId,
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			resp, err := svc.GetApplicationById(context.Background(), tC.req)

			if tC.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected Error %v, Got Nil", tC.expectedErr)
				}
				if err.Error() != tC.expectedErr.Error() {
					t.Errorf(
						"GetApplicationById Error: Expected: %q, Actual: %q",
						tC.expectedErr.Error(),
						err.Error(),
					)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected Error: %v", err)
			}

			if resp.Status != tC.expectedState {
				t.Errorf(
					"Status: Expected %d, Actual %d",
					tC.expectedState, resp.Status,
				)
				t.Errorf(
					"GetApplicationById Error: Status "+
						"Expected: %q, Actual: %q",
					tC.expectedState,
					resp.Status,
				)
			}
		})
	}
}

func TestService_Internal_Mappings(t *testing.T) {
	svc := NewService(nil, nil)
	now := time.Now()

	t.Run("mapApplicantDtoToEntity", func(t *testing.T) {
		t.Run("Valid Principal With Contacts Mapping", func(t *testing.T) {
			req := ApplicantRequest{
				Birthday:   "1990-01-01",
				FirstName:  "John",
				LastName:   "Doe",
				MiddleName: "Smith",
				ContactNumbers: []ContactNumberRequest{
					{Value: "09171234567", Type: 1},
				},
			}

			got, err := svc.mapApplicantDtoToEntity(req, true)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if got.FirstName != req.FirstName {
				t.Errorf("FirstName: want %q, got %q", req.FirstName, got.FirstName)
			}
			if !got.IsPrincipal {
				t.Error("Expected IsPrincipal to be true")
			}
			if len(got.ContactNumbers) != 1 {
				t.Errorf("Contacts: want 1, got %d", len(got.ContactNumbers))
			}
		})

		t.Run("Invalid Birthday Returns Error", func(t *testing.T) {
			req := ApplicantRequest{
				Birthday: "invalid-date",
			}

			_, err := svc.mapApplicantDtoToEntity(req, true)
			if err == nil {
				t.Error("Expected error for invalid birthday, got nil")
			}
		})

		t.Run("Empty Contacts Slice Mapping", func(t *testing.T) {
			req := ApplicantRequest{
				Birthday:       "1990-01-01",
				ContactNumbers: []ContactNumberRequest{},
			}

			got, err := svc.mapApplicantDtoToEntity(req, false)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(got.ContactNumbers) != 0 {
				t.Errorf("Contacts: want 0, got %d", len(got.ContactNumbers))
			}
		})
	})

	t.Run("mapApplicantEntityToDto", func(t *testing.T) {
		t.Run("Valid Entity Mapping", func(t *testing.T) {
			entity := Applicant{
				Birthday:    now,
				FirstName:   "Jane",
				LastName:    "Doe",
				IsPrincipal: false,
				ContactNumbers: []ContactNumber{
					{Value: "123456789", Type: TypeMobile},
				},
			}

			got := svc.mapApplicantEntityToDto(entity)

			if got.FirstName != entity.FirstName {
				t.Errorf("FirstName: want %q, got %q", entity.FirstName, got.FirstName)
			}
			if len(got.ContactNumbers) != 1 {
				t.Errorf("Contacts: want 1, got %d", len(got.ContactNumbers))
			} else if got.ContactNumbers[0].Value != entity.ContactNumbers[0].Value {
				// Asserting value here ensures the loop body actually executed correctly
				t.Errorf("Contact Value: want %q, got %q", entity.ContactNumbers[0].Value, got.ContactNumbers[0].Value)
			}

			if got.Birthday != now.Format(time.DateOnly) {
				t.Errorf("Birthday: want %q, got %q", now.Format(time.DateOnly), got.Birthday)
			}
		})
	})
}
