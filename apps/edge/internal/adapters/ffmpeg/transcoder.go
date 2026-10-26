package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anosh26/custardfin/apps/edge/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

var (
	ErrEmptyInputPath  = errors.New("input path cannot be empty")
	ErrEmptyOutputDir  = errors.New("output directory cannot be empty")
	ErrTranscodeFailed = errors.New("ffmpeg transcoding failed")
)

// RenditionProfile specifies resolution and bitrate constraints for an HLS rendition.
type RenditionProfile struct {
	Resolution      string
	Width           int
	Height          int
	VideoBitrate    string
	AudioBitrate    string
	TotalBitrateBps int64
}

// SubprocessTranscoder manages FFmpeg subprocess encoding with hardware acceleration and automatic CPU fallback.
type SubprocessTranscoder struct {
	ffmpegPath string
	logger     *slog.Logger
}

// NewSubprocessTranscoder constructs an instance of SubprocessTranscoder.
func NewSubprocessTranscoder(ffmpegPath string, logger *slog.Logger) *SubprocessTranscoder {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &SubprocessTranscoder{
		ffmpegPath: ffmpegPath,
		logger:     logger,
	}
}

// ProbeCapabilities inspects available hardware encoders via FFmpeg test encodes.
func (t *SubprocessTranscoder) ProbeCapabilities(ctx context.Context) contracts.GPUCapabilities {
	caps := contracts.GPUCapabilities{}
	if t.probeSingleEncoder(ctx, "h264_nvenc") {
		caps.NVENC = true
	}
	if t.probeSingleEncoder(ctx, "h264_qsv") {
		caps.QSV = true
	}
	if t.probeSingleEncoder(ctx, "h264_vaapi") {
		caps.VAAPI = true
	}
	return caps
}

func (t *SubprocessTranscoder) probeSingleEncoder(ctx context.Context, encoder string) bool {
	cmd := exec.CommandContext(ctx, t.ffmpegPath,
		"-v", "error",
		"-f", "lavfi", "-i", "testsrc=duration=1:size=256x256:rate=30",
		"-pix_fmt", "yuv420p",
		"-c:v", encoder,
		"-f", "null", "-",
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	return cmd.Run() == nil
}

// ProbeEncoder tests hardware encoders in order of preference (NVENC -> QSV -> VAAPI)
// and returns the best operational encoder, falling back to libx264.
func (t *SubprocessTranscoder) ProbeEncoder(ctx context.Context) (string, error) {
	caps := t.ProbeCapabilities(ctx)
	if caps.NVENC {
		t.logger.Info("hardware acceleration active: NVIDIA NVENC (h264_nvenc)")
		return "h264_nvenc", nil
	}
	if caps.QSV {
		t.logger.Info("hardware acceleration active: Intel QuickSync (h264_qsv)")
		return "h264_qsv", nil
	}
	if caps.VAAPI {
		t.logger.Info("hardware acceleration active: VAAPI (h264_vaapi)")
		return "h264_vaapi", nil
	}

	t.logger.Info("hardware acceleration unavailable, falling back to CPU (libx264)")
	return "libx264", nil
}

// SelectRenditions selects the appropriate quality ladder renditions without upscaling.
func SelectRenditions(srcWidth, srcHeight int) []RenditionProfile {
	var ladder []RenditionProfile

	// 1080p
	if srcHeight >= 1080 || srcWidth >= 1920 {
		ladder = append(ladder, RenditionProfile{
			Resolution:      "1080p",
			Width:           1920,
			Height:          1080,
			VideoBitrate:    "4500k",
			AudioBitrate:    "160k",
			TotalBitrateBps: 4660000,
		})
	}

	// 720p
	if srcHeight >= 720 || srcWidth >= 1280 {
		ladder = append(ladder, RenditionProfile{
			Resolution:      "720p",
			Width:           1280,
			Height:          720,
			VideoBitrate:    "2400k",
			AudioBitrate:    "128k",
			TotalBitrateBps: 2528000,
		})
	}

	// 480p (baseline)
	targetW := 854
	targetH := 480
	if srcWidth > 0 && srcHeight > 0 && (srcWidth < targetW || srcHeight < targetH) {
		targetW = srcWidth
		targetH = srcHeight
	}

	ladder = append(ladder, RenditionProfile{
		Resolution:      "480p",
		Width:           targetW,
		Height:          targetH,
		VideoBitrate:    "900k",
		AudioBitrate:    "96k",
		TotalBitrateBps: 996000,
	})

	return ladder
}

// Transcode pre-encodes the input media into 4-second HLS segments for each selected rendition.
func (t *SubprocessTranscoder) Transcode(ctx context.Context, job ports.TranscodeJob) (*ports.TranscodeResult, error) {
	if strings.TrimSpace(job.InputPath) == "" {
		return nil, ErrEmptyInputPath
	}
	if strings.TrimSpace(job.OutputDir) == "" {
		return nil, ErrEmptyOutputDir
	}

	if err := os.MkdirAll(job.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Determine initial encoder
	var encoder string
	switch job.PreferredMode {
	case ports.EncoderModeCPU:
		encoder = "libx264"
	case ports.EncoderModeGPU:
		encoder = "h264_nvenc"
	default:
		var err error
		encoder, err = t.ProbeEncoder(ctx)
		if err != nil {
			encoder = "libx264"
		}
	}

	profiles := SelectRenditions(job.SourceWidth, job.SourceHeight)
	var contractsRenditions []contracts.Rendition

	usedEncoder := encoder

	for _, profile := range profiles {
		renditionDir := filepath.Join(job.OutputDir, profile.Resolution)
		if err := os.MkdirAll(renditionDir, 0755); err != nil {
			return nil, fmt.Errorf("failed creating rendition directory: %w", err)
		}

		playlistPath := filepath.Join(renditionDir, "index.m3u8")

		// Attempt encode with primary encoder
		err := t.encodeRendition(ctx, job.InputPath, renditionDir, playlistPath, profile, usedEncoder)
		if err != nil && usedEncoder != "libx264" && job.PreferredMode != ports.EncoderModeGPU {
			// Automatic CPU fallback on any hardware failure (Invariant 6)
			t.logger.Warn("hardware encoding failed, triggering automatic CPU fallback", "encoder", usedEncoder, "err", err)
			usedEncoder = "libx264"
			err = t.encodeRendition(ctx, job.InputPath, renditionDir, playlistPath, profile, usedEncoder)
		}

		if err != nil {
			return nil, fmt.Errorf("rendition %s failed: %w", profile.Resolution, err)
		}

		contractsRenditions = append(contractsRenditions, contracts.Rendition{
			Resolution:     profile.Resolution,
			BitrateBps:     profile.TotalBitrateBps,
			PlaylistPath:   filepath.ToSlash(filepath.Join(profile.Resolution, "index.m3u8")),
			ChunkDurationS: 4, // Invariant 6: strictly 4s chunks
		})
	}

	// Generate master.m3u8
	masterPath := filepath.Join(job.OutputDir, "master.m3u8")
	if err := writeMasterPlaylist(masterPath, profiles, contractsRenditions); err != nil {
		return nil, fmt.Errorf("failed generating master playlist: %w", err)
	}

	// Generate crisp 640x360 poster thumbnail (thumb.jpg)
	thumbPath := filepath.Join(job.OutputDir, "thumb.jpg")
	if err := t.GenerateThumbnail(ctx, job.InputPath, thumbPath); err != nil {
		t.logger.Warn("thumbnail generation failed", "err", err, "input", job.InputPath)
	}

	return &ports.TranscodeResult{
		MasterPlaylistPath: masterPath,
		ThumbnailPath:      thumbPath,
		Renditions:         contractsRenditions,
		UsedEncoder:        usedEncoder,
	}, nil
}

// GenerateThumbnail extracts a 640x360 JPEG poster frame from the media input.
func (t *SubprocessTranscoder) GenerateThumbnail(ctx context.Context, inputPath, outputPath string) error {
	scaleFilter := "scale=640:360:force_original_aspect_ratio=decrease,pad=640:360:(ow-iw)/2:(oh-ih)/2"
	args := []string{
		"-y",
		"-ss", "00:00:01",
		"-i", inputPath,
		"-vframes", "1",
		"-q:v", "2",
		"-vf", scaleFilter,
		outputPath,
	}

	cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)
	if err := cmd.Run(); err != nil {
		// Fallback to offset 0s if media is shorter than 1 second
		args[2] = "00:00:00"
		cmdFallback := exec.CommandContext(ctx, t.ffmpegPath, args...)
		return cmdFallback.Run()
	}
	return nil
}

func (t *SubprocessTranscoder) encodeRendition(
	ctx context.Context,
	inputPath string,
	renditionDir string,
	playlistPath string,
	profile RenditionProfile,
	encoder string,
) error {
	scaleFilter := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2",
		profile.Width, profile.Height, profile.Width, profile.Height)

	segmentPattern := filepath.Join(renditionDir, "segment_%04d.ts")

	args := []string{
		"-y",
		"-i", inputPath,
		"-vf", scaleFilter,
		"-pix_fmt", "yuv420p",
		"-c:v", encoder,
		"-b:v", profile.VideoBitrate,
		"-maxrate", profile.VideoBitrate,
		"-bufsize", "2M",
		// Force keyframes strictly on 4-second multiples to align chunks
		"-force_key_frames", "expr:gte(t,n_forced*4)",
	}

	if encoder == "libx264" {
		args = append(args, "-preset", "veryfast")
	} else if encoder == "h264_nvenc" {
		args = append(args, "-preset", "p4", "-tune", "hq")
	}

	args = append(args,
		"-c:a", "aac",
		"-b:a", profile.AudioBitrate,
		"-ac", "2",
		"-hls_time", "4",
		"-hls_list_size", "0",
		"-hls_segment_type", "mpegts",
		"-hls_segment_filename", segmentPattern,
		playlistPath,
	)

	cmd := exec.CommandContext(ctx, t.ffmpegPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s (stderr: %s)", ErrTranscodeFailed, err, stderr.String())
	}

	return nil
}

func writeMasterPlaylist(masterPath string, profiles []RenditionProfile, renditions []contracts.Rendition) error {
	var buf strings.Builder
	buf.WriteString("#EXTM3U\n")
	buf.WriteString("#EXT-X-VERSION:3\n")

	for i, r := range renditions {
		p := profiles[i]
		buf.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n", r.BitrateBps, p.Width, p.Height))
		buf.WriteString(fmt.Sprintf("%s/index.m3u8\n", r.Resolution))
	}

	return os.WriteFile(masterPath, []byte(buf.String()), 0644)
}
