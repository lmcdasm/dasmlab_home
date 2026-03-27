<template>
  <q-page padding class="q-gutter-md">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">Project lane</div>
        <h1 class="dasm-title">Security projects</h1>
        <p class="dasm-subtitle">
          Secure-by-design experiments across SAST/DAST, cluster hardening, and practical AI-era defense patterns.
        </p>
      </div>
    </section>

    <div class="dasm-waypoint">Trust is a feature</div>
    <section class="dasm-panel">
      <ProjectCard
        v-for="(project, index) in projects"
        :key="index"
        :title="project.title"
        :description="project.description"
        :url="project.url"
        :language="project.language"
        :badge="project.badge"
        category="Security"
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
    title: 'security-suite',
    description: 'A portable DAST/SAST container with an assortment of tooling to allow for GitActions scanning (code and Application) as well as running within a K8s cluster',
    url: 'https://github.com/dasmlab/security-suite',
    language: 'Trivy / SemGrep / Nikto / Docker / GitActions',
    badge: 'Public',
  }
]
</script>

