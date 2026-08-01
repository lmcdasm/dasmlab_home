# cheapcloud footprint register

Cluster URL `https://cheapcloud-dasmlab.apps.2026-prod-1.ocp.dasmlab.org` returned **404** on
`PUT /api/v1/assets/:id` (2026-08-01). That route exists in current cheapcloud source
(`internal/api/assets_routes.go`) but not in the deployed image — **redeploy a new
cheapcloud build** that includes product-footprint + media FinOps routes.

Until then, prove locally:

```bash
set -a; source /home/dasm/r2_creds_dasmlab_surfing; set +a
cd /home/dasm/cheapcloud
CHEAPCLOUD_CONFIG=/tmp/cheapcloud-cfg/cheapcloud.yaml /tmp/cheapcloud-test
CHEAPCLOUD_URL=http://127.0.0.1:18080 ./scripts/register-cheapcloud-footprint.sh
curl -sS -X POST http://127.0.0.1:18080/api/v1/media/consumption/scrape | python3 -m json.tool
curl -sS http://127.0.0.1:18080/api/v1/media/consumption/latest | python3 -m json.tool
```
