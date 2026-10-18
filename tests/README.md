# CustardFin Test Suite

## Test Strategy (Test Pyramid)
1. **Unit tests**: Fast tests for domain logic and contracts with fakes (`libs/contracts`, `apps/edge/internal/domain`, etc.).
2. **Contract tests**: Schema and adapter compliance against all interfaces (`libs/contracts`).
3. **Integration tests**: Running against ephemeral containers (PostgreSQL, Redis, MinIO) in `deploy/compose.dev.yml`.
4. **E2E tests** (`tests/e2e/`): End-to-end media playback flows across LAN, cache, and tunnel paths.
5. **Load tests** (`tests/load/`):
   - 10 viewers on the same title (validates Origin Shield).
   - 10 viewers on different titles (validates admission control and 429 backoff).
