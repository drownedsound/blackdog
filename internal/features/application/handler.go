package app

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// Handler manages the HTTP transport for the Application feature.
type Handler struct {
	svc *Service
}

// NewHandler creates a handler with the required dependencies.
func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// RegisterRoutes registers the route patterns with the provided ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /application", h.HandleCreate)
	// mux.HandleFunc("GET /application/{id}", h.HandleGet)

	// TODO: Add HandleFunc for API
	//       1. POST /application (new application)
	//       2. GET /application/{id} (view existing application)
	//       3. PATCH /application/{id} (update existing application)

	// TODO: Add HandleFunc for web forms
	//	     1. /application/new
	//	     2. /application/100/save
	//	     3. /application/100/view
	//	     4. /application/100/edit
	//	     4. /application/100/submit
}

// HandleCreate processes the creation of a new application.
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateApplicationRequest

	// Limit payload to 1MB (adjust based on requirements)
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		slog.ErrorContext(
			r.Context(), "json decoding failed", slog.Any("error", err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)

		return
	}

	resp, err := h.svc.CreateApplication(r.Context(), req)
	if err != nil {
		// Handles infrastructure layer concerns
		if errors.Is(err, ErrConnectionRefused) {
			slog.ErrorContext(
				r.Context(), "database error", slog.Any("error", err),
			)

			http.Error(
				w, "Internal Server Error", http.StatusInternalServerError,
			)

			return
		}

		// Handles domain layer concerns
		slog.WarnContext(
			r.Context(), "domain validation failed", slog.Any("error", err),
		)

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(
			r.Context(), "json encoding failed", slog.Any("error", err),
		)

		http.Error(
			w, "Internal Server Error", http.StatusInternalServerError,
		)

		return
	}
}

// // HandleGet processes retrieving an application by Id.
// func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
// 	s := r.PathValue("id")
// 	id, err := strconv.ParseInt(s, 10, 64)
//
// 	if err != nil || id <= 0 {
// 		slog.WarnContext(
// 			r.Context(), "id validation failed", slog.Any("error", err),
// 		)
//
// 		http.Error(w, "Invalid Id", http.StatusBadRequest)
//
// 		return
// 	}
//
// 	req := GetApplicationRequest{Id: id}
// 	resp, err := h.svc.GetApplicationById(r.Context(), req)
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			slog.WarnContext(
// 				r.Context(), "unknown id", slog.Any("error", err),
// 			)
// 			http.Error(w, "Application Not Found", http.StatusNotFound)
//
// 			return
// 		}
//
// 		slog.ErrorContext(
// 			r.Context(), "database error", slog.Any("error", err),
// 		)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
//
// 		return
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(resp); err != nil {
// 		slog.ErrorContext(
// 			r.Context(), "json encoding failed", slog.Any("error", err),
// 		)
//
// 		http.Error(
// 			w, "Internal Server Error", http.StatusInternalServerError,
// 		)
//
// 		return
// 	}
// }
