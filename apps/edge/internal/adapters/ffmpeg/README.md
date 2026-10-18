# apps/edge/internal/adapters/ffmpeg

Managed subprocess FFmpeg transcoder for CustardFin Edge daemon.

## Purpose
Pre-encodes discovered local media into multi-rendition HLS streams (1080p, 720p, 480p) chunked into strict 4-second segments.
Adheres to:
- **Invariant 6**: Pre-encoded 4s HLS chunks, never live transcoding.
- **Hardware Acceleration**: Probes and leverages NVIDIA NVENC (`h264_nvenc`) with automatic CPU fallback (`libx264`).
- **Keyframe Alignment**: Forces keyframes on 4s intervals (`-force_key_frames "expr:gte(t,n_forced*4)"`) ensuring identical chunk time boundaries across all quality renditions.

## Public Interface
- `ProbeEncoder(ctx context.Context) (string, error)`: Probes hardware encoder availability.
- `SelectRenditions(srcWidth, srcHeight int) []RenditionProfile`: Generates quality ladder without upscaling.
- `Transcode(ctx context.Context, job ports.TranscodeJob) (*ports.TranscodeResult, error)`: Produces `master.m3u8`, rendition playlists (`index.m3u8`), and `.ts` chunk files.

## Testing
```bash
go test -v ./apps/edge/internal/adapters/ffmpeg
```
