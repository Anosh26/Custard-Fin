export interface CatalogItem {
  media_id: string;
  title: string;
  description: string;
  duration_s: number;
  node_id: string;
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
}

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function fetchCatalog(
  q?: string,
  limit = 24,
  offset = 0
): Promise<CatalogResponse> {
  const params = new URLSearchParams();
  if (q && q.trim()) params.set("q", q.trim());
  params.set("limit", limit.toString());
  params.set("offset", offset.toString());

  const res = await fetch(`${API_BASE_URL}/v1/catalog?${params.toString()}`, {
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch catalog: ${res.statusText}`);
  }

  return res.json();
}

export async function fetchPlayback(
  mediaId: string
): Promise<PlaybackResponse> {
  const res = await fetch(`${API_BASE_URL}/v1/media/${mediaId}/playback`, {
    cache: "no-store",
  });

  if (!res.ok) {
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

export function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${s.toString().padStart(2, "0")}`;
}
