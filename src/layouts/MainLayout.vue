<template>
  <q-layout view="lHh Lpr lFf" class="layout-shell">
    <q-header elevated class="shell-header text-white">
      <div class="scroll-progress-wrap" aria-hidden="true">
        <div class="scroll-progress" :style="{ width: `${scrollProgress}%` }" />
      </div>
      <q-toolbar class="header-toolbar">
        <q-btn
          flat
          dense
          round
          icon="menu"
          @click="leftDrawerOpen = !leftDrawerOpen"
          aria-label="Menu"
          class="menu-btn"
        />
        <q-toolbar-title class="toolbar-title">
          <div class="crumb-pill">
            <q-icon name="terminal" size="14px" class="q-mr-xs" />
            {{ currentViewLabel }}
          </div>
        </q-toolbar-title>
        <q-space />
        <q-btn flat dense round icon="tune" class="top-icon-btn">
          <q-menu class="approach-menu">
            <q-list dense dark style="min-width: 180px;">
              <q-item-label header class="approach-menu-label">View approach</q-item-label>
              <q-item
                v-for="item in approachOptions"
                :key="item.value"
                clickable
                v-close-popup
                @click="approach = item.value"
                :active="approach === item.value"
                active-class="approach-menu-item--active"
              >
                <q-item-section avatar><q-icon :name="approach === item.value ? 'radio_button_checked' : 'radio_button_unchecked'" /></q-item-section>
                <q-item-section>{{ item.label }}</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
        <div v-show="!leftDrawerOpen" class="version-chip">v{{ appVersion }}</div>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="shell-drawer">
      <q-list class="drawer-list">
        <q-item-label header class="drawer-section">Console</q-item-label>
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
  overflow-x: hidden;
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
  gap: 0.35rem;
  flex-wrap: nowrap;
  min-height: 58px;
  overflow: hidden;
}

.toolbar-title {
  display: flex;
  align-items: center;
  flex: 1 1 auto;
  min-width: 0;
  overflow: hidden;
}

.menu-btn {
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.16);
}

.crumb-pill {
  display: inline-flex;
  align-items: center;
  border: 1px solid rgba(184, 205, 209, 0.34);
  border-radius: 999px;
  padding: 0.18rem 0.6rem;
  font-size: 0.72rem;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #d9ebea;
}

.top-icon-btn {
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.16);
  flex-shrink: 0;
}

.approach-menu {
  border: 1px solid rgba(184, 205, 209, 0.26);
  background: linear-gradient(170deg, #1c1c1c, #141414);
}

.approach-menu-label {
  color: #b8cdd1;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  font-size: 0.68rem;
}

.approach-menu-item--active {
  background: rgba(164, 108, 173, 0.22);
}

.version-chip {
  border: 1px solid rgba(255, 255, 255, 0.32);
  border-radius: 999px;
  padding: 0.2rem 0.62rem;
  font-size: 0.72rem;
  letter-spacing: 0.06em;
  color: #edf6f4;
  white-space: nowrap;
  flex-shrink: 0;
}

.shell-drawer {
  background: linear-gradient(170deg, rgba(21, 21, 21, 0.98), rgba(14, 14, 14, 0.98)) !important;
  color: #d9d9d9 !important;
  border-right: 1px solid rgba(184, 205, 209, 0.2);
}

.drawer-list {
  padding-top: 0.35rem;
  background: transparent;
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
  transition: background 170ms ease, transform 170ms ease;
  border: 1px solid transparent;
  color: #d7e1df;
}

.drawer-item :deep(.q-item__section--main) {
  color: #d7e1df;
  font-size: 0.9rem;
}

.drawer-item :deep(.q-item__section--avatar .q-icon) {
  color: #9bc0b8;
}

.drawer-item:hover {
  background: rgba(90, 139, 126, 0.16);
  transform: translateX(2px);
  border-color: rgba(90, 139, 126, 0.28);
}

.drawer-item:hover :deep(.q-item__section--avatar .q-icon) {
  color: #c090c9;
}

.drawer-item--active {
  background: rgba(164, 108, 173, 0.18);
  border: 1px solid rgba(164, 108, 173, 0.48);
}

.drawer-item--active :deep(.q-item__section--main) {
  color: #f0f3f2;
  font-weight: 600;
}

.shell-footer {
  background: linear-gradient(130deg, #1d1d1d, #212121);
  border-top: 1px solid rgba(184, 205, 209, 0.26);
  color: #b9c1c3;
}

@media (max-width: 940px) {
  .crumb-pill {
    display: none;
  }
}

@media (max-width: 1320px) {
  .version-chip {
    display: none;
  }
}
</style>

