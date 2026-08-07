# 0.7.0 release plan — 5 tasks × 10 iterations

**Goal:** Running preview release at `https://dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org` (actor-scoped), then promote via PR.  
**Assumptions:** R2 Phase 1 later; WhatsNew FE done; live-cicd prune PR open; steal interview-me / product-shots preview pattern.  
**Rule:** Each task runs **10 iterations** (attempt → check → fix → next). Do not stop early on “looks fine” — burn the iterations or reinvest surplus in polish.

---

## Task 1 — Preview GitOps scaffolding (home)

Steal `interview-me` / `product-shots` preview shape for `dasmlab-home`.

| Iter | Action | Done when |
|------|--------|-----------|
| 1 | Add `k8s_envelope/dasmlab-home_preview-ocp.yaml` (NS, SA, anyuid, Deploy, Svc, Route; host `__PREVIEW_HOST__`) | File exists, placeholders match siblings |
| 2 | Add `scripts/ci/deploy-preview.sh` → GitOps `clusters/2026-prod-1/dasmlab-home/previews/${OWNER}.yaml` | Script dry-runs sed + path |
| 3 | Add `scripts/ci/bootstrap-preview-ns.sh` (copy `dasmlab-ghcr-pull` from `dasmlab-home-system`) | Script matches product-shots pattern |
| 4 | Add `scripts/ci/cleanup-preview.sh` (+ optional workflow later) | Removes owner file from previews/ |
| 5 | Add live-cicd `previews/README.md` + Argo app `dasmlab-home-previews.yaml` | Manifest valid YAML |
| 6 | PR/merge Argo Application into `dasmlab-live-cicd` (or apply Application CR) | App exists in openshift-gitops |
| 7 | Local render test: `OWNER=lmcdasm` → expect host `dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org` | Host string exact |
| 8 | Fix image name (`ghcr.io/lmcdasm/dasmlab-home`) + port `8888` + nginx emptyDirs | Matches prod envelope |
| 9 | Confirm prune/selfHeal + `include: "*.yaml"` on previews path | Same as interview-me-previews |
| 10 | Document preview URL + NS in `docs/TARGET-0.7.0.md` | Docs updated |

**Exit:** Scaffold ready; no cluster pod yet required.

---

## Task 2 — CI: branch/PR → build → preview (not prod)

Wire `.github/workflows/main.yml` so non-`main` builds preview; `main` keeps dual-cluster prod.

| Iter | Action | Done when |
|------|--------|-----------|
| 1 | Split jobs: `prepare` / `build-publish` / `preview-deploy` / `prod-gitops` | Workflow parses |
| 2 | `is_prod` only on `main` (or explicit dispatch) | Branch builds ≠ prod live/ |
| 3 | Trigger on `pull_request` + push to feature branches (incl. `2026-lab-constellation` / current FE branch) | `on:` block updated |
| 4 | Concurrency group `dasmlab-home-preview-{actor}` | No clobber across actors |
| 5 | Preview job calls `deploy-preview.sh` with `PREVIEW_ACTOR` | Job present |
| 6 | Keep Buildah/Kaniko + counter build-args | Image builds |
| 7 | Prod path unchanged for `main` → `2026-prod-1` + `2026-prod-2-1` live/ | No regression |
| 8 | Disable or gate auto-merge-to-main until preview green | Safer than old FE auto-merge |
| 9 | Step summary prints preview URL | Visible in Actions |
| 10 | Push workflow; confirm Actions run starts on self-hosted | Run ID logged |

**Exit:** CI can publish a preview manifest without touching prod `live/`.

---

## Task 3 — Ship code + open PR (running release candidate)

Land Lab map + architecture + WhatsNew removal on a release branch and open PR.

| Iter | Action | Done when |
|------|--------|-----------|
| 1 | Branch `2026-lab-constellation` (or continue FE branch) | Branch pushed |
| 2 | Commit UI + docs (exclude `in_progress:/`, node junk) | Clean commit |
| 3 | Commit preview scripts + envelope | Second commit or same |
| 4 | `gh pr create` → main | PR URL |
| 5 | Comment preview URL convention on PR | Comment posted |
| 6 | Ensure CI triggered for PR/push | Workflow queued |
| 7 | Watch build; fix Dockerfile/nginx if fail | Image in GHCR |
| 8 | Watch preview-deploy; fix token/GitOps push | Manifest in live-cicd previews/ |
| 9 | Confirm Argo sync of `dasmlab-home-previews` | Synced |
| 10 | Record image tag + PR + preview host | Written in chat + TARGET |

**Exit:** PR open; image tagged; GitOps preview file present.

---

## Task 4 — Cluster heal → working preview URL

Get pods Ready and HTTPS responding.

| Iter | Action | Done when |
|------|--------|-----------|
| 1 | `oc` context to `2026-prod-1` | `oc whoami` ok |
| 2 | Ensure Argo Application applied if not synced from git | App listed |
| 3 | Bootstrap NS secrets (`dasmlab-ghcr-pull`) | Secret in `dasmlab-home-dev-lmcdasm` |
| 4 | Check Deploy/pods; fix ImagePullBackOff | Pod Running |
| 5 | Fix SCC/anyuid if CrashLoop (nginx root) | Ready |
| 6 | Route exists for `dev-lmcdasm-dasmlab-home.apps…` | `oc get route` |
| 7 | Curl HTTPS (follow redirects); fix TLS/edge if needed | HTTP 200 |
| 8 | Browser/snapshot: Lab map + architecture visible | UI smoke pass |
| 9 | Surfing proxy via home nginx optional for preview (or stub) | Documented if broken |
| 10 | Paste working URL + any follow-ups | User-facing |

**Exit:** **Running release** = preview URL loads the new homepage.

---

## Task 5 — Parallel harden + next-loop polish (always improve)

Burn remaining iterations on quality + R2 kickoff while preview stays up.

| Iter | Action | Done when |
|------|--------|-----------|
| 1 | Contrast/edge pass on LabMap (steal latest interview-me deltas) | Diff shipped |
| 2 | Architecture diagram accuracy vs real routes | Labels match cluster |
| 3 | Mobile: map stacks, CTAs usable | Smoke on narrow width |
| 4 | Merge WhatsNew live-cicd prune PR if green; delete Application CR | Pruned |
| 5 | Preview cleanup on PR close (copy interview-me `preview-cleanup.yml`) | Workflow exists |
| 6 | Spike surfing → R2 upload path (presign design only or stub) | Doc or code stub |
| 7 | cheapcloud MediaBroker interface sketch | Interface file or issue |
| 8 | Prod promote decision: merge PR after preview sign-off | Checklist in PR |
| 9 | Tag notes for 0.7.0-rc | TARGET checkboxes |
| 10 | Retrospective: what to steal next from sibling repos | 5 bullets in TARGET |

**Exit:** Preview stable; prune done; R2/cheapcloud next work queued; promote path clear.

---

## Parallelism

```
Task1 ──┬── Task2 ──→ Task3 ──→ Task4
        │                         │
        └─────────────────────────┴──→ Task5 (overlaps 3–4)
```

Tasks **1+2** start together. **3** needs 1–2. **4** needs 3. **5** starts as soon as preview URL exists.

---

## Definition of “running release” (this plan)

Not necessarily prod `dasmlab.org` yet — **developer preview** Live:

`https://dev-<owner>-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org`

For you: **`dev-lmcdasm-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org`**.

Prod cutover = Task 5 iter 8 after preview sign-off.
