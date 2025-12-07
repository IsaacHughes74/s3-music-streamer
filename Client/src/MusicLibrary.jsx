import { useState, useEffect } from 'react'
import './MusicLibrary.css'
import { fetchArtists, fetchAlbums, fetchSongs } from './api'
import ArtistManager from './ArtistManager'
import AlbumManager from './AlbumManager'
import SongUploader from './SongUploader'

function MusicLibrary({ onPlaySong, currentSong, isPlaying }) {
  const [artists, setArtists] = useState([])
  const [selectedArtist, setSelectedArtist] = useState(null)
  const [albums, setAlbums] = useState([])
  const [selectedAlbum, setSelectedAlbum] = useState(null)
  const [songs, setSongs] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [showArtistModal, setShowArtistModal] = useState(false)
  const [showAlbumModal, setShowAlbumModal] = useState(false)
  const [showSongModal, setShowSongModal] = useState(false)

  useEffect(() => {
    loadArtists()
  }, [])

  useEffect(() => {
    if (selectedArtist) {
      loadAlbums(selectedArtist.id)
    } else {
      setAlbums([])
      setSelectedAlbum(null)
      setSongs([])
    }
  }, [selectedArtist])

  useEffect(() => {
    if (selectedAlbum) {
      loadSongs(selectedAlbum.id)
    } else {
      setSongs([])
    }
  }, [selectedAlbum])

  const loadArtists = async () => {
    try {
      setLoading(true)
      const data = await fetchArtists()
      setArtists(data || [])
      setError(null)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const loadAlbums = async (artistId) => {
    try {
      const data = await fetchAlbums(artistId)
      setAlbums(data || [])
      setError(null)
    } catch (err) {
      setError(err.message)
    }
  }

  const loadSongs = async (albumId) => {
    try {
      const data = await fetchSongs(albumId)
      setSongs(data || [])
      setError(null)
    } catch (err) {
      setError(err.message)
    }
  }

  const handleArtistCreated = (newArtist) => {
    setArtists([...artists, newArtist])
    setSelectedArtist(newArtist)
  }

  const handleAlbumCreated = (newAlbum) => {
    setAlbums([...albums, newAlbum])
    setSelectedAlbum(newAlbum)
  }

  const handleSongUploaded = (newSong) => {
    setSongs([...songs, newSong])
    loadSongs(selectedAlbum.id) // Reload to get proper ordering
  }

  if (loading) {
    return (
      <div className="library">
        <div className="loader">Loading library...</div>
      </div>
    )
  }

  return (
    <div className="library">
      <header className="library-header">
        <h1>Music Library</h1>
      </header>

      <div className="library-content">
        {/* Artists Panel */}
        <div className="panel artists-panel">
          <div className="panel-header">
            <div className="panel-header-top">
              <h2>Artists</h2>
            </div>
            <button
              className="add-button"
              onClick={() => setShowArtistModal(true)}
            >
              New Artist
            </button>
          </div>
          <div className="panel-list">
            {artists.length === 0 ? (
              <div className="empty-state">No artists yet</div>
            ) : (
              artists.map((artist) => (
                <div
                  key={artist.id}
                  className={`list-item ${selectedArtist?.id === artist.id ? 'selected' : ''}`}
                  onClick={() => setSelectedArtist(artist)}
                >
                  <div className="item-name">{artist.name}</div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Albums Panel */}
        <div className="panel albums-panel">
          <div className="panel-header">
            <div className="panel-header-top">
              <h2>Albums</h2>
            </div>
            {selectedArtist && (
              <button
                className="add-button"
                onClick={() => setShowAlbumModal(true)}
              >
                New Album
              </button>
            )}
          </div>
          <div className="panel-list">
            {!selectedArtist ? (
              <div className="empty-state">Select an artist</div>
            ) : albums.length === 0 ? (
              <div className="empty-state">No albums yet</div>
            ) : (
              albums.map((album) => (
                <div
                  key={album.id}
                  className={`list-item ${selectedAlbum?.id === album.id ? 'selected' : ''}`}
                  onClick={() => setSelectedAlbum(album)}
                >
                  <div className="item-name">{album.title}</div>
                  {album.year && (
                    <div className="item-meta">{album.year}</div>
                  )}
                </div>
              ))
            )}
          </div>
        </div>

        {/* Songs Panel */}
        <div className="panel songs-panel">
          <div className="panel-header">
            <div className="panel-header-top">
              <h2>Songs</h2>
            </div>
            {selectedAlbum && (
              <button
                className="add-button"
                onClick={() => setShowSongModal(true)}
              >
                Upload Song
              </button>
            )}
          </div>
          <div className="panel-list">
            {!selectedAlbum ? (
              <div className="empty-state">Select an album</div>
            ) : songs.length === 0 ? (
              <div className="empty-state">No songs yet</div>
            ) : (
              songs.map((song) => {
                const isCurrentSong = currentSong?.id === song.id
                return (
                  <div
                    key={song.id}
                    className={`list-item song-item ${isCurrentSong ? 'playing' : ''}`}
                    onClick={() => onPlaySong && onPlaySong(song)}
                  >
                    {song.track_number && (
                      <div className="track-number">{song.track_number}</div>
                    )}
                    <div className="item-name">{song.title}</div>
                    {isCurrentSong && (
                      <div className="now-playing-indicator">
                        {isPlaying ? '♫' : '⏸'}
                      </div>
                    )}
                    {song.duration > 0 && (
                      <div className="item-meta">
                        {Math.floor(song.duration / 60)}:{(song.duration % 60).toString().padStart(2, '0')}
                      </div>
                    )}
                  </div>
                )
              })
            )}
          </div>
        </div>
      </div>

      {error && (
        <div className="error-banner">
          {error}
        </div>
      )}

      {/* Modals */}
      {showArtistModal && (
        <ArtistManager
          onArtistCreated={handleArtistCreated}
          onClose={() => setShowArtistModal(false)}
        />
      )}

      {showAlbumModal && selectedArtist && (
        <AlbumManager
          artist={selectedArtist}
          onAlbumCreated={handleAlbumCreated}
          onClose={() => setShowAlbumModal(false)}
        />
      )}

      {showSongModal && selectedAlbum && selectedArtist && (
        <SongUploader
          album={selectedAlbum}
          artist={selectedArtist}
          onSongUploaded={handleSongUploaded}
          onClose={() => setShowSongModal(false)}
        />
      )}
    </div>
  )
}

export default MusicLibrary
