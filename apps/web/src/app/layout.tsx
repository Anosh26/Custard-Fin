import "./globals.css";
import Link from "next/link";
import React, { Suspense } from "react";
import SearchBar from "../components/SearchBar";

export const metadata = {
  title: "CustardFin - Distributed Community Media",
  description:
    "Hybrid personal media streaming and community catalog. Video binaries stay on your machine, catalog lives in the cloud.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body suppressHydrationWarning>
        <nav className="navbar">
          <Link href="/" className="nav-brand">
            <div className="brand-icon">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="#090d16"
                strokeWidth="2.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <polygon points="5 3 19 12 5 21 5 3" />
              </svg>
            </div>
            <span>CustardFin</span>
          </Link>

          <Suspense fallback={<div className="search-box" />}>
            <SearchBar />
          </Suspense>

          <div style={{ display: "flex", alignItems: "center", gap: "1rem" }}>
            <div
              style={{
                fontSize: "0.8rem",
                padding: "0.35rem 0.75rem",
                borderRadius: "9999px",
                background: "rgba(245, 158, 11, 0.1)",
                color: "var(--accent-gold)",
                border: "1px solid rgba(245, 158, 11, 0.25)",
                fontWeight: 600,
              }}
            >
              MCA Cloud Lab
            </div>
          </div>
        </nav>

        <main className="container">{children}</main>

        <footer
          style={{
            borderTop: "1px solid var(--border-subtle)",
            padding: "2rem",
            textAlign: "center",
            fontSize: "0.85rem",
            color: "var(--text-dim)",
            background: "rgba(15, 23, 42, 0.5)",
            backdropFilter: "blur(8px)",
          }}
        >
          <div style={{ marginBottom: "0.5rem" }}>
            CustardFin: Edge-first distributed media sharing. Video binaries strictly remain at home.
          </div>
          <div>
            MCA Cloud Computing Lab Project &bull; Anish Ghadi (25P0320009) &amp; Anosh Chodankar (25P0320010)
          </div>
        </footer>
      </body>
    </html>
  );
}
