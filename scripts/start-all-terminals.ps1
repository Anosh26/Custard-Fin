# CustardFin — Multi-Terminal Presentation Launcher
# Spawns all tiers across separate labeled tabs in Windows Terminal (or separate PowerShell windows)

$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host " CustardFin — Multi-Terminal Demo Launcher (Docker + Multi-Tier Stack)         " -ForegroundColor Cyan
Write-Host "================================================================================" -ForegroundColor Cyan

# 1. Ensure Docker Cloud Infrastructure is healthy
Write-Host "[1/5] Checking Docker Cloud Infrastructure..." -ForegroundColor Yellow
docker compose -f deploy/compose.dev.yml up -d
if ($LASTEXITCODE -ne 0) {
    Write-Warning "Docker compose failed. Ensure Docker Desktop is running."
} else {
    Write-Host "  -> Docker services (Postgres, Redis, MinIO S3, FRPS) are active!" -ForegroundColor Green
}

# 2. Check for Windows Terminal (wt.exe)
$hasWt = (Get-Command wt.exe -ErrorAction SilentlyContinue) -ne $null

if ($hasWt) {
    Write-Host "[2/5] Launching Windows Terminal with 5 synchronized tabs..." -ForegroundColor Cyan

    $cmd = @(
        "-w", "0",
        "nt", "--title", "1. Docker Infra", "-d", "$projectRoot", "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab1-docker.ps1",
        ";",
        "nt", "--title", "2. Cloud API (:8080)", "-d", "$projectRoot", "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab2-apid.ps1",
        ";",
        "nt", "--title", "3. Edge Daemon (:8090)", "-d", "$projectRoot", "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab3-edged.ps1",
        ";",
        "nt", "--title", "4. Next.js Web (:3000)", "-d", "$projectRoot\apps\web", "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab4-web.ps1",
        ";",
        "nt", "--title", "5. Simulation & Viva Tests", "-d", "$projectRoot", "powershell.exe", "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab5-bench.ps1"
    )

    Start-Process wt.exe -ArgumentList $cmd
    Write-Host "  -> Windows Terminal opened with 5 tabs!" -ForegroundColor Green
} else {
    Write-Host "[2/5] Windows Terminal not found. Spawning 5 separate PowerShell windows..." -ForegroundColor Yellow

    Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab1-docker.ps1"
    Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab2-apid.ps1"
    Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab3-edged.ps1"
    Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab4-web.ps1"
    Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\tab5-bench.ps1"
}

Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host " CustardFin is now active across all terminals!" -ForegroundColor Green
Write-Host " Access the application at: http://localhost:3000" -ForegroundColor Cyan
Write-Host "================================================================================" -ForegroundColor Cyan
