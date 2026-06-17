<template>
  <div class="news-shell">
    <div class="news-header q-mb-sm">
      <div class="text-h6 text-left">
        <q-icon name="newspaper" class="q-mr-sm" />
        What’s New?
      </div>
      <div v-if="normalizedNews.length > 0" class="news-counter text-caption">
        {{ slide + 1 }} / {{ normalizedNews.length }}
      </div>
    </div>
    <div v-if="normalizedNews.length === 0" class="empty-shell text-caption q-pa-md">
      Waiting for feed updates...
    </div>
    <div v-else class="news-carousel-wrap" @mouseenter="autoplay = false" @mouseleave="autoplay = true">
      <q-carousel
        v-model="slide"
        animated
        swipeable
        infinite
        :autoplay="autoplay ? 9000 : false"
        transition-prev="slide-right"
        transition-next="slide-left"
        control-color="primary"
        navigation
        arrows
        height="176px"
        class="news-carousel"
      >
        <q-carousel-slide
          v-for="(entry, idx) in normalizedNews"
          :key="entry.key"
          :name="idx"
          class="news-slide"
        >
          <div class="news-card">
            <div class="news-card-top">
              <q-chip dense color="primary" text-color="white" icon="fork_right" class="news-chip">
                {{ entry.repoLabel }}
              </q-chip>
              <div class="news-badge">New commit activity</div>
            </div>
            <div class="news-title">{{ entry.titleSummary }}</div>
            <div class="news-meta">
              <q-icon name="schedule" size="14px" class="q-mr-xs" />
              {{ entry.dateLabel }}
            </div>
          </div>
        </q-carousel-slide>
      </q-carousel>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  news: {
    type: Array,
    required: true
  }
})

const slide = ref(0)
const autoplay = ref(true)

const normalizedNews = computed(() => (props.news || []).map((item, idx) => {
  const repo = item?.project || 'dasmlab/unknown'
  const title = item?.title || 'Recent updates pushed to repository.'
  const date = item?.date || ''

  return {
    key: `${repo}-${date}-${idx}`,
    repoLabel: repo.length > 34 ? `${repo.slice(0, 34)}...` : repo,
    titleSummary: title.length > 118 ? `${title.slice(0, 118)}...` : title,
    dateLabel: date || 'just now'
  }
}))

watch(
  () => normalizedNews.value.length,
  (len) => {
    if (len === 0) {
      slide.value = 0
      return
    }
    if (slide.value > len - 1) {
      slide.value = 0
    }
  }
)
</script>

<style scoped>
.news-shell {
  width: 100%;
  max-width: 920px;
  margin: 0 auto;
}

.news-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.news-counter {
  color: #5e7484;
  letter-spacing: 0.04em;
}

.empty-shell {
  border: 1px dashed rgba(41, 72, 99, 0.24);
  border-radius: 12px;
  color: #6b7d8d;
  background: rgba(255, 255, 255, 0.72);
}

.news-carousel-wrap {
  position: relative;
}

.news-carousel {
  border: 1px solid rgba(41, 72, 99, 0.16);
  border-radius: 12px;
  background: linear-gradient(165deg, rgba(255, 255, 255, 0.98), rgba(246, 251, 254, 0.98));
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.7);
}

.news-slide {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.4rem 0.3rem 0.95rem;
}

.news-card {
  width: 100%;
  height: 100%;
  border: 1px solid rgba(41, 72, 99, 0.14);
  border-radius: 12px;
  padding: 0.72rem 0.82rem;
  background: linear-gradient(160deg, #ffffff, #f7fbfe);
  box-shadow: 0 8px 18px rgba(24, 44, 62, 0.07);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 0.4rem;
}

.news-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.45rem;
}

.news-chip {
  max-width: 78%;
}

.news-chip :deep(.q-chip__content) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.news-badge {
  flex-shrink: 0;
  font-size: 0.68rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #687d8c;
}

.news-title {
  color: #2a3d50;
  font-weight: 600;
  line-height: 1.35;
  font-size: 0.98rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.news-meta {
  color: #647888;
  font-size: 0.8rem;
  display: inline-flex;
  align-items: center;
}

@media (max-width: 700px) {
  .news-badge {
    display: none;
  }

  .news-title {
    font-size: 0.92rem;
  }
}
</style>

