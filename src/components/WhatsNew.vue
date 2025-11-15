<template>
  <div class="whats-new-container glass q-pa-md" style="width: 100%; max-width: 800px;">
    <div
      class="text-h6 q-mb-sm text-left"
      v-motion
      :initial="{ opacity: 0, x: -20 }"
      :enter="{ opacity: 1, x: 0, transition: { delay: 100 } }"
    >
      <q-icon name="newspaper" class="q-mr-sm" />
      What's New?
    </div>

    <div
      class="scroll-wrapper"
      @mouseenter="pause = true"
      @mouseleave="pause = false"
      v-motion
      :initial="{ opacity: 0, y: 10 }"
      :enter="{ opacity: 1, y: 0, transition: { delay: 200 } }"
    >
      <div
        class="scroll-content"
        :class="{ paused: pause }"
      >
        <div
          v-for="entry in doubledNews"
          :key="entry.id + '-' + entry.loop"
          class="entry glass-inner q-pa-sm q-mb-xs"
        >
          <div class="text-subtitle2 text-weight-medium">
            <b>{{ entry.project }}</b>
          </div>
          <div class="text-body2 q-mt-xs">{{ entry.title }}</div>
          <div class="text-caption text-grey-5 q-mt-xs">{{ formatDate(entry.date) }}</div>
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

const formatDate = (dateString) => {
  if (!dateString) return ''
  try {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now - date
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffMins < 1) return 'Just now'
    if (diffMins < 60) return `${diffMins}m ago`
    if (diffHours < 24) return `${diffHours}h ago`
    if (diffDays < 7) return `${diffDays}d ago`
    return date.toLocaleDateString()
  } catch {
    return dateString
  }
}
</script>

<style scoped>
.whats-new-container {
  transition: transform 0.22s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.whats-new-container:hover {
  transform: translateY(-2px);
}

.scroll-wrapper {
  height: 120px;
  overflow: hidden;
  position: relative;
  border-radius: 8px;
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
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  min-height: 60px;
  transition: background 0.2s;
}

.entry:hover {
  background: rgba(255, 255, 255, 0.05);
}

.glass-inner {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(5px);
  border-radius: 8px;
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

/* Light mode adjustments */
body.body--light .entry {
  border-bottom-color: rgba(0, 0, 0, 0.1);
}

body.body--light .glass-inner {
  background: rgba(0, 0, 0, 0.02);
}

body.body--light .entry:hover {
  background: rgba(0, 0, 0, 0.03);
}
</style>
