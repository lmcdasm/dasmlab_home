<template>
  <section class="hero-section">
    <div class="relative overflow-hidden glass q-pa-lg">
      <!-- Particles Background -->
      <vue-particles
        id="bgParticles"
        :options="particleOpts"
        :particles-loaded="particlesLoaded"
        class="absolute-full"
      />

      <!-- Hero Content -->
      <div
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 200, duration: 600 } }"
        class="relative q-gutter-md"
      >
        <div class="text-h3 text-white text-weight-bold">DASMLAB</div>
        <div class="text-subtitle1 text-grey-4 q-mt-sm">
          GPU-heavy AI/ML experiments • Go backends • Quasar frontends
        </div>
        <div class="row q-gutter-sm q-mt-md">
          <q-btn
            color="secondary"
            label="View Projects"
            to="/#/projects/frontend"
            unelevated
            size="md"
            class="q-px-lg"
          />
          <q-btn
            color="accent"
            label="What's New"
            to="/#/"
            flat
            size="md"
            class="q-px-lg"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue'

const particleOpts = ref({
  fullScreen: false,
  background: {
    color: {
      value: 'transparent'
    }
  },
  fpsLimit: 60,
  particles: {
    number: {
      value: 28,
      density: {
        enable: true,
        area: 800
      }
    },
    links: {
      enable: true,
      opacity: 0.2,
      distance: 140,
      color: '#64FFDA'
    },
    move: {
      enable: true,
      speed: 0.6
    },
    opacity: {
      value: 0.3
    },
    size: {
      value: {
        min: 1,
        max: 3
      }
    },
    color: {
      value: '#64FFDA'
    }
  }
})

// Particles are loaded via the plugin boot file (src/boot/particles.ts)
// This callback is called when particles are initialized
const particlesLoaded = async (container) => {
  if (container) {
    console.log('Particles container loaded:', container.id)
  }
}
</script>

<style scoped>
.hero-section {
  margin-bottom: 1rem;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.relative {
  position: relative;
}

.absolute-full {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.relative > * {
  position: relative;
  z-index: 1;
}
</style>

