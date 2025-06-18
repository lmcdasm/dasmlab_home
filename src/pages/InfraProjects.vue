<template>
  <q-page padding>
    <h2 class="q-mb-xl">Infrastructure Projects</h2>
    <!-- HERO: Image left, Text right -->
    <div class="row items-center q-mb-xl">
      <div class="col-12 col-md-6 flex flex-center">
        <img
          src="src/assets/cx-stuff/dasmlab_cdevelop_foundation.png"
          alt="Foundational Principles"
          class="infra-img"
        />
      </div>
      <div class="col-12 col-md-6">
        <div class="infra-text">
          <b>DASMLAB.org</b> projects (including this site) all run inside the infrastructure depicted below.
          <b>SecDevOps</b> is the guiding principle: every project is built, tested, and published through a live <b>GitHub Actions</b> pipeline (easily adaptable for GitLab-CI too).
          All releases are <b>semantically versioned</b> and published via manifest to a control SCM, which is continuously deployed to the cluster using <b>FluxCD</b> (or <b>ArgoCD</b>).
          <br><br>
          The result: you get a fast, reproducible, fully automated CI/CD workflow you can inspect, test, and extend for your own infrastructure.
        </div>
      </div>
    </div>

    <!-- INFRA: Text left, Image right (staggered) -->
    <div class="row items-center q-mb-xl">
      <div class="col-12 col-md-6 order-md-2 flex flex-center">
        <img
          src="src/assets/cx-stuff/dasmlab_cdevelop_pipeline_overview.png"
          alt="Infrastructure Overview"
          class="infra-img"
        />
      </div>
      <div class="col-12 col-md-6 order-md-1">
        <div class="infra-text">
          Everything runs on a <b>K3s Kubernetes cluster</b> (with Calico networking and MetalLB in BGP mode for load balancing), protected by off-cluster <b>HAProxy</b>, firewall, and router.
          Monitoring is integrated via <b>Grafana Cloud (Free Tier)</b> for real-time cluster/pod insights.
          <br><br>
          A self-hosted GitHub Runner handles all CI/CD jobs. Both the runner and the prod K8s cluster are VMs on a single server, joined by a virtual switch on ESXi—an ideal “sandbox” for experiments and production alike.
        </div>
      </div>
    </div>

    <div class="text-body2 q-mb-lg">
      Here we dive into orchestration, CI/CD, validation frameworks and more. Think Nephio, Cluster API, Terraform, and other infrastructure plumbing tools.
    </div>
    <ProjectCard
      v-for="(project, index) in projects"
      :key="index"
      :title="project.title"
      :description="project.description"
      :url="project.url"
      :language="project.language"
      :badge="project.badge"
    />
  </q-page>
</template>

<script setup>
import ProjectCard from 'src/components/ProjectCard.vue'

const projects = [
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
.infra-img {
  max-width: 380px;
  width: 100%;
  border-radius: 18px;
  box-shadow: 0 4px 24px #0002;
  margin: 2rem 0;
}
.infra-text {
  font-size: 1.08rem;
  max-width: 520px;
  margin: 0 auto;
  line-height: 1.6;
  padding: 1rem 0;
}
@media (max-width: 1024px) {
  .infra-img { max-width: 90vw; }
  .infra-text { max-width: 98vw; padding: 0.5rem 0;}
}
</style>

