const STORAGE_KEY = 'dasmlab_unique_visit_recorded'
const DEBUG = (import.meta.env.COUNTER_DEBUG || process.env.COUNTER_DEBUG || 'true') === 'true'
const API_BASE = (import.meta.env.COUNTER_API_BASE || process.env.COUNTER_API_BASE || 'https://api.counterapi.dev').replace(/\/$/, '')

// Supported configurations:
// 1) COUNTER_FULL_PATH: '/v2/<workspace>/<slug>' (preferred)
//    Also supports legacy '/v2/workspaces/<workspace>/<slug>'
// 2) COUNTER_WORKSPACE + COUNTER_COUNTER_SLUG
// 3) Fallback slug 'dasmlab-home-views' under provided workspace

function resolveWorkspaceAndSlug() {
	const fullPath = (import.meta.env.COUNTER_FULL_PATH || process.env.COUNTER_FULL_PATH || '').trim()
	if (fullPath && typeof fullPath === 'string') {
		// Extract workspace + slug from:
		// - '/v2/{workspace}/{slug}'
		// - '/v2/workspaces/{workspace}/{slug}' (legacy)
		// - full URLs ending in either of the above
		const matches = fullPath.match(/\/v2\/(?:workspaces\/)?([^/]+)\/([^/?#]+)/)
		if (matches && matches.length === 3) {
			return { workspace: matches[1], slug: matches[2] }
		}
	}

	const workspace = import.meta.env.COUNTER_WORKSPACE || process.env.COUNTER_WORKSPACE || 'dasmlab-home'
	const slug = import.meta.env.COUNTER_COUNTER_SLUG || process.env.COUNTER_COUNTER_SLUG || 'dasmlab-home-views'
	return { workspace, slug }
}

const { workspace: WORKSPACE, slug: COUNTER_SLUG } = resolveWorkspaceAndSlug()
const ACCESS_TOKEN = import.meta.env.COUNTER_API_TOKEN || process.env.COUNTER_API_TOKEN || undefined

function logDebug(message, payload = undefined) {
	if (!DEBUG) return
	if (payload === undefined) {
		console.log(`[counter] ${message}`)
		return
	}
	console.log(`[counter] ${message}`, payload)
}

function buildCounterUrl(action = '') {
	const suffix = action ? `/${action}` : ''
	return `${API_BASE}/v2/${WORKSPACE}/${COUNTER_SLUG}${suffix}`
}

async function callCounterApi(action = '') {
	const isBrowser = typeof window !== 'undefined'
	const url = `${buildCounterUrl(action)}?t=${Date.now()}`
	const headers = { Accept: 'application/json' }
	// Counter API preflight does not allow Authorization header in browsers.
	if (!isBrowser && ACCESS_TOKEN) {
		headers.Authorization = `Bearer ${ACCESS_TOKEN}`
	}

	const res = await fetch(url, {
		method: 'GET',
		headers,
		cache: 'no-store'
	})

	if (!res.ok) {
		throw new Error(`counter api ${action || 'get'} failed: ${res.status}`)
	}
	return res.json()
}

function extractCount(payload) {
	if (typeof payload?.value === 'number') return payload.value
	if (typeof payload?.count === 'number') return payload.count
	if (typeof payload?.up_count === 'number') return payload.up_count

	const nested = payload?.data
	if (!nested || typeof nested !== 'object') return 0
	if (typeof nested.value === 'number') return nested.value
	if (typeof nested.count === 'number') return nested.count
	if (typeof nested.up_count === 'number') return nested.up_count
	if (typeof nested.upCount === 'number') return nested.upCount
	return 0
}

export async function incrementIfFirstVisit() {
	try {
		if (typeof window === 'undefined') return
		if (localStorage.getItem(STORAGE_KEY)) return
		const payload = await callCounterApi('up')
		logDebug('up', { value: extractCount(payload) })
		localStorage.setItem(STORAGE_KEY, '1')
	} catch (error) {
		logDebug('up failed', { message: error?.message || 'unknown' })
		return
	}
}

export async function getVisitCount() {
	try {
		const payload = await callCounterApi()
		const value = extractCount(payload)
		logDebug('get resolved', { value })
		return value
	} catch (error) {
		logDebug('get failed', { message: error?.message || 'unknown' })
		return 0
	}
}


