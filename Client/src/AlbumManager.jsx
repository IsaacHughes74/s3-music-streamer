import { useState } from 'react'
import './AlbumManager.css'
import { createAlbum } from './api'

function AlbumManager({ artist, onAlbumCreated, onClose }) {
  const [title, setTitle] = useState('')
  const [year, setYear] = useState('')
  const [coverArt, setCoverArt] = useState('')
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!title.trim()) {
      setError('Album title is required')
      return
    }

    try {
      setCreating(true)
      setError(null)

      const albumData = {
        title: title.trim(),
        artist_id: artist.id,
        year: year ? parseInt(year) : undefined,
        cover_art: coverArt.trim() || undefined,
      }

      const newAlbum = await createAlbum(albumData)

      // Reset form
      setTitle('')
      setYear('')
      setCoverArt('')

      // Notify parent component
      if (onAlbumCreated) {
        onAlbumCreated(newAlbum)
      }

      // Close modal after successful creation
      if (onClose) {
        onClose()
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Create Album for {artist.name}</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="form">
          <div className="form-group">
            <label htmlFor="title">Album Title *</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="e.g., Abbey Road"
              required
              autoFocus
            />
          </div>

          <div className="form-group">
            <label htmlFor="year">Year (optional)</label>
            <input
              type="number"
              id="year"
              value={year}
              onChange={(e) => setYear(e.target.value)}
              placeholder="e.g., 1969"
              min="1900"
              max={new Date().getFullYear() + 1}
            />
          </div>

          <div className="form-group">
            <label htmlFor="coverArt">Cover Art URL (optional)</label>
            <input
              type="text"
              id="coverArt"
              value={coverArt}
              onChange={(e) => setCoverArt(e.target.value)}
              placeholder="https://..."
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
              disabled={creating}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="submit-button"
              disabled={creating || !title.trim()}
            >
              {creating ? 'Creating...' : 'Create Album'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default AlbumManager
