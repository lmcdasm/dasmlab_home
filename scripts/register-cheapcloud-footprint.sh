#!/usr/bin/env bash
# Push dasmlab-home live cost footprint to cheapcloud.
# Usage:
#   CHEAPCLOUD_URL=https://cheapcloud-dasmlab.apps.2026-prod-1.ocp.dasmlab.org \
#     ./scripts/register-cheapcloud-footprint.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
URL="${CHEAPCLOUD_URL:-https://cheapcloud-dasmlab.apps.2026-prod-1.ocp.dasmlab.org}"
BODY="${ROOT}/scripts/cheapcloud-footprint-dasmlab-home.json"

echo "PUT ${URL}/api/v1/assets/dasmlab-home"
curl -fsS -X PUT "${URL}/api/v1/assets/dasmlab-home" \
  -H 'content-type: application/json' \
  --data-binary @"${BODY}" | python3 -m json.tool

echo
echo "Media envelope:"
curl -fsS "${URL}/api/v1/media/envelope" | python3 -m json.tool | head -40 || echo "(media routes require cheapcloud image with footprint build)"
