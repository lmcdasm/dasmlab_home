import AboutPage from 'pages/AboutPage.vue'
import FrontendProjects from 'pages/FrontendProjects.vue'
import BackendProjects from 'pages/BackendProjects.vue'
import AIProjects from 'pages/AIProjects.vue'
import CloudProjects from 'pages/CloudProjects.vue'
import InfraProjects from 'pages/InfraProjects.vue'
import SecurityProjects from 'pages/SecurityProjects.vue'
import ContactPage from 'pages/ContactPage.vue'
import SurfingPage from 'pages/SurfingPage.vue'
import ActivityPage from 'pages/ActivityPage.vue'

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
      { path: 'contact', component: ContactPage },
      { path: 'surfing', component: SurfingPage },
      {
        path: 'activity',
        name: 'activity',
        component: ActivityPage,
        meta: { admin: true, activityViewer: true }
      }
    ]
  },

  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue')
  }
]

export default routes
