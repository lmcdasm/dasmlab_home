<template>
  <q-page class="lane-page">
    <header class="lane-hero">
      <div class="dasm-caps">Project lane</div>
      <h1>Backend projects</h1>
      <p>
        Go APIs, cookie sessions, and Kubernetes-native runtimes — each stage diagrams the constraint we solved.
      </p>
    </header>
    <div class="stage-stack">
      <ProjectStage v-for="p in projects" :key="p.title" v-bind="p" />
    </div>
  </q-page>
</template>

<script setup>
import ProjectStage from 'src/components/ProjectStage.vue'
import { useSeo } from 'src/composables/useSeo'

useSeo({
  title: 'Backend projects',
  description: 'Go services from DASMLAB — Surfing, design carousel, Activity CDP.',
  path: '/projects/backend'
})

const projects = [
  {
    title: 'surfing-service',
    lane: 'Backend',
    badge: 'Live',
    description: 'Gin API for Surfing media days, OIDC, upload orchestration, anonymous Activity.',
    problem: 'Auth + media + engagement split across ad-hoc scripts',
    outcome: 'One service: days, CDN hooks, public Activity POST',
    techs: ['Go', 'Gin', 'OIDC', 'JSONL', 'OpenShift'],
    accent: '#1f6f62',
    hubPath: '/projects/surfing',
    liveUrl: 'https://dasmlab.org/surfing'
  },
  {
    title: 'design-carousel-service',
    lane: 'Backend',
    badge: 'Live',
    description: 'Programmatic image carousel queue with WebP conversion for lean transfers.',
    problem: 'Heavy images kill lab homepage performance',
    outcome: 'JSON feed + WebP pipeline',
    techs: ['Go', 'Docker', 'K8s', 'WebP'],
    accent: '#2f8f7d',
    liveUrl: 'https://design-carousel.svc.dasmlab.org/carousel'
  },
  {
    title: 'Activity CDP spine',
    lane: 'Backend',
    badge: 'Live',
    description: 'Anonymous + authenticated engagement events for public sites.',
    problem: 'Public sites had no visitor signal without forcing login',
    outcome: 'First-party anon IDs + owner panel',
    techs: ['Go', 'Cookies', 'JSONL', 'UA enrich'],
    accent: '#74865c',
    hubPath: '/labs/activity-anon-cdp',
    liveUrl: 'https://dasmlab.org/labs/activity-anon-cdp'
  },
  {
    title: 'Interview-Me API',
    lane: 'Backend',
    badge: 'Demo',
    description: 'Sessions, invites, OTP guests, and public demo start — SQLite-backed.',
    problem: 'Need real interview ops + safe public try',
    outcome: 'Admin OIDC + im_demo synthetic session',
    techs: ['Go', 'chi', 'SQLite', 'SMTP'],
    accent: '#3f9f8e',
    hubPath: '/projects/interview-me',
    liveUrl: 'https://interview-me.dasmlab.org/demo'
  },
  {
    title: 'Mock-Me API',
    lane: 'Backend',
    badge: 'Demo',
    description: 'MockUp CRUD + deploy engine with DenyDemoMutate and /demo/simulate.',
    problem: 'Live deploy paths must never run for visitors',
    outcome: 'Fixture timelines + Activity demo flag',
    techs: ['Go', 'chi', 'OIDC', 'deploy engine'],
    accent: '#1f6f62',
    hubPath: '/projects/mock-me',
    liveUrl: 'https://mock-me.dasmlab.org/demo'
  },
  {
    title: 'camera-scrape',
    lane: 'Backend',
    badge: 'Live',
    description: 'Multi-source scrape workers, snapshot store, and timelapse stitch for Live Cams.',
    problem: 'Heterogeneous cams without a VMS tax',
    outcome: 'Scheduled scrape → gallery + Surfing cross-link',
    techs: ['Go', 'schedulers', 'image pipeline', 'OpenShift'],
    accent: '#2f8f7d',
    hubPath: '/projects/camera-scrape',
    liveUrl: 'https://camera-scrape.apps.2026-prod-1.ocp.dasmlab.org/'
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
  background:
    radial-gradient(circle at 12% 20%, rgba(47, 143, 125, 0.14), transparent 42%),
    linear-gradient(155deg, #fff, #eef5f2);
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
  max-width: 46rem;
  line-height: 1.5;
}
.stage-stack {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>
