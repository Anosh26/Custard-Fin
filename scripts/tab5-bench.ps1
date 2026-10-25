$host.UI.RawUI.WindowTitle = "CustardFin [5] Live Simulation & Benchmarks"
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

Write-Host "================================================================================" -ForegroundColor White
Write-Host " CUSTARDFIN DEMO SIMULATION & BENCHMARKS TERMINAL                              " -ForegroundColor White
Write-Host "================================================================================" -ForegroundColor White
Write-Host ""
Write-Host "Ready for Evaluator Presentation Tests:" -ForegroundColor Cyan
Write-Host "  1. Run GPU NVENC Benchmark:   powershell -File scripts/bench-gpu.ps1" -ForegroundColor Yellow
Write-Host "  2. Run 10-Viewer Simulation:  powershell -File scripts/load-10-viewers.ps1" -ForegroundColor Yellow
Write-Host "  3. Seed Demo Media:           powershell -File scripts/seed-media.ps1" -ForegroundColor Yellow
Write-Host ""
Write-Host "Type any command above and press Enter to execute." -ForegroundColor Gray
