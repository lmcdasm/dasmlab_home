<template>
  <q-page class="launch-page">
    <section class="launch-hero rise" style="--rise-delay: 0ms">
      <div class="launch-noise" />
      <div class="launch-orb launch-orb-a" />
      <div class="launch-orb launch-orb-b" />
      <div class="launch-inner">
        <div class="caps-label">DASMLAB</div>
        <h1 class="launch-title">
          2.0 is live.
          <span class="accent-word">Engineering Knowledge Network.</span>
        </h1>
        <p class="launch-copy">
          Answers for visitors. How-we-built for engineers. Observed by the Digital Presence Observatory we just shipped.
        </p>
        <div class="launch-cta">
          <button type="button" class="hero-btn primary" @click="goTo('/')">Explore the lab</button>
          <a
            class="hero-btn ghost"
            href="https://dpo-dasmlab.apps.2026-prod-1.ocp.dasmlab.org"
            rel="noopener"
          >Observatory</a>
          <button type="button" class="hero-btn ghost" @click="goTo('/contact')">Contact</button>
        </div>
        <p class="launch-meta">
          Campaign <code>dasmlab-2.0-launch</code> · home ADR-002 · observatory ADR-0402
        </p>
      </div>
    </section>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSeo } from 'src/composables/useSeo'

const router = useRouter()
const route = useRoute()

useSeo({
  title: 'DASMLAB 2.0 — Engineering Knowledge Network',
  description:
    'DASMLAB 2.0 launch: Engineering Knowledge Network with Digital Presence Observatory — answers for visitors, how-we-built for engineers.',
  path: '/launch'
})

const goTo = (path) => router.push(path)

onMounted(() => {
  // Ensure default campaign UTM is present for Activity if none provided.
  if (!route.query.utm_campaign) {
    const q = {
      ...route.query,
      utm_source: route.query.utm_source || 'direct',
      utm_medium: route.query.utm_medium || 'campaign',
      utm_campaign: 'dasmlab-2.0-launch'
    }
    router.replace({ path: '/launch', query: q })
  }
})
</script>

<style scoped>
.launch-page {
  padding: 0;
  min-height: calc(100vh - 64px);
}

.rise {
  animation: riseIn 560ms ease both;
  animation-delay: var(--rise-delay, 0ms);
}

@keyframes riseIn {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: none; }
}

.launch-hero {
  position: relative;
  overflow: hidden;
  min-height: calc(100vh - 64px);
  display: flex;
  align-items: center;
  background:
    radial-gradient(circle at 18% 20%, rgba(47, 143, 125, 0.22), transparent 45%),
    radial-gradient(circle at 88% 10%, rgba(18, 64, 78, 0.18), transparent 40%),
    linear-gradient(165deg, #0f1c24 0%, #163038 48%, #1f4a44 100%);
  color: #eef6f3;
}

.launch-noise {
  position: absolute;
  inset: 0;
  opacity: 0.08;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
  pointer-events: none;
}

.launch-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(40px);
  pointer-events: none;
}
.launch-orb-a {
  width: 280px; height: 280px; left: -40px; top: 10%;
  background: rgba(47, 143, 125, 0.35);
  animation: floatA 9s ease-in-out infinite;
}
.launch-orb-b {
  width: 220px; height: 220px; right: 8%; bottom: 12%;
  background: rgba(120, 180, 170, 0.22);
  animation: floatB 11s ease-in-out infinite;
}
@keyframes floatA {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(18px, -22px); }
}
@keyframes floatB {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-14px, 16px); }
}

.launch-inner {
  position: relative;
  z-index: 2;
  width: min(720px, 92vw);
  margin: 0 auto;
  padding: 2.5rem 1rem 3rem;
}

.caps-label {
  font-size: 0.72rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: rgba(190, 230, 220, 0.9);
  font-weight: 600;
  margin-bottom: 0.75rem;
}

.launch-title {
  margin: 0 0 0.85rem;
  font-size: clamp(2.1rem, 5vw, 3.4rem);
  line-height: 1.05;
  letter-spacing: -0.02em;
  font-weight: 700;
}

.accent-word {
  display: block;
  color: #8fd4c4;
  margin-top: 0.2rem;
}

.launch-copy {
  margin: 0 0 1.5rem;
  font-size: clamp(1.05rem, 2vw, 1.25rem);
  line-height: 1.45;
  max-width: 36rem;
  color: rgba(230, 242, 238, 0.92);
}

.launch-cta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem;
  margin-bottom: 1.5rem;
}

.hero-btn {
  appearance: none;
  border: 1px solid transparent;
  border-radius: 999px;
  padding: 0.65rem 1.15rem;
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
}
.hero-btn.primary {
  background: #2f8f7d;
  color: #fff;
}
.hero-btn.ghost {
  background: transparent;
  border-color: rgba(200, 230, 220, 0.45);
  color: #eef6f3;
}

.launch-meta {
  font-size: 0.8rem;
  color: rgba(190, 220, 210, 0.75);
}
.launch-meta code {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.78rem;
}
</style>
