# Architecture & Data Flows

## Two-Tier Model
CustardFin isolates the storage of large video assets on the edge while centralizing discovery and playback coordination.

```
+------------------------------------+         +-------------------------------------+
|        HOME EDGE NODE (Go)         |         |         CLOUD CONTROL (Go)          |
|                                    |         |                                     |
|  +------------------------------+  |         |  +-------------------------------+  |
|  | Scanner -> Transcoder -> HLS |  |         |  | PostgreSQL (users & catalog)  |  |
|  +------------------------------+  |         |  +-------------------------------+  |
|                 |                  |         |                 |                   |
|  +------------------------------+  | Manifest|  +-------------------------------+  |
|  | Local SQLite (watch progress)|  | ========>  | Redis (node heartbeats only)  |  |
|  +------------------------------+  |  (3 KB) |  +-------------------------------+  |
|                 |                  |         |                 |                   |
|  +------------------------------+  | Tunnel  |  +-------------------------------+  |
|  | FRP Client (outbound TLS)    |  | ========>  | FRP Relay Server              |  |
|  +------------------------------+  |         |  +-------------------------------+  |
+------------------------------------+         |                 |                   |
                                               |  +-------------------------------+  |
                                               |  | S3-compatible Origin Shield   |  |
                                               |  | (24h temporary chunk cache)   |  |
                                               |  +-------------------------------+  |
                                               +-------------------------------------+
```

## Playback Decision Flowchart
1. **Node status check**: If `heartbeat:{node_id}` expired in Redis, return offline or cached chunks only.
2. **Subnet check**: If viewer is on the same local subnet, return direct LAN URL (`http://<lan-ip>:<edge-port>/...`).
3. **Cache check**: If chunk exists in Origin Shield, return signed S3/CDN URL.
4. **Cache miss**: Pull once over FRP tunnel from edge, write to Origin Shield, return signed URL.
5. **Admission check**: If node reports 10 active remote streams, return HTTP 429 (`node busy`).
