#!/usr/bin/env bash
# ==============================================================================
# CustardFin Home Edge Node Launcher (Interactive Cloud IP Selector)
# Prompts for the GCP / Cloud VM External IP on every run, with smart default
# caching so you can simply press Enter if the IP hasn't changed.
# ==============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
IP_CACHE_FILE="${ROOT_DIR}/.last_cloud_ip"

cd "${ROOT_DIR}"

# 1. Ensure user's local bin is in PATH for frpc
export PATH="$HOME/.local/bin:$PATH"

# 2. Check for frpc binary
if ! command -v frpc &> /dev/null; then
    echo "⚠️  frpc binary not found in PATH."
    echo "   Ensure frpc is installed (e.g. in ~/.local/bin/frpc)."
fi

# 3. Read previous IP if saved
DEFAULT_IP="136.71.107.75"
if [ -f "${IP_CACHE_FILE}" ]; then
    CACHED_IP="$(cat "${IP_CACHE_FILE}" 2>/dev/null | tr -d '[:space:]')"
    if [ -n "${CACHED_IP}" ]; then
        DEFAULT_IP="${CACHED_IP}"
    fi
fi

echo "================================================================="
echo "  🍮 CustardFin Edge Node Launcher (RTX NVENC / Hybrid Edge)"
echo "================================================================="
echo ""
read -r -p "Enter Cloud VM External IP [Press Enter for ${DEFAULT_IP}]: " INPUT_IP

CLOUD_IP="${INPUT_IP:-$DEFAULT_IP}"

# Sanitize input (strip protocol http://, https://, trailing slashes, and ports)
CLOUD_IP="${CLOUD_IP#http://}"
CLOUD_IP="${CLOUD_IP#https://}"
CLOUD_IP="${CLOUD_IP%%/*}"
CLOUD_IP="${CLOUD_IP%%:*}"
CLOUD_IP="$(echo "${CLOUD_IP}" | tr -d '[:space:]')"

if [ -z "${CLOUD_IP}" ]; then
    echo "❌ Error: Cloud IP cannot be empty."
    exit 1
fi

# Cache the IP for next run
echo "${CLOUD_IP}" > "${IP_CACHE_FILE}"

echo ""
echo "🔗 Connecting Home Edge Daemon to Cloud VM:"
echo "   • Cloud Control API : http://${CLOUD_IP}:8080"
echo "   • FRP Reverse Tunnel: ${CLOUD_IP}:7000"
echo "   • Local Stream Port : 8090"
echo ""

# Pre-flight check: probe cloud API healthz
if command -v curl &> /dev/null; then
    echo "⏳ Probing Cloud API at http://${CLOUD_IP}:8080/healthz..."
    if curl -s -m 3 "http://${CLOUD_IP}:8080/healthz" | grep -q "ok"; then
        echo "   ✅ Cloud API is active and healthy!"
    else
        echo "   ⚠️  Cloud API did not respond immediately."
        echo "      (If the VM was just booted, containers may take ~10s to warm up)."
    fi
    echo ""
fi

# 4. Set Environment Variables
export EDGE_CLOUD_URL="http://${CLOUD_IP}:8080"
export EDGE_FRP_ENABLED="true"
export EDGE_FRP_SERVER_ADDR="${CLOUD_IP}"
export EDGE_FRP_SERVER_PORT="7000"
export EDGE_FRP_AUTH_TOKEN="custardfin_demo_frp_token"
export EDGE_ADDR=":8090"

# 5. Start Edge Daemon
echo "🚀 Starting CustardFin edged daemon..."
exec go run ./apps/edge/cmd/edged
