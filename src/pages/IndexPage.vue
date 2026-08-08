<template>
  <q-page class="home-page">
    <div class="home-shell">
      <!-- One composition: brand strip + Lab Map dominate the first viewport -->
      <section class="hero-map rise" style="--rise-delay: 0ms">
        <div class="hero-strip">
          <div class="hero-noise" />
          <div class="hero-orb hero-orb-a" />
          <div class="hero-orb hero-orb-b" />
          <div class="hero-content">
            <div class="caps-label">DASMLAB</div>
            <h1 class="hero-title">
              Craft over templates.
              <span class="accent-word">Systems with personality.</span>
            </h1>
            <p class="hero-copy">
              Engineering Knowledge Network — answers for visitors, how-we-built for engineers.
            </p>
            <div class="hero-cta-row">
              <button type="button" class="hero-btn primary" @click="goTo('/surfing')">Surfing</button>
              <button type="button" class="hero-btn ghost" @click="goTo('/projects/frontend')">Projects</button>
              <button type="button" class="hero-btn ghost" @click="goTo('/contact')">Contact</button>
            </div>
          </div>
        </div>

        <div class="map-block">
          <div class="map-kicker">
            <span class="caps-label">Lab map</span>
            <span class="map-hint">Click a lane — detail opens beside the map</span>
          </div>
          <LabMap />
        </div>
      </section>

      <div class="flow-rail" aria-hidden="true">
        <span>Explore</span>
        <span class="flow-dot" />
        <span>Wire</span>
        <span class="flow-dot" />
        <span>Experiments</span>
      </div>

      <section class="section-block rise" style="--rise-delay: 100ms">
        <div class="section-head">
          <div class="caps-label">How the lab is wired</div>
          <h2 class="section-title">Edge → OpenShift → GitOps → media off basement disk.</h2>
        </div>
        <LabArchitecture />
      </section>

      <section class="section-block rise" style="--rise-delay: 160ms">
        <div class="section-head">
          <div class="caps-label">Featured experiments</div>
          <h2 class="section-title">Live carousel feed from the lab.</h2>
        </div>
        <div class="content-frame">
          <DesignCarousel :entries="carouselEntries" />
        </div>
      </section>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DesignCarousel from 'src/components/DesignCarousel.vue'
import LabMap from 'src/components/LabMap.vue'
import LabArchitecture from 'src/components/LabArchitecture.vue'
import { useSeo } from 'src/composables/useSeo'
import axios from 'axios'

const router = useRouter()
useSeo({
  title: 'Engineering Knowledge Network',
  description:
    'Technologies DASMLAB Inc. living lab — cloud-native, AI, and infrastructure builds with how-we-built-it depth.',
  path: '/'
})

const carouselEntries = ref([])
const goTo = (route) => router.push(route)

onMounted(async () => {
  try {
    const res = await axios.get('https://design-carousel.svc.dasmlab.org/carousel')
    carouselEntries.value = res.data
  } catch (err) {
    console.error('Failed to Fetch from design-carousel-service:', err)
    carouselEntries.value = [
      {
        id: 1,
        image_url: '/default.png',
        title: 'Service Offline',
        created_at: new Date().toISOString()
      }
    ]
  }
})
</script>

<style scoped>
.home-page {
  position: relative;
  overflow-x: hidden;
  padding: 0.45rem clamp(0.55rem, 1.4vw, 1rem) 1.5rem;
}

.home-shell {
  width: 100%;
  max-width: 1180px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.rise {
  animation: riseIn 520ms ease both;
  animation-delay: var(--rise-delay, 0ms);
}

.hero-map {
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
  border: 1.5px solid rgba(31, 111, 98, 0.34);
  border-radius: 20px;
  overflow: hidden;
  background:
    radial-gradient(circle at 12% 8%, rgba(47, 143, 125, 0.14), transparent 42%),
    linear-gradient(165deg, #ffffff, #eef5f2);
  box-shadow: 0 16px 34px rgba(18, 40, 52, 0.1);
  min-height: min(78vh, 760px);
}

.hero-strip {
  position: relative;
  padding: 0.75rem 1rem 0.55rem;
  border-bottom: 1px solid rgba(31, 111, 98, 0.14);
}

.hero-content {
  position: relative;
  z-index: 2;
  max-width: 40rem;
}

.hero-title {
  margin: 0 0 0.3rem;
  line-height: 1.08;
  font-size: clamp(1.35rem, 2.4vw, 1.95rem);
  font-weight: 700;
  color: #12202c;
  letter-spacing: -0.015em;
}

.hero-copy {
  margin: 0 0 0.55rem;
  color: #3d5263;
  font-size: 0.95rem;
  line-height: 1.4;
}

.accent-word {
  display: inline;
  color: #1f6f62;
  text-shadow: none;
}

.hero-cta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
}

.hero-btn {
  border-radius: 10px;
  padding: 0.42rem 0.8rem;
  font-weight: 700;
  font-size: 0.88rem;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.hero-btn.primary {
  border: none;
  background: linear-gradient(135deg, #1f6f62, #2f8f7d);
  color: #fff;
  box-shadow: 0 8px 18px rgba(31, 111, 98, 0.28);
}

.hero-btn.ghost {
  border: 1.5px solid rgba(31, 111, 98, 0.35);
  background: rgba(255, 255, 255, 0.8);
  color: #1d2b36;
}

.hero-btn:hover {
  transform: translateY(-1px);
}

.hero-orb {
  position: absolute;
  z-index: 1;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  filter: blur(18px);
  opacity: 0.28;
}

.hero-noise {
  position: absolute;
  inset: 0;
  z-index: 1;
  opacity: 0.1;
  background-image: radial-gradient(rgba(29, 43, 54, 0.28) 0.6px, transparent 0.6px);
  background-size: 4px 4px;
}

.hero-orb-a {
  background: rgba(47, 143, 125, 0.45);
  top: -40px;
  right: -20px;
}

.hero-orb-b {
  background: rgba(116, 134, 92, 0.4);
  left: -30px;
  bottom: -40px;
}

.map-block {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0 0.75rem 0.85rem;
  min-height: 0;
}

.map-kicker {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 0.35rem;
  padding: 0 0.15rem;
}

.map-hint {
  font-size: 0.78rem;
  color: #5a6f80;
}

.flow-rail {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.65rem;
  padding: 0.35rem 0.5rem;
  font-size: 0.72rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  font-weight: 700;
  color: #5a6f80;
}

.flow-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #1f6f62;
  opacity: 0.55;
}

.caps-label {
  display: inline-block;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  font-size: 0.68rem;
  color: #4d6575;
  font-weight: 700;
  margin-bottom: 0.15rem;
}

.section-block {
  width: 100%;
}

.section-head {
  margin-bottom: 0.55rem;
}

.section-title {
  margin: 0.2rem 0 0;
  font-size: clamp(1.05rem, 1.8vw, 1.35rem);
  font-weight: 700;
  color: #12202c;
  line-height: 1.25;
}

.content-frame {
  border: 1.5px solid rgba(31, 111, 98, 0.22);
  border-radius: 16px;
  background: linear-gradient(170deg, #fff, #f3f8f6);
  box-shadow: 0 10px 24px rgba(18, 40, 52, 0.07);
  padding: 0.35rem;
}

@keyframes riseIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 720px) {
  .hero-map {
    min-height: auto;
  }
}
</style>
