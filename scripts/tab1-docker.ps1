$host.UI.RawUI.WindowTitle = "CustardFin [1] Docker Cloud Infrastructure"
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host " CUSTARDFIN DOCKER CLOUD INFRASTRUCTURE (PostgreSQL, Redis, MinIO S3, FRPS)      " -ForegroundColor Cyan
Write-Host "================================================================================" -ForegroundColor Cyan

docker compose -f deploy/compose.dev.yml up -d
Write-Host ""
docker compose -f deploy/compose.dev.yml ps
Write-Host ""
Write-Host "Streaming live container logs (Press Ctrl+C to stop log streaming):" -ForegroundColor Gray
docker compose -f deploy/compose.dev.yml logs -f --tail=30
