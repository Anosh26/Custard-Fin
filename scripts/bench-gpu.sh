#!/usr/bin/env bash
# ==============================================================================
# CustardFin GPU vs CPU Transcoder Benchmark Script (Milestone 9)
# Measures encoding speed, frames per second, and speedup ratio for HLS chunks.
# ==============================================================================

set -euo pipefail

# ANSI color codes
BOLD="\033[1m"
GREEN="\033[32m"
CYAN="\033[36m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

DURATION=${1:-15} # Default test duration: 15 seconds
WORKDIR=$(mktemp -d /tmp/custardfin_bench_XXXXXX)
cleanup() {
    rm -rf "$WORKDIR"
}
trap cleanup EXIT

echo -e "${BOLD}${CYAN}================================================================================"
echo -e "                   CustardFin Transcoder Benchmark (Milestone 9)"
echo -e "================================================================================${RESET}"

# 1. Dependency checks
if ! command -v ffmpeg >/dev/null 2>&1; then
    echo -e "${RED}Error: ffmpeg is not installed or not in PATH.${RESET}"
    exit 1
fi

FFMPEG_VER=$(ffmpeg -version | head -n 1)
echo -e "FFmpeg:   ${FFMPEG_VER}"

# 2. Hardware / GPU detection
GPU_NAME="None detected"
NVENC_AVAILABLE=false
if command -v nvidia-smi >/dev/null 2>&1; then
    GPU_INFO=$(nvidia-smi --query-gpu=name,driver_version --format=csv,noheader 2>/dev/null || true)
    if [ -n "$GPU_INFO" ]; then
        GPU_NAME="$GPU_INFO"
    fi
fi

if ffmpeg -v error -f lavfi -i testsrc=duration=1:size=256x256:rate=30 -pix_fmt yuv420p -c:v h264_nvenc -f null - >/dev/null 2>&1; then
    NVENC_AVAILABLE=true
    echo -e "GPU:      ${GREEN}${GPU_NAME} [h264_nvenc operational]${RESET}"
else
    echo -e "GPU:      ${YELLOW}${GPU_NAME} [h264_nvenc unavailable]${RESET}"
fi

echo -e "Clip:     ${BOLD}${DURATION}s 1920x1080 @ 30fps${RESET}"
echo -e "Ladders:  1080p (4500k), 720p (2400k), 480p (1000k) | 4s HLS Segments"
echo -e "Workdir:  ${WORKDIR}"
echo ""

# 3. Generate synthetic source media
SRC_CLIP="${WORKDIR}/input_1080p.mp4"
echo -e "${CYAN}[1/3] Generating synthetic 1080p test video (${DURATION}s)...${RESET}"
ffmpeg -y -v error \
    -f lavfi -i "testsrc2=duration=${DURATION}:size=1920x1080:rate=30" \
    -pix_fmt yuv420p \
    -c:v libx264 -preset ultrafast \
    "${SRC_CLIP}"

# Helper function to run HLS transcode of 3 renditions
run_hls_encode() {
    local encoder=$1
    local outdir=$2
    local extra_flags=()

    if [ "$encoder" = "libx264" ]; then
        extra_flags=("-preset" "veryfast")
    elif [ "$encoder" = "h264_nvenc" ]; then
        extra_flags=("-preset" "p4" "-tune" "hq")
    fi

    mkdir -p "${outdir}/1080p" "${outdir}/720p" "${outdir}/480p"

    # 1080p
    ffmpeg -y -v error -i "${SRC_CLIP}" \
        -vf "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2" \
        -pix_fmt yuv420p -c:v "${encoder}" "${extra_flags[@]}" \
        -b:v 4500k -maxrate 4500k -bufsize 2M \
        -force_key_frames "expr:gte(t,n_forced*4)" \
        -c:a aac -b:a 160k -ac 2 \
        -hls_time 4 -hls_list_size 0 -hls_segment_type mpegts \
        -hls_segment_filename "${outdir}/1080p/segment_%04d.ts" \
        "${outdir}/1080p/index.m3u8"

    # 720p
    ffmpeg -y -v error -i "${SRC_CLIP}" \
        -vf "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2" \
        -pix_fmt yuv420p -c:v "${encoder}" "${extra_flags[@]}" \
        -b:v 2400k -maxrate 2400k -bufsize 2M \
        -force_key_frames "expr:gte(t,n_forced*4)" \
        -c:a aac -b:a 128k -ac 2 \
        -hls_time 4 -hls_list_size 0 -hls_segment_type mpegts \
        -hls_segment_filename "${outdir}/720p/segment_%04d.ts" \
        "${outdir}/720p/index.m3u8"

    # 480p
    ffmpeg -y -v error -i "${SRC_CLIP}" \
        -vf "scale=854:480:force_original_aspect_ratio=decrease,pad=854:480:(ow-iw)/2:(oh-ih)/2" \
        -pix_fmt yuv420p -c:v "${encoder}" "${extra_flags[@]}" \
        -b:v 1000k -maxrate 1000k -bufsize 2M \
        -force_key_frames "expr:gte(t,n_forced*4)" \
        -c:a aac -b:a 96k -ac 2 \
        -hls_time 4 -hls_list_size 0 -hls_segment_type mpegts \
        -hls_segment_filename "${outdir}/480p/segment_%04d.ts" \
        "${outdir}/480p/index.m3u8"
}

# 4. CPU Benchmark
echo -e "${CYAN}[2/3] Running CPU Baseline Transcode (libx264)...${RESET}"
START_CPU=$(date +%s%N)
run_hls_encode "libx264" "${WORKDIR}/cpu_out"
END_CPU=$(date +%s%N)
CPU_NANOS=$((END_CPU - START_CPU))
CPU_SEC=$(awk "BEGIN {print ${CPU_NANOS}/1000000000}")
TOTAL_CONTENT_SEC=$((DURATION * 3))
CPU_SPEED=$(awk "BEGIN {printf \"%.2f\", ${TOTAL_CONTENT_SEC}/${CPU_SEC}}")
echo -e "CPU completed in: ${BOLD}${CPU_SEC}s${RESET} (${CPU_SPEED}x realtime)"

# 5. GPU Benchmark (if operational)
GPU_SEC="N/A"
GPU_SPEED="N/A"
SPEEDUP="N/A"
PERCENT_SAVED="N/A"

if [ "$NVENC_AVAILABLE" = true ]; then
    echo -e "${CYAN}[3/3] Running GPU Accelerated Transcode (h264_nvenc)...${RESET}"
    START_GPU=$(date +%s%N)
    run_hls_encode "h264_nvenc" "${WORKDIR}/gpu_out"
    END_GPU=$(date +%s%N)
    GPU_NANOS=$((END_GPU - START_GPU))
    GPU_SEC=$(awk "BEGIN {print ${GPU_NANOS}/1000000000}")
    GPU_SPEED=$(awk "BEGIN {printf \"%.2f\", ${TOTAL_CONTENT_SEC}/${GPU_SEC}}")
    SPEEDUP=$(awk "BEGIN {printf \"%.2f\", ${CPU_SEC}/${GPU_SEC}}")
    PERCENT_SAVED=$(awk "BEGIN {printf \"%.1f\", ((${CPU_SEC} - ${GPU_SEC}) / ${CPU_SEC}) * 100}")
    echo -e "GPU completed in: ${BOLD}${GPU_SEC}s${RESET} (${GPU_SPEED}x realtime)"
else
    echo -e "${YELLOW}[3/3] Skipping GPU benchmark: hardware acceleration unavailable on this host.${RESET}"
fi

echo ""
echo -e "${BOLD}${CYAN}================================================================================"
echo -e "                     BENCHMARK COMPARISON RESULTS"
echo -e "================================================================================${RESET}"
printf "%-26s %-24s %-24s\n" "METRIC" "CPU (libx264)" "GPU (h264_nvenc)"
echo "--------------------------------------------------------------------------------"
printf "%-26s %-24s %-24s\n" "Total Content Length" "${TOTAL_CONTENT_SEC}s (3 renditions)" "${TOTAL_CONTENT_SEC}s (3 renditions)"
printf "%-26s %-24s %-24s\n" "Execution Wall Time" "${CPU_SEC}s" "${GPU_SEC}s"
printf "%-26s %-24s %-24s\n" "Processing Rate" "${CPU_SPEED}x realtime" "${GPU_SPEED}x realtime"
if [ "$NVENC_AVAILABLE" = true ]; then
    echo -e "$(printf "%-26s %-24s" "Hardware Acceleration" "Baseline (1.0x)") ${GREEN}${SPEEDUP}x FASTER (${PERCENT_SAVED}% time saved)${RESET}"
else
    echo -e "$(printf "%-26s %-24s" "Hardware Acceleration" "Baseline (1.0x)") ${YELLOW}Hardware fallback active${RESET}"
fi
printf "%-26s %-24s %-24s\n" "Keyframe Alignment" "Strict 4.0s" "Strict 4.0s"
echo "================================================================================"
echo -e "${GREEN}Benchmark finished successfully.${RESET}"
