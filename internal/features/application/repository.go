package app

import (
	"context"
	"errors"
)

var (
	ErrNotFound          = errors.New("application not found")
	ErrInsertFailed      = errors.New("record insert failed")
	ErrConnectionRefused = errors.New("database connection refused")
)

type Repository interface {
	Save(ctx context.Context, a *Application) error
	GetById(ctx context.Context, id int64) (Application, error)
}
