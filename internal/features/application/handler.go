package app

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
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
	mux.HandleFunc("GET /{$}", h.Home)
	mux.HandleFunc("GET /site/{$}", h.Home)
	mux.HandleFunc("GET /site/application/new", h.NewApplicationForm)
	mux.HandleFunc("POST /site/application/save/{id}", h.SaveApplicationForm)
	mux.HandleFunc("POST /site/application/submit/{id}", h.SubmitApplicationForm)
	mux.HandleFunc("GET /site/application/view/{id}", h.ViewApplicationForm)

	mux.HandleFunc("POST /api/application", h.Create)
	mux.HandleFunc("GET /api/application/{id}", h.Get)

	// TODO: Add HandleFunc for API
	//
	// 1. ✅ POST /api/application (new application)
	// 2. ✅ GET /api/application/{id} (view existing application)
	// 3. ❌ GET /api/application (view existing applications
	//				with paging, sorting and filtering
	// 4. ❌ PATCH /api/application/{id} (update existing application)

	// TODO: Add HandleFunc for web forms
	// 1. ✅ /site/application/new
	// 2. ✅ /site/application/save/100
	// 3. ✅ /site/application/view/100
	// 4. ❌ /site/application/100/edit
	// 5. ✅ /site/application/100/submit
}

// HandleCreate processes the creation of a new application.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateApplicationRequest

	// Limit payload to 1MB (adjust based on requirements)
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		h.svc.logger.ErrorContext(
			r.Context(), "json decoding failed", slog.Any("error", err),
		)

		http.Error(w, "invalid json", http.StatusBadRequest)

		return
	}

	resp, err := h.svc.CreateApplication(r.Context(), req)
	if err != nil {
		// Handles infrastructure failures
		if errors.Is(err, ErrConnectionRefused) {
			h.svc.logger.ErrorContext(
				r.Context(), "database error", slog.Any("error", err),
			)

			http.Error(
				w, "internal server error", http.StatusInternalServerError,
			)

			return
		}

		// Handles domain failures
		h.svc.logger.WarnContext(
			r.Context(), "domain validation failed", slog.Any("error", err),
		)

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.svc.logger.ErrorContext(
			r.Context(), "json encoding failed", slog.Any("error", err),
		)

		http.Error(
			w, "internal server error", http.StatusInternalServerError,
		)

		return
	}
}

// HandleGet processes retrieving an application by Id.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	s := r.PathValue("id")
	id, err := strconv.ParseInt(s, 10, 64)

	if err != nil || id <= 0 {
		h.svc.logger.WarnContext(
			r.Context(), "id validation failed", slog.Any("error", err),
		)

		http.Error(w, "invalid id", http.StatusBadRequest)

		return
	}

	req := GetApplicationRequest{Id: id}
	resp, err := h.svc.GetApplicationById(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			h.svc.logger.WarnContext(
				r.Context(), "unknown id", slog.Any("error", err),
			)
			http.Error(w, "application not found", http.StatusNotFound)

			return
		}

		h.svc.logger.ErrorContext(
			r.Context(), "database error", slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.svc.logger.ErrorContext(
			r.Context(), "json encoding failed", slog.Any("error", err),
		)

		http.Error(
			w, "internal server error", http.StatusInternalServerError,
		)

		return
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "BlackDog")

	ts, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		h.svc.logger.ErrorContext(
			r.Context(), "template parsing failed", slog.Any("error", err),
		)

		http.Error(
			w, "internal server error", http.StatusInternalServerError,
		)

		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		h.svc.logger.ErrorContext(
			r.Context(), "template rendering failed", slog.Any("error", err),
		)

		http.Error(
			w, "internal server error", http.StatusInternalServerError,
		)

		return
	}
}

func (h *Handler) NewApplicationForm(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "New Application")
}

func (h *Handler) SaveApplicationForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Save Application "+id)
}

func (h *Handler) SubmitApplicationForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.WriteHeader(http.StatusAccepted)
	:wa

	io.WriteString(w, "Submit Application "+id)
}

func (h *Handler) ViewApplicationForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}

	io.WriteString(w, "View Application "+strconv.Itoa(id))
}
