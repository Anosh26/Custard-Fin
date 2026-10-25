$host.UI.RawUI.WindowTitle = "CustardFin [2] Cloud Control API (:8080)"
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

Write-Host "================================================================================" -ForegroundColor Yellow
Write-Host " CUSTARDFIN CLOUD CONTROL API (apid)                                            " -ForegroundColor Yellow
Write-Host " REST & Stream Router: http://localhost:8080                                  " -ForegroundColor Yellow
Write-Host " Architecture: Hexagonal Ports & Adapters, PostgreSQL Catalog, Redis Heartbeats " -ForegroundColor Yellow
Write-Host "================================================================================" -ForegroundColor Yellow

$env:API_ADDR = ":8080"
$env:DATABASE_URL = "postgres://custard:custard_secret@127.0.0.1:5433/custardfin?sslmode=disable"
$env:REDIS_ADDR = "127.0.0.1:6379"
$env:S3_ENDPOINT = "http://127.0.0.1:9000"
$env:S3_BUCKET = "custardfin-origin-cache"
$env:S3_ACCESS_KEY = "minioadmin"
$env:S3_SECRET_KEY = "minioadmin"
$env:TUNNEL_BASE_URL = "http://localhost:8090"

go run ./apps/api/cmd/apid
