<template>
  <q-page padding class="q-gutter-md">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">Project lane</div>
        <h1 class="dasm-title">Infrastructure projects</h1>
        <p class="dasm-subtitle">
          Automation, GitOps, Kubernetes, and deployment foundations that power the lab and production paths.
        </p>
      </div>
    </section>

    <div class="dasm-grid dasm-grid--2">
      <div class="dasm-panel">
        <img
          src="/dasmlab_cdevelop_foundation.png"
          alt="Foundational Principles"
          class="infra-img clickable"
          @click="openModal('/dasmlab_cdevelop_foundation.png')"
        />
      </div>
      <div class="dasm-panel infra-text">
        <p>
          DASMLAB projects run in a SecDevOps workflow. Every release is built, tested, and published through GitHub Actions,
          then promoted by manifest into a GitOps control repository.
        </p>
        <p>
          Argo CD GitOps patterns keep deployments versioned, reproducible, and auditable while remaining easy to inspect.
        </p>
      </div>
    </div>

    <div class="dasm-grid dasm-grid--2">
      <div class="dasm-panel infra-text">
        <p>
          The environment uses K3s with Calico and MetalLB, plus HAProxy/firewall edge controls and Grafana monitoring.
        </p>
        <p>
          A self-hosted GitHub runner and production cluster VMs share a controlled virtual switch, creating a practical sandbox-to-prod path.
        </p>
      </div>
      <div class="dasm-panel">
        <img
          src="/dasmlab_cdevelop_pipeline_overview.png"
          alt="Infrastructure Overview"
          class="infra-img clickable"
          @click="openModal('/dasmlab_cdevelop_pipeline_overview.png')"
        />
      </div>
    </div>

    <q-dialog v-model="modalOpen" persistent>
      <div class="infra-modal">
        <img
          v-if="modalImage"
          :src="modalImage"
          alt="Preview"
          class="infra-modal-image"
        />
        <q-btn
          flat
          round
          dense
          icon="close"
          color="white"
          @click="modalOpen = false"
          class="infra-modal-close"
        />
      </div>
    </q-dialog>
    <div class="dasm-waypoint">Deployment systems and platform plumbing</div>
    <section class="dasm-panel">
      <ProjectCard
        v-for="(project, index) in projects"
        :key="index"
        :title="project.title"
        :description="project.description"
        :url="project.url"
        :language="project.language"
        :badge="project.badge"
        category="Infrastructure"
      />
    </section>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import ProjectCard from 'src/components/ProjectCard.vue'

const modalOpen = ref(false)
const modalImage = ref('')

function openModal(src) {
  modalImage.value = src
  modalOpen.value = true
}

const projects = [
  {
    title: 'Circle-Ci',
    description: 'Using Circle CI config.yml and a K8s Deployed container-agent to provide CX Pipelines Activities for a FastAPI, Python based application. Includes Walkthrough PPT',
    url: 'https://github.com/dasmlab/mcp-server-agent-list-service',
    language: 'CircleCI / FastAPI (Python) / Buildah',
    badge: 'Public'
  },
  {
    title: 'arq1',
    description: 'Using Terraform GCP providers, GitActions and some scaffolding to demonstrate Terraform code + Git Actions techniques.',
    url: 'https://github.com/lmcdasm/arq1',
    language: 'Terraform / Smarty / Git Actions',
    badge: 'Private'
  },
]
</script>

<style scoped>
.infra-img {
  max-width: 100%;
  width: 100%;
  border-radius: 14px;
  border: 1px solid rgba(41, 72, 99, 0.16);
  box-shadow: 0 8px 20px rgba(25, 47, 67, 0.08);
  cursor: pointer;
  transition: box-shadow 170ms ease, transform 170ms ease;
}

.infra-img.clickable:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 24px rgba(63, 122, 107, 0.18);
}

.infra-text {
  color: #4b5d6d;
  line-height: 1.58;
}

.infra-modal {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
  background: rgba(20, 31, 40, 0.92);
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.infra-modal-image {
  max-width: 88vw;
  max-height: 80vh;
  background: transparent;
}

.infra-modal-close {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 1001;
  background: rgba(0, 0, 0, 0.36);
}
</style>

