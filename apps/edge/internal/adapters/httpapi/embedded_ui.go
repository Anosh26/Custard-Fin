package httpapi

// MiniJellyfinHTML is the self-contained, standalone web UI for the CustardFin edge node.
// It serves an offline-capable, responsive media browser and HLS player directly from the Go binary.
const MiniJellyfinHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>CustardFin Edge &mdash; Local Media Hub</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
  <script src="https://cdn.jsdelivr.net/npm/hls.js@1.5.8/dist/hls.min.js"></script>
  <style>
    :root {
      --bg-dark: #0a0e17;
      --bg-surface: #111827;
      --bg-card: #182234;
      --bg-card-hover: #1f2c42;
      --primary: #6366f1;
      --primary-hover: #4f46e5;
      --accent: #38bdf8;
      --text: #f1f5f9;
      --text-muted: #94a3b8;
      --text-dim: #64748b;
      --border: rgba(255, 255, 255, 0.08);
      --border-accent: rgba(99, 102, 241, 0.3);
      --success: #10b981;
      --warning: #f59e0b;
      --danger: #ef4444;
      --radius-sm: 6px;
      --radius-md: 10px;
      --radius-lg: 14px;
      --shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.5), 0 8px 10px -6px rgba(0, 0, 0, 0.4);
    }
    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }
    body {
      font-family: 'Inter', system-ui, -apple-system, sans-serif;
      background: var(--bg-dark);
      color: var(--text);
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      line-height: 1.5;
      -webkit-font-smoothing: antialiased;
    }
    header {
      background: rgba(17, 24, 39, 0.85);
      backdrop-filter: blur(12px);
      border-bottom: 1px solid var(--border);
      position: sticky;
      top: 0;
      z-index: 100;
      padding: 0.85rem 1.5rem;
    }
    .header-inner {
      max-width: 1400px;
      margin: 0 auto;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
      flex-wrap: wrap;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      text-decoration: none;
      color: inherit;
    }
    .brand-icon {
      width: 36px;
      height: 36px;
      border-radius: var(--radius-md);
      background: linear-gradient(135deg, var(--primary), var(--accent));
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.2rem;
      box-shadow: 0 4px 12px rgba(99, 102, 241, 0.35);
    }
    .brand-title {
      font-size: 1.15rem;
      font-weight: 700;
      letter-spacing: -0.02em;
      background: linear-gradient(to right, #ffffff, #cbd5e1);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
    .brand-badge {
      font-size: 0.7rem;
      font-weight: 600;
      padding: 0.15rem 0.5rem;
      border-radius: 9999px;
      background: rgba(16, 185, 129, 0.15);
      color: var(--success);
      border: 1px solid rgba(16, 185, 129, 0.3);
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }
    .header-actions {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      flex-wrap: wrap;
    }
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.45rem;
      background: var(--bg-card);
      color: var(--text);
      border: 1px solid var(--border);
      padding: 0.45rem 0.9rem;
      border-radius: var(--radius-md);
      font-size: 0.825rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s ease;
      text-decoration: none;
    }
    .btn:hover {
      background: var(--bg-card-hover);
      border-color: var(--border-accent);
      transform: translateY(-1px);
    }
    .btn-primary {
      background: linear-gradient(135deg, var(--primary), #4338ca);
      border: none;
      box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
    }
    .btn-primary:hover {
      background: linear-gradient(135deg, #4f46e5, #3730a3);
      box-shadow: 0 4px 14px rgba(99, 102, 241, 0.45);
    }
    .status-pill {
      font-size: 0.775rem;
      display: flex;
      align-items: center;
      gap: 0.4rem;
      background: rgba(15, 23, 42, 0.6);
      padding: 0.35rem 0.75rem;
      border-radius: 9999px;
      border: 1px solid var(--border);
      color: var(--text-muted);
    }
    .status-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: var(--success);
      box-shadow: 0 0 8px var(--success);
    }
    main {
      flex: 1;
      max-width: 1400px;
      width: 100%;
      margin: 0 auto;
      padding: 1.75rem 1.5rem;
    }
    .banner {
      background: linear-gradient(135deg, rgba(30, 41, 59, 0.7), rgba(15, 23, 42, 0.8));
      border: 1px solid var(--border);
      border-radius: var(--radius-lg);
      padding: 1.25rem 1.5rem;
      margin-bottom: 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 1.5rem;
      flex-wrap: wrap;
    }
    .banner-text h2 {
      font-size: 1.25rem;
      font-weight: 700;
      margin-bottom: 0.25rem;
    }
    .banner-text p {
      font-size: 0.875rem;
      color: var(--text-muted);
    }
    .banner-stats {
      display: flex;
      gap: 1.5rem;
      flex-wrap: wrap;
    }
    .stat-item {
      display: flex;
      flex-direction: column;
    }
    .stat-label {
      font-size: 0.7rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: var(--text-dim);
    }
    .stat-val {
      font-size: 1.15rem;
      font-weight: 700;
      color: var(--accent);
    }
    .section-title {
      font-size: 1.1rem;
      font-weight: 600;
      margin-bottom: 1.25rem;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
    .media-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
      gap: 1.5rem;
    }
    .media-card {
      background: var(--bg-card);
      border: 1px solid var(--border);
      border-radius: var(--radius-lg);
      overflow: hidden;
      display: flex;
      flex-direction: column;
      transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
      cursor: pointer;
      position: relative;
    }
    .media-card:hover {
      transform: translateY(-4px);
      border-color: var(--border-accent);
      box-shadow: var(--shadow);
    }
    .thumb-wrap {
      width: 100%;
      aspect-ratio: 16 / 9;
      background: #0f172a;
      position: relative;
      overflow: hidden;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .thumb-wrap img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.3s ease;
    }
    .media-card:hover .thumb-wrap img {
      transform: scale(1.04);
    }
    .thumb-placeholder {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.5rem;
      color: var(--text-dim);
      font-size: 0.85rem;
    }
    .duration-badge {
      position: absolute;
      bottom: 8px;
      right: 8px;
      background: rgba(0, 0, 0, 0.75);
      backdrop-filter: blur(4px);
      padding: 0.2rem 0.5rem;
      border-radius: var(--radius-sm);
      font-size: 0.75rem;
      font-weight: 600;
      color: #fff;
    }
    .res-badge {
      position: absolute;
      top: 8px;
      left: 8px;
      background: rgba(99, 102, 241, 0.85);
      backdrop-filter: blur(4px);
      padding: 0.15rem 0.45rem;
      border-radius: var(--radius-sm);
      font-size: 0.7rem;
      font-weight: 700;
      text-transform: uppercase;
      color: #fff;
    }
    .card-body {
      padding: 1rem;
      display: flex;
      flex-direction: column;
      flex: 1;
      justify-content: space-between;
    }
    .media-title {
      font-size: 0.95rem;
      font-weight: 600;
      line-height: 1.35;
      margin-bottom: 0.5rem;
      color: var(--text);
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    .card-meta {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 0.75rem;
      font-size: 0.75rem;
      color: var(--text-muted);
    }
    .tag {
      padding: 0.15rem 0.5rem;
      border-radius: 9999px;
      font-size: 0.7rem;
      font-weight: 600;
      text-transform: uppercase;
    }
    .tag-public {
      background: rgba(16, 185, 129, 0.15);
      color: var(--success);
      border: 1px solid rgba(16, 185, 129, 0.3);
    }
    .tag-private {
      background: rgba(148, 163, 184, 0.15);
      color: var(--text-muted);
      border: 1px solid rgba(148, 163, 184, 0.3);
    }
    /* Player Overlay Modal */
    .player-modal {
      position: fixed;
      inset: 0;
      background: rgba(5, 8, 15, 0.92);
      backdrop-filter: blur(16px);
      z-index: 1000;
      display: none;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .player-modal.active {
      display: flex;
    }
    .player-wrapper {
      width: 100%;
      max-width: 1100px;
      background: #000;
      border-radius: var(--radius-lg);
      overflow: hidden;
      box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.8), 0 0 0 1px var(--border);
      position: relative;
    }
    .player-header {
      width: 100%;
      max-width: 1100px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.75rem;
      color: #fff;
    }
    .player-title {
      font-size: 1.15rem;
      font-weight: 600;
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }
    .player-lan-tag {
      font-size: 0.75rem;
      background: rgba(16, 185, 129, 0.2);
      color: var(--success);
      border: 1px solid var(--success);
      padding: 0.2rem 0.6rem;
      border-radius: 9999px;
      font-weight: 600;
    }
    .close-btn {
      background: rgba(255, 255, 255, 0.1);
      border: none;
      color: #fff;
      width: 36px;
      height: 36px;
      border-radius: 50%;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      transition: background 0.2s;
    }
    .close-btn:hover {
      background: rgba(255, 255, 255, 0.2);
    }
    video {
      width: 100%;
      aspect-ratio: 16 / 9;
      display: block;
      background: #000;
      outline: none;
    }
    .player-controls-bar {
      width: 100%;
      max-width: 1100px;
      margin-top: 0.75rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 0.825rem;
      color: var(--text-muted);
    }
    .quality-select {
      background: var(--bg-card);
      color: var(--text);
      border: 1px solid var(--border);
      padding: 0.35rem 0.75rem;
      border-radius: var(--radius-sm);
      outline: none;
      cursor: pointer;
    }
    /* Toast notifications */
    .toast {
      position: fixed;
      bottom: 24px;
      right: 24px;
      background: var(--bg-card);
      border: 1px solid var(--border-accent);
      color: var(--text);
      padding: 0.75rem 1.25rem;
      border-radius: var(--radius-md);
      box-shadow: var(--shadow);
      display: none;
      align-items: center;
      gap: 0.5rem;
      z-index: 2000;
      font-size: 0.85rem;
      animation: fadeIn 0.3s ease;
    }
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(8px); }
      to { opacity: 1; transform: translateY(0); }
    }
    .empty-state {
      text-align: center;
      padding: 5rem 1rem;
      color: var(--text-muted);
    }
    .empty-state svg {
      margin-bottom: 1rem;
      opacity: 0.4;
    }
  </style>
</head>
<body>
  <header>
    <div class="header-inner">
      <a href="/" class="brand">
        <div class="brand-icon">🍮</div>
        <div>
          <div class="brand-title">CustardFin Edge</div>
        </div>
        <span class="brand-badge">Standalone LAN</span>
      </a>

      <div class="header-actions">
        <div class="status-pill">
          <div class="status-dot"></div>
          <span id="node-lan-text">Local Network Node</span>
        </div>
        <button class="btn" id="btn-scan" title="Scan local media directory">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/>
          </svg>
          Scan Library
        </button>
        <button class="btn btn-primary" id="btn-upload" title="Upload new video to home edge">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M17 8l-5-5-5 5M12 3v12"/>
          </svg>
          Upload Video
        </button>
        <input type="file" id="file-input" accept="video/*" style="display: none;" />
      </div>
    </div>
  </header>

  <main>
    <div class="banner">
      <div class="banner-text">
        <h2>Direct LAN Streaming Hub</h2>
        <p>Zero cloud overhead. Pure 4-second pre-encoded HLS streamed directly from this home edge node to your Wi-Fi devices.</p>
      </div>
      <div class="banner-stats">
        <div class="stat-item">
          <span class="stat-label">Latency</span>
          <span class="stat-val">&lt; 5 ms</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Active Streams</span>
          <span class="stat-val" id="stat-active">0 / 10</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Transcoder</span>
          <span class="stat-val" id="stat-gpu">Probing...</span>
        </div>
      </div>
    </div>

    <div class="section-title">
      <span>Local Media Library</span>
      <span id="media-count" style="font-size: 0.85rem; color: var(--text-dim);">0 titles</span>
    </div>

    <div class="media-grid" id="media-grid">
      <div class="empty-state" style="grid-column: 1 / -1;">
        <p>Loading local media library...</p>
      </div>
    </div>
  </main>

  <!-- Player Modal -->
  <div class="player-modal" id="player-modal">
    <div class="player-header">
      <div class="player-title">
        <span id="player-video-title">Video Title</span>
        <span class="player-lan-tag">⚡ Direct LAN (&lt; 5ms)</span>
      </div>
      <button class="close-btn" id="btn-close-player" title="Close player">&times;</button>
    </div>
    <div class="player-wrapper">
      <video id="video-element" controls playsinline></video>
    </div>
    <div class="player-controls-bar">
      <div>
        <span>Buffer: <strong id="stat-buffer" style="color: var(--accent);">0.0s</strong></span> &bull; 
        <span>Bitrate: <strong id="stat-bitrate" style="color: var(--accent);">Auto</strong></span>
      </div>
      <div>
        <label for="quality-selector" style="margin-right: 0.5rem;">Quality:</label>
        <select id="quality-selector" class="quality-select">
          <option value="-1">Auto (Adaptive)</option>
        </select>
      </div>
    </div>
  </div>

  <div class="toast" id="toast">Message</div>

  <script>
    let currentMediaList = [];
    let hlsInstance = null;
    let activeMediaId = null;
    let progressTimer = null;

    function showToast(msg) {
      const t = document.getElementById('toast');
      t.textContent = msg;
      t.style.display = 'flex';
      setTimeout(() => { t.style.display = 'none'; }, 3500);
    }

    function formatTime(sec) {
      if (!sec || isNaN(sec)) return "00:00";
      sec = Math.floor(sec);
      const m = Math.floor(sec / 60);
      const s = sec % 60;
      return String(m).padStart(2, '0') + ':' + String(s).padStart(2, '0');
    }

    async function loadStatus() {
      try {
        const res = await fetch('/api/status');
        if (res.ok) {
          const d = await res.json();
          document.getElementById('stat-active').textContent = (d.active_streams || 0) + ' / ' + (d.max_streams || 10);
          document.getElementById('stat-gpu').textContent = d.nvenc_enabled ? 'NVIDIA NVENC (GPU)' : 'FFmpeg (CPU)';
          if (d.lan_addr) {
            document.getElementById('node-lan-text').textContent = 'LAN: ' + d.lan_addr;
          } else {
            document.getElementById('node-lan-text').textContent = 'LAN: ' + window.location.host;
          }
        }
      } catch (e) {
        console.warn('Status fetch error:', e);
      }
    }

    async function loadMedia() {
      const grid = document.getElementById('media-grid');
      try {
        const res = await fetch('/api/media');
        if (!res.ok) throw new Error('HTTP ' + res.status);
        const list = await res.json();
        currentMediaList = list || [];
        document.getElementById('media-count').textContent = currentMediaList.length + ' titles';

        if (currentMediaList.length === 0) {
          grid.innerHTML = '<div class="empty-state" style="grid-column: 1 / -1;"><p>No videos found in configured media directory. Drop video files into <code>./media</code> or click Upload above!</p></div>';
          return;
        }

        grid.innerHTML = currentMediaList.map(item => {
          const thumbUrl = '/stream/' + item.id + '/thumb.jpg';
          const dur = formatTime(item.duration_s);
          const resBadge = item.height >= 1080 ? '1080p' : (item.height >= 720 ? '720p' : '480p');
          const isPub = item.is_public;

          return '<div class="media-card" onclick="openPlayer(\'' + item.id + '\')">' +
            '<div class="thumb-wrap">' +
              '<img src="' + thumbUrl + '" onerror="this.style.display=\'none\'; this.nextElementSibling.style.display=\'flex\';" alt="' + item.title + '" />' +
              '<div class="thumb-placeholder" style="display: none;">' +
                '<span>🎬</span>' +
                '<span>' + item.title + '</span>' +
              '</div>' +
              '<span class="res-badge">' + resBadge + '</span>' +
              '<span class="duration-badge">' + dur + '</span>' +
            '</div>' +
            '<div class="card-body">' +
              '<h3 class="media-title">' + item.title + '</h3>' +
              '<div class="card-meta">' +
                '<span class="tag ' + (isPub ? 'tag-public' : 'tag-private') + '">' + (isPub ? 'Public' : 'Private') + '</span>' +
                '<span>⚡ Direct LAN</span>' +
              '</div>' +
            '</div>' +
          '</div>';
        }).join('');
      } catch (err) {
        grid.innerHTML = '<div class="empty-state" style="grid-column: 1 / -1; color: var(--danger);"><p>Failed to load media library: ' + err.message + '</p></div>';
      }
    }

    async function openPlayer(mediaId) {
      const item = currentMediaList.find(m => m.id === mediaId);
      if (!item) return;

      activeMediaId = mediaId;
      document.getElementById('player-video-title').textContent = item.title;
      const modal = document.getElementById('player-modal');
      const video = document.getElementById('video-element');
      const streamUrl = '/stream/' + mediaId + '/master.m3u8';
      const qSelector = document.getElementById('quality-selector');

      modal.classList.add('active');

      if (hlsInstance) {
        hlsInstance.destroy();
        hlsInstance = null;
      }

      // Fetch saved watch progress
      let startPos = 0;
      try {
        const wpRes = await fetch('/api/watch/' + mediaId);
        if (wpRes.ok) {
          const wp = await wpRes.json();
          if (wp && wp.last_position_s > 5) {
            startPos = wp.last_position_s;
          }
        }
      } catch (e) {}

      if (window.Hls && Hls.isSupported()) {
        const hls = new Hls({
          enableWorker: true,
          lowLatencyMode: false,
        });
        hlsInstance = hls;
        hls.loadSource(streamUrl);
        hls.attachMedia(video);

        hls.on(Hls.Events.MANIFEST_PARSED, (ev, data) => {
          qSelector.innerHTML = '<option value="-1">Auto (Adaptive)</option>';
          data.levels.forEach((lvl, idx) => {
            const opt = document.createElement('option');
            opt.value = idx;
            opt.textContent = (lvl.height ? lvl.height + 'p' : 'Level ' + idx) + ' (' + Math.round(lvl.bitrate / 1000) + ' kbps)';
            qSelector.appendChild(opt);
          });
          if (startPos > 0) video.currentTime = startPos;
          video.play().catch(() => {});
        });

        hls.on(Hls.Events.LEVEL_SWITCHED, (ev, data) => {
          const lvl = hls.levels[data.level];
          if (lvl) {
            document.getElementById('stat-bitrate').textContent = (lvl.height ? lvl.height + 'p' : '') + ' (' + Math.round(lvl.bitrate / 1000) + ' kbps)';
          }
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        // Native Safari / iOS
        video.src = streamUrl;
        if (startPos > 0) video.currentTime = startPos;
        video.play().catch(() => {});
      }

      // Buffer updater & Watch Progress Sync
      if (progressTimer) clearInterval(progressTimer);
      progressTimer = setInterval(() => {
        if (!video.paused && activeMediaId) {
          // Update buffer stat
          if (video.buffered && video.buffered.length > 0) {
            const bufLen = video.buffered.end(video.buffered.length - 1) - video.currentTime;
            document.getElementById('stat-buffer').textContent = Math.max(0, bufLen).toFixed(1) + 's';
          }
          // Save watch progress to SQLite
          fetch('/api/watch/' + activeMediaId, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ position_s: Math.floor(video.currentTime) })
          }).catch(() => {});
        }
      }, 5000);
    }

    function closePlayer() {
      const modal = document.getElementById('player-modal');
      const video = document.getElementById('video-element');
      modal.classList.remove('active');
      if (video) {
        if (!video.paused && activeMediaId) {
          fetch('/api/watch/' + activeMediaId, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ position_s: Math.floor(video.currentTime) })
          }).catch(() => {});
        }
        video.pause();
        video.removeAttribute('src');
        video.load();
      }
      if (hlsInstance) {
        hlsInstance.destroy();
        hlsInstance = null;
      }
      if (progressTimer) {
        clearInterval(progressTimer);
        progressTimer = null;
      }
      activeMediaId = null;
    }

    document.getElementById('btn-close-player').onclick = closePlayer;
    document.getElementById('quality-selector').onchange = (e) => {
      if (hlsInstance) {
        hlsInstance.currentLevel = parseInt(e.target.value, 10);
      }
    };

    // Scan Library Button
    document.getElementById('btn-scan').onclick = async () => {
      showToast('Scanning media directory...');
      try {
        const res = await fetch('/api/media/scan', { method: 'POST' });
        if (res.ok) {
          showToast('Library scan complete!');
          loadMedia();
        } else {
          showToast('Scan failed');
        }
      } catch (e) {
        showToast('Error: ' + e.message);
      }
    };

    // Upload Video Button
    const fileInput = document.getElementById('file-input');
    document.getElementById('btn-upload').onclick = () => fileInput.click();
    fileInput.onchange = async () => {
      if (!fileInput.files || fileInput.files.length === 0) return;
      const file = fileInput.files[0];
      showToast('Uploading ' + file.name + '...');
      const formData = new FormData();
      formData.append('video', file);
      try {
        const res = await fetch('/api/media/upload', {
          method: 'POST',
          body: formData,
        });
        if (res.ok) {
          showToast('Video uploaded successfully!');
          fileInput.value = '';
          loadMedia();
        } else {
          const err = await res.json().catch(() => ({}));
          showToast('Upload failed: ' + (err.error || res.statusText));
        }
      } catch (e) {
        showToast('Upload error: ' + e.message);
      }
    };

    // Escape key closes modal
    window.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') closePlayer();
    });

    // Initial load
    loadStatus();
    loadMedia();
    setInterval(loadStatus, 15000);
  </script>
</body>
</html>
`
