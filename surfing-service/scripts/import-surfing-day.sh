#!/usr/bin/env bash
# Bulk-import a local media folder into the surfing-service PVC preload area.
# On pod restart the service imports preload/* into persistent storage (idempotent).
#
# Usage: ./import-surfing-day.sh "BONAIRE-J-2026" /home/dasm/BONAIRE-MEDIA [date] [location]
set -euo pipefail

DAY_TITLE="${1:?day title required}"
SOURCE_DIR="${2:?source directory required}"
DAY_DATE="${3:-2026-01-01}"
DAY_LOCATION="${4:-Bonaire}"
NAMESPACE="${SURFING_NAMESPACE:-surfing-service-system}"
KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}"

if [[ ! -d "$SOURCE_DIR" ]]; then
  echo "ERROR: source directory not found: $SOURCE_DIR" >&2
  exit 1
fi

POD="$(oc --kubeconfig="$KUBECONFIG" get pods -n "$NAMESPACE" -l app=surfing-service -o jsonpath='{.items[0].metadata.name}')"
if [[ -z "$POD" ]]; then
  echo "ERROR: no surfing-service pod in $NAMESPACE" >&2
  exit 1
fi

PRELOAD_PATH="/data/preload/${DAY_TITLE}"
echo "Importing into pod $POD at ${PRELOAD_PATH} ..."

oc --kubeconfig="$KUBECONFIG" exec -n "$NAMESPACE" "$POD" -- mkdir -p "$PRELOAD_PATH"

echo "Syncing media files (this may take a while for large folders) ..."
oc --kubeconfig="$KUBECONFIG" rsync "$SOURCE_DIR/" "$NAMESPACE/$POD:${PRELOAD_PATH}/"

TMP_META="$(mktemp)"
cat >"$TMP_META" <<EOF
{
  "title": "$DAY_TITLE",
  "date": "$DAY_DATE",
  "location": "$DAY_LOCATION"
}
EOF
oc --kubeconfig="$KUBECONFIG" cp "$TMP_META" "$NAMESPACE/$POD:${PRELOAD_PATH}/day.json"
rm -f "$TMP_META"

echo "Restarting surfing-service to run preload import ..."
oc --kubeconfig="$KUBECONFIG" rollout restart deployment/surfing-service -n "$NAMESPACE"
oc --kubeconfig="$KUBECONFIG" rollout status deployment/surfing-service -n "$NAMESPACE" --timeout=300s

echo "Done. Day '${DAY_TITLE}' should appear on /surfing after the pod restarts."
