#!/usr/bin/env bash
# Chunked PVC/draft → R2 promote from inside the surfing-service pod (loopback, no OIDC).
# Safe: no cleanup_pvc. Persists every few items server-side.
#
# Usage:
#   ./scripts/promote-album-chunked.sh 5bd87267-225f-4f46-bf71-e6ff8a76463c
#   LIMIT=10 ./scripts/promote-album-chunked.sh <dayId>
set -euo pipefail

NAMESPACE="${SURFING_NAMESPACE:-surfing-service-system}"
DAY_ID="${1:?day id required}"
LIMIT="${LIMIT:-15}"
MAX_ROUNDS="${MAX_ROUNDS:-50}"

POD="$(oc -n "${NAMESPACE}" get pods -l app=surfing-service -o jsonpath='{.items[0].metadata.name}')"
echo "Pod=${POD} day=${DAY_ID} limit=${LIMIT}"

for round in $(seq 1 "${MAX_ROUNDS}"); do
  echo "── round ${round} ──"
  OUT="$(oc -n "${NAMESPACE}" exec "${POD}" -- \
    curl -sS -m 600 -X POST "http://127.0.0.1:10023/days/${DAY_ID}/publish?limit=${LIMIT}" || true)"
  if [[ -z "${OUT}" ]]; then
    echo "empty/failed response — retrying shortly"
    sleep 5
    continue
  fi
  echo "${OUT}" | python3 -c '
import sys, json
r=json.load(sys.stdin)
print("published=%s skipped=%s failed=%s remaining=%s" % (
    r.get("published"), r.get("skipped"), r.get("failed"), r.get("remaining")))
errs=r.get("errors") or []
for e in errs[:8]:
    print("  err:", e)
if len(errs) > 8:
    print("  … +%d more" % (len(errs)-8))
open("/tmp/surf-promote-remaining","w").write(str(r.get("remaining", -1)))
'
  REM="$(cat /tmp/surf-promote-remaining 2>/dev/null || echo -1)"
  if [[ "${REM}" == "0" ]]; then
    echo "Album fully on CDN."
    exit 0
  fi
  sleep 1
done

echo "Stopped after ${MAX_ROUNDS} rounds (remaining may still be >0)."
exit 1
