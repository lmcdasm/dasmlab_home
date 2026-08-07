<template>
  <q-page padding class="surfing-page q-gutter-md">
    <section
      class="dasm-shell surfing-hero"
      :class="{ 'surfing-hero--themed': !!activeTheme?.banner_url }"
      :style="heroThemeStyle"
    >
      <div class="surfing-hero__wash" aria-hidden="true" />
      <div class="surfing-hero__glow" aria-hidden="true" />
      <div class="dasm-shell__content surfing-hero__content">
        <div class="dasm-caps surfing-caps">{{ activeTheme?.label || 'Windsurfing Trips' }}</div>
        <h1 class="dasm-title surfing-hero__title">Sessions around the world</h1>
        <p class="dasm-subtitle">
          Your curated box — photos, videos, and activity shares on storage you control.
          Publish to the edge CDN, add notes, share a link. Friends open your gallery;
          Meta never gets the source.
        </p>
        <p v-if="activeTheme?.style_brief" class="surfing-hero__brief">{{ activeTheme.style_brief }}</p>
        <div class="surfing-hero__actions">
          <q-btn
            unelevated
            color="primary"
            icon="ios_share"
            label="Share"
            class="surfing-cta"
            @click="openShareSheet()"
          />
          <q-btn
            outline
            color="primary"
            icon="auto_awesome"
            label="Generate theme"
            class="surfing-cta surfing-cta--ghost"
            :loading="generatingTheme"
            @click="runGenerateTheme"
          >
            <q-tooltip>Sample album photos → AI banner + palette (OpenAI now; cheapcloud farm later)</q-tooltip>
          </q-btn>
          <span class="surfing-hero__hint">Videos · Photos · More — one publishable album</span>
        </div>
      </div>
    </section>

    <a
      class="dasm-panel surfing-related"
      href="https://camera-scrape.apps.2026-prod-1.ocp.dasmlab.org/"
      target="_blank"
      rel="noopener noreferrer"
    >
      <div class="surfing-related__icon" aria-hidden="true">
        <q-icon name="videocam" size="26px" />
      </div>
      <div class="surfing-related__body">
        <div class="surfing-related__title">Live Cams</div>
        <div class="text-caption text-grey-7">
          Ste-Agathe · Jibe City — scrape snapshots and stitch timelapses from the spots you ride.
        </div>
      </div>
      <q-icon name="open_in_new" size="18px" color="grey-6" />
    </a>

    <div v-if="loadError" class="dasm-panel surfing-alert">
      <q-icon name="cloud_off" size="28px" class="q-mb-sm" />
      <div class="text-subtitle2">Surfing service is offline</div>
      <div class="text-caption">{{ loadError }}</div>
    </div>

    <section
      v-else
      class="dasm-panel surfing-workspace"
      :style="workspaceThemeStyle"
    >
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
          class="q-pa-md surfing-day-panel"
        >
          <div class="day-header">
            <div>
              <div class="day-title-row">
                <h2 class="day-title">{{ day.title }}</h2>
                <q-btn
                  v-if="isAdmin"
                  flat
                  dense
                  round
                  size="sm"
                  icon="edit"
                  color="primary"
                  @click="openEditDay(day)"
                >
                  <q-tooltip>Edit title &amp; date display</q-tooltip>
                </q-btn>
              </div>
              <div class="day-meta">
                <template v-if="formatAlbumDate(day)">
                  <q-icon name="event" size="16px" class="q-mr-xs" />
                  {{ formatAlbumDate(day) }}
                </template>
                <span v-if="day.location" :class="{ 'q-ml-md': !!formatAlbumDate(day) }">
                  <q-icon name="place" size="16px" class="q-mr-xs" />
                  {{ day.location }}
                </span>
                <span v-if="day.published" class="day-pill q-ml-md">
                  <q-icon name="cloud_done" size="14px" class="q-mr-xs" />
                  On CDN
                </span>
              </div>
            </div>
            <div class="day-header__actions">
              <q-btn
                v-if="isAdmin"
                unelevated
                dense
                color="primary"
                icon="publish"
                :label="publishButtonLabel(day)"
                :loading="publishingDayId === day.id"
                @click="curateAndPublish(day)"
              >
                <q-tooltip>{{ publishButtonHint(day) }}</q-tooltip>
              </q-btn>
              <q-btn
                v-if="isAdmin"
                outline
                dense
                color="primary"
                icon="auto_awesome"
                label="AI Curate"
                @click="runAICurate(day)"
              />
              <q-btn flat dense round icon="ios_share" color="primary" @click="openShareSheet(day)">
                <q-tooltip>Share album</q-tooltip>
              </q-btn>
              <q-btn flat dense round icon="link" color="primary" @click="copyShareLink(day)">
                <q-tooltip>Copy album page link</q-tooltip>
              </q-btn>
              <q-btn flat dense round icon="add_link" color="primary" @click="openAddLink(day)">
                <q-tooltip>Add activity / share link</q-tooltip>
              </q-btn>
              <q-btn flat dense round icon="delete_outline" color="negative" @click="confirmDeleteDay(day)">
                <q-tooltip>Delete day</q-tooltip>
              </q-btn>
            </div>
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
            <q-icon name="waves" size="40px" class="q-mb-sm drop-zone__icon" />
            <div class="text-subtitle2">Drop photos or videos into this session</div>
            <div class="text-caption">or click to browse · uploads go direct-to-CDN · activity shares via + link</div>
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
            <div class="queue-title row items-center no-wrap">
              <span class="col">
                Upload queue
                <span class="queue-title__stats">
                  · {{ queueStats(day.id).active }} active
                  <template v-if="queueStats(day.id).done"> · {{ queueStats(day.id).done }} done</template>
                  <template v-if="queueStats(day.id).error"> · {{ queueStats(day.id).error }} failed</template>
                </span>
              </span>
              <q-btn
                v-if="queueStats(day.id).done || queueStats(day.id).error"
                flat
                dense
                size="sm"
                color="primary"
                label="Clear finished"
                @click="clearFinishedUploads(day.id)"
              />
            </div>
            <div class="upload-queue__scroll">
              <div
                v-for="item in queueForDay(day.id)"
                :key="item.id"
                class="queue-item"
                :class="`queue-item--${item.status}`"
              >
                <q-icon :name="item.file.type.startsWith('video/') ? 'movie' : 'image'" />
                <div class="queue-item__meta">
                  <div class="queue-item__name">{{ item.file.name }}</div>
                  <q-linear-progress
                    v-if="item.status === 'uploading' || item.status === 'queued'"
                    :value="item.status === 'queued' ? 0 : item.progress / 100"
                    :indeterminate="item.status === 'queued'"
                    color="primary"
                    class="q-mt-xs"
                  />
                  <div v-else-if="item.status === 'done'" class="text-positive text-caption">Uploaded</div>
                  <div v-else-if="item.status === 'error'" class="text-negative text-caption">{{ item.error }}</div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="!day.media?.length" class="empty-gallery q-mt-lg">
            <q-icon name="surfing" size="48px" class="q-mb-sm" />
            <div>No media yet — drop files or add a Garmin / iPhone share.</div>
          </div>

          <template v-else>
            <div
              v-for="section in mediaSections(day)"
              :key="section.key"
              class="media-section"
              :class="`media-section--${section.key}`"
            >
              <div class="media-section__head">
                <div class="media-section__label">
                  <q-icon :name="section.icon" size="20px" />
                  <h3 class="media-section__title">{{ section.title }}</h3>
                  <span class="media-section__count">{{ section.items.length }}</span>
                </div>
                <p class="media-section__blurb">{{ section.blurb }}</p>
              </div>

              <!-- Videos: Discovery-style map + card browser -->
              <div v-if="section.key === 'videos'">
                <div class="video-sort-bar row items-center q-mb-sm">
                  <span class="text-caption text-grey-7 q-mr-sm">Sort</span>
                  <q-btn-toggle
                    v-model="videoSort"
                    dense
                    toggle-color="primary"
                    :options="videoSortOptions"
                  />
                </div>
                <VideoAlbumMap
                  :items="section.items"
                  :album-title="day.title"
                  :day-id="day.id"
                  :sort-mode="videoSort"
                  :publishing="publishingDayId === day.id"
                  :tagging="taggingMedia"
                  @play="(item) => openViewer(item, day)"
                  @share="(item) => openShareSheet(day, item)"
                  @cdn="openOutbound"
                  @publish="() => publishAlbum(day)"
                  @propose-tag="(item, name) => proposeTag(day, item, name)"
                  @approve-tag="(item, tag) => moderateTag(day, item, tag, 'approve')"
                  @reject-tag="(item, tag) => moderateTag(day, item, tag, 'reject')"
                >
                  <template #grid>
                    <div class="video-rail">
                      <article
                        v-for="item in section.items"
                        :key="item.id"
                        class="video-card"
                        :class="{ 'media-card--hidden': item.hidden }"
                      >
                        <button type="button" class="video-card__stage" @click="openViewer(item, day)">
                          <video
                            :src="mediaUrl(item.url)"
                            muted
                            playsinline
                            preload="metadata"
                            class="video-card__asset"
                            @loadedmetadata="onVideoMeta(day, item, $event)"
                          />
                          <span v-if="clipDuration(item)" class="video-card__dur">{{ formatDuration(clipDuration(item)) }}</span>
                          <span class="video-card__play">
                            <q-icon name="play_arrow" size="36px" />
                          </span>
                          <span class="video-card__open">Play muted</span>
                        </button>
                        <div class="video-card__body">
                          <div class="video-card__title">{{ item.caption || item.filename }}</div>
                          <div class="video-card__meta text-caption">
                            <template v-if="clipDuration(item)">{{ formatDuration(clipDuration(item)) }} · </template>
                            {{ item.play_count || 0 }} plays
                            · {{ (item.tags || []).filter((t) => t.status === 'approved').length }} tags
                          </div>
                          <p v-if="item.notes" class="video-card__notes">{{ item.notes }}</p>
                          <div class="video-card__actions">
                            <q-btn flat dense size="sm" icon="ios_share" label="Share" @click="openShareSheet(day, item)" />
                            <q-btn
                              v-if="canDownloadMedia(item)"
                              flat
                              dense
                              size="sm"
                              icon="download"
                              label="Download"
                              tag="a"
                              :href="mediaDownloadUrl(item, day.id)"
                              :download="item.filename"
                              rel="noopener"
                            />
                            <q-btn flat dense size="sm" icon="open_in_new" label="CDN" @click="openOutbound(item)" />
                            <q-btn flat dense size="sm" icon="edit_note" label="Notes" @click="openNotesEditor(day, item)" />
                            <q-btn flat dense size="sm" icon="visibility_off" @click="removeMedia(day, item)" />
                          </div>
                        </div>
                      </article>
                    </div>
                  </template>
                </VideoAlbumMap>
              </div>

              <!-- Photos: Cards / Cabinet (parity with videos) -->
              <div v-else-if="section.key === 'photos'">
                <PhotoAlbumMap
                  :items="section.items"
                  :album-title="day.title"
                  :day-id="day.id"
                  :publishing="publishingDayId === day.id"
                  :tagging="taggingMedia"
                  :is-admin="isAdmin"
                  @view="(item) => openViewer(item, day)"
                  @share="(item) => openShareSheet(day, item)"
                  @cdn="openOutbound"
                  @publish="() => publishAlbum(day)"
                  @notes="(item) => openNotesEditor(day, item)"
                  @hide="(item) => removeMedia(day, item)"
                  @restore="(item) => restoreMedia(day, item)"
                  @propose-tag="(item, name) => proposeTag(day, item, name)"
                  @approve-tag="(item, tag) => moderateTag(day, item, tag, 'approve')"
                  @reject-tag="(item, tag) => moderateTag(day, item, tag, 'reject')"
                >
                  <template #grid>
                    <div class="photo-rail">
                      <article
                        v-for="item in section.items"
                        :key="item.id"
                        class="photo-card"
                        :class="{ 'media-card--hidden': item.hidden }"
                      >
                        <button type="button" class="photo-card__stage" @click="openViewer(item, day)">
                          <img
                            :src="mediaUrl(item.url)"
                            :alt="item.caption || item.filename"
                            class="photo-card__asset"
                            loading="lazy"
                          />
                          <span class="photo-card__glow" aria-hidden="true" />
                          <span class="photo-card__open">
                            <q-icon name="zoom_in" size="28px" />
                          </span>
                          <span v-if="item.published" class="photo-card__badge">CDN</span>
                          <span v-else-if="item.notes" class="photo-card__badge photo-card__badge--notes">Notes</span>
                        </button>
                        <div class="photo-card__body">
                          <div class="photo-card__title">{{ item.caption || item.filename }}</div>
                          <div class="photo-card__meta text-caption">
                            {{ (item.tags || []).filter((t) => t.status === 'approved').length }} tags
                            · {{ item.published ? 'CDN' : 'Draft' }}
                          </div>
                          <p v-if="item.notes" class="photo-card__notes">{{ item.notes }}</p>
                          <div class="photo-card__actions">
                            <q-btn flat dense size="sm" icon="ios_share" label="Share" @click="openShareSheet(day, item)" />
                            <q-btn flat dense size="sm" icon="edit_note" label="Notes" @click="openNotesEditor(day, item)" />
                            <q-btn
                              v-if="canDownloadMedia(item)"
                              flat
                              dense
                              size="sm"
                              icon="download"
                              label="Download"
                              tag="a"
                              :href="mediaDownloadUrl(item, day.id)"
                              :download="item.filename"
                              rel="noopener"
                            />
                            <q-btn flat dense size="sm" icon="open_in_new" label="CDN" @click="openOutbound(item)" />
                            <q-btn flat dense size="sm" icon="visibility_off" @click="removeMedia(day, item)" />
                            <q-btn
                              v-if="isAdmin && item.hidden"
                              flat
                              dense
                              size="sm"
                              icon="visibility"
                              color="primary"
                              label="Show"
                              @click="restoreMedia(day, item)"
                            />
                          </div>
                        </div>
                      </article>
                    </div>
                  </template>
                </PhotoAlbumMap>
              </div>

              <!-- Other: Garmin / iPhone / shares -->
              <div v-else class="other-list">
                <article
                  v-for="item in section.items"
                  :key="item.id"
                  class="other-card"
                >
                  <div class="other-card__icon">
                    <q-icon :name="otherIcon(item)" size="28px" />
                  </div>
                  <div class="other-card__body">
                    <div class="other-card__title">{{ item.caption || item.filename }}</div>
                    <p v-if="item.notes" class="other-card__notes">{{ item.notes }}</p>
                    <a
                      v-if="mediaOutboundUrl(item)"
                      class="other-card__link"
                      :href="mediaOutboundUrl(item)"
                      target="_blank"
                      rel="noopener"
                    >Open share</a>
                  </div>
                  <div class="other-card__actions">
                    <q-btn flat dense round icon="ios_share" @click="openShareSheet(day, item)" />
                    <q-btn
                      v-if="canDownloadMedia(item)"
                      flat
                      dense
                      round
                      icon="download"
                      tag="a"
                      :href="mediaDownloadUrl(item, day.id)"
                      :download="item.filename"
                      rel="noopener"
                    />
                    <q-btn flat dense round icon="edit_note" @click="openNotesEditor(day, item)" />
                    <q-btn flat dense round icon="visibility_off" @click="removeMedia(day, item)" />
                  </div>
                </article>
              </div>
            </div>
          </template>
        </q-tab-panel>

        <q-tab-panel v-if="!days.length" name="" class="q-pa-lg text-center">
          <q-icon name="sailing" size="48px" class="q-mb-md" />
          <div class="text-h6 q-mb-sm">No sessions yet</div>
          <div class="text-body2 q-mb-md">Create your first day, then drop photos, clips, or activity links into it.</div>
          <q-btn color="primary" icon="add" label="Create first day" @click="openCreateDay" />
        </q-tab-panel>
      </q-tab-panels>
    </section>

    <ShareSheet
      v-model="shareOpen"
      :day-id="shareDayId"
      :media-id="shareMediaId"
      :album-page="shareAlbumPage"
    />

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

    <q-dialog v-model="editDayOpen">
      <q-card style="min-width: min(380px, 92vw)">
        <q-card-section>
          <div class="text-h6">Edit album</div>
          <div class="text-caption text-grey-7">Title and how the date appears on the page</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-input v-model="editDay.title" label="Title" filled dense autofocus />
          <q-input v-model="editDay.date" label="Date (source)" type="date" filled dense />
          <q-input v-model="editDay.location" label="Location (optional)" filled dense />
          <q-select
            v-model="editDay.date_precision"
            :options="datePrecisionOptions"
            label="Date display"
            filled
            dense
            emit-value
            map-options
            hint="Hide, year only, year + month, or full date"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn color="primary" label="Save" :loading="savingDay" @click="submitEditDay" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="linkOpen">
      <q-card style="min-width: min(420px, 92vw)">
        <q-card-section>
          <div class="text-h6">Add activity / share</div>
          <div class="text-caption text-grey-7">Garmin Connect, iPhone share sheets, Strava, drive links…</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-input v-model="newLink.title" label="Title" filled dense autofocus />
          <q-input v-model="newLink.url" label="URL" filled dense type="url" hint="https://…" />
          <q-input v-model="newLink.source_label" label="Source label (optional)" filled dense placeholder="Garmin · iPhone · Drive" />
          <q-input v-model="newLink.notes" label="Notes" filled dense type="textarea" autogrow />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn color="primary" label="Add" :loading="savingLink" @click="submitAddLink" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="notesOpen">
      <q-card style="min-width: min(440px, 92vw)">
        <q-card-section>
          <div class="text-h6">Notes</div>
          <div class="text-caption text-grey-7">{{ notesDraft.filename }}</div>
        </q-card-section>
        <q-card-section class="q-gutter-sm">
          <q-input v-model="notesDraft.caption" label="Caption" filled dense />
          <q-input
            v-model="notesDraft.notes"
            label="Notes"
            filled
            type="textarea"
            autogrow
            hint="Session feel, wind, who was there…"
          />
          <q-select
            v-model="notesDraft.kind"
            :options="kindOptions"
            label="Section"
            filled
            dense
            emit-value
            map-options
          />
          <q-select
            v-model="notesDraft.notes_visibility"
            :options="notesVisibilityOptions"
            label="Notes visibility"
            filled
            dense
            emit-value
            map-options
            hint="Public / private (you) / group (signed-up members later)"
          />
          <q-select
            v-model="notesDraft.download_visibility"
            :options="downloadVisibilityOptions"
            label="Download access"
            filled
            dense
            emit-value
            map-options
            hint="Who can use Download — gated by DASMLAB, not a raw CDN guess"
          />
          <q-input
            v-model="notesDraft.external_url"
            label="Outbound link (optional)"
            filled
            dense
            type="url"
            hint="Videos & shares can open out to this URL"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn color="primary" label="Save" :loading="savingNotes" @click="saveNotes" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="viewerOpen" maximized @hide="onViewerHide">
      <q-card class="viewer-card">
        <q-bar class="viewer-bar">
          <div>{{ viewerItem?.caption || viewerItem?.filename }}</div>
          <q-space />
          <q-btn
            v-if="viewerItem?.media_type === 'video'"
            dense
            flat
            :icon="viewerMuted ? 'volume_off' : 'volume_up'"
            @click="toggleViewerMute"
          >
            <q-tooltip>
              {{ viewerMuted ? 'Unmute (starts muted — loud/unmastered clips)' : 'Mute' }}
            </q-tooltip>
          </q-btn>
          <q-btn
            v-if="viewerItem?.media_type === 'video'"
            dense
            flat
            icon="auto_fix"
            @click="audioCleanupTeaser"
          >
            <q-tooltip>Premium foreshadow: AI clean / overlay audio — you keep adapter control</q-tooltip>
          </q-btn>
          <q-btn
            v-if="viewerItem"
            dense
            flat
            icon="edit_note"
            @click="openNotesEditor(activeDay, viewerItem)"
          />
          <q-btn dense flat icon="close" v-close-popup />
        </q-bar>
        <q-card-section class="viewer-body">
          <div v-if="viewerItem?.media_type === 'video'" class="viewer-video-wrap">
            <video
              ref="viewerVideoEl"
              :src="mediaUrl(viewerItem.url)"
              controls
              autoplay
              playsinline
              :muted="viewerMuted"
              class="viewer-asset"
              @volumechange="onViewerVolumeChange"
            />
            <button
              v-if="viewerMuted"
              type="button"
              class="viewer-unmute"
              @click="toggleViewerMute"
            >
              <q-icon name="volume_off" size="22px" />
              Sound off — tap to unmute
            </button>
          </div>
          <img
            v-else-if="viewerItem"
            :src="mediaUrl(viewerItem.url)"
            :alt="viewerItem.caption || viewerItem.filename"
            class="viewer-asset"
          />
        </q-card-section>
        <q-card-section v-if="viewerItem?.notes" class="viewer-notes">
          {{ viewerItem.notes }}
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useQuasar } from 'quasar'
import ShareSheet from 'src/components/ShareSheet.vue'
import PhotoAlbumMap from 'src/components/PhotoAlbumMap.vue'
import VideoAlbumMap from 'src/components/VideoAlbumMap.vue'
import { useAuth } from 'src/composables/useAuth'
import {
  addMediaLink,
  aiCurate,
  createDay,
  deleteDay,
  deleteMedia,
  fetchDays,
  generateTheme,
  mediaDownloadUrl,
  mediaKind,
  mediaOutboundUrl,
  mediaUrl,
  canDownloadMedia,
  moderateMediaTag,
  patchDay,
  proposeMediaTag,
  publishDayUntilDone,
  recordMediaPlay,
  unhideMedia,
  updateMedia,
  uploadMedia
} from 'src/services/surfingApi'

const $q = useQuasar()
const { isAdmin, oidcEnabled, login } = useAuth()

const days = ref([])
const activeDayId = ref('')
const loadError = ref('')
const createDayOpen = ref(false)
const creatingDay = ref(false)
const editDayOpen = ref(false)
const savingDay = ref(false)
const editDayId = ref('')
const editDay = reactive({
  title: '',
  date: '',
  location: '',
  date_precision: 'day'
})

const datePrecisionOptions = [
  { label: 'Hide date', value: 'hide' },
  { label: 'Year only', value: 'year' },
  { label: 'Year and month', value: 'month' },
  { label: 'Full date (year · month · day)', value: 'day' }
]
const viewerOpen = ref(false)
const viewerItem = ref(null)
const viewerMuted = ref(true)
const viewerVideoEl = ref(null)
const dragOverDayId = ref('')
const fileInputs = ref({})
const uploadQueue = ref([])
const generatingTheme = ref(false)
const shareOpen = ref(false)
const shareDayId = ref('')
const shareMediaId = ref('')
const shareAlbumPage = ref('')
const publishingDayId = ref('')
const taggingMedia = ref(false)

const linkOpen = ref(false)
const savingLink = ref(false)
const linkDayId = ref('')
const newLink = reactive({
  title: '',
  url: '',
  notes: '',
  source_label: ''
})

const notesOpen = ref(false)
const savingNotes = ref(false)
const notesDayId = ref('')
const notesMediaId = ref('')
const notesDraft = reactive({
  filename: '',
  caption: '',
  notes: '',
  kind: 'photo',
  external_url: '',
  notes_visibility: 'public',
  download_visibility: 'public'
})

const kindOptions = [
  { label: 'Photos', value: 'photo' },
  { label: 'Videos', value: 'video' },
  { label: 'More (shares / other)', value: 'other' }
]

const notesVisibilityOptions = [
  { label: 'Public', value: 'public' },
  { label: 'Private (only you)', value: 'private' },
  { label: 'Group (signed-up members)', value: 'group' }
]

const downloadVisibilityOptions = [
  { label: 'Public can download', value: 'public' },
  { label: 'Private — owner only', value: 'private' },
  { label: 'Group can download', value: 'group' }
]

const newDay = reactive({
  title: '',
  date: new Date().toISOString().slice(0, 10),
  location: ''
})

const activeDay = computed(() => days.value.find((d) => d.id === activeDayId.value) || null)
const activeTheme = computed(() => activeDay.value?.theme || null)

const heroThemeStyle = computed(() => {
  const t = activeTheme.value
  if (!t) return {}
  const style = {}
  if (t.primary) style['--surf-teal'] = t.primary
  if (t.secondary) style['--surf-deep'] = t.secondary
  if (t.accent) style['--surf-horizon'] = t.accent
  if (t.banner_url) {
    style.backgroundImage = `linear-gradient(125deg, rgba(6, 47, 56, 0.72), rgba(10, 92, 88, 0.55) 48%, rgba(20, 122, 108, 0.45)), url(${mediaUrl(t.banner_url)})`
    style.backgroundSize = 'cover'
    style.backgroundPosition = 'center'
  }
  return style
})

const workspaceThemeStyle = computed(() => {
  const t = activeTheme.value
  if (!t?.background_url && !t?.primary) return {}
  const style = {}
  if (t.background_url) {
    style.backgroundImage = `linear-gradient(180deg, rgba(255,255,255,0.92), rgba(240,250,247,0.88)), url(${mediaUrl(t.background_url)})`
    style.backgroundSize = 'cover'
    style.backgroundPosition = 'center'
  }
  if (t.primary) style.borderColor = `${t.primary}55`
  return style
})

function tabLabel(day) {
  return day.title || 'Session'
}

function mediaSections(day) {
  const items = day?.media || []
  const videos = sortMediaItems(
    items.filter((m) => mediaKind(m) === 'video'),
    videoSort.value
  )
  const photos = items.filter((m) => mediaKind(m) === 'photo')
  const other = items.filter((m) => mediaKind(m) === 'other')
  const sections = []
  if (videos.length) {
    sections.push({
      key: 'videos',
      title: 'Videos',
      icon: 'movie',
      blurb: 'Longest clips first by default — best candidates to chop. Cards / Cabinet for browse.',
      items: videos
    })
  }
  if (photos.length) {
    sections.push({
      key: 'photos',
      title: 'Photos',
      icon: 'photo_camera',
      blurb: 'Session stills — Cabinet by default. Notes, tags, share, and publish live next to the frame.',
      items: photos
    })
  }
  if (other.length) {
    sections.push({
      key: 'other',
      title: 'More',
      icon: 'hub',
      blurb: 'Garmin activities, iPhone shares, and other publishable data.',
      items: other
    })
  }
  return sections
}

/** Longest-first default — length is the chop-list signal. */
const videoSort = ref('longest')
const videoSortOptions = [
  { label: 'Longest', value: 'longest' },
  { label: 'Shortest', value: 'shortest' },
  { label: 'Newest', value: 'newest' },
  { label: 'Name', value: 'name' },
  { label: 'Plays', value: 'plays' }
]

/** Client-side duration cache until persisted on the server. */
const durationCache = reactive({})

function clipDuration(item) {
  if (!item) return 0
  const cached = durationCache[item.id]
  if (cached > 0) return cached
  return Number(item.duration_sec) || 0
}

function formatDuration(sec) {
  const s = Math.round(Number(sec) || 0)
  if (s <= 0) return ''
  const m = Math.floor(s / 60)
  const r = s % 60
  if (m >= 60) {
    const h = Math.floor(m / 60)
    const mm = m % 60
    return `${h}:${String(mm).padStart(2, '0')}:${String(r).padStart(2, '0')}`
  }
  return `${m}:${String(r).padStart(2, '0')}`
}

function sortMediaItems(list, mode) {
  const items = [...(list || [])]
  const dur = (i) => clipDuration(i)
  switch (mode) {
    case 'shortest':
      return items.sort((a, b) => dur(a) - dur(b) || (b.created_at || '').localeCompare(a.created_at || ''))
    case 'newest':
      return items.sort((a, b) => (b.created_at || '').localeCompare(a.created_at || ''))
    case 'name':
      return items.sort((a, b) =>
        String(a.caption || a.filename || '').localeCompare(String(b.caption || b.filename || ''), undefined, {
          sensitivity: 'base'
        })
      )
    case 'plays':
      return items.sort((a, b) => (b.play_count || 0) - (a.play_count || 0) || dur(b) - dur(a))
    case 'longest':
    default:
      // Unknown length sinks; known long clips rise — best chop candidates first.
      return items.sort((a, b) => {
        const da = dur(a)
        const db = dur(b)
        if (da > 0 && db > 0) return db - da
        if (da > 0) return -1
        if (db > 0) return 1
        return (b.created_at || '').localeCompare(a.created_at || '')
      })
  }
}

const durationPersistPending = new Set()

async function onVideoMeta(day, item, ev) {
  const el = ev?.target
  const sec = el?.duration
  if (!item?.id || !Number.isFinite(sec) || sec <= 0) return
  durationCache[item.id] = sec
  // Persist once for owner so next load sorts without re-probing every clip.
  if (!isAdmin.value || item.duration_sec > 0 || durationPersistPending.has(item.id)) return
  durationPersistPending.add(item.id)
  try {
    await updateMedia(day.id, item.id, { duration_sec: Math.round(sec * 10) / 10 })
    item.duration_sec = Math.round(sec * 10) / 10
  } catch {
    durationPersistPending.delete(item.id)
  }
}

function otherIcon(item) {
  const hay = `${item.caption || ''} ${item.filename || ''} ${item.notes || ''}`.toLowerCase()
  if (hay.includes('garmin')) return 'directions_bike'
  if (hay.includes('iphone') || hay.includes('apple')) return 'phone_iphone'
  if (hay.includes('strava')) return 'directions_run'
  return 'link'
}

function openOutbound(item) {
  const url = mediaOutboundUrl(item)
  if (!url) return
  window.open(url, '_blank', 'noopener,noreferrer')
}

async function runGenerateTheme() {
  const day = activeDay.value
  if (!day) {
    $q.notify({ type: 'warning', message: 'Create or select a session day first' })
    return
  }
  generatingTheme.value = true
  try {
    const data = await generateTheme(day.id, {
      sport: 'Windsurfing Trips',
      prompt: `${day.title}${day.location ? ` in ${day.location}` : ''} — ocean sessions, trade winds, travel log`,
      sample_count: 3
    })
    await loadDays(day.id)
    const warning = data?.warning
    if (warning) {
      $q.notify({ type: 'warning', message: warning, timeout: 5000 })
    } else {
      $q.notify({ type: 'positive', message: 'Theme generated from your album samples' })
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err?.response?.data?.error || err?.message || 'Theme generation failed'
    })
  } finally {
    generatingTheme.value = false
  }
}

function openShareSheet(day, item) {
  const target = day || activeDay.value
  if (!target?.id) {
    $q.notify({ type: 'warning', message: 'Select a session day first' })
    return
  }
  shareDayId.value = target.id
  shareMediaId.value = item?.id || ''
  const url = new URL(window.location.href)
  url.hash = `day=${target.id}`
  shareAlbumPage.value = url.toString()
  shareOpen.value = true
}

async function copyShareLink(day) {
  const targetId = day?.id || activeDayId.value
  const url = new URL(window.location.href)
  if (targetId) {
    url.hash = `day=${targetId}`
  }
  try {
    await navigator.clipboard.writeText(url.toString())
    $q.notify({ type: 'positive', message: 'Share link copied — your box, their invite' })
  } catch {
    $q.notify({ type: 'negative', message: 'Could not copy link' })
  }
}

/** Public album date label — respects date_precision (hide|year|month|day). */
function formatAlbumDate(day) {
  if (!day?.date) return ''
  const precision = (day.date_precision || 'day').toLowerCase()
  if (precision === 'hide') return ''
  const parsed = new Date(`${day.date}T12:00:00`)
  if (Number.isNaN(parsed.getTime())) return day.date
  if (precision === 'year') {
    return parsed.toLocaleDateString(undefined, { year: 'numeric' })
  }
  if (precision === 'month') {
    return parsed.toLocaleDateString(undefined, { year: 'numeric', month: 'long' })
  }
  return parsed.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })
}

function openEditDay(day) {
  editDayId.value = day.id
  editDay.title = day.title || ''
  editDay.date = day.date || ''
  editDay.location = day.location || ''
  editDay.date_precision = day.date_precision || 'day'
  editDayOpen.value = true
}

async function submitEditDay() {
  if (!editDay.title.trim()) {
    $q.notify({ type: 'warning', message: 'Title is required' })
    return
  }
  savingDay.value = true
  try {
    const updated = await patchDay(editDayId.value, {
      title: editDay.title.trim(),
      date: editDay.date,
      location: editDay.location.trim(),
      date_precision: editDay.date_precision || 'day'
    })
    const idx = days.value.findIndex((d) => d.id === editDayId.value)
    if (idx >= 0) {
      days.value[idx] = { ...days.value[idx], ...updated, media: days.value[idx].media }
    }
    editDayOpen.value = false
    $q.notify({ type: 'positive', message: 'Album updated' })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || err?.message || 'Could not save' })
  } finally {
    savingDay.value = false
  }
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

function queueStats(dayId) {
  const items = queueForDay(dayId)
  return {
    active: items.filter((i) => i.status === 'queued' || i.status === 'uploading').length,
    done: items.filter((i) => i.status === 'done').length,
    error: items.filter((i) => i.status === 'error').length
  }
}

function clearFinishedUploads(dayId) {
  uploadQueue.value = uploadQueue.value.filter(
    (item) => item.dayId !== dayId || (item.status !== 'done' && item.status !== 'error')
  )
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
    const hashDay = (window.location.hash.match(/day=([^&]+)/) || [])[1]
    const preferred = selectId || hashDay
    if (preferred && data.some((day) => day.id === preferred)) {
      activeDayId.value = preferred
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

function openAddLink(day) {
  linkDayId.value = day.id
  newLink.title = ''
  newLink.url = ''
  newLink.notes = ''
  newLink.source_label = ''
  linkOpen.value = true
}

async function submitAddLink() {
  if (!newLink.url.trim()) {
    $q.notify({ type: 'warning', message: 'URL is required' })
    return
  }
  savingLink.value = true
  try {
    await addMediaLink(linkDayId.value, {
      title: newLink.title.trim(),
      url: newLink.url.trim(),
      notes: newLink.notes.trim(),
      source_label: newLink.source_label.trim()
    })
    linkOpen.value = false
    await loadDays(linkDayId.value)
    $q.notify({ type: 'positive', message: 'Share added under More' })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not add link' })
  } finally {
    savingLink.value = false
  }
}

function openNotesEditor(day, item) {
  if (!day || !item) return
  notesDayId.value = day.id
  notesMediaId.value = item.id
  notesDraft.filename = item.filename || ''
  notesDraft.caption = item.caption || ''
  notesDraft.notes = item.notes || ''
  notesDraft.kind = mediaKind(item)
  notesDraft.external_url = item.external_url || ''
  notesDraft.notes_visibility = item.notes_visibility || 'public'
  notesDraft.download_visibility = item.download_visibility || 'public'
  notesOpen.value = true
}

async function saveNotes() {
  if (oidcEnabled.value && !isAdmin.value) {
    $q.notify({
      type: 'warning',
      message: 'Signed in, but owner role missing — sign out/in, or ask Keycloak admin for dasmlab-home/admin'
    })
    return
  }
  savingNotes.value = true
  try {
    await updateMedia(notesDayId.value, notesMediaId.value, {
      caption: notesDraft.caption,
      notes: notesDraft.notes,
      kind: notesDraft.kind,
      external_url: notesDraft.external_url,
      notes_visibility: notesDraft.notes_visibility,
      download_visibility: notesDraft.download_visibility
    })
    notesOpen.value = false
    await loadDays(notesDayId.value)
    if (viewerItem.value?.id === notesMediaId.value) {
      const day = days.value.find((d) => d.id === notesDayId.value)
      viewerItem.value = day?.media?.find((m) => m.id === notesMediaId.value) || viewerItem.value
    }
    $q.notify({ type: 'positive', message: 'Notes saved' })
  } catch (err) {
    const msg = err?.response?.data?.error || 'Could not save notes'
    const detail = err?.response?.data?.detail
    $q.notify({ type: 'negative', message: detail ? `${msg}: ${detail}` : msg })
  } finally {
    savingNotes.value = false
  }
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
    $q.notify({ type: 'warning', message: 'Only photos and videos here — use + link for Garmin / shares' })
    return
  }
  accepted.forEach((file) => {
    const queueItem = {
      id: `${dayId}-${file.name}-${Date.now()}-${Math.random()}`,
      dayId,
      file,
      progress: 0,
      status: 'queued',
      error: '',
      attempts: 0
    }
    uploadQueue.value.push(queueItem)
  })
  pumpUploadQueue()
}

const UPLOAD_CONCURRENCY = 3
const UPLOAD_MAX_ATTEMPTS = 3
let uploadWorkers = 0
let uploadReloadTimer = null

function pumpUploadQueue() {
  while (uploadWorkers < UPLOAD_CONCURRENCY) {
    const next = uploadQueue.value.find((i) => i.status === 'queued')
    if (!next) break
    uploadWorkers++
    next.status = 'uploading'
    runUpload(next).finally(() => {
      uploadWorkers--
      pumpUploadQueue()
    })
  }
}

async function runUpload(queueItem) {
  queueItem.attempts = (queueItem.attempts || 0) + 1
  try {
    await uploadMedia(queueItem.dayId, queueItem.file, {}, (progress) => {
      queueItem.progress = progress
    })
    queueItem.status = 'done'
    queueItem.progress = 100
    scheduleUploadReload(queueItem.dayId)
    // Drop successful rows after a short beat so the box doesn't grow forever.
    setTimeout(() => {
      uploadQueue.value = uploadQueue.value.filter((i) => i.id !== queueItem.id)
    }, 3500)
  } catch (err) {
    if (queueItem.attempts < UPLOAD_MAX_ATTEMPTS) {
      queueItem.status = 'queued'
      queueItem.error = `retry ${queueItem.attempts}/${UPLOAD_MAX_ATTEMPTS}…`
      queueItem.progress = 0
      await new Promise((r) => setTimeout(r, 600 * queueItem.attempts))
      return
    }
    queueItem.status = 'error'
    queueItem.error = err?.response?.data?.error || err?.message || 'Upload failed'
  }
}

function scheduleUploadReload(dayId) {
  if (uploadReloadTimer) clearTimeout(uploadReloadTimer)
  uploadReloadTimer = setTimeout(() => {
    loadDays(dayId)
    uploadReloadTimer = null
  }, 800)
}

function openViewer(item, day) {
  viewerItem.value = item
  viewerMuted.value = item?.media_type === 'video'
  viewerOpen.value = true
  if (item?.media_type === 'video' && day?.id && item?.id) {
    recordMediaPlay(day.id, item.id)
      .then(() => loadDays(day.id))
      .catch(() => {})
  }
}

async function publishAlbum(day) {
  if (!day?.id) return
  if (oidcEnabled.value && !isAdmin.value) {
    $q.notify({ type: 'warning', message: 'Sign in as owner to publish' })
    login()
    return
  }
  publishingDayId.value = day.id
  try {
    const result = await publishDayUntilDone(day.id)
    await loadDays(day.id)
    $q.notify({ type: 'positive', message: publishResultMessage(result) })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Publish failed' })
  } finally {
    publishingDayId.value = ''
  }
}

async function curateAndPublish(day) {
  if (!day?.id) return
  if (oidcEnabled.value && !isAdmin.value) {
    $q.notify({ type: 'warning', message: 'Sign in as owner to publish' })
    login()
    return
  }
  const pending = countPendingPublish(day)
  if (pending === 0) {
    $q.notify({
      type: 'info',
      message: 'Already on CDN — nothing to promote. (Compress/transcode adapter not online yet.)'
    })
    return
  }
  publishingDayId.value = day.id
  try {
    const result = await publishDayUntilDone(day.id)
    await loadDays(day.id)
    $q.notify({ type: 'positive', message: publishResultMessage(result) })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Curate publish failed' })
  } finally {
    publishingDayId.value = ''
  }
}

function countPendingPublish(day) {
  const media = day?.media || []
  return media.filter((m) => {
    if (m.hidden) return false
    if (m.origin === 'r2-draft') return true
    if (m.origin === 'pvc' || !m.origin) return !m.published
    if (m.published && String(m.url || '').startsWith('http')) return false
    return !m.published
  }).length
}

function publishButtonLabel(day) {
  const n = countPendingPublish(day)
  if (n === 0) return 'On CDN'
  return `Publish ${n}`
}

function publishButtonHint(day) {
  const n = countPendingPublish(day)
  if (n === 0) {
    return 'All clips already on CDN (original/). Compress/AI curate are separate.'
  }
  return `Promote ${n} draft/PVC item(s) → CDN original/. Direct uploads only need this flip — no re-upload.`
}

function publishResultMessage(result) {
  const pub = result?.published || 0
  const rem = result?.remaining || 0
  if (pub === 0 && rem === 0) return 'Already on CDN — nothing left to promote'
  if (rem === 0) return `Published ${pub} to CDN — album ready`
  return `Published ${pub}; ${rem} still pending — hit Publish again`
}

async function runAICurate(day) {
  try {
    const data = await aiCurate(day.id, {
      prompt: `Curate highlights for ${day.title}${day.location ? ' in ' + day.location : ''}`,
      action: 'save'
    })
    $q.notify({
      type: data?.warning ? 'warning' : 'positive',
      message: data?.curation?.brief || data?.message || 'AI curate accepted',
      timeout: 5000
    })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'AI curate failed' })
  }
}

async function restoreMedia(day, item) {
  try {
    await unhideMedia(day.id, item.id)
    await loadDays(day.id)
    $q.notify({ type: 'positive', message: 'Visible again in public gallery' })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not unhide' })
  }
}

async function proposeTag(day, item, name) {
  taggingMedia.value = true
  try {
    await proposeMediaTag(day.id, item.id, name)
    await loadDays(day.id)
    $q.notify({ type: 'positive', message: 'Name submitted — awaiting your approval' })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not tag' })
  } finally {
    taggingMedia.value = false
  }
}

async function moderateTag(day, item, tag, action) {
  try {
    await moderateMediaTag(day.id, item.id, tag.id, action)
    await loadDays(day.id)
    $q.notify({
      type: 'positive',
      message: action === 'approve' ? `Approved “${tag.name}”` : `Rejected “${tag.name}”`
    })
  } catch (err) {
    $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not update tag' })
  }
}

function onViewerHide() {
  viewerMuted.value = true
  const el = viewerVideoEl.value
  if (el) {
    el.pause()
  }
}

function toggleViewerMute() {
  viewerMuted.value = !viewerMuted.value
  const el = viewerVideoEl.value
  if (el) {
    el.muted = viewerMuted.value
    if (!viewerMuted.value) {
      el.volume = el.volume || 1
    }
  }
}

function onViewerVolumeChange(event) {
  const el = event?.target
  if (!el) return
  // Native controls unmute should clear our overlay state.
  if (!el.muted && viewerMuted.value) {
    viewerMuted.value = false
  }
}

function audioCleanupTeaser() {
  $q.notify({
    type: 'info',
    timeout: 5500,
    message: 'Premium later: AI clean / level / overlay audio — you keep the adapters and who can hear it. Close the share, Meta loses the tollgate.'
  })
}

async function removeMedia(day, item) {
  $q.dialog({
    title: 'Hide from album?',
    message: `${item.caption || item.filename} stays on your storage — it just leaves this gallery view.`,
    cancel: true,
    ok: { label: 'Hide', color: 'primary' }
  }).onOk(async () => {
    try {
      await deleteMedia(day.id, item.id)
      await loadDays(day.id)
    } catch (err) {
      $q.notify({ type: 'negative', message: err?.response?.data?.error || 'Could not hide media' })
    }
  })
}

onMounted(() => {
  loadDays()
})
</script>

<style scoped>
.surfing-page {
  max-width: 1240px;
  margin: 0 auto;
  --surf-deep: #063642;
  --surf-teal: #0f8f7c;
  --surf-foam: #dff7f1;
  --surf-horizon: #5eb4c8;
  --surf-ink: #102833;
  --surf-sand: #f4fbf9;
}

.surfing-hero {
  position: relative;
  overflow: hidden;
  border-color: rgba(15, 143, 124, 0.5);
  background:
    linear-gradient(125deg, #062f38 0%, #0a5c58 42%, #147a6c 68%, #1a9a86 100%);
  color: #f2fffb;
  box-shadow: 0 22px 48px rgba(6, 54, 66, 0.28);
  animation: hero-rise 700ms cubic-bezier(0.22, 1, 0.36, 1) both;
}

.surfing-hero :deep(.dasm-title),
.surfing-hero :deep(.dasm-subtitle),
.surfing-hero :deep(.dasm-caps) {
  color: inherit;
}

.surfing-hero :deep(.dasm-subtitle) {
  color: rgba(232, 252, 247, 0.88);
  max-width: 720px;
}

.surfing-hero__wash,
.surfing-hero__glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.surfing-hero__wash {
  background:
    radial-gradient(ellipse 70% 80% at 8% 120%, rgba(94, 180, 200, 0.45), transparent 55%),
    radial-gradient(ellipse 45% 50% at 95% -20%, rgba(255, 255, 255, 0.18), transparent 50%),
    repeating-linear-gradient(
      -12deg,
      transparent,
      transparent 22px,
      rgba(255, 255, 255, 0.04) 22px,
      rgba(255, 255, 255, 0.04) 23px
    );
  animation: surf-wash 12s ease-in-out infinite alternate;
}

.surfing-hero__glow {
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.12), transparent);
  transform: translateX(-120%) skewX(-18deg);
  animation: surf-sheen 7s ease-in-out infinite;
}

.surfing-hero__content {
  position: relative;
  z-index: 2;
}

.surfing-caps {
  color: #9feedd !important;
  letter-spacing: 0.22em;
}

.surfing-hero__title {
  color: #fff !important;
  text-shadow: 0 8px 28px rgba(0, 20, 28, 0.35);
}

.surfing-hero--themed {
  background-color: #062f38;
}

.surfing-hero__brief {
  margin: 0.75rem 0 0;
  max-width: 640px;
  font-size: 0.9rem;
  line-height: 1.5;
  color: rgba(220, 245, 238, 0.9);
  border-left: 3px solid rgba(159, 238, 221, 0.7);
  padding-left: 0.75rem;
}

.surfing-hero__actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.75rem 1rem;
  margin-top: 1.15rem;
}

.surfing-cta {
  font-weight: 600;
  letter-spacing: 0.02em;
}

.surfing-cta--ghost {
  color: #e8fcf7 !important;
  border-color: rgba(232, 252, 247, 0.45) !important;
}

.surfing-hero__hint {
  font-size: 0.82rem;
  color: rgba(220, 245, 238, 0.78);
}

.surfing-related {
  display: flex;
  align-items: center;
  gap: 0.9rem;
  text-decoration: none;
  color: inherit;
  padding: 0.95rem 1.1rem;
  transition: transform 160ms ease, border-color 160ms ease;
}

.surfing-related:hover {
  transform: translateX(3px);
}

.surfing-related__icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  background: rgba(6, 54, 66, 0.9);
  color: #9feedd;
}

.surfing-related__body {
  flex: 1;
  min-width: 0;
}

.surfing-related__title {
  font-weight: 650;
  color: var(--surf-ink, #12202c);
}

@keyframes hero-rise {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: none; }
}

@keyframes surf-wash {
  from { opacity: 0.75; transform: translateY(0) scale(1); }
  to { opacity: 1; transform: translateY(-6px) scale(1.02); }
}

@keyframes surf-sheen {
  0%, 40% { transform: translateX(-120%) skewX(-18deg); opacity: 0; }
  50% { opacity: 0.55; }
  70%, 100% { transform: translateX(140%) skewX(-18deg); opacity: 0; }
}

.surfing-alert {
  text-align: center;
  color: #5d7283;
}

.surfing-workspace {
  overflow: hidden;
  border-color: rgba(15, 143, 124, 0.2);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(240, 250, 247, 0.94));
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

.day-header__actions {
  display: flex;
  align-items: center;
  gap: 0.15rem;
}

.day-title-row {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.day-title {
  margin: 0;
  font-size: clamp(1.4rem, 2.4vw, 1.85rem);
  color: var(--surf-ink);
  letter-spacing: -0.02em;
}

.day-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.25rem;
  color: #607384;
  margin-top: 0.35rem;
}

.day-pill {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  color: var(--surf-teal);
  background: rgba(15, 143, 124, 0.1);
  border: 1px solid rgba(15, 143, 124, 0.25);
  border-radius: 6px;
  padding: 0.12rem 0.45rem;
}

.drop-zone {
  border: 2px dashed rgba(15, 143, 124, 0.42);
  border-radius: 18px;
  padding: 1.65rem;
  text-align: center;
  color: #3f5a69;
  background:
    linear-gradient(160deg, rgba(15, 143, 124, 0.1), rgba(94, 180, 200, 0.12));
  cursor: pointer;
  transition: border-color 180ms ease, transform 180ms ease, box-shadow 180ms ease;
}

.drop-zone__icon {
  color: var(--surf-teal);
}

.drop-zone:hover,
.drop-zone--active {
  border-color: rgba(15, 143, 124, 0.9);
  transform: translateY(-2px);
  box-shadow: 0 14px 28px rgba(6, 54, 66, 0.12);
}

.hidden-input {
  display: none;
}

.upload-queue {
  border: 1px solid rgba(36, 61, 81, 0.12);
  border-radius: 12px;
  padding: 0.65rem 0.75rem 0.5rem;
  background: rgba(255, 255, 255, 0.75);
  max-width: 100%;
}

.upload-queue__scroll {
  max-height: 220px;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 0.15rem;
  scrollbar-gutter: stable;
}

.queue-title {
  font-size: 0.78rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6a7f90;
  margin-bottom: 0.45rem;
  gap: 0.5rem;
}

.queue-title__stats {
  text-transform: none;
  letter-spacing: 0;
  font-weight: 500;
  color: #8a9eac;
}

.queue-item {
  display: flex;
  align-items: flex-start;
  gap: 0.65rem;
  padding: 0.4rem 0.15rem;
}

.queue-item + .queue-item {
  border-top: 1px solid rgba(36, 61, 81, 0.08);
}

.queue-item--done {
  opacity: 0.72;
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
  padding: 2.4rem 1rem;
}

.media-section {
  margin-top: 1.75rem;
  animation: section-in 520ms ease both;
}

.media-section--videos { animation-delay: 40ms; }
.media-section--photos { animation-delay: 90ms; }
.media-section--other { animation-delay: 140ms; }

@keyframes section-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: none; }
}

.media-section__head {
  margin-bottom: 0.9rem;
  padding-bottom: 0.55rem;
  border-bottom: 2px solid rgba(15, 143, 124, 0.18);
}

.media-section__label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--surf-deep);
}

.media-section__title {
  margin: 0;
  font-size: 0.82rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  font-weight: 800;
}

.media-section__count {
  font-size: 0.75rem;
  color: #6a8090;
  background: rgba(15, 143, 124, 0.1);
  border-radius: 999px;
  padding: 0.1rem 0.5rem;
}

.media-section__blurb {
  margin: 0.35rem 0 0;
  font-size: 0.86rem;
  color: #5d7484;
}

.video-rail {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.1rem;
}

.video-card {
  border-radius: 16px;
  overflow: hidden;
  background: #061820;
  color: #e8f7f3;
  box-shadow: 0 16px 36px rgba(6, 40, 52, 0.28);
  transition: transform 220ms ease, box-shadow 220ms ease;
}

.video-card:hover {
  transform: translateY(-4px) scale(1.01);
  box-shadow: 0 22px 44px rgba(6, 40, 52, 0.38);
}

.video-card__stage {
  position: relative;
  display: block;
  width: 100%;
  aspect-ratio: 16 / 10;
  padding: 0;
  border: 0;
  cursor: pointer;
  background: #0a2430;
  overflow: hidden;
}

.video-card__asset {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  opacity: 0.92;
  transition: transform 360ms ease, opacity 220ms ease;
}

.video-card:hover .video-card__asset {
  transform: scale(1.05);
  opacity: 1;
}

.video-card__play {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  color: #fff;
  background: radial-gradient(circle, rgba(15, 143, 124, 0.35), transparent 60%);
}

.video-card__dur {
  position: absolute;
  left: 0.55rem;
  top: 0.45rem;
  z-index: 2;
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.02em;
  color: #fff;
  background: rgba(6, 24, 32, 0.78);
  border-radius: 6px;
  padding: 0.12rem 0.4rem;
  font-variant-numeric: tabular-nums;
}

.video-sort-bar {
  gap: 0.35rem;
}

.video-card__play :deep(.q-icon) {
  background: rgba(6, 24, 32, 0.72);
  border-radius: 999px;
  padding: 0.35rem;
  border: 2px solid rgba(255, 255, 255, 0.55);
}

.video-card__open {
  position: absolute;
  right: 0.7rem;
  bottom: 0.7rem;
  font-size: 0.68rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  background: rgba(15, 143, 124, 0.9);
  color: #fff;
}

.video-card__body {
  padding: 0.75rem 0.85rem 0.55rem;
}

.media-card--hidden {
  opacity: 0.45;
  filter: grayscale(0.35);
  outline: 1px dashed rgba(90, 112, 128, 0.45);
}

.media-card--hidden::after {
  content: 'Hidden from public';
  display: block;
  font-size: 0.7rem;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #6a8090;
  padding: 0.25rem 0.65rem 0.55rem;
}

.video-card__title {
  font-weight: 650;
  font-size: 0.92rem;
}

.video-card__meta {
  margin-top: 0.2rem;
  color: rgba(210, 236, 230, 0.75);
}

.video-card__notes,
.photo-card__notes,
.other-card__notes {
  margin: 0.35rem 0 0;
  font-size: 0.8rem;
  line-height: 1.45;
  color: rgba(210, 236, 230, 0.82);
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.other-card__notes {
  color: #5a7080;
}

.video-card__actions,
.photo-card__actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.15rem;
  margin-top: 0.45rem;
}

.photo-rail {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.1rem;
}

.photo-card {
  border-radius: 16px;
  overflow: hidden;
  background: linear-gradient(165deg, #10262e 0%, #18343c 55%, #1c2e28 100%);
  color: #eef7f3;
  box-shadow: 0 16px 36px rgba(6, 40, 52, 0.26);
  transition: transform 220ms ease, box-shadow 220ms ease;
}

.photo-card:hover {
  transform: translateY(-4px) scale(1.01);
  box-shadow: 0 22px 44px rgba(6, 40, 52, 0.36);
}

.photo-card__stage {
  position: relative;
  display: block;
  width: 100%;
  aspect-ratio: 4 / 3;
  padding: 0;
  border: 0;
  cursor: pointer;
  background: #0a2430;
  overflow: hidden;
}

.photo-card__asset {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  opacity: 0.94;
  transition: transform 360ms ease, opacity 220ms ease;
}

.photo-card:hover .photo-card__asset {
  transform: scale(1.06);
  opacity: 1;
}

.photo-card__glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 80% 15%, rgba(255, 186, 96, 0.28), transparent 42%),
    linear-gradient(180deg, transparent 45%, rgba(6, 24, 32, 0.55));
}

.photo-card__open {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  color: #fff;
  opacity: 0;
  background: radial-gradient(circle, rgba(15, 143, 124, 0.32), transparent 62%);
  transition: opacity 180ms ease;
}

.photo-card:hover .photo-card__open {
  opacity: 1;
}

.photo-card__badge {
  position: absolute;
  left: 0.55rem;
  top: 0.45rem;
  z-index: 2;
  font-size: 0.68rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #fff;
  background: rgba(15, 143, 124, 0.88);
  border-radius: 6px;
  padding: 0.12rem 0.4rem;
}

.photo-card__badge--notes {
  background: rgba(176, 122, 18, 0.9);
}

.photo-card__body {
  padding: 0.75rem 0.85rem 0.85rem;
}

.photo-card__title {
  font-weight: 700;
  font-size: 0.95rem;
  line-height: 1.25;
  color: #eef7f3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.photo-card__meta {
  margin-top: 0.2rem;
  color: rgba(210, 236, 230, 0.75);
}

.photo-card__notes {
  color: rgba(220, 236, 228, 0.88);
}

.photo-card__actions :deep(.q-btn) {
  color: rgba(210, 236, 230, 0.92);
}

.media-card--hidden {
  opacity: 0.55;
  filter: grayscale(0.35);
}

.other-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.other-card {
  display: flex;
  align-items: flex-start;
  gap: 0.85rem;
  padding: 0.9rem 1rem;
  border-radius: 14px;
  background: linear-gradient(120deg, rgba(15, 143, 124, 0.08), rgba(255, 255, 255, 0.95));
  border: 1px solid rgba(15, 143, 124, 0.18);
  transition: transform 180ms ease, border-color 180ms ease;
}

.other-card:hover {
  transform: translateX(3px);
  border-color: rgba(15, 143, 124, 0.4);
}

.other-card__icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  background: rgba(6, 54, 66, 0.9);
  color: #9feedd;
  flex-shrink: 0;
}

.other-card__body {
  flex: 1;
  min-width: 0;
}

.other-card__title {
  font-weight: 650;
  color: var(--surf-ink);
}

.other-card__link {
  display: inline-block;
  margin-top: 0.4rem;
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--surf-teal);
  text-decoration: none;
  letter-spacing: 0.02em;
}

.other-card__link:hover {
  text-decoration: underline;
}

.other-card__actions {
  display: flex;
  flex-shrink: 0;
}

.viewer-card {
  background: #0a141a;
  color: #fff;
}

.viewer-bar {
  background: rgba(0, 0, 0, 0.55);
}

.viewer-body {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 120px);
  background: #071018;
}

.viewer-video-wrap {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.viewer-unmute {
  position: absolute;
  left: 50%;
  bottom: 1.25rem;
  transform: translateX(-50%);
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  border: 1px solid rgba(255, 255, 255, 0.35);
  border-radius: 999px;
  padding: 0.55rem 0.95rem;
  background: rgba(6, 24, 32, 0.82);
  color: #e8fcf7;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  backdrop-filter: blur(6px);
  z-index: 2;
}

.viewer-unmute:hover {
  background: rgba(15, 143, 124, 0.85);
  border-color: rgba(255, 255, 255, 0.55);
}

.viewer-asset {
  max-width: 100%;
  max-height: calc(100vh - 140px);
  object-fit: contain;
}

.viewer-notes {
  max-width: 720px;
  margin: 0 auto;
  color: rgba(220, 240, 235, 0.88);
  line-height: 1.5;
}

@media (max-width: 760px) {
  .photo-rail,
  .video-rail {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }

  .photo-card__stage {
    aspect-ratio: 1;
  }
}
</style>
