package ffprobe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os/exec"
	"strconv"

	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

var (
	ErrNoVideoStream = errors.New("no valid video stream found in media file")
	ErrProbeFailed   = errors.New("ffprobe command failed")
)

type ffprobeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
	Streams []struct {
		CodecType string `json:"codec_type"`
		CodecName string `json:"codec_name"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Duration  string `json:"duration"`
	} `json:"streams"`
}

// SubprocessProber implements ports.Prober by executing ffprobe as a managed subprocess.
type SubprocessProber struct {
	binaryPath string
}

// NewSubprocessProber constructs a prober pointing to the system ffprobe binary.
func NewSubprocessProber(binaryPath string) *SubprocessProber {
	if binaryPath == "" {
		binaryPath = "ffprobe"
	}
	return &SubprocessProber{
		binaryPath: binaryPath,
	}
}

// Probe executes ffprobe on the target file and extracts technical video metadata.
func (p *SubprocessProber) Probe(ctx context.Context, filePath string) (*ports.ProbeResult, error) {
	args := []string{
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	}

	cmd := exec.CommandContext(ctx, p.binaryPath, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s (stderr: %s)", ErrProbeFailed, err, stderr.String())
	}

	var output ffprobeOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe json output: %w", err)
	}

	var duration float64
	if output.Format.Duration != "" {
		duration, _ = strconv.ParseFloat(output.Format.Duration, 64)
	}

	var (
		videoFound bool
		width      int
		height     int
		videoCodec string
		audioCodec string
	)

	for _, stream := range output.Streams {
		if stream.CodecType == "video" && !videoFound {
			videoFound = true
			width = stream.Width
			height = stream.Height
			videoCodec = stream.CodecName
			if duration <= 0 && stream.Duration != "" {
				duration, _ = strconv.ParseFloat(stream.Duration, 64)
			}
		} else if stream.CodecType == "audio" && audioCodec == "" {
			audioCodec = stream.CodecName
		}
	}

	if !videoFound {
		return nil, ErrNoVideoStream
	}

	durationS := int(math.Round(duration))
	if durationS <= 0 {
		durationS = 1
	}

	return &ports.ProbeResult{
		DurationS:  durationS,
		Width:      width,
		Height:     height,
		VideoCodec: videoCodec,
		AudioCodec: audioCodec,
	}, nil
}
