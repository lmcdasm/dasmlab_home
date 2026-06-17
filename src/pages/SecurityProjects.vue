<template>
  <q-page padding class="q-gutter-md security-page">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">Project lane</div>
        <h1 class="dasm-title">Security projects</h1>
        <p class="dasm-subtitle">
          Secure-by-design experiments across SAST/DAST, cluster hardening, and practical AI-era defense patterns.
        </p>
      </div>
    </section>

    <section class="security-scene dasm-panel">
      <div class="security-scene-copy">
        <div class="security-caps">Defense in motion</div>
        <h2 class="security-title">Secure systems should look and feel alive.</h2>
        <p class="security-copy">
          This lane blends shift-left scanning, runtime checks, and cluster-era hardening patterns into practical release workflows.
          The goal is not only finding issues early, but continuously proving trust through every deploy hop.
        </p>
        <div class="security-pills">
          <span v-for="pill in focusPills" :key="pill" class="security-pill">
            {{ pill }}
          </span>
        </div>
      </div>
      <div class="security-scene-signal" aria-hidden="true">
        <div class="security-signal-card">
          <div class="security-signal-kicker">Live posture</div>
          <div class="security-signal-title">Pipeline guardrails</div>
          <ul class="security-signal-list">
            <li>Static + dynamic scan gates</li>
            <li>Container image policy checks</li>
            <li>Cluster workload verification</li>
          </ul>
        </div>
      </div>
    </section>

    <div class="dasm-waypoint">Trust is a feature</div>
    <section class="security-card-rack dasm-panel">
      <ProjectCard
        v-for="(project, index) in projects"
        :key="index"
        :title="project.title"
        :description="project.description"
        :url="project.url"
        :language="project.language"
        :badge="project.badge"
        category="Security"
        :live-url="project.liveUrl"
        :tuto-url="project.tutoUrl"
      />
    </section>
  </q-page>
</template>

<script setup>
import ProjectCard from 'src/components/ProjectCard.vue'

const focusPills = [
  'SAST / DAST',
  'Supply chain checks',
  'Runtime hardening',
  'Policy as code'
]

const projects = [
  {
    title: 'security-suite',
    description: 'A portable DAST/SAST container with an assortment of tooling to allow for GitActions scanning (code and Application) as well as running within a K8s cluster',
    url: 'https://github.com/dasmlab/security-suite',
    language: 'Trivy / SemGrep / Nikto / Docker / GitActions',
    badge: 'Public',
  }
]
</script>

<style scoped>
.security-page {
  overflow-x: clip;
}

.security-scene {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(44, 78, 101, 0.2);
  background: linear-gradient(155deg, rgba(251, 254, 255, 0.98), rgba(240, 248, 254, 0.96));
  box-shadow: 0 16px 32px rgba(22, 43, 61, 0.09);
  display: grid;
  grid-template-columns: 1.08fr 0.92fr;
  gap: 0.95rem;
}

.security-scene::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: url('/media/hero/security-ambient.svg');
  background-size: cover;
  background-position: center;
  opacity: 0.36;
  pointer-events: none;
}

.security-scene-copy,
.security-scene-signal {
  position: relative;
  z-index: 1;
}

.security-scene-copy {
  padding: 1.2rem 1.24rem;
}

.security-caps {
  display: inline-flex;
  align-items: center;
  gap: 0.34rem;
  font-size: 0.72rem;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #557488;
  border: 1px solid rgba(89, 120, 143, 0.26);
  border-radius: 999px;
  padding: 0.22rem 0.62rem;
  background: rgba(255, 255, 255, 0.74);
}

.security-title {
  margin: 0.6rem 0 0.52rem;
  color: #1f3547;
  font-size: clamp(1.22rem, 2.2vw, 1.6rem);
  line-height: 1.2;
}

.security-copy {
  margin: 0;
  color: #4f6475;
  line-height: 1.66;
  font-size: 0.97rem;
  max-width: 58ch;
}

.security-pills {
  margin-top: 0.86rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
}

.security-pill {
  border-radius: 999px;
  border: 1px solid rgba(83, 120, 146, 0.28);
  background: linear-gradient(150deg, rgba(223, 239, 251, 0.7), rgba(246, 252, 255, 0.88));
  color: #3b5c75;
  font-size: 0.74rem;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  font-weight: 600;
  padding: 0.26rem 0.56rem;
}

.security-scene-signal {
  padding: 1.15rem 1.08rem 1.08rem 0.55rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.security-signal-card {
  width: min(100%, 360px);
  border-radius: 16px;
  border: 1px solid rgba(58, 89, 112, 0.22);
  background: linear-gradient(160deg, rgba(31, 56, 76, 0.88), rgba(39, 75, 102, 0.84));
  box-shadow: 0 16px 28px rgba(14, 28, 40, 0.26);
  padding: 0.9rem 0.92rem;
}

.security-signal-kicker {
  color: rgba(200, 228, 247, 0.86);
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.11em;
}

.security-signal-title {
  margin-top: 0.36rem;
  color: #ecf7ff;
  font-weight: 700;
  font-size: 1.03rem;
}

.security-signal-list {
  margin: 0.68rem 0 0;
  padding-left: 1rem;
  color: rgba(220, 241, 255, 0.9);
  line-height: 1.58;
  font-size: 0.86rem;
}

.security-card-rack {
  border: 1px solid rgba(42, 74, 97, 0.18);
  background: linear-gradient(170deg, rgba(255, 255, 255, 0.95), rgba(247, 252, 255, 0.9));
}

.security-card-rack :deep(.project-card) {
  border-color: rgba(58, 99, 128, 0.2);
}

@media (max-width: 980px) {
  .security-scene {
    grid-template-columns: 1fr;
    gap: 0.45rem;
  }

  .security-scene-signal {
    padding: 0.2rem 1rem 1.06rem;
    justify-content: flex-start;
  }
}
</style>

