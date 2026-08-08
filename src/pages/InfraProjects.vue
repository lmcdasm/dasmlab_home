<template>
  <q-page class="lane-page">
    <header class="lane-hero">
      <div class="dasm-caps">Project lane</div>
      <h1>Infrastructure projects</h1>
      <p>
        GitOps, MetalLB, etcd lab tooling — visual stages first. Diagrams stay readable (no giant cropped circles).
      </p>
    </header>

    <section class="infra-band">
      <article class="band-card">
        <div class="dasm-caps">SecDevOps</div>
        <h2>Build → sign → GitOps handoff</h2>
        <p>GitHub Actions + Argo CD path from commit to live envelope — practical, not theatrical.</p>
        <img
          class="band-img"
          src="/dasmlab_cdevelop_foundation.png"
          alt="Build and release backbone diagram"
        />
      </article>
      <article class="band-card">
        <div class="dasm-caps">Cluster fabric</div>
        <h2>K3s · Calico · MetalLB · observability</h2>
        <p>Lab topology that mirrors production runner paths without drowning the page in cropped circles.</p>
        <img
          class="band-img"
          src="/dasmlab_cdevelop_pipeline_overview.png"
          alt="Infrastructure topology diagram"
        />
      </article>
    </section>

    <div class="stage-stack">
      <ProjectStage v-for="p in projects" :key="p.title" v-bind="p" />
    </div>
  </q-page>
</template>

<script setup>
import ProjectStage from 'src/components/ProjectStage.vue'
import { useSeo } from 'src/composables/useSeo'

useSeo({
  title: 'Infrastructure projects',
  description: 'GitOps, MetalLB, etcd lab tooling, and deployment foundations from DASMLAB.',
  path: '/projects/infrastructure'
})

const projects = [
  {
    title: 'camera-scrape / Live Cams',
    lane: 'Infrastructure',
    badge: 'Live',
    description: 'Cheap real-time scrape and timelapse from field cams — multi-source capture without SaaS VMS.',
    problem: 'Keep wind / snow / garden cams on a budget stack',
    outcome: 'Snapshots + timelapses via camera-scrape OpenShift route (vanity additive)',
    techs: ['Go', 'OpenShift', 'object storage', 'cron'],
    accent: '#1f6f62',
    hubPath: '/projects/camera-scrape',
    liveUrl: 'https://camera-scrape.apps.2026-prod-1.ocp.dasmlab.org/'
  },
  {
    title: 'etcd synthetic load',
    lane: 'Infrastructure',
    badge: 'Lab',
    description: 'Stress tooling with hard operator gates — demo is dry-run only.',
    problem: 'Need to rehearse etcd pressure safely',
    outcome: 'Dry-run demo + flagged real load path',
    techs: ['Go', 'etcd', 'K8s', 'dry-run'],
    accent: '#1f6f62',
    hubPath: '/labs/demo-visitor-facade'
  },
  {
    title: 'MetalLB in the lab',
    lane: 'Infrastructure',
    badge: 'Hub',
    description: 'Bare-metal LoadBalancer patterns across OpenShift / K3s labs.',
    problem: 'No cloud LB on bare metal',
    outcome: 'Documented MetalLB + BGP playbook',
    techs: ['MetalLB', 'BGP', 'K3s'],
    accent: '#2f8f7d',
    hubPath: '/topics/metallb'
  },
  {
    title: 'Circle-CI agent path',
    lane: 'Infrastructure',
    badge: 'Public',
    description: 'CircleCI + K8s container-agent pipelines for FastAPI workloads.',
    problem: 'CI agents and cluster deploy need a clear MoP',
    outcome: 'Reusable pipeline + walkthrough',
    techs: ['CircleCI', 'FastAPI', 'Buildah'],
    accent: '#74865c',
    url: 'https://github.com/dasmlab/mcp-server-agent-list-service'
  },
  {
    title: 'arq1 Terraform lab',
    lane: 'Infrastructure',
    badge: 'Private',
    description: 'Terraform GCP providers + GitHub Actions scaffolding.',
    problem: 'IaC examples stay toy-shaped',
    outcome: 'Lab-grade Terraform + Actions pattern',
    techs: ['Terraform', 'GCP', 'GitHub Actions'],
    accent: '#3f9f8e',
    url: 'https://github.com/lmcdasm/arq1'
  },
  {
    title: 'dasmlab-live-cicd',
    lane: 'Infrastructure',
    badge: 'GitOps',
    description: 'Cluster overlays for the constellation — secrets stay private.',
    problem: 'Multi-app deploy without a single source of truth',
    outcome: 'Argo-synced envelopes per product',
    techs: ['Argo CD', 'OpenShift', 'Kustomize'],
    accent: '#12202c',
    hubPath: '/topics/openshift'
  }
]
</script>

<style scoped>
.lane-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0.75rem clamp(0.7rem, 2vw, 1.2rem) 2rem;
}
.lane-hero {
  border: 1.5px solid rgba(31, 111, 98, 0.32);
  border-radius: 20px;
  padding: 1.15rem 1.25rem;
  margin-bottom: 1rem;
  background: linear-gradient(155deg, #fff, #eef5f2);
}
.lane-hero h1 {
  margin: 0.2rem 0 0.4rem;
  font-family: 'Fraunces', Georgia, serif;
  font-size: clamp(1.45rem, 2.5vw, 2rem);
  color: #12202c;
}
.lane-hero p {
  margin: 0;
  color: #3d5263;
  line-height: 1.5;
}
.infra-band {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.85rem;
  margin-bottom: 1rem;
}
.band-card {
  border: 1.5px solid rgba(31, 111, 98, 0.22);
  border-radius: 16px;
  padding: 0.95rem;
  background: #fff;
  box-shadow: 0 10px 22px rgba(18, 40, 52, 0.06);
}
.band-card h2 {
  margin: 0.25rem 0 0.35rem;
  font-family: 'Fraunces', Georgia, serif;
  font-size: 1.15rem;
  color: #12202c;
}
.band-card p {
  margin: 0 0 0.65rem;
  color: #4a5d6d;
  font-size: 0.9rem;
  line-height: 1.45;
}
.band-img {
  width: 100%;
  height: auto;
  max-height: 220px;
  object-fit: contain;
  object-position: center;
  border-radius: 12px;
  background: #f3f8f6;
  border: 1px solid rgba(28, 52, 73, 0.12);
  display: block;
}
.stage-stack {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
@media (max-width: 860px) {
  .infra-band {
    grid-template-columns: 1fr;
  }
}
</style>
