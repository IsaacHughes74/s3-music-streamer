package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"s3-music-streamer/internal/database"
	"s3-music-streamer/internal/models"
	"s3-music-streamer/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *database.DB
	s3 *storage.S3Client
}

func New(db *database.DB, s3 *storage.S3Client) *Handler {
	return &Handler{
		db: db,
		s3: s3,
	}
}

func (h *Handler) ListSongs(w http.ResponseWriter, r *http.Request) {
	// Support filtering by artist_id or album_id via query params
	artistID := r.URL.Query().Get("artist_id")
	albumID := r.URL.Query().Get("album_id")

	query := `
		SELECT s.id, s.title, s.artist_id, s.album_id, s.track_number, s.duration, s.file_size,
		       s.content_type, s.created_at, s.updated_at,
		       ar.name as artist_name, al.title as album_title
		FROM songs s
		LEFT JOIN artists ar ON s.artist_id = ar.id
		LEFT JOIN albums al ON s.album_id = al.id
	`

	var rows *sql.Rows
	var err error

	if artistID != "" {
		query += " WHERE s.artist_id = ? ORDER BY s.created_at DESC"
		rows, err = h.db.Query(query, artistID)
	} else if albumID != "" {
		query += " WHERE s.album_id = ? ORDER BY s.track_number ASC, s.created_at DESC"
		rows, err = h.db.Query(query, albumID)
	} else {
		query += " ORDER BY s.created_at DESC"
		rows, err = h.db.Query(query)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	songs := []models.Song{}
	for rows.Next() {
		var song models.Song
		var artistName, albumTitle sql.NullString
		var trackNumber sql.NullInt64
		if err := rows.Scan(
			&song.ID, &song.Title, &song.ArtistID, &song.AlbumID, &trackNumber, &song.Duration,
			&song.FileSize, &song.ContentType, &song.CreatedAt, &song.UpdatedAt,
			&artistName, &albumTitle,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if trackNumber.Valid {
			track := int(trackNumber.Int64)
			song.TrackNumber = &track
		}
		if artistName.Valid {
			song.ArtistName = artistName.String
		}
		if albumTitle.Valid {
			song.AlbumTitle = albumTitle.String
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *Handler) GetSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid song id", http.StatusBadRequest)
		return
	}

	var song models.Song
	var artistName, albumTitle sql.NullString
	var trackNumber sql.NullInt64
	err = h.db.QueryRow(`
		SELECT s.id, s.title, s.artist_id, s.album_id, s.track_number, s.duration, s.file_size,
		       s.content_type, s.created_at, s.updated_at,
		       ar.name as artist_name, al.title as album_title
		FROM songs s
		LEFT JOIN artists ar ON s.artist_id = ar.id
		LEFT JOIN albums al ON s.album_id = al.id
		WHERE s.id = ?
	`, id).Scan(
		&song.ID, &song.Title, &song.ArtistID, &song.AlbumID, &trackNumber, &song.Duration,
		&song.FileSize, &song.ContentType, &song.CreatedAt, &song.UpdatedAt,
		&artistName, &albumTitle,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "song not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if trackNumber.Valid {
		track := int(trackNumber.Int64)
		song.TrackNumber = &track
	}
	if artistName.Valid {
		song.ArtistName = artistName.String
	}
	if albumTitle.Valid {
		song.AlbumTitle = albumTitle.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
		INSERT INTO songs (title, artist_id, album_id, track_number, duration, file_size, content_type)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, song.Title, song.ArtistID, song.AlbumID, song.TrackNumber, song.Duration, song.FileSize, song.ContentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	song.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid song id", http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
		UPDATE songs
		SET title = ?, artist_id = ?, album_id = ?, track_number = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, song.Title, song.ArtistID, song.AlbumID, song.TrackNumber, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "song not found", http.StatusNotFound)
		return
	}

	song.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid song id", http.StatusBadRequest)
		return
	}

	// Check if song exists
	var exists bool
	err = h.db.QueryRow("SELECT 1 FROM songs WHERE id = ?", id).Scan(&exists)
	if err == sql.ErrNoRows {
		http.Error(w, "song not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate S3 key from ID
	s3Key := fmt.Sprintf("songs/%d/song.mp3", id)

	if err := h.s3.DeleteObject(r.Context(), s3Key); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete from S3: %v", err), http.StatusInternalServerError)
		return
	}

	if _, err := h.db.Exec("DELETE FROM songs WHERE id = ?", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) StreamSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid song id", http.StatusBadRequest)
		return
	}

	var contentType string
	err = h.db.QueryRow("SELECT content_type FROM songs WHERE id = ?", id).Scan(&contentType)
	if err == sql.ErrNoRows {
		http.Error(w, "song not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate S3 key from ID
	s3Key := fmt.Sprintf("songs/%d/song.mp3", id)

	body, err := h.s3.GetObject(r.Context(), s3Key)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get from S3: %v", err), http.StatusInternalServerError)
		return
	}
	defer body.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")

	if _, err := io.Copy(w, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UploadSong(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100 MB max
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	artistIDStr := r.FormValue("artist_id")
	albumIDStr := r.FormValue("album_id")
	trackNumberStr := r.FormValue("track_number")
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "audio/mpeg"
	}

	var artistID, albumID *int64
	var trackNumber *int
	if artistIDStr != "" {
		id, err := strconv.ParseInt(artistIDStr, 10, 64)
		if err == nil {
			artistID = &id
		}
	}
	if albumIDStr != "" {
		id, err := strconv.ParseInt(albumIDStr, 10, 64)
		if err == nil {
			albumID = &id
		}
	}
	if trackNumberStr != "" {
		track, err := strconv.Atoi(trackNumberStr)
		if err == nil {
			trackNumber = &track
		}
	}

	// Step 1: Insert into database first to get an ID
	result, err := h.db.Exec(`
		INSERT INTO songs (title, artist_id, album_id, track_number, file_size, content_type)
		VALUES (?, ?, ?, ?, ?, ?)
	`, title, artistID, albumID, trackNumber, header.Size, contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Step 2: Generate S3 key with the ID (fixed path: songs/{id}/song.mp3)
	s3Key := fmt.Sprintf("songs/%d/song.mp3", id)

	// Step 3: Upload to S3
	if err := h.s3.PutObject(r.Context(), s3Key, file, contentType); err != nil {
		// Step 4: Delete from database if S3 upload failed
		h.db.Exec("DELETE FROM songs WHERE id = ?", id)
		http.Error(w, fmt.Sprintf("failed to upload to S3: %v", err), http.StatusInternalServerError)
		return
	}

	song := models.Song{
		ID:          id,
		Title:       title,
		ArtistID:    artistID,
		AlbumID:     albumID,
		TrackNumber: trackNumber,
		FileSize:    header.Size,
		ContentType: contentType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}
