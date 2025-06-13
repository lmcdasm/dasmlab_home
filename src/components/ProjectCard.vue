<template>
  <q-card
    class="q-mb-md q-pa-md shadow-2"
    bordered
    v-scroll="onScroll"
    :class="{ 'animated fadeInUp': isVisible }"
  >
    <div class="row items-center justify-between">
      <div class="text-h6 text-primary">
        <a :href="url" class="text-primary" target="_blank" rel="noopener">
          {{ title }}
        </a>
      </div>
      <q-badge v-if="badge" color="secondary" outline>{{ badge }}</q-badge>
    </div>

    <div class="q-mt-sm text-body2 text-grey-8">{{ description }}</div>

    <div class="q-mt-sm">
      <q-icon name="code" size="16px" class="q-mr-xs" />
      <span class="text-caption">{{ language }}</span>
    </div>

    <div class="q-mt-sm text-caption">
      <q-icon name="link" class="q-mr-xs" />
      <span class="text-grey-9">Live Example:</span>
      <template v-if="liveUrl">
        <a :href="liveUrl" class="text-secondary" target="_blank" rel="noopener">
          {{ liveUrl.replace(/^https?:\/\//, '') }}
        </a>
      </template>
      <template v-else>
        <span class="text-grey-6 q-ml-xs">Coming soon</span>
      </template>
    </div>
  </q-card>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  title: String,
  description: String,
  url: String,
  language: String,
  badge: String,
  liveUrl: String
})

const isVisible = ref(false)

// ESLint-safe scroll trigger
function onScroll() {
  isVisible.value = true
}
</script>

<style scoped>
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
.animated {
  animation-duration: 0.6s;
  animation-fill-mode: both;
}
.fadeInUp {
  animation-name: fadeInUp;
}
</style>

