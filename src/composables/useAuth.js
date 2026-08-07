import { computed, onMounted, ref } from 'vue'
import { authLoginUrl, authLogoutUrl, fetchAuthConfig, fetchAuthMe } from 'src/services/authApi'

const me = ref(null)
const config = ref({ enabled: false })
const loaded = ref(false)

export function useAuth() {
  const authenticated = computed(() => !!me.value?.authenticated)
  const isAdmin = computed(() => !!me.value?.user?.is_admin)
  const user = computed(() => me.value?.user || null)
  const oidcEnabled = computed(() => !!config.value?.enabled)

  async function refresh() {
    try {
      config.value = await fetchAuthConfig()
      me.value = await fetchAuthMe()
    } catch {
      me.value = { authenticated: false }
    } finally {
      loaded.value = true
    }
  }

  function login() {
    window.location.href = authLoginUrl()
  }

  function logout() {
    window.location.href = authLogoutUrl()
  }

  onMounted(() => {
    if (!loaded.value) refresh()
  })

  return { me, config, loaded, authenticated, isAdmin, user, oidcEnabled, refresh, login, logout }
}
