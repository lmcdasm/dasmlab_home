<template>
  <q-page class="about-page">
    <section class="about-hero">
      <div class="dasm-caps">About</div>
      <h1 class="hero-title">Living lab, lived experience</h1>
      <p class="hero-sub">
        DASMLAB is where technology gets poked, prodded, and brought to life. Every build is part experiment
        and part teach-back — demystifying the work so others can learn, remix, and iterate.
      </p>
    </section>

    <section class="about-grid">
      <div class="portrait-wrap">
        <img
          :src="meHead"
          class="about-head"
          width="420"
          height="525"
          alt="Daniel Smith (dasm), founder of Technologies DASMLAB Inc."
          @error="onMediaError"
        />
      </div>
      <div class="bio-panel">
        <h2>Who am I?</h2>
        <p>
          I’m Daniel Smith (dasm) — founder of Technologies DASMLAB Inc. I’ve spent two decades across telecom,
          automotive, cloud, and AI platforms. The throughline is practical systems thinking: shipping resilient
          solutions people can understand and extend.
        </p>
        <p>
          Work and life have taken me through <strong>60+ countries</strong>. From Outaouais/Ottawa to Fernie and
          now the Laurentians, the mission stays the same: make complex engineering feel teachable — and show the
          design authority behind each build.
        </p>
        <p class="q-mb-none">
          This site is the public Engineering Knowledge Network: answers for visitors, and “how we built it” for
          engineers — with labeled demos on private products so you can look around without Keycloak.
        </p>
      </div>
    </section>

    <div class="waypoint">Beyond the code</div>

    <section class="about-grid">
      <div class="carousel-wrap">
        <q-carousel
          v-model="carouselIndex"
          animated
          swipeable
          infinite
          autoplay
          transition-prev="slide-right"
          transition-next="slide-left"
          height="340px"
          class="about-carousel"
        >
          <q-carousel-slide
            v-for="(img, index) in lifestyleImages"
            :key="index"
            :name="index"
            class="column items-center justify-center"
          >
            <div class="q-pa-sm column items-center">
              <img
                :src="img.src"
                :alt="img.caption"
                width="560"
                height="300"
                class="about-lifestyle-img"
                @error="onMediaError"
              />
              <div class="caption-band">{{ img.caption }}</div>
            </div>
          </q-carousel-slide>
        </q-carousel>
      </div>
      <div class="bio-panel">
        <h2>Outside work</h2>
        <p>
          Snowboarding, windsurfing, gardening, house music, and pottery keep the creative loop open. They shape
          the rhythm of how I design systems: playful, grounded, and always iterative.
        </p>
        <p class="q-mb-none">
          When I’m not in a cluster or a PR, I’m usually chasing wind on the water, carving snow at Mont Ste. Anne,
          or tending the garden with Princess Leah nearby.
        </p>
      </div>
    </section>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useSeo } from 'src/composables/useSeo'
import { ABOUT_MEDIA_BASE } from 'src/data/hubs'

const carouselIndex = ref(0)

function media(name) {
  return `${ABOUT_MEDIA_BASE}/${name}`
}

function onMediaError(event) {
  const el = event.target
  const name = (el.getAttribute('src') || '').split('/').pop()
  if (!name || el.dataset.fallback === '1') return
  el.dataset.fallback = '1'
  el.src = `/media/hero/${name}`
}

const meHead = media('portrait.webp')
const lifestyleImages = [
  { src: media('dasm_ride.webp'), caption: 'Windsurfing on a summer low-wind day' },
  { src: media('me_ride_2.webp'), caption: 'Snowboarding at Mont Ste. Anne' },
  { src: media('me_home.webp'), caption: 'Country living focus' },
  { src: media('me_plant_1.webp'), caption: 'Garden therapy in action' },
  { src: media('me_leah_1.webp'), caption: 'Princess Leah herself' },
  { src: media('me_sail_1.webp'), caption: 'Current favorite sail setup' },
  { src: media('me_baking_1.webp'), caption: 'Weekend baking mode' },
  { src: media('me_baking_2.webp'), caption: 'And the reward at the end' }
]

useSeo({
  title: 'About',
  description:
    'Daniel Smith (dasm) — Technologies DASMLAB Inc. Living lab, 60+ countries, windsurfing, and engineering teach-back.',
  path: '/about',
  person: true,
  image: media('portrait.webp')
})
</script>

<style scoped>
.about-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0.75rem clamp(0.7rem, 2vw, 1.2rem) 2rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.about-hero {
  border: 1.5px solid rgba(31, 111, 98, 0.32);
  border-radius: 20px;
  padding: 1.2rem 1.25rem;
  background:
    radial-gradient(circle at 10% 20%, rgba(47, 143, 125, 0.14), transparent 42%),
    linear-gradient(155deg, #fff, #eef5f2);
}

.hero-title {
  margin: 0.25rem 0 0.4rem;
  font-family: 'Fraunces', Georgia, serif;
  font-size: clamp(1.5rem, 2.6vw, 2.15rem);
  color: #12202c;
  line-height: 1.1;
}

.hero-sub {
  margin: 0;
  max-width: 48rem;
  color: #3d5263;
  line-height: 1.5;
}

.about-grid {
  display: grid;
  grid-template-columns: 0.9fr 1.1fr;
  gap: 1rem;
}

.portrait-wrap,
.carousel-wrap,
.bio-panel {
  border: 1px solid rgba(28, 52, 73, 0.14);
  border-radius: 16px;
  background: linear-gradient(165deg, #fff, #f7fbf9);
  padding: 0.95rem;
  box-shadow: 0 10px 22px rgba(18, 40, 52, 0.06);
}

.about-head {
  width: 100%;
  max-width: 420px;
  aspect-ratio: 4 / 5;
  object-fit: cover;
  object-position: center 18%;
  border-radius: 14px;
  border: 1px solid rgba(41, 72, 99, 0.16);
  display: block;
}

.bio-panel h2 {
  margin: 0 0 0.5rem;
  font-family: 'Fraunces', Georgia, serif;
  color: #1f6f62;
  font-size: 1.2rem;
}

.bio-panel p {
  color: #4a5d6d;
  line-height: 1.6;
}

.media-note {
  margin: 0.55rem 0 0;
  font-size: 0.78rem;
  color: #5a6f80;
}

.media-note--inline {
  margin-top: 0.85rem;
  padding: 0.55rem 0.65rem;
  border-radius: 10px;
  background: rgba(47, 143, 125, 0.08);
  border: 1px dashed rgba(31, 111, 98, 0.3);
}

.media-note code {
  font-size: 0.75rem;
}

.waypoint {
  text-align: center;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  font-size: 0.72rem;
  font-weight: 700;
  color: #5a6f80;
  padding: 0.35rem;
  border-top: 1px dashed rgba(82, 108, 131, 0.35);
  border-bottom: 1px dashed rgba(82, 108, 131, 0.35);
}

.about-carousel {
  border-radius: 14px;
  overflow: hidden;
}

.about-lifestyle-img {
  width: min(100%, 560px);
  height: auto;
  max-height: 260px;
  object-fit: cover;
  border-radius: 10px;
}

.caption-band {
  margin-top: 0.5rem;
  font-size: 0.9rem;
  color: #435564;
  text-align: center;
}

@media (max-width: 860px) {
  .about-grid {
    grid-template-columns: 1fr;
  }
}
</style>
