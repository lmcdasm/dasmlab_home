import { defineRouter } from '#q-app/wrappers'
import { createRouter, createMemoryHistory, createWebHistory, createWebHashHistory } from 'vue-router'
import routes from './routes'
import { installActivityTracker } from 'src/lib/activityTracker'
import { useAuth } from 'src/composables/useAuth'

export default defineRouter(function (/* { store, ssrContext } */) {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : (process.env.VUE_ROUTER_MODE === 'history' ? createWebHistory : createWebHashHistory)

  const Router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,
    history: createHistory(process.env.VUE_ROUTER_BASE)
  })

  Router.beforeEach(async (to) => {
    if (!to.meta?.admin && !to.meta?.activityViewer) return true
    const { loaded, refresh, oidcEnabled, authenticated, isAdmin, canViewActivity, login } = useAuth()
    if (!loaded.value) await refresh()
    if (!oidcEnabled.value) return true
    if (!authenticated.value) {
      login()
      return false
    }
    if (to.meta.admin && !isAdmin.value) return { path: '/' }
    if (to.meta.activityViewer && !canViewActivity.value) return { path: '/' }
    return true
  })

  installActivityTracker(Router)
  return Router
})
