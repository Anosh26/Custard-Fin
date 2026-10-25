# ==============================================================================
# CustardFin GPU vs CPU Transcoder Benchmark Script (Milestone 9 - Windows/PS)
# Measures encoding speed, frames per second, and speedup ratio for HLS chunks.
# ==============================================================================

param (
    [int]$Duration = 15
)

$ErrorActionPreference = "Stop"

Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host "                   CustardFin Transcoder Benchmark (Milestone 9)" -ForegroundColor Cyan
Write-Host "================================================================================" -ForegroundColor Cyan

# 1. Dependency checks
if (-not (Get-Command ffmpeg -ErrorAction SilentlyContinue)) {
    Write-Host "Error: ffmpeg is not installed or not in PATH." -ForegroundColor Red
    exit 1
}

$ffmpegVer = (ffmpeg -version | Select-Object -First 1)
Write-Host "FFmpeg:   $ffmpegVer"

# 2. Hardware / GPU detection
$gpuName = "None detected"
$nvencAvailable = $false

if (Get-Command nvidia-smi -ErrorAction SilentlyContinue) {
    try {
        $gpuInfo = (nvidia-smi --query-gpu=name,driver_version --format=csv,noheader 2>$null)
        if ($gpuInfo) {
            $gpuName = $gpuInfo
        }
    } catch {}
}

& ffmpeg -v error -f lavfi -i testsrc=duration=1:size=256x256:rate=30 -pix_fmt yuv420p -c:v h264_nvenc -f null - 2>$null
if ($LASTEXITCODE -eq 0) {
    $nvencAvailable = $true
    Write-Host "GPU:      $gpuName [h264_nvenc operational]" -ForegroundColor Green
} else {
    Write-Host "GPU:      $gpuName [h264_nvenc unavailable]" -ForegroundColor Yellow
}

Write-Host "Clip:     ${Duration}s 1920x1080 @ 30fps"
Write-Host "Ladders:  1080p (4500k), 720p (2400k), 480p (1000k) | 4s HLS Segments"

$workDir = Join-Path $env:TEMP ("custardfin_bench_" + [System.Guid]::NewGuid().ToString().Substring(0,8))
New-Item -ItemType Directory -Force -Path $workDir | Out-Null

try {
    # 3. Generate synthetic source media
    $srcClip = Join-Path $workDir "input_1080p.mp4"
    Write-Host "`n[1/3] Generating synthetic 1080p test video (${Duration}s)..." -ForegroundColor Cyan
    & ffmpeg -y -v error -f lavfi -i "testsrc2=duration=${Duration}:size=1920x1080:rate=30" -pix_fmt yuv420p -c:v libx264 -preset ultrafast $srcClip

    function Run-HlsEncode {
        param (
            [string]$Encoder,
            [string]$OutDir
        )
        $extraFlags = @()
        if ($Encoder -eq "libx264") {
            $extraFlags = @("-preset", "veryfast")
        } elseif ($Encoder -eq "h264_nvenc") {
            $extraFlags = @("-preset", "p4", "-tune", "hq")
        }

        New-Item -ItemType Directory -Force -Path (Join-Path $OutDir "1080p") | Out-Null
        New-Item -ItemType Directory -Force -Path (Join-Path $OutDir "720p") | Out-Null
        New-Item -ItemType Directory -Force -Path (Join-Path $OutDir "480p") | Out-Null

        & ffmpeg -y -v error -i $srcClip -vf "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2" -pix_fmt yuv420p -c:v $Encoder @extraFlags -b:v 4500k -maxrate 4500k -bufsize 2M -force_key_frames "expr:gte(t,n_forced*4)" -c:a aac -b:a 160k -ac 2 -hls_time 4 -hls_list_size 0 -hls_segment_type mpegts -hls_segment_filename (Join-Path $OutDir "1080p\segment_%04d.ts") (Join-Path $OutDir "1080p\index.m3u8")
        & ffmpeg -y -v error -i $srcClip -vf "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2" -pix_fmt yuv420p -c:v $Encoder @extraFlags -b:v 2400k -maxrate 2400k -bufsize 2M -force_key_frames "expr:gte(t,n_forced*4)" -c:a aac -b:a 128k -ac 2 -hls_time 4 -hls_list_size 0 -hls_segment_type mpegts -hls_segment_filename (Join-Path $OutDir "720p\segment_%04d.ts") (Join-Path $OutDir "720p\index.m3u8")
        & ffmpeg -y -v error -i $srcClip -vf "scale=854:480:force_original_aspect_ratio=decrease,pad=854:480:(ow-iw)/2:(oh-ih)/2" -pix_fmt yuv420p -c:v $Encoder @extraFlags -b:v 1000k -maxrate 1000k -bufsize 2M -force_key_frames "expr:gte(t,n_forced*4)" -c:a aac -b:a 96k -ac 2 -hls_time 4 -hls_list_size 0 -hls_segment_type mpegts -hls_segment_filename (Join-Path $OutDir "480p\segment_%04d.ts") (Join-Path $OutDir "480p\index.m3u8")
    }

    # 4. CPU Benchmark
    Write-Host "[2/3] Running CPU Baseline Transcode (libx264)..." -ForegroundColor Cyan
    $cpuStart = [System.Diagnostics.Stopwatch]::StartNew()
    Run-HlsEncode -Encoder "libx264" -OutDir (Join-Path $workDir "cpu_out")
    $cpuStart.Stop()
    $cpuSec = [Math]::Round($cpuStart.Elapsed.TotalSeconds, 2)
    $totalContentSec = $Duration * 3
    $cpuSpeed = [Math]::Round($totalContentSec / $cpuSec, 2)
    Write-Host "CPU completed in: ${cpuSec}s (${cpuSpeed}x realtime)"

    # 5. GPU Benchmark
    $gpuSec = "N/A"
    $gpuSpeed = "N/A"
    $speedup = "N/A"
    $percentSaved = "N/A"

    if ($nvencAvailable) {
        Write-Host "[3/3] Running GPU Accelerated Transcode (h264_nvenc)..." -ForegroundColor Cyan
        $gpuStart = [System.Diagnostics.Stopwatch]::StartNew()
        Run-HlsEncode -Encoder "h264_nvenc" -OutDir (Join-Path $workDir "gpu_out")
        $gpuStart.Stop()
        $gpuSec = [Math]::Round($gpuStart.Elapsed.TotalSeconds, 2)
        $gpuSpeed = [Math]::Round($totalContentSec / $gpuSec, 2)
        $speedup = [Math]::Round($cpuSec / $gpuSec, 2)
        $percentSaved = [Math]::Round((($cpuSec - $gpuSec) / $cpuSec) * 100, 1)
        Write-Host "GPU completed in: ${gpuSec}s (${gpuSpeed}x realtime)"
    } else {
        Write-Host "[3/3] Skipping GPU benchmark: hardware acceleration unavailable on this host." -ForegroundColor Yellow
    }

    Write-Host "`n================================================================================" -ForegroundColor Cyan
    Write-Host "                     BENCHMARK COMPARISON RESULTS" -ForegroundColor Cyan
    Write-Host "================================================================================" -ForegroundColor Cyan
    "{0,-26} {1,-24} {2,-24}" -f "METRIC", "CPU (libx264)", "GPU (h264_nvenc)"
    Write-Host "--------------------------------------------------------------------------------"
    "{0,-26} {1,-24} {2,-24}" -f "Total Content Length", "${totalContentSec}s (3 renditions)", "${totalContentSec}s (3 renditions)"
    "{0,-26} {1,-24} {2,-24}" -f "Execution Wall Time", "${cpuSec}s", "${gpuSec}s"
    "{0,-26} {1,-24} {2,-24}" -f "Processing Rate", "${cpuSpeed}x realtime", "${gpuSpeed}x realtime"
    if ($nvencAvailable) {
        "{0,-26} {1,-24} {2,-24}" -f "Hardware Acceleration", "Baseline (1.0x)", "${speedup}x FASTER (${percentSaved}% saved)"
    } else {
        "{0,-26} {1,-24} {2,-24}" -f "Hardware Acceleration", "Baseline (1.0x)", "Hardware fallback active"
    }
    "{0,-26} {1,-24} {2,-24}" -f "Keyframe Alignment", "Strict 4.0s", "Strict 4.0s"
    Write-Host "================================================================================" -ForegroundColor Cyan
    Write-Host "Benchmark finished successfully." -ForegroundColor Green
} finally {
    Remove-Item -Recurse -Force $workDir -ErrorAction SilentlyContinue
}
