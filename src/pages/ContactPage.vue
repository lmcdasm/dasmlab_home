<template>
  <q-page class="contact-page">
    <section class="contact-hero">
      <div class="hero-glow" />
      <div class="hero-copy">
        <div class="dasm-caps">Contact</div>
        <h1 class="hero-title">From requirements to MVP — let's build it properly.</h1>
        <p class="hero-sub">
          Discovery, architecture, prototypes, and pragmatic roadmaps. Commercial licensing and
          support for DASMLAB CE products welcome.
        </p>
      </div>
      <svg class="hero-svg" viewBox="0 0 360 220" aria-hidden="true">
        <defs>
          <linearGradient id="cg" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#1f6f62" stop-opacity="0.35" />
            <stop offset="100%" stop-color="#74865c" stop-opacity="0.15" />
          </linearGradient>
        </defs>
        <rect width="360" height="220" rx="18" fill="url(#cg)" />
        <circle cx="90" cy="110" r="34" fill="#fff" stroke="#1f6f62" stroke-width="3" />
        <circle cx="180" cy="70" r="26" fill="#fff" stroke="#2f8f7d" stroke-width="2.5" />
        <circle cx="250" cy="130" r="30" fill="#fff" stroke="#74865c" stroke-width="2.5" />
        <path d="M120 105 C145 90 155 85 160 78" stroke="#1f6f62" stroke-width="2" fill="none" stroke-dasharray="5 4" />
        <path d="M200 85 C220 100 225 115 230 120" stroke="#2f8f7d" stroke-width="2" fill="none" stroke-dasharray="5 4" />
        <text x="90" y="114" text-anchor="middle" font-size="11" font-weight="700" fill="#12202c">You</text>
        <text x="180" y="74" text-anchor="middle" font-size="10" font-weight="700" fill="#12202c">Scope</text>
        <text x="250" y="134" text-anchor="middle" font-size="10" font-weight="700" fill="#12202c">Ship</text>
      </svg>
    </section>

    <section class="engage-grid">
      <article v-for="item in engagements" :key="item.title" class="engage-card">
        <div class="engage-icon" aria-hidden="true">{{ item.glyph }}</div>
        <h2>{{ item.title }}</h2>
        <p>{{ item.copy }}</p>
      </article>
    </section>

    <section class="cta-panel">
      <div>
        <h2>Start a focused conversation</h2>
        <p>
          Send context, constraints, and timeline. For CE commercial licensing, say so in the subject —
          Technologies DASMLAB Inc.
        </p>
      </div>
      <div class="cta-actions">
        <a class="cta-primary" href="mailto:lmcdasm@gmail.com?subject=DASMLAB%20engagement">
          Email lmcdasm@gmail.com
        </a>
        <router-link class="cta-ghost" to="/projects/dasmlab-home">See how we build</router-link>
      </div>
    </section>
  </q-page>
</template>

<script setup>
import { useSeo } from 'src/composables/useSeo'

useSeo({
  title: 'Contact',
  description:
    'Contact Technologies DASMLAB Inc. for architecture, prototypes, commercial licensing, and support.',
  path: '/contact'
})

const engagements = [
  {
    glyph: '01',
    title: 'Requirements vetting',
    copy: 'Pressure-test scope before you burn runway — constraints, risks, and the smallest valuable slice.'
  },
  {
    glyph: '02',
    title: 'Architecture & design',
    copy: 'Trade-offs across OpenShift, CDN origins, OIDC, and product boundaries — with Behind the Design depth.'
  },
  {
    glyph: '03',
    title: 'Prototype strategy',
    copy: 'Lab-grade MVPs that showcase without mutating live systems — demo facades included.'
  },
  {
    glyph: '04',
    title: 'Commercial & support',
    copy: 'Dual-license CE products, white-label, hosted options — see COMMERCIAL terms via project hubs.'
  }
]
</script>

<style scoped>
.contact-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0.75rem clamp(0.7rem, 2vw, 1.2rem) 2rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.contact-hero {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 1rem;
  align-items: center;
  position: relative;
  overflow: hidden;
  border: 1.5px solid rgba(31, 111, 98, 0.35);
  border-radius: 20px;
  padding: 1.25rem 1.3rem;
  background:
    radial-gradient(circle at 8% 20%, rgba(47, 143, 125, 0.16), transparent 45%),
    linear-gradient(155deg, #fff, #eef5f2);
  box-shadow: 0 16px 34px rgba(18, 40, 52, 0.1);
}

.hero-glow {
  position: absolute;
  inset: auto -10% -40% auto;
  width: 280px;
  height: 280px;
  background: rgba(116, 134, 92, 0.18);
  filter: blur(40px);
  border-radius: 50%;
}

.hero-title {
  margin: 0.25rem 0 0.45rem;
  font-family: 'Fraunces', Georgia, serif;
  font-size: clamp(1.45rem, 2.6vw, 2.1rem);
  line-height: 1.1;
  color: #12202c;
}

.hero-sub {
  margin: 0;
  color: #3d5263;
  line-height: 1.5;
  max-width: 36rem;
}

.hero-svg {
  width: 100%;
  max-width: 360px;
  justify-self: end;
  filter: drop-shadow(0 10px 18px rgba(18, 40, 52, 0.08));
}

.engage-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.75rem;
}

.engage-card {
  border: 1.5px solid rgba(31, 111, 98, 0.22);
  border-radius: 16px;
  padding: 1rem;
  background: linear-gradient(165deg, #fff, #f5f9f7);
  box-shadow: 0 10px 22px rgba(18, 40, 52, 0.06);
  transition: transform 160ms ease;
}

.engage-card:hover {
  transform: translateY(-2px);
}

.engage-icon {
  font-family: 'Fraunces', Georgia, serif;
  font-weight: 700;
  color: #1f6f62;
  font-size: 1.1rem;
  margin-bottom: 0.35rem;
}

.engage-card h2 {
  margin: 0 0 0.35rem;
  font-size: 1rem;
  font-family: 'Fraunces', Georgia, serif;
  color: #12202c;
}

.engage-card p {
  margin: 0;
  font-size: 0.88rem;
  color: #4a5d6d;
  line-height: 1.45;
}

.cta-panel {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 1rem;
  align-items: center;
  border-radius: 18px;
  padding: 1.2rem 1.3rem;
  background: linear-gradient(135deg, #1f6f62, #2f8f7d 55%, #74865c);
  color: #fff;
  box-shadow: 0 16px 32px rgba(31, 111, 98, 0.28);
}

.cta-panel h2 {
  margin: 0 0 0.35rem;
  font-family: 'Fraunces', Georgia, serif;
  font-size: 1.25rem;
}

.cta-panel p {
  margin: 0;
  opacity: 0.92;
  max-width: 36rem;
  line-height: 1.45;
}

.cta-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.cta-primary,
.cta-ghost {
  display: inline-flex;
  align-items: center;
  padding: 0.55rem 0.95rem;
  border-radius: 11px;
  font-weight: 700;
  text-decoration: none;
}

.cta-primary {
  background: #fff;
  color: #1f6f62;
}

.cta-ghost {
  border: 1.5px solid rgba(255, 255, 255, 0.65);
  color: #fff;
}

@media (max-width: 960px) {
  .engage-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 720px) {
  .contact-hero {
    grid-template-columns: 1fr;
  }

  .hero-svg {
    justify-self: center;
  }

  .engage-grid {
    grid-template-columns: 1fr;
  }
}
</style>
