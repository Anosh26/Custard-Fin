# CustardFin Presentation Deck

Comprehensive 8-slide presentation suite for the MCA Cloud Computing Lab project **CustardFin**.

- **Authors:** Anish Ghadi (25P0320009) & Anosh Chodankar (25P0320010)
- **Institution:** MCA Cloud Computing Lab, Semester End Evaluation 2026
- **Architecture:** Distributed Community Video Streaming (Zero Permanent Cloud Video Storage)

---

## Available Formats

| Format | File Path | Description |
|---|---|---|
| **PowerPoint (.pptx)** | [`CustardFin_Presentation.pptx`](CustardFin_Presentation.pptx) | Modern 16:9 widescreen presentation with native **OpenXML slide transitions**, shape entrance build animations, high-contrast pastel cards, and Custard mascot badges. Open directly in PowerPoint, Keynote, or Google Slides. |
| **OpenXML (.xml)** | [`presentation.xml`](presentation.xml) | Formatted OpenXML PresentationML structure containing the slide tree, `<p:transition>`, and `<p:timing>` animation sequences. |
| **Interactive Web Deck** | [`index.html`](index.html) | Interactive HTML5 presentation deck with smooth 60fps animations, keyboard navigation (Arrow keys / Spacebar / F for Fullscreen), progress tracking, inline SVG diagrams, and 8 unique per-slide Custard mascot animations. |

---

## 8-Slide Structure Breakdown

1. **Slide 1: What It Is (Hero & Vision)**
   - *Title:* `CUSTARDFIN: Community Video Platform`
   - *Core Message:* Keep your videos at home like Jellyfin, share them with the world like YouTube — with Zero Permanent Cloud Video Storage.
   - *Cards:* Team Credits (Anish Ghadi & Anosh Chodankar), Core Rule (Videos stay 100% on home PC, ~3 KB text sync), Simple Tech Stack (Go 1.24, Next.js 15, Postgres/SQLite/Redis, GPU + CPU encoding).
   - *Mascot Speech:* `"Welcome to CustardFin! ✨"`

2. **Slide 2: Why Current Apps Fail (The Core Problem)**
   - *Title:* `Why Current Apps Fail: Big Cloud vs Home Servers`
   - *Cards:* Big Cloud pitfalls (monthly storage bills, loss of data ownership, high bandwidth fees) vs Home Server bottlenecks (ISP CGNAT blockage, lack of public discovery, 3–4 outside viewers choking home Wi-Fi upload).
   - *Bottom Banner:* The CustardFin Fix (Videos stay at home • Cloud only stores ~3 KB info file • Cloud cache protects home Wi-Fi).
   - *Mascot Speech:* `"Best of Both Worlds! 💡"`

3. **Slide 3: How It Works (4-Step End-to-End Flow)**
   - *Title:* `How It Works: From Home Folder to Public Stream`
   - *Pipeline Steps:*
     - **Step 01 (Home):** Scan & Hash (`SHA-256` unique ID, thumbnail preview, saved in SQLite).
     - **Step 02 (Encode):** 4s Video Chunks (`4s HLS` `.ts` segments via GPU or CPU, zero lag on play).
     - **Step 03 (Connect):** Safe Outbound Tunnel (`OUTBOUND` FRP link to GCP, bypasses CGNAT, no open router ports, 15s ping).
     - **Step 04 (Share):** 3 KB Info Sync (`~3 KB` JSON manifest synced to cloud Postgres on 1-click Public toggle).
   - *Mascot Speech:* `"4 Simple Steps! ⚙️"`

4. **Slide 4: System Architecture (Two-Part Design)**
   - *Title:* `Two-Part Design: Home PC + Cloud Server`
   - *Card 1:* Part 1 • Home Edge App (Runs on any Home PC) — Tiny Go app (`< 20 MB RAM`), local SQLite WAL DB, FFmpeg worker, FRP tunnel client, direct gigabit LAN streaming.
   - *Card 2:* Part 2 • Cloud Control Server (Runs on Google Cloud GCP) — Next.js 15 web UI, Go Cloud API, PostgreSQL + Redis, 24h Origin Shield chunk cache, Docker Compose stack (`0 Permanent Videos`).
   - *Mascot Speech:* `"Home & Cloud Together! 🤝"`

5. **Slide 5: How Cloud (GCP) Is Used**
   - *Title:* `How We Use the Cloud: GCP Features & Setup`
   - *Cards:*
     - **GCP Compute Engine VM:** Linux host running Docker Compose containers with static public IP.
     - **GCP VPC & Firewall:** Ports 80/443 for web traffic, port 7000 for FRP tunnel, Caddy auto-HTTPS.
     - **In-VM S3 Object Cache:** S3-compatible MinIO container running inside the single GCP VM caching active 4s chunks and auto-deleting after 24 hours (`24h LRU`), with 0 extra cloud bucket fees.
     - **Cloud State:** PostgreSQL for ~3 KB manifests/users + Redis for 15-second live node health checks (`15s TTL`).
   - *Mascot Speech:* `"Powered by GCP! ☁️"`

6. **Slide 6: Smart Traffic Routing & Origin Shield**
   - *Title:* `Origin Shield & Router: Protecting Home Wi-Fi`
   - *Card 1:* Smart Play Router (4-Step Decision Check: Redis 15s online check, Same Wi-Fi LAN bypass, Cloud cache hit, Cache miss pull-through via tunnel capped at 10 remote streams).
   - *Card 2:* Home Wi-Fi Protection (`92.4% SAVED`) — Vector flow diagram showing `Home PC (1x Upload / 5 Mbps) -> GCP Shield (24h Cache) -> 100 Viewers (500 Mbps Out)`.
   - *Mascot Speech:* `"Home Wi-Fi Safe! 🛡️"`

7. **Slide 7: Open-Source Tech Stack & Universal Hardware**
   - *Title:* `Open-Source Tools & Universal Hardware Support`
   - *Cards:*
     - **Smart Encoding (`GPU / CPU`):** Auto-detects hardware on startup; uses NVIDIA NVENC GPU if available or falls back to CPU (`libx264`) so it runs uniformly on all team laptops.
     - **Low Memory (`< 18 MB`):** Tiny static Go binary uses under 18 MB idle RAM compared to 400+ MB on Jellyfin/Plex (95% less RAM).
     - **100% Free & Open Source (`100% FOSS`):** Go 1.24, Next.js 15, React 19, PostgreSQL, SQLite, Redis, FFmpeg, FRP, Caddy.
   - *Mascot Speech:* `"Runs on Any PC! 🚀"`

8. **Slide 8: Live Demo & Summary**
   - *Title:* `Live Demonstration: What We Show Today`
   - *Card 1:* Live Demo Checklist — User signup/login, Private vs Public toggle, auto video thumbnails, self-cleaning library on startup, LAN bypass + remote cloud streaming, 15s instant offline detection.
   - *Card 2:* Project Summary (`READY`) — Zero permanent cloud video storage, smart GCP control layer, universal GPU/CPU support, and Q&A / Viva Defense sign-off for Anish Ghadi & Anosh Chodankar.
   - *Mascot Speech:* `"Thank You! Questions? 🏆"`

---

## Rebuilding the Presentation

To re-generate the PowerPoint presentation (`CustardFin_Presentation.pptx`) and re-extract the OpenXML representation (`presentation.xml`):

```powershell
python scripts/generate_presentation.py
```
