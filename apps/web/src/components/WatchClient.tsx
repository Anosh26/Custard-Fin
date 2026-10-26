"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import Player from "./Player";
import AuthModal from "./AuthModal";
import {
  PlaybackResponse,
  fetchPlayback,
  getThumbnailUrl,
  getAuthToken,
} from "../lib/api-client";

interface WatchClientProps {
  mediaId: string;
  initialPlayback: PlaybackResponse | null;
  initialError: string | null;
}

export default function WatchClient({
  mediaId,
  initialPlayback,
  initialError,
}: WatchClientProps) {
  const [playback, setPlayback] = useState<PlaybackResponse | null>(initialPlayback);
  const [error, setError] = useState<string | null>(initialError);
  const [loading, setLoading] = useState(false);
  const [isAuthOpen, setIsAuthOpen] = useState(false);

  const resolvePlayback = async () => {
    setLoading(true);
    try {
      const data = await fetchPlayback(mediaId);
      setPlayback(data);
      setError(null);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to resolve playback");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // If there's an initial error (like 403 SSR) or client has auth token, re-resolve with token
    if (initialError || getAuthToken()) {
      resolvePlayback();
    }

    const handleAuthChange = () => {
      resolvePlayback();
    };

    window.addEventListener("custardfin_auth_change", handleAuthChange);
    return () => {
      window.removeEventListener("custardfin_auth_change", handleAuthChange);
    };
  }, [mediaId]);

  const getRouteLabel = (type: string) => {
    switch (type) {
      case "direct_lan":
        return "Direct LAN (Local Subnet - 0 ms Overhead)";
      case "origin_shield":
        return "Origin Shield (S3 Pull-Through Cache - 0 Home Upload)";
      case "cache_only":
        return "Offline Cache (Node Offline - Serving Cached Media)";
      default:
        return type;
    }
  };

  return (
    <div className="player-container">
      <Link href="/" className="back-btn">
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <line x1="19" y1="12" x2="5" y2="12" />
          <polyline points="12 19 5 12 12 5" />
        </svg>
        <span>Back to Catalog</span>
      </Link>

      {loading && !playback ? (
        <div
          style={{
            background: "var(--bg-card)",
            borderRadius: "var(--radius-lg)",
            padding: "4rem 2rem",
            textAlign: "center",
            color: "var(--text-muted)",
          }}
        >
          <div style={{ fontSize: "1.1rem", marginBottom: "0.5rem" }}>
            Resolving stream route &amp; verifying credentials...
          </div>
          <div style={{ fontSize: "0.85rem", color: "var(--text-dim)" }}>
            Checking Origin Shield cache &amp; edge tunnel
          </div>
        </div>
      ) : error && !playback ? (
        <div
          style={{
            background: "var(--bg-card)",
            border: "1px solid rgba(239, 68, 68, 0.4)",
            borderRadius: "var(--radius-lg)",
            padding: "3rem 2rem",
            textAlign: "center",
          }}
        >
          <div
            style={{
              width: "56px",
              height: "56px",
              borderRadius: "50%",
              background: "rgba(239, 68, 68, 0.15)",
              color: "var(--accent-red)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              margin: "0 auto 1.5rem auto",
            }}
          >
            <svg
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
          </div>

          <h2 style={{ fontSize: "1.4rem", marginBottom: "0.75rem", color: "var(--text-main)" }}>
            Playback Stream Unavailable
          </h2>
          <p
            style={{
              color: "var(--text-muted)",
              maxWidth: "520px",
              margin: "0 auto 1.5rem auto",
              lineHeight: 1.6,
              fontSize: "0.95rem",
            }}
          >
            {error}
          </p>

          <div
            style={{
              display: "inline-flex",
              gap: "1rem",
              marginTop: "0.5rem",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            {error.toLowerCase().includes("private") && (
              <button
                onClick={() => setIsAuthOpen(true)}
                className="nav-btn primary"
                style={{ padding: "0.6rem 1.25rem" }}
              >
                Sign In to Unlock
              </button>
            )}
            <Link
              href="/"
              style={{
                padding: "0.6rem 1.25rem",
                borderRadius: "var(--radius-sm)",
                background: "rgba(255, 255, 255, 0.1)",
                color: "var(--text-main)",
                fontSize: "0.9rem",
              }}
            >
              Browse Other Titles
            </Link>
          </div>
        </div>
      ) : playback ? (
        <div>
          <Player streamUrl={playback.stream_url} poster={getThumbnailUrl(mediaId, getAuthToken())} />

          <div
            style={{
              marginTop: "1.5rem",
              padding: "1.25rem",
              background: "var(--bg-card)",
              border: "1px solid var(--border-subtle)",
              borderRadius: "var(--radius-md)",
            }}
          >
            <div
              style={{
                display: "flex",
                flexWrap: "wrap",
                alignItems: "center",
                justifyContent: "space-between",
                gap: "1rem",
                marginBottom: "0.75rem",
              }}
            >
              <div>
                <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", marginBottom: "0.25rem" }}>
                  <div style={{ fontSize: "0.8rem", color: "var(--text-dim)" }}>
                    Active Stream Routing (Play Decision)
                  </div>
                  {playback.is_public === false && (
                    <span
                      style={{
                        fontSize: "0.7rem",
                        padding: "0.15rem 0.45rem",
                        borderRadius: "4px",
                        background: "rgba(245, 158, 11, 0.15)",
                        border: "1px solid rgba(245, 158, 11, 0.3)",
                        color: "var(--accent-gold)",
                        fontWeight: 700,
                      }}
                    >
                      🔒 Private Media (Owner Stream)
                    </span>
                  )}
                </div>
                <div
                  style={{
                    fontSize: "0.95rem",
                    fontWeight: 600,
                    color: "var(--accent-gold)",
                    display: "flex",
                    alignItems: "center",
                    gap: "0.5rem",
                  }}
                >
                  <span
                    style={{
                      width: "8px",
                      height: "8px",
                      borderRadius: "50%",
                      backgroundColor: "var(--accent-gold)",
                    }}
                  />
                  {getRouteLabel(playback.route_type)}
                </div>
              </div>

              <div style={{ textAlign: "right", fontSize: "0.8rem", color: "var(--text-dim)" }}>
                <div>Media ID: {mediaId}</div>
                <div suppressHydrationWarning>
                  URL Expiry: {new Date(playback.expires_at).toLocaleTimeString()}
                </div>
              </div>
            </div>

            {/* Uploader & Node information */}
            <div
              style={{
                borderTop: "1px solid var(--border-subtle)",
                paddingTop: "0.75rem",
                display: "flex",
                flexWrap: "wrap",
                gap: "1.25rem",
                fontSize: "0.85rem",
              }}
            >
              <div style={{ display: "flex", alignItems: "center", gap: "0.4rem" }}>
                <span style={{ color: "var(--text-dim)" }}>Uploader:</span>
                <span style={{ color: "var(--accent-gold)", fontWeight: 600 }}>
                  @{playback.uploader_username || "community"}
                </span>
              </div>
              {playback.node_name && (
                <div style={{ display: "flex", alignItems: "center", gap: "0.4rem" }}>
                  <span style={{ color: "var(--text-dim)" }}>Edge Node:</span>
                  <span style={{ color: "var(--text-main)", fontWeight: 500 }}>
                    {playback.node_name}
                  </span>
                </div>
              )}
            </div>
          </div>
        </div>
      ) : null}

      <AuthModal
        isOpen={isAuthOpen}
        onClose={() => setIsAuthOpen(false)}
        onSuccess={resolvePlayback}
      />
    </div>
  );
}
