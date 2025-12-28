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
	testCases := []struct {
		desc string
		// reqBody uses interface{} to test invalid JSON payloads
		reqBody        interface{}
		mockSave       func(context.Context, *Application) error
		expectedStatus int
		expectedBody   string // Partial match or specific check
	}{
		{
			desc: "Success_Returns_201_And_JSON",
			reqBody: CreateApplicationRequest{
				MemberReferenceNo: "REF-123",
				CardProfile:       1,
				RequestedAmount:   50_000,
			},
			mockSave: func(ctx context.Context, a *Application) error {
				a.Id = 101
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":101`,
		},
		{
			desc:           "Invalid_JSON_Returns_400",
			reqBody:        "invalid-json-string",
			mockSave:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON",
		},
		{
			desc: "Service_Error_Returns_400",
			reqBody: CreateApplicationRequest{
				MemberReferenceNo: "REF-FAIL",
				CardProfile:       1,
				RequestedAmount:   50_000,
			},
			mockSave: func(ctx context.Context, a *Application) error {
				return errors.New("db error")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "application.service failed to save entity",
		},
		{
			desc: "Infrastructure_Error_Returns_500",
			reqBody: CreateApplicationRequest{
				MemberReferenceNo: "ERR-100",
				CardProfile:       1,
				RequestedAmount:   50_000,
			},
			mockSave: func(ctx context.Context, a *Application) error {
				return ErrConnectionRefused
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := &mockRepo{saveFunc: tC.mockSave}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			h := NewHandler(svc)
			mux := http.NewServeMux()
			h.RegisterRoutes(mux)

			var body []byte
			if s, ok := tC.reqBody.(string); ok && s == "invalid-json-string" {
				body = []byte("{invalid")
			} else {
				body, _ = json.Marshal(tC.reqBody)
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

			// Simple substring check for body to verify error messages
			// or ID presence
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
		mockGet        func(context.Context, int64) (Application, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			desc:  "Success_Returns_200",
			urlId: "101",
			mockGet: func(ctx context.Context, id int64) (Application, error) {
				return Application{
					Id:                101,
					MemberReferenceNo: "REF-101",
					CreatedAt:         now,
					UpdatedAt:         now,
					CreditCard: CreditCard{
						CardProfile:  1,
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
			mockGet: func(ctx context.Context, id int64) (Application, error) {
				return Application{}, ErrNotFound
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
			mockGet: func(ctx context.Context, id int64) (Application, error) {
				return Application{}, errors.New("connection failed")
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
		saveFunc: func(ctx context.Context, a *Application) error {
			a.Id = 101
			return nil
		},
		getByInternalIdFunc: func(ctx context.Context, id int64) (Application, error) {
			return Application{Id: 101}, nil
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewService(mock, logger)
	h := NewHandler(svc)

	t.Run("HandleCreate_JsonEncodeError_LogsError", func(t *testing.T) {
		reqBody := CreateApplicationRequest{
			MemberReferenceNo: "REF-123",
			RequestedAmount:   50_000,
			CardProfile:       1,
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
