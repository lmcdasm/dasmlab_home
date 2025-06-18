<template>
  <q-page padding>
    <h2>Security Projects</h2>
    <p>Explorations of SAST, DAST, K8s sidecars, and secure LLM integration. Security is fundamental — and even more so in this interconnected AI future.</p>
        <ProjectCard
         v-for="(project, index) in projects"
         :key="index"
         :title="project.title"
         :description="project.description"
         :url="project.url"
         :language="project.language"
         :badge="project.badge"
         :liveUrl="project.liveUrl"
         :tutoUrl="project.tutoUrl"
        />
  </q-page>
</template>

<script setup>
import ProjectCard from 'src/components/ProjectCard.vue'

const projects = [
  {
    title: 'security-suite',
    description: 'A portable DAST/SAST container with an assortment of tooling to allow for GitActions scanning (code and Application) as well as running within a K8s cluster',
    url: 'https://github.com/dasmlab/security-suite',
    language: 'Trivy / SemGrep / Nikto / Docker / GitActions',
    badge: 'Public',
  }
]
</script>

