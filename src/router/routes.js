import AboutPage from 'pages/AboutPage.vue'
import FrontendProjects from 'pages/FrontendProjects.vue'
import BackendProjects from 'pages/BackendProjects.vue'
import AIProjects from 'pages/AIProjects.vue'
import CloudProjects from 'pages/CloudProjects.vue'
import InfraProjects from 'pages/InfraProjects.vue'
import SecurityProjects from 'pages/SecurityProjects.vue'
import ContactPage from 'pages/ContactPage.vue'

const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/IndexPage.vue') },
      { path: 'about', component: AboutPage },
      { path: 'projects/frontend', component: FrontendProjects },
      { path: 'projects/backend', component: BackendProjects },
      { path: 'projects/ai-ml', component: AIProjects },
      { path: 'projects/cloud', component: CloudProjects },
      { path: 'projects/infrastructure', component: InfraProjects },
      { path: 'projects/security', component: SecurityProjects },
      { path: 'contact', component: ContactPage }
    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue')
  }
]

export default routes
