<template>
  <q-page padding class="surfing-page q-gutter-md">
    <section class="dasm-shell">
      <div class="dasm-shell__content">
        <div class="dasm-caps">Windsurfing</div>
        <h1 class="dasm-title">Sessions around the world</h1>
        <p class="dasm-subtitle">
          Photo and video log from lakes, rivers, and ocean spots. Organize by day, drag files into the queue, and they land on cluster storage.
        </p>
      </div>
    </section>

    <div v-if="loadError" class="dasm-panel surfing-alert">
      <q-icon name="cloud_off" size="28px" class="q-mb-sm" />
      <div class="text-subtitle2">Surfing service is offline</div>
      <div class="text-caption">{{ loadError }}</div>
    </div>

    <section v-else class="dasm-panel surfing-workspace">
      <div class="surfing-toolbar">
        <q-tabs
          v-model="activeDayId"
          dense
          align="left"
          class="surfing-tabs"
          active-color="primary"
          indicator-color="primary"
          outside-arrows
          mobile-arrows
        >
          <q-tab
            v-for="day in days"
            :key="day.id"
            :name="day.id"
            :label="tabLabel(day)"
            class="surfing-tab"
          />
        </q-tabs>

        <q-btn
          round
          color="primary"
          icon="add"
          class="surfing-add-day-btn"
          @click="openCreateDay"
        >
          <q-tooltip>New day</q-tooltip>
        </q-btn>
      </div>

      <q-tab-panels v-model="activeDayId" animated class="surfing-panels">
        <q-tab-panel
          v-for="day in days"
          :key="day.id"
          :name="day.id"
          class="q-pa-md"
        >
          <div class="day-header">
            <div>
              <h2 class="day-title">{{ day.title }}</h2>
              <div class="day-meta">
                <q-icon name="event" size="16px" class="q-mr-xs" />
                {{ formatDate(day.date) }}
                <span v-if="day.location" class="q-ml-md">
                  <q-icon name="place" size="16px" class="q-mr-xs" />
                  {{ day.location }}
                </span>
              </div>
            </div>
            <q-btn
              flat
              dense
              round
              icon="delete_outline"
              color="negative"
              @click="confirmDeleteDay(day)"
            >
              <q-tooltip>Delete day</q-tooltip>
            </q-btn>
          </div>

          <div
            class="drop-zone"
            :class="{ 'drop-zone--active': dragOverDayId === day.id }"
            @dragenter.prevent="onDragEnter(day.id)"
            @dragover.prevent="onDragOver(day.id)"
            @dragleave.prevent="onDragLeave(day.id, $event)"
            @drop.prevent="onDrop(day.id, $event)"
            @click="openFilePicker(day.id)"
          >
            <q-icon name="add_photo_alternate" size="36px" class="q-mb-sm" />
            <div class="text-subtitle2">Drop photos or videos here</div>
            <div class="text-caption">or click to browse</div>
            <input
              :ref="(el) => setFileInput(day.id, el)"
              type="file"
              multiple
              accept="image/*,video/*"
              class="hidden-input"
              @change="onFileInput(day.id, $event)"
            />
          </div>

          <div v-if="queueForDay(day.id).length" class="upload-queue q-mt-md">
            <div class="queue-title">Upload queue</div>
            <div
              v-for="item in queueForDay(day.id)"
              :key="item.id"
              class="queue-item"
            >
              <q-icon :name="item.file.type.startsWith('video/') ? 'movie' : 'image'" />
              <div class="queue-item__meta">
                <div class="queue-item__name">{{ item.file.name }}</div>
                <q-linear-progress
                  v-if="item.status === 'uploading'"
                  :value="item.progress / 100"
                  color="primary"
                  class="q-mt-xs"
                />
                <div v-else-if="item.status === 'done'" class="text-positive text-caption">Uploaded</div>
                <div v-else-if="item.status === 'error'" class="text-negative text-caption">{{ item.error }}</div>
              </div>
            </div>
          </div>

          <div v-if="!day.media?.length" class="empty-gallery q-mt-lg">
            <q-icon name="surfing" size="42px" class="q-mb-sm" />
            <div>No media yet for this day.</div>
          </div>

          <div v-else class="media-grid q-mt-lg">
            <div
              v-for="item in day.media"
              :key="item.id"
              class="media-card"
            >
              <div class="media-card__frame" @click="openViewer(item)">
                <video
                  v-if="item.media_type === 'video'"
                  :src="mediaUrl(item.url)"
                  muted
                  playsinline
                  preload="metadata"
                  class="media-card__asset"
                />
                <img
                  v-else
                  :src="mediaUrl(item.url)"
                  :alt="item.caption || item.filename"
                  class="media-card__asset"
                  loading="lazy"
                />
                <div class="media-card__badge">
                  <q-icon :name="item.media_type === 'video' ? 'play_circle' : 'photo'" size="18px" />
                </div>
              </div>
              <div class="media-card__footer">
                <div class="media-card__caption">{{ item.caption || item.filename }}</div>
                <div class="media-card__actions">
                  <q-btn
                    flat
                    dense
                    round
                    icon="download"
                    size="sm"
                    tag="a"
                    :href="mediaDownloadUrl(item)"
                    :download="item.filename"
                    rel="noopener"
                    @click.stop
                  >
                    <q-tooltip>Download</q-tooltip>
                  </q-btn>
                  <q-btn
                    flat
                    dense
                    round
                    icon="close"
                    size="sm"
                    @click.stop="removeMedia(day, item)"
                  >
                    <q-tooltip>Remove</q-tooltip>
                  </q-btn>
                </div>
              </div>
            </div>
          </div>
        </q-tab-panel>

        <q-tab-panel v-if="!days.length" name="" class="q-pa-lg text-center">
          <q-icon name="sailing" size="48px" class="q-mb-md" />
          <div class="text-h6 q-mb-sm">No sessions yet</div>
          <div class="text-body2 q-mb-md">Create your first day, then drop photos and clips into it.</div>
          <q-btn color="primary" icon="add" label="Create first day" @click="openCreateDay" />
        </q-tab-panel>
      </q-tab-panels>
    </section>

    <q-dialog v-model="createDayOpen">
      <q-card style="min-width: 320px">
        <q-card-section>
          <div class="text-h6">New session day</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-input v-model="newDay.title" label="Title" filled dense autofocus />
          <q-input v-model="newDay.date" label="Date" type="date" filled dense />
          <q-input v-model="newDay.location" label="Location (optional)" filled dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn color="primary" label="Create" :loading="creatingDay" @click="submitCreateDay" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="viewerOpen" maximized>
      <q-card class="viewer-card">
        <q-bar class="viewer-bar">
          <div>{{ viewerItem?.caption || viewerItem?.filename }}</div>
          <q-space />
          <q-btn dense flat icon="close" v-close-popup />
        </q-bar>
        <q-card-section class="viewer-body">
          <video
            v-if="viewerItem?.media_type === 'video'"
            :src="mediaUrl(viewerItem.url)"
            controls
            autoplay
            class="viewer-asset"
          />
          <img
            v-else-if="viewerItem"
            :src="mediaUrl(viewerItem.url)"
            :alt="viewerItem.caption || viewerItem.filename"
            class="viewer-asset"
          />
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useQuasar } from 'quasar'
import {
  createDay,
  deleteDay,
  deleteMedia,
  fetchDays,
  mediaDownloadUrl,
  mediaUrl,
  uploadMedia
} from 'src/services/surfingApi'

const $q = useQuasar()

const days = ref([])
const activeDayId = ref('')
const loadError = ref('')
const createDayOpen = ref(false)
const creatingDay = ref(false)
const viewerOpen = ref(false)
const viewerItem = ref(null)
const dragOverDayId = ref('')
const fileInputs = ref({})
const uploadQueue = ref([])

const newDay = reactive({
  title: '',
  date: new Date().toISOString().slice(0, 10),
  location: ''
})

function tabLabel(day) {
  return day.title || 'Session'
}

function formatDate(value, short = false) {
  if (!value) return ''
  const parsed = new Date(`${value}T12:00:00`)
  if (Number.isNaN(parsed.getTime())) return value
  return parsed.toLocaleDateString(undefined, short
    ? { month: 'short', day: 'numeric' }
    : { year: 'numeric', month: 'long', day: 'numeric' })
}

function setFileInput(dayId, el) {
  if (el) fileInputs.value[dayId] = el
}

function openFilePicker(dayId) {
  fileInputs.value[dayId]?.click()
}

function queueForDay(dayId) {
  return uploadQueue.value.filter((item) => item.dayId === dayId)
}

async function loadDays(selectId) {
  try {
    const data = await fetchDays()
    days.value = data
    loadError.value = ''
    if (!data.length) {
      activeDayId.value = ''
      return
    }
    if (selectId && data.some((day) => day.id === selectId)) {
      activeDayId.value = selectId
    } else if (!activeDayId.value || !data.some((day) => day.id === activeDayId.value)) {
      activeDayId.value = data[0].id
    }
  } catch (err) {
    loadError.value = err?.message || 'Could not reach surfing service'
  }
}

function openCreateDay() {
  newDay.title = ''
  newDay.date = new Date().toISOString().slice(0, 10)
  newDay.location = ''
  createDayOpen.value = true
}

async function submitCreateDay() {
  if (!newDay.title.trim()) {
    $q.notify({ type: 'warning', message: 'Title is required' })
    return
  }
  creatingDay.value = true
  try {
    const day = await createDay({
      title: newDay.title.trim(),
      date: newDay.date,
      location: newDay.location.trim()
    })
    createDayOpen.value = false
    await loadDays(day.id)
    $q.notify({ type: 'positive', message: 'Day created' })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not create day' })
  } finally {
    creatingDay.value = false
  }
}

function confirmDeleteDay(day) {
  $q.dialog({
    title: 'Delete day?',
    message: `Remove "${day.title}" and all of its media?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await deleteDay(day.id)
      await loadDays()
      $q.notify({ type: 'positive', message: 'Day deleted' })
    } catch (err) {
      $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not delete day' })
    }
  })
}

function onDragEnter(dayId) {
  dragOverDayId.value = dayId
}

function onDragOver(dayId) {
  dragOverDayId.value = dayId
}

function onDragLeave(dayId, event) {
  if (event.currentTarget?.contains(event.relatedTarget)) return
  if (dragOverDayId.value === dayId) dragOverDayId.value = ''
}

function onDrop(dayId, event) {
  dragOverDayId.value = ''
  enqueueFiles(dayId, [...(event.dataTransfer?.files || [])])
}

function onFileInput(dayId, event) {
  enqueueFiles(dayId, [...(event.target.files || [])])
  event.target.value = ''
}

function enqueueFiles(dayId, files) {
  const accepted = files.filter((file) => file.type.startsWith('image/') || file.type.startsWith('video/'))
  if (!accepted.length) {
    $q.notify({ type: 'warning', message: 'Only photos and videos are supported' })
    return
  }
  accepted.forEach((file) => {
    const queueItem = {
      id: `${dayId}-${file.name}-${Date.now()}-${Math.random()}`,
      dayId,
      file,
      progress: 0,
      status: 'queued',
      error: ''
    }
    uploadQueue.value.push(queueItem)
    startUpload(queueItem)
  })
}

async function startUpload(queueItem) {
  queueItem.status = 'uploading'
  try {
    await uploadMedia(queueItem.dayId, queueItem.file, '', (progress) => {
      queueItem.progress = progress
    })
    queueItem.status = 'done'
    queueItem.progress = 100
    await loadDays(queueItem.dayId)
    setTimeout(() => {
      uploadQueue.value = uploadQueue.value.filter((item) => item.id !== queueItem.id)
    }, 1800)
  } catch (err) {
    queueItem.status = 'error'
    queueItem.error = err?.response?.data?.error || 'Upload failed'
  }
}

function openViewer(item) {
  viewerItem.value = item
  viewerOpen.value = true
}

async function removeMedia(day, item) {
  $q.dialog({
    title: 'Remove media?',
    message: item.caption || item.filename,
    cancel: true
  }).onOk(async () => {
    try {
      await deleteMedia(day.id, item.id)
      await loadDays(day.id)
    } catch (err) {
      $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not delete media' })
    }
  })
}

onMounted(() => {
  loadDays()
})
</script>

<style scoped>
.surfing-page {
  max-width: 1180px;
  margin: 0 auto;
}

.surfing-alert {
  text-align: center;
  color: #5d7283;
}

.surfing-workspace {
  overflow: hidden;
}

.surfing-toolbar {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-bottom: 1px solid rgba(36, 61, 81, 0.12);
  padding: 0.35rem 0.75rem 0;
}

.surfing-tabs {
  flex: 1;
  min-width: 0;
}

.surfing-add-day-btn {
  flex-shrink: 0;
}

.surfing-panels {
  background: transparent;
}

.day-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1rem;
}

.day-title {
  margin: 0;
  font-size: 1.35rem;
  color: #1f3344;
}

.day-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.25rem;
  color: #607384;
  margin-top: 0.35rem;
}

.drop-zone {
  border: 2px dashed rgba(63, 122, 107, 0.35);
  border-radius: 16px;
  padding: 1.5rem;
  text-align: center;
  color: #4f6879;
  background: linear-gradient(160deg, rgba(63, 122, 107, 0.06), rgba(158, 115, 178, 0.05));
  cursor: pointer;
  transition: border-color 160ms ease, transform 160ms ease, background 160ms ease;
}

.drop-zone:hover,
.drop-zone--active {
  border-color: rgba(63, 122, 107, 0.8);
  background: linear-gradient(160deg, rgba(63, 122, 107, 0.12), rgba(158, 115, 178, 0.08));
  transform: translateY(-1px);
}

.hidden-input {
  display: none;
}

.upload-queue {
  border: 1px solid rgba(36, 61, 81, 0.12);
  border-radius: 12px;
  padding: 0.75rem;
  background: rgba(255, 255, 255, 0.75);
}

.queue-title {
  font-size: 0.78rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6a7f90;
  margin-bottom: 0.5rem;
}

.queue-item {
  display: flex;
  align-items: flex-start;
  gap: 0.65rem;
  padding: 0.45rem 0;
}

.queue-item + .queue-item {
  border-top: 1px solid rgba(36, 61, 81, 0.08);
}

.queue-item__meta {
  flex: 1;
  min-width: 0;
}

.queue-item__name {
  font-size: 0.88rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-gallery {
  text-align: center;
  color: #6d8292;
  padding: 2rem 1rem;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 1rem;
}

.media-card {
  border: 1px solid rgba(36, 61, 81, 0.12);
  border-radius: 14px;
  overflow: hidden;
  background: #fff;
  box-shadow: 0 8px 18px rgba(24, 43, 60, 0.06);
}

.media-card__frame {
  position: relative;
  aspect-ratio: 4 / 3;
  background: #edf3f8;
  cursor: zoom-in;
}

.media-card__asset {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.media-card__badge {
  position: absolute;
  right: 0.55rem;
  bottom: 0.55rem;
  background: rgba(20, 30, 38, 0.62);
  color: #fff;
  border-radius: 999px;
  padding: 0.2rem 0.45rem;
}

.media-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  padding: 0.55rem 0.65rem;
}

.media-card__actions {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.media-card__caption {
  font-size: 0.82rem;
  color: #4d6273;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.viewer-card {
  background: #101820;
  color: #fff;
}

.viewer-bar {
  background: rgba(0, 0, 0, 0.55);
}

.viewer-body {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 48px);
  background: #0d141a;
}

.viewer-asset {
  max-width: 100%;
  max-height: calc(100vh - 80px);
  object-fit: contain;
}

@media (max-width: 760px) {
  .media-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  }
}
</style>
