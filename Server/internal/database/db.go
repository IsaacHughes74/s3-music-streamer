package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		artist TEXT NOT NULL,
		album TEXT,
		duration INTEGER DEFAULT 0,
		file_size INTEGER DEFAULT 0,
		content_type TEXT DEFAULT 'audio/mpeg',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);
	CREATE INDEX IF NOT EXISTS idx_songs_album ON songs(album);
	CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}
