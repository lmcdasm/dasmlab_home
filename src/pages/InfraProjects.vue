<template>
  <q-page padding>
    <h2>Infrastructure Projects</h2>
    <p>Here we dive into orchestration, CI/CD, validation frameworks and more. Think Nephio, Cluster API, Terraform, and other infrastructure plumbing tools.</p>
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
    title: 'arq1',
    description: 'Using Terraform GCP providers, GitActions and some scaffolding to demonstrate Terraform code + Git Actions techniques.',
    url: 'https://github.com/lmcdasm/arq1',
    language: ' Terraform / Smarty / Git Actions',
    badge: 'Private'
  },
]
</script>

