<template>
  <section class="behind dasm-panel" aria-label="Behind the Design">
    <div class="dasm-caps">Behind the Design</div>
    <h2 class="behind__title">How this was implemented</h2>
    <p v-if="intro" class="behind__intro">{{ intro }}</p>
    <ul class="behind__stack">
      <li v-for="item in stack" :key="item">
        <router-link
          v-if="topicSlug(item)"
          class="behind__chip"
          :to="`/topics/${topicSlug(item)}`"
        >
          {{ item }}
        </router-link>
        <span v-else class="behind__chip behind__chip--plain">{{ item }}</span>
      </li>
    </ul>
    <p class="behind__note">
      Click a technology to see everything we’ve built with it across the lab.
    </p>
  </section>
</template>

<script setup>
const props = defineProps({
  stack: { type: Array, default: () => [] },
  topics: { type: Array, default: () => [] },
  intro: String
})

const ALIASES = {
  'vue 3': 'vue',
  vue: 'vue',
  quasar: 'quasar',
  gin: 'gin',
  'gin (go)': 'gin',
  go: 'gin',
  oidc: 'oidc',
  keycloak: 'oidc',
  'oidc / keycloak': 'oidc',
  metallb: 'metallb',
  openshift: 'openshift',
  'openshift / gitops': 'openshift',
  vite: null,
  nginx: null
}

function topicSlug(label) {
  const key = String(label || '').toLowerCase()
  if (props.topics?.length) {
    const hit = props.topics.find((t) => key.includes(t) || t.includes(key.replace(/\s+/g, '')))
    if (hit) return hit
  }
  if (Object.prototype.hasOwnProperty.call(ALIASES, key)) return ALIASES[key]
  const compact = key.replace(/[^a-z0-9]+/g, '')
  for (const [k, v] of Object.entries(ALIASES)) {
    if (compact.includes(k.replace(/[^a-z0-9]+/g, '')) && v) return v
  }
  return null
}
</script>

<style scoped>
.behind__title {
  margin: 0.35rem 0 0.5rem;
  font-size: 1.2rem;
  color: #1d2b36;
}

.behind__intro {
  margin: 0 0 0.75rem;
  color: #4a5d6d;
  line-height: 1.55;
}

.behind__stack {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.behind__chip {
  display: inline-block;
  padding: 0.35rem 0.7rem;
  border: 1.5px solid var(--dasm-border-strong);
  border-radius: 10px;
  text-decoration: none;
  color: #1f6f62;
  font-weight: 600;
  font-size: 0.88rem;
  background: rgba(63, 159, 142, 0.08);
  transition: transform 160ms ease, background 160ms ease;
}

.behind__chip:hover {
  transform: translateY(-1px);
  background: rgba(63, 159, 142, 0.16);
}

.behind__chip--plain {
  color: #435564;
  font-weight: 500;
  background: #f3f8f6;
}

.behind__note {
  margin: 0.85rem 0 0;
  font-size: 0.85rem;
  color: #5a6f80;
}
</style>
