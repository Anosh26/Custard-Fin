package ffprobe

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestSubprocessProber_RealVideo(t *testing.T) {
	// Check if ffmpeg is available to create a quick test clip
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping real video probe test")
	}

	tempDir := t.TempDir()
	sampleVideo := filepath.Join(tempDir, "sample.mp4")

	// Generate a 2-second test video clip with audio
	cmd := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "testsrc=duration=2:size=640x360:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=1000:duration=2",
		"-c:v", "libx264",
		"-c:a", "aac",
		sampleVideo,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to generate test video fixture: %v (%s)", err, string(out))
	}

	prober := NewSubprocessProber("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := prober.Probe(ctx, sampleVideo)
	if err != nil {
		t.Fatalf("probe failed on generated sample: %v", err)
	}

	if res.DurationS != 2 {
		t.Errorf("expected duration 2s, got %d", res.DurationS)
	}
	if res.Width != 640 || res.Height != 360 {
		t.Errorf("expected 640x360, got %dx%d", res.Width, res.Height)
	}
	if res.VideoCodec != "h264" {
		t.Errorf("expected h264 video codec, got %s", res.VideoCodec)
	}
	if res.AudioCodec != "aac" {
		t.Errorf("expected aac audio codec, got %s", res.AudioCodec)
	}
}

func TestSubprocessProber_NonExistentFile(t *testing.T) {
	prober := NewSubprocessProber("")
	_, err := prober.Probe(context.Background(), "non_existent_video_path.mp4")
	if err == nil {
		t.Fatalf("expected error probing non-existent file, got nil")
	}
}
