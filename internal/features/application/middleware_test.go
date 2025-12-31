package app

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLogger(t *testing.T) {
	testCases := []struct {
		desc          string
		handlerStatus int
		expectedLevel string // "INFO", "WARN", "ERROR"
		// Whether the handler explicitly calls WriteHeader
		shouldWriteHdr bool
	}{
		{
			desc:           "Success_200_Logged_As_Info",
			handlerStatus:  http.StatusOK,
			expectedLevel:  "INFO",
			shouldWriteHdr: true,
		},
		{
			desc:           "Implicit_200_Logged_As_Info",
			handlerStatus:  http.StatusOK,
			expectedLevel:  "INFO",
			shouldWriteHdr: false,
		},
		{
			desc:           "Client_Error_400_Logged_As_Warn",
			handlerStatus:  http.StatusBadRequest,
			expectedLevel:  "WARN",
			shouldWriteHdr: true,
		},
		{
			desc:           "Server_Error_500_Logged_As_Error",
			handlerStatus:  http.StatusInternalServerError,
			expectedLevel:  "ERROR",
			shouldWriteHdr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuf, nil))
			middleware := RequestLogger(logger)
			handler := http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if tC.shouldWriteHdr {
						w.WriteHeader(tC.handlerStatus)
					}
					w.Write([]byte("response body"))
				})

			wrappedHandler := middleware(handler)
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rec, req)

			expectedCode := tC.handlerStatus
			if !tC.shouldWriteHdr {
				expectedCode = http.StatusOK
			}

			if rec.Code != expectedCode {
				t.Errorf(
					"HTTP Status: Expected %d, Got %d", expectedCode, rec.Code,
				)
			}

			logOutput := logBuf.String()
			if !strings.Contains(logOutput, "level="+tC.expectedLevel) {
				t.Errorf(
					"Log Level: Expected %s to be in log output: %s",
					tC.expectedLevel,
					logOutput,
				)
			}

			expectedFields := []string{
				"method=", "path=", "status=", "dur=", "bytes=",
			}
			for _, field := range expectedFields {
				if !strings.Contains(logOutput, field) {
					t.Errorf(
						"Log Output missing field %q: %s", field, logOutput,
					)
				}
			}
		})
	}
}
