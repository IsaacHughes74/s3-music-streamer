import { useState } from 'react'
import './UploadSong.css'

function UploadSong({ onUploadComplete, onClose }) {
  const [file, setFile] = useState(null)
  const [title, setTitle] = useState('')
  const [artist, setArtist] = useState('')
  const [album, setAlbum] = useState('')
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
      formData.append('artist', artist)
      formData.append('album', album)

      const response = await fetch('/api/v1/songs/upload', {
        method: 'POST',
        body: formData,
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(errorText || 'Upload failed')
      }

      const uploadedSong = await response.json()

      // Reset form
      setFile(null)
      setTitle('')
      setArtist('')
      setAlbum('')

      // Notify parent component
      if (onUploadComplete) {
        onUploadComplete(uploadedSong)
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
    <div className="upload-modal-overlay" onClick={onClose}>
      <div className="upload-modal" onClick={(e) => e.stopPropagation()}>
        <div className="upload-header">
          <h2>Upload Song</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="upload-form">
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
            <label htmlFor="artist">Artist</label>
            <input
              type="text"
              id="artist"
              value={artist}
              onChange={(e) => setArtist(e.target.value)}
              placeholder="Artist name"
            />
          </div>

          <div className="form-group">
            <label htmlFor="album">Album</label>
            <input
              type="text"
              id="album"
              value={album}
              onChange={(e) => setAlbum(e.target.value)}
              placeholder="Album name"
            />
          </div>

          {error && (
            <div className="upload-error">
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
              className="upload-button"
              disabled={uploading || !file}
            >
              {uploading ? 'Uploading...' : 'Upload'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default UploadSong
