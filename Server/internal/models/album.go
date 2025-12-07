package models

import "time"

type Album struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	ArtistID   int64     `json:"artist_id"`
	ArtistName string    `json:"artist_name,omitempty"` // For joined queries
	Year       int       `json:"year,omitempty"`
	CoverArt   string    `json:"cover_art,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
