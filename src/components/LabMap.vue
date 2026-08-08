<template>
  <div class="lab-map">
    <div class="map-toolbar">
      <span class="hint">Click a lane — detail opens beside the map, not on top of it.</span>
    </div>

    <div class="map-layout">
      <div class="map-stage">
        <svg class="map-svg" :viewBox="`0 0 ${vbW} ${vbH}`" preserveAspectRatio="xMidYMid meet">
          <defs>
            <pattern id="labGrid" width="24" height="24" patternUnits="userSpaceOnUse">
              <path d="M 24 0 L 0 0 0 24" fill="none" stroke="rgba(31, 111, 98, 0.07)" stroke-width="1" />
            </pattern>
            <filter id="labShadow" x="-40%" y="-40%" width="180%" height="180%">
              <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0a1a22" flood-opacity="0.22" />
            </filter>
            <linearGradient id="labHubGrad" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stop-color="#1f6f62" />
              <stop offset="55%" stop-color="#2f8f7d" />
              <stop offset="100%" stop-color="#976eb0" />
            </linearGradient>
            <radialGradient id="labWash" cx="50%" cy="45%" r="55%">
              <stop offset="0%" stop-color="rgba(47, 143, 125, 0.16)" />
              <stop offset="100%" stop-color="rgba(47, 143, 125, 0)" />
            </radialGradient>
          </defs>

          <rect width="100%" height="100%" fill="url(#labGrid)" />
          <rect width="100%" height="100%" fill="url(#labWash)" />

          <path
            v-for="b in branches"
            :key="'link-' + b.id"
            :d="spoke(hub.x, hub.y, b.x, b.y)"
            class="link"
            :class="{ active: selectedId === b.id }"
            :stroke="b.color"
            fill="none"
          />

          <g class="hub" filter="url(#labShadow)">
            <circle :cx="hub.x" :cy="hub.y" r="68" fill="url(#labHubGrad)" class="hub-pulse" />
            <circle :cx="hub.x" :cy="hub.y" r="58" fill="url(#labHubGrad)" />
            <text :x="hub.x" :y="hub.y - 6" text-anchor="middle" class="hub-title">DASMLAB</text>
            <text :x="hub.x" :y="hub.y + 14" text-anchor="middle" class="hub-sub">LAB</text>
          </g>

          <g
            v-for="b in branches"
            :key="'node-' + b.id"
            class="lane-node"
            :class="{ active: selectedId === b.id }"
            filter="url(#labShadow)"
            tabindex="0"
            role="button"
            :aria-label="b.title"
            @click="select(b.id)"
            @keydown.enter="select(b.id)"
          >
            <rect
              :x="b.x - b.w / 2"
              :y="b.y - 22"
              :width="b.w"
              height="44"
              rx="12"
              class="lane-rect"
              :stroke="b.color"
            />
            <text :x="b.x" :y="b.y + 5" text-anchor="middle" class="lane-text">{{ b.short }}</text>
          </g>
        </svg>
      </div>

      <aside class="detail-panel" :class="{ empty: !selected }">
        <template v-if="selected">
          <div class="detail-head">
            <div>
              <div class="detail-eyebrow" :style="{ color: selected.color }">{{ selected.eyebrow }}</div>
              <div class="detail-title">{{ selected.title }}</div>
              <p class="detail-blurb">{{ selected.blurb }}</p>
            </div>
            <button type="button" class="detail-close" aria-label="Close" @click="selectedId = null">×</button>
          </div>

          <div class="chip-row">
            <span v-for="tag in selected.tags" :key="tag" class="meta-chip">{{ tag }}</span>
          </div>

          <div class="flow-block">
            <div class="flow-label builds">Builds</div>
            <div class="chip-list">
              <span v-for="item in selected.builds" :key="item" class="asset-chip builds">{{ item }}</span>
            </div>
          </div>

          <div class="flow-block">
            <div class="flow-label teaches">Teaches</div>
            <div class="chip-list">
              <span v-for="item in selected.teaches" :key="item" class="asset-chip teaches">{{ item }}</span>
            </div>
          </div>

          <button type="button" class="detail-cta" @click="go(selected.route)">
            Open {{ selected.short }}
            <span aria-hidden="true">→</span>
          </button>
        </template>
        <template v-else>
          <div class="empty-hint">
            <div class="empty-glyph" aria-hidden="true" />
            <div class="empty-title">Select a lane</div>
            <div class="empty-copy">Builds / teaches chips and a deep-link appear here — map stays clean.</div>
          </div>
        </template>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const selectedId = ref(null)
const vbW = 640
const vbH = 640
const hub = { x: 320, y: 320 }

const lanes = [
  {
    id: 'frontend',
    short: 'Frontend',
    title: 'Frontend Projects',
    eyebrow: 'UI craft',
    blurb: 'Quasar / Vue surfaces with personality — maps, stages, motion.',
    color: '#c2185b',
    route: '/projects/frontend',
    tags: ['Vue', 'Quasar', 'SVG'],
    builds: ['Lab maps', 'Studio shells', 'Design systems'],
    teaches: ['Motion restraint', 'Brand-first heroes'],
  },
  {
    id: 'backend',
    short: 'Backend',
    title: 'Backend Projects',
    eyebrow: 'APIs',
    blurb: 'Go services, envelopes, and GitOps-ready containers.',
    color: '#ef6c00',
    route: '/projects/backend',
    tags: ['Go', 'Gin', 'GHCR'],
    builds: ['Surfing API', 'Carousel', 'Cheapcloud brokers'],
    teaches: ['Idempotent builds', 'SemVer tags'],
  },
  {
    id: 'aiml',
    short: 'AI / ML',
    title: 'AI / ML Tech Dives',
    eyebrow: 'Experiments',
    blurb: 'Field notes from models, pipelines, and practical demos.',
    color: '#7b1fa2',
    route: '/projects/ai-ml',
    tags: ['Notebooks', 'Pipelines'],
    builds: ['Eval harnesses', 'Prompt labs'],
    teaches: ['Reproducible runs', 'Cost awareness'],
  },
  {
    id: 'cloud',
    short: 'Cloud',
    title: 'Cloud Provider Techs',
    eyebrow: 'Providers',
    blurb: 'Azure / AWS / GCP patterns — and the cheap+dirty escapes.',
    color: '#1565c0',
    route: '/projects/cloud',
    tags: ['Azure', 'AWS', 'GCP'],
    builds: ['Blob origins', 'IAM drills'],
    teaches: ['Egress math', 'Region picks'],
  },
  {
    id: 'infra',
    short: 'Infra',
    title: 'Infrastructure',
    eyebrow: 'Platform',
    blurb: 'OCP clusters, GitOps, LVMS, and the wires under the lab.',
    color: '#2e7d32',
    route: '/projects/infrastructure',
    tags: ['OCP', 'Argo', 'LVMS'],
    builds: ['2026-prod lanes', 'Live-cicd'],
    teaches: ['Prune discipline', 'SCC hygiene'],
  },
  {
    id: 'security',
    short: 'Security',
    title: 'Security Projects',
    eyebrow: 'Hardening',
    blurb: 'OIDC, secrets, SCC, and “don’t leave the basement open.”',
    color: '#c62828',
    route: '/projects/security',
    tags: ['OIDC', 'Secrets', 'SCC'],
    builds: ['Pull secrets', 'Route TLS'],
    teaches: ['Least privilege', 'Public API caution'],
  },
  {
    id: 'surfing',
    short: 'Surfing',
    title: 'Surfing',
    eyebrow: 'Media',
    blurb: 'Video days off basement disk — R2 + edge cache next.',
    color: '#00838f',
    route: '/surfing',
    tags: ['VOD', 'R2', 'CDN'],
    builds: ['Day manifests', 'Media broker'],
    teaches: ['Origin vs edge', 'Egress free paths'],
  },
]

function textWidth(s) {
  return Math.min(150, Math.max(92, String(s).length * 7.4 + 30))
}

function spoke(x1, y1, x2, y2) {
  const mx = (x1 + x2) / 2
  const my = (y1 + y2) / 2
  const dx = x2 - x1
  const dy = y2 - y1
  const cx = mx - dy * 0.08
  const cy = my + dx * 0.06
  return `M ${x1} ${y1} Q ${cx} ${cy} ${x2} ${y2}`
}

const branches = computed(() => {
  const n = lanes.length
  const r = 214
  return lanes.map((lane, i) => {
    const angle = (Math.PI * 2 * i) / n - Math.PI / 2
    return {
      ...lane,
      x: hub.x + Math.cos(angle) * r,
      y: hub.y + Math.sin(angle) * r,
      w: textWidth(lane.short),
    }
  })
})

const selected = computed(() => {
  if (!selectedId.value) return null
  return branches.value.find((b) => b.id === selectedId.value) || null
})

function select(id) {
  selectedId.value = selectedId.value === id ? null : id
}

function go(route) {
  router.push(route)
}
</script>

<style scoped>
.lab-map {
  --ink: #12202c;
  --muted: #5a6f80;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.map-svg {
  width: 100%;
  height: 100%;
  display: block;
  min-height: 0;
}

.map-toolbar {
  display: none;
}

.map-layout {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: stretch;
  gap: 0.85rem;
  width: 100%;
  flex: 1;
  min-height: 0;
}

.map-stage {
  border-radius: 16px;
  overflow: hidden;
  border: 2px solid rgba(18, 72, 64, 0.45);
  background:
    radial-gradient(ellipse at 50% 42%, rgba(47, 143, 125, 0.16), transparent 56%),
    linear-gradient(165deg, #eef6f3 0%, #dfe8f0 100%);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.75),
    0 10px 22px rgba(10, 28, 36, 0.12);
  flex: 1 1 420px;
  width: min(100%, 560px);
  height: min(52vh, 520px);
  max-width: 560px;
  max-height: 520px;
  min-height: 320px;
  aspect-ratio: 1 / 1;
  display: grid;
  place-items: center;
}

.detail-panel {
  border-radius: 16px;
  border: 1.5px solid rgba(31, 111, 98, 0.22);
  background: linear-gradient(175deg, #ffffff, #f5f9fb);
  padding: 1rem 1.05rem;
  flex: 1 1 280px;
  width: min(100%, 340px);
  min-height: 280px;
  max-height: min(52vh, 520px);
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
  box-shadow: 0 10px 22px rgba(18, 40, 52, 0.08);
}

@media (max-width: 720px) {
  .detail-panel {
    width: 100%;
    max-width: none;
    max-height: none;
    min-height: 180px;
  }

  .map-stage {
    flex: 1 1 auto;
    width: min(100%, 100%);
    height: auto;
    min-height: 280px;
    max-height: min(70vw, 400px);
  }
}

.link {
  stroke-width: 2.5;
  opacity: 0.42;
  stroke-linecap: round;
  transition: opacity 0.2s ease, stroke-width 0.2s ease;
}

.link.active {
  opacity: 1;
  stroke-width: 3.5;
}

.hub-title,
.hub-sub {
  fill: #fff;
  font-family: "Segoe UI", system-ui, sans-serif;
  pointer-events: none;
}

.hub-title {
  font-size: 15px;
  font-weight: 800;
  letter-spacing: 0.08em;
}

.hub-sub {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.22em;
  opacity: 0.9;
}

.hub-pulse {
  animation: hubPulse 3.2s ease-in-out infinite;
  opacity: 0.38;
}

.lane-node {
  cursor: pointer;
  outline: none;
}

.lane-rect {
  fill: #fff;
  stroke-width: 2.75;
  transition: fill 0.15s ease, stroke-width 0.15s ease, transform 0.15s ease;
}

.lane-node:hover .lane-rect,
.lane-node.active .lane-rect {
  fill: #f0faf7;
  stroke-width: 3.5;
}

.lane-text {
  fill: var(--ink);
  font-size: 12.5px;
  font-weight: 700;
  font-family: "Segoe UI", system-ui, sans-serif;
  pointer-events: none;
}

.detail-panel.empty {
  justify-content: center;
}

.detail-head {
  display: flex;
  justify-content: space-between;
  gap: 0.75rem;
  align-items: flex-start;
}

.detail-eyebrow {
  font-size: 0.68rem;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.detail-title {
  font-size: 1.15rem;
  font-weight: 750;
  color: var(--ink);
  line-height: 1.2;
  margin-top: 0.15rem;
}

.detail-blurb {
  margin: 0.35rem 0 0;
  color: var(--muted);
  font-size: 0.9rem;
  line-height: 1.45;
}

.detail-close {
  border: 1px solid rgba(28, 52, 73, 0.18);
  background: #fff;
  width: 2rem;
  height: 2rem;
  border-radius: 999px;
  font-size: 1.2rem;
  line-height: 1;
  color: var(--muted);
  cursor: pointer;
}

.chip-row,
.chip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.meta-chip {
  font-size: 0.72rem;
  font-weight: 650;
  letter-spacing: 0.04em;
  padding: 0.28rem 0.55rem;
  border-radius: 999px;
  border: 1px solid rgba(31, 111, 98, 0.28);
  background: rgba(47, 143, 125, 0.08);
  color: #1f6f62;
}

.flow-label {
  font-size: 0.68rem;
  font-weight: 800;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  margin-bottom: 0.4rem;
}

.flow-label.builds {
  color: #1b5e20;
}

.flow-label.teaches {
  color: #1565c0;
}

.asset-chip {
  font-size: 0.78rem;
  line-height: 1.3;
  padding: 0.4rem 0.6rem;
  border-radius: 9px;
  border: 1px solid transparent;
}

.asset-chip.builds {
  background: #e8f5e9;
  border-color: #81c784;
  color: #1b5e20;
}

.asset-chip.teaches {
  background: #e3f2fd;
  border-color: #64b5f6;
  color: #0d47a1;
}

.detail-cta {
  margin-top: auto;
  border: none;
  border-radius: 12px;
  padding: 0.72rem 1rem;
  background: linear-gradient(135deg, #1f6f62, #2f8f7d);
  color: #fff;
  font-weight: 700;
  letter-spacing: 0.03em;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  box-shadow: 0 10px 22px rgba(31, 111, 98, 0.28);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.detail-cta:hover {
  transform: translateY(-1px);
  box-shadow: 0 14px 26px rgba(31, 111, 98, 0.34);
}

.empty-hint {
  text-align: center;
  padding: 1.5rem 0.75rem;
  color: var(--muted);
}

.empty-glyph {
  width: 42px;
  height: 42px;
  margin: 0 auto 0.75rem;
  border-radius: 14px;
  border: 2px dashed rgba(31, 111, 98, 0.35);
  background:
    radial-gradient(circle at 50% 50%, rgba(47, 143, 125, 0.2), transparent 65%);
}

.empty-title {
  font-weight: 750;
  color: var(--ink);
  margin-bottom: 0.25rem;
}

.empty-copy {
  font-size: 0.82rem;
  line-height: 1.45;
  max-width: 220px;
  margin: 0 auto;
}

@keyframes hubPulse {
  0%,
  100% {
    opacity: 0.28;
    transform-origin: center;
  }
  50% {
    opacity: 0.55;
  }
}
</style>
