"use client";

import React, { useEffect, useState } from "react";
import ModalPortal from "./ModalPortal";
import {
  loginUser,
  registerUser,
  logoutUser,
  getCurrentUser,
  UserAccount,
} from "../lib/api-client";

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

export default function AuthModal({
  isOpen,
  onClose,
  onSuccess,
}: AuthModalProps) {
  const [currentUser, setCurrentUser] = useState<UserAccount | null>(null);
  const [mode, setMode] = useState<"login" | "register">("login");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      setCurrentUser(getCurrentUser());
      setError(null);
    }
  }, [isOpen]);

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !password.trim()) {
      setError("Username and password are required");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      if (mode === "login") {
        await loginUser(username.trim(), password.trim());
      } else {
        await registerUser(username.trim(), password.trim());
      }
      onClose();
      if (onSuccess) onSuccess();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Authentication failed");
    } finally {
      setLoading(false);
    }
  };

  const handleDemoLogin = async () => {
    setLoading(true);
    setError(null);
    try {
      await loginUser("demo", "demo1234");
      onClose();
      if (onSuccess) onSuccess();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Demo login failed");
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logoutUser();
    setCurrentUser(null);
    onClose();
    if (onSuccess) onSuccess();
  };

  if (!isOpen) return null;

  return (
    <ModalPortal>
      <div className="modal-backdrop" onClick={onClose}>
        <div
          className="modal-card"
          style={{ maxWidth: "440px" }}
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="modal-header">
            <div style={{ display: "flex", alignItems: "center", gap: "0.6rem" }}>
              <div
                style={{
                  width: "36px",
                  height: "36px",
                  borderRadius: "var(--radius-sm)",
                  background: "linear-gradient(135deg, var(--accent-gold), #ea580c)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  color: "#000",
                }}
              >
                <svg
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                  <circle cx="12" cy="7" r="4" />
                </svg>
              </div>
              <div>
                <h2 style={{ fontSize: "1.15rem", fontWeight: 700 }}>
                  {currentUser ? "User Account" : mode === "login" ? "Sign In" : "Create Account"}
                </h2>
                <p style={{ fontSize: "0.75rem", color: "var(--text-dim)" }}>
                  {currentUser
                    ? "Manage your CustardFin cloud identity"
                    : "Access your private videos & uploaded catalogs"}
                </p>
              </div>
            </div>
            <button className="modal-close-btn" onClick={onClose}>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </div>

          {currentUser ? (
            /* Logged In View */
            <div style={{ padding: "0.5rem 0" }}>
              <div
                style={{
                  background: "rgba(15, 23, 42, 0.6)",
                  border: "1px solid var(--border-subtle)",
                  borderRadius: "var(--radius-md)",
                  padding: "1.25rem",
                  marginBottom: "1.5rem",
                }}
              >
                <div style={{ display: "flex", alignItems: "center", gap: "0.85rem", marginBottom: "0.75rem" }}>
                  <div
                    style={{
                      width: "44px",
                      height: "44px",
                      borderRadius: "50%",
                      background: "rgba(245, 158, 11, 0.2)",
                      border: "1px solid var(--accent-gold)",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                      color: "var(--accent-gold)",
                      fontWeight: 700,
                      fontSize: "1.2rem",
                    }}
                  >
                    {currentUser.username.slice(0, 1).toUpperCase()}
                  </div>
                  <div>
                    <div style={{ fontSize: "1.1rem", fontWeight: 700, color: "var(--text-main)" }}>
                      @{currentUser.username}
                    </div>
                    <div style={{ fontSize: "0.75rem", color: "var(--accent-emerald)", display: "flex", alignItems: "center", gap: "0.3rem" }}>
                      <span className="status-dot online" style={{ width: "6px", height: "6px" }} />
                      Authenticated Session
                    </div>
                  </div>
                </div>

                <div style={{ fontSize: "0.78rem", color: "var(--text-dim)", fontFamily: "monospace", wordBreak: "break-all" }}>
                  User ID: {currentUser.user_id}
                </div>
              </div>

              <div style={{ display: "flex", gap: "0.75rem", justifyContent: "flex-end" }}>
                <button
                  type="button"
                  className="nav-btn secondary"
                  onClick={onClose}
                  style={{ padding: "0.6rem 1.25rem" }}
                >
                  Close
                </button>
                <button
                  type="button"
                  className="nav-btn ghost"
                  onClick={handleLogout}
                  style={{
                    padding: "0.6rem 1.25rem",
                    color: "var(--accent-red)",
                    borderColor: "rgba(239, 68, 68, 0.4)",
                  }}
                >
                  Log Out
                </button>
              </div>
            </div>
          ) : (
            /* Logged Out / Auth Form */
            <div>
              {/* Tab Switcher */}
              <div
                style={{
                  display: "flex",
                  background: "rgba(15, 23, 42, 0.6)",
                  padding: "0.25rem",
                  borderRadius: "var(--radius-sm)",
                  marginBottom: "1.25rem",
                  border: "1px solid var(--border-subtle)",
                }}
              >
                <button
                  type="button"
                  onClick={() => {
                    setMode("login");
                    setError(null);
                  }}
                  style={{
                    flex: 1,
                    padding: "0.45rem",
                    borderRadius: "4px",
                    border: "none",
                    background: mode === "login" ? "rgba(245, 158, 11, 0.2)" : "transparent",
                    color: mode === "login" ? "var(--accent-gold)" : "var(--text-muted)",
                    fontWeight: 600,
                    fontSize: "0.85rem",
                    cursor: "pointer",
                    transition: "all 0.15s ease",
                  }}
                >
                  Sign In
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setMode("register");
                    setError(null);
                  }}
                  style={{
                    flex: 1,
                    padding: "0.45rem",
                    borderRadius: "4px",
                    border: "none",
                    background: mode === "register" ? "rgba(245, 158, 11, 0.2)" : "transparent",
                    color: mode === "register" ? "var(--accent-gold)" : "var(--text-muted)",
                    fontWeight: 600,
                    fontSize: "0.85rem",
                    cursor: "pointer",
                    transition: "all 0.15s ease",
                  }}
                >
                  Create Account
                </button>
              </div>

              {error && (
                <div
                  style={{
                    padding: "0.75rem",
                    borderRadius: "var(--radius-sm)",
                    background: "rgba(239, 68, 68, 0.15)",
                    border: "1px solid rgba(239, 68, 68, 0.3)",
                    color: "var(--accent-red)",
                    fontSize: "0.85rem",
                    marginBottom: "1rem",
                  }}
                >
                  {error}
                </div>
              )}

              <form onSubmit={handleSubmit}>
                <div style={{ marginBottom: "1rem" }}>
                  <label
                    style={{
                      display: "block",
                      fontSize: "0.8rem",
                      fontWeight: 600,
                      color: "var(--text-muted)",
                      marginBottom: "0.35rem",
                    }}
                  >
                    Username
                  </label>
                  <input
                    type="text"
                    required
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="e.g. demo"
                    style={{
                      width: "100%",
                      padding: "0.6rem 0.85rem",
                      borderRadius: "var(--radius-sm)",
                      border: "1px solid var(--border-subtle)",
                      background: "rgba(15, 23, 42, 0.8)",
                      color: "var(--text-main)",
                      fontSize: "0.9rem",
                    }}
                  />
                </div>

                <div style={{ marginBottom: "1.25rem" }}>
                  <label
                    style={{
                      display: "block",
                      fontSize: "0.8rem",
                      fontWeight: 600,
                      color: "var(--text-muted)",
                      marginBottom: "0.35rem",
                    }}
                  >
                    Password
                  </label>
                  <input
                    type="password"
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••"
                    style={{
                      width: "100%",
                      padding: "0.6rem 0.85rem",
                      borderRadius: "var(--radius-sm)",
                      border: "1px solid var(--border-subtle)",
                      background: "rgba(15, 23, 42, 0.8)",
                      color: "var(--text-main)",
                      fontSize: "0.9rem",
                    }}
                  />
                </div>

                <div style={{ display: "flex", flexDirection: "column", gap: "0.6rem" }}>
                  <button
                    type="submit"
                    disabled={loading}
                    className="nav-btn primary"
                    style={{
                      width: "100%",
                      justifyContent: "center",
                      padding: "0.7rem",
                      fontSize: "0.9rem",
                    }}
                  >
                    {loading
                      ? "Processing..."
                      : mode === "login"
                      ? "Sign In to Cloud API"
                      : "Create Account"}
                  </button>

                  <button
                    type="button"
                    disabled={loading}
                    onClick={handleDemoLogin}
                    style={{
                      width: "100%",
                      padding: "0.65rem",
                      borderRadius: "var(--radius-sm)",
                      border: "1px dashed rgba(245, 158, 11, 0.4)",
                      background: "rgba(245, 158, 11, 0.08)",
                      color: "var(--accent-gold)",
                      fontSize: "0.85rem",
                      fontWeight: 600,
                      cursor: "pointer",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                      gap: "0.4rem",
                      transition: "all 0.15s ease",
                    }}
                  >
                    <span>⚡ Quick Demo Login (as @demo)</span>
                  </button>
                </div>
              </form>
            </div>
          )}
        </div>
      </div>
    </ModalPortal>
  );
}
