<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)">
    <q-card class="share-sheet" style="min-width: min(440px, 94vw)">
      <q-card-section class="row items-start q-pb-none">
        <div class="col">
          <div class="text-h6">Share</div>
          <div class="text-caption text-grey-7">
            Social second · private first — friends get a DASMLAB tollgate link, not your source dump to Meta.
          </div>
        </div>
        <q-btn flat round dense icon="close" v-close-popup />
      </q-card-section>

      <q-card-section v-if="loading" class="text-center q-py-lg">
        <q-spinner color="primary" size="32px" />
        <div class="text-caption q-mt-sm">Minting tollgate…</div>
      </q-card-section>

      <q-card-section v-else-if="error" class="text-negative">
        {{ error }}
      </q-card-section>

      <template v-else-if="share">
        <q-card-section class="q-pt-sm">
          <div class="share-sheet__title">{{ share.title }}</div>
          <p class="share-sheet__text">{{ share.text }}</p>
          <div v-if="share.private" class="share-sheet__private">
            <q-icon name="lock" size="14px" class="q-mr-xs" />
            Private / unlock-gated (premium foreshadow)
          </div>
        </q-card-section>

        <q-card-section class="q-gutter-sm q-pt-none">
          <q-btn
            unelevated
            color="primary"
            icon="content_copy"
            label="Copy share link"
            class="full-width"
            @click="copy(share.tollgate_url, 'Tollgate link copied — hits metered for FinOps')"
          />
          <q-btn
            outline
            color="primary"
            icon="cloud"
            label="Copy CDN link"
            class="full-width"
            :disable="!share.cdn_url"
            @click="copy(share.cdn_url, 'CDN link copied (direct edge)')"
          />
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="share-sheet__section">Hottest right now</div>
          <div class="share-sheet__grid">
            <q-btn
              v-for="ch in hotChannels"
              :key="ch.id"
              outline
              dense
              no-caps
              class="share-sheet__chip"
              :icon="channelIcon(ch.id)"
              :label="ch.label"
              @click="runChannel(ch)"
            />
          </div>
          <div class="share-sheet__section q-mt-md">Also</div>
          <div class="share-sheet__grid">
            <q-btn
              v-for="ch in otherChannels"
              :key="ch.id"
              flat
              dense
              no-caps
              class="share-sheet__chip"
              :icon="channelIcon(ch.id)"
              :label="ch.label"
              @click="runChannel(ch)"
            >
              <q-tooltip v-if="ch.hint">{{ ch.hint }}</q-tooltip>
            </q-btn>
          </div>
        </q-card-section>
      </template>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useQuasar } from 'quasar'
import { createShare } from 'src/services/surfingApi'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  dayId: { type: String, default: '' },
  mediaId: { type: String, default: '' },
  albumPage: { type: String, default: '' }
})

defineEmits(['update:modelValue'])

const $q = useQuasar()
const loading = ref(false)
const error = ref('')
const share = ref(null)

const hotChannels = computed(() =>
  (share.value?.channels || []).filter((c) => c.hot && c.id !== 'copy_tollgate' && c.id !== 'copy_cdn')
)
const otherChannels = computed(() =>
  (share.value?.channels || []).filter((c) => !c.hot && !['copy_tollgate', 'copy_cdn'].includes(c.id))
)

watch(
  () => props.modelValue,
  async (open) => {
    if (!open) return
    error.value = ''
    share.value = null
    if (!props.dayId) {
      error.value = 'No album selected'
      return
    }
    loading.value = true
    try {
      share.value = await createShare({
        day_id: props.dayId,
        media_id: props.mediaId || undefined,
        album_page: props.albumPage || undefined
      })
    } catch (err) {
      error.value = err?.response?.data?.error || err?.message || 'Could not create share'
    } finally {
      loading.value = false
    }
  }
)

function channelIcon(id) {
  const map = {
    whatsapp: 'chat',
    facebook: 'public',
    threads: 'alternate_email',
    x: 'tag',
    linkedin: 'work',
    instagram: 'camera_alt',
    tiktok: 'music_note',
    native: 'ios_share'
  }
  return map[id] || 'share'
}

async function copy(text, okMessage) {
  if (!text) {
    $q.notify({ type: 'warning', message: 'Nothing to copy' })
    return
  }
  try {
    await navigator.clipboard.writeText(text)
    $q.notify({ type: 'positive', message: okMessage || 'Copied' })
  } catch {
    $q.notify({ type: 'negative', message: 'Could not copy' })
  }
}

async function runChannel(ch) {
  if (!share.value) return
  if (ch.id === 'instagram' || ch.id === 'tiktok' || ch.kind === 'copy') {
    const payload = `${share.value.text}\n${share.value.tollgate_url}`
    await copy(payload, ch.hint || 'Link copied — paste into the app')
    return
  }
  if (ch.kind === 'native') {
    if (navigator.share) {
      try {
        await navigator.share({
          title: share.value.title,
          text: share.value.text,
          url: share.value.tollgate_url
        })
      } catch {
        /* user cancelled */
      }
    } else {
      await copy(share.value.tollgate_url, 'Device share unavailable — link copied instead')
    }
    return
  }
  if (ch.url) {
    window.open(ch.url, '_blank', 'noopener,noreferrer,width=640,height=720')
  }
}
</script>

<style scoped>
.share-sheet__title {
  font-weight: 700;
  color: #102833;
}

.share-sheet__text {
  margin: 0.4rem 0 0;
  font-size: 0.88rem;
  line-height: 1.45;
  color: #5a7080;
}

.share-sheet__private {
  margin-top: 0.55rem;
  font-size: 0.75rem;
  color: #0f8f7c;
  display: flex;
  align-items: center;
}

.share-sheet__section {
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #6a8090;
  margin-bottom: 0.45rem;
}

.share-sheet__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.45rem;
}

.share-sheet__chip {
  justify-content: flex-start;
}
</style>
