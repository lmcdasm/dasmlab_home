# 🏠 dasmlab_home

**`dasmlab_home`** is the main landing page and interactive dashboard for the DASMLAB ecosystem — a living lab environment showcasing projects, infrastructure, and integrations. This project is designed as a modern, performant Vue.js-based SPA (Single Page Application) built with the [Quasar Framework](https://quasar.dev), and fully containerized for scalable deployment across Kubernetes environments.

---

## ✨ Key Features

- ⚡ **Vue.js 3 + Quasar-based SPA**  
  Leveraging Pinia for state management, Vue Router for navigation, and built-in Quasar components for fast UI scaffolding.

- 🔌 **Live Data via SSE (Server-Sent Events)**  
  Real-time updates from backend agents and controllers streamed directly into the interface.

- 🧩 **Modular Component Design**  
  Includes reusable ProjectCard, Gauge, StatusPanel, and ConsoleLog components, styled for dark/light mode flexibility.

- 📦 **Containerized Deployment**  
  Packaged via `nginx` as a static web asset container. Wrapped using `buildme.sh` and `runme.sh` for local or CI builds.

- 🚀 **CI/CD Ready with GitHub Actions & K8s Envelope**  
  Integrated into a 5-phase pipeline (Build, Run, Secure, Validate, Publish) that publishes to a GitOps target repository for live deployment via Flux or ArgoCD.

---

## 🛠️ Project Structure

```
dasmlab-home/
├── src/                # Quasar frontend sources (layouts, pages, components, router)
├── public/             # Static assets
├── buildme.sh          # Build wrapper script (Docker)
├── runme.sh            # Local runner script (Docker)
├── k8s_envelope/
│   └── dasmlab-home_deploy.yaml   # K8s manifest (Deployment + Service template)
├── .github/workflows/  # CI/CD pipeline
└── README.md           # This file
```

---

## ⚙️ System Requirements

To build and run this project locally or in CI, ensure the following are installed:

| Tool        | Version      | Link                                |
|-------------|--------------|-------------------------------------|
| Ubuntu      | 22.04+       | https://ubuntu.com/download         |
| Node.js     | ≥ 18.x       | https://nodejs.org/en/download      |
| npm         | ≥ 9.x        | https://docs.npmjs.com/downloading |
| Docker      | ≥ 24.x       | https://docs.docker.com/get-docker |
| Quasar CLI  | ≥ 2.x        | https://quasar.dev/start/installation |

Install Quasar globally:

```bash
npm install -g @quasar/cli
```

---

## 🧪 Development Setup

To run in local dev mode (hot reload, dev server):

```bash
quasar dev
```

To build the production container and run it locally:

```bash
./buildme.sh
./runme.sh
```

This will expose the containerized app at:  
👉 http://localhost:8080

---

## 🔄 CI/CD + Deployment Pipeline

This project is fully wired into a GitHub Actions pipeline that follows a **5-phase structure**:

1. **Build** – Quasar app is built into static assets  
2. **Run** – Container image is assembled using `nginx`  
3. **Secure** – Security scans (Trivy, custom DAST/SAST hooks)  
4. **Validate** – Runtime container tests & image inspection  
5. **Publish** – Final image is pushed to registry and K8s manifest is templated and sent to:

```
lmcdasm/dasmlab_live_cicd/
└── clusters/
    └── dasmlab-prod-1/
        └── dasmlab_home/<tagged_manifest>.yaml
```

This GitOps target repo is picked up by **FluxCD** or **ArgoCD**, which applies the deployment to the live cluster.

---

## 📸 Live Architecture Visuals

> (You can include your PNG architecture diagram and UI screenshot here)

![DASMLAB Architecture](resources/infra-architecture.png)  
![UI Preview](resources/homepage-preview.png)

---

## 📬 Contact / Contribution

This is a living project within the DASMLAB ecosystem. PRs and ideas are welcome — feel free to fork and customize your own lab page or contribute enhancements!

---

## 📝 License

MIT License — see [LICENSE](LICENSE)

