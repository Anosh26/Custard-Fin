# apps/edge (edged)

CustardFin Home Edge daemon running on the user's home machine.

## Purpose
- Scans local media folders and maintains local metadata/watch history in pure-Go SQLite (`modernc.org/sqlite`).
- Orchestrates FFmpeg pre-transcoding into 4s HLS segments (GPU NVENC with automatic CPU fallback).
- Serves LAN viewers directly without cloud overhead.
- Maintains outbound-only TLS tunnel to FRP relay (no inbound port forwarding needed).
- Sends periodic 15s heartbeats to Cloud Control API and syncs public manifests.

## Hexagonal Package Structure
- `cmd/edged/`: Main entrypoint and DI wiring.
- `internal/domain/`: Core entities (Media, Chunk, StreamSession, Settings).
- `internal/app/`: Use cases (ScanMedia, TranscodeJob, SyncManifest, HeartbeatRunner, LocalStreamer).
- `internal/ports/`: Interfaces for adapters (`MediaStore`, `Encoder`, `TunnelClient`, `CloudAPIClient`).
- `internal/adapters/`:
  - `sqlite/`: Local DB queries, migrations, watch progress.
  - `ffmpeg/`: Managed subprocess FFmpeg transcoder with probe order NVENC -> QSV/VAAPI -> x264.
  - `frp/`: Managed subprocess FRP client for outbound tunnel.
  - `httpapi/`: Local LAN HTTP streaming server with CORS and range support.
  - `fs/`: Local filesystem scanner.

## Running Locally
```bash
go run ./cmd/edged
```
