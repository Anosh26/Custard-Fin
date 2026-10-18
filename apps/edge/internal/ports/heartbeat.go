package ports

import (
	"context"

	"github.com/anosh26/custardfin/libs/contracts"
)

// HeartbeatReporter defines the interface for periodic node status reporting to Cloud Control.
type HeartbeatReporter interface {
	// SendHeartbeat sends an individual heartbeat payload to the cloud API.
	SendHeartbeat(ctx context.Context, hb *contracts.Heartbeat) error
	// Start starts the periodic background heartbeat loop.
	Start(ctx context.Context) error
	// Stop gracefully shuts down the background heartbeat worker.
	Stop() error
}

// StreamTracker reports the current number of active concurrent streams on the edge node.
type StreamTracker interface {
	ActiveStreams() int
}
