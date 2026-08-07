<template>
  <q-page class="activity-page" padding>
    <header class="page-head">
      <div>
        <p class="page-kicker">Internal · first-party CDP</p>
        <h1 class="page-title">Activity</h1>
        <p class="page-lead">
          Public + authenticated engagement (newest first). Dual-gated to
          <code>dasm</code>. Anonymous visitors carry <code>anonymousId</code> /
          <code>sessionId</code>; login stitches via identify/alias.
        </p>
        <p v-if="events.length" class="page-meta">
          {{ events.length }} recent event<span v-if="events.length !== 1">s</span>
          · {{ anonCount }} anon · {{ knownCount }} known
        </p>
      </div>
      <div class="head-actions">
        <q-select
          v-model="whoFilter"
          dense
          outlined
          emit-value
          map-options
          :options="whoOptions"
          label="Who"
          style="min-width: 120px"
        />
        <q-select
          v-model="typeFilter"
          dense
          outlined
          emit-value
          map-options
          :options="typeOptions"
          label="Type"
          style="min-width: 140px"
        />
        <q-btn color="primary" outline icon="refresh" label="Refresh" :loading="loading" @click="load" />
      </div>
    </header>

    <div v-if="loading && !events.length" class="row justify-center q-my-xl">
      <q-spinner size="3em" color="primary" />
    </div>

    <div v-else-if="!filtered.length" class="empty-state">
      <q-icon name="history" size="56px" color="blue-grey-5" />
      <div class="empty-title">No events yet</div>
      <div class="empty-sub">Browse the public site — page views will appear here.</div>
    </div>

    <ul v-else class="event-list">
      <li v-for="(ev, i) in filtered" :key="`${ev.ts}-${i}`" class="event-row">
        <span class="event-ts">{{ formatTs(ev.ts) }}</span>
        <span class="event-body">
          {{ formatLine(ev) }}
          <span v-if="contextBits(ev).length" class="event-chips">
            <span v-for="bit in contextBits(ev)" :key="bit" class="chip">{{ bit }}</span>
          </span>
        </span>
      </li>
    </ul>
  </q-page>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { listActivity } from 'src/services/authApi'

const events = ref([])
const loading = ref(false)
const typeFilter = ref('all')
const whoFilter = ref('all')
const typeOptions = [
  { label: 'All types', value: 'all' },
  { label: 'Page', value: 'page' },
  { label: 'Navigate', value: 'navigate' },
  { label: 'Login', value: 'login' },
  { label: 'Identify', value: 'identify' },
  { label: 'Alias', value: 'alias' },
  { label: 'Track', value: 'track' },
  { label: 'Engaged', value: 'engaged' }
]
const whoOptions = [
  { label: 'Everyone', value: 'all' },
  { label: 'Anonymous', value: 'anon' },
  { label: 'Known', value: 'known' }
]

let poll = null

const filtered = computed(() => {
  let list = events.value
  if (typeFilter.value !== 'all') {
    list = list.filter((e) => e.type === typeFilter.value || (typeFilter.value === 'page' && e.type === 'navigate'))
  }
  if (whoFilter.value === 'anon') {
    list = list.filter((e) => !e.user && !e.sub)
  } else if (whoFilter.value === 'known') {
    list = list.filter((e) => e.user || e.sub)
  }
  return list
})

const anonCount = computed(() => events.value.filter((e) => !e.user && !e.sub).length)
const knownCount = computed(() => events.value.filter((e) => e.user || e.sub).length)

function formatTs(ts) {
  if (!ts) return ''
  try {
    return new Date(ts).toLocaleString()
  } catch {
    return String(ts)
  }
}

function formatDuration(ms) {
  if (ms == null || ms <= 0) return null
  if (ms < 1000) return `${ms}ms`
  const s = ms / 1000
  if (s < 60) return `${s.toFixed(1)}s`
  const m = Math.floor(s / 60)
  const rem = Math.round(s % 60)
  return `${m}m ${rem}s`
}

function shortId(id) {
  if (!id) return ''
  return id.length > 10 ? id.slice(0, 8) + '…' : id
}

function whoLabel(ev) {
  if (ev.user) return ev.user
  if (ev.email) return ev.email
  if (ev.anonymousId) return `anon:${shortId(ev.anonymousId)}`
  return 'unknown'
}

function formatLine(ev) {
  const who = whoLabel(ev)
  if (ev.type === 'login' || ev.type === 'identify') {
    return `${who} ${ev.type === 'login' ? 'logged in' : 'identified'}`
  }
  if (ev.type === 'alias') {
    return `${who} aliased from ${shortId(ev.previousId || ev.anonymousId)}`
  }
  if (ev.type === 'track') {
    return `${who} tracked ${ev.event || 'event'}${ev.path ? ` @ ${ev.path}` : ''}`
  }
  const path = ev.path || '/'
  const parts = [`${who} viewed ${path}`]
  const dwell = formatDuration(ev.dwellMs)
  const visible = formatDuration(ev.visibleMs)
  const engaged = formatDuration(ev.engagedMs)
  const metrics = []
  if (dwell) metrics.push(`dwell ${dwell}`)
  if (visible) metrics.push(`visible ${visible}`)
  if (engaged) metrics.push(`engaged ${engaged}`)
  if (ev.scrollMaxPct > 0) metrics.push(`scroll ${ev.scrollMaxPct}%`)
  if (metrics.length) parts.push(`(${metrics.join(', ')})`)
  return parts.join(' ')
}

function contextBits(ev) {
  const bits = []
  if (ev.country) bits.push(ev.country)
  if (ev.locale) bits.push(ev.locale)
  if (ev.utmSource) bits.push(`utm:${ev.utmSource}`)
  if (ev.bot) bits.push('bot')
  if (ev.sessionId) bits.push(`sid:${shortId(ev.sessionId)}`)
  return bits
}

async function load() {
  loading.value = true
  try {
    const data = await listActivity({ limit: 300 })
    events.value = data.events || []
  } catch {
    events.value = []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await load()
  poll = setInterval(load, 15_000)
})

onUnmounted(() => {
  if (poll) clearInterval(poll)
})
</script>

<style scoped>
.activity-page {
  max-width: 960px;
  margin: 0 auto;
}
.page-kicker {
  margin: 0;
  font-size: 0.75rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #78909c;
  font-weight: 700;
}
.page-title {
  margin: 0.2rem 0 0.35rem;
  font-size: clamp(1.6rem, 3vw, 2rem);
  color: #1a2b36;
}
.page-lead {
  margin: 0;
  color: #546e7a;
  line-height: 1.45;
  max-width: 42rem;
}
.page-meta {
  margin: 0.45rem 0 0;
  font-size: 0.85rem;
  color: #90a4ae;
}
.page-head {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.25rem;
}
.head-actions {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  flex-wrap: wrap;
}
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #78909c;
}
.empty-title {
  margin-top: 0.75rem;
  font-size: 1.1rem;
  font-weight: 650;
  color: #455a64;
}
.empty-sub {
  margin-top: 0.25rem;
  font-size: 0.9rem;
}
.event-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.event-row {
  display: grid;
  grid-template-columns: 11rem 1fr;
  gap: 1rem;
  padding: 0.65rem 0;
  border-bottom: 1px solid #eceff1;
  font-size: 0.95rem;
}
.event-ts {
  color: #78909c;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}
.event-body {
  color: #263238;
  word-break: break-word;
}
.event-chips {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 0.25rem;
  margin-left: 0.35rem;
}
.chip {
  display: inline-block;
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.02em;
  color: #546e7a;
  background: #eceff1;
  padding: 0.1rem 0.35rem;
  border-radius: 3px;
}
@media (max-width: 600px) {
  .event-row {
    grid-template-columns: 1fr;
    gap: 0.2rem;
  }
}
</style>
