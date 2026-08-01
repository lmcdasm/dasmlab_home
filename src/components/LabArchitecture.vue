<template>
  <div class="arch" role="img" aria-label="DASMLAB lab wiring: edge, OpenShift apps, GitOps, and Surfing media on R2">
    <svg viewBox="0 0 920 420" xmlns="http://www.w3.org/2000/svg" class="arch-svg">
      <defs>
        <marker id="dasmArr" markerWidth="7" markerHeight="7" refX="6" refY="3.5" orient="auto">
          <path d="M0,0 L7,3.5 L0,7 Z" fill="#5a7386" />
        </marker>
        <marker id="dasmArrTeal" markerWidth="7" markerHeight="7" refX="6" refY="3.5" orient="auto">
          <path d="M0,0 L7,3.5 L0,7 Z" fill="#1f6f62" />
        </marker>
        <marker id="dasmArrAccent" markerWidth="7" markerHeight="7" refX="6" refY="3.5" orient="auto">
          <path d="M0,0 L7,3.5 L0,7 Z" fill="#976eb0" />
        </marker>
        <pattern id="archGrid" width="20" height="20" patternUnits="userSpaceOnUse">
          <path d="M 20 0 L 0 0 0 20" fill="none" stroke="rgba(31, 111, 98, 0.06)" stroke-width="1" />
        </pattern>
        <linearGradient id="edgeGrad" x1="0%" y1="0%" x2="100%" y2="0%">
          <stop offset="0%" stop-color="#1f6f62" />
          <stop offset="100%" stop-color="#2f8f7d" />
        </linearGradient>
        <linearGradient id="r2Grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="#f48120" />
          <stop offset="100%" stop-color="#faad3f" />
        </linearGradient>
      </defs>

      <rect width="920" height="420" fill="url(#archGrid)" />

      <!-- Outer frame -->
      <rect x="18" y="16" width="884" height="388" rx="16" fill="#fbfcfe" stroke="#1d2b36" stroke-width="2.25" />
      <text x="38" y="42" class="t-host">DASMLAB LAB WIRING</text>
      <text x="250" y="42" class="t-muted">edge · OCP · GitOps · media off-disk</text>

      <!-- Users / edge -->
      <rect x="40" y="68" width="120" height="56" rx="10" fill="url(#edgeGrad)" />
      <text x="100" y="93" text-anchor="middle" class="t-on-dark">YOU</text>
      <text x="100" y="110" text-anchor="middle" class="t-on-dark-xs">browser</text>

      <path d="M160 96 H 210" stroke="#1f6f62" stroke-width="2" fill="none" marker-end="url(#dasmArrTeal)" />

      <rect x="210" y="62" width="150" height="68" rx="10" fill="#12202c" stroke="#2f8f7d" stroke-width="2" />
      <text x="285" y="88" text-anchor="middle" class="t-on-dark-sm">EDGE</text>
      <text x="285" y="108" text-anchor="middle" class="t-on-dark-xs">HAProxy · DNS · CF</text>

      <path d="M360 96 H 410" stroke="#1f6f62" stroke-width="2" fill="none" marker-end="url(#dasmArrTeal)" />

      <!-- OCP apps cluster box -->
      <rect x="410" y="54" width="340" height="150" rx="12" fill="#f4faf8" stroke="#2f8f7d" stroke-width="2" stroke-dasharray="7 4" />
      <text x="580" y="76" text-anchor="middle" class="t-label-teal">OPENSHIFT · 2026-PROD</text>

      <rect x="428" y="92" width="92" height="40" rx="8" fill="#1f6f62" />
      <text x="474" y="117" text-anchor="middle" class="t-on-dark-sm">home</text>

      <rect x="534" y="92" width="92" height="40" rx="8" fill="#00838f" />
      <text x="580" y="117" text-anchor="middle" class="t-on-dark-sm">surfing</text>

      <rect x="640" y="92" width="92" height="40" rx="8" fill="#1565c0" />
      <text x="686" y="117" text-anchor="middle" class="t-on-dark-sm">apps…</text>

      <rect x="428" y="146" width="148" height="40" rx="8" fill="#37474f" />
      <text x="502" y="171" text-anchor="middle" class="t-on-dark-sm">interview-me</text>

      <rect x="592" y="146" width="140" height="40" rx="8" fill="#455a64" />
      <text x="662" y="171" text-anchor="middle" class="t-on-dark-sm">mini-acm</text>

      <!-- GitOps column -->
      <rect x="770" y="54" width="112" height="150" rx="12" fill="#f7f3fb" stroke="#976eb0" stroke-width="2" />
      <text x="826" y="78" text-anchor="middle" class="t-label-accent">GITOPS</text>
      <rect x="786" y="92" width="80" height="32" rx="6" fill="#976eb0" />
      <text x="826" y="113" text-anchor="middle" class="t-on-dark-xs">Argo CD</text>
      <rect x="786" y="136" width="80" height="48" rx="6" fill="#5e4570" />
      <text x="826" y="156" text-anchor="middle" class="t-on-dark-xs">live-cicd</text>
      <text x="826" y="172" text-anchor="middle" class="t-on-dark-xs">GHCR</text>

      <path d="M750 129 H 770" stroke="#976eb0" stroke-width="1.75" fill="none" marker-end="url(#dasmArrAccent)" />
      <text x="720" y="122" class="t-label-accent">sync</text>

      <!-- Surfing → R2 path -->
      <path d="M580 204 V 240" stroke="#00838f" stroke-width="1.75" fill="none" marker-end="url(#dasmArrTeal)" />
      <text x="590" y="228" class="t-label-teal">bytes leave PVC</text>

      <rect x="410" y="248" width="220" height="72" rx="12" fill="#fff8f0" stroke="#f48120" stroke-width="2.25" />
      <text x="520" y="278" text-anchor="middle" class="t-label-orange">CLOUDFLARE R2</text>
      <text x="520" y="298" text-anchor="middle" class="t-muted-xs">media.surfing · zero egress</text>

      <rect x="660" y="248" width="220" height="72" rx="12" fill="#eef7f5" stroke="#1f6f62" stroke-width="2" />
      <text x="770" y="278" text-anchor="middle" class="t-label-teal">CDN EDGE</text>
      <text x="770" y="298" text-anchor="middle" class="t-muted-xs">cached until free-tier ceiling</text>

      <path d="M630 284 H 660" stroke="#f48120" stroke-width="1.75" fill="none" marker-end="url(#dasmArr)" />

      <!-- cheapcloud broker note -->
      <rect x="40" y="248" width="320" height="72" rx="12" fill="#f5f8fb" stroke="#5a7386" stroke-width="1.75" />
      <text x="56" y="278" class="t-label-grey">CHEAPCLOUD · MediaBroker</text>
      <text x="56" y="300" class="t-muted-xs">r2 | azure | hetzner | bunny — dry-run → live</text>

      <path d="M360 284 H 410" stroke="#5a7386" stroke-width="1.5" stroke-dasharray="4 3" fill="none" marker-end="url(#dasmArr)" />

      <!-- Floor: basement disk (deprecated path) -->
      <rect x="40" y="348" width="840" height="36" rx="8" fill="#eceff1" stroke="#78909c" stroke-width="1.5" />
      <text x="460" y="371" text-anchor="middle" class="t-label-grey">
        BASEMENT LVMS PVC — still origin today · migrate surfing media → R2
      </text>

      <!-- Legend -->
      <text x="40" y="232" class="t-muted-xs">teal = cluster path · orange = media CDN · purple = GitOps</text>
    </svg>
  </div>
</template>

<script setup>
</script>

<style scoped>
.arch {
  width: 100%;
  border: 1.5px solid rgba(29, 43, 54, 0.22);
  border-radius: 16px;
  background: #fff;
  overflow: hidden;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.8),
    0 12px 28px rgba(18, 40, 52, 0.08);
}

.arch-svg {
  display: block;
  width: 100%;
  height: auto;
}

.t-host {
  fill: #12202c;
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.06em;
}

.t-muted {
  fill: #7a8e9c;
  font-size: 11px;
}

.t-muted-xs {
  fill: #7a8e9c;
  font-size: 10px;
}

.t-on-dark {
  fill: #fff;
  font-size: 14px;
  font-weight: 800;
}

.t-on-dark-sm {
  fill: #fff;
  font-size: 12px;
  font-weight: 700;
}

.t-on-dark-xs {
  fill: #fff;
  font-size: 10px;
  font-weight: 650;
}

.t-label-teal {
  fill: #1f6f62;
  font-size: 11px;
  font-weight: 750;
  letter-spacing: 0.04em;
}

.t-label-accent {
  fill: #7a4f96;
  font-size: 11px;
  font-weight: 750;
  letter-spacing: 0.04em;
}

.t-label-orange {
  fill: #c45f0a;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.05em;
}

.t-label-grey {
  fill: #546e7a;
  font-size: 11px;
  font-weight: 650;
}
</style>
