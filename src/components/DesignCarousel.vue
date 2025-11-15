<template>
  <q-card
    class="q-my-md flex flex-center glass"
    style="max-width: 1000px; width: 100%;"
    v-motion
    :initial="{ opacity: 0, scale: 0.95 }"
    :enter="{ opacity: 1, scale: 1, transition: { delay: 300, duration: 500 } }"
  >
    <div class="flex flex-center" style="padding: 24px;">
      <div
        style="border-radius: 16px; overflow: hidden; background: rgba(255,255,255,0.02); width: 100%; height: 100%; display: flex; align-items: center; justify-content: center;"
        class="carousel-container"
      >
        <q-carousel
          v-if="entries.length > 0"
          v-model="slide"
          animated
          height="320px"
          transition-prev="slide-right"
          transition-next="slide-left"
          arrows
          navigation
          control-color="accent"
          swipeable
          class="carousel-glass"
        >
          <q-carousel-slide
            v-for="(entry, idx) in entries"
            :key="entry.id"
            :name="idx + 1"
          >
            <div class="carousel-slide-content">
              <img
                :src="fullImageUrl(entry.image_url)"
                :alt="entry.title || 'Carousel image'"
                style="width: 100%; height: 100%; object-fit: cover; cursor: zoom-in; transition: transform 0.3s;"
                class="carousel-image"
                @click="openModal(entry)"
                @mouseenter="$event.target.style.transform = 'scale(1.02)'"
                @mouseleave="$event.target.style.transform = 'scale(1)'"
                loading="lazy"
              />
              <div v-if="entry.title" class="carousel-title-overlay">
                <div class="text-subtitle1 text-white text-weight-bold">{{ entry.title }}</div>
              </div>
            </div>
          </q-carousel-slide>
        </q-carousel>
      </div>
    </div>

    <!-- Image Modal -->
    <q-dialog v-model="modalOpen" persistent>
      <div
        style="position: relative; max-width: 90vw; max-height: 90vh; background: #222; display: flex; align-items: center; justify-content: center; border-radius: 8px;"
        class="modal-container glass"
      >
        <img
          v-if="selectedEntry"
          :src="fullImageUrl(selectedEntry.image_url)"
          :alt="selectedEntry.title || 'Image'"
          style="max-width: 88vw; max-height: 80vh; background: #222; display: block; margin: auto; border-radius: 4px;"
          loading="lazy"
        />
        <q-btn
          flat
          round
          dense
          icon="close"
          color="white"
          @click="modalOpen = false"
          style="position: absolute; top: 10px; right: 10px; z-index: 1001; background: rgba(0,0,0,0.6); backdrop-filter: blur(10px);"
          class="close-btn"
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

function fullImageUrl(path) {
  if (!path) return ''
  if (/^https?:\/\//.test(path)) return path
  return DESIGN_CAROUSEL_HOST + path
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
.carousel-container {
  transition: all 0.3s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.carousel-container:hover {
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.4);
}

.carousel-slide-content {
  position: relative;
  width: 100%;
  height: 100%;
}

.carousel-title-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
  padding: 16px;
  backdrop-filter: blur(10px);
}

.carousel-image {
  border-radius: 8px;
}

.modal-container {
  transition: all 0.3s;
}

.close-btn {
  transition: all 0.2s;
}

.close-btn:hover {
  transform: scale(1.1);
  background: rgba(0, 0, 0, 0.8) !important;
}

/* For Quasar carousel arrows - most robust selector */
:deep(.q-carousel__arrow),
:deep(.q-carousel__arrow .q-icon),
:deep(.q-carousel__arrow .q-icon svg) {
  color: #64FFDA !important;
  fill: #64FFDA !important;
  stroke: #64FFDA !important;
  opacity: 1 !important;
}

:deep(.q-carousel__arrow--right),
:deep(.q-carousel__arrow--left) {
  color: #64FFDA !important;
  fill: #64FFDA !important;
  stroke: #64FFDA !important;
  background: rgba(100, 255, 218, 0.1) !important;
  backdrop-filter: blur(10px);
  transition: all 0.2s;
}

:deep(.q-carousel__arrow:hover) {
  background: rgba(100, 255, 218, 0.2) !important;
  transform: scale(1.1);
}

:deep(.q-carousel__navigation) {
  color: #64FFDA !important;
}

:deep(.q-carousel__navigation .q-carousel__navigation-icon) {
  color: #64FFDA !important;
  opacity: 0.6;
}

:deep(.q-carousel__navigation .q-carousel__navigation-icon--active) {
  opacity: 1;
}
</style>
