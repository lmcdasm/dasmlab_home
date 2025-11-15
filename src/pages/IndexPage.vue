<template>
  <q-page class="index-page q-gutter-md">
    <!-- Hero Section -->
    <HeroDasm />

    <!-- What's New? -->
    <div class="flex flex-center q-mt-md">
      <WhatsNew :news="whatsNew" />
    </div>

    <!-- Carousel Section -->
    <div class="flex flex-center q-my-md">
      <DesignCarousel :entries="carouselEntries" />
    </div>

    <!-- Project Cards -->
    <div class="q-gutter-md row items-start justify-center q-mt-lg">
      <q-card
        v-for="(card, index) in cards"
        :key="card.title"
        v-motion
        :initial="{ scale: 0.98, opacity: 0, y: 16 }"
        :enter="{ scale: 1, opacity: 1, y: 0, transition: { delay: index * 100, duration: 400 } }"
        class="glass tilt q-pa-lg cursor-pointer col-xs-12 col-sm-6 col-md-4"
        flat
        @click="goTo(card.route)"
      >
        <q-card-section class="column items-center justify-center">
          <q-icon :name="card.icon" size="48px" class="q-mb-sm text-primary" />
          <div class="text-subtitle1 text-center text-weight-medium">{{ card.title }}</div>
          <div class="text-caption text-grey-6 text-center q-mt-xs">{{ card.description }}</div>
          <!-- Tech Badges -->
          <div class="row q-mt-md q-gutter-xs justify-center">
            <q-badge
              v-for="tech in card.techs"
              :key="tech"
              outline
              :color="tech === 'Go' ? 'accent' : tech === 'Quasar' ? 'secondary' : 'info'"
              class="tech-badge"
            >
              {{ tech }}
            </q-badge>
          </div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import WhatsNew from 'src/components/WhatsNew.vue'
import DesignCarousel from 'src/components/DesignCarousel.vue'
import HeroDasm from 'src/components/HeroDasm.vue'
import axios from 'axios'

const router = useRouter()

const carouselEntries = ref([])
const whatsNew = ref([])

const goTo = (route) => router.push(route)

const cards = [
  {
    title: 'Frontend Projects',
    icon: 'view_quilt',
    route: '/projects/frontend',
    description: 'Vue.js, Quasar, and modern UI frameworks',
    techs: ['Vue.js', 'Quasar', 'TypeScript']
  },
  {
    title: 'Backend Projects',
    icon: 'dns',
    route: '/projects/backend',
    description: 'Go, gRPC, and microservices',
    techs: ['Go', 'gRPC', 'REST']
  },
  {
    title: 'AI/ML Tech Dives',
    icon: 'psychology',
    route: '/projects/ai-ml',
    description: 'GPU-heavy experiments and ML models',
    techs: ['Python', 'CUDA', 'TensorFlow']
  },
  {
    title: 'Cloud Provider Techs.',
    icon: 'cloud',
    route: '/projects/cloud',
    description: 'Kubernetes, Docker, and cloud infrastructure',
    techs: ['K8s', 'Docker', 'AWS']
  },
  {
    title: 'Infrastructure Projects',
    icon: 'storage',
    route: '/projects/infra',
    description: 'DevOps, CI/CD, and infrastructure as code',
    techs: ['Terraform', 'Ansible', 'GitOps']
  },
  {
    title: 'Security Projects',
    icon: 'shield',
    route: '/projects/security',
    description: 'Security best practices and tools',
    techs: ['OAuth', 'TLS', 'Secrets']
  }
]

onMounted(async () => {
  // Fetch WhatsNew
  try {
    const res = await axios.get('https://whatsnew.svc.dasmlab.org/get')
    whatsNew.value = res.data?.latest_commits || []
  } catch (err) {
    console.error('Failed to Fetch from whatsnew-service:', err)
    whatsNew.value = [
      {
        id: 1,
        project: 'Out of Service',
        title: 'Connection to WhatsNew Service is down.',
        date: new Date().toISOString().replace('T', ' ').substring(0, 19)
      }
    ]
  }

  // Fetch Carousel
  try {
    const res = await axios.get('https://design-carousel.svc.dasmlab.org/carousel')
    carouselEntries.value = res.data
  } catch (err) {
    console.error('Failed to Fetch from design-carousel-service:', err)
    carouselEntries.value = [
      {
        id: 1,
        image_url: '/default.png',
        title: 'Service Offline',
        created_at: new Date().toISOString()
      }
    ]
  }
})
</script>

<style scoped>
.index-page {
  padding: 16px;
}

@media (min-width: 600px) {
  .index-page {
    padding: 24px;
  }
}

.tech-badge {
  transition: all 0.22s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.tech-badge:hover {
  transform: scale(1.1);
  box-shadow: 0 0 8px currentColor;
}
</style>
