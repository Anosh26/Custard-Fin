# CustardFin Glossary

- **Edge Node / Home Node**: The user's local PC running `edged` storing media binaries.
- **Cloud Control Tier**: The centralized API and database orchestrating users, catalog discovery, and streaming routing.
- **Manifest**: A compact JSON document (≤ 3 KB) containing media metadata, content hash, duration, and rendition details.
- **Origin Shield**: Temporary cache layer (S3-compatible) in the cloud that caches hot HLS segments to prevent duplicate pulls across the home upload link.
- **FRP (Fast Reverse Proxy)**: High-performance reverse proxy used for outbound TLS tunnel creation to bypass CGNAT.
- **HLS (HTTP Live Streaming)**: Video streaming protocol splitting media into playlist (`.m3u8`) and 4-second transport chunks (`.ts` / `.m4s`).
- **Heartbeat**: 15-second periodic ping from home node to cloud API registered in Redis key with TTL.
- **Admission Control**: Mechanism enforcing a hard cap of 10 concurrent streams per home node.
