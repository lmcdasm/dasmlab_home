<template>
  <article class="stage" :style="{ '--stage-accent': accent }">
    <div class="stage__visual">
      <ProjectVisual
        :problem="problem"
        :outcome="outcome"
        :techs="techs"
        :accent="accent"
      />
    </div>
    <div class="stage__body">
      <div class="stage__badges">
        <span v-if="lane" class="badge badge--lane">{{ lane }}</span>
        <span v-if="badge" class="badge badge--status">{{ badge }}</span>
      </div>
      <h3 class="stage__title">
        <router-link v-if="hubPath" :to="hubPath" class="stage__link">{{ title }}</router-link>
        <a v-else-if="url" :href="url" class="stage__link" target="_blank" rel="noopener">{{ title }}</a>
        <span v-else>{{ title }}</span>
      </h3>
      <p class="stage__desc">{{ description }}</p>
      <div class="stage__meta">
        <span class="meta"><strong>Why</strong> {{ problem }}</span>
        <span class="meta"><strong>Stack</strong> {{ techs.join(' · ') }}</span>
      </div>
      <div class="stage__actions">
        <router-link v-if="hubPath" class="action action--primary" :to="hubPath">Behind the Design</router-link>
        <a v-if="liveUrl" class="action" :href="liveUrl" target="_blank" rel="noopener">Live</a>
        <a v-if="url && !hubPath" class="action" :href="url" target="_blank" rel="noopener">Source</a>
      </div>
    </div>
  </article>
</template>

<script setup>
import ProjectVisual from 'src/components/ProjectVisual.vue'

defineProps({
  title: String,
  description: String,
  problem: { type: String, default: 'Real-world constraint' },
  outcome: { type: String, default: 'Shippable pattern' },
  techs: { type: Array, default: () => [] },
  lane: String,
  badge: String,
  hubPath: String,
  liveUrl: String,
  url: String,
  accent: { type: String, default: '#1f6f62' }
})
</script>

<style scoped>
.stage {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 0.95fr);
  gap: 0;
  border: 1.5px solid rgba(31, 111, 98, 0.28);
  border-radius: 18px;
  overflow: hidden;
  background: linear-gradient(160deg, #ffffff, #f3f8f6);
  box-shadow: 0 14px 32px rgba(18, 40, 52, 0.08);
  transition: transform 180ms ease, box-shadow 180ms ease;
}

.stage:hover {
  transform: translateY(-2px);
  box-shadow: 0 18px 36px rgba(18, 40, 52, 0.12);
}

.stage__visual {
  padding: 0.65rem;
  background:
    radial-gradient(circle at 20% 20%, rgba(31, 111, 98, 0.16), transparent 50%),
    linear-gradient(165deg, #eef6f3, #e4edf2);
  border-right: 1px solid rgba(28, 52, 73, 0.1);
  min-height: 210px;
  display: grid;
  align-items: center;
}

.stage__body {
  padding: 1.05rem 1.15rem 1.15rem;
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
}

.stage__badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
}

.badge {
  font-size: 0.68rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  padding: 0.28rem 0.55rem;
  border-radius: 8px;
}

.badge--lane {
  background: rgba(31, 111, 98, 0.12);
  color: var(--stage-accent);
  border: 1px solid rgba(31, 111, 98, 0.35);
}

.badge--status {
  background: #fff;
  border: 1px solid rgba(28, 52, 73, 0.18);
  color: #435564;
}

.stage__title {
  margin: 0;
  font-family: 'Fraunces', Georgia, serif;
  font-size: clamp(1.15rem, 1.8vw, 1.4rem);
  font-weight: 700;
  color: #12202c;
  line-height: 1.15;
}

.stage__link {
  color: inherit;
  text-decoration: none;
}

.stage__link:hover {
  color: var(--stage-accent);
}

.stage__desc {
  margin: 0;
  color: #3d5263;
  line-height: 1.5;
  font-size: 0.95rem;
}

.stage__meta {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.meta {
  font-size: 0.82rem;
  color: #5a6f80;
  line-height: 1.4;
}

.meta strong {
  color: #1d2b36;
  margin-right: 0.35rem;
  font-weight: 700;
}

.stage__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-top: 0.25rem;
}

.action {
  display: inline-flex;
  align-items: center;
  padding: 0.42rem 0.75rem;
  border-radius: 10px;
  border: 1.5px solid rgba(31, 111, 98, 0.3);
  color: #1d2b36;
  text-decoration: none;
  font-weight: 650;
  font-size: 0.86rem;
  background: rgba(255, 255, 255, 0.85);
}

.action--primary {
  background: linear-gradient(135deg, #1f6f62, #2f8f7d);
  border-color: transparent;
  color: #fff;
}

@media (max-width: 860px) {
  .stage {
    grid-template-columns: 1fr;
  }

  .stage__visual {
    border-right: none;
    border-bottom: 1px solid rgba(28, 52, 73, 0.1);
    min-height: 180px;
  }
}
</style>
