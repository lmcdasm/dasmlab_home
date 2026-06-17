<template>
  <q-page padding class="q-gutter-md ai-page">
    <section class="dasm-shell">
      <div class="dasm-floating-grid" />
      <div class="dasm-shell__content">
        <div class="dasm-caps">Project lane</div>
        <h1 class="dasm-title">AI/ML projects</h1>
        <p class="dasm-subtitle">
          Model experimentation, deployment patterns, and practical AI architecture across DASMLAB and cloud ecosystems.
        </p>
      </div>
    </section>
    <div class="dasm-waypoint">From prototypes to production lanes</div>

    <section class="ai-command-deck dasm-panel">
      <div class="ai-command-copy">
        <div class="ai-command-caps">Flagship lane</div>
        <h2 class="ai-command-title">Designing the next wave of practical AI products.</h2>
        <p class="ai-command-subtitle">
          This page is built as a living program map: from MCP-native tooling and applied product experiments to cloud
          model platforms that support scaling, governance, and deployment discipline.
        </p>
        <div class="ai-signal-row">
          <span v-for="signal in aiSignals" :key="signal" class="ai-signal-pill">{{ signal }}</span>
        </div>
      </div>

      <div class="ai-metric-board">
        <article class="ai-metric-card">
          <div class="ai-metric-value">{{ totalProjects }}</div>
          <div class="ai-metric-label">Current mapped projects</div>
        </article>
        <article class="ai-metric-card">
          <div class="ai-metric-value">{{ dasmlabProjects.length }}</div>
          <div class="ai-metric-label">DASMLAB originals</div>
        </article>
        <article class="ai-metric-card">
          <div class="ai-metric-value">3 Clouds</div>
          <div class="ai-metric-label">Provider acceleration lanes</div>
        </article>
      </div>
    </section>

    <section class="ai-spotlight dasm-panel">
      <header class="ai-spotlight-head">
        <div class="ai-spotlight-icon">
          <q-icon name="auto_awesome" size="18px" />
        </div>
        <div>
          <h3 class="ai-spotlight-title">DASMLAB AI/ML spotlight</h3>
          <p class="ai-spotlight-subtitle">Core initiatives that define our own platform identity and experimentation pace.</p>
        </div>
      </header>
      <div class="ai-spotlight-grid">
        <ProjectCard
          v-for="(project, idx) in dasmlabProjects"
          :key="'dasmlab-' + idx"
          v-bind="project"
          category="DASMLAB"
        />
      </div>
    </section>

    <section class="ai-provider-stage">
      <article
        v-for="(lane, laneIndex) in providerLanes"
        :key="lane.key"
        class="ai-provider-lane dasm-panel"
        :class="[
          `ai-provider-lane--${lane.key}`,
          { 'ai-provider-lane--offset': laneIndex === 1 }
        ]"
      >
        <header class="ai-provider-head">
          <div class="ai-provider-emblem">
            <q-icon :name="lane.icon" size="18px" />
          </div>
          <div>
            <h3 class="ai-provider-title">{{ lane.title }}</h3>
            <p class="ai-provider-subtitle">{{ lane.subtitle }}</p>
          </div>
        </header>
        <div class="ai-provider-cards">
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

// DASMLAB (custom) projects
const dasmlabProjects = [
  {
    title: 'mcp-tour',
    description: 'A Tour of Model Context Protocol (MCP) using Go-Lang implementation. Visual authoring of MCP clients, servers, and things.',
    url: 'https://github.com/lmcdasm/mcp-tour',
    language: 'MCP / GoLang / Vue.js / Quasar',
    badge: 'Public',
    liveUrl: 'https://mcp-explorer.dasmlab.org'
  },
  {
    title: 'lumenovus',
    description: 'Fintech experimentation combining data and education for youth.',
    url: 'https://github.com/lmcdasm/lumenovus',
    language: 'Vue + Firebase',
    badge: 'Private'
  }
]

// AWS AI/ML "Hottest"
const awsProjects = [
  {
    title: 'Amazon Bedrock',
    description: 'A fully managed service to build and scale generative AI applications with leading foundation models from AWS and partners.',
    url: 'https://aws.amazon.com/bedrock/',
    language: 'Bedrock / GenAI',
    badge: 'AWS'
  },
  {
    title: 'Amazon SageMaker',
    description: 'Complete machine learning platform for building, training, and deploying models at scale.',
    url: 'https://aws.amazon.com/sagemaker/',
    language: 'SageMaker / ML',
    badge: 'AWS'
  },
  {
    title: 'AWS Lambda AI Inference',
    description: 'Run ML inference workloads at scale serverlessly, with pre-built container images for AI/ML inference.',
    url: 'https://aws.amazon.com/lambda/features/#Machine_learning',
    language: 'Lambda / Inference',
    badge: 'AWS'
  }
]

// Azure AI/ML "Hottest"
const azureProjects = [
  {
    title: 'Azure OpenAI Service',
    description: 'Access OpenAI’s GPT-4, GPT-3, and DALL·E models through Azure APIs for secure, scalable enterprise workloads.',
    url: 'https://azure.microsoft.com/en-us/products/ai-services/openai-service/',
    language: 'Azure OpenAI',
    badge: 'Azure'
  },
  {
    title: 'Azure Machine Learning',
    description: 'Build, deploy, and manage ML models with an end-to-end MLOps platform in Azure.',
    url: 'https://azure.microsoft.com/en-us/products/machine-learning/',
    language: 'Azure ML',
    badge: 'Azure'
  },
  {
    title: 'Azure Cognitive Services',
    description: 'Prebuilt AI APIs for vision, speech, language, and decision tasks.',
    url: 'https://azure.microsoft.com/en-us/products/cognitive-services/',
    language: 'Cognitive Services',
    badge: 'Azure'
  }
]

// GCP AI/ML "Hottest"
const gcpProjects = [
  {
    title: 'Vertex AI',
    description: 'End-to-end ML platform for data scientists and ML engineers to build, deploy, and scale models on Google Cloud.',
    url: 'https://cloud.google.com/vertex-ai',
    language: 'Vertex AI',
    badge: 'GCP'
  },
  {
    title: 'Generative AI Studio',
    description: 'Rapidly prototype and deploy generative AI apps with Gemini models, embedding, text, and multimodal APIs.',
    url: 'https://cloud.google.com/vertex-ai/generative-ai/docs/overview',
    language: 'GenAI Studio / Gemini',
    badge: 'GCP'
  },
  {
    title: 'AutoML by Google',
    description: 'Train high-quality custom models with minimal ML expertise using Google Cloud AutoML.',
    url: 'https://cloud.google.com/automl',
    language: 'AutoML',
    badge: 'GCP'
  }
]

const aiSignals = [
  'GenAI product paths',
  'MCP-native integration',
  'ModelOps + observability',
  'Production-ready deployment'
]

const totalProjects = dasmlabProjects.length + awsProjects.length + azureProjects.length + gcpProjects.length

const providerLanes = [
  {
    key: 'aws',
    title: 'AWS AI/ML',
    subtitle: 'Managed foundation-model and MLOps infrastructure for scale-first deployment.',
    icon: 'cloud_done',
    category: 'AWS',
    projects: awsProjects
  },
  {
    key: 'azure',
    title: 'Azure AI/ML',
    subtitle: 'Enterprise-oriented model access and lifecycle orchestration patterns.',
    icon: 'account_tree',
    category: 'Azure',
    projects: azureProjects
  },
  {
    key: 'gcp',
    title: 'GCP AI/ML',
    subtitle: 'Unified AI platform primitives for multimodal prototyping and production.',
    icon: 'memory',
    category: 'GCP',
    projects: gcpProjects
  }
]
</script>

<style scoped>
.ai-page {
  overflow-x: clip;
}

.ai-command-deck {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(43, 74, 97, 0.2);
  background: linear-gradient(160deg, rgba(251, 254, 255, 0.98), rgba(238, 247, 254, 0.95));
  box-shadow: 0 16px 30px rgba(24, 45, 62, 0.1);
  display: grid;
  grid-template-columns: 1.08fr 0.92fr;
  gap: 1rem;
}

.ai-command-deck::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: url('/media/hero/ai-ambient.svg');
  background-size: cover;
  background-position: center;
  opacity: 0.36;
  pointer-events: none;
}

.ai-command-copy,
.ai-metric-board {
  position: relative;
  z-index: 1;
}

.ai-command-copy {
  padding: 1.2rem 1.2rem 1.14rem;
}

.ai-command-caps {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  border: 1px solid rgba(82, 119, 148, 0.24);
  background: rgba(255, 255, 255, 0.78);
  color: #4f6d84;
  padding: 0.22rem 0.62rem;
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-weight: 600;
}

.ai-command-title {
  margin: 0.6rem 0 0.5rem;
  color: #1f3549;
  font-size: clamp(1.28rem, 2.3vw, 1.68rem);
  line-height: 1.2;
}

.ai-command-subtitle {
  margin: 0;
  color: #506476;
  line-height: 1.66;
  font-size: 0.97rem;
  max-width: 60ch;
}

.ai-signal-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-top: 0.86rem;
}

.ai-signal-pill {
  border-radius: 999px;
  border: 1px solid rgba(89, 125, 152, 0.3);
  background: linear-gradient(150deg, rgba(223, 239, 251, 0.72), rgba(247, 252, 255, 0.88));
  color: #375770;
  padding: 0.24rem 0.56rem;
  font-size: 0.74rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
}

.ai-metric-board {
  padding: 1.02rem 1.08rem 1.02rem 0.42rem;
  display: grid;
  gap: 0.56rem;
}

.ai-metric-card {
  border: 1px solid rgba(51, 84, 109, 0.2);
  border-radius: 14px;
  background: linear-gradient(160deg, rgba(34, 61, 83, 0.86), rgba(42, 75, 102, 0.84));
  box-shadow: 0 14px 24px rgba(15, 29, 43, 0.22);
  padding: 0.74rem 0.84rem;
}

.ai-metric-value {
  color: #ecf6ff;
  font-size: clamp(1.14rem, 2.1vw, 1.5rem);
  line-height: 1.15;
  font-weight: 700;
}

.ai-metric-label {
  margin-top: 0.28rem;
  color: rgba(210, 234, 251, 0.88);
  font-size: 0.78rem;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.ai-spotlight {
  border: 1px solid rgba(43, 75, 100, 0.18);
  background: linear-gradient(168deg, rgba(255, 255, 255, 0.96), rgba(246, 251, 255, 0.9));
}

.ai-spotlight-head {
  display: flex;
  gap: 0.72rem;
  align-items: flex-start;
  margin-bottom: 0.8rem;
}

.ai-spotlight-icon {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  border: 1px solid rgba(69, 104, 131, 0.26);
  background: linear-gradient(155deg, rgba(217, 236, 248, 0.94), rgba(241, 250, 255, 0.8));
  color: #365875;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ai-spotlight-title {
  margin: 0;
  color: #23394c;
  font-size: clamp(1.1rem, 1.9vw, 1.36rem);
  line-height: 1.24;
}

.ai-spotlight-subtitle {
  margin: 0.34rem 0 0;
  color: #607485;
  line-height: 1.52;
  font-size: 0.87rem;
}

.ai-spotlight-grid {
  display: grid;
  gap: 0.76rem;
}

.ai-provider-stage {
  width: min(1260px, 100%);
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 1rem;
}

.ai-provider-lane {
  grid-column: span 4;
  min-width: 0;
  padding: 1rem 1rem 1.06rem;
  border: 1px solid rgba(41, 72, 99, 0.16);
  background: linear-gradient(165deg, rgba(255, 255, 255, 0.97), rgba(245, 251, 255, 0.9));
  box-shadow: 0 14px 30px rgba(23, 43, 60, 0.08);
  position: relative;
  overflow: hidden;
}

.ai-provider-lane::after {
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

.ai-provider-lane--offset {
  transform: translateY(15px);
}

.ai-provider-head {
  display: flex;
  gap: 0.72rem;
  align-items: flex-start;
  margin-bottom: 0.85rem;
}

.ai-provider-emblem {
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

.ai-provider-title {
  margin: 0;
  color: #253a4d;
  font-size: clamp(1.08rem, 1.75vw, 1.3rem);
  line-height: 1.24;
}

.ai-provider-subtitle {
  margin: 0.34rem 0 0;
  color: #617587;
  font-size: 0.86rem;
  line-height: 1.48;
}

.ai-provider-cards {
  display: grid;
  gap: 0.72rem;
}

.ai-provider-lane--aws :deep(.project-card) {
  border-color: rgba(57, 103, 141, 0.2);
}

.ai-provider-lane--azure :deep(.project-card) {
  border-color: rgba(84, 126, 156, 0.2);
}

.ai-provider-lane--gcp :deep(.project-card) {
  border-color: rgba(65, 118, 129, 0.2);
}

@media (max-width: 1180px) {
  .ai-command-deck {
    grid-template-columns: 1fr;
    gap: 0.35rem;
  }

  .ai-metric-board {
    padding: 0.15rem 1rem 1rem;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .ai-provider-lane {
    grid-column: span 6;
  }

  .ai-provider-lane--offset {
    transform: none;
  }
}

@media (max-width: 780px) {
  .ai-metric-board {
    grid-template-columns: 1fr;
  }

  .ai-provider-stage {
    gap: 0.82rem;
  }

  .ai-provider-lane {
    grid-column: span 12;
    padding: 0.84rem 0.82rem 0.9rem;
  }
}
</style>

