<template>
  <q-card class="q-my-md flex flex-center" style="max-width: 700px; width: 100%;">
    <div class="flex flex-center" style="padding: 24px;">
      <div style="border-radius: 16px; overflow: hidden; background: #fff; width: 100%; height: 100%; display: flex; align-items: center; justify-content: center;">
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
              style="width: 100%; height: 100%; object-fit: cover; cursor: zoom-in;"
              @click="openModal(entry)"
            />
          </q-carousel-slide>
        </q-carousel>
      </div>
    </div>

<!-- DEBUG    <div v-if="modalOpen" style="position: fixed; top: 10px; left: 10px; background: yellow; color: black; z-index: 2000;">
      Modal Open! Entry: {{ selectedEntry && selectedEntry.title }}
    </div> -->

    <!-- Image Modal -->
    <q-dialog v-model="modalOpen" persistent>
     <div style="position: relative; max-width: 90vw; max-height: 90vh; background: #222; display: flex; align-items: center; justify-content: center;">
       <img
         v-if="selectedEntry"
         :src="fullImageUrl(selectedEntry.image_url)"
         :alt="selectedEntry.title"
         style="max-width: 88vw; max-height: 80vh; background: #222; display: block; margin: auto;"
       />
       <q-btn
         flat
         round
         dense
         icon="close"
         color="white"
         @click="modalOpen = false"
         style="position: absolute; top: 10px; right: 10px; z-index: 1001; background: rgba(0,0,0,0.4);"
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
/* For Quasar carousel arrows - most robust selector */
.q-carousel__arrow,
.q-carousel__arrow .q-icon,
.q-carousel__arrow .q-icon svg {
  color: #22746b !important;    /* Use your brand's hex here */
  fill: #22746b !important;
  stroke: #22746b !important;
  opacity: 1 !important;
}
.q-carousel__arrow--right,
.q-carousel__arrow--left {
  color: #22746b !important;
  fill: #22746b !important;
  stroke: #22746b !important;
}
</style>

