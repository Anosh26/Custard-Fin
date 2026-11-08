# CustardFin Project Brief

## Overview
CustardFin is a hybrid personal media streaming and community sharing platform created for the MCA Cloud Computing Lab project.
- **Master Report Title:** Cloud Computing Project Report
- **Team (Identified Strictly by Seat Numbers):**
  - Anish Ghadi — Seat No.: `25P0320009`
  - Anosh Chodankar — Seat No.: `25P0320010`
- **Academic Deliverable:** 62-page compiled LaTeX report (`report/main.pdf`) with 0 errors.
- **Financial Validation:** Verified 95.89% cloud expenditure reduction (slashing monthly costs from ₹1,00,380.00/mo to ₹4,128.60/mo, saving ₹11.55 Lakhs/year at benchmark rate $1.00 USD = ₹84.00 INR).

## Core Idea
Video files remain on the owner's local machine ("home edge node"). Only a compact metadata manifest (~3 KB JSON) is synced to the central cloud catalog when marked "Public".

## Key Components
1. **Home Edge Daemon (`apps/edge`)**:
   - Media scanner & pure-Go SQLite store for local metadata and private watch history.
   - Subprocess FFmpeg pre-transcoding (HLS 4-second segments, CPU or NVENC).
   - Local direct streaming server for LAN clients.
   - FRP client establishing outbound TLS tunnel to relay.
   - Background sync worker and 15s heartbeat sender.
2. **Cloud Control Tier (`apps/api`)**:
   - Go HTTP API (`net/http`, `log/slog`, `pgx` + `sqlc`).
   - PostgreSQL: users, registered nodes, catalog metadata.
   - Redis: heartbeat TTL tracking only (`heartbeat:{node_id}`).
   - Stream Router & Origin Shield pull-through cache: fetches chunks over tunnel upon viewer demand, caches temporarily in S3-compatible object storage.
   - Token & short-lived signed URL generation.
3. **Web Client (`apps/web`)**:
   - Next.js App Router (TypeScript strict, SSR/RSC default, `hls.js` client component).
   - Talks strictly to Go API; no direct database connection.
4. **Relay & Edge Networking**:
   - Outbound FRP tunnel for CGNAT traversal.
   - Signed URLs with node/media binding and short TTLs.
