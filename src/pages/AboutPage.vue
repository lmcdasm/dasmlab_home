<template>
  <q-page padding class="q-gutter-md">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">About</div>
        <h1 class="dasm-title">Living lab, lived experience</h1>
        <p class="dasm-subtitle">
          DASMLAB is where technology gets poked, prodded, and brought to life. Every build is part experiment and part teach-back.
        </p>
      </div>
    </section>

    <section class="dasm-grid dasm-grid--2">
      <div class="dasm-panel">
        <img :src="meHead" class="about-head" alt="DASMLAB founder portrait" />
      </div>
      <div class="dasm-panel">
        <h3 class="text-primary q-mt-none">Who am I?</h3>
        <p class="about-copy">
          I have spent two decades across telecom, automotive, cloud, and AI platforms. The throughline is practical systems
          thinking and shipping resilient solutions that people can understand and extend.
        </p>
        <p class="about-copy q-mb-none">
          From Outaouais/Ottawa to Fernie and now the Laurentians, the mission stays the same: make complex engineering feel teachable.
        </p>
      </div>
    </section>

    <div class="dasm-waypoint">Beyond the code</div>
    <section class="dasm-grid dasm-grid--2">
      <div class="dasm-panel">
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
              <img :src="img.src" :alt="img.caption" class="about-lifestyle-img" />
              <div class="caption-band">{{ img.caption }}</div>
            </div>
          </q-carousel-slide>
        </q-carousel>
      </div>
      <div class="dasm-panel">
        <h3 class="text-primary q-mt-none">Outside work</h3>
        <p class="about-copy">
          Snowboarding, windsurfing, gardening, house music, and pottery keep the creative loop open. They shape the rhythm
          of how I design systems: playful, grounded, and always iterative.
        </p>
      </div>
    </section>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'

const carouselIndex = ref(0)
const meHead = '/media/hero/profile-placeholder.svg'

const lifestyleImages = [
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'Windsurfing on a summer low-wind day' },
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'Snowboarding at Mont Ste. Anne' },
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'Country living focus' },
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'Garden therapy in action' },
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'Princess Leah herself' },
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'Current favorite sail setup' },
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'Weekend baking mode' },
  { src: '/media/hero/lifestyle-placeholder.svg', caption: 'And the reward at the end' }
]
</script>

<style scoped>
.about-head {
  width: 100%;
  max-width: 420px;
  border-radius: 14px;
  border: 1px solid rgba(41, 72, 99, 0.16);
}

.about-copy {
  color: #4a5d6d;
  line-height: 1.6;
}

.about-carousel {
  border-radius: 14px;
  background: linear-gradient(165deg, #ffffff, #f7fbfe);
  border: 1px solid rgba(41, 72, 99, 0.14);
}

.about-lifestyle-img {
  max-width: 100%;
  max-height: 236px;
  object-fit: contain;
  border-radius: 12px;
}

.caption-band {
  margin-top: 0.6rem;
  width: 100%;
  text-align: center;
  color: #4a5d6d;
  font-size: 0.85rem;
  background: rgba(242, 248, 253, 0.95);
  border: 1px solid rgba(41, 72, 99, 0.14);
  padding: 0.4rem 0.55rem;
  border-radius: 10px;
}
</style>
