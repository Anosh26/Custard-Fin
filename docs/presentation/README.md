# CustardFin Presentation Deck

This directory contains the complete presentation assets for the **CustardFin** MCA Cloud Computing Lab project.

## Formats Available

| Format | File Path | Description |
|---|---|---|
| **PowerPoint (.pptx)** | [`CustardFin_Presentation.pptx`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/presentation/CustardFin_Presentation.pptx) | Modern 16:9 widescreen presentation with native **OpenXML slide transitions** and **shape entrance build animations**. Open directly in Microsoft PowerPoint, Keynote, or Google Slides. |
| **OpenXML (.xml)** | [`presentation.xml`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/presentation/presentation.xml) | Formatted OpenXML PresentationML structure containing full `<PresentationML>`, slide trees, `<p:transition>`, and `<p:timing>` animation sequences. |
| **Interactive Web Deck** | [`index.html`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/presentation/index.html) | Interactive HTML5 presentation deck with smooth 60fps animations, keyboard navigation (Arrow keys / Spacebar / F for Fullscreen), progress bar, and card hover effects. |

---

## Slide Breakdown (12 Slides)

1. **Slide 1: Title & Hero** &bull; CustardFin: Distributed Community Media Streaming Platform (Authors: Anish Ghadi & Anosh Chodankar).
2. **Slide 2: The Core Problem** &bull; Why traditional solutions fail (Jellyfin/Plex private silos vs YouTube centralized storage costs & privacy loss).
3. **Slide 3: Invariant 1** &bull; Zero Cloud Video Storage: Videos strictly live at home; syncing a ~3 KB JSON manifest.
4. **Slide 4: Two-Tier Architecture** &bull; Home Edge Go Daemon (<20 MB RAM, SQLite, NVENC) + Cloud Control Plane (Next.js 15, Go API, Postgres, Redis, S3).
5. **Slide 5: What It Uses** &bull; 100% Free and Open-Source stack breakdown across Edge, Cloud, Network, and Frontend.
6. **Slide 6: What It Does** &bull; Ingestion, FFprobe metadata, 4-second keyframe HLS transcoding, and 1-click manifest sync.
7. **Slide 7: Stream Router** &bull; Intelligent Play Decision Engine (Redis liveness, Direct LAN subnet bypass, Origin Shield cache hit/miss).
8. **Slide 8: Scalability & Limits** &bull; Invariant 7 (Origin Shield pull-through cache absorbs 90%+ viral traffic) & Invariant 10 (10-stream remote cap with HTTP 429 Retry-After).
9. **Slide 9: Hardware Acceleration** &bull; Empirical benchmarks on RTX 4060: NVENC hardware encoding at 12.4x real-time speed vs CPU libx264.
10. **Slide 10: Multi-Topology Deployment** &bull; Mode 1 (Pure Local), Mode 2 (Dev Containers), Mode 3 (Floci Cloud Emulator), and Mode 4 (GCP Live Cloud VM).
11. **Slide 11: Live Demonstration** &bull; 5-terminal synchronized orchestration launcher (`start-all-terminals.ps1`) and 6 core viva scenarios.
12. **Slide 12: Summary & Q&A** &bull; Project conclusion, architectural wins, engineering rigor, and viva evaluation readiness.

---

## Re-generating the Presentation

To re-generate or modify the PowerPoint and OpenXML files automatically:

```powershell
python scripts/generate_presentation.py
```
