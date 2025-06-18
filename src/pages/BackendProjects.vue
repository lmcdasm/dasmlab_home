<template>
  <q-page padding>
    <h2>Backend Projects</h2>
    <p>This is where we explore different Backend Service implementations. My preferred stack is GoLang, wrapped in a OCI  (Docker) and served up in a K8s Deployment for a basic MicroService layout. </p>
	<ProjectCard
     	 v-for="(project, index) in projects"
     	 :key="index"
     	 :title="project.title"
     	 :description="project.description"
     	 :url="project.url"
     	 :language="project.language"
     	 :badge="project.badge"
         :liveUrl="project.liveUrl"
    	/>
  </q-page>
</template>


<script setup>
import ProjectCard from 'src/components/ProjectCard.vue'

const projects = [
  {
    title: 'whatsnew-service',
    description: 'This services monitors all GitRepos in a given organization (or selection) and then provides a JSON output via API endpoint that is consumed by the DASMLAB HOME and displayed, giving the illusion of real-time updates',
    url: 'https://github.com/dasmlab/whatsnew-service',
    language: 'GoLang / Docker / K8s / GitActions',
    badge: 'Public',
    liveUrl: 'https://whatsnew.svc.dasmlab.org/get'
  },
  {
    title: 'design-carousel-service',
    description: 'This service provides a "carousel" of Images, providing a "queue" of images that can be read programatically (JSON) and then consumed to provide a live "feed" of images to the end user.  Allows for Adding, Deleting and updating.  This also converts all images to WEPB for optimization in transfers.',
    url: 'https://github.com/dasmlab/tutorial-generator',
    language: 'GoLang / Docker / K8s / GitActions', 
    badge: 'Public',
    liveUrl: 'https://design-carousel.svc.dasmlab.org/carousel' 

  }
]
</script>
