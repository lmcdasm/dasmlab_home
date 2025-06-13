<template>
  <q-page padding>
    <div class="row q-col-gutter-md justify-center q-mt-lg">

      <q-card
        v-for="card in projectCards"
        :key="card.label"
        class="col-xs-12 col-sm-5 col-md-3 q-pa-md cursor-pointer hoverable-card"
        @click="$router.push(card.route)"
        bordered
        flat
      >
        <q-card-section class="flex flex-center column text-center">
          <q-icon :name="card.icon" size="48px" color="primary" />
          <div class="text-h6 q-mt-md">{{ card.label }}</div>
        </q-card-section>
      </q-card>

    </div>
  </q-page>
</template>

<script setup>
const projectCards = [
  { label: 'Frontend Projects', route: '/projects/frontend', icon: 'web' },
  { label: 'Backend Projects', route: '/projects/backend', icon: 'storage' },
  { label: 'AI/ML Tech Dives', route: '/projects/ai-ml', icon: 'psychology' },
  { label: 'Cloud Provider Techs.', route: '/projects/cloud', icon: 'cloud' },
  { label: 'Infrastructure Projects', route: '/projects/infrastructure', icon: 'dns' },
  { label: 'Security Projects', route: '/projects/security', icon: 'security' }
]
</script>

<style scoped>
.hoverable-card {
  transition: box-shadow 0.3s ease;
}
.hoverable-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}
</style>

