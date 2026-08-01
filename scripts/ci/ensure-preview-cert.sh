#!/usr/bin/env bash
# Ensure HAProxy on 10.20.1.10 has a Let's Encrypt CERT entry for FQDN.
# Hands-free pattern from mini-acm scripts/ci/ensure-prod-cert.sh.
# Idempotent: if FQDN already listed in runme.sh, do nothing.
set -euo pipefail

FQDN="${1:?usage: ensure-preview-cert.sh <fqdn>}"
PROXY_HOST="${PREVIEW_PROXY_HOST:-10.20.1.10}"
PROXY_USER="${PREVIEW_PROXY_USER:-dasm}"
PROXY_DIR="${PREVIEW_PROXY_DIR:-/home/dasm/dasmlab-internal/new_haproxy}"

if ! command -v ssh >/dev/null 2>&1; then
  echo "ERROR: ssh not on PATH — cannot update HAProxy CERT for ${FQDN}" >&2
  exit 1
fi

ssh -o BatchMode=yes -o StrictHostKeyChecking=accept-new \
  "${PROXY_USER}@${PROXY_HOST}" bash -s -- "${FQDN}" "${PROXY_DIR}" <<'EOS'
set -euo pipefail
FQDN="$1"
DIR="$2"
cd "$DIR"
if grep -Fq "=${FQDN}" runme.sh; then
  echo "HAProxy CERT already present for ${FQDN}"
  exit 0
fi

last="$(grep -oE 'CERT[0-9]+=' runme.sh | grep -oE '[0-9]+' | sort -n | tail -1)"
next=$((last + 1))
echo "Adding CERT${next}=${FQDN} to runme.sh and recreating new-haproxy"

tmp="$(mktemp)"
awk -v n="$next" -v h="$FQDN" '
  / -e EMAIL=/ && !done {
    print "    -e CERT" n "=" h " \\"
    done=1
  }
  { print }
' runme.sh > "$tmp"
mv "$tmp" runme.sh
chmod +x runme.sh

./runme.sh
echo "HAProxy updated for ${FQDN} (CERT${next})"
EOS
