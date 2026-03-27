<template>
  <q-page padding class="q-gutter-md">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">Project lane</div>
        <h1 class="dasm-title">Cloud provider projects</h1>
        <p class="dasm-subtitle">
          Serverless, platform automation, and cost-aware patterns across AWS, Azure, and GCP.
        </p>
      </div>
    </section>
    <div class="dasm-waypoint">Provider-by-provider map</div>

    <div class="row q-col-gutter-lg">
      <div class="col-12 col-md-4 dasm-panel">
        <h4 class="q-mb-md">AWS Projects</h4>
        <ProjectCard
          v-for="(project, idx) in awsProjects"
          :key="'aws-' + idx"
          v-bind="project"
          category="AWS"
        />
      </div>
      <div class="col-12 col-md-4 dasm-panel">
        <h4 class="q-mb-md">Azure Projects</h4>
        <ProjectCard
          v-for="(project, idx) in azureProjects"
          :key="'azure-' + idx"
          v-bind="project"
          category="Azure"
        />
      </div>
      <div class="col-12 col-md-4 dasm-panel">
        <h4 class="q-mb-md">GCP Projects</h4>
        <ProjectCard
          v-for="(project, idx) in gcpProjects"
          :key="'gcp-' + idx"
          v-bind="project"
          category="GCP"
        />
      </div>
    </div>
  </q-page>
</template>

<script setup>
import ProjectCard from 'src/components/ProjectCard.vue'

// Static AWS projects
const awsProjects = [
  {
    title: 'AWS Lambda Demo',
    description: 'A serverless API example using AWS Lambda + API Gateway.',
    url: 'https://github.com/dasmlab/aws-lambda-demo',
    language: 'AWS Lambda / Node.js',
    badge: 'Public'
  },
  {
    title: 'ECS CI Pipeline',
    description: 'Build, test, and deploy Docker containers to ECS with GitHub Actions.',
    url: '',
    language: 'ECS / GitHub Actions',
    badge: 'Coming soon'
  }
]

// Static Azure projects
const azureProjects = [
  {
    title: 'Azure Function Quickstart',
    description: 'Deploy Python serverless workloads using Azure Functions and Durable Tasks.',
    url: '',
    language: 'Azure Functions / Python',
    badge: 'Coming soon'
  }
]

// Static GCP projects
const gcpProjects = [
  {
    title: 'GKE Node Pool Toolkit',
    description: 'Utility for managing GKE clusters, auto-scaling node pools, and network policies.',
    url: '',
    language: 'GKE / Golang',
    badge: 'Coming soon'
  },
  {
    title: 'BigQuery Cost Explorer',
    description: 'Visualize and optimize BigQuery usage and cost.',
    url: '',
    language: 'BigQuery / React',
    badge: 'Coming soon'
  }
]
</script>

