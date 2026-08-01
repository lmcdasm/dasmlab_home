#!/usr/bin/env bash
# Publish one Surfing album (day) from PVC staging → R2 CDN.
# Usage:
#   ./scripts/publish-album.sh                 # publishes first/only day
#   ./scripts/publish-album.sh <dayId>
#   CLEANUP_PVC=1 ./scripts/publish-album.sh   # also delete local bytes after put
set -euo pipefail

NAMESPACE="${SURFING_NAMESPACE:-surfing-service-system}"
DAY_ID="${1:-}"
CLEANUP="${CLEANUP_PVC:-0}"
QS=""
if [[ "${CLEANUP}" == "1" || "${CLEANUP}" == "true" ]]; then
  QS="?cleanup_pvc=true"
fi

POD="$(oc -n "${NAMESPACE}" get pods -l app=surfing-service -o jsonpath='{.items[0].metadata.name}')"
echo "Using pod ${POD}"

if [[ -z "${DAY_ID}" ]]; then
  DAY_ID="$(oc -n "${NAMESPACE}" exec "${POD}" -- \
    curl -fsS http://127.0.0.1:10023/days | python3 -c 'import sys,json; d=json.load(sys.stdin); print(d[0]["id"] if d else "")')"
  echo "Auto-selected day ${DAY_ID}"
fi
[[ -n "${DAY_ID}" ]] || { echo "No day id"; exit 1; }

echo "Publishing day ${DAY_ID} cleanup_pvc=${CLEANUP}…"
oc -n "${NAMESPACE}" exec "${POD}" -- \
  curl -fsS -X POST "http://127.0.0.1:10023/days/${DAY_ID}/publish${QS}" | python3 -m json.tool

echo "Done. Verify a media URL hits pub-….r2.dev (or custom domain)."
