package ports

import (
	"context"
	"time"

	"github.com/anosh26/custardfin/libs/contracts"
)

// HeartbeatStore persists ephemeral node liveness heartbeats with TTL.
// Strictly adheres to Invariant 5 (Redis is used for heartbeats only).
type HeartbeatStore interface {
	// RecordHeartbeat writes the heartbeat data under heartbeat:{node_id} with the specified TTL.
	RecordHeartbeat(ctx context.Context, hb *contracts.Heartbeat, ttl time.Duration) error

	// IsNodeOnline checks whether the heartbeat key for nodeID exists and is unexpired.
	IsNodeOnline(ctx context.Context, nodeID string) (bool, error)

	// GetHeartbeat retrieves the active heartbeat payload if present and unexpired.
	GetHeartbeat(ctx context.Context, nodeID string) (*contracts.Heartbeat, bool, error)

	// BatchIsNodeOnline checks online status for a slice of node IDs.
	BatchIsNodeOnline(ctx context.Context, nodeIDs []string) (map[string]bool, error)
}
