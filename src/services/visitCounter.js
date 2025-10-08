import axios from 'axios'

const NAMESPACE = 'dasmlab_home'
const KEY = 'unique_visits'
const STORAGE_KEY = 'dasmlab_unique_visit_recorded'

// Multiple public base URLs for resilience (DNS or routing issues on some networks)
const BASE_URLS = [
	'https://api.countapi.xyz',
	'https://countapi.xyz'
]

const REQUEST_OPTS = { timeout: 5000 }

function encodedPath(parts) {
	const [ns, key] = [encodeURIComponent(NAMESPACE), encodeURIComponent(KEY)]
	return parts.replace(':ns', ns).replace(':key', key)
}

async function tryGetAcrossBases(path) {
	for (const base of BASE_URLS) {
		try {
			const { data } = await axios.get(`${base}${path}`, REQUEST_OPTS)
			return data
		} catch {
			// try next base
		}
	}
	throw new Error('All endpoints failed')
}

async function ensureCounterExists() {
	const createPath = encodedPath('/create?namespace=:ns&key=:key&value=0')
	for (const base of BASE_URLS) {
		try {
			await axios.get(`${base}${createPath}`, REQUEST_OPTS)
			return
		} catch {
			// ignore and try next base; /hit usually autocreates, this is just proactive
		}
	}
}

export async function incrementIfFirstVisit() {
	try {
		if (typeof window === 'undefined') return
		if (localStorage.getItem(STORAGE_KEY)) return

		await ensureCounterExists()
		await tryGetAcrossBases(encodedPath('/hit/:ns/:key'))
		localStorage.setItem(STORAGE_KEY, '1')
	} catch {
		// ignore network errors and keep UI stable
	}
}

export async function getVisitCount() {
	try {
		const data = await tryGetAcrossBases(encodedPath('/get/:ns/:key'))
		return data?.value ?? 0
	} catch {
		return 0
	}
}


