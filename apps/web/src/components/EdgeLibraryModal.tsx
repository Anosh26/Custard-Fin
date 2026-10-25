"use client";

import React, { useEffect, useState } from "react";
import {
  fetchEdgeMedia,
  toggleMediaPublic,
  EdgeMediaItem,
  formatDuration,
} from "../lib/api-client";

interface EdgeLibraryModalProps {
  isOpen: boolean;
  onClose: () => void;
  onChanged?: () => void;
}

export default function EdgeLibraryModal({
  isOpen,
  onClose,
  onChanged,
}: EdgeLibraryModalProps) {
  const [items, setItems] = useState<EdgeMediaItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [updatingId, setUpdatingId] = useState<string | null>(null);

  const loadMedia = async () => {
    setLoading(true);
    setError(null);
    try {
      const list = await fetchEdgeMedia();
      setItems(list || []);
    } catch (err: unknown) {
      setError(
        err instanceof Error
          ? err.message
          : "Cannot connect to local edge daemon at :8090"
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      loadMedia();
    }
  }, [isOpen]);

  const handleToggle = async (id: string, currentPublic: boolean) => {
    setUpdatingId(id);
    try {
      await toggleMediaPublic(id, !currentPublic);
      setItems((prev) =>
        prev.map((item) =>
          item.id === id ? { ...item, is_public: !currentPublic } : item
        )
      );
      if (onChanged) onChanged();
    } catch (err: unknown) {
      setError(
        err instanceof Error ? err.message : "Failed to toggle public status"
      );
    } finally {
      setUpdatingId(null);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="modal-card wide" onClick={(e) => e.stopPropagation()}>
        {/* Header */}
        <div className="modal-header">
          <div style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
            <div className="modal-icon-badge">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <polygon points="12 2 2 7 12 12 22 7 12 2" />
                <polyline points="2 17 12 22 22 17" />
                <polyline points="2 12 12 17 22 12" />
              </svg>
            </div>
            <div>
              <h3 style={{ fontSize: "1.15rem", fontWeight: 700 }}>
                My Edge Node Media Library
              </h3>
              <p style={{ fontSize: "0.8rem", color: "var(--text-dim)" }}>
                Video files stored in local SQLite database (Strictly on this PC).
              </p>
            </div>
          </div>
          <button className="modal-close-btn" onClick={onClose} aria-label="Close modal">
            &times;
          </button>
        </div>

        {/* Content */}
        {error && (
          <div className="alert-box error" style={{ marginBottom: "1rem" }}>
            {error}
          </div>
        )}

        {loading ? (
          <div style={{ textAlign: "center", padding: "2.5rem 0", color: "var(--text-dim)" }}>
            Loading edge library from SQLite database...
          </div>
        ) : items.length === 0 ? (
          <div style={{ textAlign: "center", padding: "2.5rem 0" }}>
            <p style={{ color: "var(--text-dim)", marginBottom: "1rem" }}>
              No videos found in your local edge library yet.
            </p>
          </div>
        ) : (
          <div className="table-container">
            <table className="dev-table">
              <thead>
                <tr>
                  <th>Title</th>
                  <th>Resolution</th>
                  <th>Duration</th>
                  <th>Storage Location</th>
                  <th>Community Share</th>
                </tr>
              </thead>
              <tbody>
                {items.map((item) => (
                  <tr key={item.id}>
                    <td>
                      <strong style={{ color: "var(--text-main)", fontSize: "0.9rem" }}>
                        {item.title}
                      </strong>
                      <div style={{ fontSize: "0.75rem", color: "var(--text-dim)", fontFamily: "monospace" }}>
                        ID: {item.id.slice(0, 13)}...
                      </div>
                    </td>
                    <td>
                      <span className="badge badge-resolution">
                        {item.width && item.height ? `${item.width}x${item.height}` : "720p"}
                      </span>
                    </td>
                    <td style={{ fontSize: "0.85rem", color: "var(--text-muted)" }}>
                      {formatDuration(item.duration_s)}
                    </td>
                    <td style={{ fontSize: "0.78rem", color: "var(--text-dim)", maxWidth: "160px", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                      {item.path}
                    </td>
                    <td>
                      <button
                        className={`toggle-switch-btn ${item.is_public ? "active" : ""}`}
                        disabled={updatingId === item.id}
                        onClick={() => handleToggle(item.id, item.is_public)}
                        title={item.is_public ? "Public: Manifest synced to cloud" : "Private: Local storage only"}
                      >
                        <span className="toggle-switch-knob" />
                        <span style={{ fontSize: "0.75rem", fontWeight: 600, paddingLeft: item.is_public ? "0.3rem" : "1.5rem" }}>
                          {item.is_public ? "Public" : "Private"}
                        </span>
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div style={{ marginTop: "1.25rem", padding: "0.85rem", background: "rgba(15, 23, 42, 0.4)", borderRadius: "var(--radius-md)", fontSize: "0.8rem", color: "var(--text-dim)", display: "flex", alignItems: "center", gap: "0.5rem" }}>
          <span style={{ color: "var(--accent-gold)", fontWeight: 700 }}>Invariant 1:</span>
          Toggling &quot;Public&quot; syncs a lightweight 3 KB JSON manifest to the Cloud API. The video binary never leaves your PC.
        </div>
      </div>
    </div>
  );
}
