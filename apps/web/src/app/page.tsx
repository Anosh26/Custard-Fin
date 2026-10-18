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
          <h1 className="hero-title">Distributed Community Media</h1>
          <p className="hero-desc">
            Video binaries stay on the owner&apos;s home machine. Discover and stream public catalogs coordinated through lightweight 3 KB manifests at near-zero cloud cost.
          </p>
        </div>
        <div className="hero-stat">
          <div style={{ fontSize: "1.5rem", fontWeight: 700, color: "var(--accent-gold)" }}>
            {total}
          </div>
          <div style={{ fontSize: "0.8rem", color: "var(--text-muted)" }}>
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
        <h2 style={{ fontSize: "1.25rem", fontWeight: 600 }}>
          {q ? `Search Results for "${q}"` : "Community Catalog"}
        </h2>
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
                    <span>Node: {item.node_id.slice(0, 8)}...</span>
                  </div>
                  <span>
                    {new Date(item.updated_at).toLocaleDateString(undefined, {
                      month: "short",
                      day: "numeric",
                    })}
                  </span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      ) : (
        !fetchError && (
          <div className="empty-state">
            <svg
              width="48"
              height="48"
              viewBox="0 0 24 24"
              fill="none"
              stroke="var(--text-dim)"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
              style={{ margin: "0 auto 1rem auto", display: "block" }}
            >
              <rect x="2" y="2" width="20" height="20" rx="2.18" ry="2.18" />
              <line x1="7" y1="2" x2="7" y2="22" />
              <line x1="17" y1="2" x2="17" y2="22" />
              <line x1="2" y1="12" x2="22" y2="12" />
              <line x1="2" y1="7" x2="7" y2="7" />
              <line x1="2" y1="17" x2="7" y2="17" />
              <line x1="17" y1="17" x2="22" y2="17" />
              <line x1="17" y1="7" x2="22" y2="7" />
            </svg>
            <h3 style={{ fontSize: "1.2rem", color: "var(--text-main)", marginBottom: "0.5rem" }}>
              No Community Media Found
            </h3>
            <p style={{ maxWidth: "480px", margin: "0 auto", fontSize: "0.9rem" }}>
              Home edge nodes can toggle any scanned media as &quot;Public&quot; to publish a lightweight 3 KB manifest into this global discovery catalog.
            </p>
          </div>
        )
      )}
    </div>
  );
}
