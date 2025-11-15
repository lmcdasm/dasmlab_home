<template>
  <q-card
    class="q-mb-md q-pa-md glass tilt"
    v-motion
    :initial="{ opacity: 0, y: 20 }"
    :enter="{ opacity: 1, y: 0, transition: { delay: 100 } }"
    bordered
  >
    <div class="row items-center justify-between">
      <div class="text-h6 text-primary">
        <a
          :href="url"
          class="text-primary text-weight-medium"
          target="_blank"
          rel="noopener"
          style="text-decoration: none; transition: color 0.2s;"
        >
          {{ title }}
        </a>
      </div>
      <q-badge v-if="badge" color="secondary" outline class="badge-glow">{{ badge }}</q-badge>
    </div>

    <div class="q-mt-sm text-body2 text-grey-8">{{ description }}</div>

    <div class="q-mt-sm">
      <q-icon name="code" size="16px" class="q-mr-xs text-accent" />
      <span class="text-caption">{{ language }}</span>
    </div>

    <div class="q-mt-sm text-caption">
      <q-icon name="link" class="q-mr-xs text-secondary" />
      <span class="text-grey-9">Live Example:</span>
      <template v-if="liveUrl">
        <a
          :href="liveUrl"
          class="text-secondary"
          target="_blank"
          rel="noopener"
          style="text-decoration: none; transition: color 0.2s;"
        >
          {{ liveUrl.replace(/^https?:\/\//, '') }}
        </a>
      </template>
      <template v-else>
        <span class="text-grey-6 q-ml-xs">Coming soon</span>
      </template>
    </div>

    <!-- Tutorial Link -->
    <div v-if="tutoUrl" class="q-mt-sm text-caption">
      <q-icon name="school" class="q-mr-xs text-accent" />
      <span class="text-grey-9">Tutorial:</span>
      <a
        href="#"
        class="text-secondary"
        @click.prevent="showTutorial"
        style="text-decoration: none; transition: color 0.2s; cursor: pointer;"
      >
        View Tutorial
      </a>
    </div>
  </q-card>
</template>

<script setup>
import { useQuasar } from 'quasar'
import TutorialPlayer from './TutorialPlayer.vue'

const props = defineProps({
  title: String,
  description: String,
  url: String,
  language: String,
  badge: String,
  liveUrl: String,
  tutoUrl: String
})

const $q = useQuasar()

const showTutorial = () => {
  if (!props.tutoUrl) return

  $q.dialog({
    component: TutorialPlayer,
    componentProps: {
      tutorialUrl: props.tutoUrl
    }
  })
}
</script>

<style scoped>
a:hover {
  color: var(--q-accent) !important;
}

.badge-glow {
  transition: all 0.22s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.badge-glow:hover {
  box-shadow: 0 0 8px currentColor;
  transform: scale(1.05);
}
</style>
