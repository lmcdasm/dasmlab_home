<template>
  <q-layout view="lHh Lpr lFf" class="layout-shell">
    <q-header elevated class="shell-header text-white">
      <div class="scroll-progress-wrap" aria-hidden="true">
        <div class="scroll-progress" :style="{ width: `${scrollProgress}%` }" />
      </div>
      <q-toolbar class="header-toolbar">
        <q-btn flat dense round icon="menu" @click="leftDrawerOpen = !leftDrawerOpen" aria-label="Menu" />
        <q-toolbar-title class="toolbar-title">
          <span class="title-primary">DASMLAB</span>
          <span class="title-secondary">{{ currentViewLabel }}</span>
        </q-toolbar-title>
        <q-space />
        <q-select
          v-model="approach"
          :options="approachOptions"
          option-value="value"
          option-label="label"
          emit-value
          map-options
          dense
          borderless
          dark
          class="approach-select q-mr-md"
        >
          <template #prepend>
            <q-icon name="tune" size="xs" />
          </template>
        </q-select>
        <div class="version-chip">v{{ appVersion }}</div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="shell-drawer">
      <q-list class="drawer-list">
        <q-item-label header class="drawer-section">Studio</q-item-label>
        <q-item
          v-for="item in topNav"
          :key="item.to || item.href"
          clickable
          v-ripple
          :to="item.to"
          :tag="item.href ? 'a' : 'div'"
          :href="item.href"
          :target="item.href ? '_blank' : undefined"
          class="drawer-item"
        >
          <q-item-section avatar><q-icon :name="item.icon" /></q-item-section>
          <q-item-section>{{ item.label }}</q-item-section>
        </q-item>

        <q-separator dark inset class="q-my-md" />
        <q-item-label header class="drawer-section">Project Lanes</q-item-label>
        <q-item
          v-for="item in projectNav"
          :key="item.to"
          clickable
          v-ripple
          :to="item.to"
          class="drawer-item"
          active-class="drawer-item--active"
        >
          <q-item-section avatar><q-icon :name="item.icon" /></q-item-section>
          <q-item-section>{{ item.label }}</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer class="shell-footer text-center q-pa-sm">
      © 2026 DASMLAB Inc.
      <span class="q-ml-md">Living lab and portfolio.</span>
      <span class="q-ml-md"><VisitCounter /></span>
    </q-footer>
  </q-layout>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import VisitCounter from 'src/components/VisitCounter.vue'
import { useApproach } from 'src/composables/useApproach'

const route = useRoute()
const leftDrawerOpen = ref(true)
const scrollProgress = ref(0)
const appVersion = (process.env.APP_VERSION || import.meta.env.APP_VERSION || 'dev')
const { approach, approachOptions } = useApproach()

const topNav = [
  { label: 'Home', icon: 'home', to: '/' },
  { label: 'About DASMLAB', icon: 'school', to: '/about' },
  { label: 'Contact', icon: 'mail', to: '/contact' },
  { label: 'GitHub', icon: 'code', href: 'https://github.com/lmcdasm' }
]

const projectNav = [
  { label: 'Frontend', icon: 'web', to: '/projects/frontend' },
  { label: 'Backend', icon: 'storage', to: '/projects/backend' },
  { label: 'AI/ML', icon: 'psychology', to: '/projects/ai-ml' },
  { label: 'Cloud', icon: 'cloud', to: '/projects/cloud' },
  { label: 'Infrastructure', icon: 'dns', to: '/projects/infrastructure' },
  { label: 'Security', icon: 'security', to: '/projects/security' }
]

const routeLabelMap = {
  '/': 'Home',
  '/about': 'About',
  '/contact': 'Contact',
  '/projects/frontend': 'Frontend',
  '/projects/backend': 'Backend',
  '/projects/ai-ml': 'AI/ML',
  '/projects/cloud': 'Cloud',
  '/projects/infrastructure': 'Infrastructure',
  '/projects/security': 'Security'
}

const currentViewLabel = computed(() => routeLabelMap[route.path] || 'Lab')

const updateScrollProgress = () => {
  const doc = document.documentElement
  const max = doc.scrollHeight - doc.clientHeight
  if (max <= 0) {
    scrollProgress.value = 0
    return
  }
  scrollProgress.value = Math.max(0, Math.min(100, (doc.scrollTop / max) * 100))
}

onMounted(() => {
  updateScrollProgress()
  window.addEventListener('scroll', updateScrollProgress, { passive: true })
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateScrollProgress)
})
</script>

<style scoped>
.layout-shell {
  background:
    radial-gradient(circle at 3% 7%, rgba(60, 107, 94, 0.16), transparent 32%),
    radial-gradient(circle at 96% 10%, rgba(164, 108, 173, 0.14), transparent 34%),
    #121212;
}

.shell-header {
  position: relative;
  background: linear-gradient(140deg, #2d5248, #3c6b5e, #2f564d);
  border-bottom: 1px solid rgba(184, 205, 209, 0.28);
}

.scroll-progress-wrap {
  position: absolute;
  inset: 0 auto auto 0;
  width: 100%;
  height: 2px;
  background: rgba(0, 0, 0, 0.18);
}

.scroll-progress {
  height: 100%;
  background: linear-gradient(90deg, #c090c9, #b8cdd1);
  transition: width 80ms linear;
}

.header-toolbar {
  backdrop-filter: blur(8px);
}

.toolbar-title {
  display: flex;
  flex-direction: column;
  gap: 0.05rem;
}

.title-primary {
  font-size: 0.98rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  font-weight: 700;
}

.title-secondary {
  font-size: 0.76rem;
  color: #deece9;
  letter-spacing: 0.09em;
  text-transform: uppercase;
}

.approach-select {
  min-width: 154px;
  border: 1px solid rgba(255, 255, 255, 0.26);
  border-radius: 10px;
  padding: 0 0.35rem;
  background: rgba(0, 0, 0, 0.15);
}

.version-chip {
  border: 1px solid rgba(255, 255, 255, 0.32);
  border-radius: 999px;
  padding: 0.2rem 0.62rem;
  font-size: 0.72rem;
  letter-spacing: 0.06em;
  color: #edf6f4;
}

.shell-drawer {
  background: linear-gradient(170deg, rgba(21, 21, 21, 0.98), rgba(14, 14, 14, 0.98));
  color: #d9d9d9;
}

.drawer-list {
  padding-top: 0.5rem;
}

.drawer-section {
  color: #a5bbb6;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  font-size: 0.7rem;
}

.drawer-item {
  border-radius: 10px;
  margin: 0.15rem 0.35rem;
  transition: background 170ms ease;
}

.drawer-item:hover {
  background: rgba(90, 139, 126, 0.16);
}

.drawer-item--active {
  background: rgba(164, 108, 173, 0.18);
  border: 1px solid rgba(164, 108, 173, 0.48);
}

.shell-footer {
  background: linear-gradient(130deg, #1d1d1d, #212121);
  border-top: 1px solid rgba(184, 205, 209, 0.26);
  color: #b9c1c3;
}
</style>

