package infra

import (
	"database/sql"
	"fmt"
	"time"
)

// NewSQLiteConnection creates a connection optimized for multi-instance use.
// Configuration Rationale:
//  1. _journal_mode=WAL: Essential for concurrency.
//  2. _busy_timeout=10000: Crucial for multiple instances. Waits 10s for
//     file locks before failing.
//  3. _synchronous=NORMAL: Faster writes, sacrificing durability only on
//     OS crash (not app crash).
//  4. _txlock=immediate: Prevents deadlocks between instances by acquiring
//     write locks upfront.
func NewSQLiteConnection(dbPath string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"file:%s?_"+
			"journal_mode=WAL&_"+
			"busy_timeout=10000&_"+
			"synchronous=NORMAL&_"+
			"foreign_keys=on&_"+
			"txlock=immediate",
		dbPath,
	)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db connection failed: %w", err)
	}

	// Verify connection immediately (Fail Fast).
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("pinging sqlite db failed: %w", err)
	}

	// Even with multiple instances, restrict process to 1 connection.
	// Allowing 10 connections in this process results in fighting over
	// a single file lock.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}
