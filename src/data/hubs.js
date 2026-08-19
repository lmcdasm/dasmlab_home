/**
 * Engineering Knowledge Network — hub inventory for dasmlab.org 2.0
 */

/** Public vanity hosts only — never apps.* or dev-* on visitor CTAs. */
export const PRODUCT_URLS = {
  home: 'https://dasmlab.org',
  surfing: 'https://dasmlab.org/surfing',
  mockMe: 'https://mock-me.dasmlab.org',
  mockMeDemo: 'https://mock-me.dasmlab.org/demo',
  interviewMe: 'https://interview-me.dasmlab.org',
  interviewMeDemo: 'https://interview-me.dasmlab.org/demo',
  cameraScrape: 'https://camera-scrape.dasmlab.org/',
  designCarousel: 'https://design-carousel.svc.dasmlab.org/carousel'
}

export const SITE = {
  name: 'DASMLAB',
  legalName: 'Technologies DASMLAB Inc.',
  url: 'https://dasmlab.org',
  description:
    'Living lab and engineering knowledge network — cloud-native, AI, and infrastructure builds with how-we-built-it depth.',
  contactPath: '/contact',
  sameAs: [
    'https://github.com/lmcdasm',
    'https://github.com/dasmlab'
  ]
}

/** About photos on Cloudflare R2 (bucket dasmlab-home). Override with VITE_ABOUT_MEDIA_BASE. */
export const ABOUT_MEDIA_BASE = (
  import.meta.env.VITE_ABOUT_MEDIA_BASE ||
  'https://pub-29bde7a836c744729bebe74bfd4008a2.r2.dev/about'
).replace(/\/$/, '')

export const AUTHOR = {
  name: 'Daniel Smith',
  alternateName: 'dasm',
  jobTitle: 'Founder / Engineer',
  description:
    'Two decades across telecom, automotive, cloud, and AI. Based in the Laurentians after Outaouais/Ottawa and Fernie. Windsurfer, snowboarder, gardener.',
  url: 'https://dasmlab.org/about',
  image: `${ABOUT_MEDIA_BASE}/portrait.webp`
}

/** @typedef {{ question: string, answer: string }} FaqItem */

/**
 * @typedef {object} ProjectHub
 * @property {string} slug
 * @property {string} title
 * @property {string} lane
 * @property {string} summary
 * @property {string} answer
 * @property {string[]} stack
 * @property {string[]} topics
 * @property {string} [liveUrl]
 * @property {string} [demoUrl]
 * @property {string} [sourceUrl]
 * @property {string} [ceUrl]
 * @property {string} architecture
 * @property {string} howWeBuilt
 * @property {FaqItem[]} faq
 */

/** @type {Record<string, ProjectHub>} */
export const projectHubs = {
  'dasmlab-home': {
    slug: 'dasmlab-home',
    title: 'dasmlab-home',
    lane: 'Frontend',
    summary: 'The public face of Technologies DASMLAB Inc. — portfolio, Surfing CDN client, and Engineering Knowledge Network.',
    answer:
      'dasmlab-home is a Vue 3 + Quasar SPA served by nginx on OpenShift, backed by surfing-service for media and OIDC, designed as a two-facet site: visitor answers and engineer-level how-we-built depth.',
    stack: ['Vue 3', 'Quasar', 'Vite', 'nginx', 'OpenShift'],
    topics: ['vue', 'quasar', 'oidc', 'openshift'],
    liveUrl: PRODUCT_URLS.home,
    sourceUrl: 'https://github.com/lmcdasm/dasmlab_home',
    architecture:
      'Static Quasar build in nginx; `/api/surfing` reverse-proxied to surfing-service; Keycloak realm `dasmlab` for owner features; public Activity POST for anonymous engagement; R2/CDN path for media bytes.',
    howWeBuilt:
      'We iterated from card-grid portfolio to Lab Constellation (hub-spoke map), then 2.0 hubs for topical authority. Auth follows the same OIDC cookie pattern as mock-me / interview-me. SEO uses history routing, sitemap, and JSON-LD.',
    faq: [
      {
        question: 'Is dasmlab-home open source?',
        answer:
          'The site repository is available for study. Production secrets, Keycloak ops, and CDN credentials stay private. Shared libraries move to Apache 2.0 CE over time.'
      },
      {
        question: 'How do I explore sibling products without Keycloak?',
        answer:
          'Use the labeled Demo / fake mode on each product (see project hubs). Demos never mutate live systems.'
      }
    ]
  },
  surfing: {
    slug: 'surfing',
    title: 'Surfing CDN',
    lane: 'Backend',
    summary: 'Personal media days, albums, and CDN-backed playback — the golden client for DASMLAB media.',
    answer:
      'Surfing is a Go API plus Vue client that stores media metadata on OpenShift and serves bytes from object storage / CDN, with public browse and owner publish controls via OIDC.',
    stack: ['Go', 'Gin', 'Vue 3', 'Quasar', 'Object storage', 'OIDC'],
    topics: ['gin', 'vue', 'oidc', 'openshift'],
    liveUrl: PRODUCT_URLS.surfing,
    architecture:
      'surfing-service (Gin) owns auth, days, upload orchestration, and Activity. Bytes migrate off PVC to object store + CDN URLs in manifests. Home nginx proxies `/api/surfing`.',
    howWeBuilt:
      'Started as basement PVC streaming; 0.7+ work moves origin to cheap object storage. Activity Phase A added anonymous visitor cookies and public POST for engagement without forcing login.',
    faq: [
      {
        question: 'Can visitors upload?',
        answer: 'No. Browse is public; publish and admin require the owner OIDC session.'
      },
      {
        question: 'Where is the source?',
        answer:
          'Core lives inside dasmlab_home/surfing-service. A future CE extract may publish libraries under Apache 2.0; production deploy stays private.'
      }
    ]
  },
  'mock-me': {
    slug: 'mock-me',
    title: 'Mock-Me',
    lane: 'Cloud',
    summary: 'Orchestrate and mock CDN / platform workflows — live ops behind Keycloak, showcase via fake demo mode.',
    answer:
      'Mock-Me lets operators mock and orchestrate workflows that need a CDN or platform path. Unknown visitors use Demo / fake mode: scripted steps with zero live node deploys.',
    stack: ['Go', 'Vue 3', 'Quasar', 'Keycloak', 'OpenShift'],
    topics: ['gin', 'vue', 'oidc', 'openshift'],
    liveUrl: PRODUCT_URLS.mockMe,
    demoUrl: PRODUCT_URLS.mockMeDemo,
    architecture:
      'OIDC-gated API for real orchestration; parallel `/demo` simulate endpoints return fixture timelines. Activity events tagged demo=true for operator filters.',
    howWeBuilt:
      'Built as a private IP-protecting product. 2.0 adds the shared Demo Visitor Contract so we can showcase without opening live mutate paths.',
    faq: [
      {
        question: 'Will the demo deploy to a real node?',
        answer: 'Never. Demo mode only returns scripted step JSON and is labeled in the UI.'
      },
      {
        question: 'Is there a Community Edition?',
        answer:
          'CE is planned (MPL 2.0 + commercial dual license). Scaffolds live under open-core/mock-me-ce until github.com/dasmlab/mock-me-ce is published. SaaS and secrets stay private.'
      }
    ]
  },
  'interview-me': {
    slug: 'interview-me',
    title: 'Interview-Me',
    lane: 'AI/ML',
    summary: 'Structured technical interviews with admin OIDC and invite guests — plus a public synthetic demo.',
    answer:
      'Interview-Me runs invite-based interview sessions with OTP guest cookies for candidates, while admins use Keycloak. A public demo session uses synthetic content with no PII writes.',
    stack: ['Go', 'Vue 3', 'Quasar', 'Keycloak', 'SMTP'],
    topics: ['vue', 'gin', 'oidc', 'quasar'],
    liveUrl: PRODUCT_URLS.interviewMe,
    demoUrl: PRODUCT_URLS.interviewMeDemo,
    architecture:
      'Admin OIDC + invite OTP → im_guest cookie scoped to one session. Public demo mints im_demo with a fixed synthetic session id.',
    howWeBuilt:
      'Guest invite flow was the first safe “non-admin” path. 2.0 extends that idea to a Keycloak-free public demo for showcase and SEO funnels from dasmlab.org.',
    faq: [
      {
        question: 'Does the public demo store candidate data?',
        answer: 'No. It uses synthetic fixtures only; answers are not persisted as PII.'
      },
      {
        question: 'Is there a Community Edition?',
        answer:
          'CE is planned under open-core/interview-me-ce (MPL + commercial). Public GitHub publish comes after the dual-license extract is ready.'
      }
    ]
  },
  cheapcloud: {
    slug: 'cheapcloud',
    title: 'CheapCloud',
    lane: 'Cloud',
    summary: 'Spend envelopes, free-tier burn, and provider recommendations — operator UI behind edge auth today.',
    answer:
      'CheapCloud helps pick cost-aware cloud origins and envelopes. The public site documents the product; the live operator UI stays behind edge authentication until a labeled public demo path is opened.',
    stack: ['Go', 'Vue', 'Azure/AWS/GCP abstractions'],
    topics: ['openshift', 'gin'],
    architecture:
      'API + UI with provider interfaces. Production sits behind edge basic auth. Future demo will use AllowDemoRead + fixture recommend with no cloud credentials.',
    howWeBuilt:
      'Started Azure-shaped / dry-run; evolving into cheap origin broker for Surfing and friends. CE will publish abstractions without prod keys.',
    faq: [
      {
        question: 'Can I try CheapCloud without login?',
        answer:
          'Not yet on the public edge — operator login is required. Use this hub for the how-we-built story; a labeled fake-mode demo is on the roadmap.'
      },
      {
        question: 'Does demo call my cloud account?',
        answer: 'No. When public demo ships it will never inject credentials or mutate cloud resources.'
      }
    ]
  },
  'camera-scrape': {
    slug: 'camera-scrape',
    title: 'Camera Scrape / Live Cams',
    lane: 'Infrastructure',
    summary:
      'Cheap multi-source camera scrape, real-time snapshots, and timelapse stitch — field cams for spots you ride and lab scenes you watch.',
    answer:
      'camera-scrape pulls frames from heterogeneous camera sources on a budget stack, stores snapshots, and stitches timelapses so Surfing / Live Cams can show Ste-Agathe, Jibe City, and other spots without a heavyweight VMS.',
    stack: ['Go', 'OpenShift', 'object storage', 'cron/schedulers', 'image pipeline'],
    topics: ['openshift', 'gin', 'oidc'],
    liveUrl: PRODUCT_URLS.cameraScrape,
    architecture:
      'Scrapers + schedule workers on OpenShift; public gallery for sample / published cams; private camera URLs and operator controls stay proprietary. Identity boundary separates public gallery from private scrapes.',
    howWeBuilt:
      'Started as a practical way to keep wind / snow / garden cams without SaaS camera platforms. 2.0 adds a project hub; vanity host is additive once DNS+CERT are proven — Live CTAs keep the working apps.* route until then.',
    faq: [
      {
        question: 'Are private camera URLs public?',
        answer: 'No. Public gallery shows allowed samples; private endpoints and schedules stay operator-gated.'
      },
      {
        question: 'Is there a Community Edition?',
        answer:
          'Scaffolded under open-core/camera-scrape-ce (MPL + commercial). Sample gallery sketches publish; production schedules and camera credentials never leave the private tree.'
      }
    ]
  }
}

/** @type {Record<string, { slug: string, title: string, summary: string, answer: string, projects: string[], faq: FaqItem[] }>} */
export const topicHubs = {
  vue: {
    slug: 'vue',
    title: 'Vue 3',
    summary: 'Composition API SPAs across DASMLAB products.',
    answer:
      'We standardize on Vue 3 + Composition API with Quasar for product UIs (home, interview-me, mock-me) so patterns like demo mode and Activity trackers port cleanly.',
    projects: ['dasmlab-home', 'surfing', 'mock-me', 'interview-me'],
    faq: [
      {
        question: 'Why Quasar instead of bare Vite?',
        answer: 'Shared layout, router, and mobile-ready primitives across multi-sites without rewriting chrome each time.'
      }
    ]
  },
  quasar: {
    slug: 'quasar',
    title: 'Quasar',
    summary: 'UI kit and build pipeline for DASMLAB frontends.',
    answer: 'Quasar gives us Vite builds, nginx-static deploy, and consistent components across the constellation.',
    projects: ['dasmlab-home', 'interview-me', 'mock-me'],
    faq: []
  },
  gin: {
    slug: 'gin',
    title: 'Gin (Go)',
    summary: 'HTTP APIs for surfing, mock-me, and interview-me.',
    answer:
      'Gin handlers plus cookie sessions keep OIDC and demo middleware explicit: RequireAuth, AllowDemoRead, DenyDemoMutate.',
    projects: ['surfing', 'mock-me', 'interview-me', 'cheapcloud', 'camera-scrape'],
    faq: []
  },
  oidc: {
    slug: 'oidc',
    title: 'OIDC / Keycloak',
    summary: 'Identity for live ops — never required for labeled demos.',
    answer:
      'Live mutate paths use Keycloak (realm dasmlab). Demo visitors use product demo cookies and never receive admin roles.',
    projects: ['dasmlab-home', 'surfing', 'mock-me', 'interview-me', 'camera-scrape'],
    faq: [
      {
        question: 'Can I use the products without Keycloak?',
        answer: 'Yes — via Demo / fake mode or (interview-me) invite guest links for real sessions.'
      }
    ]
  },
  metallb: {
    slug: 'metallb',
    title: 'MetalLB',
    summary: 'Bare-metal LoadBalancer patterns in the DASMLAB lab.',
    answer:
      'Lab clusters use MetalLB for service exposure where cloud LBs are unavailable; docs link from infra projects and OpenShift lab notes.',
    projects: ['dasmlab-home', 'camera-scrape'],
    faq: []
  },
  openshift: {
    slug: 'openshift',
    title: 'OpenShift / GitOps',
    summary: 'Deploy envelope for the constellation.',
    answer:
      'Apps ship via container images and GitOps overlays (dasmlab-live-cicd). Production overlays and secrets stay private; CE shows example envelopes only.',
    projects: ['dasmlab-home', 'surfing', 'mock-me', 'cheapcloud', 'camera-scrape'],
    faq: []
  }
}

export const labs = {
  'activity-anon-cdp': {
    slug: 'activity-anon-cdp',
    title: 'Anonymous visitor CDP',
    summary: 'Public Activity spine for sites that are not login-walled.',
    answer:
      'Anonymous IDs and session cookies plus public POST /activity let dasmlab.org measure engagement without Keycloak, while GET/panel stay owner-gated.',
    relatedProjects: ['dasmlab-home', 'surfing'],
    faq: []
  },
  'demo-visitor-facade': {
    slug: 'demo-visitor-facade',
    title: 'Demo visitor facade',
    summary: 'Cross-product fake mode that never mutates live systems.',
    answer:
      'Shared contract: labeled banner, demo cookie, simulate-only writes, DenyDemoMutate on live APIs. Reference implementation starts on interview-me then mock-me.',
    relatedProjects: ['interview-me', 'mock-me', 'cheapcloud'],
    faq: []
  },
  'surfing-r2-origin': {
    slug: 'surfing-r2-origin',
    title: 'Surfing object origin',
    summary: 'Move media bytes off basement PVC to object storage + CDN.',
    answer:
      'API stays on OpenShift; manifests carry CDN URLs so playback no longer depends on cluster disk for hot paths.',
    relatedProjects: ['surfing', 'cheapcloud', 'camera-scrape'],
    faq: []
  }
}

export function listProjectHubs() {
  return Object.values(projectHubs)
}

export function listTopicHubs() {
  return Object.values(topicHubs)
}

export function listLabs() {
  return Object.values(labs)
}
