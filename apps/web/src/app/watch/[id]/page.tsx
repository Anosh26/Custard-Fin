import Link from "next/link";
import Player from "../../../components/Player";
import { fetchPlayback, PlaybackResponse } from "../../../lib/api-client";

interface WatchPageProps {
  params: Promise<{ id: string }>;
}

export default async function WatchPage({ params }: WatchPageProps) {
  const { id } = await params;

  let playback: PlaybackResponse | null = null;
  let playbackError: string | null = null;

  try {
    playback = await fetchPlayback(id);
  } catch (err: unknown) {
    playbackError = err instanceof Error ? err.message : "Failed to resolve playback route";
  }

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

      {playbackError ? (
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
            {playbackError}
          </p>

          <div
            style={{
              display: "inline-flex",
              gap: "1rem",
              marginTop: "0.5rem",
            }}
          >
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
          <Player streamUrl={playback.stream_url} />

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
              }}
            >
              <div>
                <div style={{ fontSize: "0.8rem", color: "var(--text-dim)", marginBottom: "0.25rem" }}>
                  Active Stream Routing (Play Decision)
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
                <div>Media ID: {id}</div>
                <div suppressHydrationWarning>URL Expiry: {new Date(playback.expires_at).toLocaleTimeString()}</div>
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
