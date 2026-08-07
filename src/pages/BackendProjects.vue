<template>
  <q-page padding class="q-gutter-md">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">Project lane</div>
        <h1 class="dasm-title">Backend projects</h1>
        <p class="dasm-subtitle">
          Services built for reliability: Go-based APIs, containerized runtime, and Kubernetes-native deployment patterns.
        </p>
      </div>
    </section>

    <div class="dasm-waypoint">Service endpoints and APIs</div>
    <section class="dasm-panel">
      <ProjectCard
        v-for="(project, index) in projects"
        :key="index"
        :title="project.title"
        :description="project.description"
        :url="project.url"
        :language="project.language"
        :badge="project.badge"
        category="Backend"
        :live-url="project.liveUrl"
      />
    </section>
  </q-page>
</template>


<script setup>
import ProjectCard from 'src/components/ProjectCard.vue'

const projects = [
  {
    title: 'design-carousel-service',
    description: 'This service provides a "carousel" of Images, providing a "queue" of images that can be read programatically (JSON) and then consumed to provide a live "feed" of images to the end user.  Allows for Adding, Deleting and updating.  This also converts all images to WEPB for optimization in transfers.',
    url: 'https://github.com/dasmlab/tutorial-generator',
    language: 'GoLang / Docker / K8s / GitActions',
    badge: 'Public',
    liveUrl: 'https://design-carousel.svc.dasmlab.org/carousel'
  }
]
</script>
