import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

# 1. HTML for Community Catalog Feed
catalog_html = """<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>CustardFin - Community Catalog</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
  body { background: #0b0f19; color: #f3f4f6; padding: 0; }
  header { background: #111827; border-bottom: 1px solid #1f2937; padding: 14px 32px; display: flex; align-items: center; justify-content: space-between; }
  .logo-box { display: flex; align-items: center; gap: 12px; }
  .logo-icon { width: 36px; height: 36px; background: linear-gradient(135deg, #eab308, #ca8a04); border-radius: 8px; display: flex; align-items: center; justify-content: center; font-weight: 700; color: #000; font-size: 20px; }
  .logo-text { font-size: 20px; font-weight: 700; color: #fff; letter-spacing: -0.5px; }
  .logo-sub { font-size: 11px; color: #eab308; text-transform: uppercase; letter-spacing: 1.5px; font-weight: 600; }
  .search-bar { flex: 1; max-width: 540px; margin: 0 40px; position: relative; }
  .search-input { width: 100%; background: #1f2937; border: 1px solid #374151; padding: 10px 16px 10px 40px; border-radius: 20px; color: #fff; font-size: 14px; outline: none; }
  .search-icon { position: absolute; left: 14px; top: 11px; color: #9ca3af; font-size: 14px; }
  .nav-actions { display: flex; align-items: center; gap: 14px; }
  .btn { padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; transition: 0.2s; }
  .btn-gold { background: #eab308; color: #000; }
  .btn-outline { background: transparent; border: 1px solid #374151; color: #d1d5db; }
  .badge-user { background: #1f2937; padding: 6px 12px; border-radius: 20px; border: 1px solid #374151; font-size: 13px; display: flex; align-items: center; gap: 8px; }
  .user-dot { width: 8px; height: 8px; border-radius: 50%; background: #10b981; }

  main { padding: 32px; max-width: 1400px; margin: 0 auto; }
  .hero-banner { background: linear-gradient(135deg, #1e1b4b 0%, #0f172a 100%); border: 1px solid #312e81; border-radius: 14px; padding: 24px 32px; margin-bottom: 32px; display: flex; justify-content: space-between; align-items: center; }
  .hero-title { font-size: 22px; font-weight: 700; color: #fff; margin-bottom: 6px; }
  .hero-desc { color: #94a3b8; font-size: 14px; max-width: 700px; line-height: 1.5; }
  .banner-stats { display: flex; gap: 24px; text-align: center; }
  .stat-val { font-size: 22px; font-weight: 700; color: #eab308; }
  .stat-lbl { font-size: 11px; color: #94a3b8; text-transform: uppercase; font-weight: 600; }

  .section-hdr { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
  .section-title { font-size: 18px; font-weight: 700; }
  .grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 24px; }
  .card { background: #131b2e; border: 1px solid #1f2937; border-radius: 12px; overflow: hidden; transition: transform 0.2s, box-shadow 0.2s; position: relative; }
  .thumb-box { position: relative; width: 100%; height: 180px; background: #1e293b; overflow: hidden; display: flex; align-items: center; justify-content: center; }
  .thumb-img { width: 100%; height: 100%; object-fit: cover; }
  .duration-badge { position: absolute; bottom: 8px; right: 8px; background: rgba(0,0,0,0.85); color: #fff; font-size: 11px; font-weight: 600; padding: 2px 6px; border-radius: 4px; }
  .badge-origin { position: absolute; top: 8px; left: 8px; background: rgba(16, 185, 129, 0.9); color: #fff; font-size: 10px; font-weight: 700; padding: 3px 8px; border-radius: 4px; text-transform: uppercase; letter-spacing: 0.5px; }
  .card-body { padding: 14px; }
  .card-title { font-size: 14px; font-weight: 600; color: #fff; margin-bottom: 6px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .card-meta { display: flex; justify-content: space-between; font-size: 12px; color: #94a3b8; margin-bottom: 8px; }
  .uploader-pill { font-size: 11px; background: #1f2937; padding: 2px 8px; border-radius: 12px; color: #cbd5e1; display: inline-flex; align-items: center; gap: 4px; }
</style>
</head>
<body>
<header>
  <div class="logo-box">
    <div class="logo-icon">&#9658;</div>
    <div>
      <div class="logo-text">CustardFin</div>
      <div class="logo-sub">Community Cloud</div>
    </div>
  </div>
  <div class="search-bar">
    <span class="search-icon">&#128269;</span>
    <input type="text" class="search-input" placeholder="Search distributed catalog across all home nodes...">
  </div>
  <div class="nav-actions">
    <button class="btn btn-outline" style="border-color: #6366f1; color: #a5b4fc;">&#128187; My Edge Library</button>
    <button class="btn btn-outline">&#9881; Dev Settings</button>
    <div class="badge-user">
      <div class="user-dot"></div>
      <span>Anosh (Node: RTX-4060)</span>
    </div>
  </div>
</header>

<main>
  <div class="hero-banner">
    <div>
      <div class="hero-title">Distributed Community Media Exchange</div>
      <div class="hero-desc">Videos are hosted securely on owners' edge nodes and streamed via Origin Shield pull-through cache. Zero permanent cloud storage; zero port forwarding.</div>
    </div>
    <div class="banner-stats">
      <div><div class="stat-val">3</div><div class="stat-lbl">Active Nodes</div></div>
      <div><div class="stat-val">12</div><div class="stat-lbl">Public Titles</div></div>
      <div><div class="stat-val">99.2%</div><div class="stat-lbl">Cache Hit Ratio</div></div>
      <div><div class="stat-val">&lt; 15s</div><div class="stat-lbl">Heartbeat TTL</div></div>
    </div>
  </div>

  <div class="section-hdr">
    <div class="section-title">Trending Community Publications</div>
    <div style="font-size: 13px; color: #94a3b8;">Showing 4 of 12 discoverable streams</div>
  </div>

  <div class="grid">
    <div class="card">
      <div class="thumb-box">
        <div style="font-size: 40px; color: #64748b;">&#127916;</div>
        <div class="badge-origin">Origin Shield OK</div>
        <div class="duration-badge">12:14</div>
      </div>
      <div class="card-body">
        <div class="card-title">Tears of Steel (4K Master)</div>
        <div class="card-meta"><span>1080p / 720p ABR</span><span>1.4k views</span></div>
        <div class="uploader-pill">&#128100; Anosh (25P0320010)</div>
      </div>
    </div>

    <div class="card">
      <div class="thumb-box">
        <div style="font-size: 40px; color: #64748b;">&#128031;</div>
        <div class="badge-origin">Direct LAN Ready</div>
        <div class="duration-badge">09:56</div>
      </div>
      <div class="card-body">
        <div class="card-title">Big Buck Bunny (1080p 60fps)</div>
        <div class="card-meta"><span>1080p / 720p / 480p</span><span>3.8k views</span></div>
        <div class="uploader-pill">&#128100; Anish (25P0320009)</div>
      </div>
    </div>

    <div class="card">
      <div class="thumb-box">
        <div style="font-size: 40px; color: #64748b;">&#127912;</div>
        <div class="badge-origin">Origin Shield OK</div>
        <div class="duration-badge">10:53</div>
      </div>
      <div class="card-body">
        <div class="card-title">Cosmos Laundromat (First Cycle)</div>
        <div class="card-meta"><span>1080p NVENC HLS</span><span>942 views</span></div>
        <div class="uploader-pill">&#128100; Anosh (25P0320010)</div>
      </div>
    </div>

    <div class="card">
      <div class="thumb-box">
        <div style="font-size: 40px; color: #64748b;">&#9876;</div>
        <div class="badge-origin">Origin Shield OK</div>
        <div class="duration-badge">14:48</div>
      </div>
      <div class="card-body">
        <div class="card-title">Sintel (Open Movie Foundation)</div>
        <div class="card-meta"><span>1080p / 720p ABR</span><span>2.1k views</span></div>
        <div class="uploader-pill">&#128100; Anish (25P0320009)</div>
      </div>
    </div>
  </div>
</main>
</body>
</html>
"""

# 2. HTML for Player with Stats for Nerds HUD
player_html = """<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>CustardFin - Player with Stats for Nerds HUD</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
  body { background: #070a13; color: #f3f4f6; padding: 24px 32px; }
  .top-nav { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
  .back-btn { display: flex; align-items: center; gap: 8px; color: #94a3b8; font-size: 14px; text-decoration: none; cursor: pointer; }
  .player-container { position: relative; width: 100%; height: 580px; background: #000; border-radius: 12px; overflow: hidden; border: 1px solid #1f2937; box-shadow: 0 20px 40px rgba(0,0,0,0.8); display: flex; align-items: center; justify-content: center; }
  .video-canvas { width: 100%; height: 100%; background: radial-gradient(circle at center, #1e1b4b 0%, #000 70%); display: flex; flex-direction: column; align-items: center; justify-content: center; }
  .video-title-watermark { font-size: 32px; font-weight: 700; color: rgba(255,255,255,0.4); text-transform: uppercase; letter-spacing: 2px; }
  .routing-indicator { position: absolute; top: 16px; left: 16px; background: rgba(16, 185, 129, 0.9); backdrop-filter: blur(8px); padding: 6px 14px; border-radius: 6px; font-size: 12px; font-weight: 700; color: #fff; display: flex; align-items: center; gap: 8px; }
  .indicator-dot { width: 8px; height: 8px; border-radius: 50%; background: #fff; animation: pulse 2s infinite; }

  /* Stats for Nerds HUD Overlay */
  .stats-hud { position: absolute; top: 16px; right: 16px; background: rgba(15, 23, 42, 0.88); backdrop-filter: blur(10px); border: 1px solid rgba(255,255,255,0.15); border-radius: 8px; padding: 14px 18px; width: 380px; font-family: 'JetBrains Mono', monospace; font-size: 11px; line-height: 1.6; color: #cbd5e1; box-shadow: 0 10px 25px rgba(0,0,0,0.5); }
  .hud-header { display: flex; justify-content: space-between; border-bottom: 1px solid rgba(255,255,255,0.1); padding-bottom: 6px; margin-bottom: 8px; font-weight: 700; color: #38bdf8; font-size: 12px; }
  .hud-row { display: flex; justify-content: space-between; }
  .hud-lbl { color: #94a3b8; }
  .hud-val { color: #f1f5f9; font-weight: 600; }
  .hud-val.hit { color: #10b981; }
  .hud-val.gold { color: #eab308; }

  /* Video Player Controls */
  .player-controls { position: absolute; bottom: 0; left: 0; right: 0; background: linear-gradient(to top, rgba(0,0,0,0.9), transparent); padding: 20px 24px; display: flex; flex-direction: column; gap: 10px; }
  .scrubber-bar { width: 100%; height: 5px; background: rgba(255,255,255,0.2); border-radius: 3px; position: relative; }
  .scrubber-fill { width: 45%; height: 100%; background: #eab308; border-radius: 3px; }
  .scrubber-buffer { width: 62%; height: 100%; background: rgba(255,255,255,0.4); border-radius: 3px; position: absolute; top: 0; left: 0; z-index: -1; }
  .controls-row { display: flex; justify-content: space-between; align-items: center; }
  .ctrl-btn-group { display: flex; align-items: center; gap: 16px; font-size: 14px; }
  .btn-play { font-size: 18px; color: #eab308; background: none; border: none; cursor: pointer; }
  .time-stamp { font-size: 12px; color: #94a3b8; font-family: 'JetBrains Mono', monospace; }
  .quality-badge { background: #1f2937; border: 1px solid #374151; padding: 4px 10px; border-radius: 4px; font-size: 11px; font-weight: 700; color: #eab308; }
</style>
</head>
<body>
<div class="top-nav">
  <div class="back-btn">&#8592; Back to Community Catalog</div>
  <div style="font-size: 14px; font-weight: 600; color: #fff;">Tears of Steel 4K Master (m-8f7a2c1e4d9b)</div>
  <div style="font-size: 12px; color: #94a3b8;">Node: node-4060-desktop (Online)</div>
</div>

<div class="player-container">
  <div class="video-canvas">
    <div class="video-title-watermark">TEARS OF STEEL</div>
    <div style="font-size: 14px; color: #64748b; margin-top: 8px;">CustardFin Origin Shield Stream (RFC 8216 HLS)</div>
  </div>

  <div class="routing-indicator">
    <div class="indicator-dot"></div>
    <span>Stream Mode: Origin Shield Cache (HIT)</span>
  </div>

  <div class="stats-hud">
    <div class="hud-header">
      <span>STATS FOR NERDS (Telemetry HUD)</span>
      <span style="color: #94a3b8; cursor: pointer;">&times;</span>
    </div>
    <div class="hud-row"><span class="hud-lbl">Video Resolution</span><span class="hud-val">1920x1080@30fps (1080p)</span></div>
    <div class="hud-row"><span class="hud-lbl">Codec / Profile</span><span class="hud-val">H.264 High (avc1.640028) / AAC</span></div>
    <div class="hud-row"><span class="hud-lbl">Active Delivery Path</span><span class="hud-val hit">S3 Origin Shield (Single Ingest)</span></div>
    <div class="hud-row"><span class="hud-lbl">Edge Node Uplink</span><span class="hud-val">0 bps (Cached By Shield)</span></div>
    <div class="hud-row"><span class="hud-lbl">Current Bitrate</span><span class="hud-val gold">4,982 kbps (Adaptive Ladder)</span></div>
    <div class="hud-row"><span class="hud-lbl">Forward Buffer</span><span class="hud-val">16.8 seconds (4 Chunks)</span></div>
    <div class="hud-row"><span class="hud-lbl">Dropped Frames</span><span class="hud-val">0 / 1420 (0.00%)</span></div>
    <div class="hud-row"><span class="hud-lbl">HLS Segment Length</span><span class="hud-val">4.0s (Strict GOP g=120)</span></div>
    <div class="hud-row"><span class="hud-lbl">HMAC Token Expiry</span><span class="hud-val">248s remaining (Signed URL)</span></div>
    <div class="hud-row"><span class="hud-lbl">Chunk Request Latency</span><span class="hud-val hit">14 ms (PoP Cache)</span></div>
  </div>

  <div class="player-controls">
    <div class="scrubber-bar">
      <div class="scrubber-fill"></div>
      <div class="scrubber-buffer"></div>
    </div>
    <div class="controls-row">
      <div class="ctrl-btn-group">
        <button class="btn-play">&#10074;&#10074;</button>
        <span class="time-stamp">05:32 / 12:14</span>
      </div>
      <div class="ctrl-btn-group">
        <span class="quality-badge">1080p ABR</span>
        <span style="font-size: 12px; color: #38bdf8; font-weight: 600;">[Stats HUD Active]</span>
      </div>
    </div>
  </div>
</div>
</body>
</html>
"""

# 3. HTML for Edge Library Management Modal
modal_html = """<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>CustardFin - My Edge Library Modal</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
  body { background: rgba(11, 15, 25, 0.95); color: #f3f4f6; display: flex; align-items: center; justify-content: center; height: 100vh; padding: 20px; }
  .modal { width: 880px; background: #111827; border: 1px solid #1f2937; border-radius: 14px; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.8); overflow: hidden; }
  .modal-hdr { background: #182234; padding: 18px 24px; border-bottom: 1px solid #1f2937; display: flex; justify-content: space-between; align-items: center; }
  .modal-title { font-size: 17px; font-weight: 700; color: #fff; display: flex; align-items: center; gap: 10px; }
  .node-pill { background: #064e3b; color: #34d399; font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 12px; text-transform: uppercase; }

  .node-status-bar { padding: 14px 24px; background: #0f172a; border-bottom: 1px solid #1f2937; display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; font-size: 12px; }
  .status-item { display: flex; flex-direction: column; gap: 3px; }
  .status-lbl { color: #94a3b8; font-size: 11px; text-transform: uppercase; font-weight: 600; }
  .status-val { font-weight: 600; color: #fff; font-family: 'JetBrains Mono', monospace; }

  .table-box { padding: 20px 24px; max-height: 400px; overflow-y: auto; }
  table { width: 100%; border-collapse: collapse; text-align: left; font-size: 13px; }
  th { padding: 10px 12px; color: #94a3b8; font-size: 11px; text-transform: uppercase; border-bottom: 1px solid #1f2937; font-weight: 600; }
  td { padding: 12px; border-bottom: 1px solid #1f2937; }
  .media-title { font-weight: 600; color: #fff; }
  .media-meta { font-size: 11px; color: #94a3b8; font-family: 'JetBrains Mono', monospace; }

  /* Switch */
  .switch-box { display: flex; align-items: center; gap: 8px; }
  .switch { position: relative; display: inline-block; width: 36px; height: 20px; }
  .slider { position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0; background-color: #374151; transition: .3s; border-radius: 20px; }
  .slider:before { position: absolute; content: ""; height: 14px; width: 14px; left: 3px; bottom: 3px; background-color: white; transition: .3s; border-radius: 50%; }
  .switch.active .slider { background-color: #eab308; }
  .switch.active .slider:before { transform: translateX(16px); }

  .btn-lan { background: #312e81; color: #a5b4fc; border: 1px solid #4338ca; padding: 4px 10px; border-radius: 6px; font-size: 11px; font-weight: 600; cursor: pointer; text-decoration: none; }
  .modal-ftr { padding: 16px 24px; background: #182234; border-top: 1px solid #1f2937; display: flex; justify-content: space-between; align-items: center; }
  .btn-gold { background: #eab308; color: #000; padding: 8px 18px; border-radius: 6px; font-size: 13px; font-weight: 600; border: none; cursor: pointer; }
</style>
</head>
<body>
<div class="modal">
  <div class="modal-hdr">
    <div class="modal-title">
      <span>&#128187; My Edge Node Library Manager</span>
      <span class="node-pill">Node: node-4060-desktop (Online)</span>
    </div>
    <span style="color: #94a3b8; cursor: pointer; font-size: 18px;">&times;</span>
  </div>

  <div class="node-status-bar">
    <div class="status-item">
      <span class="status-lbl">Transcode Engine</span>
      <span class="status-val" style="color: #38bdf8;">NVIDIA NVENC (HW Active)</span>
    </div>
    <div class="status-item">
      <span class="status-lbl">Edge Local Storage</span>
      <span class="status-val">42.8 GB / 1.0 TB (4.2%)</span>
    </div>
    <div class="status-item">
      <span class="status-lbl">Active Streams</span>
      <span class="status-val" style="color: #10b981;">2 / 10 (Within Cap)</span>
    </div>
    <div class="status-item">
      <span class="status-lbl">Heartbeat Beacon</span>
      <span class="status-val" style="color: #eab308;">Last sent: 4s ago</span>
    </div>
  </div>

  <div class="table-box">
    <table>
      <thead>
        <tr>
          <th>Media Title & Specs</th>
          <th>Ingest Status</th>
          <th>Public Manifest Sync</th>
          <th>LAN Direct Access</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>
            <div class="media-title">Tears of Steel 4K Master</div>
            <div class="media-meta">SHA-256: e3b0c442... | 1080p HLS (12:14)</div>
          </td>
          <td><span style="color: #10b981; font-weight: 600;">&#10003; READY</span></td>
          <td>
            <div class="switch-box">
              <div class="switch active"><div class="slider"></div></div>
              <span style="font-size: 11px; color: #eab308; font-weight: 600;">PUBLIC (Synced ~3KB)</span>
            </div>
          </td>
          <td><button class="btn-lan">&#128640; Launch LAN Stream</button></td>
        </tr>
        <tr>
          <td>
            <div class="media-title">Big Buck Bunny 1080p 60fps</div>
            <div class="media-meta">SHA-256: 7f83b165... | Multi-bitrate HLS (09:56)</div>
          </td>
          <td><span style="color: #10b981; font-weight: 600;">&#10003; READY</span></td>
          <td>
            <div class="switch-box">
              <div class="switch active"><div class="slider"></div></div>
              <span style="font-size: 11px; color: #eab308; font-weight: 600;">PUBLIC (Synced ~3KB)</span>
            </div>
          </td>
          <td><button class="btn-lan">&#128640; Launch LAN Stream</button></td>
        </tr>
        <tr>
          <td>
            <div class="media-title">Family Vacation Archive 2024 (Private)</div>
            <div class="media-meta">SHA-256: a14f89d0... | 1080p HLS (34:20)</div>
          </td>
          <td><span style="color: #10b981; font-weight: 600;">&#10003; READY</span></td>
          <td>
            <div class="switch-box">
              <div class="switch"><div class="slider"></div></div>
              <span style="font-size: 11px; color: #94a3b8;">PRIVATE (Home Only)</span>
            </div>
          </td>
          <td><button class="btn-lan">&#128640; Launch LAN Stream</button></td>
        </tr>
      </tbody>
    </table>
  </div>

  <div class="modal-ftr">
    <div style="font-size: 12px; color: #94a3b8;">Private items never sync metadata or chunks to cloud.</div>
    <div style="display: flex; gap: 12px;">
      <button style="background: transparent; border: 1px solid #374151; color: #d1d5db; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600;">Scan Directory</button>
      <button class="btn-gold">&#8635; Force Sync Manifests</button>
    </div>
  </div>
</div>
</body>
</html>
"""

# 4. HTML for Standalone Mini-Jellyfin Offline Web Hub
mini_jellyfin_custom_html = """<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>CustardFin Edge - Local Media Hub (Mini-Jellyfin)</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
  body { background: #0a0e17; color: #f1f5f9; padding: 0; }
  header { background: #111827; border-bottom: 1px solid #1f2937; padding: 14px 28px; display: flex; justify-content: space-between; align-items: center; }
  .brand { display: flex; align-items: center; gap: 12px; }
  .brand-badge { background: #6366f1; color: #fff; font-size: 11px; font-weight: 700; padding: 4px 8px; border-radius: 6px; text-transform: uppercase; }
  .brand-title { font-size: 18px; font-weight: 700; color: #fff; }
  .lan-pill { background: #064e3b; border: 1px solid #059669; color: #34d399; font-size: 11px; font-weight: 700; padding: 4px 12px; border-radius: 20px; display: flex; align-items: center; gap: 6px; }
  .dot { width: 6px; height: 6px; background: #34d399; border-radius: 50%; }

  .status-strip { background: #182234; border-bottom: 1px solid #1f2937; padding: 12px 28px; display: flex; gap: 32px; font-size: 12px; font-family: 'JetBrains Mono', monospace; }
  .stat-block span:first-child { color: #94a3b8; margin-right: 6px; }
  .stat-block span:last-child { color: #38bdf8; font-weight: 600; }

  main { padding: 28px; max-width: 1300px; margin: 0 auto; }
  .banner-lan { background: linear-gradient(135deg, #1e1b4b 0%, #172554 100%); border: 1px solid #3b82f6; border-radius: 10px; padding: 18px 24px; margin-bottom: 24px; display: flex; justify-content: space-between; align-items: center; }
  .banner-lan h3 { font-size: 16px; font-weight: 700; margin-bottom: 4px; }
  .banner-lan p { font-size: 13px; color: #93c5fd; }

  .grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; }
  .card { background: #111827; border: 1px solid #1f2937; border-radius: 10px; overflow: hidden; }
  .card-poster { height: 160px; background: #1e293b; display: flex; align-items: center; justify-content: center; font-size: 44px; color: #64748b; position: relative; }
  .card-body { padding: 14px; }
  .card-title { font-size: 14px; font-weight: 700; color: #fff; margin-bottom: 4px; }
  .card-info { font-size: 12px; color: #94a3b8; font-family: 'JetBrains Mono', monospace; margin-bottom: 12px; }
  .btn-play-lan { width: 100%; background: #6366f1; color: #fff; border: none; padding: 8px; border-radius: 6px; font-size: 12px; font-weight: 600; cursor: pointer; }
</style>
</head>
<body>
<header>
  <div class="brand">
    <div class="brand-badge">MINI-JELLYFIN</div>
    <div class="brand-title">CustardFin Edge &mdash; Local Media Hub</div>
  </div>
  <div class="lan-pill">
    <div class="dot"></div>
    <span>STANDALONE OFFLINE MODE (PORT 8089)</span>
  </div>
</header>

<div class="status-strip">
  <div class="stat-block"><span>NODE:</span><span>node-4060-desktop</span></div>
  <div class="stat-block"><span>TRANSCODER:</span><span style="color:#10b981;">NVIDIA NVENC (HW Probe OK)</span></div>
  <div class="stat-block"><span>SQLITE:</span><span>WAL Mode (0 Locks)</span></div>
  <div class="stat-block"><span>CLOUD RELAY:</span><span>Connected (FRP Port 7000)</span></div>
</div>

<main>
  <div class="banner-lan">
    <div>
      <h3>Direct Local Area Network Playback Engine</h3>
      <p>Streaming directly from host NVMe storage via Go HTTP engine. 0ms cloud egress, 100% data privacy.</p>
    </div>
    <div style="font-family:'JetBrains Mono',monospace; font-size:12px; background:rgba(0,0,0,0.4); padding:8px 14px; border-radius:6px; color:#6ee7b7;">
      LAN IP: 192.168.1.105:8089
    </div>
  </div>

  <div class="grid">
    <div class="card">
      <div class="card-poster">&#127916;</div>
      <div class="card-body">
        <div class="card-title">Tears of Steel 4K Master</div>
        <div class="card-info">1080p 60fps &bull; H.264 &bull; 4s HLS</div>
        <button class="btn-play-lan">&#9658; Direct LAN Stream (0 Cloud Cost)</button>
      </div>
    </div>

    <div class="card">
      <div class="card-poster">&#128031;</div>
      <div class="card-body">
        <div class="card-title">Big Buck Bunny (1080p)</div>
        <div class="card-info">1080p 30fps &bull; Multi-Ladder HLS</div>
        <button class="btn-play-lan">&#9658; Direct LAN Stream (0 Cloud Cost)</button>
      </div>
    </div>

    <div class="card">
      <div class="card-poster">&#128274;</div>
      <div class="card-body">
        <div class="card-title">Family Vacation Archive 2024</div>
        <div class="card-info">PRIVATE &bull; Stored Locally in SQLite</div>
        <button class="btn-play-lan" style="background:#059669;">&#9658; Play Private Local Media</button>
      </div>
    </div>
  </div>
</main>
</body>
</html>
"""

# 5. HTML for Developer Settings & Topology Modal
dev_settings_html = """<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>CustardFin - Dev Settings & Grading Modal</title>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Inter', sans-serif; }
  body { background: rgba(11, 15, 25, 0.95); color: #f3f4f6; display: flex; align-items: center; justify-content: center; height: 100vh; padding: 20px; }
  .modal { width: 720px; background: #111827; border: 1px solid #1f2937; border-radius: 14px; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.8); overflow: hidden; }
  .modal-hdr { background: #182234; padding: 16px 22px; border-bottom: 1px solid #1f2937; display: flex; justify-content: space-between; align-items: center; }
  .modal-title { font-size: 16px; font-weight: 700; color: #fff; }
  .modal-body { padding: 20px 24px; font-size: 13px; display: flex; flex-direction: column; gap: 18px; }
  .form-group { display: flex; flex-direction: column; gap: 6px; }
  label { font-size: 11px; text-transform: uppercase; font-weight: 600; color: #94a3b8; }
  .input { background: #1f2937; border: 1px solid #374151; padding: 9px 12px; border-radius: 6px; color: #fff; font-family: 'JetBrains Mono', monospace; font-size: 12px; }
  .topology-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
  .topo-card { border: 1px solid #374151; background: #182234; padding: 12px; border-radius: 8px; text-align: center; cursor: pointer; }
  .topo-card.active { border-color: #eab308; background: #232219; }
  .topo-name { font-weight: 700; font-size: 13px; color: #fff; margin-bottom: 4px; }
  .topo-desc { font-size: 11px; color: #94a3b8; }
  .btn-gold { background: #eab308; color: #000; padding: 9px 18px; border-radius: 6px; font-size: 13px; font-weight: 600; border: none; cursor: pointer; }
  .modal-ftr { padding: 14px 22px; background: #182234; border-top: 1px solid #1f2937; display: flex; justify-content: space-between; align-items: center; }
</style>
</head>
<body>
<div class="modal">
  <div class="modal-hdr">
    <div class="modal-title">&#9881; Developer Diagnostics & Grading Topology</div>
    <span style="color: #94a3b8; font-size: 18px; cursor: pointer;">&times;</span>
  </div>
  <div class="modal-body">
    <div class="form-group">
      <label>Active Deployment Profile</label>
      <div class="topology-grid">
        <div class="topo-card">
          <div class="topo-name">Profile 1: Local Dev</div>
          <div class="topo-desc">Local Go API + MinIO S3 + Redis</div>
        </div>
        <div class="topo-card active">
          <div class="topo-name">Profile 2: Floci Lab</div>
          <div class="topo-desc">Local AWS S3 Emulator (:4566)</div>
        </div>
        <div class="topo-card">
          <div class="topo-name">Profile 3: Cloud VM</div>
          <div class="topo-desc">Self-Hosted Static Relay + Caddy</div>
        </div>
      </div>
    </div>

    <div class="form-group">
      <label>Cloud API Control Plane Endpoint</label>
      <input type="text" class="input" value="http://localhost:8080/v1">
    </div>

    <div class="form-group">
      <label>S3 Origin Shield Emulator Endpoint</label>
      <input type="text" class="input" value="http://localhost:4566/custard-shield-cache">
    </div>

    <div class="form-group">
      <label>Admission Control Simulator Status</label>
      <div style="background:#0f172a; padding:10px 14px; border-radius:6px; font-family:'JetBrains Mono',monospace; font-size:12px; display:flex; justify-content:space-between;">
        <span>Current Active Streams: 2 / 10</span>
        <span style="color:#10b981; font-weight:700;">Cap State: SAFE (HTTP 200 OK)</span>
      </div>
    </div>
  </div>
  <div class="modal-ftr">
    <span style="font-size:11px; color:#94a3b8;">Config hot-reloaded into Next.js React Context</span>
    <button class="btn-gold">Apply Diagnostics & Save</button>
  </div>
</div>
</body>
</html>
"""

pages = [
    ("ui_catalog_feed", catalog_html),
    ("ui_player_stats", player_html),
    ("ui_edge_library", modal_html),
    ("ui_mini_jellyfin", mini_jellyfin_custom_html),
    ("ui_dev_settings", dev_settings_html),
]

for name, html_content in pages:
    html_file = os.path.join(fig_dir, f"{name}.html")
    png_file = os.path.join(fig_dir, f"{name}.png")
    with open(html_file, "w", encoding="utf-8") as f:
        f.write(html_content)
    
    args = [
        chrome_path,
        "--headless=old",
        "--disable-gpu",
        f"--screenshot={png_file}",
        "--window-size=1280,800",
        html_file
    ]
    print(f"Generating {name}.png...")
    res = subprocess.run(args, capture_output=True, text=True)
    if os.path.exists(png_file):
        print(f"  -> SUCCESS: {name}.png ({os.path.getsize(png_file)} bytes)")
    else:
        print(f"  -> FAILED: {res.stderr}")

print("All screenshots generated successfully.")
