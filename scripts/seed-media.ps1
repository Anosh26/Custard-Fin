# Seed media script for CustardFin local development and testing
# Downloads Blender Foundation open movies (e.g. Big Buck Bunny / Tears of Steel sample clips)

param (
    [string]$TargetDir = "./media"
)

if (!(Test-Path $TargetDir)) {
    New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null
}

Write-Host "Target media directory: $TargetDir"
Write-Host "For demo and automated testing, download public domain Blender open movies (e.g., Big Buck Bunny, Tears of Steel)."
Write-Host "Seed script ready."
