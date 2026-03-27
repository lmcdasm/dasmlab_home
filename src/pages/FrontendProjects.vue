<template>
  <q-page padding class="q-gutter-md">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">Project lane</div>
        <h1 class="dasm-title">Frontend projects</h1>
        <p class="dasm-subtitle">
          UI engineering across Vue, Quasar, mobile, voice, and accessibility-first experiences. Built for clarity and real-world delivery.
        </p>
      </div>
    </section>

    <div class="dasm-waypoint">Scroll and explore active builds</div>
    <section class="dasm-panel">
      <ProjectCard
        v-for="(project, index) in projects"
        :key="index"
        :title="project.title"
        :description="project.description"
        :url="project.url"
        :language="project.language"
        :badge="project.badge"
        category="Frontend"
        :live-url="project.liveUrl"
        :tuto-url="project.tutoUrl"
      />
    </section>
  </q-page>
</template>


<script setup>
import ProjectCard from 'src/components/ProjectCard.vue'

const projects = [
  {
    title: 'dasmlab-home',
    description: 'This Site! - Core UI + UX exploration for DASMLAB living lab site. Native Quasar/Vue.js/Vite application served with simple Nginx',
    url: 'https://github.com/lmcdasm/dasmlab_home',
    language: 'Vue.js / Quasar',
    badge: 'Public',
    liveUrl: 'https://dasmlab.org'
  },
  {
    title: 'tutorial-generator',
    description: 'A programmatic tutorial generator for dasmlab projects, blending animation, narration, and branding — all defined in YAML. Each tutorial helps explain complex devops, GitOps, or AI workflows in a visual, animated, and reusable way.',
    url: 'https://github.com/dasmlab/tutorial-generator',
    language: 'GoLang / Python / D3.js / Vue.js / NLP', 
    badge: 'Public',
    liveUrl: '',
    tutoUrl: 'https://dasmlab.org/tutorials/demo1.html' // New: to standalonetuto

  }
]
</script>
