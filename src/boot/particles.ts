import { boot } from 'quasar/wrappers'
import Particles from '@tsparticles/vue3'
import { loadSlim } from '@tsparticles/slim'

export default boot(({ app }) => {
  app.use(Particles, {
    init: async (engine) => {
      // Load the slim version of tsparticles (lightweight, no shapes/animations)
      await loadSlim(engine)
    }
  })
})

