package models

import (
	"fmt"
	"time"
)

type Song struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	Album       string    `json:"album"`
	Duration    int       `json:"duration"` // duration in seconds
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetS3Key returns the S3 key for this song based on its ID
func (s *Song) GetS3Key() string {
	return fmt.Sprintf("songs/%d/song.mp3", s.ID)
}
