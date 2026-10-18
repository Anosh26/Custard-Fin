package ports

import (
	"context"

	"github.com/anosh26/custardfin/libs/contracts"
)

// EncoderMode defines the user preference or node capability for encoding.
type EncoderMode string

const (
	EncoderModeAuto EncoderMode = "auto"
	EncoderModeGPU  EncoderMode = "gpu"
	EncoderModeCPU  EncoderMode = "cpu"
)

// TranscodeJob defines the specifications for an HLS encoding task.
type TranscodeJob struct {
	MediaID         string
	InputPath       string
	OutputDir       string
	SourceWidth     int
	SourceHeight    int
	PreferredMode   EncoderMode
}

// TranscodeResult returns details about the completed HLS renditions.
type TranscodeResult struct {
	MasterPlaylistPath string
	Renditions         []contracts.Rendition
	UsedEncoder        string
}

// Transcoder defines the managed subprocess interface for pre-encoding media to HLS.
type Transcoder interface {
	ProbeEncoder(ctx context.Context) (string, error)
	Transcode(ctx context.Context, job TranscodeJob) (*TranscodeResult, error)
}
