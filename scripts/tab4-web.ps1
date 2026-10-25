$host.UI.RawUI.WindowTitle = "CustardFin [4] Next.js Web Client (:3000)"
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location "$projectRoot\apps\web"

Write-Host "================================================================================" -ForegroundColor Magenta
Write-Host " CUSTARDFIN NEXT.JS WEB CLIENT                                                  " -ForegroundColor Magenta
Write-Host " URL: http://localhost:3000                                                     " -ForegroundColor Magenta
Write-Host " Features: HLS.js Player, 4s Segments, Stats for Nerds HUD, Dev Settings Modal " -ForegroundColor Magenta
Write-Host "================================================================================" -ForegroundColor Magenta

npm run dev
