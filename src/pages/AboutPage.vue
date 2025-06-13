<template>
  <q-page padding class="q-pa-md q-gutter-lg">
    <div class="column q-col-gutter-lg q-ml-md">

      <!-- Welcome Section -->
      <div class="text-center">
        <q-icon name="emoji_objects" size="64px" color="primary" />
        <div class="text-h4 q-mt-md text-primary">Welcome to DASMLAB!</div>
        <p class="q-mt-md text-body1 text-grey-8">
          This site is a living lab — a place where technology gets poked, prodded, and brought to life. From raw prototypes to polished tools, every project here is part experiment, part education. I believe in demystifying the build process and sharing “how I did it” so others can learn, remix, and iterate.
        </p>
      </div>

      <!-- Bio Section -->
      <q-separator />
      <div class="row items-center q-col-gutter-xl">
        <div class="col-12 col-md-6">
          <img
            :src="meHead"
            class="rounded-borders full-width"
            style="max-width: 400px; width: 100%; height: auto;"
            alt="Me"
          />
        </div>
        <div class="col-12 col-md-6">
          <div class="text-h6 text-primary">Who am I?</div>
          <p class="text-body1 text-grey-9 q-mt-sm">
            I’m someone who's spent over two decades navigating the shifting tides of technology — from cloud automation to high-performance AI platforms. My journey has taken me through telecom, automotive, AI, and beyond, working across continents and learning from over 60 countries and their cultures. I grew up near the Outaouais/Ottawa region, spent adventurous years in Fernie, and now call the Laurentians home.
          </p>
        </div>
      </div>

	<!-- Lifestyle Carousel Section -->
	<q-separator />
	<div class="row items-center q-col-gutter-xl reverse">
	  <div class="col-12 col-md-6">
	     <q-carousel
		  v-model="carouselIndex"
		  animated
		  swipeable
		  infinite
		  autoplay
		  transition-prev="slide-right"
		  transition-next="slide-left"
		  height="320px"
		  class="rounded-borders shadow-2"
	        >
		<q-carousel-slide 
		  :name="index" 
		  v-for="(img, index) in lifestyleImages" 
	          :key="index"
		  class="column items-center justify-center bg-grey-2"
                 >
		<div class="q-pa-md column items-center">
		    <img 
			:src="img.src" 
			:alt="img.caption" 
			style="max-width: 100%; max-height: 220px; object-fit: contain;" 
			class="rounded-borders"
                     />

		    <div class="text-subtitle2 text-white text-center q-mt-sm bg-black bg-opacity-50 q-pa-sm" 
			style="background: rgba(0, 0, 0, 0.7); padding: 6px 12px; width: 100%;">
		      {{ img.caption }}
		    </div>
		</div>
		</q-carousel-slide>
	     </q-carousel>
	  </div>

	  <div class="col-12 col-md-6">
	    <div class="text-h6 text-primary">Beyond the Code</div>
	    <p class="text-body1 text-grey-9 q-mt-sm">
	      Life’s not all tech. When not building or breaking things, you’ll find me snowboarding in winter, windsurfing in summer, digging in the garden, vibing to house music, or shaping clay into something passable in pottery class. These passions fuel the curiosity and balance behind everything I do.
	    </p>
	  </div>
	</div>

    </div>
  </q-page>
</template>

<script setup>

import { ref } from 'vue'
import meHead from 'src/assets/dasm_head.png'
import meRide from 'src/assets/dasm_ride.jpg'
import meRide2 from 'src/assets/me_ride_2.jpg'
import meHome from 'src/assets/me_home.jpg'
import mePlant from 'src/assets/me_plant_1.jpg'
import meLeah from 'src/assets/me_leah_1.jpg'
import meSail from 'src/assets/me_sail_1.jpg'
import meBake1 from 'src/assets/me_baking_1.jpg'
import meBake2 from 'src/assets/me_baking_2.jpg'

const carouselIndex = ref(0)  // start at slide 0


const lifestyleImages = [
  {
    src: meRide,
    caption: 'Windsurfing on a perfect (but low wind!) summer day'
  },
  {
    src: meRide2,
    caption: 'Snowboarding at Mont Ste. Anne'
  },
  {
    src: meHome,
    caption: 'Country Living is the best!'
  },
  {
    src: mePlant,
    caption: 'The magic that is garden therapy! :)'
  },
  {
    src: meLeah,
    caption: 'The Pricess Leah herself!!'
  },
  {
    src: meSail,
    caption: 'My current favourite sail!'
  },
  {
    src: meBake1,
    caption: 'Love to bake!!'
  },
  {
    src: meBake2,
    caption: '... and see the results!!'
  }
]
</script>

<style scoped>
::-webkit-scrollbar {
  display: none;
}

.text-body1 {
  font-family: 'Helvetica Neue', sans-serif;
  font-size: 1.1rem;
}
</style>

