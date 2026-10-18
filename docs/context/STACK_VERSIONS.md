# Stack Versions and References

Record all pinned dependency versions, verification dates, and official documentation URLs here before writing code.

| Component | Target Version | Documentation URL | Notes / Verification Date |
|---|---|---|---|
| Go | 1.27.2 | https://go.dev/doc/ | Installed at `C:\Program Files\Go\bin\go.exe`; tests verified |
| SQLite Driver | `modernc.org/sqlite` v1.60.1 | https://pkg.go.dev/modernc.org/sqlite | Pure-Go, cgo-free |
| PostgreSQL Driver | `github.com/jackc/pgx/v5` v5.11.0 | https://pkg.go.dev/github.com/jackc/pgx/v5 | High-performance Go PostgreSQL driver |
| Password Hashing | `golang.org/x/crypto/bcrypt` v0.58.0 | https://pkg.go.dev/golang.org/x/crypto/bcrypt | Secure password hashing |
| JWT Tokens | `github.com/golang-jwt/jwt/v5` v5.3.1 | https://pkg.go.dev/github.com/golang-jwt/jwt/v5 | HMAC-SHA256 bearer tokens |
| Redis Client | `go-redis/v9` | https://redis.uptrace.dev/ | Used for node heartbeats only |
| SQL Generator | `sqlc` | https://docs.sqlc.dev/ | Type-safe SQL code generation |
| Node.js | v24.14.0 (npm 11.9.0) | https://nodejs.org/docs/latest/ | Node runtime only (no Bun) |
| Next.js | 15.1.0 | https://nextjs.org/docs | App router, React Server Components (React 19.0.0) |
| Video Player | `hls.js` v1.5.17 | https://github.com/video-dev/hls.js/ | Client-side HLS playback with rendition switching |
| Reverse Proxy | Caddy 2.x | https://caddyserver.com/docs/ | Automated TLS for demo profile |
| Tunnel | FRP v0.52+ (TOML config) | https://github.com/fatedier/frp | Managed subprocess frpc client adapter, outbound TLS, CGNAT bypass |
| Transcoder | FFmpeg 8.0.1 | https://ffmpeg.org/documentation.html | Verified locally: `h264_nvenc`, `hevc_nvenc`, `av1_nvenc` active on RTX |
