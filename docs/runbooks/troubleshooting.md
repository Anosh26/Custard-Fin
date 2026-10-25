# CustardFin Troubleshooting & Emergency Recovery Runbook

This runbook provides diagnostic commands and quick resolutions for potential issues encountered during development, local testing, or live demonstration.

---

## 1. Port Conflicts (Address Already in Use)

### Symptom
Error message: `bind: address already in use` or Docker container fails to start because a port is occupied.

### Common Port Assignments
- **5432**: PostgreSQL
- **6379**: Redis
- **7000**: FRP Server Control Port
- **8080**: CustardFin Cloud API (`apid`)
- **8090**: CustardFin Edge Daemon (`edged`)
- **9000 / 9001**: MinIO S3 API & Web Console
- **3000**: Next.js Web Frontend
- **80 / 443**: Caddy Reverse Proxy

### Resolution
1. Identify the process occupying the port:
   ```bash
   # Linux / macOS
   sudo lsof -i :<PORT> || sudo netstat -nlp | grep :<PORT>

   # Windows (PowerShell)
   Get-NetTCPConnection -LocalPort <PORT>
   ```
2. Terminate the conflicting process or change the port mapping in `.env`.
3. If an orphaned Docker container is holding the port:
   ```bash
   docker ps -a
   docker rm -f $(docker ps -aq)
   ```

---

## 2. FRP Reverse Tunnel Connection Failures

### Symptom
Edge daemon logs: `failed to connect to frp server: dial tcp: connection refused` or `token mismatch`.

### Diagnosis & Fix
1. **Check FRPS container status**:
   ```bash
   docker logs custardfin-frps
   ```
2. **Verify Auth Token Match**:
   Ensure `auth.token` in `deploy/frps.toml` exactly matches `EDGE_FRP_AUTH_TOKEN` in edge configuration:
   - Default: `custardfin_demo_frp_token`
3. **Check Firewall Rules (GCP Deployment)**:
   Ensure TCP port `7000` and ports `8000-8020` are permitted by the GCP VPC firewall (see `docs/runbooks/deploy.md`).
4. **Local Fallback**:
   If testing strictly on localhost, set:
   ```ini
   EDGE_FRP_SERVER_ADDR=localhost
   EDGE_FRP_SERVER_PORT=7000
   ```

---

## 3. Hardware Transcoding & NVENC Fallback

### Symptom
FFmpeg reports `Cannot load libcuda.so.1` or `No NVENC capable devices found`.

### Diagnosis
Run the benchmark script to check hardware detection:
```bash
make bench
```

### Resolution
- CustardFin's managed transcoder automatically detects hardware capabilities (`h264_nvenc`, `h264_qsv`, `h264_vaapi`) with a 1-second test probe.
- If no supported GPU or driver is detected, CustardFin automatically falls back to **CPU encoding (`libx264`)** without error (Invariant 6).
- If forcing CPU encoding is desired:
  ```ini
  EDGE_TRANSCODE_PRESET=cpu
  ```

---

## 4. Node Heartbeat & "Node Offline" Status

### Symptom
Video catalog shows `NODE_OFFLINE` or playback returns `HTTP 503`.

### Diagnosis & Fix
1. Check Redis heartbeat status:
   ```bash
   # Connect to Redis CLI
   docker exec -it custardfin-redis redis-cli

   # List active heartbeat keys
   KEYS heartbeat:*

   # Inspect TTL for your node ID
   TTL heartbeat:<NODE_ID>
   ```
2. If TTL is -2 (key expired), verify `edged` is running and able to send POST requests to `/v1/nodes/{id}/heartbeat`.
3. Verify network latency: edge node sends heartbeats every 15s; TTL default is 30s.

---

## 5. Clean Slate Reset (Demo Day Recovery)

If database state or test artifacts become corrupted, perform a 100% clean reset:

```bash
# 1. Stop all containers and remove volumes
make demo-down
docker volume prune -f

# 2. Clean temporary transcode directories
rm -rf ./transcoded/* ./data/edge.db*

# 3. Re-seed demo media fixtures
make seed

# 4. Bring up fresh demo stack
make demo-up

# 5. Verify health
make demo-check
```
