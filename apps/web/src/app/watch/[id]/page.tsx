import { fetchPlayback, PlaybackResponse } from "../../../lib/api-client";
import WatchClient from "../../../components/WatchClient";

interface WatchPageProps {
  params: Promise<{ id: string }>;
  searchParams?: Promise<{ lan?: string }>;
}

export default async function WatchPage({ params, searchParams }: WatchPageProps) {
  const { id } = await params;
  const sp = searchParams ? await searchParams : {};
  const isLan = sp?.lan === "true";

  let playback: PlaybackResponse | null = null;
  let playbackError: string | null = null;

  try {
    playback = await fetchPlayback(id, isLan);
  } catch (err: unknown) {
    playbackError = err instanceof Error ? err.message : "Failed to resolve playback route";
  }

  return (
    <WatchClient
      mediaId={id}
      initialPlayback={playback}
      initialError={playbackError}
    />
  );
}
