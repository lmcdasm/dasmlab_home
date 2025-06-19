<template>
  <div style="width: 100%; max-width: 500px;">
    <div class="text-h6 q-mb-sm text-left">
      <q-icon name="newspaper" class="q-mr-sm" />
      What’s New?
    </div>

    <div class="scroll-wrapper" @mouseenter="pause = true" @mouseleave="pause = false">
      <div
        class="scroll-content"
        :class="{ paused: pause }"
      >
        <div
          v-for="entry in doubledNews"
          :key="entry.id + '-' + entry.loop"
          class="entry"
        >
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
.scroll-wrapper {
  height: 100px;
  overflow: hidden;
  position: relative;
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
  border-bottom: 1px solid #ccc;
}

/* Keyframes for infinite vertical scrolling */
@keyframes scroll-up {
  0% {
    transform: translateY(0%);
  }
  100% {
    transform: translateY(-50%);
  }
}
</style>

