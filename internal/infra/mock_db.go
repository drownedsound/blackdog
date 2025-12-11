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
func (r *MockDb) Save(ctx context.Context, app *app.Application) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Simulate Auto-Increment
	if app.Id == 0 {
		r.lastId++
		app.Id = r.lastId
		app.CreatedAt = time.Now().UTC()
	}
	app.UpdatedAt = time.Now().UTC()

	// Serialize to simulate DB storage
	data, err := json.Marshal(app)
	if err != nil {
		return err
	}

	r.store[app.Id] = data
	return nil
}

// GetById retrieves an application by id.
func (r *MockDb) GetById(ctx context.Context, id int64) (
	*app.Application, error,
) {
	if id <= 0 {
		return nil, app.ErrInvalidId
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	data, ok := r.store[id]
	if !ok {
		return nil, app.ErrNotFound
	}

	// Deserialize to return a fresh instance
	var app app.Application
	if err := json.Unmarshal(data, &app); err != nil {
		return nil, err
	}

	return &app, nil
}
