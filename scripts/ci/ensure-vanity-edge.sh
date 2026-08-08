#!/usr/bin/env bash
# Ensure vanity product hosts are matched on the edge HAProxy and have CERTs.
# Host-first: edits live files on 10.20.1.10 only (never overwrite from moldy git).
#
# Products (vanity → OCP route Host rewrite):
#   mock-me.dasmlab.org        → mock-me.apps.2026-prod-1.ocp.dasmlab.org
#   interview-me.dasmlab.org   → interview-me.apps.2026-prod-1.ocp.dasmlab.org
#   camera-scrape.dasmlab.org  → camera-scrape.apps.2026-prod-1.ocp.dasmlab.org
#
# Usage: ./scripts/ci/ensure-vanity-edge.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
PROXY_HOST="${PREVIEW_PROXY_HOST:-10.20.1.10}"
PROXY_USER="${PREVIEW_PROXY_USER:-dasm}"
PROXY_DIR="${PREVIEW_PROXY_DIR:-/home/dasm/dasmlab-internal/new_haproxy}"

if ! command -v ssh >/dev/null 2>&1; then
  echo "ERROR: ssh not on PATH — cannot update HAProxy vanity edge" >&2
  exit 1
fi

echo "Ensuring vanity host match + backends on ${PROXY_USER}@${PROXY_HOST}:${PROXY_DIR}"

ssh -o BatchMode=yes -o StrictHostKeyChecking=accept-new \
  "${PROXY_USER}@${PROXY_HOST}" bash -s -- "${PROXY_DIR}" <<'EOS'
set -euo pipefail
DIR="$1"
CFG="${DIR}/haproxycfg/haproxy.cfg"

if [[ ! -f "$CFG" ]]; then
  echo "ERROR: missing ${CFG}" >&2
  exit 1
fi

if grep -Fq "acl host_mock_me hdr(host)" "$CFG"; then
  echo "Vanity host ACLs already present in haproxy.cfg"
else
  echo "Inserting vanity host ACLs + backends into haproxy.cfg"
  cp -a "$CFG" "${CFG}.bak-vanity-$(date +%Y%m%d%H%M%S)"
  python3 - "$CFG" <<'PY'
import sys
from pathlib import Path
cfg_path = Path(sys.argv[1])
text = cfg_path.read_text()

acl = """
    # BEGIN dasmlab vanity product hosts (managed by ensure-vanity-edge.sh)
    acl host_mock_me hdr(host) -i mock-me.dasmlab.org
    use_backend mock_me_backend if host_mock_me

    acl host_interview_me hdr(host) -i interview-me.dasmlab.org
    use_backend interview_me_backend if host_interview_me

    acl host_camera_scrape hdr(host) -i camera-scrape.dasmlab.org
    use_backend camera_scrape_backend if host_camera_scrape
    # END dasmlab vanity product hosts ACL
"""

backends = """
# BEGIN dasmlab vanity product hosts (managed by ensure-vanity-edge.sh)
backend mock_me_backend
    http-request set-log-level info
    http-request set-var(txn.custom_message) str("MOCK-ME VANITY")
    option httpclose
    option forwardfor except 127.0.0.1
    http-request set-header Host mock-me.apps.2026-prod-1.ocp.dasmlab.org
    http-request set-header X-Forwarded-For %[src]
    http-request set-header X-Forwarded-Proto https if { ssl_fc }
    http-request set-header X-Forwarded-Host %[req.hdr(Host)]
    server ocp2026prod1-mock-me 10.20.1.70:443 ssl verify none check sni str(mock-me.apps.2026-prod-1.ocp.dasmlab.org)

backend interview_me_backend
    http-request set-log-level info
    http-request set-var(txn.custom_message) str("INTERVIEW-ME VANITY")
    option httpclose
    option forwardfor except 127.0.0.1
    http-request set-header Host interview-me.apps.2026-prod-1.ocp.dasmlab.org
    http-request set-header X-Forwarded-For %[src]
    http-request set-header X-Forwarded-Proto https if { ssl_fc }
    http-request set-header X-Forwarded-Host %[req.hdr(Host)]
    server ocp2026prod1-interview-me 10.20.1.70:443 ssl verify none check sni str(interview-me.apps.2026-prod-1.ocp.dasmlab.org)

backend camera_scrape_backend
    http-request set-log-level info
    http-request set-var(txn.custom_message) str("CAMERA-SCRAPE VANITY")
    option httpclose
    option forwardfor except 127.0.0.1
    http-request set-header Host camera-scrape.apps.2026-prod-1.ocp.dasmlab.org
    http-request set-header X-Forwarded-For %[src]
    http-request set-header X-Forwarded-Proto https if { ssl_fc }
    http-request set-header X-Forwarded-Host %[req.hdr(Host)]
    server ocp2026prod1-camera-scrape 10.20.1.70:443 ssl verify none check sni str(camera-scrape.apps.2026-prod-1.ocp.dasmlab.org)
# END dasmlab vanity product hosts

"""

needle = "    # Fallback default backend (dasmlab-home)"
if needle not in text:
    raise SystemExit("Could not find default_backend insertion point")
# Avoid matching use_backend fsapi_backend — only the real backend definition.
be_needle = "\nbackend fsapi_backend"
if be_needle not in text:
    raise SystemExit("Could not find backend fsapi_backend definition")
if "backend mock_me_backend" in text.split("frontend https")[-1].split("backend fsapi_backend")[0] if False else False:
    pass
text = text.replace(needle, acl + "\n" + needle, 1)
text = text.replace(be_needle, "\n" + backends + "backend fsapi_backend", 1)
if "use_\n" in text or "use_#" in text:
    raise SystemExit("Patch corrupted use_backend line — abort")
if "use_backend fsapi_backend if fsapi_host" not in text:
    raise SystemExit("fsapi use_backend missing after patch — abort")
cfg_path.write_text(text)
print("haproxy.cfg patched")
PY
fi
EOS

need_runme=0
for fqdn in mock-me.dasmlab.org interview-me.dasmlab.org camera-scrape.dasmlab.org; do
  echo "--- CERT ensure: ${fqdn} ---"
  if ! getent hosts "$fqdn" >/dev/null 2>&1; then
    echo "SKIP CERT ${fqdn} — no DNS yet (add A/AAAA to 209.15.95.244, then re-run)"
    continue
  fi
  if ssh -o BatchMode=yes "${PROXY_USER}@${PROXY_HOST}" "grep -Fq '=${fqdn}' '${PROXY_DIR}/runme.sh'"; then
    echo "HAProxy CERT already present for ${fqdn}"
  else
    need_runme=1
    bash "${ROOT}/scripts/ci/ensure-preview-cert.sh" "${fqdn}"
  fi
done

# Soft reload when CERTs already present (avoid full certbot walk when possible)
if [[ "$need_runme" -eq 0 ]]; then
  echo "Reloading HAProxy (SIGHUP) to pick up cfg if needed"
  ssh -o BatchMode=yes -o StrictHostKeyChecking=accept-new \
    "${PROXY_USER}@${PROXY_HOST}" bash -s -- "${PROXY_DIR}" <<'EOS'
set -euo pipefail
DIR="$1"
cd "$DIR"
if docker ps --format '{{.Names}}' 2>/dev/null | grep -qx new-haproxy; then
  # Only HUP if haproxy master is already running (not still in certbot)
  if docker exec new-haproxy pgrep -x haproxy >/dev/null 2>&1; then
    docker kill -s HUP new-haproxy
    echo "SIGHUP sent"
  elif docker exec new-haproxy pgrep -f certbot >/dev/null 2>&1 || docker exec new-haproxy pgrep -f certs.sh >/dev/null 2>&1; then
    echo "new-haproxy is mid certbot walk via authoritative runme.sh — do NOT restart"
  else
    echo "haproxy not ready yet inside container — skip HUP (do not re-run runme)"
  fi
else
  echo "new-haproxy not running — starting via ./runme.sh"
  ./runme.sh
fi
EOS
fi

echo "--- wait for :443 ---"
for i in $(seq 1 90); do
  if curl -sk --connect-timeout 2 -o /dev/null -w '' https://dasmlab.org/ 2>/dev/null; then
    echo "edge listening (try $i)"
    break
  fi
  sleep 10
done

echo "--- verify vanity hosts ---"
fail=0
for url in \
  https://mock-me.dasmlab.org/ \
  https://mock-me.dasmlab.org/demo \
  https://interview-me.dasmlab.org/ \
  https://interview-me.dasmlab.org/demo \
  https://camera-scrape.dasmlab.org/
do
  code="$(curl -skL -o /dev/null -w '%{http_code}' --connect-timeout 8 --max-time 20 "$url" || echo ERR)"
  echo "${code} ${url}"
  if [[ "$code" == "503" || "$code" == "000" || "$code" == "ERR" ]]; then
    fail=1
  fi
done
if [[ "$fail" -ne 0 ]]; then
  echo "ERROR: one or more vanity hosts still failing" >&2
  exit 1
fi
echo "Vanity edge OK"
