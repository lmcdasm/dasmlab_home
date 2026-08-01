<template>
  <q-page class="home-page">
    <div class="home-shell">

    <!-- Approach: Nav TL;DR — compact summary at top, then explore -->
    <template v-if="isNavTldr">
      <div class="hero-shell rise" style="--rise-delay: 0ms">
        <div class="hero-noise" />
        <div class="hero-orb hero-orb-a" />
        <div class="hero-orb hero-orb-b" />
        <div class="hero-content">
          <div class="caps-label">TL;DR</div>
          <h1 class="hero-title">
            Living lab.
            <span class="accent-word">Real builds.</span>
            Shared craft.
          </h1>
          <p class="hero-copy">
            DASMLAB is where ideas get poked, prodded, and brought to life.
            Build notes, experiments, and production lessons in one place.
          </p>
          <div class="hero-cta-row">
            <button type="button" class="hero-btn primary" @click="goTo('/surfing')">Surfing</button>
            <button type="button" class="hero-btn ghost" @click="goTo('/projects/infrastructure')">Infra</button>
            <button type="button" class="hero-btn ghost" @click="goTo('/contact')">Contact</button>
          </div>
        </div>
      </div>
      <div class="scanline">
        <q-icon name="south" class="q-mr-xs" />
        Pick a project lane below
      </div>
    </template>

    <!-- Approach: Hero + lab map + architecture -->
    <template v-else>
      <div class="hero-shell rise" style="--rise-delay: 0ms">
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
            Portfolio plus workshop: production projects, field notes, and practical demos from across the stack.
          </p>
          <div class="hero-cta-row">
            <button type="button" class="hero-btn primary" @click="goTo('/surfing')">Surfing</button>
            <button type="button" class="hero-btn ghost" @click="goTo('/projects/frontend')">Projects</button>
            <button type="button" class="hero-btn ghost" @click="goTo('/contact')">Contact</button>
          </div>
        </div>
      </div>

      <section class="section-block rise" style="--rise-delay: 80ms">
        <div class="section-head">
          <div class="caps-label">Lab map</div>
          <h2 class="section-title">Click a lane for builds, teaches, and a deep link.</h2>
        </div>
        <LabMap />
      </section>

      <section class="section-block rise" style="--rise-delay: 140ms">
        <div class="section-head section-head--center">
          <div class="caps-label">Project lanes</div>
          <h2 class="section-title">Dive by category.</h2>
        </div>
        <div class="lane-row">
          <button
            v-for="card in cards"
            :key="card.title"
            type="button"
            class="lane-chip"
            @click="goTo(card.route)"
          >
            <q-icon :name="card.icon" size="20px" />
            <span>{{ card.title }}</span>
          </button>
        </div>
      </section>

      <section class="section-block rise" style="--rise-delay: 180ms">
        <div class="section-head">
          <div class="caps-label">How the lab is wired</div>
          <h2 class="section-title">Edge → OCP → GitOps → media off basement disk.</h2>
        </div>
        <LabArchitecture />
      </section>

      <section class="section-block rise" style="--rise-delay: 220ms">
        <div class="section-head">
          <div class="caps-label">Featured experiments</div>
          <h2 class="section-title">Live carousel feed.</h2>
        </div>
        <div class="content-frame">
          <DesignCarousel :entries="carouselEntries" />
        </div>
      </section>
    </template>

    <!-- TL;DR approach: lanes still available below -->
    <section v-if="isNavTldr" class="section-block rise" style="--rise-delay: 120ms">
      <div class="section-head section-head--center">
        <div class="caps-label">Project lanes</div>
        <h2 class="section-title">Dive by category.</h2>
      </div>
      <div class="lane-row">
        <button
          v-for="card in cards"
          :key="card.title"
          type="button"
          class="lane-chip"
          @click="goTo(card.route)"
        >
          <q-icon :name="card.icon" size="20px" />
          <span>{{ card.title }}</span>
        </button>
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
import { useApproach } from 'src/composables/useApproach'
import axios from 'axios'

const router = useRouter()
const { isNavTldr } = useApproach()

const carouselEntries = ref([])

const goTo = (route) => router.push(route)

const cards = [
  { title: 'Frontend', icon: 'view_quilt', route: '/projects/frontend' },
  { title: 'Backend', icon: 'dns', route: '/projects/backend' },
  { title: 'AI / ML', icon: 'psychology', route: '/projects/ai-ml' },
  { title: 'Cloud', icon: 'cloud', route: '/projects/cloud' },
  { title: 'Infra', icon: 'storage', route: '/projects/infrastructure' },
  { title: 'Security', icon: 'shield', route: '/projects/security' }
]

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
  padding: 0.7rem clamp(0.65rem, 1.6vw, 1.1rem) 1.4rem;
}

.home-shell {
  width: 100%;
  max-width: 920px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

@media (min-width: 1400px) {
  .home-shell {
    max-width: 980px;
  }
}

.rise {
  animation: riseIn 520ms ease both;
  animation-delay: var(--rise-delay, 0ms);
}

.hero-shell {
  position: relative;
  overflow: hidden;
  width: 100%;
  border: 1.5px solid rgba(31, 111, 98, 0.32);
  border-radius: 20px;
  background:
    radial-gradient(circle at 12% 20%, rgba(47, 143, 125, 0.16), transparent 46%),
    radial-gradient(circle at 82% 30%, rgba(151, 110, 176, 0.14), transparent 48%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.99), rgba(236, 245, 242, 0.98));
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.85),
    0 16px 34px rgba(18, 40, 52, 0.1);
}

.hero-content {
  position: relative;
  z-index: 2;
  padding: 1rem 1.15rem 1.05rem;
  max-width: 36rem;
}

.hero-title {
  margin: 0 0 0.45rem;
  line-height: 1.08;
  font-size: clamp(1.25rem, 2vw, 1.7rem);
  font-weight: 800;
  color: #12202c;
  letter-spacing: -0.01em;
}

.hero-copy {
  margin: 0 0 0.75rem;
  max-width: 34rem;
  color: #3d5263;
  font-size: 0.92rem;
  line-height: 1.45;
}

.hero-orb {
  position: absolute;
  z-index: 1;
  width: 140px;
  height: 140px;
  border-radius: 50%;
  filter: blur(18px);
  opacity: 0.32;
  animation: float-orb 7s ease-in-out infinite;
}

.hero-noise {
  position: absolute;
  inset: 0;
  z-index: 1;
  opacity: 0.12;
  background-image: radial-gradient(rgba(29, 43, 54, 0.28) 0.6px, transparent 0.6px);
  background-size: 4px 4px;
}

.hero-orb-a {
  background: rgba(47, 143, 125, 0.42);
  top: -50px;
  right: -36px;
}

.hero-orb-b {
  background: rgba(151, 110, 176, 0.34);
  left: -40px;
  bottom: -56px;
  animation-delay: 1.5s;
}

.caps-label {
  display: inline-block;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  font-size: 0.68rem;
  color: #4d6575;
  font-weight: 700;
  margin-bottom: 0.25rem;
}

.accent-word {
  display: block;
  color: #7a4f96;
  text-shadow: 0 0 14px rgba(151, 110, 176, 0.25);
  animation: accent-pulse 3.2s ease-in-out infinite;
}

.hero-cta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
}

.hero-btn {
  border-radius: 10px;
  padding: 0.48rem 0.85rem;
  font-weight: 720;
  font-size: 0.9rem;
  letter-spacing: 0.02em;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease, background 0.15s ease;
}

.hero-btn.primary {
  border: none;
  background: linear-gradient(135deg, #1f6f62, #2f8f7d);
  color: #fff;
  box-shadow: 0 10px 22px rgba(31, 111, 98, 0.28);
}

.hero-btn.ghost {
  border: 1.5px solid rgba(31, 111, 98, 0.35);
  background: rgba(255, 255, 255, 0.75);
  color: #1d2b36;
}

.hero-btn:hover {
  transform: translateY(-1px);
}

.scanline {
  border: 1px dashed rgba(90, 117, 139, 0.4);
  border-radius: 10px;
  text-align: center;
  color: #4d6575;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  font-size: 0.77rem;
  background: rgba(255, 255, 255, 0.85);
  padding: 0.55rem 0.75rem;
}

.section-block {
  width: 100%;
}

.section-head {
  margin-bottom: 0.7rem;
}

.section-head--center {
  text-align: center;
}

.section-title {
  margin: 0.25rem 0 0;
  font-size: clamp(1.05rem, 1.8vw, 1.35rem);
  font-weight: 750;
  color: #12202c;
  line-height: 1.25;
}

.content-frame {
  position: relative;
  width: 100%;
  overflow: hidden;
  border: 1.5px solid rgba(31, 111, 98, 0.22);
  border-radius: 16px;
  background: linear-gradient(170deg, rgba(255, 255, 255, 0.98), rgba(243, 248, 246, 0.98));
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.55), 0 10px 24px rgba(18, 40, 52, 0.07);
  padding: 0.35rem;
}

.lane-row {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.65rem;
  width: 100%;
  padding: 0.35rem 0 0.15rem;
}

.lane-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  border: 1.5px solid rgba(31, 111, 98, 0.28);
  border-radius: 999px;
  padding: 0.58rem 1.05rem;
  background: linear-gradient(160deg, #ffffff, #f3f8f6);
  color: #12202c;
  font-weight: 650;
  cursor: pointer;
  transition: transform 0.15s ease, border-color 0.15s ease, box-shadow 0.15s ease;
}

.lane-chip:hover {
  transform: translateY(-2px);
  border-color: rgba(151, 110, 176, 0.55);
  box-shadow: 0 10px 20px rgba(18, 40, 52, 0.1);
}

@keyframes float-orb {
  0%,
  100% { transform: translate3d(0, 0, 0) scale(1); }
  50% { transform: translate3d(0, 9px, 0) scale(1.06); }
}

@keyframes accent-pulse {
  0%,
  100% { opacity: 0.88; }
  50% { opacity: 1; }
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
</style>
