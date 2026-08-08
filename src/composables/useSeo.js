import { onBeforeUnmount, watch } from 'vue'
import { SITE, AUTHOR } from 'src/data/hubs'

const DEFAULT_OG = `${SITE.url}/icons/dasmlab-icon-128x128.png`

function ensureMeta(selector, attrs) {
  let el = document.head.querySelector(selector)
  if (!el) {
    el = document.createElement('meta')
    document.head.appendChild(el)
  }
  Object.entries(attrs).forEach(([k, v]) => {
    if (v != null) el.setAttribute(k, v)
  })
  return el
}

function ensureLink(rel, href) {
  let el = document.head.querySelector(`link[rel="${rel}"]`)
  if (!el) {
    el = document.createElement('link')
    el.setAttribute('rel', rel)
    document.head.appendChild(el)
  }
  el.setAttribute('href', href)
  return el
}

function setJsonLd(id, data) {
  const sid = `jsonld-${id}`
  let el = document.getElementById(sid)
  if (!data) {
    if (el) el.remove()
    return
  }
  if (!el) {
    el = document.createElement('script')
    el.type = 'application/ld+json'
    el.id = sid
    document.head.appendChild(el)
  }
  el.textContent = JSON.stringify(data)
}

/**
 * Route-level SEO: title, description, canonical, Open Graph, Twitter, JSON-LD.
 * @param {import('vue').Ref|import('vue').ComputedRef|object} source reactive or plain seo object
 */
export function useSeo(source) {
  const apply = (seo) => {
    if (typeof document === 'undefined' || !seo) return
    const title = seo.title ? `${seo.title} · ${SITE.name}` : `${SITE.name} · Engineering lab`
    const description = seo.description || SITE.description
    const path = seo.path || '/'
    const url = `${SITE.url}${path.startsWith('/') ? path : `/${path}`}`
    const image = seo.image || DEFAULT_OG

    document.title = title
    ensureMeta('meta[name="description"]', { name: 'description', content: description })
    ensureLink('canonical', url)

    ensureMeta('meta[property="og:type"]', { property: 'og:type', content: seo.type || 'website' })
    ensureMeta('meta[property="og:site_name"]', { property: 'og:site_name', content: SITE.legalName })
    ensureMeta('meta[property="og:title"]', { property: 'og:title', content: title })
    ensureMeta('meta[property="og:description"]', { property: 'og:description', content: description })
    ensureMeta('meta[property="og:url"]', { property: 'og:url', content: url })
    ensureMeta('meta[property="og:image"]', { property: 'og:image', content: image })

    ensureMeta('meta[name="twitter:card"]', { name: 'twitter:card', content: 'summary_large_image' })
    ensureMeta('meta[name="twitter:title"]', { name: 'twitter:title', content: title })
    ensureMeta('meta[name="twitter:description"]', { name: 'twitter:description', content: description })
    ensureMeta('meta[name="twitter:image"]', { name: 'twitter:image', content: image })

    setJsonLd('org', {
      '@context': 'https://schema.org',
      '@type': 'Organization',
      name: SITE.legalName,
      alternateName: SITE.name,
      url: SITE.url,
      logo: DEFAULT_OG,
      sameAs: SITE.sameAs
    })

    if (seo.person) {
      setJsonLd('person', {
        '@context': 'https://schema.org',
        '@type': 'Person',
        name: AUTHOR.name,
        alternateName: AUTHOR.alternateName,
        jobTitle: AUTHOR.jobTitle,
        description: AUTHOR.description,
        url: AUTHOR.url,
        image: AUTHOR.image,
        worksFor: { '@type': 'Organization', name: SITE.legalName, url: SITE.url }
      })
    } else {
      setJsonLd('person', null)
    }

    if (seo.jsonLd) {
      setJsonLd('page', seo.jsonLd)
    } else {
      setJsonLd('page', null)
    }

    if (seo.faq?.length) {
      setJsonLd('faq', {
        '@context': 'https://schema.org',
        '@type': 'FAQPage',
        mainEntity: seo.faq.map((f) => ({
          '@type': 'Question',
          name: f.question,
          acceptedAnswer: { '@type': 'Answer', text: f.answer }
        }))
      })
    } else {
      setJsonLd('faq', null)
    }
  }

  if (source && typeof source === 'object' && 'value' in source) {
    watch(source, (v) => apply(v), { immediate: true, deep: true })
  } else {
    apply(source)
  }

  onBeforeUnmount(() => {
    setJsonLd('page', null)
    setJsonLd('faq', null)
  })
}

export function faqToJsonLd(faq) {
  if (!faq?.length) return null
  return {
    '@context': 'https://schema.org',
    '@type': 'FAQPage',
    mainEntity: faq.map((f) => ({
      '@type': 'Question',
      name: f.question,
      acceptedAnswer: { '@type': 'Answer', text: f.answer }
    }))
  }
}
