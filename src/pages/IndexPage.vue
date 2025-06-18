<template>
  <q-page padding class="q-gutter-md">

    <!-- What's New? -->
    <div class="q-pa-sm flex flex-center">
      <WhatsNew :news="whatsNew" />
    </div>

    <!-- Carousel Section -->
    <div class="flex flex-center q-my-md">
      <q-card flat bordered style="max-width: 700px; width: 100%;">
        <!-- Carousel "picture frame" -->
        <div class="flex flex-center" style="padding: 24px;">
          <div style="border-radius: 16px; overflow: hidden; background: #fff; width: 100%; height: 100%;  display: flex; align-items: center; justify-content: center;">
            <q-carousel
              v-model="slide"
              animated
              height="300px"
              :controls="false"
              :navigation="false"
              :swipeable="true"
              :interval="5000"
              transition-prev="fade"
              transition-next="fade"
            >
              <q-carousel-slide name="1">
                <img
                  src="https://upload.wikimedia.org/wikipedia/commons/3/3f/Fronalpstock_big.jpg"
                  alt="Slide 1"
                  style="width: 100%; height: 100%; object-fit: cover;"
                />
              </q-carousel-slide>
              <q-carousel-slide name="2">
                <img
                  src="https://upload.wikimedia.org/wikipedia/commons/thumb/e/ec/Clouds_over_the_Atlantic_Ocean.jpg/1280px-Clouds_over_the_Atlantic_Ocean.jpg"
                  alt="Slide 2"
                  style="width: 100%; height: 100%; object-fit: cover;"
                />
              </q-carousel-slide>
            </q-carousel>
          </div>
        </div>
        <!-- Controls/dots below the image -->
        <div class="flex flex-center q-pb-md q-pt-sm">
          <q-carousel-control
            :total="2"
            :current="slide"
            @update:current="slide = $event"
            color="accent"
          />
        </div>
      </q-card>
    </div>

    <!-- Custom dot controls below carousel -->
    <div class="flex flex-center q-pb-md q-pt-sm">
      <q-btn
        v-for="n in 2"
        :key="n"
        round
        size="sm"
        :color="slide === n ? 'accent' : 'grey-5'"
        @click="slide = n"
        class="q-mx-xs"
      >
        <span style="width: 8px; height: 8px; display: block; border-radius: 50%; background: currentColor;"></span>
      </q-btn>
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
import axios from 'axios'

const router = useRouter()
const pauseScroll = ref(false)

// Used to pass hold the Slide Carousel (center screen)
const slide = ref('1')

// Used for our WhatsNew Service inputs (API call onMounted())
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
  try {
     const res = await axios.get(
	'https://whatsnew.svc.dasmlab.org/get'
     )
     // comes as JSON { latest_commits: [{elm-1},{elm-2}, ...]}
     whatsNew.value = res.data?.latest_commits || []
  } catch (err) {
     console.error('Failed to Fetch from whatsnew-service:', err)
     // Fall back to OutOfService Messagse on Err
     whatsNew.value = [
      {
        id: 1,
        project: 'Out of Service',
        title: 'Connection to WhatsNew Service is down.',
        date: new Date().toISOString().replace('T', ' ').substring(0, 19) // e.g. 2025-06-15 20:29:15
      }
     ]
  }

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
