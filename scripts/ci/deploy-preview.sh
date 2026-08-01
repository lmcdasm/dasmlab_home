#!/usr/bin/env bash
# Publish a per-developer dasmlab-home preview via GitOps (dasmlab-live-cicd).
# Does NOT oc-apply from the runner — Argo CD syncs the previews/ path.
set -euo pipefail

VERSION_TAG="${VERSION_TAG:?}"
ACTOR="${PREVIEW_ACTOR:?}"
CLUSTER_APPS="${CLUSTER_APPS_DOMAIN:-apps.2026-prod-1.ocp.dasmlab.org}"

OWNER="$(echo "${ACTOR}" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g; s/^-+//; s/-+$//; s/-+/-/g' | cut -c1-20)"
OWNER="${OWNER:-dev}"
NS="dasmlab-home-dev-${OWNER}"
HOST="dev-${OWNER}-dasmlab-home.${CLUSTER_APPS}"
PREVIEW_URL="https://${HOST}"

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
RENDERED="$(mktemp)"
sed \
  -e "s|__VERSION__|${VERSION_TAG}|g" \
  -e "s|__PREVIEW_NS__|${NS}|g" \
  -e "s|__PREVIEW_HOST__|${HOST}|g" \
  -e "s|__PREVIEW_OWNER__|${OWNER}|g" \
  "${ROOT}/k8s_envelope/dasmlab-home_preview-ocp.yaml" > "${RENDERED}"

echo "Preview owner=${OWNER} ns=${NS} host=${HOST} version=${VERSION_TAG}"

if [[ "${SKIP_PREVIEW_BOOTSTRAP:-}" == "true" ]]; then
  echo "SKIP_PREVIEW_BOOTSTRAP=true — not copying preview secrets"
elif command -v oc >/dev/null 2>&1; then
  if oc whoami >/dev/null 2>&1; then
    echo "oc available — ensuring preview NS secrets via bootstrap-preview-ns.sh"
    bash "${ROOT}/scripts/ci/bootstrap-preview-ns.sh" "${OWNER}" || {
      echo "WARN: bootstrap-preview-ns.sh failed; pod may stay Pending until secrets exist" >&2
    }
  else
    echo "WARN: oc present but not logged in — skip preview secret bootstrap" >&2
  fi
else
  echo "WARN: oc not on PATH — skip preview secret bootstrap" >&2
fi

DEPLOY_TOKEN=""
if [ -f "/home/dasm/gh_token" ]; then
  DEPLOY_TOKEN="$(tr -d '\n\r' < /home/dasm/gh_token)"
fi
if [ -z "${DEPLOY_TOKEN}" ]; then
  DEPLOY_TOKEN="${DASMLAB_GHCR_PAT:-${GH_TOKEN:-}}"
fi
if [ -z "${DEPLOY_TOKEN}" ]; then
  echo "ERROR: deploy token not set (gh_token / DASMLAB_GHCR_PAT / GH_TOKEN)" >&2
  exit 1
fi

WORK="$(mktemp -d)"
trap 'rm -rf "${WORK}" "${RENDERED}"' EXIT
git clone --depth 1 "https://x-access-token:${DEPLOY_TOKEN}@github.com/lmcdasm/dasmlab-live-cicd.git" "${WORK}/live-cicd"
PREVIEW_DIR="${WORK}/live-cicd/clusters/2026-prod-1/dasmlab-home/previews"
mkdir -p "${PREVIEW_DIR}"
cp "${RENDERED}" "${PREVIEW_DIR}/${OWNER}.yaml"

cd "${WORK}/live-cicd"
git config user.name "dasmlab-bot"
git config user.email "ci@dasmlab.org"
git add "clusters/2026-prod-1/dasmlab-home/previews/${OWNER}.yaml"
if git diff --cached --quiet; then
  echo "No GitOps preview changes (manifest identical)"
else
  git commit -m "preview(${OWNER}): dasmlab-home ${VERSION_TAG}"
  git push
fi

{
  echo "## dasmlab-home preview (GitOps)"
  echo ""
  echo "- **URL:** ${PREVIEW_URL}"
  echo "- **Namespace:** \`${NS}\`"
  echo "- **GitOps file:** \`clusters/2026-prod-1/dasmlab-home/previews/${OWNER}.yaml\`"
  echo "- **Argo app:** \`dasmlab-home-previews\` (auto-sync)"
  echo "- **Image:** \`ghcr.io/lmcdasm/dasmlab-home:${VERSION_TAG}\`"
  echo "- **Owner:** \`${OWNER}\`"
} >> "${GITHUB_STEP_SUMMARY:-/dev/null}"

echo "PREVIEW_URL=${PREVIEW_URL}"
echo "Preview GitOps published: ${PREVIEW_URL}"
