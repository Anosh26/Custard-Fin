package ffmpeg

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

func TestNewSubprocessTranscoder(t *testing.T) {
	tc1 := NewSubprocessTranscoder("", nil)
	if tc1.ffmpegPath != "ffmpeg" {
		t.Errorf("expected default 'ffmpeg', got %q", tc1.ffmpegPath)
	}
	if tc1.logger == nil {
		t.Error("expected non-nil default logger")
	}

	tc2 := NewSubprocessTranscoder("custom/ffmpeg", nil)
	if tc2.ffmpegPath != "custom/ffmpeg" {
		t.Errorf("expected 'custom/ffmpeg', got %q", tc2.ffmpegPath)
	}
}

func TestSelectRenditions(t *testing.T) {
	tests := []struct {
		name          string
		w             int
		h             int
		expectedCount int
		expectedFirst string
	}{
		{"4K UltraHD", 3840, 2160, 3, "1080p"},
		{"1080p FullHD", 1920, 1080, 3, "1080p"},
		{"720p HD", 1280, 720, 2, "720p"},
		{"480p SD", 854, 480, 1, "480p"},
		{"360p Small (no upscale)", 640, 360, 1, "480p"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renditions := SelectRenditions(tt.w, tt.h)
			if len(renditions) != tt.expectedCount {
				t.Fatalf("expected %d renditions, got %d", tt.expectedCount, len(renditions))
			}
			if renditions[0].Resolution != tt.expectedFirst {
				t.Errorf("expected first rendition %s, got %s", tt.expectedFirst, renditions[0].Resolution)
			}
		})
	}
}

func TestWriteMasterPlaylist(t *testing.T) {
	tempDir := t.TempDir()
	masterPath := filepath.Join(tempDir, "master.m3u8")

	profiles := []RenditionProfile{
		{Resolution: "1080p", Width: 1920, Height: 1080, TotalBitrateBps: 4660000},
		{Resolution: "720p", Width: 1280, Height: 720, TotalBitrateBps: 2528000},
	}
	renditions := []contracts.Rendition{
		{Resolution: "1080p", BitrateBps: 4660000, PlaylistPath: "1080p/index.m3u8", ChunkDurationS: 4},
		{Resolution: "720p", BitrateBps: 2528000, PlaylistPath: "720p/index.m3u8", ChunkDurationS: 4},
	}

	if err := writeMasterPlaylist(masterPath, profiles, renditions); err != nil {
		t.Fatalf("failed writing master playlist: %v", err)
	}

	content, err := os.ReadFile(masterPath)
	if err != nil {
		t.Fatalf("failed reading written master playlist: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "#EXTM3U") {
		t.Error("missing #EXTM3U header")
	}
	if !strings.Contains(str, "BANDWIDTH=4660000,RESOLUTION=1920x1080") {
		t.Error("missing 1080p stream inf")
	}
	if !strings.Contains(str, "1080p/index.m3u8") {
		t.Error("missing 1080p path")
	}
}

func TestSubprocessTranscoder_ValidationErrors(t *testing.T) {
	tc := NewSubprocessTranscoder("", nil)
	ctx := context.Background()

	// Empty input path
	_, err := tc.Transcode(ctx, ports.TranscodeJob{
		InputPath: "",
		OutputDir: "some/out",
	})
	if !errors.Is(err, ErrEmptyInputPath) {
		t.Errorf("expected ErrEmptyInputPath, got: %v", err)
	}

	// Empty output dir
	_, err = tc.Transcode(ctx, ports.TranscodeJob{
		InputPath: "video.mp4",
		OutputDir: "",
	})
	if !errors.Is(err, ErrEmptyOutputDir) {
		t.Errorf("expected ErrEmptyOutputDir, got: %v", err)
	}
}

func TestSubprocessTranscoder_ProbeEncoder(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed, skipping probe test")
	}

	tc := NewSubprocessTranscoder("", nil)
	encoder, err := tc.ProbeEncoder(context.Background())
	if err != nil {
		t.Fatalf("ProbeEncoder returned error: %v", err)
	}

	validEncoders := map[string]bool{
		"h264_nvenc": true,
		"h264_qsv":   true,
		"h264_vaapi": true,
		"libx264":    true,
	}
	if !validEncoders[encoder] {
		t.Errorf("unexpected probed encoder: %s", encoder)
	}
}

func TestSubprocessTranscoder_ProbeCapabilities(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed, skipping probe test")
	}

	tc := NewSubprocessTranscoder("", nil)
	caps := tc.ProbeCapabilities(context.Background())
	// Assert caps fields are readable
	t.Logf("Probed GPU capabilities: NVENC=%v, QSV=%v, VAAPI=%v", caps.NVENC, caps.QSV, caps.VAAPI)
}

func TestSubprocessTranscoder_RealTranscode_Integration(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping real transcode integration test")
	}

	tempDir := t.TempDir()
	sourceVideo := filepath.Join(tempDir, "source.mp4")
	outputDir := filepath.Join(tempDir, "hls_output")

	// Generate a 5-second test video (guarantees at least two 4-second chunk intervals)
	cmd := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "testsrc=duration=5:size=640x360:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=1000:duration=5",
		"-c:v", "libx264", "-c:a", "aac",
		sourceVideo,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed generating test clip: %v (%s)", err, string(out))
	}

	tc := NewSubprocessTranscoder("", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	job := ports.TranscodeJob{
		MediaID:       "550e8400-e29b-41d4-a716-446655440000",
		InputPath:     sourceVideo,
		OutputDir:     outputDir,
		SourceWidth:   640,
		SourceHeight:  360,
		PreferredMode: ports.EncoderModeAuto,
	}

	result, err := tc.Transcode(ctx, job)
	if err != nil {
		t.Fatalf("transcode failed: %v", err)
	}

	// 1. Verify master playlist
	if _, err := os.Stat(result.MasterPlaylistPath); err != nil {
		t.Errorf("master.m3u8 was not created: %v", err)
	}

	// 2. Verify renditions (480p baseline for 640x360 source)
	if len(result.Renditions) == 0 {
		t.Fatal("expected at least 1 rendition generated")
	}

	for _, r := range result.Renditions {
		// Invariant 6: must be strictly 4s chunks
		if r.ChunkDurationS != 4 {
			t.Errorf("expected 4s chunk duration invariant, got %d", r.ChunkDurationS)
		}

		playlistOnDisk := filepath.Join(outputDir, filepath.FromSlash(r.PlaylistPath))
		content, err := os.ReadFile(playlistOnDisk)
		if err != nil {
			t.Fatalf("failed reading rendition playlist %s: %v", playlistOnDisk, err)
		}

		playlistStr := string(content)
		if !strings.Contains(playlistStr, "#EXT-X-TARGETDURATION") {
			t.Errorf("playlist missing TARGETDURATION: %s", playlistStr)
		}
		if !strings.Contains(playlistStr, "segment_0000.ts") {
			t.Errorf("playlist missing first chunk segment_0000.ts: %s", playlistStr)
		}

		// Verify chunk file exists on disk
		firstChunk := filepath.Join(filepath.Dir(playlistOnDisk), "segment_0000.ts")
		if _, err := os.Stat(firstChunk); err != nil {
			t.Errorf("segment_0000.ts does not exist on disk: %v", err)
		}
	}
}

func TestSubprocessTranscoder_ExplicitCPU(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping CPU transcode test")
	}

	tempDir := t.TempDir()
	sourceVideo := filepath.Join(tempDir, "source_cpu.mp4")
	outputDir := filepath.Join(tempDir, "hls_cpu_output")

	cmd := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "testsrc=duration=2:size=320x240:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=1000:duration=2",
		"-c:v", "libx264", "-c:a", "aac",
		sourceVideo,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed generating test clip: %v (%s)", err, string(out))
	}

	tc := NewSubprocessTranscoder("", nil)
	result, err := tc.Transcode(context.Background(), ports.TranscodeJob{
		MediaID:       "550e8400-e29b-41d4-a716-446655440001",
		InputPath:     sourceVideo,
		OutputDir:     outputDir,
		SourceWidth:   320,
		SourceHeight:  240,
		PreferredMode: ports.EncoderModeCPU,
	})
	if err != nil {
		t.Fatalf("cpu transcode failed: %v", err)
	}

	if result.UsedEncoder != "libx264" {
		t.Errorf("expected used encoder 'libx264', got %q", result.UsedEncoder)
	}
}
