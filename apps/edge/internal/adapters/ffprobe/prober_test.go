package ffprobe

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSubprocessProber(t *testing.T) {
	pDefault := NewSubprocessProber("")
	if pDefault.binaryPath != "ffprobe" {
		t.Errorf("expected default 'ffprobe', got %q", pDefault.binaryPath)
	}

	pCustom := NewSubprocessProber("custom-ffprobe")
	if pCustom.binaryPath != "custom-ffprobe" {
		t.Errorf("expected custom-ffprobe, got %q", pCustom.binaryPath)
	}
}

func TestSubprocessProber_RealVideo(t *testing.T) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func TestSubprocessProber_AudioOnlyFile_ReturnsNoVideoStream(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping audio-only probe test")
	}

	tempDir := t.TempDir()
	sampleAudio := filepath.Join(tempDir, "audio_only.mp4")

	// Generate audio only (no video stream)
	cmd := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "sine=frequency=440:duration=1",
		"-c:a", "aac",
		sampleAudio,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to generate audio fixture: %v (%s)", err, string(out))
	}

	prober := NewSubprocessProber("")
	_, err := prober.Probe(context.Background(), sampleAudio)
	if !errors.Is(err, ErrNoVideoStream) {
		t.Errorf("expected ErrNoVideoStream on audio-only file, got: %v", err)
	}
}

func TestSubprocessProber_CorruptedFile(t *testing.T) {
	tempDir := t.TempDir()
	corrupted := filepath.Join(tempDir, "corrupted.mp4")
	if err := os.WriteFile(corrupted, []byte("this is plain text not video"), 0644); err != nil {
		t.Fatal(err)
	}

	prober := NewSubprocessProber("")
	_, err := prober.Probe(context.Background(), corrupted)
	if !errors.Is(err, ErrProbeFailed) {
		t.Errorf("expected ErrProbeFailed on corrupted media file, got: %v", err)
	}
}

func TestSubprocessProber_NonExistentFile(t *testing.T) {
	prober := NewSubprocessProber("")
	_, err := prober.Probe(context.Background(), "non_existent_video_path.mp4")
	if err == nil {
		t.Fatalf("expected error probing non-existent file, got nil")
	}
	if !errors.Is(err, ErrProbeFailed) {
		t.Errorf("expected ErrProbeFailed, got: %v", err)
	}
}

func TestSubprocessProber_ContextTimeout(t *testing.T) {
	prober := NewSubprocessProber("")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel context

	_, err := prober.Probe(ctx, "any_path.mp4")
	if err == nil {
		t.Fatal("expected error on cancelled context, got nil")
	}
}
