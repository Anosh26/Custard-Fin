## 2026-10-25 - Live Cloud Deployment (GCP VM), FRP Reverse Tunnel & Mobile End-to-End Verification
- **Author/Agents:** CustardFin Pair Programming Agent (Anish Ghadi, Anosh Chodankar & Agent)
- **Milestone:** Milestone 10 (Live Cloud Deployment & End-to-End Production Verification)
- **Summary:** Successfully deployed the CustardFin cloud tier on a live GCP Compute Engine VM (`34.132.188.14`), connected the home edge daemon via outbound FRP reverse tunnel, and verified end-to-end streaming:
  1. Cloud Tier Live Orchestration: Stood up all 7 cloud services (`postgres`, `redis`, `minio/cache`, `frps`, `api`, `web`, `caddy`) on GCP using Docker Compose. Resolved Docker Hub MinIO pull access by switching to Chainguard/pgsty MinIO images, and resolved multi-module Go build issues in Docker by setting `GOWORK=off`.
  2. FRP Reverse Tunnel Routing: Fixed FRP `customDomains` configuration and added `req.Host = nodeID` header propagation in `apps/api/internal/adapters/tunnel/puller.go` to properly route HTTP tunnel requests through `frps` to the home edge daemon without requiring open router ports (Invariant 2).
  3. Private Network Access (PNA) & CORS: Added `Access-Control-Allow-Private-Network: true` and comprehensive CORS headers in edge HTTP streaming daemon so web clients accessing from public cloud IPs can interact with local edge services.
  4. Dynamic Player Routing: Configured client-side origin normalization in `apps/web/src/components/Player.tsx` and Docker SSR URL configuration (`NEXT_PUBLIC_API_URL`) so stream and playback metadata seamlessly resolve across desktop and remote mobile devices.
  5. Live Verification Across Devices & Scenarios: Confirmed flawless operation on desktop and mobile (Brave browser on Android):
     - Hardware acceleration: NVIDIA NVENC active on home edge daemon (`h264_nvenc`).
     - HLS RFC 8216: 4.0s chunks generated and aligned (Invariant 6).
     - Origin Shield Cache: Pull-through cache active (`HIT`), serving repeat chunk requests from cloud MinIO S3 with 0 home upload (Invariant 7).
     - Security: Short-lived HMAC-SHA256 signed playback URLs verified (Invariant 9).
     - Streaming Resilience: Smooth initial loading, player buffering recovery, timeline seeking/resuming, drag-and-drop video upload, and instant Public/Private catalog synchronization.
- **Files Touched:**
  - `apps/api/Dockerfile`
  - `apps/api/internal/adapters/tunnel/puller.go`
  - `apps/edge/internal/adapters/frp/config.go`
  - `apps/edge/internal/adapters/frp/config_test.go`
  - `apps/edge/internal/adapters/httpapi/server.go`
  - `apps/web/src/components/Player.tsx`
  - `apps/web/src/lib/api-client.ts`
  - `deploy/compose.demo.yml`
  - `docs/context/SESSION_LOG.md`
  - `docs/context/STATE.md`

## 2026-10-25 - Pastel Color Palette, Typography Hierarchy & Animated Custard Mascot
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Demo Presentation Aesthetics & Mascot
- **Summary:** Upgraded the CustardFin presentation system with rich pastel aesthetics, enhanced typography, OpenXML timing, and an interactive animated mascot:
  1. Custard Vector Mascot (`docs/presentation/custard_mascot.svg`): designed scalable kawaii custard mascot with caramel drip, glossy red cherry, mint leaf, serving dish, ambient glow, and sparkles.
  2. Mascot Rasterizer (`scripts/render_mascot_png.py`): created high-resolution, anti-aliased 512x512 transparent PNG (`docs/presentation/custard_mascot.png`) using Pillow with 2x supersampling and Lanczos downsampling.
  3. Curated Pastel Color Palette: implemented pastel accents (Warm Custard Cream `#FEF08A`, Sage Mint `#A7F3D0`, Soft Lavender `#DDD6FE`, Powder Sky `#BAE6FD`, Warm Peach `#FED7AA`, Gentle Coral Rose `#FECDD3`) over deep obsidian canvas (`#090D16`), providing striking contrast and readability.
  4. Typography & Layout Enhancement: integrated Google Font Plus Jakarta Sans, increased title sizes to 28pt / 50pt, card headers to 15pt, and body text to 12pt with 6pt spacing. Added card glide build entrance animations.
  5. Animated Corner Custard Mascot: mounted floating/wobbling mascot in bottom-right corner of interactive web deck (`docs/presentation/index.html`) featuring slide-specific contextual speech bubbles (e.g. "Zero Cloud Storage! ✨", "No Open Router Ports! 🔒", "NVENC 12.4x Speed! 🚀") and an interactive click easter egg with random trivia quotes.
  6. PowerPoint OpenXML Generator (`scripts/generate_presentation.py`): updated generator to embed corner mascot graphics on all 12 slides, apply pastel palette tokens and typography hierarchy, and regenerated both `docs/presentation/CustardFin_Presentation.pptx` and `docs/presentation/presentation.xml`.
  7. Visual Verification: verified all 12 slides in browser subagent, capturing photographic proof confirming typography hierarchy, pastel borders, and interactive corner mascot.
- **Files Touched:**
  - `docs/presentation/custard_mascot.svg` (new)
  - `docs/presentation/custard_mascot.png` (new)
  - `scripts/render_mascot_png.py` (new)
  - `docs/presentation/index.html`
  - `scripts/generate_presentation.py`
  - `docs/presentation/CustardFin_Presentation.pptx`
  - `docs/presentation/presentation.xml`
  - `docs/context/SESSION_LOG.md`
  - `docs/context/STATE.md`
  - `CHANGELOG.md`

## 2026-10-25 - Modern OpenXML PowerPoint Presentation & Interactive Animated Deck
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Demo Presentation Assets
- **Summary:** Engineered modern 16:9 widescreen presentation deck with OpenXML slide transitions, entrance animations, and full project lifecycle overview:
  1. OpenXML PowerPoint Generator (`scripts/generate_presentation.py`): built automated Python script using `python-pptx` and direct PresentationML OpenXML manipulation to generate `docs/presentation/CustardFin_Presentation.pptx` (12 slides, dark theme `#090D16`, cards, accent borders, `<p:transition>` slide fades/pushes, and `<p:timing>` entrance animations).
  2. OpenXML PresentationML Export (`docs/presentation/presentation.xml`): exported complete OpenXML document structure with `<PresentationML>`, slide trees, shape layouts, and animation timing nodes.
  3. Interactive Animated Web Presentation (`docs/presentation/index.html`): built standalone HTML5/CSS3 presentation deck with keyboard navigation (Arrows/Spacebar), fullscreen support (`F`), progress bar, and card hover effects.
  4. 12 Comprehensive Slides: covering Project Introduction & Team Credits, Core Problem (Plex vs YouTube), Invariant 1 (Zero Cloud Video Storage), Two-Tier Architecture, Technology Stack, Ingestion & 4s HLS Pipeline, Stream Router (Play Decision Engine), Scalability & Admission Control, NVENC Hardware Benchmarks, Multi-Topology Deployment, and Demonstration Scenarios.
- **Files Touched:**
  - `scripts/generate_presentation.py` (new)
  - `docs/presentation/CustardFin_Presentation.pptx` (new)
  - `docs/presentation/presentation.xml` (new)
  - `docs/presentation/index.html` (new)
  - `docs/presentation/README.md` (new)
  - `docs/context/SESSION_LOG.md`
  - `docs/context/STATE.md`
  - `CHANGELOG.md`

## 2026-10-25 - Fix Modal Stacking Context Overflow, Portaling & Responsive Polish
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 UI Hardening & Bug Fixes
- **Summary:** Resolved modal overflow bug, backdrop stacking containment, and UI edge cases across the web application:
  1. Root Cause Identified & Resolved: `.navbar` utilizes `backdrop-filter: blur(12px)`. In standard CSS specification, `backdrop-filter` creates an isolated stacking and containing block for all descendants, trapping `position: fixed` elements inside the 55px-tall navbar and vertically centering them at Y ≈ 27px (cutting off headers and close buttons above the top of the viewport).
  2. ModalPortal Architecture: created `apps/web/src/components/ModalPortal.tsx` to safely portal all dialog backdrops (`AddVideoModal`, `EdgeLibraryModal`, `DevSettingsModal`) directly to `document.body`, escaping any ancestor CSS containment.
  3. Overflow & Alignment Hardening: updated `.modal-backdrop` and `.modal-card` in `globals.css` with `margin: auto;`, `box-sizing: border-box;`, `max-height: calc(100vh - 5rem);`, and custom theme-matched scrollbars to prevent clipping on small viewports and mobile screens.
  4. Accessibility & UX: added Escape key listeners (`keydown`) and background scroll locking (`document.body.style.overflow = "hidden"`) while any modal is active. Added backdrop drop event prevention to prevent accidental browser file navigation. Added readability delay for upload status messages before catalog auto-refresh.
  5. Search & Watch Polish: added `useEffect` synchronization in `SearchBar.tsx` to keep the search input updated on "Clear Search" or back navigation; added `suppressHydrationWarning` on watch page URL expiry time; added `position: relative` to `.player-container` to lock the "Stats for Nerds" HUD position.
  6. Verified via automated browser subagent: captured photographic proof confirming all modals are perfectly centered, headers and close buttons fully visible, with zero overflow on home and watch pages.
- **Files Touched:**
  - `apps/web/src/components/ModalPortal.tsx` (new)
  - `apps/web/src/components/AddVideoModal.tsx`
  - `apps/web/src/components/EdgeLibraryModal.tsx`
  - `apps/web/src/components/DevSettingsModal.tsx`
  - `apps/web/src/components/SearchBar.tsx`
  - `apps/web/src/app/watch/[id]/page.tsx`
  - `apps/web/src/app/globals.css`
  - `docs/context/SESSION_LOG.md`
  - `docs/context/STATE.md`
  - `CHANGELOG.md`
- **Follow-ups:**
  - All modals and responsive views validated.

## 2026-10-25 - Multi-Terminal Docker Demo Launcher & Live Browser Demonstration
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Multi-Terminal Presentation & Demo Hardening
- **Summary:** Engineered multi-terminal presentation environment and executed live end-to-end browser walkthrough:
  1. Multi-Terminal Presentation Launcher: created `scripts/start-all-terminals.ps1`, `scripts/start-all-terminals.bat`, and 5 dedicated tab scripts (`scripts/tab1-docker.ps1` through `scripts/tab5-bench.ps1`) to launch synchronized terminal tabs in Windows Terminal (or separate PowerShell windows) with formatted banners, clear titles, and colored status logs.
  2. Live Browser Demonstration: executed browser subagent walkthrough verifying community video catalog, Dev Settings Console (1ms latency ping, Invariant compliance checklist), Edge Library manager, Add Video modal, and HLS video playback with the "Stats for Nerds" diagnostic overlay.
  3. Resolved stale webpack cache issue (`Cannot find module './68.js'`) by purging `.next` build cache and restarting dev server.
  4. Generated comprehensive walkthrough artifact `demo_live_walkthrough.md` with embedded screenshots of the Dev Settings console and watch page HUD.
- **Files Touched:**
  - `scripts/start-all-terminals.ps1`
  - `scripts/start-all-terminals.bat`
  - `scripts/tab1-docker.ps1`
  - `scripts/tab2-apid.ps1`
  - `scripts/tab3-edged.ps1`
  - `scripts/tab4-web.ps1`
  - `scripts/tab5-bench.ps1`
  - `docs/context/SESSION_LOG.md`
  - `docs/context/STATE.md`
  - `CHANGELOG.md`
- **Follow-ups:**
  - Demo stack live and ready for final evaluation and viva defense.

## 2026-10-25 - Interactive Video Management, Dev Console Settings & Player Stats HUD
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Web UI & Edge Ergonomics
- **Summary:** Enhanced CustardFin with user-friendly video management and an interactive Dev Settings Console:
  1. Edge Daemon Video Ingestion API: added `POST /api/media/upload` (multipart video upload with auto-transcode and manifest sync), `POST /api/media/scan` (on-demand directory re-scan), `POST /api/media/generate-sample` (1-click 15s synthetic FFmpeg test card generator), and `GET /api/status` (live node stream count and NVENC hardware probe).
  2. Web UI "+ Add Video" Modal: created `AddVideoModal.tsx` featuring drag-and-drop file upload, folder scan trigger, and 1-click test clip synthesizer.
  3. Edge Library Manager Modal: created `EdgeLibraryModal.tsx` showing all local edge videos with resolution, duration, storage path, and 1-click "Public / Private" community sharing switch.
  4. Dev Settings Console: created `DevSettingsModal.tsx` with dynamic Cloud API and Edge Node endpoint switcher, live latency ping tester, Invariant verification checklist (Invariants 1, 2, 5, 6, 7, 9, 10), and live admission control concurrency simulation button.
  5. Player "Stats for Nerds" HUD: added live diagnostic overlay in `Player.tsx` displaying stream protocol (HLS RFC 8216), 4s segment enforcement, active rendition/bitrate, Origin Shield cache status, and HMAC token validity.
  6. Design system upgrade in `globals.css` with glassmorphic cards, glowing status badges, and responsive navigation controls.
- **Files Touched:**
  - `apps/edge/internal/adapters/httpapi/server.go`
  - `apps/edge/cmd/edged/main.go`
  - `apps/web/src/lib/api-client.ts`
  - `apps/web/src/components/AddVideoModal.tsx`
  - `apps/web/src/components/EdgeLibraryModal.tsx`
  - `apps/web/src/components/DevSettingsModal.tsx`
  - `apps/web/src/components/NavbarClientControls.tsx`
  - `apps/web/src/components/Player.tsx`
  - `apps/web/src/app/layout.tsx`
  - `apps/web/src/app/page.tsx`
  - `apps/web/src/app/globals.css`
  - `docs/context/SESSION_LOG.md`
- **Follow-ups:**
  - Push commit to GitHub master.

## 2026-10-25 - Fix Edge Node 401 Heartbeat Rejection via Credential Pre-Seeding & Port Conflict Resolution
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Demo Reliability & Node Heartbeat Fix
- **Summary:** Diagnosed and resolved repeated `401 Unauthorized (Invalid node token)` heartbeat rejection on `apid`:
  1. Identified root cause: `apid` was attempting to connect to port 5432, which was captured by a native Windows PostgreSQL service (`postgres.exe`, PID 8292), causing SASL auth failure and silently falling back to in-memory repositories with an empty node table.
  2. Updated `deploy/compose.dev.yml` to map Docker PostgreSQL to host port `${POSTGRES_PORT:-5433}:5432`, eliminating host port collision.
  3. Updated `apps/api/cmd/apid/main.go` default `DATABASE_URL` to connect to `127.0.0.1:5433` (with fallback to 5432) and automatically run schema migrations.
  4. Added automatic pre-seeding of the development demo user (`usr-demo-...`) and edge node (`898ea0c7-f7a3-48e8-ba59-b378a219f9ee`, token `1c94b7b8fd38a97cd634b4894358fc1feed28f3a180a104aa0b83f9710190de2`) in `apid` during startup in both PostgreSQL and in-memory modes.
  5. Updated `apps/edge/cmd/edged/main.go` to automatically publish and sync scanned library manifests to the cloud API upon transcode verification.
  6. Verified end-to-end: periodic heartbeats return HTTP 200 OK, Redis stores `heartbeat:{node_id}` with NVENC GPU caps, `GET /v1/catalog` lists all 3 items with `is_node_online: true`, `GET /v1/media/{id}/playback` returns signed stream URLs, and video chunks stream with HTTP 200 via Origin Shield.
- **Files Touched:**
  - `apps/api/cmd/apid/main.go`
  - `apps/api/internal/adapters/postgres/postgres_test.go`
  - `apps/edge/cmd/edged/main.go`
  - `deploy/compose.dev.yml`
  - `docs/context/SESSION_LOG.md`
- **Follow-ups:**
  - Push master commit to GitHub.

## 2026-10-25 - Seamless Edge Node Pairing, Auto-Transcoding & Port Conflict Resolution
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Demo Polishing & Edge Integration
- **Summary:** Configured standard paired development credentials (`EDGE_NODE_ID` and `EDGE_NODE_TOKEN`) as defaults in `apps/edge/cmd/edged/main.go`, enabling `edged` to automatically dispatch periodic 10s heartbeats to Redis without requiring complex environment variable export. Integrated automatic 4-second HLS transcoding upon media library scanning using `ffmpeg` with NVENC hardware acceleration. Configured `apid` with default `TUNNEL_BASE_URL=http://localhost:8090` for seamless local dev pull-through caching. Resolved socket port binding conflicts by clearing orphaned background daemon tasks and cleaning `.next` cache. All unit, contract, and adapter test suites pass at 100%.
- **Files Touched:**
  - `apps/edge/cmd/edged/main.go`
  - `apps/api/cmd/apid/main.go`
  - `docs/context/SESSION_LOG.md`
- **Follow-ups:**
  - System fully primed for seamless 1-command execution during viva demonstration.

## 2026-10-25 - Step-by-Step Viva Demo Guides for All 4 Modes
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Viva Presentation Runbooks
- **Summary:** Authored comprehensive step-by-step viva presentation runbooks covering all 4 operational topologies of CustardFin: Mode 1 (Pure Local Zero-Docker Standalone), Mode 2 (Local Containerized with MinIO), Mode 3 (Floci Local Cloud Emulator Multi-PC Lab), and Mode 4 (Google Cloud Platform Live Cloud Evaluation). Each mode includes exact terminal configurations, copy-pasteable PowerShell/Bash commands, live speaker scripts ("what to say" to evaluators), screen output checkpoints, and viva defense Q&A. Created `docs/runbooks/DEMO_ALL_MODES.md` and `docs/runbooks/DOCKER_DEMO_GUIDE.md`, and updated `RUNNING.md` and `docs/runbooks/README.md`.
- **Files Touched:**
  - `docs/runbooks/DEMO_ALL_MODES.md`
  - `docs/runbooks/DOCKER_DEMO_GUIDE.md`
  - `docs/runbooks/README.md`
  - `RUNNING.md`
  - `docs/context/SESSION_LOG.md`
- **Follow-ups:**
  - Ready for viva evaluation across any chosen lab topology.

## 2026-10-25 - Fix S3 Cache Docker Image & Healthcheck for `make dev-up`
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Dev Tooling & Compose Fix
- **Summary:** Resolved `make dev-up` failure (`pull access denied for minio/minio, repository does not exist or may require 'docker login'`) caused by MinIO's upstream removal/access restriction of public container images from Docker Hub and Quay.io. Switched S3 cache image in `deploy/compose.dev.yml` and `deploy/compose.demo.yml` to `pgsty/minio:latest`, a fully compatible open-source MinIO distribution. Updated container healthcheck from `["CMD", "mc", "ready", "local"]` to `["CMD-SHELL", "curl -f http://localhost:9000/minio/health/live || exit 1"]`. Verified `make dev-up` runs cleanly and all 4 containers (`custardfin-cache`, `custardfin-frps`, `custardfin-postgres`, and `custardfin-redis`) are healthy. Ran unit, domain, and adapter test suites with 100% pass rate.
- **Files Touched:**
  - `deploy/compose.dev.yml`
  - `deploy/compose.demo.yml`
  - `docs/context/STACK_VERSIONS.md`
  - `docs/context/SESSION_LOG.md`
- **Follow-ups:**
  - None; stack is running and healthy.

## 2026-10-25 - PR #1 Merge Sync, Windows RTX 4060 Rehearsal & Universal Running Guide
- **Author/Agents:** CustardFin Pair Programming Agent (Anosh Chodankar & Agent)
- **Milestone:** Post-Milestone 10 Rehearsal, Multi-Topology Deployment & Universal Documentation
- **Summary:** Synchronized local repository with remote `origin/master` following the merge of PR #1 (`feat/m8-admission-prefetch`). Fast-forwarded 7 commits cleanly and verified unit and e2e test suites. Executed the complete live viva demonstration workflow on the Windows host equipped with NVIDIA GeForce RTX 4060 Laptop GPU. Upgraded `scripts/seed-media.ps1` to automatically synthesize 1080p, 720p, and 480p test video fixtures using FFmpeg. Executed `scripts/bench-gpu.ps1` verifying 16.48x realtime processing on RTX 4060 NVENC (2.73s vs 3.63s CPU baseline, 4.0s keyframe alignment). Fixed PowerShell string interpolation syntax in `scripts/load-10-viewers.ps1`. Tested `edged` daemon startup and confirmed automatic directory scanning, SHA-256 content hashing, FFprobe dimension/codec extraction, and SQLite persistence for all seeded media fixtures. Tested `apid` cloud API in local development mode and successfully ran the 10-viewer admission and Origin Shield load simulation across all 4 scenarios. Authored universal running guide `docs/runbooks/RUNNING_GUIDE.md` and `RUNNING.md` detailing all four deployment topologies (Pure Local Zero-Docker, Local Containerized, Floci Local Cloud Emulator Multi-PC Lab setup, and Google Cloud Platform). Created Floci cloud emulator compose profile (`deploy/floci/compose.floci.yml`) and guide (`deploy/floci/README.md`), adding `make floci-up` and `make floci-down` Makefile targets.
- **Files Touched:**
  - `docs/runbooks/RUNNING_GUIDE.md`
  - `RUNNING.md`
  - `deploy/floci/{compose.floci.yml, README.md}`
  - `Makefile`
  - `README.md`
  - `docs/runbooks/README.md`
  - `scripts/{seed-media.ps1, load-10-viewers.ps1}`
  - `docs/context/{STATE.md, SESSION_LOG.md, STACK_VERSIONS.md}`
- **Follow-ups:**
  - Rehearse viva presentation Q&A ahead of final evaluation.

## 2026-10-25 - Milestone 10: Demo Profile, Load Simulation Scripts, GCP Deployment & Viva Runbooks
- **Author/Agents:** CustardFin Pair Programming Agent (Anish Ghadi & Agent)
- **Milestone:** Milestone 10 (Hardening, demo profile, rehearsal, runbooks)
- **Summary:** Completed the final milestone of the CustardFin project roadmap. Built cross-platform 10-viewer load and admission simulation scripts (`scripts/load-10-viewers.sh` and `scripts/load-10-viewers.ps1`) executing live assertions for Origin Shield cache hit offload (90% reduction in home edge upload bandwidth), dynamic master playlist rendition pruning under load, the 10-stream per-node hard cap (HTTP 429 Retry-After: 30), and offline playback fallback. Created production-ready multi-stage Dockerfiles for cloud control API (`apps/api/Dockerfile`) and Next.js frontend (`apps/web/Dockerfile`). Assembled unified demo orchestration profile in `deploy/compose.demo.yml` with Caddy reverse proxy (`deploy/Caddyfile`) managing ports 80/443. Added `make load-test` target in `Makefile`. Authored comprehensive GCP deployment runbook (`docs/runbooks/deploy.md`) with VPC firewall rules, ports table, VM sizing ($40 credits budget), and home edge client configuration. Authored semester viva presentation script and defense guide (`docs/runbooks/demo.md`) with step-by-step walkthroughs and anticipated Q&A. Authored troubleshooting and emergency clean-slate recovery guide (`docs/runbooks/troubleshooting.md`).
- **Files Touched:**
  - `scripts/{load-10-viewers.sh, load-10-viewers.ps1}`
  - `apps/api/Dockerfile`
  - `apps/web/Dockerfile`
  - `deploy/{compose.demo.yml, Caddyfile}`
  - `Makefile`
  - `docs/runbooks/{README.md, deploy.md, demo.md, troubleshooting.md}`
  - `docs/context/{STATE.md, SESSION_LOG.md}`
- **Follow-ups:**
  - Rehearse live demonstration flow with Anosh Chodankar ahead of final lab evaluation.

## 2026-10-25 - Milestone 9: NVENC Hardware Path, Probe, Fallback & Benchmark Tooling
- **Author/Agents:** CustardFin Pair Programming Agent (Anish Ghadi & Agent)
- **Milestone:** Milestone 9 (NVENC path with probe, fallback and benchmark script)
- **Summary:** Completed hardware acceleration capabilities detection, multi-encoder fallback cascade, and automated GPU/CPU benchmarking tooling adhering strictly to Invariant 6 (Pre-encoded 4s HLS chunks, NVENC optional with automatic CPU fallback). In `apps/edge/internal/adapters/ffmpeg`, implemented `ProbeCapabilities` checking NVENC, QSV, and VAAPI via 1s null test encodes, refined `ProbeEncoder` to check hardware encoders before falling back to `libx264`, and updated `Transcode` with generalized automatic CPU fallback upon any hardware encoding failure. Updated `edged/main.go` to report all detected hardware acceleration flags in periodic node heartbeats (`contracts.GPUCapabilities`). Created `scripts/bench-gpu.sh` and `scripts/bench-gpu.ps1` measuring execution wall-clock time, encoding speed multiplier (achieved 10.7x realtime on GeForce GTX), and speedup ratio. Built automated media seeder `scripts/seed-media.sh` generating standard 1080p, 720p, and 480p test fixtures into `./media`. Updated `.gitignore` to prevent committing media binaries and updated `Makefile` with `make bench` and `make seed`.
- **Files Touched:**
  - `apps/edge/internal/adapters/ffmpeg/{transcoder.go, transcoder_test.go}`
  - `apps/edge/cmd/edged/main.go`
  - `scripts/{bench-gpu.sh, bench-gpu.ps1, seed-media.sh}`
  - `.gitignore`
  - `Makefile`
  - `docs/context/{STATE.md, SESSION_LOG.md}`
- **Follow-ups:**
  - Milestone 10: Hardening, demo profile, rehearsal, runbooks.

## 2026-10-24 - Milestone 8: Admission Control, Adaptive Bitrate & Prefetch
- **Author/Agents:** CustardFin Pair Programming Agent (Anish Ghadi & Agent)
- **Milestone:** Milestone 8 (Admission control, adaptive bitrate, prefetch)
- **Summary:** Delivered node saturation estimation, adaptive master playlist rendition throttling, and initial chunk prefetch worker adhering strictly to Invariant 7 (Origin Shield pull-through cache), Invariant 9 (ephemeral HMAC signed URLs), and Invariant 10 (10-stream concurrency limit). In `apps/api`, defined `ports.SaturationEstimator` and `ports.Prefetcher` in `apps/api/internal/ports/admission.go`, implemented `domain.LoadLevel` (`LoadNormal`, `LoadModerate`, `LoadSaturated`, `LoadExhausted`) with resolution filtering rules in `apps/api/internal/domain/load.go`. Implemented `admission.DefaultEstimator` (`apps/api/internal/adapters/admission`) mapping heartbeat active stream counts to load tiers. Built `prefetch.BackgroundPrefetcher` (`apps/api/internal/adapters/prefetch`) proactively warming Origin Shield cache with initial `.ts` video chunks ahead of time. Enhanced `httpapi.Server` master playlist rewriting with `rewritePlaylistWithLoad`, dynamically pruning 1080p under moderate load and 1080p+720p under saturated load (with single-rendition safety fallback). Wired asynchronous prefetch into `handleGetPlayback` and saturation estimator into `handleStreamChunk`. Added comprehensive unit test coverage and updated ADR 0002.
- **Files Touched:**
  - `apps/api/internal/domain/{load.go, load_test.go}`
  - `apps/api/internal/ports/admission.go`
  - `apps/api/internal/adapters/admission/{estimator.go, estimator_test.go, README.md}`
  - `apps/api/internal/adapters/prefetch/{prefetcher.go, prefetcher_test.go, README.md}`
  - `apps/api/internal/adapters/httpapi/{server.go, server_test.go}`
  - `tests/e2e/e2e_test.go`
  - `docs/adr/0002-admission-control-and-adaptive-bitrate.md`
  - `docs/context/{STATE.md, SESSION_LOG.md}`
- **Follow-ups:**
  - Milestone 9: NVENC path with probe, fallback and benchmark script.

## 2026-10-19 - Milestone 7: Origin Shield Cache, Signed URLs & Stream Router
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 7 (Origin Shield cache + signed URLs + stream router)
- **Summary:** Delivered the complete Origin Shield pull-through cache, HMAC-SHA256 signed URL stream enforcement, and Play Decision Tree stream router adhering strictly to Invariant 1 (temporary cache only, 24h retention tagging), Invariant 7 (hot chunks stored once to shield home node upload), Invariant 9 (ephemeral HMAC signed URLs), and Invariant 10 (10-stream concurrency limit). In `apps/api`, defined `ports.ChunkCache` and `ports.TunnelPuller`, implemented `S3ChunkCache` (`apps/api/internal/adapters/s3`) using AWS SDK for Go v2 with path-style support (MinIO/Garage/Floci) and 24h retention tagging, built `MemoryChunkCache` (`apps/api/internal/adapters/memory`) for zero-dependency test isolation, and built `HTTPTunnelPuller` (`apps/api/internal/adapters/tunnel`). In `httpapi.Server`, registered `GET /v1/stream/{node_id}/{media_id}/{chunk_path...}` enforcing HMAC signatures via `contracts.VerifyChunkSignature`, serving cached chunks on cache HIT (`X-Cache: HIT`), and pulling uncached chunks through the reverse tunnel on cache MISS (`X-Cache: MISS`). Implemented dynamic HLS `.m3u8` playlist rewriter signing relative rendition and segment URLs. Upgraded Play Decision Tree in `GET /v1/media/{id}/playback` to route between `direct_lan`, `origin_shield`, and `cache_only` (allowing cached chunk playback when node is offline per FR-7). Updated Next.js web client (`apps/web`) to display routing badges. 100% test pass rate across all unit tests and live cross-subsystem E2E test suite.
- **Files Touched:**
  - `apps/api/go.mod`, `apps/api/go.sum`
  - `apps/api/internal/ports/chunk_cache.go`
  - `apps/api/internal/adapters/s3/{chunk_cache.go, chunk_cache_test.go, README.md}`
  - `apps/api/internal/adapters/memory/{chunk_cache.go, chunk_cache_test.go, heartbeat_store.go}`
  - `apps/api/internal/adapters/tunnel/{puller.go, puller_test.go, README.md}`
  - `apps/api/internal/adapters/httpapi/{server.go, server_test.go}`
  - `apps/api/cmd/apid/main.go`
  - `apps/web/src/lib/api-client.ts`
  - `apps/web/src/app/watch/[id]/page.tsx`
  - `docs/api/openapi.yaml`
  - `tests/e2e/e2e_test.go`
  - `docs/context/{STATE.md, STACK_VERSIONS.md, SESSION_LOG.md}`
- **Follow-ups:**
  - Milestone 8: Admission control, adaptive bitrate, prefetch.

## 2026-10-18 - Milestone 6: Redis Heartbeat and Offline Logic
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 6 (Redis heartbeat and offline logic)
- **Summary:** Delivered complete end-to-end node liveness tracking, capacity admission enforcement, and offline detection adhering to Invariant 5 (Redis strictly used for `heartbeat:{node_id}` TTL keys) and Invariant 10 (10-stream concurrency limit). In `apps/api`, defined `ports.HeartbeatStore`, implemented Redis adapter (`apps/api/internal/adapters/redis`) with `go-redis/v9 v9.23.0` and fallback in-memory store (`apps/api/internal/adapters/memory`). Registered authenticated `POST /v1/nodes/{id}/heartbeat`, enriched `GET /v1/catalog` with batch node liveness (`is_node_online`), and updated `GET /v1/media/{id}/playback` router to return HTTP 503 `NODE_OFFLINE` when node key is absent/expired, and HTTP 429 `NODE_BUSY` with `Retry-After: 30` when streams reach 10-cap. In `apps/edge`, defined `ports.HeartbeatReporter` and `ports.StreamTracker`, implemented `apps/edge/internal/adapters/heartbeat/worker.go` with 15s ticker, GPU capabilities reporting (NVENC), and dynamic stream tracking. Wired both `cmd/apid` and `cmd/edged` with 12-factor configuration and graceful shutdown. All unit tests and live cross-subsystem E2E test suite passing 100%.
- **Files Touched:**
  - `apps/api/go.mod`, `apps/api/go.sum`
  - `apps/api/internal/ports/heartbeat_store.go`
  - `apps/api/internal/adapters/redis/{heartbeat_store.go, heartbeat_store_test.go, README.md}`
  - `apps/api/internal/adapters/memory/{heartbeat_store.go, heartbeat_store_test.go}`
  - `apps/api/internal/adapters/httpapi/{server.go, server_test.go}`
  - `apps/api/cmd/apid/main.go`
  - `apps/edge/internal/ports/heartbeat.go`
  - `apps/edge/internal/adapters/heartbeat/{worker.go, worker_test.go, README.md}`
  - `apps/edge/internal/adapters/httpapi/{server.go, server_test.go}`
  - `apps/edge/cmd/edged/main.go`
  - `tests/e2e/e2e_test.go`
  - `docs/context/{STACK_VERSIONS.md, STATE.md, SESSION_LOG.md}`
- **Follow-ups:**
  - Milestone 7: Origin Shield cache + signed URLs + stream router.

## 2026-10-18 - Milestone 5: FRP Reverse Tunnel Client Adapter & Deployment
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 5 (FRP tunnel: remote playback end to end)
- **Summary:** Implemented outbound-only reverse tunnel adapter `apps/edge/internal/adapters/frp` bypassing ISP CGNAT (Invariant 2). Defined `ports.TunnelClient` port. Implemented `frp.Config` with input validation, official FRP v0.52+ client TOML generator (`GenerateTOML`), and secure `0600` temp file creation. Built `ManagedFRPClient` supervising `frpc` child process with stdout/stderr line streaming, supervisory auto-reconnect backoff (1s doubling up to 30s), and clean shutdown. Wired into `apps/edge/cmd/edged/main.go` with 12-factor configuration (`EDGE_FRP_ENABLED`, `EDGE_FRP_SERVER_ADDR`, `EDGE_FRP_SERVER_PORT`, `EDGE_FRP_AUTH_TOKEN`, `EDGE_FRP_BIN`) and graceful SIGTERM shutdown. All unit tests passing with mock subprocess simulation.
- **Files Touched:**
  - `apps/edge/internal/ports/tunnel.go`
  - `apps/edge/internal/adapters/frp/{config.go, config_test.go, client.go, client_test.go, README.md}`
  - `apps/edge/cmd/edged/main.go`
  - `deploy/env.example`
  - `docs/context/{STATE.md, STACK_VERSIONS.md, SESSION_LOG.md}`
- **Follow-ups:**
  - Milestone 6: Redis Heartbeat and Offline Logic.

## 2026-10-18 - Milestone 4: Public Toggle, Manifest Sync Worker, and Next.js Web UI
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 4 (Public toggle -> manifest sync worker -> Next.js community catalog UI)
- **Summary:** Completed both Edge manifest synchronization and Next.js Web UI tiers for Milestone 4. In `apps/edge`, implemented `HTTPManifestSyncer` enforcing strict <= 3 KB JSON manifest cap and 4s segment durations, wired `PUT /api/media/{id}/public` toggle and SQLite persistence. In `apps/api`, registered `GET /v1/media/{id}/playback` route resolving stream routing decisions. In `apps/web`, built Next.js 15.1.0 App Router UI with HLS player client component (`hls.js` quality rendition switcher and Safari fallback), rich dark theme glassmorphism aesthetic (`globals.css`), community catalog discovery grid (`/`), search filtering (`?q=`), and watch route (`/watch/[id]`). Passing 100% of quality gates: `go test -v ./...` on edge and api, `npm run build`, `tsc --noEmit`, and `eslint`.
- **Files Touched:**
  - `apps/edge/internal/ports/syncer.go`
  - `apps/edge/internal/adapters/sync/{syncer.go, syncer_test.go}`
  - `apps/edge/internal/adapters/sqlite/media_repo.go`
  - `apps/edge/internal/adapters/httpapi/{server.go, server_test.go}`
  - `apps/edge/cmd/edged/main.go`
  - `apps/api/internal/adapters/httpapi/{server.go, server_test.go}`
  - `apps/web/{package.json, tsconfig.json, next.config.js, .eslintrc.json}`
  - `apps/web/src/components/{Player.tsx, SearchBar.tsx}`
  - `apps/web/src/app/{globals.css, layout.tsx, page.tsx, watch/[id]/page.tsx}`
  - `apps/web/src/lib/api-client.ts`
  - `docs/context/{STATE.md, STACK_VERSIONS.md, SESSION_LOG.md}`
- **Follow-ups:**
  - Milestone 5: FRP Tunnel — Remote Playback End-to-End.

## 2026-10-18 - Milestone 3: Cloud Control API (Auth, PostgreSQL Catalog, Node Pairing & Manifest Sync)
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 3 (Cloud API: Auth, Postgres, manifest and catalog endpoints)
- **Summary:** Built `apps/api` following Hexagonal Architecture and the Mandatory Granular Testing Policy. Implemented bcrypt password hashing, HMAC-SHA256 JWT bearer token issuance/validation, versioned PostgreSQL schema migrations (`users`, `nodes`, `catalog_items`, `catalog_renditions`), parameterized pgx query repositories, in-memory repositories for rapid isolated testing, and full OpenAPI-compliant HTTP endpoints (`/v1/auth/register`, `/v1/auth/login`, `/v1/nodes`, `/v1/nodes/{id}/media/{mid}/manifest`, `/v1/catalog`, `/healthz`). Enforced strict 3 KB manifest cap (Invariant 1) and node stream security. Wired `cmd/apid/main.go` daemon with graceful shutdown. All granular unit and integration tests passing.
- **Files Touched:**
  - `apps/api/cmd/apid/main.go`
  - `apps/api/internal/domain/{user.go, user_test.go, node.go, node_test.go, catalog.go, catalog_test.go, errors.go}`
  - `apps/api/internal/ports/{user_repository.go, node_repository.go, catalog_repository.go, auth.go}`
  - `apps/api/internal/adapters/auth/{bcrypt_hasher.go, bcrypt_hasher_test.go, jwt_service.go, jwt_service_test.go}`
  - `apps/api/internal/adapters/postgres/{db.go, migrations.go, user_repo.go, node_repo.go, catalog_repo.go, migrations/001_initial_schema.sql}`
  - `apps/api/internal/adapters/memory/{user_repo.go, node_repo.go, catalog_repo.go, memory_test.go}`
  - `apps/api/internal/adapters/httpapi/{server.go, server_test.go, middleware.go, README.md}`
  - `docs/context/STACK_VERSIONS.md`
  - `docs/context/STATE.md`
- **Follow-ups:**
  - Milestone 4: Edge public toggle -> manifest sync worker -> Next.js catalog UI.

## 2026-10-18 - Milestone 2: Edge Local HTTP Streaming Server (Standalone Mini-Jellyfin)
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Implemented `apps/edge/internal/adapters/httpapi` standard library `net/http` streaming server for direct LAN playback. Delivered CORS headers, structured logging, panic recovery, static HLS file serving with byte ranges (`Accept-Ranges: bytes`), media catalog endpoints (`GET /api/media`, `GET /api/media/{id}`), local watch progress tracking (`GET /api/watch/{id}`, `POST /api/watch/{id}`), and path traversal prevention. Wired HTTP server into `apps/edge/cmd/edged/main.go` with graceful context shutdown. Comprehensive unit and integration test suite passing. Milestone 2 is fully complete.
- **Files Touched:**
  - `apps/edge/internal/adapters/httpapi/server.go`
  - `apps/edge/internal/adapters/httpapi/server_test.go`
  - `apps/edge/internal/adapters/httpapi/README.md`
  - `apps/edge/cmd/edged/main.go`
  - `docs/context/STATE.md`
- **Follow-ups:**
  - Milestone 3: Cloud Control API (`apps/api`) setup (PostgreSQL schemas, `pgx`, `sqlc`, JWT auth, manifest sync, and catalog discovery).

## 2026-10-18 - Milestone 2: Managed Subprocess FFmpeg 4s HLS Transcoder
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Implemented `apps/edge/internal/adapters/ffmpeg` subprocess transcoder adhering to Invariant 6 (strict 4-second HLS segments). Added GPU probe utilizing NVIDIA NVENC (`h264_nvenc`) with verified minimum dimension constraints, automatic CPU fallback to `libx264`, forced keyframe multiples (`expr:gte(t,n_forced*4)`), multi-rendition quality ladder calculation without upscaling, and master playlist generation (`master.m3u8`). Verified with granular unit tests and real GPU encode integration test on RTX card (0.96s encode).
- **Files Touched:**
  - `apps/edge/internal/ports/transcoder.go`
  - `apps/edge/internal/adapters/ffmpeg/transcoder.go`
  - `apps/edge/internal/adapters/ffmpeg/transcoder_test.go`
  - `apps/edge/internal/adapters/ffmpeg/README.md`
  - `docs/context/STATE.md`

## 2026-10-18 - Milestone 2: Edge SQLite Store & Schema Migrations
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Implemented pure-Go SQLite store (`modernc.org/sqlite` v1.60.1) adhering to Invariant 8 (private data stays home). Delivered versioned schema migrations for `media` and `watch_progress` tables, WAL journal mode with `busy_timeout=5000`, single-writer connection pooling, `MediaRepository`, and `WatchProgressRepository` with cascading deletes. Wired scanner with SQLite repository in `edged/main.go`. Added comprehensive unit tests and verified with full integration test `TestScanner_IntegrationWithRealSQLite`.
- **Files Touched:**
  - `apps/edge/internal/adapters/sqlite/db.go`
  - `apps/edge/internal/adapters/sqlite/db_test.go`
  - `apps/edge/internal/adapters/sqlite/migrations.go`
  - `apps/edge/internal/adapters/sqlite/media_repo.go`
  - `apps/edge/internal/adapters/sqlite/media_repo_test.go`
  - `apps/edge/internal/adapters/sqlite/watch_repo.go`
  - `apps/edge/internal/adapters/sqlite/watch_repo_test.go`
  - `apps/edge/internal/adapters/sqlite/README.md`
  - `apps/edge/internal/domain/watch.go`
  - `apps/edge/internal/domain/watch_test.go`
  - `apps/edge/internal/ports/scanner.go`
  - `apps/edge/internal/app/scanner_integration_test.go`
  - `apps/edge/cmd/edged/main.go`
  - `docs/context/STATE.md`

## 2026-10-18 - Milestone 2: Mandatory Testing Policy & Comprehensive Test Suite
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Updated `AGENTS.md` with the Mandatory Granular Testing Policy requiring comprehensive unit tests for every package, module, and function, plus integration tests. Built and ran exhaustive test suites across `libs/contracts`, `domain`, `adapters/fs`, `adapters/ffprobe`, and `app`, including real-world integration testing with generated ffmpeg clips and subprocess ffprobe. Caught and resolved edge cases in empty/whitespace root directory handling and test variable declarations. All quality gates and tests are green.
- **Files Touched:**
  - `AGENTS.md`
  - `libs/contracts/manifest_test.go`
  - `libs/contracts/heartbeat_test.go`
  - `libs/contracts/signing_test.go`
  - `apps/edge/internal/domain/media_test.go`
  - `apps/edge/internal/adapters/fs/hasher_test.go`
  - `apps/edge/internal/adapters/fs/scanner.go`
  - `apps/edge/internal/adapters/fs/scanner_test.go`
  - `apps/edge/internal/adapters/ffprobe/prober_test.go`
  - `apps/edge/internal/app/scan_service_test.go`
  - `apps/edge/internal/app/scanner_integration_test.go`
  - `docs/context/STATE.md`
  - `apps/edge/internal/domain/media.go`
  - `apps/edge/internal/domain/media_test.go`
  - `apps/edge/internal/ports/scanner.go`
  - `apps/edge/internal/adapters/fs/hasher.go`
  - `apps/edge/internal/adapters/fs/hasher_test.go`
  - `apps/edge/internal/adapters/fs/scanner.go`
  - `apps/edge/internal/adapters/fs/scanner_test.go`
  - `apps/edge/internal/adapters/ffprobe/prober.go`
  - `apps/edge/internal/adapters/ffprobe/prober_test.go`
  - `apps/edge/internal/app/scan_service.go`
  - `apps/edge/internal/app/scan_service_test.go`
  - `docs/context/STATE.md`
- **Follow-ups:**
  - Implement edge SQLite store (`modernc.org/sqlite`) for local metadata and watch progress.

## 2026-10-18 - Session Start & Milestone 1 Initialization
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 1 (Contracts and repo skeleton, CI, docs scaffolding, memory files)
- **Summary:** Initialized `AGENTS.md` master prompt, `CLAUDE.md`, and Context Tier 0/1 memory system (`INVARIANTS.md`, `STATE.md`, `PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`).
- **Files Touched:**
  - `AGENTS.md`
  - `CLAUDE.md`
  - `docs/context/INVARIANTS.md`
  - `docs/context/STATE.md`
  - `docs/context/PROJECT_BRIEF.md`
  - `docs/context/GLOSSARY.md`
  - `docs/context/STACK_VERSIONS.md`
  - `docs/context/SESSION_LOG.md`
- **Follow-ups:**
  - Establish `libs/contracts` (Manifest & Heartbeat schema).
  - Create directory skeleton (`apps/edge`, `apps/api`, `apps/web`, `deploy`, `scripts`, `tests`).
  - Set up `go.work` and root `Makefile`.
