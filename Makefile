.PHONY: help dev-up dev-down floci-up floci-down demo-up demo-down demo-check seed bench load-test test lint fmt run-edge

help:
	@echo "CustardFin Development Targets:"
	@echo "  make dev-up     - Start dev dependencies via docker-compose (MinIO)"
	@echo "  make dev-down   - Tear down dev dependencies"
	@echo "  make floci-up   - Start local cloud emulator stack via Floci (port 4566)"
	@echo "  make floci-down - Tear down Floci local cloud emulator stack"
	@echo "  make demo-up    - Run stack with demo profile and Caddy"
	@echo "  make demo-down  - Stop demo stack"
	@echo "  make demo-check - Run health checks on all active components"
	@echo "  make run-edge   - Interactive Edge daemon launcher (prompts for Cloud IP)"
	@echo "  make seed       - Seed media fixtures for testing and demo"
	@echo "  make bench      - Run GPU vs CPU transcoding benchmark"
	@echo "  make load-test  - Run 10-viewer load and admission simulation"
	@echo "  make test       - Run unit and contract tests"
	@echo "  make lint       - Run quality gates across all modules"
	@echo "  make fmt        - Format Go and TS code"

dev-up:
	docker compose -f deploy/compose.dev.yml up -d

dev-down:
	docker compose -f deploy/compose.dev.yml down

floci-up:
	docker compose -f deploy/floci/compose.floci.yml up -d

floci-down:
	docker compose -f deploy/floci/compose.floci.yml down

demo-up:
	docker compose -f deploy/compose.demo.yml up -d

demo-down:
	docker compose -f deploy/compose.demo.yml down

demo-check:
	@echo "Checking services..."
	curl -f http://localhost:8080/healthz || exit 1
	@echo "All services healthy!"

seed:
	@bash scripts/seed-media.sh || powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1

bench:
	@bash scripts/bench-gpu.sh || powershell -ExecutionPolicy Bypass -File scripts/bench-gpu.ps1

load-test:
	@bash scripts/load-10-viewers.sh || powershell -ExecutionPolicy Bypass -File scripts/load-10-viewers.ps1

run-edge:
	@bash scripts/run-edge.sh || powershell -ExecutionPolicy Bypass -File scripts/run-edge.ps1

test:
	go test -v ./libs/contracts/... ./apps/edge/... ./apps/api/...

lint:
	go vet ./...
	cd apps/web && npm run lint

fmt:
	gofmt -s -w .

