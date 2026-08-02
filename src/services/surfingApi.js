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

export function mediaDownloadUrl(item, dayId) {
  // Prefer DASMLAB-gated download path (auth + visibility). Falls back to legacy CDN.
  if (item?.download_path) {
    return `${SURFING_HOST}${item.download_path}`
  }
  if (dayId && item?.id) {
    return `${SURFING_HOST}/days/${dayId}/media/${item.id}/download`
  }
  if (item?.can_download === false) {
    return ''
  }
  const url = mediaUrl(item?.url)
  if (!url) return ''
  if (/^https?:\/\//.test(url)) return url
  const joiner = url.includes('?') ? '&' : '?'
  return `${url}${joiner}download=1`
}

export function canDownloadMedia(item) {
  if (!item) return false
  if (typeof item.can_download === 'boolean') return item.can_download
  const vis = (item.download_visibility || 'public').toLowerCase()
  return vis === 'public' || vis === ''
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
  // Prefer direct-to-R2 draft (presign → PUT → complete). Falls back to cluster multipart.
  try {
    return await uploadMediaDirect(dayId, file, meta, onProgress)
  } catch (err) {
    const status = err?.response?.status
    const fallback = err?.response?.data?.fallback === 'multipart' || status === 503 || status === 404
    if (!fallback && err?.directFailed) {
      // Presign worked but R2 PUT failed (often CORS) — try multipart once.
      return uploadMediaMultipart(dayId, file, meta, onProgress)
    }
    if (fallback || !err?.response) {
      return uploadMediaMultipart(dayId, file, meta, onProgress)
    }
    throw err
  }
}

async function uploadMediaDirect(dayId, file, meta = {}, onProgress) {
  const filename = file.name || 'upload.bin'
  const contentType = file.type || 'application/octet-stream'
  const body = {
    filename,
    content_type: contentType,
    size: file.size
  }
  if (typeof meta === 'string') {
    if (meta) body.caption = meta
  } else {
    if (meta.caption) body.caption = meta.caption
    if (meta.kind) body.kind = meta.kind
  }

  let presign
  try {
    ;({ data: presign } = await client.post(`/days/${dayId}/media/presign`, body, { timeout: 30000 }))
  } catch (err) {
    err.directFailed = false
    throw err
  }

  const putURL = presign.upload_url
  const headers = { ...(presign.headers || {}) }
  if (!headers['Content-Type'] && !headers['content-type']) {
    headers['Content-Type'] = contentType
  }

  try {
    await axios.put(putURL, file, {
      headers,
      timeout: 0, // large videos — no axios timeout on CDN put
      onUploadProgress: (event) => {
        if (!onProgress || !event.total) return
        // Reserve last 5% for complete handshake
        onProgress(Math.min(95, Math.round((event.loaded / event.total) * 95)))
      },
      // Do not send surfing cookies to R2
      withCredentials: false,
      // Avoid transforming File/Blob
      transformRequest: [(d) => d]
    })
  } catch (err) {
    err.directFailed = true
    throw err
  }

  const mediaId = presign.media?.id
  if (!mediaId) {
    const e = new Error('presign missing media id')
    e.directFailed = true
    throw e
  }

  // Brief settle for R2 listing consistency, then confirm.
  let lastErr
  for (let attempt = 0; attempt < 4; attempt++) {
    try {
      if (attempt) await new Promise((r) => setTimeout(r, 400 * attempt))
      const { data } = await client.post(`/days/${dayId}/media/${mediaId}/complete`, null, { timeout: 30000 })
      if (onProgress) onProgress(100)
      return data.media || data
    } catch (err) {
      lastErr = err
      if (err?.response?.status !== 409) break
    }
  }
  lastErr.directFailed = true
  throw lastErr
}

async function uploadMediaMultipart(dayId, file, meta = {}, onProgress) {
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
    timeout: 300000,
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
