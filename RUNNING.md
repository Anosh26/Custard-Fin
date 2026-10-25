# Running CustardFin

CustardFin supports 4 operational modes:
1. **[Mode 1: Pure Local (Zero Docker)](docs/runbooks/RUNNING_GUIDE.md#3-mode-1-pure-local-zero-docker--standalone-mode)**: In-memory fallback stores; zero external services needed.
2. **[Mode 2: Local Containerized](docs/runbooks/RUNNING_GUIDE.md#4-mode-2-local-containerized-dev-profile-with-minio)**: Full local Docker stack (MinIO, Postgres, Redis, FRP).
3. **[Mode 3: Floci Local Cloud Emulator (Multi-PC Lab)](docs/runbooks/RUNNING_GUIDE.md#5-mode-3-floci-local-cloud-emulator-multi-pc-lab-setup)**: Multi-computer lab demo with Floci S3 emulator (:4566), RTX 4060 edge node, and remote viewer clients.
4. **[Mode 4: Google Cloud Platform (GCP Live Cloud)](docs/runbooks/RUNNING_GUIDE.md#6-mode-4-google-cloud-platform-gcp-live-cloud-way)**: Real cloud deployment on GCP VM with Caddy HTTPS and outbound reverse tunnel punching through CGNAT.

👉 **Viva Presentation Guide (All 4 Modes)**: [**`docs/runbooks/DEMO_ALL_MODES.md`**](docs/runbooks/DEMO_ALL_MODES.md)  
👉 **Docker Presentation Guide**: [**`docs/runbooks/DOCKER_DEMO_GUIDE.md`**](docs/runbooks/DOCKER_DEMO_GUIDE.md)  
👉 **Complete Master Running Guide**: [**`docs/runbooks/RUNNING_GUIDE.md`**](docs/runbooks/RUNNING_GUIDE.md)
