package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"s3-music-streamer/internal/models"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListArtists(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT id, name, bio, created_at, updated_at
		FROM artists
		ORDER BY name ASC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	artists := []models.Artist{}
	for rows.Next() {
		var artist models.Artist
		var bio sql.NullString
		if err := rows.Scan(
			&artist.ID, &artist.Name, &bio, &artist.CreatedAt, &artist.UpdatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if bio.Valid {
			artist.Bio = bio.String
		}
		artists = append(artists, artist)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

func (h *Handler) GetArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid artist id", http.StatusBadRequest)
		return
	}

	var artist models.Artist
	var bio sql.NullString
	err = h.db.QueryRow(`
		SELECT id, name, bio, created_at, updated_at
		FROM artists WHERE id = ?
	`, id).Scan(
		&artist.ID, &artist.Name, &bio, &artist.CreatedAt, &artist.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "artist not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bio.Valid {
		artist.Bio = bio.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

func (h *Handler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	var artist models.Artist
	if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if artist.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	var bio *string
	if artist.Bio != "" {
		bio = &artist.Bio
	}

	result, err := h.db.Exec(`
		INSERT INTO artists (name, bio)
		VALUES (?, ?)
	`, artist.Name, bio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	artist.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

func (h *Handler) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid artist id", http.StatusBadRequest)
		return
	}

	var artist models.Artist
	if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var bio *string
	if artist.Bio != "" {
		bio = &artist.Bio
	}

	result, err := h.db.Exec(`
		UPDATE artists
		SET name = ?, bio = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, artist.Name, bio, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "artist not found", http.StatusNotFound)
		return
	}

	artist.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

func (h *Handler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid artist id", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("DELETE FROM artists WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "artist not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
