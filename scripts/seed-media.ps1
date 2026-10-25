# Seed media script for CustardFin local development and testing
# Downloads Blender Foundation open movies (e.g. Big Buck Bunny / Tears of Steel sample clips)

param (
    [string]$TargetDir = "./media"
)

if (!(Test-Path $TargetDir)) {
    New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null
}

Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host "                   CustardFin Media Library Seeder" -ForegroundColor Cyan
Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host "Target media directory: $TargetDir"

if (-not (Get-Command ffmpeg -ErrorAction SilentlyContinue)) {
    Write-Host "Error: ffmpeg is required to synthesize media fixtures." -ForegroundColor Red
    exit 1
}

# 1. Synthesize 1080p clip with audio
$clip1080 = Join-Path $TargetDir "cosmos_1080p.mp4"
if (-not (Test-Path $clip1080)) {
    Write-Host "[1/3] Generating 1080p media fixture (15s)..." -ForegroundColor Cyan
    & ffmpeg -y -v error `
        -f lavfi -i "testsrc2=duration=15:size=1920x1080:rate=30" `
        -f lavfi -i "sine=frequency=1000:duration=15" `
        -pix_fmt yuv420p -c:v libx264 -preset fast `
        -c:a aac -b:a 160k -shortest `
        $clip1080
} else {
    Write-Host "[1/3] 1080p fixture exists: $clip1080"
}

# 2. Synthesize 720p clip with audio
$clip720 = Join-Path $TargetDir "trailer_720p.mp4"
if (-not (Test-Path $clip720)) {
    Write-Host "[2/3] Generating 720p media fixture (20s)..." -ForegroundColor Cyan
    & ffmpeg -y -v error `
        -f lavfi -i "testsrc2=duration=20:size=1280x720:rate=30" `
        -f lavfi -i "sine=frequency=440:duration=20" `
        -pix_fmt yuv420p -c:v libx264 -preset fast `
        -c:a aac -b:a 128k -shortest `
        $clip720
} else {
    Write-Host "[2/3] 720p fixture exists: $clip720"
}

# 3. Synthesize 480p clip with audio
$clip480 = Join-Path $TargetDir "vintage_480p.mp4"
if (-not (Test-Path $clip480)) {
    Write-Host "[3/3] Generating 480p media fixture (30s)..." -ForegroundColor Cyan
    & ffmpeg -y -v error `
        -f lavfi -i "testsrc2=duration=30:size=854x480:rate=30" `
        -f lavfi -i "sine=frequency=880:duration=30" `
        -pix_fmt yuv420p -c:v libx264 -preset fast `
        -c:a aac -b:a 96k -shortest `
        $clip480
} else {
    Write-Host "[3/3] 480p fixture exists: $clip480"
}

Write-Host "--------------------------------------------------------------------------------"
Write-Host "Media library successfully seeded in $TargetDir" -ForegroundColor Green
Get-ChildItem -Path $TargetDir -Filter "*.mp4" | Select-Object Name, Length, LastWriteTime | Format-Table -AutoSize
Write-Host "================================================================================" -ForegroundColor Cyan

