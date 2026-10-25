$host.UI.RawUI.WindowTitle = "CustardFin [3] Home Edge Daemon (:8090) [NVENC GPU]"
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

Write-Host "================================================================================" -ForegroundColor Green
Write-Host " CUSTARDFIN HOME EDGE DAEMON (edged)                                            " -ForegroundColor Green
Write-Host " Local Streaming: http://localhost:8090 | GPU: NVIDIA NVENC Hardware Accelerated" -ForegroundColor Green
Write-Host " Invariants: Outbound-Only TLS, 4s HLS Segments, Private Edge SQLite WAL DB    " -ForegroundColor Green
Write-Host "================================================================================" -ForegroundColor Green

$env:EDGE_ADDR = ":8090"
$env:EDGE_MEDIA_DIR = "./media"
$env:EDGE_TRANSCODE_DIR = "./transcoded"
$env:EDGE_DB_PATH = "./data/edge.db"
$env:EDGE_CLOUD_URL = "http://localhost:8080"
$env:EDGE_NODE_ID = "898ea0c7-f7a3-48e8-ba59-b378a219f9ee"
$env:EDGE_NODE_TOKEN = "1c94b7b8fd38a97cd634b4894358fc1feed28f3a180a104aa0b83f9710190de2"

go run ./apps/edge/cmd/edged
