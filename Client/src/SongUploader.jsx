import { useState } from 'react'
import './SongUploader.css'
import { uploadSong } from './api'

function SongUploader({ album, artist, onSongUploaded, onClose }) {
  const [file, setFile] = useState(null)
  const [title, setTitle] = useState('')
  const [trackNumber, setTrackNumber] = useState('')
  const [uploading, setUploading] = useState(false)
  const [error, setError] = useState(null)

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0]
    if (selectedFile) {
      setFile(selectedFile)
      // Auto-fill title from filename if empty
      if (!title) {
        const filename = selectedFile.name.replace(/\.[^/.]+$/, '') // Remove extension
        setTitle(filename)
      }
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!file) {
      setError('Please select a file')
      return
    }

    try {
      setUploading(true)
      setError(null)

      const formData = new FormData()
      formData.append('file', file)
      formData.append('title', title)
      formData.append('artist_id', artist.id.toString())
      formData.append('album_id', album.id.toString())
      if (trackNumber) {
        formData.append('track_number', trackNumber)
      }

      const uploadedSong = await uploadSong(formData)

      // Reset form
      setFile(null)
      setTitle('')
      setTrackNumber('')

      // Notify parent component
      if (onSongUploaded) {
        onSongUploaded(uploadedSong)
      }

      // Close modal after successful upload
      if (onClose) {
        onClose()
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Upload Song to {album.title}</h2>
          <p className="modal-subtitle">{artist.name}</p>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="form">
          <div className="form-group">
            <label htmlFor="file">Audio File *</label>
            <input
              type="file"
              id="file"
              accept="audio/*"
              onChange={handleFileChange}
              required
            />
            {file && (
              <div className="file-info">
                Selected: {file.name} ({(file.size / (1024 * 1024)).toFixed(2)} MB)
              </div>
            )}
          </div>

          <div className="form-group">
            <label htmlFor="title">Title *</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Song title"
              required
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
              disabled={uploading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="submit-button"
              disabled={uploading || !file}
            >
              {uploading ? 'Uploading...' : 'Upload Song'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default SongUploader
