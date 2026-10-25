"use client";

import React, { useEffect, useState } from "react";
import ModalPortal from "./ModalPortal";
import {
  uploadEdgeMedia,
  triggerEdgeScan,
  generateSampleClip,
} from "../lib/api-client";

interface AddVideoModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

export default function AddVideoModal({
  isOpen,
  onClose,
  onSuccess,
}: AddVideoModalProps) {
  const [tab, setTab] = useState<"upload" | "scan" | "sample">("upload");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadPercent, setUploadPercent] = useState<number>(0);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);

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

  if (!isOpen) return null;

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
      setStatusMessage(null);
      setErrorMessage(null);
    }
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      setSelectedFile(e.dataTransfer.files[0]);
      setStatusMessage(null);
      setErrorMessage(null);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) return;
    setIsUploading(true);
    setUploadPercent(0);
    setStatusMessage(null);
    setErrorMessage(null);

    try {
      const res = await uploadEdgeMedia(selectedFile, (pct) => {
        setUploadPercent(pct);
      });
      setStatusMessage(
        res.message || "Video uploaded! Transcoding to 4s HLS and syncing manifest."
      );
      setSelectedFile(null);
      if (onSuccess) {
        setTimeout(() => onSuccess(), 1200);
      }
    } catch (err: unknown) {
      setErrorMessage(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setIsUploading(false);
    }
  };

  const handleScanDirectory = async () => {
    setIsUploading(true);
    setStatusMessage(null);
    setErrorMessage(null);
    try {
      const res = await triggerEdgeScan();
      setStatusMessage(
        `Directory scan complete: ${res.total} total items detected (${res.added} added). Untranscoded items are encoding.`
      );
      if (onSuccess) {
        setTimeout(() => onSuccess(), 1200);
      }
    } catch (err: unknown) {
      setErrorMessage(err instanceof Error ? err.message : "Scan failed");
    } finally {
      setIsUploading(false);
    }
  };

  const handleGenerateSample = async () => {
    setIsUploading(true);
    setStatusMessage(null);
    setErrorMessage(null);
    try {
      const res = await generateSampleClip();
      setStatusMessage(
        res.message || "Sample test video synthesized, transcoded, and published to catalog!"
      );
      if (onSuccess) {
        setTimeout(() => onSuccess(), 1200);
      }
    } catch (err: unknown) {
      setErrorMessage(
        err instanceof Error ? err.message : "Failed generating sample clip"
      );
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <ModalPortal>
      <div
        className="modal-backdrop"
        onClick={onClose}
        onDragOver={(e) => e.preventDefault()}
        onDrop={(e) => e.preventDefault()}
      >
        <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        {/* Modal Header */}
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
                <path d="M12 5v14M5 12h14" />
              </svg>
            </div>
            <div>
              <h3 style={{ fontSize: "1.15rem", fontWeight: 700 }}>Add Video to Edge Node</h3>
              <p style={{ fontSize: "0.8rem", color: "var(--text-dim)" }}>
                Video stays on your machine; only a ~3 KB manifest syncs to the cloud.
              </p>
            </div>
          </div>
          <button className="modal-close-btn" onClick={onClose} aria-label="Close modal">
            &times;
          </button>
        </div>

        {/* Tab Switcher */}
        <div className="modal-tabs">
          <button
            className={`modal-tab ${tab === "upload" ? "active" : ""}`}
            onClick={() => {
              setTab("upload");
              setStatusMessage(null);
              setErrorMessage(null);
            }}
          >
            Upload File
          </button>
          <button
            className={`modal-tab ${tab === "scan" ? "active" : ""}`}
            onClick={() => {
              setTab("scan");
              setStatusMessage(null);
              setErrorMessage(null);
            }}
          >
            Scan Folder
          </button>
          <button
            className={`modal-tab ${tab === "sample" ? "active" : ""}`}
            onClick={() => {
              setTab("sample");
              setStatusMessage(null);
              setErrorMessage(null);
            }}
          >
            1-Click Sample Clip
          </button>
        </div>

        {/* Tab 1: File Upload */}
        {tab === "upload" && (
          <div>
            <div
              className={`dropzone ${isDragging ? "dragging" : ""}`}
              onDragOver={(e) => {
                e.preventDefault();
                setIsDragging(true);
              }}
              onDragLeave={() => setIsDragging(false)}
              onDrop={handleDrop}
            >
              <input
                type="file"
                id="video-upload-input"
                accept="video/mp4,video/mkv,video/webm,video/quicktime,video/x-msvideo"
                onChange={handleFileChange}
                style={{ display: "none" }}
              />
              <label htmlFor="video-upload-input" style={{ cursor: "pointer", display: "block" }}>
                <div style={{ marginBottom: "0.75rem", color: "var(--accent-gold)" }}>
                  <svg
                    width="36"
                    height="36"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="1.75"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    style={{ margin: "0 auto", display: "block" }}
                  >
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                    <polyline points="17 8 12 3 7 8" />
                    <line x1="12" y1="3" x2="12" y2="15" />
                  </svg>
                </div>
                {selectedFile ? (
                  <div>
                    <strong style={{ color: "var(--text-main)", fontSize: "0.95rem" }}>
                      {selectedFile.name}
                    </strong>
                    <div style={{ fontSize: "0.8rem", color: "var(--text-muted)", marginTop: "0.25rem" }}>
                      {(selectedFile.size / (1024 * 1024)).toFixed(1)} MB &bull; Ready to upload
                    </div>
                  </div>
                ) : (
                  <div>
                    <strong style={{ color: "var(--text-main)", fontSize: "0.95rem" }}>
                      Drop your video file here, or click to browse
                    </strong>
                    <div style={{ fontSize: "0.8rem", color: "var(--text-dim)", marginTop: "0.25rem" }}>
                      Supports MP4, MKV, MOV, WEBM (Hardware NVENC accelerated)
                    </div>
                  </div>
                )}
              </label>
            </div>

            {isUploading && uploadPercent > 0 && uploadPercent < 100 && (
              <div style={{ marginTop: "1rem" }}>
                <div className="progress-bar-bg">
                  <div className="progress-bar-fill" style={{ width: `${uploadPercent}%` }} />
                </div>
                <div style={{ fontSize: "0.75rem", color: "var(--text-dim)", textAlign: "right", marginTop: "0.25rem" }}>
                  Uploading to Home Edge: {uploadPercent}%
                </div>
              </div>
            )}

            <button
              className="action-btn primary"
              style={{ width: "100%", marginTop: "1.25rem" }}
              disabled={!selectedFile || isUploading}
              onClick={handleUpload}
            >
              {isUploading ? "Uploading & Transcoding..." : "Upload & Publish to Edge Library"}
            </button>
          </div>
        )}

        {/* Tab 2: Folder Scan */}
        {tab === "scan" && (
          <div style={{ padding: "0.5rem 0" }}>
            <p style={{ fontSize: "0.88rem", color: "var(--text-muted)", lineHeight: 1.6, marginBottom: "1.25rem" }}>
              Have files already on your hard drive? Simply drop or copy your video files into the Edge daemon&apos;s media folder:
            </p>
            <div className="code-pill" style={{ marginBottom: "1.5rem" }}>
              ./media/your_video.mp4
            </div>
            <p style={{ fontSize: "0.82rem", color: "var(--text-dim)", marginBottom: "1.5rem" }}>
              Clicking below triggers the Edge daemon to hash new files, extract codec metadata via FFprobe, generate 4-second HLS segments, and publish manifests to the Cloud catalog.
            </p>
            <button
              className="action-btn primary"
              style={{ width: "100%" }}
              disabled={isUploading}
              onClick={handleScanDirectory}
            >
              {isUploading ? "Scanning Directory..." : "Scan Media Folder Now"}
            </button>
          </div>
        )}

        {/* Tab 3: Sample Test Clip */}
        {tab === "sample" && (
          <div style={{ padding: "0.5rem 0" }}>
            <p style={{ fontSize: "0.88rem", color: "var(--text-muted)", lineHeight: 1.6, marginBottom: "1.25rem" }}>
              Need to test or present right away without downloading a video? Synthesize a clean 720p 30fps test card clip with stereo audio directly using FFmpeg.
            </p>
            <div
              style={{
                background: "rgba(15, 23, 42, 0.6)",
                border: "1px solid var(--border-subtle)",
                borderRadius: "var(--radius-md)",
                padding: "1rem",
                marginBottom: "1.5rem",
              }}
            >
              <div style={{ display: "flex", justifyContent: "space-between", fontSize: "0.8rem", color: "var(--text-muted)", marginBottom: "0.35rem" }}>
                <span>Resolution:</span>
                <span style={{ color: "var(--text-main)", fontWeight: 600 }}>1280x720 (720p)</span>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between", fontSize: "0.8rem", color: "var(--text-muted)", marginBottom: "0.35rem" }}>
                <span>Duration:</span>
                <span style={{ color: "var(--text-main)", fontWeight: 600 }}>15 seconds</span>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between", fontSize: "0.8rem", color: "var(--text-muted)" }}>
                <span>Encoding:</span>
                <span style={{ color: "var(--accent-gold)", fontWeight: 600 }}>H.264 4s HLS chunks</span>
              </div>
            </div>
            <button
              className="action-btn gold"
              style={{ width: "100%" }}
              disabled={isUploading}
              onClick={handleGenerateSample}
            >
              {isUploading ? "Synthesizing & Encoding..." : "Generate & Publish Sample Test Clip"}
            </button>
          </div>
        )}

        {/* Status / Error feedback */}
        {statusMessage && (
          <div className="alert-box success" style={{ marginTop: "1rem" }}>
            {statusMessage}
          </div>
        )}
        {errorMessage && (
          <div className="alert-box error" style={{ marginTop: "1rem" }}>
            {errorMessage}
          </div>
        )}
      </div>
    </div>
    </ModalPortal>
  );
}
