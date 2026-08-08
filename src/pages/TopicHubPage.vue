<template>
  <q-page v-if="hub" padding class="q-gutter-md">
    <AnswerFirstArticle
      eyebrow="Topic hub"
      :title="hub.title"
      :answer="hub.answer"
      :faq="hub.faq"
    >
      <p class="topic-summary">{{ hub.summary }}</p>

      <section class="dasm-panel">
        <h2 class="topic-h2">Projects using {{ hub.title }}</h2>
        <ul class="topic-list">
          <li v-for="slug in hub.projects" :key="slug">
            <router-link :to="`/projects/${slug}`">{{ projectTitle(slug) }}</router-link>
          </li>
        </ul>
      </section>
    </AnswerFirstArticle>
  </q-page>
  <q-page v-else padding>
    <div class="dasm-panel">
      <h1 class="dasm-title">Topic not found</h1>
      <p><router-link to="/">Back to home</router-link></p>
    </div>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AnswerFirstArticle from 'src/components/AnswerFirstArticle.vue'
import { topicHubs, projectHubs, SITE } from 'src/data/hubs'
import { useSeo } from 'src/composables/useSeo'

const route = useRoute()
const hub = computed(() => topicHubs[route.params.tech] || null)

function projectTitle(slug) {
  return projectHubs[slug]?.title || slug
}

const seo = computed(() => {
  if (!hub.value) {
    return { title: 'Not found', description: 'Topic hub not found', path: route.path }
  }
  return {
    title: hub.value.title,
    description: hub.value.summary,
    path: `/topics/${hub.value.slug}`,
    type: 'article',
    faq: hub.value.faq,
    jsonLd: {
      '@context': 'https://schema.org',
      '@type': 'TechArticle',
      headline: hub.value.title,
      description: hub.value.summary,
      author: { '@type': 'Organization', name: SITE.legalName, url: SITE.url }
    }
  }
})

useSeo(seo)
</script>

<style scoped>
.topic-summary {
  margin: 0;
  color: #4a5d6d;
  line-height: 1.55;
}

.topic-h2 {
  margin: 0 0 0.5rem;
  font-size: 1.15rem;
  color: #1d2b36;
}

.topic-list {
  margin: 0;
  padding-left: 1.2rem;
  line-height: 1.8;
}

.topic-list a {
  color: #1f6f62;
  font-weight: 600;
  text-decoration: none;
}

.topic-list a:hover {
  text-decoration: underline;
}
</style>
