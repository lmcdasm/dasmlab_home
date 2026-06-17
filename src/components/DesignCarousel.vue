<template>
  <q-card class="carousel-shell q-my-md flex flex-center">
    <div class="carousel-shell__inner flex flex-center">
      <div class="carousel-frame">
        <q-carousel
          v-if="entries.length > 0"
          v-model="slide"
          animated
          height="300px"
          transition-prev="fade"
          transition-next="fade"
          arrows
          navigation
          color="primary"
        >
          <q-carousel-slide
            v-for="(entry, idx) in entries"
            :key="entry.id"
            :name="idx + 1"
          >
            <img
              :src="fullImageUrl(entry.image_url)"
              :alt="entry.title"
              class="carousel-image"
              @error="onImageError"
              @click="openModal(entry)"
            />
          </q-carousel-slide>
        </q-carousel>
        <div v-else class="carousel-empty">
          <q-icon name="imagesmode" size="34px" class="q-mb-sm" />
          <div class="text-subtitle2">Featured experiments are loading</div>
          <div class="text-caption">Design carousel service not reachable right now.</div>
        </div>
      </div>
    </div>

    <q-dialog v-model="modalOpen" persistent>
      <div class="carousel-modal">
        <img
          v-if="selectedEntry"
          :src="fullImageUrl(selectedEntry.image_url)"
          :alt="selectedEntry.title"
          class="carousel-modal__image"
        />
        <q-btn
          flat
          round
          dense
          icon="close"
          color="white"
          @click="modalOpen = false"
          class="carousel-modal__close"
        />
      </div>
    </q-dialog>
  </q-card>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  entries: {
    type: Array,
    default: () => []
  }
})

const interval = 5000
const slide = ref(1)
const DESIGN_CAROUSEL_HOST = 'https://design-carousel.svc.dasmlab.org'
const FALLBACK_IMAGE = '/media/hero/lifestyle-placeholder.svg'

function fullImageUrl(path) {
  if (!path) return FALLBACK_IMAGE
  if (/^https?:\/\//.test(path)) return path
  return DESIGN_CAROUSEL_HOST + path
}

function onImageError(event) {
  event.target.src = FALLBACK_IMAGE
}

// Modal dialog state
const modalOpen = ref(false)
const selectedEntry = ref(null)

function openModal(entry) {
  selectedEntry.value = entry
  modalOpen.value = true
}

// Always reset slide to 1 when entries load (important for auto-interval to work)
watch(
  () => props.entries.length,
  async (n) => {
    if (n > 0) {
      await nextTick()
      slide.value = 1
    }
  }
)

// ---- MANUAL AUTO-ROTATION ----
let autoTimer = null

onMounted(() => {
  autoTimer = setInterval(() => {
    if (props.entries.length > 1) {
      slide.value = slide.value >= props.entries.length ? 1 : slide.value + 1
    }
  }, interval)
})

onBeforeUnmount(() => {
  if (autoTimer) clearInterval(autoTimer)
})
</script>

<style scoped>
.carousel-shell {
  max-width: 760px;
  width: 100%;
  border: 1px solid rgba(41, 72, 99, 0.15);
  background: linear-gradient(160deg, #ffffff, #f7fbfe);
  box-shadow: 0 10px 22px rgba(25, 47, 67, 0.08);
}

.carousel-shell__inner {
  padding: 18px;
  width: 100%;
}

.carousel-frame {
  border-radius: 16px;
  overflow: hidden;
  background: #f4f9fc;
  width: 100%;
  height: 100%;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: zoom-in;
}

.carousel-empty {
  min-height: 260px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #54697c;
  background: radial-gradient(circle at 50% 35%, rgba(63, 122, 107, 0.16), rgba(246, 250, 254, 0.95));
}

.carousel-modal {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
  background: rgba(20, 30, 38, 0.92);
  display: flex;
  align-items: center;
  justify-content: center;
}

.carousel-modal__image {
  max-width: 88vw;
  max-height: 80vh;
  background: transparent;
  display: block;
  margin: auto;
}

.carousel-modal__close {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 1001;
  background: rgba(0, 0, 0, 0.4);
}

.q-carousel__arrow,
.q-carousel__arrow .q-icon,
.q-carousel__arrow .q-icon svg {
  color: #2c8f7e !important;
  fill: #2c8f7e !important;
  stroke: #2c8f7e !important;
  opacity: 1 !important;
}

.q-carousel__arrow--right,
.q-carousel__arrow--left {
  color: #2c8f7e !important;
  fill: #2c8f7e !important;
  stroke: #2c8f7e !important;
}
</style>

