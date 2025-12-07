import { useState } from 'react'
import './ArtistManager.css'
import { createArtist } from './api'

function ArtistManager({ onArtistCreated, onClose }) {
  const [name, setName] = useState('')
  const [bio, setBio] = useState('')
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!name.trim()) {
      setError('Artist name is required')
      return
    }

    try {
      setCreating(true)
      setError(null)

      const artistData = {
        name: name.trim(),
        bio: bio.trim() || undefined,
      }

      const newArtist = await createArtist(artistData)

      // Reset form
      setName('')
      setBio('')

      // Notify parent component
      if (onArtistCreated) {
        onArtistCreated(newArtist)
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
          <h2>Create Artist</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="form">
          <div className="form-group">
            <label htmlFor="name">Artist Name *</label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g., The Beatles"
              required
              autoFocus
            />
          </div>

          <div className="form-group">
            <label htmlFor="bio">Bio (optional)</label>
            <textarea
              id="bio"
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              placeholder="Brief artist biography..."
              rows="4"
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
              disabled={creating || !name.trim()}
            >
              {creating ? 'Creating...' : 'Create Artist'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default ArtistManager
