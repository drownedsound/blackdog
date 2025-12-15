package infra

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
)

type MockDb struct {
	mu     sync.RWMutex
	store  map[int64][]byte // Store as []bytes to force deep-copy
	lastId int64
}

// NewMockDb creates a thread-safe in-memory store.
func NewMockDb() *MockDb {
	return &MockDb{
		store:  make(map[int64][]byte),
		lastId: 0,
	}
}

// Save persists the application.
// If Id is 0, it acts as an INSERT and assigns a new Id.
// If Id > 0, it acts as an UPDATE.
func (r *MockDb) Save(ctx context.Context, a *app.Application) error {
	// Create error scenario to force an HTTP 500 on the handler
	if a.MemberReferenceNo == "ERR-100" {
		return app.ErrConnectionRefused
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Simulate Auto-Increment
	if a.Id == 0 {
		r.lastId++
		a.Id = r.lastId
		a.CreatedAt = time.Now().UTC()
	}
	a.UpdatedAt = time.Now().UTC()

	// Serialize to simulate DB storage
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	r.store[a.Id] = data
	return nil
}

// GetById retrieves an application by id.
func (r *MockDb) GetById(ctx context.Context, id int64) (
	app.Application, error,
) {
	if id <= 0 {
		return app.Application{}, app.ErrInvalidId
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	data, ok := r.store[id]
	if !ok {
		return app.Application{}, app.ErrNotFound
	}

	// Deserialize to return a fresh instance
	var a app.Application
	if err := json.Unmarshal(data, &a); err != nil {
		return app.Application{}, err
	}

	return a, nil
}
