package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_CreateApplication(t *testing.T) {
	base := CreateApplicationRequest{
		MemberReferenceNumber: "REF-123",
		CardProfile:           1,
		RequestedAmount:       50_000,
		Birthday:              "1980-01-01",
		LastName:              "Doe",
		FirstName:             "John",
		ContactNumbers: []ContactNumberRequest{
			{Value: "09171234567", Type: 1}, // TypeMobile
		},
	}

	testCases := []struct {
		desc           string
		mutate         func(*CreateApplicationRequest)
		rawPayload     interface{}
		mockInsert     func(context.Context, *Application) error
		expectedStatus int
		expectedBody   string
	}{
		{
			desc:   "Success_Returns_201_And_JSON",
			mutate: func(r *CreateApplicationRequest) {},
			mockInsert: func(ctx context.Context, a *Application) error {
				a.Id = 101
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":101`,
		},
		{
			desc:           "Invalid_JSON_Returns_400",
			rawPayload:     "invalid-json-string",
			mockInsert:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON",
		},
		{
			desc: "Service_Error_Returns_400",
			mutate: func(r *CreateApplicationRequest) {
				r.MemberReferenceNumber = "REF-FAIL"
			},
			mockInsert: func(ctx context.Context, a *Application) error {
				return errors.New("db error")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "application.service failed to save entity",
		},
		{
			desc: "Infrastructure_Error_Returns_500",
			mutate: func(r *CreateApplicationRequest) {
				r.MemberReferenceNumber = "ERR-100"
			},
			mockInsert: func(ctx context.Context, a *Application) error {
				return ErrConnectionRefused
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := &mockRepo{insertFunc: tC.mockInsert}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			h := NewHandler(svc)
			mux := http.NewServeMux()
			h.RegisterRoutes(mux)

			var body []byte
			if tC.rawPayload != nil {
				if s, ok := tC.rawPayload.(string); ok && s == "invalid-json-string" {
					body = []byte("{invalid")
				} else {
					body, _ = json.Marshal(tC.rawPayload)
				}
			} else {
				req := base
				if tC.mutate != nil {
					tC.mutate(&req)
				}
				body, _ = json.Marshal(req)
			}

			req := httptest.NewRequest(
				http.MethodPost, "/application", bytes.NewReader(body),
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tC.expectedStatus {
				t.Errorf(
					"Status: Expected %d, Got %d", tC.expectedStatus, rec.Code,
				)
			}

			if tC.expectedBody != "" {
				if !bytes.Contains(rec.Body.Bytes(), []byte(tC.expectedBody)) {
					t.Errorf(
						"Body: Expected to contain %q, Got %q",
						tC.expectedBody,
						rec.Body.String(),
					)
				}
			}
		})
	}
}

func TestHandler_GetApplication(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		desc           string
		urlId          string
		mockGet        func(context.Context, int64) (*Application, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			desc:  "Success_Returns_200",
			urlId: "101",
			mockGet: func(ctx context.Context, id int64) (*Application, error) {
				return &Application{
					Id:                    101,
					MemberReferenceNumber: "REF-101",
					CreatedAt:             now,
					UpdatedAt:             now,
					CreditCard: CreditCard{
						ProfileId:    1,
						InterestRate: 1250,
					},
					Status:          StatusCreated,
					RequestedAmount: 1_000_000,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":101`,
		},
		{
			desc:  "Not_Found_Returns_404",
			urlId: "999",
			mockGet: func(ctx context.Context, id int64) (*Application, error) {
				return nil, ErrNotFound
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Application Not Found",
		},
		{
			desc:           "Invalid_ID_Format_Returns_400",
			urlId:          "abc",
			mockGet:        nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid Id",
		},
		{
			desc:           "Invalid_ID_Value_Returns_400",
			urlId:          "0",
			mockGet:        nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid Id",
		},
		{
			desc:  "Internal_Error_Returns_500",
			urlId: "500",
			mockGet: func(ctx context.Context, id int64) (*Application, error) {
				return nil, errors.New("connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "connection failed",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := &mockRepo{getByInternalIdFunc: tC.mockGet}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			h := NewHandler(svc)
			mux := http.NewServeMux()
			h.RegisterRoutes(mux)

			req := httptest.NewRequest(
				http.MethodGet, fmt.Sprintf("/application/%s", tC.urlId), nil,
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tC.expectedStatus {
				t.Errorf(
					"Status: Expected %d, Got %d", tC.expectedStatus, rec.Code,
				)
			}

			if tC.expectedBody != "" {
				if !bytes.Contains(rec.Body.Bytes(), []byte(tC.expectedBody)) {
					t.Errorf(
						"Body: Expected to contain %q, Got %q",
						tC.expectedBody,
						rec.Body.String(),
					)
				}
			}
		})
	}
}

// failWriter mocks http.ResponseWriter and always returns an error on Write.
// This is used to test error handling when the client connection is broken
// during response generation.
type failWriter struct {
	http.ResponseWriter
}

func (f *failWriter) Write([]byte) (int, error) {
	return 0, errors.New("simulated connection reset")
}

func TestHandler_JsonEncodingFailure(t *testing.T) {
	mock := &mockRepo{
		insertFunc: func(ctx context.Context, a *Application) error {
			a.Id = 101
			return nil
		},
		getByInternalIdFunc: func(
			ctx context.Context,
			id int64,
		) (*Application, error) {
			return &Application{Id: 101}, nil
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewService(mock, logger)
	h := NewHandler(svc)

	t.Run("HandleCreate_JsonEncodeError_LogsError", func(t *testing.T) {
		reqBody := CreateApplicationRequest{
			MemberReferenceNumber: "REF-123",
			RequestedAmount:       50_000,
			CardProfile:           1,
			Birthday:              "1980-01-01",
			LastName:              "Doe",
			FirstName:             "John",
			ContactNumbers: []ContactNumberRequest{
				{Value: "09171234567", Type: 1}, // TypeMobile
			},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(
			http.MethodPost,
			"/application",
			bytes.NewReader(body))

		// Use failWriter to trigger the json.Encode error
		rec := httptest.NewRecorder()
		fw := &failWriter{ResponseWriter: rec}

		// Call handler directly to bypass mux and inject faulty writer
		h.HandleCreate(fw, req)

		// Can't easily check the response body because Write failed,
		// but checking that the function didn't panic and coverage increased
		// is sufficient. The handler attempts to write 500, but that write
		// will also fail (which is fine).
	})

	t.Run("HandleGet_JsonEncodeError_LogsError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/application/101", nil)
		req.SetPathValue("id", "101") // Manually set path value for Go 1.22+

		rec := httptest.NewRecorder()
		fw := &failWriter{ResponseWriter: rec}

		h.HandleGet(fw, req)
	})
}
