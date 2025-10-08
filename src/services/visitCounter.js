import axios from 'axios'

const NAMESPACE = 'dasmlab_home'
const KEY = 'unique_visits'
const STORAGE_KEY = 'dasmlab_unique_visit_recorded'

export async function incrementIfFirstVisit() {
	try {
		if (localStorage.getItem(STORAGE_KEY)) return
		await axios.get(`https://api.countapi.xyz/hit/${NAMESPACE}/${KEY}`)
		localStorage.setItem(STORAGE_KEY, '1')
	} catch {
		// ignore network errors and keep UI stable
	}
}

export async function getVisitCount() {
	try {
		const { data } = await axios.get(`https://api.countapi.xyz/get/${NAMESPACE}/${KEY}`)
		return data?.value ?? 0
	} catch {
		return 0
	}
}


