import { useState, useEffect, useRef } from 'react'
import './App.css'
import UploadSong from './UploadSong'

function App() {
  const [songs, setSongs] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [currentSong, setCurrentSong] = useState(null)
  const [isPlaying, setIsPlaying] = useState(false)
  const [showUpload, setShowUpload] = useState(false)
  const [currentTime, setCurrentTime] = useState(0)
  const [duration, setDuration] = useState(0)
  const audioRef = useRef(null)

  useEffect(() => {
    fetchSongs()
  }, [])

  const fetchSongs = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/v1/songs')
      if (!response.ok) {
        throw new Error('Failed to fetch songs')
      }
      const data = await response.json()
      setSongs(data)
      setError(null)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const playSong = (song) => {
    if (currentSong?.id === song.id) {
      if (isPlaying) {
        audioRef.current.pause()
        setIsPlaying(false)
      } else {
        audioRef.current.play()
        setIsPlaying(true)
      }
    } else {
      setCurrentSong(song)
      setIsPlaying(true)
      setTimeout(() => {
        audioRef.current.play()
      }, 100)
    }
  }

  const formatDuration = (seconds) => {
    if (!seconds) return '0:00'
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, '0')}`
  }

  const formatFileSize = (bytes) => {
    if (!bytes) return '0 MB'
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  const handleUploadComplete = (uploadedSong) => {
    // Add the new song to the list
    setSongs([uploadedSong, ...songs])
  }

  const playNextSong = () => {
    if (!currentSong || songs.length === 0) return

    const currentIndex = songs.findIndex(song => song.id === currentSong.id)
    const nextIndex = (currentIndex + 1) % songs.length // Loop back to start
    playSong(songs[nextIndex])
  }

  const playPreviousSong = () => {
    if (!currentSong || songs.length === 0) return

    const currentIndex = songs.findIndex(song => song.id === currentSong.id)
    const previousIndex = currentIndex === 0 ? songs.length - 1 : currentIndex - 1
    playSong(songs[previousIndex])
  }

  if (loading) {
    return (
      <div className="app">
        <div className="loader">Loading songs...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="app">
        <div className="error">
          <p>Error: {error}</p>
          <button onClick={fetchSongs}>Retry</button>
        </div>
      </div>
    )
  }

  return (
    <div className="app">
      <header className="header">
        <h1>Music Streamer</h1>
        <p className="song-count">{songs.length} songs</p>
        <button className="upload-trigger-button" onClick={() => setShowUpload(true)}>
          + Upload Song
        </button>
      </header>

      {songs.length === 0 ? (
        <div className="empty-state">
          <p>No songs available yet</p>
        </div>
      ) : (
        <div className="song-list">
          {songs.map((song) => {
            const isCurrentSong = currentSong?.id === song.id
            const progress = isCurrentSong && duration > 0 ? currentTime / duration : 0

            return (
              <div
                key={song.id}
                className={`song-item ${isCurrentSong ? 'playing' : ''}`}
                onClick={() => playSong(song)}
              >
                <div className="song-play-icon-container">
                  <svg width="44" height="44" viewBox="0 0 44 44" className="song-progress-ring">
                    <circle
                      cx="22"
                      cy="22"
                      r="18"
                      stroke="rgba(255, 255, 255, 0.15)"
                      strokeWidth="2.5"
                      fill="none"
                    />
                    {isCurrentSong && (
                      <circle
                        cx="22"
                        cy="22"
                        r="18"
                        stroke="white"
                        strokeWidth="2.5"
                        fill="none"
                        strokeDasharray={`${2 * Math.PI * 18}`}
                        strokeDashoffset={`${2 * Math.PI * 18 * (1 - progress)}`}
                        strokeLinecap="round"
                        transform="rotate(-90 22 22)"
                        style={{ transition: 'stroke-dashoffset 0.1s linear' }}
                      />
                    )}
                  </svg>
                  <div className="song-play-icon">
                    {isCurrentSong && isPlaying ? '⏸' : '▶'}
                  </div>
                </div>
                <div className="song-info">
                  <div className="song-title">{song.title || 'Untitled'}</div>
                  <div className="song-meta">
                    <span className="song-artist">{song.artist || 'Unknown Artist'}</span>
                    {song.album && <span className="song-album"> • {song.album}</span>}
                  </div>
                  <div className="song-details">
                    {isCurrentSong && duration > 0 ? (
                      <span className="song-duration">
                        {formatDuration(Math.floor(currentTime))} / {formatDuration(Math.floor(duration))}
                      </span>
                    ) : song.duration > 0 ? (
                      <span className="song-duration">{formatDuration(song.duration)}</span>
                    ) : null}
                    <span className="song-size">{formatFileSize(song.file_size)}</span>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}

      {currentSong && (
        <div className="player">
          <audio
            ref={audioRef}
            src={`/api/v1/songs/${currentSong.id}/stream`}
            onEnded={() => playNextSong()}
            onPause={() => setIsPlaying(false)}
            onPlay={() => setIsPlaying(true)}
            onTimeUpdate={(e) => setCurrentTime(e.target.currentTime)}
            onLoadedMetadata={(e) => setDuration(e.target.duration)}
          />
          <div className="player-progress-ring">
            <svg width="56" height="56" viewBox="0 0 56 56">
              <circle
                cx="28"
                cy="28"
                r="24"
                stroke="rgba(255, 255, 255, 0.1)"
                strokeWidth="3"
                fill="none"
              />
              <circle
                cx="28"
                cy="28"
                r="24"
                stroke="white"
                strokeWidth="3"
                fill="none"
                strokeDasharray={`${2 * Math.PI * 24}`}
                strokeDashoffset={`${2 * Math.PI * 24 * (1 - (duration > 0 ? currentTime / duration : 0))}`}
                strokeLinecap="round"
                transform="rotate(-90 28 28)"
                style={{ transition: 'stroke-dashoffset 0.1s linear' }}
              />
            </svg>
            <button
              className="player-button-icon"
              onClick={() => playSong(currentSong)}
              aria-label={isPlaying ? 'Pause' : 'Play'}
            >
              {isPlaying ? '⏸' : '▶'}
            </button>
          </div>
          <div className="player-info">
            <div className="player-title">{currentSong.title || 'Untitled'}</div>
            <div className="player-artist">{currentSong.artist || 'Unknown Artist'}</div>
            <div className="player-time">
              {formatDuration(Math.floor(currentTime))} / {formatDuration(Math.floor(duration))}
            </div>
          </div>
          <div className="player-controls">
            <button
              className="player-control-button"
              onClick={playPreviousSong}
              aria-label="Previous song"
              disabled={songs.length <= 1}
            >
              ⏮
            </button>
            <button
              className="player-control-button"
              onClick={playNextSong}
              aria-label="Next song"
              disabled={songs.length <= 1}
            >
              ⏭
            </button>
          </div>
        </div>
      )}

      {showUpload && (
        <UploadSong
          onUploadComplete={handleUploadComplete}
          onClose={() => setShowUpload(false)}
        />
      )}
    </div>
  )
}

export default App
