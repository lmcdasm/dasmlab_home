<template>
  <q-page v-if="lab" padding class="q-gutter-md">
    <AnswerFirstArticle
      eyebrow="Lab / experiment"
      :title="lab.title"
      :answer="lab.answer"
      :faq="lab.faq"
    >
      <p class="lab-summary">{{ lab.summary }}</p>
      <section class="dasm-panel">
        <h2 class="lab-h2">Related projects</h2>
        <ul class="lab-list">
          <li v-for="slug in lab.relatedProjects" :key="slug">
            <router-link :to="`/projects/${slug}`">{{ projectTitle(slug) }}</router-link>
          </li>
        </ul>
      </section>
    </AnswerFirstArticle>
  </q-page>
  <q-page v-else padding>
    <div class="dasm-panel">
      <h1 class="dasm-title">Lab not found</h1>
      <p><router-link to="/">Back to home</router-link></p>
    </div>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AnswerFirstArticle from 'src/components/AnswerFirstArticle.vue'
import { labs, projectHubs, SITE } from 'src/data/hubs'
import { useSeo } from 'src/composables/useSeo'

const route = useRoute()
const lab = computed(() => labs[route.params.slug] || null)

function projectTitle(slug) {
  return projectHubs[slug]?.title || slug
}

const seo = computed(() => {
  if (!lab.value) {
    return { title: 'Not found', description: 'Lab not found', path: route.path }
  }
  return {
    title: lab.value.title,
    description: lab.value.summary,
    path: `/labs/${lab.value.slug}`,
    type: 'article',
    faq: lab.value.faq,
    jsonLd: {
      '@context': 'https://schema.org',
      '@type': 'TechArticle',
      headline: lab.value.title,
      description: lab.value.summary,
      author: { '@type': 'Organization', name: SITE.legalName, url: SITE.url }
    }
  }
})

useSeo(seo)
</script>

<style scoped>
.lab-summary {
  margin: 0;
  color: #4a5d6d;
  line-height: 1.55;
}

.lab-h2 {
  margin: 0 0 0.5rem;
  font-size: 1.15rem;
  color: #1d2b36;
}

.lab-list {
  margin: 0;
  padding-left: 1.2rem;
  line-height: 1.8;
}

.lab-list a {
  color: #1f6f62;
  font-weight: 600;
  text-decoration: none;
}
</style>
