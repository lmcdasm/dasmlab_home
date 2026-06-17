<template>
  <div class="news-shell">
    <div class="news-header q-mb-sm">
      <div class="text-h6 text-left news-heading">
        <span class="news-heading-icon">
          <q-icon name="newspaper" />
        </span>
        <span>What’s New?</span>
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
        height="224px"
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
            <div class="news-title-bubble">
              <div class="news-title">{{ entry.titleSummary }}</div>
            </div>
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
  gap: 0.75rem;
}

.news-heading {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
}

.news-heading-icon {
  width: 1.85rem;
  height: 1.85rem;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #3d5f78;
  background: linear-gradient(155deg, rgba(168, 207, 235, 0.45), rgba(126, 172, 205, 0.25));
  box-shadow: inset 0 0 0 1px rgba(62, 97, 122, 0.2), 0 4px 12px rgba(18, 43, 61, 0.08);
}

.news-counter {
  color: #5e7484;
  letter-spacing: 0.04em;
  font-weight: 600;
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
  overflow: hidden;
}

.news-slide {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.45rem 2.45rem 2.2rem;
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

.news-title-bubble {
  border-radius: 10px;
  border: 1px solid rgba(79, 116, 143, 0.2);
  background: linear-gradient(150deg, rgba(232, 242, 249, 0.7), rgba(240, 248, 253, 0.92));
  padding: 0.58rem 0.72rem;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.76);
}

.news-title {
  color: #2a3d50;
  font-weight: 600;
  line-height: 1.45;
  font-size: 0.96rem;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.news-meta {
  color: #5a6f80;
  font-size: 0.82rem;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  margin-top: 0.2rem;
}

.news-carousel :deep(.q-carousel__prev-arrow),
.news-carousel :deep(.q-carousel__next-arrow) {
  width: 2rem;
  height: 2rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.95);
  color: #365871;
  box-shadow: 0 8px 16px rgba(23, 48, 66, 0.14);
  top: 47%;
}

.news-carousel :deep(.q-carousel__prev-arrow) {
  left: 0.5rem;
}

.news-carousel :deep(.q-carousel__next-arrow) {
  right: 0.5rem;
}

.news-carousel :deep(.q-carousel__navigation) {
  bottom: 0.42rem;
}

.news-carousel :deep(.q-carousel__navigation .q-btn) {
  margin: 0 2px;
  min-width: 14px;
  min-height: 14px;
}

@media (max-width: 700px) {
  .news-header {
    align-items: flex-start;
  }

  .news-counter {
    margin-top: 0.22rem;
  }

  .news-badge {
    display: none;
  }

  .news-slide {
    padding: 0.35rem 1.95rem 2.1rem;
  }

  .news-heading-icon {
    width: 1.7rem;
    height: 1.7rem;
  }

  .news-title {
    font-size: 0.92rem;
  }
}
</style>

