# Stack Versions and References

Record all pinned dependency versions, verification dates, and official documentation URLs here before writing code.

| Component | Target Version | Documentation URL | Notes / Verification Date |
|---|---|---|---|
| Go | 1.27.2 | https://go.dev/doc/ | Installed at `C:\Program Files\Go\bin\go.exe`; tests verified |
| SQLite Driver | `modernc.org/sqlite` v1.60.1 | https://pkg.go.dev/modernc.org/sqlite | Pure-Go, cgo-free |
| PostgreSQL Driver | `pgx/v5` | https://pkg.go.dev/github.com/jackc/pgx/v5 | High-performance Go PostgreSQL driver |
| Redis Client | `go-redis/v9` | https://redis.uptrace.dev/ | Used for node heartbeats only |
| SQL Generator | `sqlc` | https://docs.sqlc.dev/ | Type-safe SQL code generation |
| Node.js | LTS (20.x or 22.x) | https://nodejs.org/docs/latest/ | Node runtime only (no Bun) |
| Next.js | 15.x / latest LTS App Router | https://nextjs.org/docs | App router, React Server Components |
| Video Player | `hls.js` | https://github.com/video-dev/hls.js/ | Client-side HLS playback |
| Reverse Proxy | Caddy 2.x | https://caddyserver.com/docs/ | Automated TLS for demo profile |
| Tunnel | FRP | https://github.com/fatedier/frp | Outbound TLS tunnel |
| Transcoder | FFmpeg 8.0.1 | https://ffmpeg.org/documentation.html | Verified locally: `h264_nvenc`, `hevc_nvenc`, `av1_nvenc` active on RTX |
