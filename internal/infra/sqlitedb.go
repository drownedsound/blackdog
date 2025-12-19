package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
	"github.com/mattn/go-sqlite3"
)

type SqliteDb struct {
	db *sql.DB
}

func NewSqliteDb(db *sql.DB) *SqliteDb {
	return &SqliteDb{
		db: db,
	}
}

func NewSQLiteConnection(dbPath string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"file:%s?"+
			"_journal_mode=WAL&"+
			"_busy_timeout=5000&"+
			"_synchronous=NORMAL&"+
			"_foreign_keys=on",
		dbPath,
	)

	// sql.Open only validates arguments. With a hardcoded driver and DSN format,
	// this cannot fail. We ignore the error and rely on Ping to verify the
	// actual connection and file permissions.
	db, _ := sql.Open("sqlite3", dsn)

	// Connection Pool Settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (r *SqliteDb) Save(ctx context.Context, a *app.Application) error {
	const query = `
		INSERT INTO APPLICATION (
		member_reference_no, category_id, status_id, requested_amount, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)`

	// Use ExecContext directly. Previously, PrepareContext+Close
	// incurred overhead without statement reuse.
	// This also removes the unreachable 'Prepare' error branch.
	res, err := r.db.ExecContext(
		ctx,
		query,
		a.MemberReferenceNo,
		int(a.CategoryCode),
		int(a.StatusCode),
		a.RequestedAmount,
		a.CreatedAt.Unix(),
		a.UpdatedAt.Unix(),
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("duplicate application: %w", err)
			}
		}
		return fmt.Errorf("%w: %v", app.ErrInsertFailed, err)
	}

	id, _ := res.LastInsertId()
	a.Id = id
	return nil
}

func (r *SqliteDb) GetById(ctx context.Context, id int64) (app.Application, error) {
	const query = `
		SELECT 
		id, member_reference_no, category_id, status_id, 
		requested_amount, created_at, updated_at

		FROM APPLICATION 

		WHERE id = ?`

	var (
		appEntity     app.Application
		categoryId    int
		statusId      int
		createdAtUnix int64
		updatedAtUnix int64
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&appEntity.Id,
		&appEntity.MemberReferenceNo,
		&categoryId,
		&statusId,
		&appEntity.RequestedAmount,
		&createdAtUnix,
		&updatedAtUnix,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.Application{}, app.ErrNotFound
		}
		return app.Application{}, fmt.Errorf("query failed: %w", err)
	}

	appEntity.CreatedAt = time.Unix(createdAtUnix, 0).UTC()
	appEntity.UpdatedAt = time.Unix(updatedAtUnix, 0).UTC()
	appEntity.CategoryCode = app.ProductCategory(categoryId)
	appEntity.StatusCode = app.ApplicationStatus(statusId)

	return appEntity, nil
}
