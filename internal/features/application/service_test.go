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

type mockRepo struct {
	saveFunc    func(ctx context.Context, a *Application) error
	getByIdFunc func(ctx context.Context, id int64) (Application, error)
}

func (m *mockRepo) Save(ctx context.Context, a *Application) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, a)
	}
	return nil
}

func (m *mockRepo) GetById(ctx context.Context, id int64) (Application, error) {
	if m.getByIdFunc != nil {
		return m.getByIdFunc(ctx, id)
	}
	return Application{}, ErrNotFound
}

func TestService_CreateApplication(t *testing.T) {
	testCases := []struct {
		desc          string
		req           CreateApplicationRequest
		mockSave      func(ctx context.Context, a *Application) error
		expectedErr   error
		expectedId    int64
		expectedState ApplicationStatus
	}{
		{
			desc: "Repository Save Succeeds",
			req: CreateApplicationRequest{
				MemberReferenceNo: "APP-001",
				// CategoryCode:      "PERSONAL_LOAN",
				RequestedAmount: 500_000,
			},
			mockSave: func(ctx context.Context, a *Application) error {
				a.Id = 101
				return nil
			},
			expectedErr:   nil,
			expectedId:    101,
			expectedState: StatusCreated,
		},
		{
			desc: "Repository Save Fails",
			req: CreateApplicationRequest{
				MemberReferenceNo: "APP-002",
				// CategoryCode:      "CREDIT_CARD",
				RequestedAmount: 100_000,
			},
			mockSave: func(ctx context.Context, a *Application) error {
				return ErrInsertFailed
			},
			expectedErr: fmt.Errorf(
				"application.service failed to save entity: %w",
				ErrInsertFailed,
			),
		},
		// {
		// 	desc: "Domain Validation Failure (Invalid Category)",
		// 	req: CreateApplicationRequest{
		// 		MemberReferenceNo: "APP-INVALID",
		// 		// CategoryCode:      "INVALID",
		// 		RequestedAmount: 100_000,
		// 	},
		// 	mockSave: func(ctx context.Context, a *Application) error {
		// 		t.Error(
		// 			"Repository Save should not be called" +
		// 				" on validation failure")
		// 		return nil
		// 	},
		// 	expectedErr: fmt.Errorf(
		// 		"application.service stopped saving entity: %w",
		// 		ErrInvalidCategoryCode,
		// 	),
		// 	expectedId:    0,
		// 	expectedState: "",
		// },
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := &mockRepo{
				saveFunc: tC.mockSave,
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			resp, err := svc.CreateApplication(context.Background(), tC.req)

			if tC.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tC.expectedErr)
				}
				// Compare error strings because wrapped errors
				// have different pointers
				if err.Error() != tC.expectedErr.Error() {
					t.Errorf(
						"CreateApplication Error: Expected: %q, Actual: %q",
						tC.expectedErr.Error(),
						err.Error(),
					)
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
			if resp.StatusCode != tC.expectedState {
				t.Errorf(
					"CreateApplication Error: StatusCode Expected %s,"+
						"Actual %s",
					tC.expectedState,
					resp.StatusCode,
				)
			}
		})
	}
}

func TestService_GetApplicationById(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		desc          string
		req           GetApplicationRequest
		mockGetById   func(ctx context.Context, id int64) (Application, error)
		expectedErr   error
		expectedState ApplicationStatus
	}{
		{
			desc: "Get Existing Application Succeeds",
			req:  GetApplicationRequest{Id: 99},
			mockGetById: func(ctx context.Context, id int64) (
				Application, error,
			) {
				if id == 99 {
					return Application{
						CreatedAt:         now,
						UpdatedAt:         now,
						MemberReferenceNo: "APP-99",
						// CategoryCode:      CategoryCard,
						StatusCode:      StatusApproved,
						Id:              99,
						RequestedAmount: 1_000_000,
					}, nil
				}
				return Application{}, ErrNotFound
			},
			expectedErr:   nil,
			expectedState: StatusApproved,
		},
		{
			desc: "Get Non-Existent Application Fails",
			req:  GetApplicationRequest{Id: 999},
			mockGetById: func(ctx context.Context, id int64) (
				Application, error,
			) {
				return Application{}, ErrNotFound
			},
			expectedErr: fmt.Errorf(
				"application.service failed to get application: %w",
				ErrNotFound,
			),
		},
		{
			desc: "Repository Failure Returns Error",
			req:  GetApplicationRequest{Id: 50},
			mockGetById: func(ctx context.Context, id int64) (
				Application, error,
			) {
				return Application{}, errors.New("connection refused")
			},
			expectedErr: fmt.Errorf(
				"application.service failed to get application: %w",
				errors.New("connection refused")),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := &mockRepo{
				getByIdFunc: tC.mockGetById,
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

			if resp.StatusCode != tC.expectedState {
				t.Errorf(
					"Status: Expected %s, Actual %s",
					tC.expectedState, resp.StatusCode,
				)
				t.Errorf(
					"GetApplicationById Error: StatusCode "+
						"Expected: %q, Actual: %q",
					tC.expectedState,
					resp.StatusCode,
				)
			}
		})
	}
}
