import { postActivity } from 'src/services/authApi'

const INPUT_IDLE_MS = 5000
const SURFING_HOST = import.meta.env.VITE_SURFING_API_HOST || '/api/surfing'
const COOKIE_AID = 'surf_aid'
const COOKIE_SID = 'surf_sid'

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
let scrollMaxPct = 0
let landingReferrer = ''
let landingUTM = { utmSource: '', utmMedium: '', utmCampaign: '' }

function now() {
  return typeof performance !== 'undefined' && performance.now ? performance.now() : Date.now()
}

function readCookie(name) {
  if (typeof document === 'undefined') return ''
  const m = document.cookie.match(new RegExp('(?:^|; )' + name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + '=([^;]*)'))
  return m ? decodeURIComponent(m[1]) : ''
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

function updateScrollMax() {
  if (typeof window === 'undefined' || typeof document === 'undefined') return
  const doc = document.documentElement
  const max = doc.scrollHeight - doc.clientHeight
  if (max <= 0) {
    scrollMaxPct = Math.max(scrollMaxPct, 100)
    return
  }
  const pct = Math.min(100, Math.round((window.scrollY / max) * 100))
  if (pct > scrollMaxPct) scrollMaxPct = pct
}

function snapshotAndReset() {
  syncVisibility()
  closeEngagedSegment()
  closeVisibleSegment()
  updateScrollMax()
  const dwellMs = pageEnteredAt > 0 ? Math.round(Math.max(0, now() - pageEnteredAt)) : 0
  const visibleMs = Math.round(Math.max(0, visibleAccumMs))
  const engagedMs = Math.round(Math.max(0, engagedAccumMs))
  const scroll = scrollMaxPct
  visibleAccumMs = 0
  engagedAccumMs = 0
  scrollMaxPct = 0
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
  return { dwellMs, visibleMs, engagedMs, scrollMaxPct: scroll }
}

function campaignFromSearch(search) {
  const q = new URLSearchParams(search || '')
  return {
    utmSource: q.get('utm_source') || '',
    utmMedium: q.get('utm_medium') || '',
    utmCampaign: q.get('utm_campaign') || ''
  }
}

function buildPayload(path, metrics) {
  return {
    type: 'page',
    path,
    title: typeof document !== 'undefined' ? document.title || '' : '',
    referrer: landingReferrer,
    utmSource: landingUTM.utmSource,
    utmMedium: landingUTM.utmMedium,
    utmCampaign: landingUTM.utmCampaign,
    dwellMs: metrics.dwellMs,
    visibleMs: metrics.visibleMs,
    engagedMs: metrics.engagedMs,
    scrollMaxPct: metrics.scrollMaxPct,
    anonymousId: readCookie(COOKIE_AID),
    sessionId: readCookie(COOKIE_SID)
  }
}

async function flushPage(path, metrics) {
  try {
    await postActivity(buildPayload(path, metrics))
  } catch {
    // Tracking must never break navigation.
  }
}

function onRouteChange(toPath, toFullPath) {
  const next = toPath || '/'
  if (currentPath && currentPath !== next) {
    const metrics = snapshotAndReset()
    void flushPage(currentPath, metrics)
  } else if (!currentPath) {
    pageEnteredAt = now()
    visibleAccumMs = 0
    engagedAccumMs = 0
    scrollMaxPct = 0
    if (foregroundActive()) {
      visible = true
      visibleSegmentStart = now()
    }
    // First paint: emit a page view with zero dwell so anon traffic appears immediately.
    void flushPage(next, { dwellMs: 0, visibleMs: 0, engagedMs: 0, scrollMaxPct: 0 })
  }
  currentPath = next
  if (toFullPath && toFullPath.includes('?')) {
    const utm = campaignFromSearch(toFullPath.slice(toFullPath.indexOf('?')))
    if (utm.utmSource || utm.utmMedium || utm.utmCampaign) {
      landingUTM = utm
    }
  }
}

function onPageHide() {
  if (!currentPath) return
  const metrics = snapshotAndReset()
  const body = JSON.stringify(buildPayload(currentPath, metrics))
  try {
    if (navigator.sendBeacon) {
      const blob = new Blob([body], { type: 'application/json' })
      navigator.sendBeacon(`${SURFING_HOST}/activity`, blob)
    } else {
      void flushPage(currentPath, metrics)
    }
  } catch {
    /* ignore */
  }
}

/**
 * Site-wide first-party engagement collector (anonymous + authenticated).
 * Activity Panel view remains dual-gated to dasm.
 */
export function installActivityTracker(router) {
  if (started || !router) return
  started = true

  if (typeof window !== 'undefined') {
    landingReferrer = document.referrer || ''
    landingUTM = campaignFromSearch(window.location.search)
    window.addEventListener('focus', syncVisibility)
    window.addEventListener('blur', syncVisibility)
    document.addEventListener('visibilitychange', syncVisibility)
    ;['mousemove', 'scroll', 'keydown', 'touchstart', 'click'].forEach((evt) => {
      window.addEventListener(evt, onInput, { passive: true })
    })
    window.addEventListener('scroll', updateScrollMax, { passive: true })
    window.addEventListener('pagehide', onPageHide)
    setInterval(tickEngagementIdle, 1000)
  }

  router.afterEach((to) => {
    onRouteChange(to.path || '/', to.fullPath || to.path)
  })
}
