#!/usr/bin/env bash
# ==============================================================================
# CustardFin Media Seeder Script
# Prepares open-license test media fixtures for edge scanner and transcoding.
# ==============================================================================

set -euo pipefail

TARGET_DIR="${1:-./media}"
mkdir -p "$TARGET_DIR"

echo "================================================================================"
echo "                   CustardFin Media Library Seeder"
echo "================================================================================"
echo "Target directory: $TARGET_DIR"

if ! command -v ffmpeg >/dev/null 2>&1; then
    echo "Error: ffmpeg is required to synthesize media fixtures."
    exit 1
fi

# 1. Synthesize 1080p clip with audio
CLIP_1080="${TARGET_DIR}/cosmos_1080p.mp4"
if [ ! -f "$CLIP_1080" ]; then
    echo "[1/3] Generating 1080p media fixture (15s)..."
    ffmpeg -y -v error \
        -f lavfi -i "testsrc2=duration=15:size=1920x1080:rate=30" \
        -f lavfi -i "sine=frequency=1000:duration=15" \
        -pix_fmt yuv420p -c:v libx264 -preset fast \
        -c:a aac -b:a 160k -shortest \
        "$CLIP_1080"
else
    echo "[1/3] 1080p fixture exists: $CLIP_1080"
fi

# 2. Synthesize 720p clip with audio
CLIP_720="${TARGET_DIR}/trailer_720p.mp4"
if [ ! -f "$CLIP_720" ]; then
    echo "[2/3] Generating 720p media fixture (20s)..."
    ffmpeg -y -v error \
        -f lavfi -i "testsrc2=duration=20:size=1280x720:rate=30" \
        -f lavfi -i "sine=frequency=440:duration=20" \
        -pix_fmt yuv420p -c:v libx264 -preset fast \
        -c:a aac -b:a 128k -shortest \
        "$CLIP_720"
else
    echo "[2/3] 720p fixture exists: $CLIP_720"
fi

# 3. Synthesize 480p clip with audio
CLIP_480="${TARGET_DIR}/vintage_480p.mp4"
if [ ! -f "$CLIP_480" ]; then
    echo "[3/3] Generating 480p media fixture (30s)..."
    ffmpeg -y -v error \
        -f lavfi -i "testsrc2=duration=30:size=854x480:rate=30" \
        -f lavfi -i "sine=frequency=880:duration=30" \
        -pix_fmt yuv420p -c:v libx264 -preset fast \
        -c:a aac -b:a 96k -shortest \
        "$CLIP_480"
else
    echo "[3/3] 480p fixture exists: $CLIP_480"
fi

echo "--------------------------------------------------------------------------------"
echo "Media library successfully seeded in ${TARGET_DIR}:"
ls -lh "$TARGET_DIR"/*.mp4 2>/dev/null || true
echo "================================================================================"
