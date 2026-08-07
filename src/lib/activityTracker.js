import { postActivity } from 'src/services/authApi'
import { useAuth } from 'src/composables/useAuth'

const INPUT_IDLE_MS = 5000
const SURFING_HOST = import.meta.env.VITE_SURFING_API_HOST || '/api/surfing'

let started = false
let pageEnteredAt = 0
let currentPath = ''
let visibleAccumMs = 0
let engagedAccumMs = 0
let visibleSegmentStart = 0
let engagedSegmentStart = 0
let lastInputAt = 0
let visible = true
let engaged = false
function now() {
  return typeof performance !== 'undefined' && performance.now ? performance.now() : Date.now()
}

function isDocumentVisible() {
  if (typeof document === 'undefined') return true
  return document.visibilityState !== 'hidden'
}

function isWindowFocused() {
  if (typeof document === 'undefined') return true
  return typeof document.hasFocus !== 'function' || document.hasFocus()
}

function foregroundActive() {
  return isDocumentVisible() && isWindowFocused()
}

function closeVisibleSegment() {
  if (visibleSegmentStart > 0) {
    visibleAccumMs += Math.max(0, now() - visibleSegmentStart)
    visibleSegmentStart = 0
  }
}

function closeEngagedSegment() {
  if (engagedSegmentStart > 0) {
    engagedAccumMs += Math.max(0, now() - engagedSegmentStart)
    engagedSegmentStart = 0
  }
  engaged = false
}

function syncVisibility() {
  const fg = foregroundActive()
  if (fg && !visible) {
    visible = true
    visibleSegmentStart = now()
    if (now() - lastInputAt < INPUT_IDLE_MS) {
      engaged = true
      engagedSegmentStart = now()
    }
  } else if (!fg && visible) {
    closeEngagedSegment()
    closeVisibleSegment()
    visible = false
  }
}

function onInput() {
  lastInputAt = now()
  if (!foregroundActive()) return
  if (!engaged) {
    engaged = true
    engagedSegmentStart = now()
  }
}

function tickEngagementIdle() {
  if (!engaged) return
  if (now() - lastInputAt >= INPUT_IDLE_MS) {
    closeEngagedSegment()
  }
}

function snapshotAndReset() {
  syncVisibility()
  closeEngagedSegment()
  closeVisibleSegment()
  const dwellMs = pageEnteredAt > 0 ? Math.round(Math.max(0, now() - pageEnteredAt)) : 0
  const visibleMs = Math.round(Math.max(0, visibleAccumMs))
  const engagedMs = Math.round(Math.round(Math.max(0, engagedAccumMs)))
  visibleAccumMs = 0
  engagedAccumMs = 0
  if (foregroundActive()) {
    visible = true
    visibleSegmentStart = now()
  } else {
    visible = false
    visibleSegmentStart = 0
  }
  engaged = false
  engagedSegmentStart = 0
  pageEnteredAt = now()
  return { dwellMs, visibleMs, engagedMs }
}

async function flushNavigate(path, metrics) {
  const auth = useAuth()
  if (!auth.authenticated.value) return
  try {
    await postActivity({
      type: 'navigate',
      path,
      dwellMs: metrics.dwellMs,
      visibleMs: metrics.visibleMs,
      engagedMs: metrics.engagedMs
    })
  } catch {
    // Tracking must never break navigation.
  }
}

function onRouteChange(toPath) {
  const next = toPath || '/'
  if (currentPath && currentPath !== next) {
    const metrics = snapshotAndReset()
    void flushNavigate(currentPath, metrics)
  } else if (!currentPath) {
    pageEnteredAt = now()
    visibleAccumMs = 0
    engagedAccumMs = 0
    if (foregroundActive()) {
      visible = true
      visibleSegmentStart = now()
    }
  }
  currentPath = next
}

function onPageHide() {
  if (!currentPath) return
  const metrics = snapshotAndReset()
  const auth = useAuth()
  if (!auth.authenticated.value) return
  const body = JSON.stringify({
    type: 'navigate',
    path: currentPath,
    dwellMs: metrics.dwellMs,
    visibleMs: metrics.visibleMs,
    engagedMs: metrics.engagedMs
  })
  try {
    if (navigator.sendBeacon) {
      const blob = new Blob([body], { type: 'application/json' })
      navigator.sendBeacon(`${SURFING_HOST}/activity`, blob)
    } else {
      void flushNavigate(currentPath, metrics)
    }
  } catch {
    /* ignore */
  }
}

/**
 * Install page dwell / engagement tracking for authenticated users (site-wide).
 * View of the Activity panel remains dual-gated to dasm.
 */
export function installActivityTracker(router) {
  if (started || !router) return
  started = true

  if (typeof window !== 'undefined') {
    window.addEventListener('focus', syncVisibility)
    window.addEventListener('blur', syncVisibility)
    document.addEventListener('visibilitychange', syncVisibility)
    ;['mousemove', 'scroll', 'keydown', 'touchstart', 'click'].forEach((evt) => {
      window.addEventListener(evt, onInput, { passive: true })
    })
    window.addEventListener('pagehide', onPageHide)
    setInterval(tickEngagementIdle, 1000)
  }

  router.afterEach((to) => {
    onRouteChange(to.fullPath || to.path)
  })
}
