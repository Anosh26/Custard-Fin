# ==============================================================================
# CustardFin Home Edge Node Launcher (Interactive Cloud IP Selector - PowerShell)
# Prompts for the GCP / Cloud VM External IP on every run, with smart default
# caching so you can simply press Enter if the IP hasn't changed.
# ==============================================================================

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir
$CacheFile = Join-Path $RootDir ".last_cloud_ip"

Set-Location $RootDir

$DefaultIP = "136.71.107.75"
if (Test-Path $CacheFile) {
    $Cached = (Get-Content $CacheFile -ErrorAction SilentlyContinue).Trim()
    if ($Cached) {
        $DefaultIP = $Cached
    }
}

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host "  🍮 CustardFin Edge Node Launcher (RTX NVENC / Hybrid Edge)" -ForegroundColor Yellow
Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host ""

$Prompt = "Enter Cloud VM External IP [Press Enter for $DefaultIP]: "
$InputIP = Read-Host -Prompt $Prompt

if ([string]::IsNullOrWhiteSpace($InputIP)) {
    $CloudIP = $DefaultIP
} else {
    $CloudIP = $InputIP.Trim()
}

# Sanitize input
$CloudIP = $CloudIP -replace "^https?://", ""
$CloudIP = $CloudIP -replace "/.*$", ""
$CloudIP = $CloudIP -replace ":.*$", ""

if ([string]::IsNullOrWhiteSpace($CloudIP)) {
    Write-Host "❌ Error: Cloud IP cannot be empty." -ForegroundColor Red
    exit 1
}

$CloudIP | Out-File -FilePath $CacheFile -Encoding ascii -NoNewline

Write-Host ""
Write-Host "🔗 Connecting Home Edge Daemon to Cloud VM:" -ForegroundColor Green
Write-Host "   • Cloud Control API : http://${CloudIP}:8080"
Write-Host "   • FRP Reverse Tunnel: ${CloudIP}:7000"
Write-Host "   • Local Stream Port : 8090"
Write-Host ""

# Environment variables
$env:EDGE_CLOUD_URL = "http://${CloudIP}:8080"
$env:EDGE_FRP_ENABLED = "true"
$env:EDGE_FRP_SERVER_ADDR = "$CloudIP"
$env:EDGE_FRP_SERVER_PORT = "7000"
$env:EDGE_FRP_AUTH_TOKEN = "custardfin_demo_frp_token"
$env:EDGE_ADDR = ":8090"

Write-Host "🚀 Starting CustardFin edged daemon..." -ForegroundColor Yellow
go run ./apps/edge/cmd/edged
