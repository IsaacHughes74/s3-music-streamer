package models

import "time"

type Artist struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
