# S3 Chunk Cache Adapter

`apps/api/internal/adapters/s3` implements `ports.ChunkCache` using the official AWS SDK for Go v2.

## Purpose
Provides the Origin Shield temporary object cache for CustardFin:
- **Zero Permanent Cloud Storage (Invariant 1)**: Objects are tagged with `retention=24h&managed-by=custardfin` for automatic LRU/lifecycle purging.
- **Origin Shield (Invariant 7)**: Hot HLS chunks requested by remote viewers are pulled once from the home node and cached in S3.
- **Cloud-Agnostic**: Compatible with AWS S3, MinIO, Garage, SeaweedFS, Cloudflare R2, and Floci via path-style endpoint configuration.

## Public Interface
- `Config`: Connection configuration (`Endpoint`, `Bucket`, `Region`, `AccessKeyID`, `SecretAccessKey`, `UsePathStyle`).
- `NewS3ChunkCache(cfg Config) (*S3ChunkCache, error)`: Standard constructor.
- `NewS3ChunkCacheWithClient(client S3API, bucket string) *S3ChunkCache`: Constructor with dependency injection for testing.

## Testing
Unit tests run against a mock `S3API` verifying lifecycle, tagging, and error paths:
```bash
go test -v ./apps/api/internal/adapters/s3/...
```
