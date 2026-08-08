<template>
  <q-card
    class="project-card q-mb-md q-pa-md"
    bordered
    v-scroll="onScroll"
    :class="{ 'animated fadeInUp': isVisible }"
  >
    <div class="row items-start justify-between no-wrap q-col-gutter-sm">
      <div class="col">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge v-if="category" color="primary" class="token-chip">{{ category }}</q-badge>
          <q-badge v-if="badge" color="secondary" outline>{{ badge }}</q-badge>
        </div>
        <div class="text-h6 card-title">
          <router-link
            v-if="hubPath"
            :to="hubPath"
            class="text-primary card-link"
          >
            {{ title }}
          </router-link>
          <a
            v-else
            :href="url || '#'"
            class="text-primary card-link"
            target="_blank"
            rel="noopener"
            @click.prevent="onTitleClick"
          >
            {{ title }}
          </a>
        </div>
      </div>
      <div class="icon-bubble">
        <q-icon name="rocket_launch" size="18px" />
      </div>
    </div>

    <div class="q-mt-sm text-body2 card-description">{{ description }}</div>

    <div class="q-mt-sm detail-line">
      <q-icon name="code" size="16px" class="q-mr-xs" />
      <span class="text-caption">{{ language || 'Stack not listed yet' }}</span>
    </div>

    <div v-if="hubPath" class="q-mt-sm text-caption detail-line">
      <q-icon name="hub" class="q-mr-xs" />
      <span class="detail-label">Hub:</span>
      <router-link :to="hubPath" class="text-secondary link-inline">Behind the Design</router-link>
    </div>

    <div class="q-mt-sm text-caption detail-line">
      <q-icon name="link" class="q-mr-xs" />
      <span class="detail-label">Live:</span>
      <template v-if="liveUrl">
        <a :href="liveUrl" class="text-secondary link-inline" target="_blank" rel="noopener">
          {{ liveUrl.replace(/^https?:\/\//, '') }}
        </a>
      </template>
      <template v-else>
        <span class="text-grey-6 q-ml-xs">coming soon</span>
      </template>
    </div>

    <div v-if="tutoUrl" class="q-mt-sm text-caption detail-line">
      <q-icon name="school" class="q-mr-xs" />
      <span class="detail-label">Tutorial:</span>
      <a href="#" class="text-secondary link-inline" @click.prevent="showTutorial">
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
  category: String,
  liveUrl: String,
  tutoUrl: String,
  hubPath: String
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

const onTitleClick = () => {
  if (!props.url) return
  window.open(props.url, '_blank', 'noopener')
}
</script>

<style scoped>
.project-card {
  border: 1px solid var(--dasm-border-soft);
  border-radius: 14px;
  background: linear-gradient(165deg, #ffffff, #f7fbfe);
  box-shadow: 0 10px 22px rgba(24, 44, 62, 0.08);
  transition: transform 170ms ease, border-color 170ms ease, box-shadow 170ms ease;
}

.project-card:hover {
  transform: translateY(-2px);
  border-color: rgba(151, 110, 176, 0.42);
  box-shadow: 0 16px 26px rgba(24, 44, 62, 0.14);
}

.token-chip {
  letter-spacing: 0.06em;
}

.card-title {
  margin-bottom: 0.2rem;
  line-height: 1.18;
}

.card-link {
  text-decoration: none;
}

.card-link:hover {
  text-decoration: underline;
}

.card-description {
  color: #4b5d6d;
  line-height: 1.45;
}

.detail-line {
  color: #647686;
}

.detail-label {
  color: #5c6f80;
  margin-right: 0.3rem;
}

.link-inline {
  text-decoration: none;
}

.link-inline:hover {
  text-decoration: underline;
}

.icon-bubble {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(63, 122, 107, 0.32);
  background: radial-gradient(circle at 35% 30%, rgba(63, 122, 107, 0.16), rgba(255, 255, 255, 0.9));
  color: #4d8f82;
}

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

