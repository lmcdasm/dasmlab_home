import axios from 'axios'

const SURFING_HOST = import.meta.env.VITE_SURFING_API_HOST || '/api/surfing'

const client = axios.create({
  baseURL: SURFING_HOST,
  timeout: 120000
})

export function mediaUrl(path) {
  if (!path) return ''
  if (/^https?:\/\//.test(path)) return path
  return `${SURFING_HOST}${path}`
}

export function mediaDownloadUrl(item) {
  const url = mediaUrl(item?.url)
  if (!url) return ''
  const joiner = url.includes('?') ? '&' : '?'
  return `${url}${joiner}download=1`
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

export async function uploadMedia(dayId, file, caption = '', onProgress) {
  const form = new FormData()
  form.append('file', file)
  if (caption) form.append('caption', caption)

  const { data } = await client.post(`/days/${dayId}/media`, form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (event) => {
      if (!onProgress || !event.total) return
      onProgress(Math.round((event.loaded / event.total) * 100))
    }
  })
  return data
}

export async function deleteMedia(dayId, mediaId) {
  await client.delete(`/days/${dayId}/media/${mediaId}`)
}

function sortDays(days) {
  return [...days].sort((a, b) => {
    const dateCmp = (b.date || '').localeCompare(a.date || '')
    if (dateCmp !== 0) return dateCmp
    return (b.created_at || '').localeCompare(a.created_at || '')
  })
}
