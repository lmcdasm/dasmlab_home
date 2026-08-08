<template>
  <q-page class="lane-page">
    <header class="lane-hero">
      <div class="lane-hero__grid">
        <div>
          <div class="dasm-caps">Defense in motion</div>
          <h1>Security projects</h1>
          <p>
            Secure-by-design experiments that feel alive — shift-left gates, OIDC boundaries, and demo facades that
            never get live privileges.
          </p>
          <div class="pill-row">
            <span v-for="p in pills" :key="p">{{ p }}</span>
          </div>
        </div>
        <aside class="posture">
          <div class="posture__kicker">Live posture</div>
          <div class="posture__title">Pipeline guardrails</div>
          <ul>
            <li>Static + dynamic scan gates</li>
            <li>Container image policy checks</li>
            <li>Demo DenyMutate on live APIs</li>
            <li>OIDC for operators only</li>
          </ul>
        </aside>
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
  title: 'Security projects',
  description: 'Security suite, OIDC patterns, and demo visitor safeguards from DASMLAB.',
  path: '/projects/security'
})

const pills = ['SAST / DAST', 'Supply chain', 'Runtime hardening', 'OIDC', 'Demo facades']

const projects = [
  {
    title: 'security-suite',
    lane: 'Security',
    badge: 'Public',
    description: 'Portable DAST/SAST container for GitHub Actions and in-cluster scanning.',
    problem: 'Scan tooling is fragmented across pipelines',
    outcome: 'One container for Trivy / Semgrep / Nikto paths',
    techs: ['Trivy', 'Semgrep', 'Nikto', 'Docker'],
    accent: '#1f6f62',
    url: 'https://github.com/dasmlab/security-suite'
  },
  {
    title: 'OIDC across the constellation',
    lane: 'Security',
    badge: 'Hub',
    description: 'Keycloak-backed live ops; demos never receive admin roles.',
    problem: 'Showcase vs operator privilege collision',
    outcome: 'Clear persona split: demo / guest / admin',
    techs: ['OIDC', 'Keycloak', 'cookies'],
    accent: '#2f8f7d',
    hubPath: '/topics/oidc'
  },
  {
    title: 'Demo visitor facade',
    lane: 'Security',
    badge: 'Lab',
    description: 'Labeled fake mode contract — simulate-only writes, DenyDemoMutate.',
    problem: 'Public interest vs IP-protected mutate paths',
    outcome: 'Cross-product demo cookie pattern',
    techs: ['authz', 'fixtures', 'Activity'],
    accent: '#74865c',
    hubPath: '/labs/demo-visitor-facade'
  },
  {
    title: 'Activity owner gate',
    lane: 'Security',
    badge: 'Live',
    description: 'Public POST for anon engagement; GET/panel dual-gated to owner viewers.',
    problem: 'Engagement data must not leak broadly',
    outcome: 'Write-open / read-gated Activity spine',
    techs: ['Go', 'allowlist', 'cookies'],
    accent: '#3f9f8e',
    hubPath: '/labs/activity-anon-cdp'
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
    radial-gradient(circle at 90% 10%, rgba(47, 143, 125, 0.12), transparent 40%),
    linear-gradient(155deg, #fff, #e8f0f4);
}
.lane-hero__grid {
  display: grid;
  grid-template-columns: 1.25fr 0.75fr;
  gap: 1rem;
  align-items: stretch;
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
.pill-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
}
.pill-row span {
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  padding: 0.3rem 0.55rem;
  border-radius: 8px;
  border: 1px solid rgba(31, 111, 98, 0.28);
  background: rgba(255, 255, 255, 0.9);
  color: #1f6f62;
}
.posture {
  border-radius: 16px;
  padding: 1rem;
  background: linear-gradient(150deg, #12202c, #1a3532);
  color: #fff;
  box-shadow: 0 12px 28px rgba(18, 40, 52, 0.2);
}
.posture__kicker {
  font-size: 0.68rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  opacity: 0.75;
}
.posture__title {
  font-family: 'Fraunces', Georgia, serif;
  font-size: 1.2rem;
  margin: 0.25rem 0 0.55rem;
}
.posture ul {
  margin: 0;
  padding-left: 1.1rem;
  line-height: 1.55;
  font-size: 0.9rem;
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
