#!/usr/bin/env bash
# Copy GHCR pull secret from dasmlab-home-system into a preview namespace.
set -euo pipefail

OWNER="${1:?owner}"
OWNER="$(echo "${OWNER}" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g; s/^-+//; s/-+$//; s/-+/-/g' | cut -c1-20)"
NS="dasmlab-home-dev-${OWNER}"
SRC_NS="${SRC_NS:-dasmlab-home-system}"

oc get ns "${NS}" >/dev/null 2>&1 || oc create ns "${NS}"

if oc get secret dasmlab-ghcr-pull -n "${SRC_NS}" >/dev/null 2>&1; then
  oc get secret dasmlab-ghcr-pull -n "${SRC_NS}" -o yaml \
    | sed -e "/namespace:/d" -e "/resourceVersion:/d" -e "/uid:/d" -e "/creationTimestamp:/d" \
    | oc apply -n "${NS}" -f -
else
  echo "WARN: dasmlab-ghcr-pull missing in ${SRC_NS}"
fi

oc -n "${NS}" apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: openshift-gitops-argocd-application-controller
  namespace: ${NS}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
  - kind: ServiceAccount
    name: openshift-gitops-argocd-application-controller
    namespace: openshift-gitops
EOF

oc adm policy add-scc-to-user anyuid -z dasmlab-home-sa -n "${NS}" 2>/dev/null || true

echo "Bootstrap complete for ${NS}"
