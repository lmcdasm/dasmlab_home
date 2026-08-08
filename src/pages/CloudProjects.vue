<template>
  <q-page class="lane-page">
    <header class="lane-hero">
      <div class="dasm-caps">Project lane</div>
      <h1>Cloud provider projects</h1>
      <p>
        Cost-aware origins, orchestration, and multi-cloud patterns — standardized visual stages (same language as
        Frontend/Backend), not a mismatched card soup.
      </p>
    </header>

    <div class="provider-rail" aria-label="Providers">
      <button
        v-for="p in providers"
        :key="p.id"
        type="button"
        class="provider-chip"
        :class="{ active: filter === p.id }"
        @click="filter = p.id"
      >
        {{ p.label }}
      </button>
    </div>

    <div class="stage-stack">
      <ProjectStage v-for="p in visible" :key="p.title" v-bind="p" />
    </div>
  </q-page>
</template>

<script setup>
import { computed, ref } from 'vue'
import ProjectStage from 'src/components/ProjectStage.vue'
import { useSeo } from 'src/composables/useSeo'

useSeo({
  title: 'Cloud projects',
  description: 'CheapCloud, Mock-Me, and multi-cloud lab patterns from Technologies DASMLAB Inc.',
  path: '/projects/cloud'
})

const filter = ref('all')
const providers = [
  { id: 'all', label: 'All' },
  { id: 'aws', label: 'AWS' },
  { id: 'azure', label: 'Azure' },
  { id: 'gcp', label: 'GCP' },
  { id: 'multi', label: 'Multi' }
]

const projects = [
  {
    provider: 'multi',
    title: 'CheapCloud',
    lane: 'Cloud',
    badge: 'Demo',
    description: 'Spend envelopes and provider recommendations — demo is readonly fixtures.',
    problem: 'Picking cheap origins is tribal knowledge',
    outcome: 'Recommend + envelope UX without credentials in demo',
    techs: ['Go', 'Gin', 'AWS', 'Azure', 'GCP'],
    accent: '#1f6f62',
    hubPath: '/projects/cheapcloud',
    liveUrl: 'https://cheapcloud.dasmlab.org/demo'
  },
  {
    provider: 'aws',
    title: 'Mock-Me orchestration',
    lane: 'Cloud',
    badge: 'Demo',
    description: 'CDN/platform workflow orchestration — fake mode never deploys live nodes.',
    problem: 'Showcase vs production deploy collision',
    outcome: 'Scripted assembly line for visitors',
    techs: ['Go', 'Vue', 'OpenShift', 'OIDC'],
    accent: '#2f8f7d',
    hubPath: '/projects/mock-me',
    liveUrl: 'https://mock-me.dasmlab.org/demo'
  },
  {
    provider: 'aws',
    title: 'AWS Lambda Demo',
    lane: 'Cloud',
    badge: 'Public',
    description: 'Serverless API example using Lambda + API Gateway.',
    problem: 'Need a minimal serverless reference',
    outcome: 'Public demo repo + pattern notes',
    techs: ['Lambda', 'API Gateway', 'Node.js'],
    accent: '#74865c',
    url: 'https://github.com/dasmlab/aws-lambda-demo'
  },
  {
    provider: 'azure',
    title: 'CheapCloud Azure adapters',
    lane: 'Cloud',
    badge: 'Lab',
    description: 'Azure-shaped dry-run and cost envelope patterns.',
    problem: 'Azure SKUs obscure free-tier burn',
    outcome: 'Adapter layer for recommend path',
    techs: ['Azure', 'Go', 'dry-run'],
    accent: '#1f6f62',
    hubPath: '/projects/cheapcloud'
  },
  {
    provider: 'azure',
    title: 'OpenShift lab (Azure-adjacent)',
    lane: 'Cloud',
    badge: 'Hub',
    description: 'GitOps envelopes and origin broker notes for lab clusters.',
    problem: 'Cluster ops docs scatter across repos',
    outcome: 'Topic hub for OpenShift patterns',
    techs: ['OpenShift', 'GitOps', 'Argo'],
    accent: '#2f8f7d',
    hubPath: '/topics/openshift'
  },
  {
    provider: 'gcp',
    title: 'Surfing object origin',
    lane: 'Cloud',
    badge: 'Lab',
    description: 'Move Surfing bytes off PVC to object storage + CDN.',
    problem: 'Basement disk as origin does not scale',
    outcome: 'Manifest CDN URLs + cheap origin path',
    techs: ['Object storage', 'CDN', 'Go'],
    accent: '#3f9f8e',
    hubPath: '/labs/surfing-r2-origin'
  },
  {
    provider: 'multi',
    title: 'dasmlab-cdn-mgr (radar)',
    lane: 'Cloud',
    badge: '3.0',
    description: 'Future multi-realm CDN manager for GEO/SEO/engagement — see ADR-001.',
    problem: 'Per-site nginx proxies won’t carry GEO + engagement',
    outcome: 'Platform roadmap without blocking 2.0 polish',
    techs: ['CDN', 'realms', 'GEO', 'Activity'],
    accent: '#12202c',
    hubPath: '/labs/surfing-r2-origin'
  }
]

const visible = computed(() =>
  filter.value === 'all' ? projects : projects.filter((p) => p.provider === filter.value)
)
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
  margin-bottom: 0.85rem;
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
.provider-rail {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-bottom: 1rem;
}
.provider-chip {
  border: 1.5px solid rgba(31, 111, 98, 0.28);
  background: #fff;
  border-radius: 10px;
  padding: 0.4rem 0.75rem;
  font-weight: 700;
  cursor: pointer;
  color: #1d2b36;
}
.provider-chip.active {
  background: linear-gradient(135deg, #1f6f62, #2f8f7d);
  color: #fff;
  border-color: transparent;
}
.stage-stack {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>
