"use client";

import React, { useEffect, useState } from "react";
import AddVideoModal from "./AddVideoModal";
import EdgeLibraryModal from "./EdgeLibraryModal";
import DevSettingsModal from "./DevSettingsModal";
import AuthModal from "./AuthModal";
import { getCurrentUser, fetchCurrentUser, UserAccount } from "../lib/api-client";

export default function NavbarClientControls() {
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [isLibraryOpen, setIsLibraryOpen] = useState(false);
  const [isDevOpen, setIsDevOpen] = useState(false);
  const [isAuthOpen, setIsAuthOpen] = useState(false);
  const [currentUser, setCurrentUser] = useState<UserAccount | null>(null);

  useEffect(() => {
    // Initial user load
    setCurrentUser(getCurrentUser());
    fetchCurrentUser().then((u) => {
      if (u) setCurrentUser(u);
    });

    const handleAuthChange = () => {
      setCurrentUser(getCurrentUser());
    };

    window.addEventListener("custardfin_auth_change", handleAuthChange);
    return () => {
      window.removeEventListener("custardfin_auth_change", handleAuthChange);
    };
  }, []);

  const handleRefresh = () => {
    if (typeof window !== "undefined") {
      window.location.reload();
    }
  };

  return (
    <div style={{ display: "flex", alignItems: "center", gap: "0.6rem" }}>
      {/* User Account / Sign In Button */}
      {currentUser ? (
        <button
          className="nav-btn ghost"
          onClick={() => setIsAuthOpen(true)}
          title={`Signed in as @${currentUser.username} (Click to manage account)`}
          style={{
            borderColor: "rgba(245, 158, 11, 0.4)",
            background: "rgba(245, 158, 11, 0.08)",
            color: "var(--accent-gold)",
          }}
        >
          <span
            style={{
              width: "20px",
              height: "20px",
              borderRadius: "50%",
              background: "rgba(245, 158, 11, 0.25)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              fontSize: "0.75rem",
              fontWeight: 700,
            }}
          >
            {currentUser.username.slice(0, 1).toUpperCase()}
          </span>
          <span style={{ fontWeight: 600 }}>@{currentUser.username}</span>
        </button>
      ) : (
        <button
          className="nav-btn ghost"
          onClick={() => setIsAuthOpen(true)}
          title="Sign in to CustardFin Cloud to access private videos"
          style={{
            borderColor: "rgba(245, 158, 11, 0.35)",
            color: "var(--accent-gold)",
          }}
        >
          <svg
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>
          <span>Sign In</span>
        </button>
      )}

      {/* Add Video Button */}
      <button
        className="nav-btn primary"
        onClick={() => setIsAddOpen(true)}
        title="Add or upload a video to your home edge node"
      >
        <svg
          width="15"
          height="15"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M12 5v14M5 12h14" />
        </svg>
        <span>Add Video</span>
      </button>

      {/* Edge Library Button */}
      <button
        className="nav-btn secondary"
        onClick={() => setIsLibraryOpen(true)}
        title="Manage your local edge library and community sharing"
      >
        <svg
          width="15"
          height="15"
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
        <span>My Library</span>
      </button>

      {/* Dev Settings Button */}
      <button
        className="nav-btn ghost"
        onClick={() => setIsDevOpen(true)}
        title="Open Dev Console & Architecture Settings"
      >
        <svg
          width="15"
          height="15"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <circle cx="12" cy="12" r="3" />
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
        </svg>
        <span>Dev Console</span>
      </button>

      {/* Modals */}
      <AuthModal
        isOpen={isAuthOpen}
        onClose={() => setIsAuthOpen(false)}
        onSuccess={handleRefresh}
      />

      <AddVideoModal
        isOpen={isAddOpen}
        onClose={() => setIsAddOpen(false)}
        onSuccess={handleRefresh}
      />

      <EdgeLibraryModal
        isOpen={isLibraryOpen}
        onClose={() => setIsLibraryOpen(false)}
        onChanged={handleRefresh}
      />

      <DevSettingsModal
        isOpen={isDevOpen}
        onClose={() => setIsDevOpen(false)}
        onConfigChanged={handleRefresh}
      />
    </div>
  );
}
