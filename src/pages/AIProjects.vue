<template>
  <q-page padding>
    <h2 class="q-mb-xl">AI/ML Projects</h2>
    <div class="text-body1 q-mb-lg">
      With AI/ML technologies exploding all around us, this is where we explore models, deployment frameworks, and experimentation. Current focus includes MCP, A2A, and XAgentX stylings.
    </div>

    <!-- DASMLAB Projects (full-width row) -->
    <div class="row q-col-gutter-lg q-mb-xl">
      <div class="col-12">
        <h4 class="q-mb-md">DASMLAB AI/ML Projects</h4>
        <ProjectCard
          v-for="(project, idx) in dasmlabProjects"
          :key="'dasmlab-' + idx"
          v-bind="project"
        />
      </div>
    </div>

    <!-- Cloud Provider Columns -->
    <div class="row q-col-gutter-lg">
      <!-- AWS AI/ML -->
      <div class="col-12 col-md-4">
        <h4 class="q-mb-md">AWS AI/ML</h4>
        <ProjectCard
          v-for="(project, idx) in awsProjects"
          :key="'aws-' + idx"
          v-bind="project"
        />
      </div>
      <!-- Azure AI/ML -->
      <div class="col-12 col-md-4">
        <h4 class="q-mb-md">Azure AI/ML</h4>
        <ProjectCard
          v-for="(project, idx) in azureProjects"
          :key="'azure-' + idx"
          v-bind="project"
        />
      </div>
      <!-- GCP AI/ML -->
      <div class="col-12 col-md-4">
        <h4 class="q-mb-md">GCP AI/ML</h4>
        <ProjectCard
          v-for="(project, idx) in gcpProjects"
          :key="'gcp-' + idx"
          v-bind="project"
        />
      </div>
    </div>
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
</script>

