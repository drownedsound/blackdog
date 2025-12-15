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

// NewSqliteDb creates a new instance of the repository.
func NewSqliteDb(db *sql.DB) *SqliteDb {
	return &SqliteDb{
		db: db,
	}
}

// NewSQLiteConnection initializes the database connection
// with performance-tuned PRAGMAs.
func NewSQLiteConnection(dbPath string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"file:%s"+"?"+
			// _journal_mode=WAL: Enables Write-Ahead Logging. Non-blocking reads.
			"_journal_mode=WAL&"+"&"+
			// _busy_timeout=5000:  Waits 5s for the write lock before failing.
			"_busy_timeout=5000"+"&"+
			// _synchronous=NORMAL: Faster writes, safe against app crashes.
			"_synchronous=NORMAL"+"&"+
			// _foreign_keys=on:    Enforce schema constraints.
			"_foreign_keys=on",
		dbPath,
	)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Connection Pool Settings for High Throughput
	// Limit open connections to prevent "database is locked" contention
	// even with WAL mode.
	db.SetMaxOpenConns(10) // TODO: Create configuration item
	db.SetMaxIdleConns(5)  // TODO: Create configuration item
	// Reuse connections indefinitely
	db.SetConnMaxLifetime(0) // TODO: Create configuration item

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Save persists the Application entity in the SQLite database
func (r *SqliteDb) Save(ctx context.Context, a *app.Application) error {
	const query = `
		INSERT INTO APPLICATION (
		member_reference_no,
		category_id,
		status_id,
		requested_amount,
		created_at,
		updated_at
		) VALUES (?, ?, ?, ?, ?, ?)`

	// FIXME: Use transactions
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	// Need to case enums to int and time.Time to Unix Epoch (int64)
	res, err := stmt.ExecContext(
		ctx,
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

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	a.Id = id
	return nil
}

// GetById retrieves the Application entity by ID.
func (r *SqliteDb) GetById(ctx context.Context, id int64) (app.Application, error) {
	const query = `
		SELECT 
		id,
		member_reference_no,
		category_id,
		status_id,
		requested_amount,
		created_at,
		updated_at
		FROM APPLICATION
		WHERE id = ?`

	// Scan directly into stack-allocated variables.
	// This avoids allocating a map or intermediate interface{}.
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

	// Rehydrate domain objects
	// Convert optimized storage types (int/int64) back to Domain types.
	appEntity.CreatedAt = time.Unix(createdAtUnix, 0).UTC()
	appEntity.UpdatedAt = time.Unix(updatedAtUnix, 0).UTC()
	appEntity.CategoryCode = app.ProductCategory(categoryId)
	appEntity.StatusCode = app.ApplicationStatus(statusId)

	return appEntity, nil
}
