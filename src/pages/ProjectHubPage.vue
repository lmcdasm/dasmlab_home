<template>
  <q-page v-if="hub" padding class="q-gutter-md">
    <AnswerFirstArticle
      eyebrow="Project hub"
      :title="hub.title"
      :answer="hub.answer"
      :faq="hub.faq"
    >
      <p class="hub-summary">{{ hub.summary }}</p>

      <div class="hub-ctas">
        <q-btn
          v-if="hub.liveUrl"
          color="primary"
          unelevated
          no-caps
          label="Open live"
          :href="hub.liveUrl"
          target="_blank"
          rel="noopener"
        />
        <q-btn
          v-if="hub.demoUrl"
          color="secondary"
          outline
          no-caps
          label="Try demo (fake mode)"
          :href="hub.demoUrl"
          target="_blank"
          rel="noopener"
        />
        <q-btn
          v-if="hub.sourceUrl || hub.ceUrl"
          flat
          no-caps
          color="primary"
          :label="hub.ceUrl ? 'Community Edition' : 'Source'"
          :href="hub.ceUrl || hub.sourceUrl"
          target="_blank"
          rel="noopener"
        />
        <q-btn flat no-caps :to="lanePath" :label="`${hub.lane} lane`" />
      </div>

      <section class="dasm-panel">
        <h2 class="hub-h2">Architecture</h2>
        <p class="hub-body">{{ hub.architecture }}</p>
      </section>

      <section class="dasm-panel">
        <h2 class="hub-h2">How we built this</h2>
        <p class="hub-body">{{ hub.howWeBuilt }}</p>
      </section>

      <BehindTheDesign :stack="hub.stack" :topics="hub.topics" />
    </AnswerFirstArticle>
  </q-page>
  <q-page v-else padding>
    <div class="dasm-panel">
      <h1 class="dasm-title">Project not found</h1>
      <p><router-link to="/">Back to home</router-link></p>
    </div>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AnswerFirstArticle from 'src/components/AnswerFirstArticle.vue'
import BehindTheDesign from 'src/components/BehindTheDesign.vue'
import { projectHubs, SITE } from 'src/data/hubs'
import { useSeo } from 'src/composables/useSeo'

const route = useRoute()
const hub = computed(() => projectHubs[route.params.slug] || null)

const lanePath = computed(() => {
  const map = {
    Frontend: '/projects/frontend',
    Backend: '/projects/backend',
    'AI/ML': '/projects/ai-ml',
    Cloud: '/projects/cloud',
    Infrastructure: '/projects/infrastructure',
    Security: '/projects/security'
  }
  return map[hub.value?.lane] || '/projects/frontend'
})

const seo = computed(() => {
  if (!hub.value) {
    return { title: 'Not found', description: 'Project hub not found', path: route.path }
  }
  return {
    title: hub.value.title,
    description: hub.value.summary,
    path: `/projects/${hub.value.slug}`,
    type: 'article',
    faq: hub.value.faq,
    jsonLd: {
      '@context': 'https://schema.org',
      '@type': 'SoftwareApplication',
      name: hub.value.title,
      description: hub.value.summary,
      applicationCategory: 'DeveloperApplication',
      url: `${SITE.url}/projects/${hub.value.slug}`,
      author: { '@type': 'Organization', name: SITE.legalName, url: SITE.url }
    }
  }
})

useSeo(seo)
</script>

<style scoped>
.hub-summary {
  margin: 0;
  color: #4a5d6d;
  line-height: 1.55;
}

.hub-ctas {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.hub-h2 {
  margin: 0 0 0.5rem;
  font-size: 1.15rem;
  color: #1d2b36;
}

.hub-body {
  margin: 0;
  color: #4a5d6d;
  line-height: 1.6;
}
</style>
