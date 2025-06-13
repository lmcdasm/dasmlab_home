<template>
  <q-layout view="lHh Lpr lFf">
    <!-- Top Toolbar -->
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn flat dense round icon="menu" @click="leftDrawerOpen = !leftDrawerOpen" aria-label="Menu" />
        <q-toolbar-title>
          DASMLAB Home Page
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <!-- Sidebar Navigation -->
    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label header class="text-bold text-primary">DASMLAB.org</q-item-label>

        <q-item clickable to="/about" v-ripple>
          <q-item-section avatar><q-icon name="school" /></q-item-section>
          <q-item-section>About DASMLAB.org</q-item-section>
        </q-item>

        <q-item clickable tag="a" href="https://github.com/lmcdasm" target="_blank" v-ripple>
          <q-item-section avatar><q-icon name="code" /></q-item-section>
          <q-item-section>GitHub</q-item-section>
        </q-item>

        <q-separator spaced />

        <q-item-label header class="text-bold text-primary">Projects</q-item-label>

	<q-item clickable to="/projects/frontend" v-ripple>
  		<q-item-section avatar><q-icon name="web" /></q-item-section>
  		<q-item-section>Frontend Projects</q-item-section>
	</q-item>

	<q-item clickable to="/projects/backend" v-ripple>
  		<q-item-section avatar><q-icon name="storage" /></q-item-section>
  		<q-item-section>Backend Projects</q-item-section>
	</q-item>

	<q-item clickable to="/projects/ai-ml" v-ripple>
  		<q-item-section avatar><q-icon name="psychology" /></q-item-section>
  		<q-item-section>AI/ML Projects</q-item-section>
	</q-item>

	<q-item clickable to="/projects/cloud" v-ripple>
  		<q-item-section avatar><q-icon name="cloud" /></q-item-section>
  		<q-item-section>Cloud Provider Projects</q-item-section>
	</q-item>

	<q-item clickable to="/projects/infrastructure" v-ripple>
  		<q-item-section avatar><q-icon name="dns" /></q-item-section>
  		<q-item-section>Infrastructure Projects</q-item-section>
	</q-item>

	<q-item clickable to="/projects/security" v-ripple>
  		<q-item-section avatar><q-icon name="security" /></q-item-section>
  		<q-item-section>Security Projects</q-item-section>
	</q-item>

        <q-separator spaced />

        <q-item clickable to="/contact" v-ripple>
          <q-item-section avatar><q-icon name="mail" /></q-item-section>
          <q-item-section>Contact</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <!-- Footer -->
    <q-footer class="bg-primary text-white text-center q-pa-sm">
           © 2025 DASMLAB Inc. All rights reserved.
    </q-footer>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue';
const leftDrawerOpen = ref(true);
</script>

