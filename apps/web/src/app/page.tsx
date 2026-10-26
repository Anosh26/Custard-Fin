import Link from "next/link";
import { fetchCatalog, CatalogItem } from "../lib/api-client";
import CatalogFeed from "../components/CatalogFeed";

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

      {/* Dynamic Catalog Feed */}
      <CatalogFeed
        initialItems={items}
        initialTotal={total}
        searchQuery={q}
      />
    </div>
  );
}

