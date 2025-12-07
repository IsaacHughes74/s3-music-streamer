package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"s3-music-streamer/internal/models"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListAlbums(w http.ResponseWriter, r *http.Request) {
	// Support filtering by artist_id via query param
	artistID := r.URL.Query().Get("artist_id")

	query := `
		SELECT a.id, a.title, a.artist_id, a.year, a.cover_art,
		       a.created_at, a.updated_at, ar.name as artist_name
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.id
	`

	var rows *sql.Rows
	var err error

	if artistID != "" {
		query += " WHERE a.artist_id = ? ORDER BY a.year DESC, a.title ASC"
		rows, err = h.db.Query(query, artistID)
	} else {
		query += " ORDER BY a.created_at DESC"
		rows, err = h.db.Query(query)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	albums := []models.Album{}
	for rows.Next() {
		var album models.Album
		var year sql.NullInt64
		var coverArt sql.NullString
		if err := rows.Scan(
			&album.ID, &album.Title, &album.ArtistID, &year, &coverArt,
			&album.CreatedAt, &album.UpdatedAt, &album.ArtistName,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if year.Valid {
			album.Year = int(year.Int64)
		}
		if coverArt.Valid {
			album.CoverArt = coverArt.String
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

func (h *Handler) GetAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid album id", http.StatusBadRequest)
		return
	}

	var album models.Album
	var year sql.NullInt64
	var coverArt sql.NullString
	err = h.db.QueryRow(`
		SELECT a.id, a.title, a.artist_id, a.year, a.cover_art,
		       a.created_at, a.updated_at, ar.name as artist_name
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.id
		WHERE a.id = ?
	`, id).Scan(
		&album.ID, &album.Title, &album.ArtistID, &year, &coverArt,
		&album.CreatedAt, &album.UpdatedAt, &album.ArtistName,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "album not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if year.Valid {
		album.Year = int(year.Int64)
	}
	if coverArt.Valid {
		album.CoverArt = coverArt.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func (h *Handler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var album models.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if album.Title == "" || album.ArtistID == 0 {
		http.Error(w, "title and artist_id are required", http.StatusBadRequest)
		return
	}

	var year *int
	if album.Year > 0 {
		year = &album.Year
	}
	var coverArt *string
	if album.CoverArt != "" {
		coverArt = &album.CoverArt
	}

	result, err := h.db.Exec(`
		INSERT INTO albums (title, artist_id, year, cover_art)
		VALUES (?, ?, ?, ?)
	`, album.Title, album.ArtistID, year, coverArt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	album.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(album)
}

func (h *Handler) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid album id", http.StatusBadRequest)
		return
	}

	var album models.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var year *int
	if album.Year > 0 {
		year = &album.Year
	}
	var coverArt *string
	if album.CoverArt != "" {
		coverArt = &album.CoverArt
	}

	result, err := h.db.Exec(`
		UPDATE albums
		SET title = ?, artist_id = ?, year = ?, cover_art = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, album.Title, album.ArtistID, year, coverArt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "album not found", http.StatusNotFound)
		return
	}

	album.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func (h *Handler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid album id", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("DELETE FROM albums WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "album not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
