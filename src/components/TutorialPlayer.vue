<template>
  <q-dialog v-model="isOpen" persistent full-width>
    <q-card class="bg-dark text-white">
      <q-card-section>
        <iframe
          :src="resolvedUrl"
          width="100%"
          height="600px"
          frameborder="0"
          title="Tutorial player"
        />
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat label="Close" @click="isOpen = false" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  tutorialUrl: {
    type: String,
    default: ''
  }
})

const isOpen = ref(true)
const resolvedUrl = computed(() => props.tutorialUrl || 'about:blank')
</script>
