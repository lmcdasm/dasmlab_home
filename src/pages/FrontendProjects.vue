<template>
  <q-page padding>
    <h2>Frontend Projects</h2>
    <p>This is where we explore different User Interface types and technologies, from Vue.js / Quasar (my current fav) to Voice, Android/iOS, Biometric based applications. This area is where we play with the "view" from the user’s point of view.</p>
	<ProjectCard
     	 v-for="(project, index) in projects"
     	 :key="index"
     	 :title="project.title"
     	 :description="project.description"
     	 :url="project.url"
     	 :language="project.language"
     	 :badge="project.badge"
    	/>
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
    badge: 'Public'
  }
]
</script>
