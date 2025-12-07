// API service for communicating with the backend

const API_BASE = '/api/v1'

// Artists
export const fetchArtists = async () => {
  const response = await fetch(`${API_BASE}/artists`)
  if (!response.ok) throw new Error('Failed to fetch artists')
  return response.json()
}

export const createArtist = async (artistData) => {
  const response = await fetch(`${API_BASE}/artists`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(artistData),
  })
  if (!response.ok) throw new Error('Failed to create artist')
  return response.json()
}

export const updateArtist = async (id, artistData) => {
  const response = await fetch(`${API_BASE}/artists/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(artistData),
  })
  if (!response.ok) throw new Error('Failed to update artist')
  return response.json()
}

export const deleteArtist = async (id) => {
  const response = await fetch(`${API_BASE}/artists/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) throw new Error('Failed to delete artist')
}

// Albums
export const fetchAlbums = async (artistId = null) => {
  const url = artistId
    ? `${API_BASE}/albums?artist_id=${artistId}`
    : `${API_BASE}/albums`
  const response = await fetch(url)
  if (!response.ok) throw new Error('Failed to fetch albums')
  return response.json()
}

export const createAlbum = async (albumData) => {
  const response = await fetch(`${API_BASE}/albums`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(albumData),
  })
  if (!response.ok) throw new Error('Failed to create album')
  return response.json()
}

export const updateAlbum = async (id, albumData) => {
  const response = await fetch(`${API_BASE}/albums/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(albumData),
  })
  if (!response.ok) throw new Error('Failed to update album')
  return response.json()
}

export const deleteAlbum = async (id) => {
  const response = await fetch(`${API_BASE}/albums/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) throw new Error('Failed to delete album')
}

// Songs
export const fetchSongs = async (albumId = null, artistId = null) => {
  let url = `${API_BASE}/songs`
  const params = new URLSearchParams()
  if (albumId) params.append('album_id', albumId)
  else if (artistId) params.append('artist_id', artistId)
  if (params.toString()) url += `?${params.toString()}`

  const response = await fetch(url)
  if (!response.ok) throw new Error('Failed to fetch songs')
  return response.json()
}

export const uploadSong = async (formData) => {
  const response = await fetch(`${API_BASE}/songs/upload`, {
    method: 'POST',
    body: formData,
  })
  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(errorText || 'Upload failed')
  }
  return response.json()
}

export const updateSong = async (id, songData) => {
  const response = await fetch(`${API_BASE}/songs/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(songData),
  })
  if (!response.ok) throw new Error('Failed to update song')
  return response.json()
}

export const deleteSong = async (id) => {
  const response = await fetch(`${API_BASE}/songs/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) throw new Error('Failed to delete song')
}
