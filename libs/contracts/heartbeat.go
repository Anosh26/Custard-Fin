package contracts

import (
	"errors"
	"fmt"
	"time"
)

// MaxConcurrentStreamsPerNode is a non-negotiable invariant cap.
const MaxConcurrentStreamsPerNode = 10

// GPUCapabilities defines the hardware encoder capabilities detected on the home node.
type GPUCapabilities struct {
	NVENC bool `json:"nvenc"`
	QSV   bool `json:"qsv"`
	VAAPI bool `json:"vaapi"`
}

// Heartbeat represents a node status ping sent every 15s to the cloud control tier.
type Heartbeat struct {
	NodeID        string           `json:"node_id"`
	SentAt        time.Time        `json:"sent_at"`
	Version       string           `json:"version"`
	GPUCaps       *GPUCapabilities `json:"gpu_caps,omitempty"`
	ActiveStreams int              `json:"active_streams"`
}

var (
	ErrInvalidNodeID        = errors.New("node_id is required")
	ErrExceededStreamCap    = errors.New("active streams exceeds node capacity")
	ErrNegativeStreamCount  = errors.New("active streams cannot be negative")
)

// Validate checks the heartbeat values for sanity and adherence to invariants.
func (h *Heartbeat) Validate() error {
	if h.NodeID == "" {
		return ErrInvalidNodeID
	}
	if h.ActiveStreams < 0 {
		return ErrNegativeStreamCount
	}
	if h.ActiveStreams > MaxConcurrentStreamsPerNode {
		return fmt.Errorf("%w: max %d allowed, reported %d", ErrExceededStreamCap, MaxConcurrentStreamsPerNode, h.ActiveStreams)
	}
	return nil
}
