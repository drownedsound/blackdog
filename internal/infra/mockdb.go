package infra

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
)

type MockDb struct {
	mu    sync.RWMutex
	store map[int64][]byte // Store as []bytes to force deep-copy
}

// NewMockDb creates a thread-safe in-memory store.
func NewMockDb() *MockDb {
	return &MockDb{
		store: make(map[int64][]byte),
	}
}

// Save persists the application.
func (r *MockDb) Insert(ctx context.Context, a *app.Application) error {
	// Create error scenario to force an HTTP 500 on the handler
	if a.MemberReferenceNo == "ERR-100" {
		return app.ErrConnectionRefused
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if a.Id <= 0 {
		return app.ErrInvalidUUID
	}

	now := time.Now().UTC()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	// Serialize to simulate DB storage
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	r.store[a.Id] = data
	return nil
}

// GetByInternalId retrieves an application by id.
func (r *MockDb) GetByInternalId(ctx context.Context, id int64) (
	app.Application, error,
) {
	if id <= 0 {
		return app.Application{}, app.ErrInvalidUUID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	data, ok := r.store[id]
	if !ok {
		return app.Application{}, app.ErrNotFound
	}

	var a app.Application
	if err := json.Unmarshal(data, &a); err != nil {
		return app.Application{}, err
	}

	return a, nil
}
