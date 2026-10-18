# libs/contracts

Shared domain contracts, schemas, and cryptographic helpers for CustardFin.

## Purpose
Enforces system-wide invariants across the Edge daemon, Cloud API, and Web frontend:
1. Compact media manifest (<= 3 KB payload, pre-encoded 4s HLS chunks).
2. Node heartbeat payload with active streams limit (max 10 concurrent streams).
3. HMAC-SHA256 chunk URL signing and verification with expiration and node/media binding.

## Public Interface
- `Manifest`: Schema validation (`Validate() error`) and serialization.
- `Heartbeat`: Status reporting and admission cap validation (`Validate() error`).
- `SignChunkURL` / `VerifyChunkSignature`: Ephemeral signed URL generation and verification.

## Testing
```bash
go test -v ./...
```
