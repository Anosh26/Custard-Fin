import Link from "next/link";
import { fetchCatalog, formatDuration, CatalogItem } from "../lib/api-client";

interface PageProps {
  searchParams: Promise<{ q?: string }>;
}

export default async function HomePage({ searchParams }: PageProps) {
  const { q } = await searchParams;

  let items: CatalogItem[] = [];
  let total = 0;
  let fetchError: string | null = null;

  try {
    const res = await fetchCatalog(q);
    items = res.items || [];
    total = res.total || 0;
  } catch (err: unknown) {
    fetchError = err instanceof Error ? err.message : "Unable to reach Cloud API";
  }

  return (
    <div>
      {/* Hero Section */}
      <section className="hero">
        <div>
          <div style={{ display: "inline-flex", alignItems: "center", gap: "0.5rem", padding: "0.25rem 0.65rem", background: "rgba(245, 158, 11, 0.12)", border: "1px solid rgba(245, 158, 11, 0.3)", borderRadius: "9999px", fontSize: "0.75rem", color: "var(--accent-gold)", fontWeight: 600, marginBottom: "0.75rem" }}>
            <span className="status-dot online" />
            <span>Hybrid Architecture &bull; 0 Cloud Video Storage</span>
          </div>
          <h1 className="hero-title">Distributed Community Media</h1>
          <p className="hero-desc">
            Video binaries stay on the owner&apos;s home machine. Discover and stream public catalogs coordinated through lightweight 3 KB manifests at near-zero cloud cost.
          </p>
        </div>
        <div className="hero-stat">
          <div style={{ fontSize: "1.75rem", fontWeight: 800, color: "var(--accent-gold)", letterSpacing: "-0.03em" }}>
            {total}
          </div>
          <div style={{ fontSize: "0.8rem", color: "var(--text-muted)", fontWeight: 500 }}>
            Public Titles
          </div>
        </div>
      </section>

      {/* Catalog Header */}
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: "1.5rem",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
          <h2 style={{ fontSize: "1.25rem", fontWeight: 700 }}>
            {q ? `Search Results for "${q}"` : "Community Catalog"}
          </h2>
          <span style={{ fontSize: "0.75rem", padding: "0.2rem 0.5rem", background: "rgba(255, 255, 255, 0.05)", borderRadius: "4px", color: "var(--text-dim)" }}>
            {items.length} {items.length === 1 ? "video" : "videos"} available
          </span>
        </div>
        {q && (
          <Link
            href="/"
            style={{
              fontSize: "0.85rem",
              color: "var(--accent-gold)",
              textDecoration: "underline",
            }}
          >
            Clear Search
          </Link>
        )}
      </div>

      {fetchError && (
        <div
          style={{
            padding: "1.25rem",
            background: "rgba(239, 68, 68, 0.1)",
            border: "1px solid rgba(239, 68, 68, 0.3)",
            borderRadius: "var(--radius-md)",
            color: "var(--accent-red)",
            marginBottom: "2rem",
            fontSize: "0.9rem",
          }}
        >
          <strong>Cloud Catalog Notice:</strong> {fetchError}. Ensure the Cloud API is running.
        </div>
      )}

      {/* Media Grid */}
      {items.length > 0 ? (
        <div className="catalog-grid">
          {items.map((item) => (
            <Link
              key={item.media_id}
              href={`/watch/${item.media_id}`}
              className="video-card"
            >
              <div className="card-thumbnail">
                <img
                  src={`/v1/media/${item.media_id}/thumbnail`}
                  alt={item.title}
                  className="card-thumb-img"
                  loading="lazy"
                />
                <div className="card-play-overlay">
                  <svg
                    width="22"
                    height="22"
                    viewBox="0 0 24 24"
                    fill="#ffffff"
                    stroke="#ffffff"
                    strokeWidth="1"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <polygon points="5 3 19 12 5 21 5 3" />
                  </svg>
                </div>
                <div className="duration-badge">
                  {formatDuration(item.duration_s)}
                </div>
              </div>

              <div className="card-content">
                <h3 className="card-title" title={item.title}>
                  {item.title}
                </h3>
                <div className="card-meta">
                  <div className="node-status">
                    <span
                      className={`status-dot ${item.is_node_online ? "online" : ""}`}
                    />
                    <span>{item.is_node_online ? "Online" : "Cached"}</span>
                  </div>
                  <div style={{ display: "flex", gap: "0.35rem", alignItems: "center" }}>
                    <span style={{ fontSize: "0.7rem", padding: "0.15rem 0.35rem", borderRadius: "3px", background: "rgba(245, 158, 11, 0.15)", color: "var(--accent-gold)", fontWeight: 600 }}>
                      4s HLS
                    </span>
                    <span suppressHydrationWarning style={{ fontSize: "0.75rem", color: "var(--text-dim)" }}>
                      {new Date(item.updated_at).toLocaleDateString("en-US", {
                        month: "short",
                        day: "numeric",
                      })}
                    </span>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      ) : (
        !fetchError && (
          <div className="empty-state">
            <svg
              width="54"
              height="54"
              viewBox="0 0 24 24"
              fill="none"
              stroke="var(--accent-gold)"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
              style={{ margin: "0 auto 1.25rem auto", display: "block", opacity: 0.8 }}
            >
              <rect x="2" y="2" width="20" height="20" rx="2.18" ry="2.18" />
              <line x1="7" y1="2" x2="7" y2="22" />
              <line x1="17" y1="2" x2="17" y2="22" />
              <line x1="2" y1="12" x2="22" y2="12" />
            </svg>
            <h3 style={{ fontSize: "1.2rem", fontWeight: 700, marginBottom: "0.5rem", color: "var(--text-main)" }}>
              No Public Videos in Community Catalog
            </h3>
            <p style={{ maxWidth: "460px", margin: "0 auto 1.5rem auto", fontSize: "0.88rem", lineHeight: 1.6 }}>
              Videos remain on the owner&apos;s machine. Click <strong>&quot;+ Add Video&quot;</strong> in the top navigation to upload a video or generate an instant test card fixture.
            </p>
          </div>
        )
      )}
    </div>
  );
}

