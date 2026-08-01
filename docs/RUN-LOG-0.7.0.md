# 0.7.0 run log — fixes / improvements as we execute PLAN-0.7.0-5x10

Rolling log of **≥10 improvements/fixes per task** while shipping the preview release.

---

## Task 1 — Preview scaffolding

| # | Fix / improvement |
|---|-------------------|
| 1 | Add `k8s_envelope/dasmlab-home_preview-ocp.yaml` (NS + SA + anyuid + Deploy/Svc/Route) |
| 2 | Host pattern `dev-__OWNER__-dasmlab-home.apps.2026-prod-1.ocp.dasmlab.org` |
| 3 | Port 8888 + nginx emptyDirs (match prod; avoid CrashLoop) |
| 4 | `scripts/ci/deploy-preview.sh` → live-cicd `dasmlab-home/previews/${OWNER}.yaml` |
| 5 | `scripts/ci/bootstrap-preview-ns.sh` copies `dasmlab-ghcr-pull` |
| 6 | `scripts/ci/cleanup-preview.sh` for PR-close prune |
| 7 | Argo Application `dasmlab-home-previews` (prune + selfHeal) |
| 8 | live-cicd `previews/README.md` |
| 9 | Image `ghcr.io/lmcdasm/dasmlab-home` (not dasmlab/ org mismatch) |
| 10 | Document preview URL in TARGET + this run log |

---

## Task 2 — CI

| # | Fix / improvement |
|---|-------------------|
| 1 | Split prepare / build / preview-deploy / prod-gitops |
| 2 | `is_prod` only on `main` |
| 3 | Trigger on `pull_request` + feature branches |
| 4 | Concurrency `preview-{actor}` |
| 5 | Wire `deploy-preview.sh` |
| 6 | Keep Buildah/Kaniko + counter args |
| 7 | Prod dual-cluster path untouched for main |
| 8 | Gate auto-merge until safer |
| 9 | Step summary with preview URL |
| 10 | Branch name includes lab-constellation |

---

## Task 3 — Ship + PR

| # | Fix / improvement |
|---|-------------------|
| 1 | Branch `2026-lab-constellation` |
| 2 | Commit UI (map/arch/WhatsNew out) |
| 3 | Commit preview scaffolding |
| 4 | Exclude `in_progress:/` |
| 5 | `gh pr create` |
| 6 | PR body with preview URL |
| 7 | Push live-cicd Argo app |
| 8 | Trigger / wait CI |
| 9 | Fix build failures if any |
| 10 | Record tag + URLs |

---

## Task 4 — Cluster heal

| # | Fix / improvement |
|---|-------------------|
| 1 | Apply / sync Argo previews app |
| 2 | Bootstrap pull secret |
| 3 | Wait pods Ready |
| 4 | Fix ImagePull / SCC as needed |
| 5 | Confirm Route |
| 6 | Curl HTTPS 200 |
| 7 | Smoke Lab map |
| 8 | Smoke architecture |
| 9 | Note surfing proxy behavior on preview |
| 10 | Publish working URL |

---

## Task 5 — Harden

| # | Fix / improvement |
|---|-------------------|
| 1 | Extra contrast on LabMap spokes |
| 2 | Arch diagram label accuracy |
| 3 | Mobile stack check |
| 4 | WhatsNew prune PR follow-up |
| 5 | preview-cleanup workflow |
| 6 | R2 spike note |
| 7 | cheapcloud broker note |
| 8 | Promote checklist on PR |
| 9 | TARGET checkboxes |
| 10 | Next-steal list from siblings |

---

## Live notes (filled while running)

_Start: 2026-08-01_
