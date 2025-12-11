package app

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
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
	mux.HandleFunc("GET /application/{id}", h.HandleGet)
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// TODO: Log error
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.CreateApplication(r.Context(), req)
	if err != nil {
		// TODO: Add ERROR logging
		// TODO: Return HTTP 400 for validation failures
		// TODO: Return HTTP 500 for unexpected errors
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: Add ERROR logging
	}

	h.svc.logger.Info(
		"Flushing response body",
		slog.Any("body", resp),
	)
}

// HandleGet processes retrieving an application by Id.
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		// TODO: Add ERROR logging
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	req := GetApplicationRequest{Id: id}
	resp, err := h.svc.GetApplicationById(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// TODO: Add ERROR logging
			http.Error(w, "Application not found", http.StatusNotFound)
			return
		}

		// TODO: Add ERROR logging
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: Add ERROR logging
	}

	h.svc.logger.Info(
		"Flushing response body",
		slog.Any("body", resp),
	)
}
