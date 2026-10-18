# CustardFin Project State

## Current Milestone
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.

## Done
- Initialized `AGENTS.md` and `CLAUDE.md`.
- Initialized `docs/context/INVARIANTS.md`.
- Initialized `docs/context/STATE.md`.

## In Progress
- Scaffolding repository structure (`apps/`, `libs/contracts/`, `docs/`, `deploy/`, `scripts/`, `tests/`).
- Initializing Go workspace (`go.work`), contracts, and memory files (`PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`).

## Blocked
- None.

## Next 3 Tasks
1. Complete Tier 1 context files (`PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`) and basic docs scaffolding (`docs/api/openapi.yaml`, `docs/adr/`).
2. Scaffold `go.work` and modules: `libs/contracts`, `apps/edge`, `apps/api`.
3. Set up CI workflow (`.github/workflows/ci.yml`), `Makefile`, and `README.md`.

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
