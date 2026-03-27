/**
 * Approach switcher — live-switch between presentation "approaches" (e.g. hero+block vs nav-level TL;DR).
 * Persisted in localStorage so the choice survives refresh. Reusable for future A/B-style layout toggles.
 *
 * @see docs/TARGET-0.6.0.md — Approach switcher
 * @see docs/VOICE-AND-PRIORITIES.md
 */

import { ref, watch, computed } from 'vue'

const STORAGE_KEY = 'dasmlab-home-approach'

export const APPROACH_HERO_BLOCK = 'hero-block'
export const APPROACH_NAV_TLDR = 'nav-tldr'

const approachOptions = [
  { value: APPROACH_HERO_BLOCK, label: 'Hero + block' },
  { value: APPROACH_NAV_TLDR, label: 'Nav TL;DR' }
]

function readStored () {
  try {
    const v = localStorage.getItem(STORAGE_KEY)
    return v === APPROACH_NAV_TLDR ? APPROACH_NAV_TLDR : APPROACH_HERO_BLOCK
  } catch {
    return APPROACH_HERO_BLOCK
  }
}

const approach = ref(readStored())

watch(approach, (val) => {
  try {
    localStorage.setItem(STORAGE_KEY, val)
  } catch (_) {}
}, { immediate: false })

export function useApproach () {
  function setApproach (value) {
    approach.value = value === APPROACH_NAV_TLDR ? APPROACH_NAV_TLDR : APPROACH_HERO_BLOCK
  }

  const isHeroBlock = computed(() => approach.value === APPROACH_HERO_BLOCK)
  const isNavTldr = computed(() => approach.value === APPROACH_NAV_TLDR)

  return {
    approach,
    setApproach,
    approachOptions,
    isHeroBlock,
    isNavTldr
  }
}
