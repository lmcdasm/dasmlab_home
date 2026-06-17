<template>
  <q-page padding class="q-gutter-md home-page">

    <!-- Approach: Nav TL;DR — compact summary at top, then explore -->
    <template v-if="isNavTldr">
      <div class="hero-shell q-mb-lg">
        <div class="hero-noise" />
        <div class="hero-orb hero-orb-a" />
        <div class="hero-orb hero-orb-b" />
        <div class="hero-content">
          <div class="caps-label">TL;DR</div>
          <h1 class="hero-title q-mt-sm q-mb-sm">
            Living lab.
            <span class="accent-word">Real builds.</span>
            Shared craft.
          </h1>
          <p class="hero-copy q-mb-md">
            DASMLAB is where ideas get poked, prodded, and brought to life.
            Build notes, experiments, and production lessons in one place.
          </p>
          <div class="hero-chips">
            <q-chip dense color="primary" text-color="white" icon="science" label="Experiment" />
            <q-chip dense color="secondary" text-color="white" icon="construction" label="Build" />
            <q-chip dense color="accent" text-color="white" icon="school" label="Teach back" />
          </div>
        </div>
      </div>
      <div class="scanline q-pa-sm q-mb-md">
        <q-icon name="south" class="q-mr-xs" />
        Pick a project lane below
      </div>
    </template>

    <!-- Approach: Hero + block — WhatsNew + Carousel + cards -->
    <template v-else>
      <div class="hero-shell q-mb-lg">
        <div class="hero-noise" />
        <div class="hero-orb hero-orb-a" />
        <div class="hero-orb hero-orb-b" />
        <div class="hero-content">
          <div class="caps-label">DASMLAB</div>
          <h1 class="hero-title q-mt-sm q-mb-sm">
            Craft over templates.
            <span class="accent-word">Systems with personality.</span>
          </h1>
          <p class="hero-copy q-mb-md">
            Portfolio plus workshop: production projects, field notes, and practical demos from across the stack.
          </p>
          <div class="hero-chips">
            <q-chip dense color="primary" text-color="white" icon="palette" label="Design systems" />
            <q-chip dense color="secondary" text-color="white" icon="terminal" label="Engineering" />
            <q-chip dense color="accent" text-color="white" icon="hub" label="Infra + AI" />
          </div>
        </div>
      </div>

      <!-- What's New? -->
      <div class="q-pa-sm flex flex-center content-frame q-mb-md">
        <div class="frame-label">Live signal</div>
        <WhatsNew :news="whatsNew" />
      </div>

      <!-- Carousel Section -->
      <div class="flex flex-center q-my-md content-frame">
        <div class="frame-label">Featured experiments</div>
        <DesignCarousel :entries="carouselEntries" />
      </div>
    </template>

    <!-- Project Cards (shared) -->
    <div class="q-gutter-md row items-start justify-center q-mt-sm">
      <q-card
        v-for="card in cards"
        :key="card.title"
        class="q-pa-lg cursor-pointer col-xs-12 col-sm-4 col-md-3 project-card animate__animated animate__fadeInUp"
        flat
        bordered
        @click="goTo(card.route)"
      >
        <q-card-section class="column items-center justify-center card-inner">
          <div class="icon-wrap q-mb-sm">
            <q-icon :name="card.icon" size="34px" class="text-primary" />
          </div>
          <div class="text-subtitle1 text-center card-title">{{ card.title }}</div>
          <div class="card-cta q-mt-xs">Explore</div>
        </q-card-section>
      </q-card>
    </div>

  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import WhatsNew from 'src/components/WhatsNew.vue'
import DesignCarousel from 'src/components/DesignCarousel.vue'
import { useApproach } from 'src/composables/useApproach'
import axios from 'axios'

const router = useRouter()
const { isNavTldr } = useApproach()

const carouselEntries = ref([])
const whatsNew = ref([])

const goTo = (route) => router.push(route)

const cards = [
  { title: 'Frontend Projects', icon: 'view_quilt', route: '/projects/frontend' },
  { title: 'Backend Projects', icon: 'dns', route: '/projects/backend' },
  { title: 'AI/ML Tech Dives', icon: 'psychology', route: '/projects/ai-ml' },
  { title: 'Cloud Provider Techs.', icon: 'cloud', route: '/projects/cloud' },
  { title: 'Infrastructure Projects', icon: 'storage', route: '/projects/infrastructure' },
  { title: 'Security Projects', icon: 'shield', route: '/projects/security' }
]

onMounted(async () => {
  // Fetch WhatsNew
  try {
    const res = await axios.get('https://whatsnew.svc.dasmlab.org/get')
    whatsNew.value = res.data?.latest_commits || []
  } catch (err) {
    console.error('Failed to Fetch from whatsnew-service:', err)
    whatsNew.value = [
      {
        id: 1,
        project: 'Out of Service',
        title: 'Connection to WhatsNew Service is down.',
        date: new Date().toISOString().replace('T', ' ').substring(0, 19)
      }
    ]
  }

  // Fetch Carousel
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
}

.hero-shell {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(43, 76, 103, 0.16);
  border-radius: 18px;
  background:
    radial-gradient(circle at 12% 20%, rgba(63, 122, 107, 0.14), transparent 46%),
    radial-gradient(circle at 82% 30%, rgba(158, 115, 178, 0.12), transparent 48%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.98), rgba(246, 251, 255, 0.98));
  box-shadow: 0 14px 28px rgba(25, 47, 67, 0.08);
}

.hero-content {
  position: relative;
  z-index: 2;
  padding: 1.35rem 1.2rem 1.25rem;
}

.hero-noise {
  position: absolute;
  inset: 0;
  z-index: 1;
  opacity: 0.1;
  background-image: radial-gradient(rgba(44, 66, 86, 0.25) 0.6px, transparent 0.6px);
  background-size: 4px 4px;
}

.hero-orb {
  position: absolute;
  z-index: 1;
  width: 180px;
  height: 180px;
  border-radius: 50%;
  filter: blur(20px);
  opacity: 0.34;
  animation: float-orb 7s ease-in-out infinite;
}

.hero-orb-a {
  background: rgba(63, 122, 107, 0.35);
  top: -60px;
  right: -40px;
}

.hero-orb-b {
  background: rgba(158, 115, 178, 0.3);
  left: -50px;
  bottom: -70px;
  animation-delay: 1.5s;
}

.caps-label {
  display: inline-block;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  font-size: 0.73rem;
  color: #607587;
  font-weight: 600;
}

.hero-title {
  max-width: 700px;
  line-height: 1.08;
  font-size: clamp(1.4rem, 2.7vw, 2.2rem);
  font-weight: 700;
  color: #1f2f3e;
}

.accent-word {
  display: block;
  color: #885fa0;
  text-shadow: 0 0 12px rgba(158, 115, 178, 0.22);
  animation: accent-pulse 3.2s ease-in-out infinite;
}

.hero-copy {
  max-width: 780px;
  color: #445769;
}

.hero-chips :deep(.q-chip) {
  margin-right: 0.45rem;
  margin-bottom: 0.35rem;
}

.scanline {
  border: 1px dashed rgba(90, 117, 139, 0.35);
  border-radius: 10px;
  text-align: center;
  color: #607587;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  font-size: 0.77rem;
  background: rgba(255, 255, 255, 0.8);
}

.content-frame {
  position: relative;
  width: 100%;
  overflow: hidden;
  border: 1px solid rgba(43, 76, 103, 0.14);
  border-radius: 14px;
  background: linear-gradient(170deg, rgba(255, 255, 255, 0.96), rgba(247, 251, 254, 0.98));
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.55), 0 8px 20px rgba(25, 47, 67, 0.06);
}

.content-frame::after {
  content: '';
  position: absolute;
  inset: 0 auto 0 0;
  width: 4px;
  background: linear-gradient(180deg, rgba(63, 122, 107, 0.64), rgba(158, 115, 178, 0.55));
  pointer-events: none;
}

.content-frame::before {
  display: none;
}

.frame-label {
  position: absolute;
  top: 8px;
  right: 10px;
  font-size: 0.68rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: rgba(84, 106, 124, 0.9);
}

.project-card {
  position: relative;
  border: 1px solid rgba(43, 76, 103, 0.15) !important;
  border-radius: 14px;
  background: linear-gradient(160deg, #ffffff, #f7fbfe);
  transition: transform 180ms ease, border-color 180ms ease, box-shadow 180ms ease;
}

.project-card:hover {
  transform: translateY(-3px);
  border-color: rgba(151, 110, 176, 0.5) !important;
  box-shadow: 0 12px 24px rgba(25, 47, 67, 0.12);
}

.card-inner {
  min-height: 128px;
}

.icon-wrap {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: radial-gradient(circle at 30% 25%, rgba(63, 122, 107, 0.18), rgba(255, 255, 255, 0.9));
  border: 1px solid rgba(63, 122, 107, 0.25);
}

.card-title {
  font-weight: 600;
  color: #24384a;
  letter-spacing: 0.01em;
}

.card-cta {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: #607587;
}

@media (max-width: 900px) {
  .project-card {
    width: 100%;
  }
}

@keyframes float-orb {
  0%,
  100% { transform: translate3d(0, 0, 0) scale(1); }
  50% { transform: translate3d(0, 9px, 0) scale(1.06); }
}

@keyframes accent-pulse {
  0%,
  100% { opacity: 0.86; }
  50% { opacity: 1; }
}
</style>

