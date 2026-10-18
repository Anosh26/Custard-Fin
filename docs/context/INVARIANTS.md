# Non-Negotiable Invariants

> **IMPORTANT:** Only a human may edit this file.

1. **No video binaries stored permanently in the cloud.** Only a temporary cache (auto-purge 24 h, LRU).
2. **Home nodes only make outbound connections.** Never require an open inbound router port (ISP CGNAT).
3. **100% free and open-source dependencies.** No paid licences, no proprietary SDKs required to run.
4. **Backend is Go. Frontend is Next.js (UI only, Node runtime, no Bun).** No BullMQ, no extra queue libraries.
5. **Redis is used for heartbeats only.** Background jobs are Go goroutines with state in PostgreSQL.
6. **Pre-encoded HLS (4 s chunks), not live transcoding.** GPU (NVENC) is optional with automatic CPU fallback.
7. **Origin Shield:** home uploads each hot chunk once; the cloud cache serves viewers.
8. **Private data stays home:** watch history lives in the edge SQLite.
9. **All chunk access uses short-lived signed URLs.** All API calls are authenticated.
10. **Per-node cap of 10 concurrent distinct remote streams**, with a clear "node busy" response.
