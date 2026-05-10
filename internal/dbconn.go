package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DataPath = defaultDataPath()

func defaultDataPath() string {
	if cfg, err := os.UserConfigDir(); err == nil {
		return filepath.Join(cfg, "argonvault", "storage.db")
	}
	return "./data/storage.db"
}

// SQLite tuning: WAL for concurrent reads, busy timeout to ride out contention,
// foreign keys on (off by default in SQLite).
const dsnOptions = "?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on"

func OpenDB() (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(DataPath), 0o700); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	db, err := sql.Open("sqlite3", DataPath+dsnOptions)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return db, nil
}
