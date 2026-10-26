# CustardFin Deployment Guide: GCP Compute Engine & Cloud Tier

This guide provides instructions for deploying the **CustardFin Cloud Tier** onto **Google Cloud Platform (GCP) Compute Engine** using free/promotional credits (e.g. $40 credits), along with complete network firewall configuration, environment variable alignment, and home edge client settings.

---

## Architecture Topology

```
+-------------------------------------------------------------+
| HOME EDGE NODE (Anish & Anosh's RTX Laptop / Home PC)       |
|  - Video Binaries Stay Here (Invariant 1)                   |
|  - Inbound Router Ports: NONE REQUIRED (ISP CGNAT)          |
|  - Outbound FRP Client -> Connects to GCP VM Port 7000      |
+-------------------------------------------------------------+
                               |
                               | Outbound TLS/TCP Tunnel
                               v
+-------------------------------------------------------------+
| GCP COMPUTE ENGINE VM (Cloud Tier)                          |
|  - Static External IP: [GCP_EXTERNAL_IP]                    |
|  - Docker Compose (Postgres, Redis, MinIO, FRPS, API, Web)  |
|  - Caddy Reverse Proxy (Ports 80 / 443)                     |
+-------------------------------------------------------------+
                               ^
                               | HTTP/HLS Playback
+-------------------------------------------------------------+
| WEB VIEWERS / PROFESSORS / EVALUATION DEVICES               |
|  - Browser connects to http://[GCP_EXTERNAL_IP]             |
|  - Chunks cached in MinIO Origin Shield (Invariant 7)       |
+-------------------------------------------------------------+
```

---

## 1. GCP VM Provisioning (Within $40 Budget)

1. Open the [GCP Console](https://console.cloud.google.com/) and navigate to **Compute Engine > VM Instances**.
2. Click **Create Instance**:
   - **Name**: `custardfin-cloud`
   - **Region**: Choose a low-latency or cost-effective region (e.g., `asia-south1` Mumbai/Delhi or `us-central1`).
   - **Machine Configuration**:
     - Series: **E2**
     - Machine type: **e2-medium** (2 vCPU, 4 GB RAM) — costs approximately ~$25/month (~$0.034/hour), easily covered by your $40 credit for several weeks of demonstration.
   - **Boot Disk**:
     - OS: **Ubuntu 24.04 LTS** or **Debian 12**
     - Size: **25 GB** (Standard Persistent Disk).
   - **Firewall**: Check **Allow HTTP traffic** and **Allow HTTPS traffic**.
   - Click **Create**.

3. **Reserve a Static External IP** (Recommended so IP does not change upon VM reboot):
   - Go to **VPC Network > IP addresses**.
   - Locate the ephemeral IP assigned to `custardfin-cloud` and click **Promote to static IP**.
   - Note down your `GCP_EXTERNAL_IP`.

---

## 2. GCP Firewall Rules Configuration

By default, GCP blocks all non-standard ingress ports. We must create firewall rules for the FRP tunnel and web services.

Go to **VPC Network > Firewall > Create Firewall Rule**:

### Rule A: CustardFin FRP Reverse Tunnel
- **Name**: `custardfin-allow-frp`
- **Network**: `default`
- **Direction of traffic**: `Ingress`
- **Action on match**: `Allow`
- **Targets**: `All instances in the network` (or apply target tag `custardfin-cloud`)
- **Source IPv4 ranges**: `0.0.0.0/0`
- **Protocols and ports**:
  - Specified protocols and ports:
    - `tcp:7000` (FRP control channel)
    - `tcp:8000-8020` (FRP vhost and streaming ports)
    - `tcp:7500` (FRP web dashboard, optional)

### Rule B: Web & API Ingress
- If not already permitted by default HTTP/HTTPS:
  - **Name**: `custardfin-allow-web`
  - **Protocols and ports**: `tcp:80`, `tcp:443`, `tcp:8080`, `tcp:9001` (MinIO console)
  - **Source IPv4 ranges**: `0.0.0.0/0`

### Summary of Ports Reference Table

| Port | Protocol | Direction | Purpose | Required Open on GCP VM? | Required Open on Home Router? |
|---|---|---|---|---|---|
| **80** | TCP | Ingress | HTTP (Caddy Proxy & Web UI) | **Yes** | **No** (0 inbound) |
| **443** | TCP | Ingress | HTTPS (Caddy auto-TLS) | **Yes** | **No** (0 inbound) |
| **7000** | TCP | Ingress | FRP Tunnel Server Port | **Yes** | **No** (Edge connects outbound) |
| **8000-8020** | TCP | Ingress | FRP Proxied Stream Ports | **Yes** | **No** |
| **8080** | TCP | Ingress | Cloud API direct port | Optional (proxied via Caddy) | **No** |
| **9001** | TCP | Ingress | MinIO Web Console | Optional (for demo dashboard) | **No** |
| **5432** | TCP | Internal | PostgreSQL | No (Docker internal only) | **No** |
| **6379** | TCP | Internal | Redis | No (Docker internal only) | **No** |

---

## 3. Server Setup on GCP VM

SSH into your GCP VM (via the GCP browser SSH button or `gcloud compute ssh custardfin-cloud`):

### 3.1 Install Docker & Docker Compose
```bash
sudo apt update && sudo apt install -y curl git ca-certificates

# Install Docker engine
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Log out and log back in, or run:
newgrp docker

# Verify Docker version
docker --version
docker compose version
```

### 3.2 Clone Repository
```bash
git clone https://github.com/anosh26/custardfin.git
cd custardfin
```

### 3.3 Configure Environment Variables for GCP
Create `deploy/.env` based on `deploy/env.example`:

```bash
cp deploy/env.example deploy/.env
```

Edit `deploy/.env` (using `nano deploy/.env`):
```ini
# Replace <GCP_EXTERNAL_IP> with your VM's public IP
NEXT_PUBLIC_API_URL=http://<GCP_EXTERNAL_IP>/v1

# S3 Origin Cache Configuration
S3_ENDPOINT=http://s3-cache:9000
S3_BUCKET=custardfin-origin-cache
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin

# FRP Tunnel Relay Configuration
FRP_SERVER_ADDR=frps
FRP_SERVER_PORT=7000
FRP_AUTH_TOKEN=custardfin_demo_frp_token

# Cloud API Security
API_ADDR=:8080
JWT_SECRET=custardfin_super_secure_gcp_production_key_32_bytes!
```

### 3.4 Launch Cloud Tier Services
```bash
docker compose -f deploy/compose.demo.yml up -d
```

Check service health:
```bash
docker compose -f deploy/compose.demo.yml ps
curl -i http://localhost/healthz
```
Expected output: `HTTP/1.1 200 OK` with JSON `{"status":"ok","timestamp":...}`.

---

## 4. Edge Client Configuration (Home RTX Machine)

On your local home PC/laptop where your media library is located, configure the edge daemon to connect **outbound** to your GCP instance:

### 4.1 Update Local `.env`
Set the following environment variables for `edged`:

```ini
# Address of your newly deployed GCP Cloud VM (port 8080 for API, or port 80 if Caddy is running)
EDGE_CLOUD_URL=http://<GCP_EXTERNAL_IP>:8080

# Enable FRP Reverse Tunnel for remote stream playback
EDGE_FRP_ENABLED=true
EDGE_FRP_SERVER_ADDR=<GCP_EXTERNAL_IP>
EDGE_FRP_SERVER_PORT=7000
EDGE_FRP_AUTH_TOKEN=custardfin_demo_frp_token

# Local library paths
EDGE_MEDIA_DIR=./media
EDGE_TRANSCODE_DIR=./transcoded
EDGE_DB_PATH=./data/edge.db
```

Or export them directly in your shell:
```bash
export PATH="$HOME/.local/bin:$PATH"
export EDGE_CLOUD_URL="http://<GCP_EXTERNAL_IP>:8080"
export EDGE_FRP_ENABLED="true"
export EDGE_FRP_SERVER_ADDR="<GCP_EXTERNAL_IP>"
export EDGE_FRP_SERVER_PORT="7000"
export EDGE_FRP_AUTH_TOKEN="custardfin_demo_frp_token"

go run ./apps/edge/cmd/edged
```

### 4.2 Run the Edge Node
```bash
# Seed demo video fixtures if not already present
make seed

# Run edge daemon
go run ./apps/edge/cmd/edged
```

### 4.3 Verify Connection to GCP Cloud
In the edge logs, verify:
1. `registered node with cloud api` -> `HTTP 201 Created`
2. `started heartbeat worker` -> sent every 15s to `http://<GCP_EXTERNAL_IP>/v1/nodes/{id}/heartbeat`
3. `frp reverse tunnel established` -> outbound link to `<GCP_EXTERNAL_IP>:7000` active.

---

## 5. Verification Checklist

1. **Browser Test**: Open `http://<GCP_EXTERNAL_IP>` in your phone or another laptop disconnected from home Wi-Fi (e.g. mobile data).
2. **Catalog Visibility**: Verify that videos published from home edge node appear in the public catalog.
3. **Remote Playback**: Click Play; verify stream starts smoothly over the FRP tunnel without needing port-forwarding on home router.
4. **Origin Shield Cache**: Check MinIO console at `http://<GCP_EXTERNAL_IP>:9001` (user: `minioadmin`, pass: `minioadmin`) to observe cached `.ts` chunks being written.
