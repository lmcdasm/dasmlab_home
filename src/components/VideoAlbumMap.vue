<template>
  <div class="video-map">
    <div class="video-map__toolbar row items-center q-mb-sm">
      <div class="col text-caption text-grey-7">
        {{ modeHint }}
      </div>
      <div class="row no-wrap q-gutter-xs">
        <q-btn
          dense
          unelevated
          :outline="mode !== 'grid'"
          :color="mode === 'grid' ? 'primary' : 'grey-7'"
          icon="grid_view"
          label="Cards"
          @click="mode = 'grid'"
        >
          <q-tooltip anchor="top middle" self="bottom middle" :delay="250">
            Cards shows large thumbnails
          </q-tooltip>
        </q-btn>
        <q-btn
          dense
          unelevated
          :outline="mode !== 'cabinet'"
          :color="mode === 'cabinet' ? 'primary' : 'grey-7'"
          icon="folder_open"
          label="Cabinet"
          @click="mode = 'cabinet'"
        >
          <q-tooltip anchor="top middle" self="bottom middle" :delay="250">
            This shows a smaller player view
          </q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Cards = default browser grid (slot from parent) -->
    <div v-if="mode === 'grid'" class="video-map__grid-slot">
      <slot name="grid" />
    </div>

    <!-- Cabinet = chronological filing list + tight panel (Immich-style timeline, not radial hub) -->
    <div v-else class="video-map__layout">
      <div class="video-map__cabinet">
        <div class="cabinet-head row items-center q-mb-sm">
          <q-icon name="schedule" size="18px" color="primary" class="q-mr-xs" />
          <span class="text-caption text-weight-medium">By capture / upload time</span>
          <q-space />
          <span class="text-caption text-grey-6">{{ items.length }} clips</span>
        </div>

        <div v-if="!drawers.length" class="cabinet-empty text-caption text-grey-6">
          No videos in this album yet.
        </div>

        <section
          v-for="drawer in drawers"
          :key="drawer.key"
          class="cabinet-drawer"
        >
          <button type="button" class="drawer-tab" @click="toggleDrawer(drawer.key)">
            <q-icon :name="collapsed[drawer.key] ? 'chevron_right' : 'expand_more'" size="18px" />
            <span class="drawer-tab__date">{{ drawer.label }}</span>
            <span class="drawer-tab__count">{{ drawer.items.length }}</span>
          </button>

          <ul v-show="!collapsed[drawer.key]" class="drawer-list">
            <li
              v-for="row in drawer.items"
              :key="row.item.id"
              class="drawer-row"
              :class="{ active: selectedId === row.item.id }"
              @click="select(row.item.id)"
              @dblclick="openBrowser(row.item.id)"
            >
              <div class="drawer-row__thumb">
                <video
                  :src="mediaUrl(row.item.url)"
                  muted
                  playsinline
                  preload="metadata"
                  class="drawer-row__video"
                />
              </div>
              <div class="drawer-row__body">
                <div class="drawer-row__title">{{ row.item.caption || row.item.filename }}</div>
                <div class="drawer-row__meta">
                  <span v-if="row.timeLabel !== '—'">{{ row.timeLabel }}</span>
                  <span v-if="row.durationLabel">· {{ row.durationLabel }}</span>
                  <span v-if="row.item.play_count">· {{ row.item.play_count }} plays</span>
                  <span>· {{ row.item.published ? 'CDN' : 'Draft' }}</span>
                </div>
              </div>
              <q-icon
                v-if="selectedId === row.item.id"
                name="chevron_right"
                color="primary"
                size="20px"
              />
            </li>
          </ul>
        </section>
      </div>

      <aside class="video-map__panel" :class="{ empty: !selected }">
        <template v-if="selected">
          <div class="panel-head row items-start no-wrap">
            <div class="col">
              <div class="text-subtitle1">{{ selected.caption || selected.filename }}</div>
              <div class="text-caption text-grey-7">Clip in {{ albumTitle }}</div>
            </div>
            <q-btn flat dense round icon="close" @click="selectedId = null" />
          </div>

          <!-- Small preview box — click plays muted in viewer -->
          <button type="button" class="panel-preview" @click="emitPlay(selected)">
            <video
              :src="mediaUrl(selected.url)"
              muted
              playsinline
              preload="metadata"
              class="panel-preview__video"
            />
            <span class="panel-preview__play">
              <q-icon name="play_arrow" size="28px" />
            </span>
            <span class="panel-preview__hint">Play muted</span>
          </button>

          <div class="stat-row">
            <div class="stat"><span>Plays</span><strong>{{ selected.play_count || 0 }}</strong></div>
            <div class="stat"><span>Length</span><strong>{{ selectedDuration || '—' }}</strong></div>
            <div class="stat"><span>Tags</span><strong>{{ approvedTags.length }}</strong></div>
            <div class="stat">
              <span>CDN</span>
              <strong>{{ selected.published ? 'Yes' : 'Draft' }}</strong>
            </div>
          </div>

          <div class="panel-actions q-gutter-xs q-mb-md">
            <q-btn unelevated color="primary" dense icon="play_arrow" label="Play muted" @click="emitPlay(selected)" />
            <q-btn outline color="primary" dense icon="ios_share" label="Share" @click="$emit('share', selected)" />
            <q-btn
              v-if="canDownloadSelected"
              outline
              color="primary"
              dense
              icon="download"
              label="Download"
              tag="a"
              :href="downloadHref"
              :download="selected.filename"
              rel="noopener"
            />
            <q-btn
              v-else
              outline
              color="primary"
              dense
              icon="download"
              label="Download"
              disable
            >
              <q-tooltip>Restricted ({{ selected.download_visibility || 'private' }})</q-tooltip>
            </q-btn>
            <q-btn
              outline
              color="primary"
              dense
              icon="cloud_upload"
              label="Publish"
              :loading="publishing"
              :disable="!!selected.published"
              @click="emitPublish"
            />
            <q-btn flat dense icon="open_in_new" @click="$emit('cdn', selected)">
              <q-tooltip>Open on CDN</q-tooltip>
            </q-btn>
            <q-btn flat dense icon="grid_view" @click="openBrowser(selected.id)">
              <q-tooltip>Card browser</q-tooltip>
            </q-btn>
          </div>

          <div class="flow-block">
            <div class="flow-label">Approved names</div>
            <div v-if="approvedTags.length" class="chip-list">
              <span v-for="t in approvedTags" :key="t.id" class="name-chip approved">{{ t.name }}</span>
            </div>
            <div v-else class="text-caption text-grey-6">No approved name tags yet.</div>
          </div>

          <div class="flow-block">
            <div class="flow-label pending">Pending approval</div>
            <div v-if="pendingTags.length" class="chip-list">
              <div v-for="t in pendingTags" :key="t.id" class="pending-row">
                <span class="name-chip pending">{{ t.name }}</span>
                <q-btn flat dense size="sm" color="positive" icon="check" @click="$emit('approve-tag', selected, t)" />
                <q-btn flat dense size="sm" color="negative" icon="close" @click="$emit('reject-tag', selected, t)" />
              </div>
            </div>
            <div v-else class="text-caption text-grey-6">Nothing waiting.</div>
          </div>

          <div class="flow-block">
            <div class="flow-label">Propose a name</div>
            <div class="text-caption text-grey-7 q-mb-xs">
              Plain name only — no links. Owner must approve. Not a social network.
            </div>
            <div class="row q-gutter-sm items-center">
              <q-input
                v-model="tagDraft"
                dense
                filled
                class="col"
                placeholder="e.g. Alex"
                maxlength="40"
                @keyup.enter="submitTag"
              />
              <q-btn color="primary" dense label="Tag" :loading="tagging" @click="submitTag" />
            </div>
          </div>

          <p v-if="selected.notes" class="panel-notes">{{ selected.notes }}</p>
          <p v-if="selectedTimeLabel" class="panel-notes panel-notes--meta">{{ selectedTimeLabel }}</p>
        </template>
        <template v-else>
          <div class="text-subtitle2">Album videos</div>
          <div class="text-caption text-grey-7 q-mb-md">
            Pick a clip from the cabinet — preview and actions stay here.
          </div>
          <div class="stat-row">
            <div class="stat"><span>Clips</span><strong>{{ items.length }}</strong></div>
            <div class="stat"><span>Plays</span><strong>{{ totalPlays }}</strong></div>
            <div class="stat"><span>On CDN</span><strong>{{ publishedCount }}</strong></div>
            <div class="stat"><span>Tags</span><strong>{{ totalApprovedTags }}</strong></div>
          </div>
        </template>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { canDownloadMedia, mediaDownloadUrl, mediaUrl } from 'src/services/surfingApi'

const props = defineProps({
  items: { type: Array, default: () => [] },
  albumTitle: { type: String, default: 'Session' },
  dayId: { type: String, default: '' },
  sortMode: { type: String, default: 'longest' },
  publishing: { type: Boolean, default: false },
  tagging: { type: Boolean, default: false }
})

const emit = defineEmits(['play', 'share', 'cdn', 'publish', 'propose-tag', 'approve-tag', 'reject-tag', 'browse'])

/** Cabinet first — compact timeline + side player; Cards for large thumbnails. */
const mode = ref('cabinet')
const selectedId = ref(null)
const tagDraft = ref('')
const collapsed = reactive({})

const selected = computed(() => props.items.find((i) => i.id === selectedId.value) || null)

const canDownloadSelected = computed(() => canDownloadMedia(selected.value))
const downloadHref = computed(() => mediaDownloadUrl(selected.value, props.dayId))
const selectedDuration = computed(() => {
  const s = Math.round(Number(selected.value?.duration_sec) || 0)
  if (s <= 0) return ''
  const m = Math.floor(s / 60)
  const r = s % 60
  return `${m}:${String(r).padStart(2, '0')}`
})

const approvedTags = computed(() => (selected.value?.tags || []).filter((t) => t.status === 'approved'))
const pendingTags = computed(() => (selected.value?.tags || []).filter((t) => t.status === 'pending'))

const totalPlays = computed(() => props.items.reduce((n, i) => n + (i.play_count || 0), 0))
const publishedCount = computed(() => props.items.filter((i) => i.published).length)
const totalApprovedTags = computed(() =>
  props.items.reduce((n, i) => n + (i.tags || []).filter((t) => t.status === 'approved').length, 0)
)

const modeHint = computed(() =>
  mode.value === 'grid'
    ? 'Cards — large thumbnails. Open a clip for the full viewer.'
    : 'Cabinet — smaller player view · day drawers by time · double-click a row → Cards'
)

const selectedTimeLabel = computed(() => {
  if (!selected.value) return ''
  const ts = resolveTimestamp(selected.value)
  if (!ts) return ''
  return formatLong(ts)
})

/** Group clips into day drawers (newest day first). Uses created_at; filename date as hint. */
const drawers = computed(() => {
  const buckets = new Map()
  for (const item of props.items || []) {
    const ts = resolveTimestamp(item)
    const key = ts ? dateKey(ts) : 'undated'
    const label = ts ? formatDay(ts) : 'Undated'
    if (!buckets.has(key)) {
      buckets.set(key, { key, label, sort: ts ? ts.getTime() : 0, items: [] })
    }
    buckets.get(key).items.push({
      item,
      ts,
      timeLabel: ts ? formatTime(ts) : '—',
      durationLabel: formatDur(item.duration_sec)
    })
  }
  const out = [...buckets.values()].sort((a, b) => b.sort - a.sort)
  for (const d of out) {
    d.items.sort((a, b) => compareClips(a, b, props.sortMode))
  }
  return out
})

function formatDur(sec) {
  const s = Math.round(Number(sec) || 0)
  if (s <= 0) return ''
  const m = Math.floor(s / 60)
  const r = s % 60
  return `${m}:${String(r).padStart(2, '0')}`
}

function compareClips(a, b, mode) {
  const da = Number(a.item?.duration_sec) || 0
  const db = Number(b.item?.duration_sec) || 0
  switch (mode) {
    case 'shortest':
      return da - db || (b.ts?.getTime() || 0) - (a.ts?.getTime() || 0)
    case 'newest':
      return (b.ts?.getTime() || 0) - (a.ts?.getTime() || 0)
    case 'name':
      return String(a.item?.caption || a.item?.filename || '').localeCompare(
        String(b.item?.caption || b.item?.filename || ''),
        undefined,
        { sensitivity: 'base' }
      )
    case 'plays':
      return (b.item?.play_count || 0) - (a.item?.play_count || 0) || db - da
    case 'longest':
    default:
      if (da > 0 && db > 0) return db - da
      if (da > 0) return -1
      if (db > 0) return 1
      return (b.ts?.getTime() || 0) - (a.ts?.getTime() || 0)
  }
}
watch(
  () => props.items,
  (list) => {
    if (selectedId.value && !list.some((i) => i.id === selectedId.value)) {
      selectedId.value = null
    }
  }
)

watch(
  drawers,
  (list) => {
    // Keep first drawer open; collapse none by default for small albums
    for (const d of list) {
      if (collapsed[d.key] === undefined) {
        collapsed[d.key] = false
      }
    }
  },
  { immediate: true }
)

function resolveTimestamp(item) {
  // Camera filenames often keep capture day better than bulk-upload created_at.
  const fromName = parseCameraFilename(item)
  if (fromName) return fromName
  if (item?.created_at) {
    const d = new Date(item.created_at)
    if (!Number.isNaN(d.getTime())) return d
  }
  return null
}

/** Olympus/Panasonic-ish: P7070149 → Jul 7, P7080025 → Jul 8 (month + day + seq). */
function parseCameraFilename(item) {
  const name = (item?.filename || item?.caption || '').replace(/\.[^.]+$/, '')
  const m = name.match(/^P(1[0-2]|[1-9])(\d{2})\d+/i)
  if (!m) return null
  const month = Number(m[1])
  const day = Number(m[2])
  if (month < 1 || month > 12 || day < 1 || day > 31) return null
  const yearHint = item?.created_at ? new Date(item.created_at).getFullYear() : new Date().getFullYear()
  const year = Number.isFinite(yearHint) ? yearHint : new Date().getFullYear()
  return new Date(year, month - 1, day, 12, 0, 0)
}

function dateKey(d) {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function formatDay(d) {
  return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
}

function formatTime(d) {
  return d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
}

function formatLong(d) {
  return d.toLocaleString(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function toggleDrawer(key) {
  collapsed[key] = !collapsed[key]
}

function select(id) {
  selectedId.value = id
}

function openBrowser(id) {
  if (id) selectedId.value = id
  mode.value = 'grid'
  emit('browse', id || selectedId.value)
}

function emitPlay(item) {
  emit('play', item)
}

function emitPublish() {
  emit('publish', selected.value)
}

function submitTag() {
  const name = tagDraft.value.trim()
  if (!name || !selected.value) return
  emit('propose-tag', selected.value, name)
  tagDraft.value = ''
}

defineExpose({ mode, selectedId, select, openBrowser })
</script>

<style scoped>
.video-map__layout {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(280px, 0.95fr);
  gap: 0.85rem;
  min-height: 420px;
}

.video-map__cabinet {
  border: 1px solid rgba(15, 143, 124, 0.18);
  border-radius: 14px;
  background: linear-gradient(165deg, #f7fbfa, #eef6f4);
  padding: 0.65rem 0.7rem 0.85rem;
  max-height: 520px;
  overflow: auto;
}

.cabinet-head {
  position: sticky;
  top: 0;
  z-index: 1;
  background: linear-gradient(180deg, #f7fbfa 70%, transparent);
  padding-bottom: 0.25rem;
}

.cabinet-drawer + .cabinet-drawer {
  margin-top: 0.45rem;
}

.drawer-tab {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  width: 100%;
  border: none;
  background: rgba(6, 54, 66, 0.06);
  border-radius: 8px;
  padding: 0.4rem 0.55rem;
  cursor: pointer;
  text-align: left;
  color: #102833;
  font: inherit;
}

.drawer-tab:hover {
  background: rgba(15, 143, 124, 0.12);
}

.drawer-tab__date {
  font-size: 0.82rem;
  font-weight: 650;
  flex: 1;
}

.drawer-tab__count {
  font-size: 0.72rem;
  font-weight: 700;
  color: #0f8f7c;
  background: rgba(15, 143, 124, 0.12);
  border-radius: 999px;
  padding: 0.1rem 0.45rem;
}

.drawer-list {
  list-style: none;
  margin: 0.35rem 0 0;
  padding: 0;
}

.drawer-row {
  display: flex;
  align-items: center;
  gap: 0.55rem;
  padding: 0.4rem 0.45rem;
  border-radius: 10px;
  cursor: pointer;
  border: 1px solid transparent;
  transition: background 140ms ease, border-color 140ms ease;
}

.drawer-row:hover {
  background: rgba(255, 255, 255, 0.85);
}

.drawer-row.active {
  background: #fff;
  border-color: rgba(15, 143, 124, 0.45);
  box-shadow: 0 4px 12px rgba(6, 54, 66, 0.06);
}

.drawer-row__thumb {
  width: 56px;
  height: 40px;
  border-radius: 6px;
  overflow: hidden;
  background: #063642;
  flex-shrink: 0;
}

.drawer-row__video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.drawer-row__body {
  min-width: 0;
  flex: 1;
}

.drawer-row__title {
  font-size: 0.84rem;
  font-weight: 650;
  color: #102833;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.drawer-row__meta {
  font-size: 0.72rem;
  color: #6a8090;
  margin-top: 0.1rem;
}

.video-map__panel {
  border: 1px solid rgba(15, 143, 124, 0.18);
  border-radius: 14px;
  background: #fff;
  padding: 0.9rem 1rem;
  box-shadow: 0 10px 24px rgba(6, 54, 66, 0.06);
  max-height: 520px;
  overflow: auto;
}

.video-map__panel.empty {
  background: linear-gradient(165deg, #fff, #f4faf8);
}

.panel-preview {
  position: relative;
  display: block;
  width: 100%;
  aspect-ratio: 16 / 10;
  max-height: 160px;
  margin: 0.65rem 0 0.25rem;
  border: none;
  padding: 0;
  border-radius: 12px;
  overflow: hidden;
  background: #063642;
  cursor: pointer;
}

.panel-preview__video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  opacity: 0.92;
}

.panel-preview__play {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  color: #fff;
  background: rgba(6, 54, 66, 0.28);
  transition: background 160ms ease;
}

.panel-preview:hover .panel-preview__play {
  background: rgba(6, 54, 66, 0.4);
}

.panel-preview__hint {
  position: absolute;
  left: 0.55rem;
  bottom: 0.45rem;
  font-size: 0.68rem;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 650;
}

.stat-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.45rem;
  margin: 0.75rem 0 1rem;
}

.stat {
  background: #f3f8f6;
  border-radius: 10px;
  padding: 0.45rem 0.55rem;
}

.stat span {
  display: block;
  font-size: 0.68rem;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #6a8090;
}

.stat strong {
  font-size: 1rem;
  color: #102833;
}

.flow-block {
  margin-bottom: 0.85rem;
}

.flow-label {
  font-size: 0.7rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #0f8f7c;
  font-weight: 700;
  margin-bottom: 0.35rem;
}

.flow-label.pending {
  color: #b07a12;
}

.chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
}

.name-chip {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 0.18rem 0.55rem;
  font-size: 0.78rem;
  font-weight: 600;
}

.name-chip.approved {
  background: rgba(15, 143, 124, 0.12);
  color: #0a5c52;
  border: 1px solid rgba(15, 143, 124, 0.28);
}

.name-chip.pending {
  background: rgba(240, 180, 60, 0.15);
  color: #8a5a00;
  border: 1px solid rgba(180, 120, 20, 0.28);
}

.pending-row {
  display: flex;
  align-items: center;
  gap: 0.15rem;
}

.panel-notes {
  margin: 0.5rem 0 0;
  font-size: 0.85rem;
  color: #5a7080;
  line-height: 1.45;
}

.panel-notes--meta {
  font-size: 0.75rem;
  color: #8a9aaa;
}

@media (max-width: 900px) {
  .video-map__layout {
    grid-template-columns: 1fr;
  }

  .video-map__cabinet,
  .video-map__panel {
    max-height: none;
  }
}
</style>
