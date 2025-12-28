package app

import "context"

type Repository interface {
	Insert(ctx context.Context, a *Application) error
	// BulkInsert(ctx context.Context, a []*Application) error
	// Update(ctx context.Context, a *Application) error
	GetByInternalId(ctx context.Context, id int64) (*Application, error)
	// GetByExternalId(ctx context.Context, mrn int64) (*Application, error)
}
