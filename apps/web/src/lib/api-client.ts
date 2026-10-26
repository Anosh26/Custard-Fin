export interface CatalogItem {
  media_id: string;
  title: string;
  description: string;
  duration_s: number;
  node_id: string;
  node_name?: string;
  uploader_username?: string;
  is_public: boolean;
  is_node_online: boolean;
  updated_at: string;
}

export interface CatalogResponse {
  items: CatalogItem[];
  total: number;
}

export interface PlaybackResponse {
  stream_url: string;
  route_type: "direct_lan" | "origin_shield" | "cache_only";
  expires_at: string;
  media_id?: string;
  title?: string;
  uploader_username?: string;
  node_name?: string;
  is_public?: boolean;
}

export interface UserAccount {
  user_id: string;
  username: string;
}

export interface EdgeMediaItem {
  id: string;
  title: string;
  path: string;
  content_hash: string;
  duration_s: number;
  width: number;
  height: number;
  is_public: boolean;
  created_at: string;
  updated_at: string;
}

export interface EdgeNodeStatus {
  node_id: string;
  status: string;
  media_dir: string;
  transcode_dir: string;
  media_count: number;
  active_streams: number;
  max_streams: number;
  nvenc_enabled: boolean;
}

export function getAuthToken(): string | null {
  if (typeof window !== "undefined") {
    const t = localStorage.getItem("custardfin_jwt_token");
    if (!t || t === "undefined" || t === "null") return null;
    return t;
  }
  return null;
}

export function setAuthToken(token: string | null) {
  if (typeof window !== "undefined") {
    if (token && token !== "undefined" && token !== "null") {
      localStorage.setItem("custardfin_jwt_token", token);
    } else {
      localStorage.removeItem("custardfin_jwt_token");
    }
  }
}

export function getCurrentUser(): UserAccount | null {
  if (typeof window !== "undefined") {
    const raw = localStorage.getItem("custardfin_user");
    if (raw) {
      try {
        return JSON.parse(raw);
      } catch {
        return null;
      }
    }
  }
  return null;
}

export function setCurrentUser(user: UserAccount | null) {
  if (typeof window !== "undefined") {
    if (user) {
      localStorage.setItem("custardfin_user", JSON.stringify(user));
    } else {
      localStorage.removeItem("custardfin_user");
    }
  }
}

export async function loginUser(
  username: string,
  password: string
): Promise<{ token: string; user_id: string; username: string }> {
  const baseUrl = getApiBaseUrl();
  const res = await fetch(`${baseUrl}/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || "Invalid username or password");
  }
  const data = await res.json();
  const token = data.access_token || data.token;
  setAuthToken(token);
  setCurrentUser({ user_id: data.user_id, username: data.username });
  if (typeof window !== "undefined") {
    window.dispatchEvent(new Event("custardfin_auth_change"));
  }
  return { token, user_id: data.user_id, username: data.username };
}

export async function registerUser(
  username: string,
  password: string
): Promise<{ token: string; user_id: string; username: string }> {
  const baseUrl = getApiBaseUrl();
  const res = await fetch(`${baseUrl}/v1/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || "Registration failed");
  }
  const data = await res.json();
  const token = data.access_token || data.token;
  setAuthToken(token);
  setCurrentUser({ user_id: data.user_id, username: data.username });
  if (typeof window !== "undefined") {
    window.dispatchEvent(new Event("custardfin_auth_change"));
  }
  return { token, user_id: data.user_id, username: data.username };
}

export function logoutUser() {
  setAuthToken(null);
  setCurrentUser(null);
  if (typeof window !== "undefined") {
    window.dispatchEvent(new Event("custardfin_auth_change"));
  }
}

export async function fetchCurrentUser(): Promise<UserAccount | null> {
  const token = getAuthToken();
  if (!token) return null;
  const baseUrl = getApiBaseUrl();
  try {
    const res = await fetch(`${baseUrl}/v1/auth/me`, {
      headers: { Authorization: `Bearer ${token}` },
      cache: "no-store",
    });
    if (!res.ok) {
      logoutUser();
      return null;
    }
    const data = await res.json();
    const user: UserAccount = {
      user_id: data.user_id || data.id,
      username: data.username,
    };
    setCurrentUser(user);
    return user;
  } catch {
    return getCurrentUser();
  }
}

export function getThumbnailUrl(mediaId: string, token?: string | null): string {
  const query = token ? `?token=${encodeURIComponent(token)}` : "";
  return `/v1/media/${mediaId}/thumbnail${query}`;
}

export function getApiBaseUrl(): string {
  if (typeof window !== "undefined") {
    const saved = localStorage.getItem("custardfin_cloud_api_url");
    if (saved && saved.trim()) return saved.trim().replace(/\/v1\/?$/, "");
  }
  const url = process.env.NEXT_PUBLIC_API_URL || "http://api:8080";
  return url.trim().replace(/\/v1\/?$/, "");
}

export function setApiBaseUrl(url: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem("custardfin_cloud_api_url", url.trim());
  }
}

export function getEdgeBaseUrl(): string {
  if (typeof window !== "undefined") {
    const saved = localStorage.getItem("custardfin_edge_api_url");
    if (saved && saved.trim()) return saved.trim();
  }
  return "http://localhost:8090";
}

export function setEdgeBaseUrl(url: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem("custardfin_edge_api_url", url.trim());
  }
}

export async function fetchCatalog(
  q?: string,
  limit = 24,
  offset = 0
): Promise<CatalogResponse> {
  const baseUrl = getApiBaseUrl();
  const params = new URLSearchParams();
  if (q && q.trim()) params.set("q", q.trim());
  params.set("limit", limit.toString());
  params.set("offset", offset.toString());

  const headers: Record<string, string> = {};
  const token = getAuthToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${baseUrl}/v1/catalog?${params.toString()}`, {
    cache: "no-store",
    headers,
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch catalog: ${res.statusText}`);
  }

  return res.json();
}

export async function fetchPlayback(
  mediaId: string
): Promise<PlaybackResponse> {
  const baseUrl = getApiBaseUrl();
  const headers: Record<string, string> = {};
  const token = getAuthToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${baseUrl}/v1/media/${mediaId}/playback`, {
    cache: "no-store",
    headers,
  });

  if (!res.ok) {
    if (res.status === 403) {
      throw new Error("Private video: requires owner sign-in to watch");
    }
    if (res.status === 429) {
      throw new Error("Node is busy (capacity limit exceeded: 10 active streams)");
    }
    if (res.status === 503) {
      throw new Error("Home edge node is currently offline");
    }
    throw new Error(`Failed to resolve playback route: ${res.statusText}`);
  }

  return res.json();
}

export async function fetchEdgeStatus(): Promise<EdgeNodeStatus> {
  const edgeUrl = getEdgeBaseUrl();
  try {
    const res = await fetch(`${edgeUrl}/api/status`, { cache: "no-store" });
    if (res.ok) {
      return await res.json();
    }
  } catch {
    // Local direct fetch blocked (e.g. Private Network Access security or remote viewer)
  }

  // Seamless fallback to Cloud Redis Heartbeat
  const cloudUrl = getApiBaseUrl();
  const res = await fetch(`${cloudUrl}/v1/nodes/demo/status`, { cache: "no-store" });
  if (!res.ok) throw new Error(`Edge node status failed: ${res.statusText}`);
  const data = await res.json();
  return {
    node_id: data.node_id || "898ea0c7-f7a3-48e8-ba59-b378a219f9ee",
    status: data.status || (data.is_online ? "online" : "offline"),
    media_dir: "media",
    transcode_dir: "transcoded",
    media_count: 4,
    active_streams: data.active_streams || 0,
    max_streams: data.max_streams || 10,
    nvenc_enabled: !!data.nvenc_enabled,
  };
}

export async function fetchEdgeMedia(): Promise<EdgeMediaItem[]> {
  const edgeUrl = getEdgeBaseUrl();
  const res = await fetch(`${edgeUrl}/api/media`, { cache: "no-store" });
  if (!res.ok) throw new Error(`Edge node media failed: ${res.statusText}`);
  return res.json();
}

export async function toggleMediaPublic(mediaId: string, isPublic: boolean): Promise<void> {
  const edgeUrl = getEdgeBaseUrl();
  const res = await fetch(`${edgeUrl}/api/media/${mediaId}/public`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ is_public: isPublic }),
  });
  if (!res.ok) throw new Error(`Failed updating public status: ${res.statusText}`);
}

export async function deleteEdgeMedia(mediaId: string, deleteFile = true): Promise<void> {
  const edgeUrl = getEdgeBaseUrl();
  const res = await fetch(`${edgeUrl}/api/media/${mediaId}?delete_file=${deleteFile}`, {
    method: "DELETE",
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || `Failed deleting media: ${res.statusText}`);
  }
}

export async function uploadEdgeMedia(
  file: File,
  onProgress?: (percent: number) => void
): Promise<{ status: string; filename: string; message: string }> {
  const edgeUrl = getEdgeBaseUrl();
  const formData = new FormData();
  formData.append("video", file);

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", `${edgeUrl}/api/media/upload`);

    if (xhr.upload && onProgress) {
      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          const percent = Math.round((event.loaded / event.total) * 100);
          onProgress(percent);
        }
      };
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const data = JSON.parse(xhr.responseText);
          resolve(data);
        } catch {
          resolve({ status: "ok", filename: file.name, message: "Uploaded successfully" });
        }
      } else {
        try {
          const errData = JSON.parse(xhr.responseText);
          reject(new Error(errData.error || `Upload failed: ${xhr.statusText}`));
        } catch {
          reject(new Error(`Upload failed with status ${xhr.status}`));
        }
      }
    };

    xhr.onerror = () => reject(new Error("Network error during file upload to edge node"));
    xhr.send(formData);
  });
}

export async function triggerEdgeScan(): Promise<{ total: number; added: number }> {
  const edgeUrl = getEdgeBaseUrl();
  const res = await fetch(`${edgeUrl}/api/media/scan`, {
    method: "POST",
  });
  if (!res.ok) throw new Error(`Library scan failed: ${res.statusText}`);
  return res.json();
}

export async function generateSampleClip(): Promise<{ filename: string; message: string }> {
  const edgeUrl = getEdgeBaseUrl();
  const res = await fetch(`${edgeUrl}/api/media/generate-sample`, {
    method: "POST",
  });
  if (!res.ok) throw new Error(`Sample generation failed: ${res.statusText}`);
  return res.json();
}

export async function pingEndpoint(url: string): Promise<{ ok: boolean; latencyMs: number; viaCloud?: boolean }> {
  const start = performance.now();
  try {
    const res = await fetch(`${url}/healthz`, { method: "GET", cache: "no-store" });
    const latencyMs = Math.round(performance.now() - start);
    return { ok: res.ok, latencyMs };
  } catch {
    // If probing edge daemon and local connection was blocked by browser PNA or remote access:
    if (url.includes("8090") || url.includes("edge")) {
      try {
        const cloudUrl = getApiBaseUrl();
        const res = await fetch(`${cloudUrl}/v1/nodes/demo/status`, { cache: "no-store" });
        if (res.ok) {
          const data = await res.json();
          if (data.is_online) {
            const latencyMs = Math.round(performance.now() - start);
            return { ok: true, latencyMs, viaCloud: true };
          }
        }
      } catch {
        // Fall through to offline
      }
    }
    const latencyMs = Math.round(performance.now() - start);
    return { ok: false, latencyMs };
  }
}

export function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${s.toString().padStart(2, "0")}`;
}

export interface CloudNodeStatus {
  node_id: string;
  node_name?: string;
  status: string;
  is_online: boolean;
  active_streams: number;
  max_streams: number;
  nvenc_enabled: boolean;
  version: string;
  last_heartbeat: string;
  owner_user_id?: string;
  owner_username?: string;
}

export async function fetchNodeStatus(
  nodeId: string = "898ea0c7-f7a3-48e8-ba59-b378a219f9ee"
): Promise<CloudNodeStatus> {
  const baseUrl = getApiBaseUrl();
  const res = await fetch(`${baseUrl}/v1/nodes/${nodeId}/status`, {
    cache: "no-store",
  });
  if (!res.ok) throw new Error("Failed to fetch node status");
  return res.json();
}

export async function claimNode(
  nodeId: string = "898ea0c7-f7a3-48e8-ba59-b378a219f9ee"
): Promise<{ status: string; message: string; owner_username: string }> {
  const baseUrl = getApiBaseUrl();
  const token = getAuthToken();
  if (!token) {
    throw new Error("Session expired or token missing. Please sign in again.");
  }
  const res = await fetch(`${baseUrl}/v1/nodes/${nodeId}/claim`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || "Failed to link edge node to your account");
  }
  return res.json();
}

export async function updateNodeName(
  nodeId: string = "898ea0c7-f7a3-48e8-ba59-b378a219f9ee",
  name: string
): Promise<{ status: string; node_name: string; message: string }> {
  const baseUrl = getApiBaseUrl();
  const token = getAuthToken();
  if (!token) {
    throw new Error("Session expired or token missing. Please sign in again.");
  }
  const res = await fetch(`${baseUrl}/v1/nodes/${nodeId}/name`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ name }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || "Failed to update node name");
  }
  return res.json();
}


