<template>
  <q-page padding class="q-gutter-md cloud-page">
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

    <section class="cloud-stage">
      <article
        v-for="(lane, laneIndex) in providerLanes"
        :key="lane.key"
        class="cloud-lane dasm-panel"
        :class="[
          `cloud-lane--${lane.key}`,
          { 'cloud-lane--offset': laneIndex === 1 }
        ]"
      >
        <header class="cloud-lane-head">
          <div class="cloud-lane-emblem">
            <q-icon :name="lane.icon" size="18px" />
          </div>
          <div>
            <h3 class="cloud-lane-title">{{ lane.title }}</h3>
            <p class="cloud-lane-subtitle">{{ lane.subtitle }}</p>
          </div>
        </header>

        <div class="cloud-lane-cards">
          <ProjectCard
            v-for="(project, idx) in lane.projects"
            :key="`${lane.key}-${idx}`"
            v-bind="project"
            :category="lane.category"
          />
        </div>
      </article>
    </section>
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

const providerLanes = [
  {
    key: 'aws',
    title: 'AWS Projects',
    subtitle: 'Serverless and container delivery paths tuned for resilient release velocity.',
    icon: 'cloud_done',
    category: 'AWS',
    projects: awsProjects
  },
  {
    key: 'azure',
    title: 'Azure Projects',
    subtitle: 'Functions-first automation with strong workflow orchestration patterns.',
    icon: 'account_tree',
    category: 'Azure',
    projects: azureProjects
  },
  {
    key: 'gcp',
    title: 'GCP Projects',
    subtitle: 'Cluster operations and cost-visibility tooling for practical platform ownership.',
    icon: 'storage',
    category: 'GCP',
    projects: gcpProjects
  }
]
</script>

<style scoped>
.cloud-page {
  overflow-x: clip;
}

.cloud-stage {
  width: min(1240px, 100%);
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 1rem;
}

.cloud-lane {
  grid-column: span 4;
  min-width: 0;
  padding: 1rem 1rem 1.05rem;
  border: 1px solid rgba(41, 72, 99, 0.16);
  background: linear-gradient(165deg, rgba(255, 255, 255, 0.97), rgba(246, 251, 255, 0.9));
  box-shadow: 0 14px 30px rgba(23, 43, 60, 0.08);
  position: relative;
  overflow: hidden;
}

.cloud-lane::after {
  content: '';
  position: absolute;
  width: 210px;
  height: 210px;
  border-radius: 999px;
  top: -140px;
  right: -105px;
  background: radial-gradient(circle, rgba(120, 163, 193, 0.24), transparent 70%);
  pointer-events: none;
}

.cloud-lane--offset {
  transform: translateY(16px);
}

.cloud-lane-head {
  display: flex;
  gap: 0.72rem;
  align-items: flex-start;
  margin-bottom: 0.86rem;
}

.cloud-lane-emblem {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(61, 99, 130, 0.24);
  color: #355774;
  background: linear-gradient(155deg, rgba(216, 235, 247, 0.95), rgba(239, 248, 253, 0.8));
  flex-shrink: 0;
}

.cloud-lane-title {
  margin: 0;
  color: #253a4d;
  font-size: clamp(1.15rem, 1.8vw, 1.34rem);
  line-height: 1.25;
}

.cloud-lane-subtitle {
  margin: 0.34rem 0 0;
  color: #617587;
  font-size: 0.86rem;
  line-height: 1.48;
}

.cloud-lane-cards {
  display: grid;
  gap: 0.72rem;
}

.cloud-lane--aws :deep(.project-card) {
  border-color: rgba(57, 103, 141, 0.2);
}

.cloud-lane--azure :deep(.project-card) {
  border-color: rgba(84, 126, 156, 0.2);
}

.cloud-lane--gcp :deep(.project-card) {
  border-color: rgba(65, 118, 129, 0.2);
}

@media (max-width: 1180px) {
  .cloud-lane {
    grid-column: span 6;
  }

  .cloud-lane--offset {
    transform: none;
  }
}

@media (max-width: 780px) {
  .cloud-stage {
    gap: 0.82rem;
  }

  .cloud-lane {
    grid-column: span 12;
    padding: 0.86rem 0.82rem 0.92rem;
  }
}
</style>

