<template>
  <q-page class="lane-page">
    <header class="lane-hero">
      <div class="dasm-caps">Project lane</div>
      <h1>Frontend projects</h1>
      <p>
        UI craft across Vue, Quasar, motion, and accessibility — each stage shows the problem, stack, and outcome
        as a live SVG (D3), then links into Behind the Design hubs.
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
  title: 'Frontend projects',
  description: 'Vue and Quasar builds from Technologies DASMLAB Inc. — visual stages, hubs, Surfing client.',
  path: '/projects/frontend'
})

const projects = [
  {
    title: 'dasmlab-home',
    lane: 'Frontend',
    badge: 'Public',
    description: 'Engineering Knowledge Network for Technologies DASMLAB Inc. — portfolio, Surfing client, hubs, SEO spine.',
    problem: 'Card-soup portfolios don’t show systems thinking',
    outcome: 'Interactive lab map + hubs + GEO-ready structure',
    techs: ['Vue 3', 'Quasar', 'Vite', 'D3', 'nginx'],
    accent: '#1f6f62',
    hubPath: '/projects/dasmlab-home',
    liveUrl: 'https://dasmlab.org',
    url: 'https://github.com/lmcdasm/dasmlab_home'
  },
  {
    title: 'Surfing (client)',
    lane: 'Frontend',
    badge: 'Live',
    description: 'Golden CDN-backed media UX — albums, days, owner publish, public browse.',
    problem: 'Media trapped on basement PVC / opaque UX',
    outcome: 'Day-based albums with CDN-ready playback paths',
    techs: ['Vue 3', 'Quasar', 'OIDC', 'R2/CDN'],
    accent: '#2f8f7d',
    hubPath: '/projects/surfing',
    liveUrl: 'https://dasmlab.org/surfing'
  },
  {
    title: 'Interview-Me UI',
    lane: 'Frontend',
    badge: 'Demo',
    description: 'Admin + invite guest + public demo session surfaces for structured interviews.',
    problem: 'Interview tools force SSO before anyone can try',
    outcome: 'Labeled demo facade + guest OTP flows',
    techs: ['Vue 3', 'Quasar', 'OIDC', 'OTP'],
    accent: '#3f9f8e',
    hubPath: '/projects/interview-me',
    liveUrl: 'https://interview-me.dasmlab.org/demo'
  },
  {
    title: 'Mock-Me UI',
    lane: 'Frontend',
    badge: 'Demo',
    description: 'Topology / deploy assembly-line UI with fake-mode timeline for showcase.',
    problem: 'Showcase needed without live node risk',
    outcome: 'Scripted deploy stages visitors can run safely',
    techs: ['Vue 3', 'Quasar', 'SVG', 'Activity'],
    accent: '#74865c',
    hubPath: '/projects/mock-me',
    liveUrl: 'https://mock-me.dasmlab.org/demo'
  },
  {
    title: 'Live Cams (camera-scrape)',
    lane: 'Frontend',
    badge: 'Live',
    description: 'Public gallery for multi-source scrape snapshots and timelapses — field cams without a heavyweight VMS.',
    problem: 'Keep wind / snow / garden cams without SaaS camera platforms',
    outcome: 'Cheap scrape → snapshot → timelapse gallery',
    techs: ['Vue', 'Go', 'image pipeline', 'OpenShift'],
    accent: '#2a6f62',
    hubPath: '/projects/camera-scrape',
    liveUrl: 'https://camera-scrape.dasmlab.org/'
  },
  {
    title: 'tutorial-generator',
    lane: 'Frontend',
    badge: 'Public',
    description: 'YAML-defined animated teach-backs for GitOps and AI workflows.',
    problem: 'Docs alone don’t show motion of a workflow',
    outcome: 'Reusable narrated tutorial player',
    techs: ['D3', 'Vue', 'Go', 'Python', 'NLP'],
    accent: '#1f6f62',
    url: 'https://github.com/dasmlab/tutorial-generator',
    liveUrl: ''
  },
  {
    title: 'mcp-tour UI',
    lane: 'Frontend',
    badge: 'Lab',
    description: 'Visual authoring for MCP clients/servers — explore Model Context Protocol hands-on (deployment on radar).',
    problem: 'MCP concepts stay abstract without a playground',
    outcome: 'Interactive MCP tour UI (code / repo — live explorer not public yet)',
    techs: ['Vue', 'Quasar', 'MCP', 'Go'],
    accent: '#2a6f62',
    url: 'https://github.com/lmcdasm/mcp-tour'
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
