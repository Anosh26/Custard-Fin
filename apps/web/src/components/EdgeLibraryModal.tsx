"use client";

import React, { useEffect, useState } from "react";
import ModalPortal from "./ModalPortal";
import AuthModal from "./AuthModal";
import {
  fetchEdgeMedia,
  toggleMediaPublic,
  deleteEdgeMedia,
  EdgeMediaItem,
  formatDuration,
  getCurrentUser,
  getAuthToken,
  fetchNodeStatus,
  claimNode,
  CloudNodeStatus,
  UserAccount,
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
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [nodeStatus, setNodeStatus] = useState<CloudNodeStatus | null>(null);
  const [currentUser, setCurrentUser] = useState<UserAccount | null>(null);
  const [claiming, setClaiming] = useState(false);
  const [claimSuccess, setClaimSuccess] = useState<string | null>(null);
  const [isAuthOpen, setIsAuthOpen] = useState(false);

  // Lock body scroll and register Escape key listener while modal is open
  useEffect(() => {
    if (!isOpen) return;

    const originalOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => {
      document.body.style.overflow = originalOverflow;
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isOpen, onClose]);

  const loadMedia = async () => {
    setLoading(true);
    setError(null);
    const user = getCurrentUser();
    setCurrentUser(user);

    try {
      const [listResult, statusResult] = await Promise.allSettled([
        fetchEdgeMedia(),
        fetchNodeStatus("898ea0c7-f7a3-48e8-ba59-b378a219f9ee"),
      ]);
      if (listResult.status === "fulfilled") {
        setItems(listResult.value || []);
      } else {
        setError("Cannot connect to local edge daemon at :8090");
      }
      if (statusResult.status === "fulfilled") {
        setNodeStatus(statusResult.value);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      loadMedia();
    }
  }, [isOpen]);

  const handleClaim = async () => {
    const token = getAuthToken();
    if (!token || !currentUser) {
      setIsAuthOpen(true);
      return;
    }
    setClaiming(true);
    setError(null);
    try {
      const res = await claimNode("898ea0c7-f7a3-48e8-ba59-b378a219f9ee");
      setClaimSuccess(res.message || `Node linked to @${currentUser.username}`);
      if (nodeStatus) {
        setNodeStatus({
          ...nodeStatus,
          owner_username: currentUser.username,
          owner_user_id: currentUser.user_id,
        });
      }
      if (onChanged) onChanged();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to claim edge node";
      if (
        msg.includes("Authorization") ||
        msg.includes("token") ||
        msg.includes("Session") ||
        msg.includes("expired") ||
        msg.includes("UNAUTHORIZED")
      ) {
        setError("Session expired or authentication required. Please sign in.");
        setIsAuthOpen(true);
      } else {
        setError(msg);
      }
    } finally {
      setClaiming(false);
    }
  };

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

  const handleDelete = async (id: string, title: string) => {
    const confirmed = window.confirm(
      `Remove "${title}" from your library?\n\nThis will remove it from your edge node, unsync it from the cloud catalog, and clean up its transcoded HLS files.`
    );
    if (!confirmed) return;

    setDeletingId(id);
    setError(null);
    try {
      await deleteEdgeMedia(id, true);
      setItems((prev) => prev.filter((item) => item.id !== id));
      if (onChanged) onChanged();
    } catch (err: unknown) {
      setError(
        err instanceof Error ? err.message : "Failed to remove video from library"
      );
    } finally {
      setDeletingId(null);
    }
  };

  if (!isOpen) return null;

  return (
    <ModalPortal>
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

        {/* Node Ownership & Attribution Banner */}
        {currentUser ? (
          nodeStatus?.owner_username && nodeStatus.owner_username !== currentUser.username ? (
            <div
              style={{
                marginBottom: "1rem",
                padding: "0.75rem 1rem",
                borderRadius: "var(--radius-md)",
                background: "rgba(245, 158, 11, 0.12)",
                border: "1px solid rgba(245, 158, 11, 0.35)",
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                gap: "1rem",
              }}
            >
              <div>
                <div style={{ fontSize: "0.85rem", color: "var(--accent-gold)", fontWeight: 700 }}>
                  Edge Node Linked to @{nodeStatus.owner_username}
                </div>
                <div style={{ fontSize: "0.75rem", color: "var(--text-muted)" }}>
                  Videos uploaded here appear in the Community Catalog under @{nodeStatus.owner_username}.
                </div>
              </div>
              <button
                className="nav-btn primary"
                onClick={handleClaim}
                disabled={claiming}
                style={{ padding: "0.4rem 0.85rem", fontSize: "0.8rem", whiteSpace: "nowrap" }}
              >
                {claiming ? "Linking..." : `Link to @${currentUser.username}`}
              </button>
            </div>
          ) : (
            <div
              style={{
                marginBottom: "1rem",
                padding: "0.55rem 0.85rem",
                borderRadius: "var(--radius-md)",
                background: "rgba(16, 185, 129, 0.1)",
                border: "1px solid rgba(16, 185, 129, 0.25)",
                display: "flex",
                alignItems: "center",
                gap: "0.5rem",
                fontSize: "0.8rem",
                color: "var(--accent-green)",
              }}
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                <polyline points="20 6 9 17 4 12" />
              </svg>
              <span>
                {claimSuccess || (
                  <>
                    Edge node is linked to your account <strong>@{currentUser.username}</strong>. All catalog uploads belong to you.
                  </>
                )}
              </span>
            </div>
          )
        ) : (
          <div
            style={{
              marginBottom: "1rem",
              padding: "0.55rem 0.85rem",
              borderRadius: "var(--radius-md)",
              background: "rgba(59, 130, 246, 0.08)",
              border: "1px solid rgba(59, 130, 246, 0.2)",
              fontSize: "0.78rem",
              color: "var(--text-dim)",
            }}
          >
            💡 Sign in with your account to link this edge node and attribute your uploads to your username.
          </div>
        )}

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
                  <th style={{ textAlign: "right" }}>Actions</th>
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
                    <td style={{ textAlign: "right" }}>
                      <div style={{ display: "flex", gap: "0.4rem", justifyContent: "flex-end" }}>
                        <a
                          href={`/watch/${item.id}`}
                          style={{
                            display: "inline-flex",
                            alignItems: "center",
                            gap: "0.3rem",
                            padding: "0.35rem 0.65rem",
                            borderRadius: "var(--radius-sm)",
                            fontSize: "0.75rem",
                            fontWeight: 600,
                            background: "rgba(245, 158, 11, 0.15)",
                            color: "var(--accent-gold)",
                            border: "1px solid rgba(245, 158, 11, 0.3)",
                            textDecoration: "none",
                            transition: "all 0.15s ease",
                          }}
                          title="Watch video"
                        >
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                            <polygon points="5 3 19 12 5 21 5 3" />
                          </svg>
                          <span>Watch</span>
                        </a>
                        <button
                          className="btn-delete-item"
                          disabled={deletingId === item.id}
                          onClick={() => handleDelete(item.id, item.title)}
                          title="Remove video from library & cloud catalog"
                        >
                          <svg
                            width="14"
                            height="14"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="2"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                          >
                            <polyline points="3 6 5 6 21 6" />
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                            <line x1="10" y1="11" x2="10" y2="17" />
                            <line x1="14" y1="11" x2="14" y2="17" />
                          </svg>
                          {deletingId === item.id ? "Removing..." : "Remove"}
                        </button>
                      </div>
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

    <AuthModal
      isOpen={isAuthOpen}
      onClose={() => setIsAuthOpen(false)}
      onSuccess={() => {
        setIsAuthOpen(false);
        const user = getCurrentUser();
        setCurrentUser(user);
        handleClaim();
      }}
    />
    </ModalPortal>
  );
}
