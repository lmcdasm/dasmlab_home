<template>
  <q-page padding>
    <h2>AI/ML Projects</h2>
    <p>With AI/ML technologies exploding all around us, this is where we explore models, deployment frameworks, and experimentation. Current focus includes MCP, A2A, and XAgentX stylings.</p>

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
    title: 'mcp-tour',
    description: 'A Tour of Model Context Protocol (MCP) using Go-Lang implementation. Allowing for Visual Authoring of MCP Clients, Servers and things..',
    url: 'https://github.com/lmcdasm/mcp-tour',
    language: 'MCP / GoLang / Vue.js / Quasar ',
    badge: 'Public',
    liveUrl: 'https://mcp-explorer.dasmlab.org'
  },
  {
    title: 'lumenovus',
    description: 'Fintech experimentation combining data and education for youth.',
    url: 'https://github.com/lmcdasm/lumenovus',
    language: 'Vue + Firebase',
    badge: 'Private'
  }
]
</script>

