"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import {
  CatalogItem,
  fetchCatalog,
  formatDuration,
  getThumbnailUrl,
  getAuthToken,
} from "../lib/api-client";

interface CatalogFeedProps {
  initialItems: CatalogItem[];
  initialTotal: number;
  searchQuery?: string;
}

export default function CatalogFeed({
  initialItems,
  initialTotal,
  searchQuery,
}: CatalogFeedProps) {
  const [items, setItems] = useState<CatalogItem[]>(initialItems);
  const [total, setTotal] = useState<number>(initialTotal);
  const [loading, setLoading] = useState(false);

  const loadCatalog = async () => {
    try {
      setLoading(true);
      const res = await fetchCatalog(searchQuery);
      setItems(res.items || []);
      setTotal(res.total || 0);
    } catch {
      // Fallback to initial items on error
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // If the client has an auth token, re-fetch to load private items immediately
    if (getAuthToken()) {
      loadCatalog();
    }

    const handleAuthChange = () => {
      loadCatalog();
    };

    window.addEventListener("custardfin_auth_change", handleAuthChange);
    return () => {
      window.removeEventListener("custardfin_auth_change", handleAuthChange);
    };
  }, [searchQuery]);

  if (items.length === 0 && !loading) {
    return (
      <div className="empty-state">
        <svg
          width="54"
          height="54"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          style={{ margin: "0 auto 1.25rem auto", color: "var(--text-dim)" }}
        >
          <path d="m22 8-6 4 6 4V8Z" />
          <rect width="14" height="12" x="2" y="6" rx="2" ry="2" />
        </svg>
        <h3 style={{ fontSize: "1.2rem", marginBottom: "0.5rem" }}>
          No Media Available
        </h3>
        <p style={{ color: "var(--text-muted)", maxWidth: "440px", margin: "0 auto 1.5rem auto", fontSize: "0.9rem" }}>
          {searchQuery
            ? `No media found matching "${searchQuery}". Try a different keyword.`
            : "No videos found in community catalog. Toggle Public in My Library to publish your home media."}
        </p>
      </div>
    );
  }

  return (
    <div className="catalog-grid">
      {items.map((item) => (
        <Link
          key={item.media_id}
          href={`/watch/${item.media_id}`}
          className="video-card"
        >
          <div className="card-thumbnail">
            <img
              src={getThumbnailUrl(item.media_id)}
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
            
            {/* Private Badge overlay */}
            {!item.is_public && (
              <div
                style={{
                  position: "absolute",
                  top: "0.5rem",
                  left: "0.5rem",
                  background: "rgba(15, 23, 42, 0.85)",
                  backdropFilter: "blur(4px)",
                  border: "1px solid rgba(245, 158, 11, 0.5)",
                  borderRadius: "4px",
                  padding: "0.2rem 0.5rem",
                  fontSize: "0.7rem",
                  fontWeight: 700,
                  color: "var(--accent-gold)",
                  display: "flex",
                  alignItems: "center",
                  gap: "0.3rem",
                  zIndex: 2,
                }}
              >
                <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                  <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
                <span>Private</span>
              </div>
            )}

            <div className="duration-badge">
              {formatDuration(item.duration_s)}
            </div>
          </div>

          <div className="card-content">
            <h3 className="card-title" title={item.title}>
              {item.title}
            </h3>

            {/* Uploader & Node info */}
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "0.4rem",
                fontSize: "0.78rem",
                color: "var(--text-muted)",
                marginBottom: "0.5rem",
                whiteSpace: "nowrap",
                overflow: "hidden",
                textOverflow: "ellipsis",
              }}
            >
              <span style={{ color: "var(--accent-gold)", fontWeight: 600 }}>
                @{item.uploader_username || "community"}
              </span>
              {item.node_name && (
                <>
                  <span style={{ color: "var(--text-dim)" }}>&bull;</span>
                  <span style={{ color: "var(--text-dim)", fontSize: "0.72rem" }} title={`Node: ${item.node_name}`}>
                    {item.node_name}
                  </span>
                </>
              )}
            </div>

            <div className="card-meta">
              <div className="node-status">
                <span
                  className={`status-dot ${item.is_node_online ? "online" : ""}`}
                />
                <span>{item.is_node_online ? "Online" : "Cached"}</span>
              </div>
              <div style={{ display: "flex", gap: "0.35rem", alignItems: "center" }}>
                <span
                  style={{
                    fontSize: "0.7rem",
                    padding: "0.15rem 0.35rem",
                    borderRadius: "3px",
                    background: "rgba(245, 158, 11, 0.15)",
                    color: "var(--accent-gold)",
                    fontWeight: 600,
                  }}
                >
                  4s HLS
                </span>
                <span
                  suppressHydrationWarning
                  style={{ fontSize: "0.75rem", color: "var(--text-dim)" }}
                >
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
  );
}
