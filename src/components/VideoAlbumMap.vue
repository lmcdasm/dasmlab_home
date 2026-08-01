<template>
  <div class="video-map">
    <div class="video-map__toolbar row items-center q-mb-sm">
      <div class="col text-caption text-grey-7">
        Click a clip for the tight panel · double-click to open the card browser
      </div>
      <q-btn-toggle
        v-model="mode"
        dense
        toggle-color="primary"
        :options="[
          { label: 'Map', value: 'map', icon: 'hub' },
          { label: 'Cards', value: 'grid', icon: 'grid_view' }
        ]"
      />
    </div>

    <div v-if="mode === 'map'" class="video-map__layout">
      <div class="video-map__stage">
        <svg class="video-map__svg" :viewBox="`0 0 ${vbW} ${vbH}`" preserveAspectRatio="xMidYMid meet">
          <defs>
            <pattern id="surfVidGrid" width="28" height="28" patternUnits="userSpaceOnUse">
              <path d="M 28 0 L 0 0 0 28" fill="none" stroke="rgba(6,54,66,0.07)" stroke-width="1" />
            </pattern>
            <filter id="surfVidShadow" x="-40%" y="-40%" width="180%" height="180%">
              <feDropShadow dx="0" dy="3" stdDeviation="5" flood-color="#063642" flood-opacity="0.18" />
            </filter>
            <linearGradient id="surfVidHub" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stop-color="#063642" />
              <stop offset="100%" stop-color="#0f8f7c" />
            </linearGradient>
          </defs>

          <rect width="100%" height="100%" fill="url(#surfVidGrid)" />

          <path
            v-for="n in nodes"
            :key="'link-' + n.item.id"
            :d="spoke(hub.x, hub.y, n.x, n.y)"
            class="video-map__link"
            :class="{ active: selectedId === n.item.id }"
            :stroke="n.color"
            fill="none"
          />

          <g class="video-map__hub" filter="url(#surfVidShadow)">
            <circle :cx="hub.x" :cy="hub.y" r="58" fill="url(#surfVidHub)" class="hub-pulse" />
            <circle :cx="hub.x" :cy="hub.y" r="50" fill="url(#surfVidHub)" />
            <text :x="hub.x" :y="hub.y - 6" text-anchor="middle" class="hub-title">Videos</text>
            <text :x="hub.x" :y="hub.y + 14" text-anchor="middle" class="hub-sub">{{ items.length }}</text>
          </g>

          <g
            v-for="n in nodes"
            :key="'node-' + n.item.id"
            class="video-map__node"
            :class="{ active: selectedId === n.item.id }"
            filter="url(#surfVidShadow)"
            @click="select(n.item.id)"
            @dblclick.stop="openBrowser(n.item.id)"
            tabindex="0"
            @keydown.enter="select(n.item.id)"
          >
            <rect
              :x="n.x - n.w / 2"
              :y="n.y - 18"
              :width="n.w"
              height="36"
              rx="9"
              class="node-rect"
              :stroke="n.color"
            />
            <text :x="n.x" :y="n.y + 5" text-anchor="middle" class="node-text">{{ n.label }}</text>
          </g>
        </svg>
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

          <div class="stat-row">
            <div class="stat"><span>Plays</span><strong>{{ selected.play_count || 0 }}</strong></div>
            <div class="stat"><span>Tags</span><strong>{{ approvedTags.length }}</strong></div>
            <div class="stat"><span>Pending</span><strong>{{ pendingTags.length }}</strong></div>
            <div class="stat">
              <span>CDN</span>
              <strong>{{ selected.published ? 'Yes' : 'Draft' }}</strong>
            </div>
          </div>

          <div class="panel-actions q-gutter-xs q-mb-md">
            <q-btn unelevated color="primary" dense icon="play_arrow" label="Play muted" @click="emitPlay(selected)" />
            <q-btn outline color="primary" dense icon="ios_share" label="Share" @click="$emit('share', selected)" />
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
        </template>
        <template v-else>
          <div class="text-subtitle2">Album videos</div>
          <div class="text-caption text-grey-7 q-mb-md">
            Select a satellite clip. Summary stays here — no popup chrome.
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

    <div v-else class="video-map__grid-slot">
      <slot name="grid" />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  albumTitle: { type: String, default: 'Session' },
  publishing: { type: Boolean, default: false },
  tagging: { type: Boolean, default: false }
})

const emit = defineEmits(['play', 'share', 'cdn', 'publish', 'propose-tag', 'approve-tag', 'reject-tag', 'browse'])

const mode = ref('map')
const selectedId = ref(null)
const tagDraft = ref('')
const vbW = 720
const vbH = 420
const hub = { x: 360, y: 210 }

const palette = ['#0f8f7c', '#1a7a6d', '#0aa3a0', '#5eb4c8', '#2f6f8f', '#063642', '#3d8b7a', '#148f9c']

const selected = computed(() => props.items.find((i) => i.id === selectedId.value) || null)

const approvedTags = computed(() => (selected.value?.tags || []).filter((t) => t.status === 'approved'))
const pendingTags = computed(() => (selected.value?.tags || []).filter((t) => t.status === 'pending'))

const totalPlays = computed(() => props.items.reduce((n, i) => n + (i.play_count || 0), 0))
const publishedCount = computed(() => props.items.filter((i) => i.published).length)
const totalApprovedTags = computed(() =>
  props.items.reduce((n, i) => n + (i.tags || []).filter((t) => t.status === 'approved').length, 0)
)

const nodes = computed(() => {
  const list = props.items || []
  const n = list.length || 1
  const radius = Math.min(155, 70 + n * 6)
  return list.map((item, i) => {
    const angle = (Math.PI * 2 * i) / n - Math.PI / 2
    const label = shortLabel(item)
    const w = Math.min(132, 28 + label.length * 7.2)
    return {
      item,
      label,
      color: palette[i % palette.length],
      x: hub.x + Math.cos(angle) * radius,
      y: hub.y + Math.sin(angle) * radius,
      w
    }
  })
})

watch(
  () => props.items,
  (list) => {
    if (selectedId.value && !list.some((i) => i.id === selectedId.value)) {
      selectedId.value = null
    }
  }
)

function shortLabel(item) {
  const raw = (item.caption || item.filename || 'clip').replace(/\.[^.]+$/, '')
  return raw.length > 14 ? raw.slice(0, 12) + '…' : raw
}

function spoke(x1, y1, x2, y2) {
  const mx = (x1 + x2) / 2
  const my = (y1 + y2) / 2
  return `M ${x1} ${y1} Q ${mx} ${my} ${x2} ${y2}`
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
  grid-template-columns: minmax(0, 1.35fr) minmax(260px, 0.9fr);
  gap: 0.85rem;
  min-height: 420px;
}

.video-map__stage {
  border: 1px solid rgba(15, 143, 124, 0.18);
  border-radius: 14px;
  background: linear-gradient(165deg, #f7fbfa, #eef6f4);
  overflow: hidden;
}

.video-map__svg {
  width: 100%;
  height: 420px;
  display: block;
}

.video-map__link {
  stroke-width: 2.2;
  opacity: 0.45;
  transition: opacity 160ms ease, stroke-width 160ms ease;
}

.video-map__link.active {
  opacity: 1;
  stroke-width: 3.4;
}

.hub-title {
  fill: #fff;
  font-size: 15px;
  font-weight: 700;
}

.hub-sub {
  fill: rgba(232, 252, 247, 0.9);
  font-size: 13px;
}

.hub-pulse {
  opacity: 0.35;
  animation: hubPulse 2.8s ease-in-out infinite;
}

@keyframes hubPulse {
  0%, 100% { opacity: 0.28; }
  50% { opacity: 0.45; }
}

.video-map__node {
  cursor: pointer;
}

.node-rect {
  fill: #fff;
  stroke-width: 2.2;
}

.video-map__node.active .node-rect {
  fill: #e7f7f3;
  stroke-width: 3;
}

.node-text {
  fill: #102833;
  font-size: 12px;
  font-weight: 650;
  pointer-events: none;
}

.video-map__panel {
  border: 1px solid rgba(15, 143, 124, 0.18);
  border-radius: 14px;
  background: #fff;
  padding: 0.9rem 1rem;
  box-shadow: 0 10px 24px rgba(6, 54, 66, 0.06);
}

.video-map__panel.empty {
  background: linear-gradient(165deg, #fff, #f4faf8);
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

@media (max-width: 900px) {
  .video-map__layout {
    grid-template-columns: 1fr;
  }

  .video-map__svg {
    height: 340px;
  }
}
</style>
