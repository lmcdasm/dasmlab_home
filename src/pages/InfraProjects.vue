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

    <section class="infra-journey">
      <article
        v-for="(item, index) in infraNarratives"
        :key="item.title"
        class="infra-rhythm"
        :class="{ 'infra-rhythm--right': index % 2 === 1 }"
      >
        <div class="infra-node" aria-hidden="true">
          {{ String(index + 1).padStart(2, '0') }}
        </div>

        <div class="infra-media-wrap dasm-panel">
          <button
            type="button"
            class="infra-media-orb"
            :aria-label="`Open ${item.title} diagram`"
            @click="openModal(item.image)"
          >
            <img
              :src="item.image"
              :alt="item.alt"
              class="infra-img"
            />
          </button>
          <p class="infra-media-caption">{{ item.caption }}</p>
        </div>

        <div class="infra-copy dasm-panel">
          <span class="infra-chip">{{ item.tag }}</span>
          <h3 class="infra-copy-title">{{ item.title }}</h3>
          <p class="infra-copy-main">{{ item.main }}</p>
          <p class="infra-copy-sub">{{ item.sub }}</p>
        </div>
      </article>
    </section>

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

const infraNarratives = [
  {
    title: 'Build and release backbone',
    tag: 'SecDevOps flow',
    image: '/dasmlab_cdevelop_foundation.png',
    alt: 'Foundational Principles',
    caption: 'Commit -> CI checks -> signed artifact -> GitOps manifest handoff',
    main: 'DASMLAB projects move through a SecDevOps lane where every release is built, tested, and published through GitHub Actions.',
    sub: 'Versioned manifests are promoted into the GitOps control repository so Argo CD can keep runtime state reproducible and easy to audit.'
  },
  {
    title: 'Runtime and operations fabric',
    tag: 'Cluster topology',
    image: '/dasmlab_cdevelop_pipeline_overview.png',
    alt: 'Infrastructure Overview',
    caption: 'K3s + network controls + observability + production runner integration',
    main: 'The environment combines K3s with Calico and MetalLB, edge controls through HAProxy/firewall policy, and Grafana-based monitoring.',
    sub: 'A self-hosted GitHub runner and production cluster VMs share a controlled virtual switch, creating a practical sandbox-to-prod path.'
  }
]

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
.infra-journey {
  position: relative;
  display: grid;
  gap: 1.4rem;
}

.infra-rhythm {
  position: relative;
  display: grid;
  grid-template-columns: 52px minmax(300px, 0.95fr) minmax(300px, 1.05fr);
  align-items: center;
  gap: 1.1rem;
}

.infra-rhythm:nth-child(odd) .infra-media-wrap {
  transform: translateY(-6px);
}

.infra-rhythm:nth-child(even) .infra-copy {
  transform: translateY(8px);
}

.infra-rhythm--right .infra-media-wrap {
  grid-column: 3;
}

.infra-rhythm--right .infra-copy {
  grid-column: 2;
}

.infra-node {
  width: 46px;
  height: 46px;
  border-radius: 999px;
  border: 1px solid rgba(55, 90, 117, 0.28);
  background: linear-gradient(150deg, rgba(227, 241, 250, 0.9), rgba(213, 232, 245, 0.68));
  box-shadow: 0 9px 20px rgba(30, 56, 74, 0.08);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #33526a;
  font-size: 0.75rem;
  letter-spacing: 0.08em;
  font-weight: 700;
}

.infra-media-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.7rem;
}

.infra-media-orb {
  width: min(100%, 370px);
  aspect-ratio: 1 / 1;
  border: 1px solid rgba(40, 72, 96, 0.2);
  border-radius: 999px;
  background: radial-gradient(circle at 35% 25%, rgba(255, 255, 255, 0.96), rgba(230, 242, 250, 0.86));
  box-shadow: 0 16px 30px rgba(24, 45, 62, 0.12);
  cursor: pointer;
  padding: 1.35rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 180ms ease, box-shadow 180ms ease;
}

.infra-media-orb:hover {
  transform: translateY(-3px) scale(1.01);
  box-shadow: 0 18px 36px rgba(34, 74, 106, 0.18);
}

.infra-img {
  max-width: 100%;
  max-height: 100%;
  width: 100%;
  object-fit: contain;
  border-radius: 14px;
}

.infra-media-caption {
  margin: 0;
  text-align: center;
  color: #607588;
  font-size: 0.79rem;
  line-height: 1.45;
  max-width: 32ch;
}

.infra-copy {
  padding: 1.12rem 1.18rem;
  border: 1px solid rgba(42, 75, 102, 0.15);
  background: linear-gradient(165deg, rgba(255, 255, 255, 0.94), rgba(245, 251, 255, 0.9));
}

.infra-chip {
  display: inline-flex;
  align-items: center;
  padding: 0.24rem 0.58rem;
  border-radius: 999px;
  border: 1px solid rgba(90, 126, 153, 0.3);
  color: #4f6a7f;
  font-size: 0.7rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-weight: 600;
  margin-bottom: 0.58rem;
}

.infra-copy-title {
  margin: 0 0 0.5rem;
  color: #263b4e;
  font-size: clamp(1.05rem, 1.8vw, 1.35rem);
  line-height: 1.24;
}

.infra-copy-main,
.infra-copy-sub {
  margin: 0;
  color: #4c6070;
  line-height: 1.62;
  font-size: 0.98rem;
}

.infra-copy-sub {
  margin-top: 0.62rem;
  color: #5f7282;
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

@media (max-width: 1220px) {
  .infra-rhythm {
    grid-template-columns: 44px minmax(270px, 1fr) minmax(270px, 1fr);
    gap: 0.9rem;
  }

  .infra-media-orb {
    width: min(100%, 330px);
  }
}

@media (max-width: 980px) {
  .infra-rhythm,
  .infra-rhythm--right {
    grid-template-columns: 1fr;
    gap: 0.85rem;
  }

  .infra-rhythm:nth-child(odd) .infra-media-wrap,
  .infra-rhythm:nth-child(even) .infra-copy {
    transform: none;
  }

  .infra-rhythm--right .infra-media-wrap,
  .infra-rhythm--right .infra-copy {
    grid-column: auto;
  }

  .infra-node {
    width: 40px;
    height: 40px;
    justify-self: start;
  }

  .infra-media-wrap {
    align-items: flex-start;
  }

  .infra-media-caption {
    text-align: left;
    max-width: 48ch;
  }
}
</style>

