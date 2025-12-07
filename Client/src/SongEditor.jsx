import { useState } from 'react'
import './SongEditor.css'
import { updateSong } from './api'

function SongEditor({ song, onSongUpdated, onClose }) {
  const [title, setTitle] = useState(song.title || '')
  const [trackNumber, setTrackNumber] = useState(song.track_number || '')
  const [updating, setUpdating] = useState(false)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!title.trim()) {
      setError('Title is required')
      return
    }

    try {
      setUpdating(true)
      setError(null)

      const songData = {
        title: title.trim(),
        artist_id: song.artist_id,
        album_id: song.album_id,
        track_number: trackNumber ? parseInt(trackNumber) : null,
      }

      const updatedSong = await updateSong(song.id, songData)

      // Notify parent component
      if (onSongUpdated) {
        onSongUpdated(updatedSong)
      }

      // Close modal after successful update
      if (onClose) {
        onClose()
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setUpdating(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Edit Song</h2>
          <p className="modal-subtitle">{song.artist_name} - {song.album_title}</p>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="form">
          <div className="form-group">
            <label htmlFor="title">Title *</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Song title"
              required
              autoFocus
            />
          </div>

          <div className="form-group">
            <label htmlFor="trackNumber">Track Number (optional)</label>
            <input
              type="number"
              id="trackNumber"
              value={trackNumber}
              onChange={(e) => setTrackNumber(e.target.value)}
              placeholder="e.g., 1"
              min="1"
            />
          </div>

          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <div className="form-actions">
            <button
              type="button"
              onClick={onClose}
              className="cancel-button"
              disabled={updating}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="submit-button"
              disabled={updating || !title.trim()}
            >
              {updating ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default SongEditor
