.PHONY: help dev-up dev-down demo-up demo-down demo-check seed test lint fmt

help:
	@echo "CustardFin Development Targets:"
	@echo "  make dev-up     - Start dev dependencies via docker-compose"
	@echo "  make dev-down   - Tear down dev dependencies"
	@echo "  make demo-up    - Run stack with demo profile and Caddy"
	@echo "  make demo-down  - Stop demo stack"
	@echo "  make demo-check - Run health checks on all active components"
	@echo "  make seed       - Seed media fixtures for testing and demo"
	@echo "  make test       - Run unit and contract tests"
	@echo "  make lint       - Run quality gates across all modules"
	@echo "  make fmt        - Format Go and TS code"

dev-up:
	docker compose -f deploy/compose.dev.yml up -d

dev-down:
	docker compose -f deploy/compose.dev.yml down

demo-up:
	docker compose -f deploy/compose.demo.yml up -d

demo-down:
	docker compose -f deploy/compose.demo.yml down

demo-check:
	@echo "Checking services..."
	curl -f http://localhost:8080/healthz || exit 1
	@echo "All services healthy!"

seed:
	powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1

test:
	go test -v ./libs/contracts/... ./apps/edge/... ./apps/api/...

lint:
	go vet ./...
	cd apps/web && npm run lint

fmt:
	gofmt -s -w .
