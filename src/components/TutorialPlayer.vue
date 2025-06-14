<template>
  <q-dialog v-model="isOpen" persistent full-width>
    <q-card class="bg-dark text-white">
      <q-card-section>
        <iframe src="http://localhost:5173/player/index.html"
                width="100%" height="600px" frameborder="0"></iframe>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat label="Close" @click="isOpen = false" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref } from 'vue'
const isOpen = ref(true)
</script>
