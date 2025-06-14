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

    <!-- Tutorial Link -->
    <div v-if="tutoUrl" class="q-mt-sm text-caption">
      <q-icon name="school" class="q-mr-xs" />
      <span class="text-grey-9">Tutorial:</span>
      <a href="#" class="text-secondary" @click.prevent="showTutorial">
        View Tutorial
      </a>
    </div>
  </q-card>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import TutorialPlayer from './TutorialPlayer.vue' // adjust if needed

// ✅ This defines the incoming props
const props = defineProps({
  title: String,
  description: String,
  url: String,
  language: String,
  badge: String,
  liveUrl: String,
  tutoUrl: String
})

// ✅ Access Quasar Dialog plugin
const $q = useQuasar()

// ✅ Show tutorial modal (called on link click)
const showTutorial = () => {
  if (!props.tutoUrl) return

  $q.dialog({
    component: TutorialPlayer,
    componentProps: {
      tutorialUrl: props.tutoUrl
    }
  })
}

// ✅ Fade-in effect
const isVisible = ref(false)
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

