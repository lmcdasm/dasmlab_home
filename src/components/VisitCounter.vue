<template>
	<span>Visits: {{ count }}</span>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { incrementIfFirstVisit, getVisitCount } from 'src/services/visitCounter'

const count = ref(0)
const DEBUG = (import.meta.env.COUNTER_DEBUG || 'false') === 'true'

onMounted(async () => {
	await incrementIfFirstVisit()
	const value = await getVisitCount()
	if (DEBUG) console.log('[counter] mount display value', value)
	count.value = value
})
</script>


