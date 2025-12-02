package common

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

type Services struct {
	Logger *slog.Logger
}

func (s *Services) ServerError(
	w http.ResponseWriter, r *http.Request, err error,
) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	s.Logger.Error(
		err.Error(),
		slog.String("method", method),
		slog.String("uri", uri),
		slog.String("trace", trace),
	)

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
