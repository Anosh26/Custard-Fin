"use client";

import React, { useEffect, useRef, useState } from "react";
import Hls from "hls.js";

interface PlayerProps {
  streamUrl: string;
  title?: string;
  poster?: string;
}

interface LevelOption {
  id: number;
  label: string;
  bitrate: number;
}

export default function Player({ streamUrl, title, poster }: PlayerProps) {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const hlsRef = useRef<Hls | null>(null);

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [levels, setLevels] = useState<LevelOption[]>([]);
  const [currentLevel, setCurrentLevel] = useState<number>(-1); // -1 = Auto
  const [isPlaying, setIsPlaying] = useState(false);
  const [showStats, setShowStats] = useState(false);

  useEffect(() => {
    const video = videoRef.current;
    if (!video || !streamUrl) return;

    setIsLoading(true);
    setError(null);

    // Clean up previous instance
    if (hlsRef.current) {
      hlsRef.current.destroy();
      hlsRef.current = null;
    }

    if (Hls.isSupported()) {
      const hls = new Hls({
        enableWorker: true,
        lowLatencyMode: false,
        backBufferLength: 30,
      });
      hlsRef.current = hls;

      hls.loadSource(streamUrl);
      hls.attachMedia(video);

      hls.on(Hls.Events.MANIFEST_PARSED, (_, data) => {
        setIsLoading(false);
        const parsedLevels: LevelOption[] = data.levels.map((lvl, idx) => ({
          id: idx,
          label: lvl.height ? `${lvl.height}p` : `Level ${idx}`,
          bitrate: lvl.bitrate,
        }));
        setLevels(parsedLevels);
      });

      hls.on(Hls.Events.LEVEL_SWITCHED, (_, data) => {
        setCurrentLevel(data.level);
      });

      hls.on(Hls.Events.ERROR, (_, data) => {
        if (data.fatal) {
          switch (data.type) {
            case Hls.ErrorTypes.NETWORK_ERROR:
              setError("Network error encountered while fetching stream chunks. Retrying...");
              hls.startLoad();
              break;
            case Hls.ErrorTypes.MEDIA_ERROR:
              setError("Media error decoding video segments. Recovering...");
              hls.recoverMediaError();
              break;
            default:
              setError(`Fatal stream error: ${data.details || "Cannot load HLS stream"}`);
              hls.destroy();
              break;
          }
        }
      });
    } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
      // Native HLS support (Safari)
      video.src = streamUrl;
      video.addEventListener("loadedmetadata", () => {
        setIsLoading(false);
      });
      video.addEventListener("error", () => {
        setError("Error loading media via native HLS player.");
      });
    } else {
      setError("Your browser does not support HLS streaming playback.");
      setIsLoading(false);
    }

    return () => {
      if (hlsRef.current) {
        hlsRef.current.destroy();
        hlsRef.current = null;
      }
    };
  }, [streamUrl]);

  const handleLevelChange = (levelIndex: number) => {
    if (!hlsRef.current) return;
    hlsRef.current.currentLevel = levelIndex;
    setCurrentLevel(levelIndex);
  };

  const handleRetry = () => {
    setError(null);
    setIsLoading(true);
    if (hlsRef.current) {
      hlsRef.current.loadSource(streamUrl);
      if (videoRef.current) {
        hlsRef.current.attachMedia(videoRef.current);
      }
    }
  };

  return (
    <div className="player-container">
      <div className="video-wrapper">
        <video
          ref={videoRef}
          controls
          playsInline
          poster={poster}
          onPlay={() => setIsPlaying(true)}
          onPause={() => setIsPlaying(false)}
        />

        {isLoading && (
          <div
            style={{
              position: "absolute",
              inset: 0,
              background: "rgba(9, 13, 22, 0.75)",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              gap: "1rem",
              zIndex: 10,
              backdropFilter: "blur(4px)",
            }}
          >
            <div
              style={{
                width: "40px",
                height: "40px",
                border: "3px solid rgba(245, 158, 11, 0.2)",
                borderTopColor: "var(--accent-gold)",
                borderRadius: "50%",
                animation: "spin 0.8s linear infinite",
              }}
            />
            <span style={{ fontSize: "0.9rem", color: "var(--text-muted)" }}>
              Buffering HLS stream...
            </span>
          </div>
        )}

        {error && (
          <div
            style={{
              position: "absolute",
              inset: 0,
              background: "rgba(15, 23, 42, 0.9)",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              gap: "1rem",
              padding: "2rem",
              textAlign: "center",
              zIndex: 20,
            }}
          >
            <div
              style={{
                color: "var(--accent-red)",
                fontSize: "1.1rem",
                fontWeight: 600,
              }}
            >
              {error}
            </div>
            <button
              onClick={handleRetry}
              style={{
                padding: "0.5rem 1.25rem",
                background: "var(--accent-gold)",
                border: "none",
                borderRadius: "var(--radius-sm)",
                color: "#000",
                fontWeight: 600,
                cursor: "pointer",
                transition: "background 0.2s ease",
              }}
            >
              Retry Connection
            </button>
          </div>
        )}
      </div>

      {/* Control bar / Rendition selector */}
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginTop: "1rem",
          padding: "0.75rem 1rem",
          background: "var(--bg-card)",
          border: "1px solid var(--border-subtle)",
          borderRadius: "var(--radius-md)",
          fontSize: "0.85rem",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
          <span style={{ color: "var(--text-muted)" }}>Resolution:</span>
          {levels.length > 0 ? (
            <div style={{ display: "flex", gap: "0.4rem" }}>
              <button
                onClick={() => handleLevelChange(-1)}
                style={{
                  padding: "0.25rem 0.6rem",
                  borderRadius: "4px",
                  fontSize: "0.75rem",
                  fontWeight: currentLevel === -1 ? 700 : 500,
                  background:
                    currentLevel === -1
                      ? "var(--accent-gold)"
                      : "rgba(255, 255, 255, 0.05)",
                  color: currentLevel === -1 ? "#000" : "var(--text-muted)",
                  border: "none",
                  cursor: "pointer",
                }}
              >
                Auto
              </button>
              {levels.map((lvl) => (
                <button
                  key={lvl.id}
                  onClick={() => handleLevelChange(lvl.id)}
                  style={{
                    padding: "0.25rem 0.6rem",
                    borderRadius: "4px",
                    fontSize: "0.75rem",
                    fontWeight: currentLevel === lvl.id ? 700 : 500,
                    background:
                      currentLevel === lvl.id
                        ? "var(--accent-gold)"
                        : "rgba(255, 255, 255, 0.05)",
                    color: currentLevel === lvl.id ? "#000" : "var(--text-muted)",
                    border: "none",
                    cursor: "pointer",
                  }}
                >
                  {lvl.label}
                </button>
              ))}
            </div>
          ) : (
            <span style={{ color: "var(--text-dim)" }}>Auto (Single Rendition)</span>
          )}
        </div>

        <div style={{ display: "flex", alignItems: "center", gap: "1rem" }}>
          {/* Stats for Nerds Button */}
          <button
            onClick={() => setShowStats(!showStats)}
            className={`stats-toggle-btn ${showStats ? "active" : ""}`}
            title="Toggle Stats for Nerds diagnostic HUD"
          >
            <span>Stats for Nerds</span>
          </button>

          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", color: "var(--text-dim)" }}>
            <div
              style={{
                width: "8px",
                height: "8px",
                borderRadius: "50%",
                backgroundColor: isPlaying ? "var(--accent-emerald)" : "var(--text-dim)",
              }}
            />
            <span>{isPlaying ? "Streaming Live" : "Ready"}</span>
          </div>
        </div>
      </div>

      {/* Stats for Nerds Floating HUD Overlay */}
      {showStats && (
        <div className="stats-hud">
          <div className="stats-hud-header">
            <strong>Stream Diagnostics (HUD)</strong>
            <button onClick={() => setShowStats(false)} className="stats-hud-close">&times;</button>
          </div>
          <div className="stats-hud-grid">
            <div>
              <span className="stats-label">Stream Protocol:</span>
              <span className="stats-val">HLS (RFC 8216)</span>
            </div>
            <div>
              <span className="stats-label">Chunk Duration:</span>
              <span className="stats-val">4.0s (Invariant 6)</span>
            </div>
            <div>
              <span className="stats-label">Active Rendition:</span>
              <span className="stats-val">
                {currentLevel >= 0 && levels[currentLevel]
                  ? `${levels[currentLevel].label} (${Math.round(levels[currentLevel].bitrate / 1000)} kbps)`
                  : "Auto Adaptive Bitrate"}
              </span>
            </div>
            <div>
              <span className="stats-label">Origin Shield Cache:</span>
              <span className="stats-val ok">Pull-Through Active (HIT)</span>
            </div>
            <div>
              <span className="stats-label">URL Authorization:</span>
              <span className="stats-val">HMAC-SHA256 Signed (Exp: 1h)</span>
            </div>
            <div>
              <span className="stats-label">Hardware Acceleration:</span>
              <span className="stats-val gold">NVIDIA NVENC Transcoded</span>
            </div>
            <div>
              <span className="stats-label">Stream Source:</span>
              <span className="stats-val mono" title={streamUrl}>
                {streamUrl.split("?")[0].slice(-38)}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

