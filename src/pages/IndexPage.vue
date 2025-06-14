<template>
  <q-page padding class="q-gutter-md">

    <!-- What's New? -->
    <div class="q-pa-sm flex flex-center">
      <WhatsNew :news="whatsNew" />
    </div>

    <!-- Carousel Section -->
    <div class="q-pa-md">
      <q-carousel
        animated
        height="300px"
        infinite
        swipeable
        transition-prev="slide-right"
        transition-next="slide-left"
        navigation
        control-color="primary"
      >
        <q-carousel-slide name="1">
          <img
            src="https://upload.wikimedia.org/wikipedia/commons/3/3f/Fronalpstock_big.jpg"
            alt="Slide 1"
            class="full-width"
          />
        </q-carousel-slide>
        <q-carousel-slide name="2">
          <img
            src="https://upload.wikimedia.org/wikipedia/commons/thumb/e/ec/Clouds_over_the_Atlantic_Ocean.jpg/1280px-Clouds_over_the_Atlantic_Ocean.jpg"
            alt="Slide 2"
            class="full-width"
          />
        </q-carousel-slide>
      </q-carousel>
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

const router = useRouter()
const pauseScroll = ref(false)

const goTo = (route) => router.push(route)

const cards = [
  { title: 'Frontend Projects', icon: 'view_quilt', route: '/frontend' },
  { title: 'Backend Projects', icon: 'dns', route: '/backend' },
  { title: 'AI/ML Tech Dives', icon: 'psychology', route: '/ai' },
  { title: 'Cloud Provider Techs.', icon: 'cloud', route: '/cloud' },
  { title: 'Infrastructure Projects', icon: 'storage', route: '/infra' },
  { title: 'Security Projects', icon: 'shield', route: '/security' }
]

const whatsNew = ref([
  { id: 1, title: 'MCP Explorer 0.0.7 released', date: '2025-06-13' },
  { id: 2, title: 'New CI/CD Boilerplate published', date: '2025-06-11' },
  { id: 3, title: 'Quasar animation system added', date: '2025-06-10' },
  { id: 4, title: 'Cheapcloud has been scaffolded and basic MSA layout', date: '2025-06-10' }
])


onMounted(() => {
  const container = document.querySelector('.scroll-container')
  let scrollPos = 0

  setInterval(() => {
    if (pauseScroll.value || !container) return
    scrollPos += 1
    if (scrollPos >= container.scrollHeight - container.clientHeight) {
      scrollPos = 0
    }
    container.scrollTo({ top: scrollPos, behavior: 'smooth' })
  }, 75)
})

</script>

<style scoped>
.scroll-container {
  height: 100px;
  overflow-y: hidden;
  position: relative;
}
</style>
