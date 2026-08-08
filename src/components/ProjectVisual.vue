<template>
  <svg
    ref="svgEl"
    class="project-visual"
    :viewBox="`0 0 ${width} ${height}`"
    role="img"
    :aria-label="ariaLabel"
  />
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref, watch } from 'vue'
import * as d3 from 'd3'

const props = defineProps({
  problem: { type: String, required: true },
  outcome: { type: String, required: true },
  techs: { type: Array, default: () => [] },
  accent: { type: String, default: '#1f6f62' },
  width: { type: Number, default: 420 },
  height: { type: Number, default: 200 }
})

const svgEl = ref(null)
let ro = null

const ariaLabel = `Diagram: problem ${props.problem}; outcome ${props.outcome}; tech ${props.techs.join(', ')}`

function draw() {
  const el = svgEl.value
  if (!el) return
  const w = props.width
  const h = props.height
  const svg = d3.select(el)
  svg.selectAll('*').remove()

  const g = svg.append('g')

  // Atmosphere
  const defs = svg.append('defs')
  const grad = defs
    .append('linearGradient')
    .attr('id', `pv-grad-${props.accent.replace('#', '')}`)
    .attr('x1', '0%')
    .attr('y1', '0%')
    .attr('x2', '100%')
    .attr('y2', '100%')
  grad.append('stop').attr('offset', '0%').attr('stop-color', props.accent).attr('stop-opacity', 0.18)
  grad.append('stop').attr('offset', '100%').attr('stop-color', '#12202c').attr('stop-opacity', 0.06)

  g.append('rect')
    .attr('width', w)
    .attr('height', h)
    .attr('rx', 14)
    .attr('fill', `url(#pv-grad-${props.accent.replace('#', '')})`)

  // Flow nodes: Problem → Tech cluster → Outcome
  const nodes = [
    { id: 'problem', label: 'Problem', x: 70, y: h * 0.42, detail: props.problem },
    { id: 'tech', label: 'Stack', x: w * 0.5, y: h * 0.38, detail: props.techs.slice(0, 4).join(' · ') },
    { id: 'outcome', label: 'Outcome', x: w - 70, y: h * 0.42, detail: props.outcome }
  ]

  const linkGen = d3
    .linkHorizontal()
    .x((d) => d.x)
    .y((d) => d.y)

  const path = g
    .append('path')
    .attr(
      'd',
      linkGen({
        source: { x: nodes[0].x + 36, y: nodes[0].y },
        target: { x: nodes[1].x - 36, y: nodes[1].y }
      })
    )
    .attr('fill', 'none')
    .attr('stroke', props.accent)
    .attr('stroke-width', 2.2)
    .attr('stroke-opacity', 0.45)
    .attr('stroke-dasharray', '6 5')

  const path2 = g
    .append('path')
    .attr(
      'd',
      linkGen({
        source: { x: nodes[1].x + 36, y: nodes[1].y },
        target: { x: nodes[2].x - 36, y: nodes[2].y }
      })
    )
    .attr('fill', 'none')
    .attr('stroke', props.accent)
    .attr('stroke-width', 2.2)
    .attr('stroke-opacity', 0.45)
    .attr('stroke-dasharray', '6 5')

  // Animate dash
  const len = 120
  path
    .attr('stroke-dashoffset', len)
    .transition()
    .duration(1400)
    .ease(d3.easeLinear)
    .attr('stroke-dashoffset', 0)
  path2
    .attr('stroke-dashoffset', len)
    .transition()
    .delay(200)
    .duration(1400)
    .ease(d3.easeLinear)
    .attr('stroke-dashoffset', 0)

  const node = g
    .selectAll('.pv-node')
    .data(nodes)
    .join('g')
    .attr('class', 'pv-node')
    .attr('transform', (d) => `translate(${d.x},${d.y})`)

  node
    .append('circle')
    .attr('r', 0)
    .attr('fill', '#fff')
    .attr('stroke', props.accent)
    .attr('stroke-width', 2.5)
    .transition()
    .duration(500)
    .attr('r', 28)

  node
    .append('text')
    .attr('text-anchor', 'middle')
    .attr('dy', '0.35em')
    .attr('fill', '#12202c')
    .attr('font-size', 10)
    .attr('font-weight', 700)
    .attr('font-family', 'Fraunces, Georgia, serif')
    .text((d) => d.label)

  node
    .append('text')
    .attr('text-anchor', 'middle')
    .attr('y', 48)
    .attr('fill', '#3d5263')
    .attr('font-size', 9.5)
    .attr('font-family', 'Source Sans 3, system-ui, sans-serif')
    .each(function (d) {
      const words = String(d.detail || '')
        .split(/\s+/)
        .filter(Boolean)
      const lines = []
      let line = ''
      words.forEach((word) => {
        const next = line ? `${line} ${word}` : word
        if (next.length > 18) {
          if (line) lines.push(line)
          line = word
        } else {
          line = next
        }
      })
      if (line) lines.push(line)
      d3.select(this)
        .selectAll('tspan')
        .data(lines.slice(0, 3))
        .join('tspan')
        .attr('x', 0)
        .attr('dy', (__, i) => (i === 0 ? 0 : 11))
        .text((t) => t)
    })

  // Orbiting tech chips
  const chips = props.techs.slice(0, 5)
  const orbit = g.append('g').attr('transform', `translate(${w * 0.5},${h * 0.38})`)
  chips.forEach((tech, i) => {
    const angle = (Math.PI * 2 * i) / Math.max(chips.length, 1) - Math.PI / 2
    const r = 58
    const x = Math.cos(angle) * r
    const y = Math.sin(angle) * r * 0.55
    const chip = orbit.append('g').attr('transform', `translate(${x},${y})`).attr('opacity', 0)
    chip
      .append('rect')
      .attr('x', -28)
      .attr('y', -10)
      .attr('width', 56)
      .attr('height', 20)
      .attr('rx', 8)
      .attr('fill', '#fff')
      .attr('stroke', props.accent)
      .attr('stroke-opacity', 0.35)
    chip
      .append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', '0.35em')
      .attr('font-size', 8.5)
      .attr('font-weight', 650)
      .attr('fill', props.accent)
      .attr('font-family', 'Source Sans 3, system-ui, sans-serif')
      .text(String(tech).slice(0, 10))
    chip.transition().delay(400 + i * 90).duration(400).attr('opacity', 1)
  })
}

onMounted(() => {
  draw()
  if (typeof ResizeObserver !== 'undefined' && svgEl.value?.parentElement) {
    ro = new ResizeObserver(() => draw())
    ro.observe(svgEl.value.parentElement)
  }
})

onBeforeUnmount(() => {
  if (ro) ro.disconnect()
})

watch(
  () => [props.problem, props.outcome, props.techs, props.accent],
  () => draw(),
  { deep: true }
)
</script>

<style scoped>
.project-visual {
  width: 100%;
  height: auto;
  display: block;
  min-height: 160px;
}
</style>
