<template>
  <div class="news-shell">
    <div class="text-h6 q-mb-sm text-left">
      <q-icon name="newspaper" class="q-mr-sm" />
      What’s New?
    </div>
    <div v-if="news.length === 0" class="text-caption text-grey-5 q-pa-md">
      Waiting for feed updates...
    </div>
    <div v-else class="scroll-wrapper" @mouseenter="pause = true" @mouseleave="pause = false">
      <div class="scroll-content" :class="{ paused: pause }">
        <div v-for="entry in doubledNews" :key="entry.id + '-' + entry.loop" class="entry">
          <div class="text-subtitle1"><b>{{ entry.project }}</b></div>
          <div class="text-subtitle2">{{ entry.title }}</div>
          <div class="text-caption text-grey-7">{{ entry.date }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  news: {
    type: Array,
    required: true
  }
})

const pause = ref(false)

const doubledNews = computed(() =>
  [...props.news, ...props.news].map((item, idx) => ({
    ...item,
    loop: Math.floor(idx / props.news.length)
  }))
)
</script>

<style scoped>
.news-shell {
  width: 100%;
  max-width: 520px;
}

.scroll-wrapper {
  height: 108px;
  overflow: hidden;
  position: relative;
  border: 1px solid rgba(41, 72, 99, 0.16);
  border-radius: 12px;
  background: linear-gradient(165deg, rgba(255, 255, 255, 0.98), rgba(246, 251, 254, 0.98));
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.7);
  padding: 0 0.55rem;
}

.scroll-content {
  display: flex;
  flex-direction: column;
  animation: scroll-up 30s linear infinite;
}

.scroll-content.paused {
  animation-play-state: paused;
}

.entry {
  padding: 0.5em 0;
  border-bottom: 1px solid rgba(41, 72, 99, 0.12);
}

@keyframes scroll-up {
  0% {
    transform: translateY(0%);
  }
  100% {
    transform: translateY(-50%);
  }
}
</style>

