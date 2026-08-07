import axios from 'axios'

const SURFING_HOST = import.meta.env.VITE_SURFING_API_HOST || '/api/surfing'

const authClient = axios.create({
  baseURL: SURFING_HOST,
  timeout: 30000,
  withCredentials: true
})

export async function fetchAuthConfig() {
  const { data } = await authClient.get('/auth/config')
  return data
}

export async function fetchAuthMe() {
  const { data } = await authClient.get('/auth/me')
  return data
}

export function authLoginUrl() {
  return `${SURFING_HOST}/auth/login`
}

export function authLogoutUrl() {
  return `${SURFING_HOST}/auth/logout`
}

export async function postActivity(payload) {
  const { data } = await authClient.post('/activity', payload)
  return data
}

export async function listActivity({ limit = 200 } = {}) {
  const { data } = await authClient.get('/activity', { params: { limit } })
  return data
}
