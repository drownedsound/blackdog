package app

import "context"

type Repository interface {
	Save(ctx context.Context, a *Application) error
	GetById(ctx context.Context, id int64) (Application, error)
}
