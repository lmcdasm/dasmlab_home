<template>
  <q-layout view="lHh Lpr lFf" class="layout-shell">
    <q-header elevated class="shell-header">
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
        <div class="header-links">
          <q-btn
            v-for="item in topNav.filter((entry) => !!entry.to)"
            :key="`top-${item.to}`"
            flat
            dense
            no-caps
            class="header-link-btn"
            :label="item.label"
            :to="item.to"
          />
          <q-btn-dropdown
            flat
            dense
            no-caps
            class="header-link-btn header-projects-btn"
            label="Projects"
            dropdown-icon="expand_more"
          >
            <q-list>
              <q-item
                v-for="item in projectNav"
                :key="`project-${item.to}`"
                clickable
                v-close-popup
                :to="item.to"
              >
                <q-item-section avatar><q-icon :name="item.icon" /></q-item-section>
                <q-item-section>{{ item.label }}</q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </div>
        <div class="right-controls">
          <q-btn
            v-if="oidcEnabled && !authenticated"
            flat
            dense
            no-caps
            class="header-link-btn"
            icon="login"
            label="Sign in"
            @click="login"
          />
          <q-btn
            v-else-if="oidcEnabled && authenticated"
            flat
            dense
            no-caps
            class="header-link-btn"
            icon="logout"
            :label="user?.preferred_username || 'Sign out'"
            @click="logout"
          >
            <q-tooltip>{{ isAdmin ? 'Admin (owner)' : 'Signed in' }}</q-tooltip>
          </q-btn>
          <q-btn flat dense round icon="tune" class="top-icon-btn">
            <q-menu class="approach-menu">
              <q-list dense style="min-width: 180px;">
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
          <div
            class="version-chip"
            data-build-tag
            :title="'Build ' + appVersion"
          >{{ appVersion }}</div>
        </div>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      bordered
      overlay
      behavior="desktop"
      content-class="shell-drawer-content"
      class="shell-drawer"
    >
      <q-list class="drawer-list">
        <q-item class="drawer-topbar">
          <q-item-section class="text-caption drawer-topbar-label">Navigation</q-item-section>
          <q-item-section side>
            <q-btn
              flat
              dense
              round
              icon="close"
              size="sm"
              class="drawer-close-btn"
              @click="leftDrawerOpen = false"
            />
          </q-item-section>
        </q-item>
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

        <q-separator inset class="q-my-md" />
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
import { useAuth } from 'src/composables/useAuth'

const { oidcEnabled, authenticated, isAdmin, user, login, logout } = useAuth()
import { useRoute } from 'vue-router'
import VisitCounter from 'src/components/VisitCounter.vue'
import { useApproach } from 'src/composables/useApproach'

const route = useRoute()
const leftDrawerOpen = ref(false)
const scrollProgress = ref(0)
const appVersion = ref(process.env.APP_VERSION || import.meta.env.APP_VERSION || 'dev')
const { approach, approachOptions } = useApproach()

const topNav = [
  { label: 'Home', icon: 'home', to: '/' },
  { label: 'About DASMLAB', icon: 'school', to: '/about' },
  { label: 'Surfing', icon: 'sailing', to: '/surfing' },
  { label: 'Live Cams', icon: 'videocam', href: 'https://camera-scrape.apps.2026-prod-1.ocp.dasmlab.org/' },
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
  '/surfing': 'Surfing',
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
  // Prefer runtime /version.json (same pattern as cheapcloud data-build-tag) so chip matches image tag.
  fetch(`/version.json?${Date.now()}`, { cache: 'no-store' })
    .then((r) => (r.ok ? r.json() : null))
    .then((info) => {
      const build = info?.build || info?.version
      if (build) appVersion.value = build
    })
    .catch(() => {})

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
    radial-gradient(circle at 4% 8%, rgba(63, 122, 107, 0.1), transparent 34%),
    radial-gradient(circle at 96% 10%, rgba(158, 115, 178, 0.08), transparent 34%),
    linear-gradient(180deg, #f7fafd, #edf4f8);
}

.shell-header {
  position: relative;
  background: rgba(255, 255, 255, 0.96);
  border-bottom: 1px solid rgba(36, 61, 81, 0.14);
  z-index: 4000;
  color: #243344;
}

.scroll-progress-wrap {
  position: absolute;
  inset: 0 auto auto 0;
  width: 100%;
  height: 2px;
  background: rgba(34, 63, 90, 0.08);
}

.scroll-progress {
  height: 100%;
  background: linear-gradient(90deg, #33a08a, #9f74b5);
  transition: width 80ms linear;
}

.header-toolbar {
  backdrop-filter: blur(6px);
  gap: 0.42rem;
  flex-wrap: nowrap;
  min-height: 58px;
  overflow: hidden;
}

.toolbar-title {
  display: flex;
  align-items: center;
  flex: 0 1 auto;
  max-width: 26%;
  min-width: 0;
  overflow: hidden;
}

.header-links {
  display: flex;
  align-items: center;
  gap: 0.2rem;
}

.header-link-btn {
  color: #304356;
  border-radius: 10px;
  padding: 0 0.4rem;
}

.header-link-btn:hover {
  background: rgba(63, 122, 107, 0.1);
}

.header-link-btn.q-router-link--active {
  color: #216f62;
  background: rgba(63, 122, 107, 0.13);
  font-weight: 600;
}

.header-projects-btn :deep(.q-btn__dropdown-icon) {
  color: #607587;
}

.menu-btn {
  border: 1px solid rgba(36, 61, 81, 0.18);
  color: #2a3d50;
  background: #ffffff;
}

.crumb-pill {
  display: inline-flex;
  align-items: center;
  border: 1px solid rgba(36, 61, 81, 0.18);
  border-radius: 999px;
  padding: 0.18rem 0.6rem;
  font-size: 0.72rem;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #2a3d50;
  background: rgba(246, 251, 255, 0.88);
}

.top-icon-btn {
  border: 1px solid rgba(36, 61, 81, 0.18);
  background: #ffffff;
  color: #2a3d50;
  flex-shrink: 0;
}

.right-controls {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 0.42rem;
  flex-shrink: 0;
}

.approach-menu {
  border: 1px solid rgba(36, 61, 81, 0.16);
  background: linear-gradient(170deg, #ffffff, #f3f8fc);
}

.approach-menu-label {
  color: #5f7383;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  font-size: 0.68rem;
}

.approach-menu-item--active {
  background: rgba(151, 110, 176, 0.14);
}

.version-chip {
  border: 1px solid rgba(36, 61, 81, 0.2);
  border-radius: 999px;
  padding: 0.2rem 0.62rem;
  font-size: 0.72rem;
  letter-spacing: 0.06em;
  color: #2a3d50;
  background: rgba(246, 251, 255, 0.92);
  white-space: nowrap;
  flex-shrink: 0;
}

.shell-drawer {
  background: linear-gradient(170deg, rgba(255, 255, 255, 0.99), rgba(245, 250, 253, 0.99)) !important;
  color: #26384a !important;
  border-right: 1px solid rgba(36, 61, 81, 0.16);
  top: 58px !important;
  height: calc(100% - 58px) !important;
  z-index: 3000 !important;
}

.shell-drawer :deep(.shell-drawer-content) {
  background: linear-gradient(178deg, #ffffff, #f4f9fc) !important;
  color: #253749 !important;
}

.drawer-list {
  padding-top: 0.35rem;
  background: transparent;
}

.drawer-topbar {
  min-height: 44px;
  border-bottom: 1px solid rgba(36, 61, 81, 0.12);
}

.drawer-topbar-label {
  color: #5e7587;
  letter-spacing: 0.09em;
  text-transform: uppercase;
}

.drawer-close-btn {
  color: #5b7284;
  border: 1px solid rgba(36, 61, 81, 0.18);
  background: rgba(255, 255, 255, 0.9);
}

.drawer-section {
  color: #61788a;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  font-size: 0.7rem;
  font-weight: 600;
}

.drawer-item {
  border-radius: 11px;
  margin: 0.18rem 0.35rem;
  transition: background 170ms ease, transform 170ms ease, border-color 170ms ease;
  border: 1px solid rgba(36, 61, 81, 0.12);
  color: #2b3e51;
  background: rgba(255, 255, 255, 0.86);
}

.drawer-item :deep(.q-item__section--main) {
  color: #2b3e51;
  font-size: 0.9rem;
  font-weight: 500;
}

.drawer-item :deep(.q-item__section--avatar .q-icon) {
  color: #4e8f82;
}

.drawer-item:hover {
  background: linear-gradient(120deg, rgba(63, 122, 107, 0.14), rgba(158, 115, 178, 0.12));
  transform: translateX(2px);
  border-color: rgba(152, 111, 178, 0.34);
}

.drawer-item:hover :deep(.q-item__section--avatar .q-icon) {
  color: #8f65a8;
}

.drawer-item--active {
  background: linear-gradient(120deg, rgba(63, 122, 107, 0.18), rgba(158, 115, 178, 0.16));
  border: 1px solid rgba(152, 111, 178, 0.45);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.55);
}

.drawer-item--active :deep(.q-item__section--main) {
  color: #213445;
  font-weight: 600;
}

.drawer-item--active :deep(.q-item__section--avatar .q-icon) {
  color: #825f9b;
}

.shell-footer {
  background: linear-gradient(130deg, #f7fafd, #eef4f8);
  border-top: 1px solid rgba(36, 61, 81, 0.14);
  color: #607384;
}

@media (max-width: 940px) {
  .crumb-pill {
    display: none;
  }

  .header-links {
    display: none;
  }
}

@media (max-width: 1320px) {
  .toolbar-title {
    max-width: 18%;
  }
}

@media (max-width: 1240px) {
  .header-links {
    display: none;
  }
}

@media (max-width: 1140px) {
  .version-chip {
    display: none;
  }
}
</style>

