import { fetchPlayback, PlaybackResponse } from "../../../lib/api-client";
import WatchClient from "../../../components/WatchClient";

interface WatchPageProps {
  params: Promise<{ id: string }>;
}

export default async function WatchPage({ params }: WatchPageProps) {
  const { id } = await params;

  let playback: PlaybackResponse | null = null;
  let playbackError: string | null = null;

  try {
    playback = await fetchPlayback(id);
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
