import { Counter } from 'counterapi'

const STORAGE_KEY = 'dasmlab_unique_visit_recorded'
const DEBUG = (import.meta.env.COUNTER_DEBUG || process.env.COUNTER_DEBUG || 'true') === 'true'

// Supported configurations:
// 1) COUNTER_FULL_PATH: '/v2/workspaces/<workspace>/<slug>' (preferred)
// 2) COUNTER_WORKSPACE + COUNTER_COUNTER_SLUG
// 3) Fallback slug 'unique_visits' under provided workspace

function resolveWorkspaceAndSlug() {
	const fullPath = import.meta.env.COUNTER_FULL_PATH || process.env.COUNTER_FULL_PATH
	if (fullPath && typeof fullPath === 'string') {
		// Extract workspace and slug from '/v2/workspaces/{ws}/{slug}'
		const matches = fullPath.match(/\/v2\/workspaces\/([^/]+)\/([^/]+)/)
		if (matches && matches.length === 3) {
			return { workspace: matches[1], slug: matches[2] }
		}
	}

	const workspace = import.meta.env.COUNTER_WORKSPACE || process.env.COUNTER_WORKSPACE || 'dasmlab-home'
	const slug = import.meta.env.COUNTER_COUNTER_SLUG || process.env.COUNTER_COUNTER_SLUG || 'unique_visits'
	return { workspace, slug }
}

const { workspace: WORKSPACE, slug: COUNTER_SLUG } = resolveWorkspaceAndSlug()
const ACCESS_TOKEN = import.meta.env.COUNTER_API_TOKEN || process.env.COUNTER_API_TOKEN || undefined

let client
function getClient() {
	if (!client) {
		const isBrowser = typeof window !== 'undefined'
		const options = { workspace: WORKSPACE, timeout: 8000 }
		// Avoid Authorization header in browser to prevent CORS preflight failures
		if (!isBrowser && ACCESS_TOKEN) {
			options.accessToken = ACCESS_TOKEN
		}
		client = new Counter(options)
		console.log('[counter] init', { workspace: WORKSPACE, slug: COUNTER_SLUG, withToken: Boolean(options.accessToken), debug: DEBUG })
	}
	return client
}

export async function incrementIfFirstVisit() {
	try {
		if (typeof window === 'undefined') return
		if (localStorage.getItem(STORAGE_KEY)) return
		const counter = getClient()
		const res = await counter.up(COUNTER_SLUG)
		console.log('[counter] up', { value: res?.value })
		localStorage.setItem(STORAGE_KEY, '1')
	} catch {
		return
	}
}

export async function getVisitCount() {
	try {
		const counter = getClient()
		const res = await counter.get(COUNTER_SLUG)
		const value = res?.value ?? 0
		console.log('[counter] get', { value })
		return value
	} catch {
		return 0
	}
}


