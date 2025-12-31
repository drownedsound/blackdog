package app

import (
	"context"
)

type mockRepo struct {
	insertFunc          func(ctx context.Context, a *Application) error
	getByInternalIdFunc func(ctx context.Context, id int64) (
		*Application,
		error,
	)
}

func (m *mockRepo) Insert(ctx context.Context, a *Application) error {
	if m.insertFunc != nil {
		return m.insertFunc(ctx, a)
	}
	return nil
}

func (m *mockRepo) GetByInternalId(ctx context.Context, id int64) (
	*Application,
	error,
) {
	if m.getByInternalIdFunc != nil {
		return m.getByInternalIdFunc(ctx, id)
	}

	return nil, ErrNotFound
}
