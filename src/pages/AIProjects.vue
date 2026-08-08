<template>
  <q-page class="lane-page">
    <header class="lane-hero">
      <div class="lane-hero__grid">
        <div>
          <div class="dasm-caps">Project lane</div>
          <h1>AI / ML projects</h1>
          <p>
            Practical AI products and MCP-native experiments — flagship metrics up top, D3 stages below (not a boring list).
          </p>
          <div class="signal-row">
            <span v-for="s in signals" :key="s">{{ s }}</span>
          </div>
        </div>
        <div class="metric-board">
          <article v-for="m in metrics" :key="m.label" class="metric">
            <div class="metric__value">{{ m.value }}</div>
            <div class="metric__label">{{ m.label }}</div>
          </article>
        </div>
      </div>
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
  title: 'AI/ML projects',
  description: 'Interview-Me, MCP tour, and AI platform experiments from Technologies DASMLAB Inc.',
  path: '/projects/ai-ml'
})

const signals = ['GenAI product paths', 'MCP-native', 'ModelOps', 'Demo-safe UX']
const metrics = [
  { value: String(5), label: 'Spotlight stages' },
  { value: '3', label: 'DASMLAB originals' },
  { value: 'Demo', label: 'Public try paths' }
]

const projects = [
  {
    title: 'Interview-Me',
    lane: 'AI/ML',
    badge: 'Demo',
    description: 'Structured knowledge-flow interviews with admin OIDC, invite guests, and synthetic public demo.',
    problem: 'Discovery interviews scatter across docs and meetings',
    outcome: 'Session engine + safe showcase without PII',
    techs: ['Go', 'Vue', 'Quasar', 'OIDC', 'SMTP'],
    accent: '#1f6f62',
    hubPath: '/projects/interview-me',
    liveUrl: 'https://interview-me.dasmlab.org/demo'
  },
  {
    title: 'mcp-tour',
    lane: 'AI/ML',
    badge: 'Public',
    description: 'Hands-on Model Context Protocol tour — clients, servers, and visual authoring.',
    problem: 'MCP stays abstract without a playground',
    outcome: 'Interactive explorer for agents & tools',
    techs: ['MCP', 'Go', 'Vue', 'Quasar'],
    accent: '#2f8f7d',
    liveUrl: 'https://mcp-explorer.dasmlab.org',
    url: 'https://github.com/lmcdasm/mcp-tour'
  },
  {
    title: 'lumenovus',
    lane: 'AI/ML',
    badge: 'Private',
    description: 'Fintech + education experiments for youth — data storytelling with care.',
    problem: 'Financial literacy tools feel corporate and cold',
    outcome: 'Youth-first narrative product experiments',
    techs: ['Vue', 'Firebase', 'data viz'],
    accent: '#74865c',
    url: 'https://github.com/lmcdasm/lumenovus'
  },
  {
    title: 'Activity → engagement intel',
    lane: 'AI/ML',
    badge: 'Lab',
    description: 'Anonymous CDP spine feeding future GEO/engagement scoring (cdn-mgr 2.x/3.0 radar).',
    problem: 'No first-party signal on public lab pages',
    outcome: 'Anon sessions ready for later ranking features',
    techs: ['Go', 'JSONL', 'UA enrich'],
    accent: '#3f9f8e',
    hubPath: '/labs/activity-anon-cdp'
  },
  {
    title: 'tutorial-generator (NLP)',
    lane: 'AI/ML',
    badge: 'Public',
    description: 'YAML + NLP assisted narration for animated teach-backs.',
    problem: 'Static docs don’t show procedure motion',
    outcome: 'Reusable narrated tutorials',
    techs: ['NLP', 'Python', 'D3', 'Vue'],
    accent: '#1f6f62',
    url: 'https://github.com/dasmlab/tutorial-generator'
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
    radial-gradient(circle at 88% 20%, rgba(47, 143, 125, 0.12), transparent 40%),
    linear-gradient(155deg, #fff, #e8f1ef);
}
.lane-hero__grid {
  display: grid;
  grid-template-columns: 1.3fr 0.7fr;
  gap: 1rem;
  align-items: center;
}
.lane-hero h1 {
  margin: 0.2rem 0 0.4rem;
  font-family: 'Fraunces', Georgia, serif;
  font-size: clamp(1.45rem, 2.5vw, 2rem);
  color: #12202c;
}
.lane-hero p {
  margin: 0 0 0.65rem;
  color: #3d5263;
  line-height: 1.5;
}
.signal-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
}
.signal-row span {
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  padding: 0.3rem 0.55rem;
  border-radius: 8px;
  border: 1px solid rgba(31, 111, 98, 0.28);
  background: rgba(255, 255, 255, 0.85);
  color: #1f6f62;
}
.metric-board {
  display: grid;
  gap: 0.55rem;
}
.metric {
  border-radius: 14px;
  padding: 0.75rem 0.85rem;
  background: linear-gradient(145deg, #12202c, #1d3a36);
  color: #fff;
}
.metric__value {
  font-family: 'Fraunces', Georgia, serif;
  font-size: 1.6rem;
  font-weight: 700;
}
.metric__label {
  font-size: 0.72rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  opacity: 0.8;
}
.stage-stack {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
@media (max-width: 860px) {
  .lane-hero__grid {
    grid-template-columns: 1fr;
  }
}
</style>
