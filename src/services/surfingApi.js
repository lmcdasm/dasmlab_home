import axios from 'axios'

const SURFING_HOST = import.meta.env.VITE_SURFING_API_HOST || '/api/surfing'

const client = axios.create({
  baseURL: SURFING_HOST,
  timeout: 120000,
  withCredentials: true
})

export function mediaUrl(path) {
  if (!path) return ''
  if (/^https?:\/\//.test(path)) return path
  return `${SURFING_HOST}${path}`
}

export function mediaOutboundUrl(item) {
  if (item?.external_url) return item.external_url
  return mediaUrl(item?.url)
}

export function mediaDownloadUrl(item) {
  // Absolute CDN URLs: hit the edge directly (no /api proxy, no ?download=1 on R2).
  const url = mediaUrl(item?.url)
  if (!url) return ''
  if (/^https?:\/\//.test(url)) return url
  const joiner = url.includes('?') ? '&' : '?'
  return `${url}${joiner}download=1`
}

export function mediaKind(item) {
  const kind = (item?.kind || '').toLowerCase()
  if (kind === 'photo' || kind === 'video' || kind === 'other') return kind
  if (item?.media_type === 'video') return 'video'
  if (item?.media_type === 'other' || item?.external_url) return 'other'
  return 'photo'
}

export async function fetchDays() {
  const { data } = await client.get('/days')
  return sortDays(data || [])
}

export async function createDay(payload) {
  const { data } = await client.post('/days', payload)
  return data
}

export async function deleteDay(dayId) {
  await client.delete(`/days/${dayId}`)
}

export async function uploadMedia(dayId, file, meta = {}, onProgress) {
  const form = new FormData()
  form.append('file', file)
  if (typeof meta === 'string') {
    if (meta) form.append('caption', meta)
  } else {
    if (meta.caption) form.append('caption', meta.caption)
    if (meta.notes) form.append('notes', meta.notes)
    if (meta.kind) form.append('kind', meta.kind)
  }

  const { data } = await client.post(`/days/${dayId}/media`, form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (event) => {
      if (!onProgress || !event.total) return
      onProgress(Math.round((event.loaded / event.total) * 100))
    }
  })
  return data
}

export async function addMediaLink(dayId, payload) {
  const { data } = await client.post(`/days/${dayId}/media/link`, payload)
  return data
}

export async function updateMedia(dayId, mediaId, patch) {
  const { data } = await client.patch(`/days/${dayId}/media/${mediaId}`, patch)
  return data
}

export async function deleteMedia(dayId, mediaId) {
  await client.delete(`/days/${dayId}/media/${mediaId}`)
}

export async function generateTheme(dayId, payload = {}) {
  const { data } = await client.post(`/days/${dayId}/theme/generate`, payload, {
    timeout: 180000
  })
  return data
}

export async function createShare(payload) {
  const { data } = await client.post('/shares', payload)
  return data
}

export async function recordMediaPlay(dayId, mediaId) {
  const { data } = await client.post(`/days/${dayId}/media/${mediaId}/play`)
  return data
}

export async function proposeMediaTag(dayId, mediaId, name) {
  const { data } = await client.post(`/days/${dayId}/media/${mediaId}/tags`, { name })
  return data
}

export async function moderateMediaTag(dayId, mediaId, tagId, action) {
  const { data } = await client.post(`/days/${dayId}/media/${mediaId}/tags/${tagId}/${action}`)
  return data
}

export async function publishDay(dayId, cleanupPvc = false) {
  const q = cleanupPvc ? '?cleanup_pvc=true' : ''
  const { data } = await client.post(`/days/${dayId}/publish${q}`, null, { timeout: 300000 })
  return data
}

export async function curatePublish(dayId, payload = {}) {
  const { data } = await client.post(`/days/${dayId}/curate/publish`, payload, { timeout: 300000 })
  return data
}

export async function aiCurate(dayId, payload = {}) {
  const { data } = await client.post(`/days/${dayId}/ai/curate`, payload, { timeout: 180000 })
  return data
}

export async function unhideMedia(dayId, mediaId) {
  const { data } = await client.post(`/days/${dayId}/media/${mediaId}/unhide`)
  return data
}

export async function patchDay(dayId, payload) {
  const { data } = await client.patch(`/days/${dayId}`, payload)
  return data
}

export async function fetchShareMeta(token) {
  const { data } = await client.get(`/shares/${token}`)
  return data
}

function sortDays(days) {
  return [...days].sort((a, b) => {
    const dateCmp = (b.date || '').localeCompare(a.date || '')
    if (dateCmp !== 0) return dateCmp
    return (b.created_at || '').localeCompare(a.created_at || '')
  })
}
