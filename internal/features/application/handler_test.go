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
	// Performance: Reuse the encoder/decoder buffer to minimize allocations in high-throughput tests
	// though for unit tests, clarity > micro-optimization.

	testCases := []struct {
		desc           string
		reqBody        interface{} // Use interface{} to test invalid JSON payloads
		mockSave       func(context.Context, *Application) error
		expectedStatus int
		expectedBody   string // Partial match or specific check
	}{
		{
			desc: "Success_Returns_201_And_JSON",
			reqBody: CreateApplicationRequest{
				MemberReferenceNo: "REF-123",
				CategoryCode:      "loan",
				RequestedAmount:   50000,
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
			desc: "Service_Error_Returns_400", // Based on your TODO in handler.go
			reqBody: CreateApplicationRequest{
				MemberReferenceNo: "REF-FAIL",
				CategoryCode:      "loan",
				RequestedAmount:   50000,
			},
			mockSave: func(ctx context.Context, a *Application) error {
				return errors.New("db error")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to save application",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// 1. Setup Dependencies
			mock := &mockRepo{saveFunc: tC.mockSave}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			h := NewHandler(svc)
			mux := http.NewServeMux()
			h.RegisterRoutes(mux)

			// 2. Prepare Request
			var body []byte
			if s, ok := tC.reqBody.(string); ok && s == "invalid-json-string" {
				body = []byte("{invalid")
			} else {
				body, _ = json.Marshal(tC.reqBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/application", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			// 3. Execute (ServeHTTP to test routing + handler)
			mux.ServeHTTP(rec, req)

			// 4. Verify
			if rec.Code != tC.expectedStatus {
				t.Errorf("Status: Expected %d, Got %d", tC.expectedStatus, rec.Code)
			}

			// Simple substring check for body to verify error messages or ID presence
			if tC.expectedBody != "" {
				if !bytes.Contains(rec.Body.Bytes(), []byte(tC.expectedBody)) {
					t.Errorf("Body: Expected to contain %q, Got %q", tC.expectedBody, rec.Body.String())
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
					Id:                101,
					MemberReferenceNo: "REF-101",
					CreatedAt:         now,
					UpdatedAt:         now,
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
			expectedBody:   "Application not found",
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
			mock := &mockRepo{getByIdFunc: tC.mockGet}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			svc := NewService(mock, logger)
			h := NewHandler(svc)
			mux := http.NewServeMux()
			h.RegisterRoutes(mux)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/application/%s", tC.urlId), nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tC.expectedStatus {
				t.Errorf("Status: Expected %d, Got %d", tC.expectedStatus, rec.Code)
			}

			if tC.expectedBody != "" {
				if !bytes.Contains(rec.Body.Bytes(), []byte(tC.expectedBody)) {
					t.Errorf("Body: Expected to contain %q, Got %q", tC.expectedBody, rec.Body.String())
				}
			}
		})
	}
}
