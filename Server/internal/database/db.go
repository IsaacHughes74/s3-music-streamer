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
	CREATE TABLE IF NOT EXISTS artists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		bio TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		artist_id INTEGER NOT NULL,
		year INTEGER,
		cover_art TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
		UNIQUE(title, artist_id)
	);

	CREATE TABLE IF NOT EXISTS songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		artist_id INTEGER,
		album_id INTEGER,
		track_number INTEGER,
		duration INTEGER DEFAULT 0,
		file_size INTEGER DEFAULT 0,
		content_type TEXT DEFAULT 'audio/mpeg',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE SET NULL,
		FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE SET NULL
	);

	CREATE INDEX IF NOT EXISTS idx_songs_artist_id ON songs(artist_id);
	CREATE INDEX IF NOT EXISTS idx_songs_album_id ON songs(album_id);
	CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
	CREATE INDEX IF NOT EXISTS idx_albums_artist_id ON albums(artist_id);
	CREATE INDEX IF NOT EXISTS idx_artists_name ON artists(name);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}
