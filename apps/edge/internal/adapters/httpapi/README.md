# apps/edge/internal/adapters/httpapi

Local HTTP streaming and metadata API for the CustardFin Edge daemon.

## Purpose
Provides standalone LAN media streaming (mini-Jellyfin style):
- Serves local HLS master playlists (`master.m3u8`), rendition playlists (`index.m3u8`), and 4-second `.ts` chunks.
- Enables direct LAN playback with CORS headers (`Access-Control-Allow-Origin: *`) and byte-range seeking (`Accept-Ranges: bytes`).
- Exposes local media library and manages private watch progress stored solely in the local SQLite database.

## Public Interface
- `NewServer(cfg Config, mediaRepo ports.MediaRepository, watchRepo ports.WatchProgressRepository, logger *slog.Logger) *Server`: Initializes server with CORS, logging, and recovery middleware.
- `Server.Start() error`: Starts the HTTP listener on the configured address.
- `Server.Shutdown(ctx context.Context) error`: Performs graceful shutdown.

## Endpoints
- `GET /healthz`: Health and liveness probe.
- `GET /api/media`: Lists all cataloged media items.
- `GET /api/media/{id}`: Retrieves details for a specific media item.
- `GET /stream/{id}/{path...}`: Streams HLS playlists and video chunks.
- `GET /api/watch/{id}`: Retrieves private playback position.
- `POST /api/watch/{id}`: Updates private playback position.

## Testing
```bash
go test -v ./apps/edge/internal/adapters/httpapi
```
