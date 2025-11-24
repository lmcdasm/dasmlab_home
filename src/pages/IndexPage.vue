<template>
  <q-page padding class="q-gutter-md">

    <!-- What's New? -->
    <div class="q-pa-sm flex flex-center">
      <WhatsNew :news="whatsNew" />
    </div>

    <!-- Carousel Section -->
    <div class="flex flex-center q-my-md">
      <DesignCarousel :entries="carouselEntries" />
    </div>

    <!-- Project Cards -->
    <div class="q-gutter-md row items-start justify-center">
      <q-card
        v-for="card in cards"
        :key="card.title"
        class="q-pa-lg cursor-pointer col-xs-12 col-sm-4 col-md-3 animate__animated animate__fadeInUp"
        flat
        bordered
        @click="goTo(card.route)"
      >
        <q-card-section class="column items-center justify-center">
          <q-icon :name="card.icon" size="48px" class="q-mb-sm text-primary" />
          <div class="text-subtitle1 text-center">{{ card.title }}</div>
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
import axios from 'axios'

const router = useRouter()

const carouselEntries = ref([])
const whatsNew = ref([])

const goTo = (route) => router.push(route)

const cards = [
  { title: 'Frontend Projects', icon: 'view_quilt', route: '/projects/frontend' },
  { title: 'Backend Projects', icon: 'dns', route: '/projects/backend' },
  { title: 'AI/ML Tech Dives', icon: 'psychology', route: '/projects/ai-ml' },
  { title: 'Cloud Provider Techs.', icon: 'cloud', route: '/projects/cloud' },
  { title: 'Infrastructure Projects', icon: 'storage', route: '/projects/infra' },
  { title: 'Security Projects', icon: 'shield', route: '/projects/security' }
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

