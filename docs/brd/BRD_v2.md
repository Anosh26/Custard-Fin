# CustardFin Business Requirements Document (BRD_v2)

Team: Anish Ghadi (25P0320009), Anosh Chodankar (25P0320010)
Subject: MCA Cloud Computing Lab Project - End Semester Demo

---

## 1. Executive Summary
CustardFin is a hybrid personal media server (Jellyfin style) and decentralized community sharing platform (YouTube style). Video binaries remain on the owner's home machine; only metadata manifests (~3 KB JSON) are synced to the cloud catalog.

---

## 2. Functional Requirements (FR-1 to FR-10)

### FR-1: Local Media Library & SQLite Storage
- Home edge daemon (`apps/edge`) recursively scans configured media directories.
- Maintains video file paths, content hashes (SHA-256), and private watch progress in a local pure-Go SQLite database.
- Private watch history must never leave the home node (Invariant 8).

### FR-2: Automated HLS Transcoding & Quality Ladders
- Pre-encodes media into 4-second HLS chunks (Invariant 6).
- Supports quality renditions: 1080p, 720p, 480p.
- Probe order: NVENC hardware acceleration -> QSV/VAAPI -> CPU fallback (x264).

### FR-3: Local Direct LAN Streaming
- When viewers are on the same local subnet as the home node, stream directly over LAN without traversing the cloud or relay tunnel.

### FR-4: Centralized Cloud Control & User Auth
- Cloud API (`apps/api`) provides user registration and authentication via JWT.
- Stores registered home nodes and public catalog metadata in PostgreSQL.

### FR-5: Public Sync & Compact Manifests
- Marking media "Public" generates and syncs a <= 3 KB JSON manifest to the cloud catalog.
- Video files are never uploaded to permanent cloud storage (Invariant 1).

### FR-6: NAT / CGNAT Traversal via Outbound Reverse Tunnel
- Home edge nodes connect outbound-only using FRP client to the central FRP relay server (Invariant 2).
- Never requires open inbound router ports.

### FR-7: Node Heartbeat & Liveness State
- Home nodes send a heartbeat every 15 seconds to `POST /v1/nodes/{id}/heartbeat`.
- Redis stores `heartbeat:{node_id}` with TTL (Invariant 5).
- If the node is offline, indicate "node offline" and only serve previously cached chunks.

### FR-8: Origin Shield Pull-Through Object Cache
- S3-compatible temporary cache (e.g. MinIO) caches hot HLS segments requested by remote viewers.
- Home node uploads hot chunks once; subsequent viewers are served from the cache (Invariant 7).
- Auto-purge after 24 hours (LRU).

### FR-9: Ephemeral Signed URLs & Security
- All chunk access URLs are signed with HMAC-SHA256, bound to node ID, media ID, chunk path, and short expiry (Invariant 9).

### FR-10: Admission Control & Bandwidth Protection
- Hard per-node cap of 10 concurrent distinct remote streams (Invariant 10).
- When saturated, return HTTP 429 with `Retry-After: 30` ("node busy").
- Dynamic fallback to lower renditions (720p/480p) and first-chunk prefetching.
