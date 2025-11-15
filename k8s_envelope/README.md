# Kubernetes Deployment Manifests

This directory contains Kubernetes deployment manifests for `dasmlab-home`.

## Files

- **`dasmlab-home_deploy.yaml`** - k3s deployment with MetalLB LoadBalancer
  - Uses MetalLB BGP pool for external IP allocation
  - Service type: LoadBalancer
  - External port: 8111
  - LoadBalancer IP: 192.168.19.100

- **`dasmlab-home_deploy-ocp.yaml`** - OpenShift Container Platform deployment
  - Uses OpenShift Route for external access
  - Service type: ClusterIP
  - Route host: dasmlab.org
  - TLS termination: edge
  - Includes imagePullSecrets for GHCR authentication

## Usage

The GitHub Actions workflow automatically selects the appropriate manifest based on the `CLUSTER_TYPE` environment variable or workflow input:

- **k3s** (default): Uses `dasmlab-home_deploy.yaml`
- **ocp**: Uses `dasmlab-home_deploy-ocp.yaml`

### Manual Selection

To manually select which manifest to use, set the `CLUSTER_TYPE` environment variable:

```bash
export CLUSTER_TYPE=ocp  # or k3s
```

Or use workflow_dispatch with the `cluster_type` input parameter.

## Prerequisites

### For k3s:
- MetalLB configured with BGP pool
- Network access to LoadBalancer IP range

### For OCP:
- OpenShift cluster with Route controller enabled
- Namespace `dasmlab-home-system` created
- Image pull secret `dasmlab-ghcr-pull` created in namespace
- DNS configured to point `dasmlab.org` to the OpenShift router

## Namespace Setup

For OCP, ensure the namespace and image pull secret exist:

```bash
# Create namespace
oc create namespace dasmlab-home-system

# Create image pull secret (if not already exists)
oc create secret docker-registry dasmlab-ghcr-pull \
  --docker-server=ghcr.io \
  --docker-username=lmcdasm \
  --docker-password="<your-token>" \
  --docker-email=dasmlab-bot@dasmlab.org \
  --namespace=dasmlab-home-system
```

