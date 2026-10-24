package admission

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

var (
	ErrNodeOffline = errors.New("node is offline")
)

// Config configures the saturation estimator thresholds.
type Config struct {
	ModerateThreshold  int // default 4
	SaturatedThreshold int // default 7
	ExhaustedThreshold int // default 10 (Invariant 10)
}

// DefaultEstimator evaluates edge node saturation based on heartbeats recorded in Redis.
type DefaultEstimator struct {
	cfg            Config
	heartbeatStore ports.HeartbeatStore
	logger         *slog.Logger
}

var _ ports.SaturationEstimator = (*DefaultEstimator)(nil)

// NewDefaultEstimator constructs a new DefaultEstimator.
func NewDefaultEstimator(cfg Config, heartbeatStore ports.HeartbeatStore, logger *slog.Logger) *DefaultEstimator {
	if cfg.ModerateThreshold <= 0 {
		cfg.ModerateThreshold = 4
	}
	if cfg.SaturatedThreshold <= 0 {
		cfg.SaturatedThreshold = 7
	}
	if cfg.ExhaustedThreshold <= 0 {
		cfg.ExhaustedThreshold = contracts.MaxConcurrentStreamsPerNode
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &DefaultEstimator{
		cfg:            cfg,
		heartbeatStore: heartbeatStore,
		logger:         logger,
	}
}

// EstimateLoad inspects the node heartbeat in Redis and returns the corresponding LoadLevel and active stream count.
func (e *DefaultEstimator) EstimateLoad(ctx context.Context, nodeID string) (domain.LoadLevel, int, error) {
	if e.heartbeatStore == nil {
		// If heartbeat store is not configured, assume normal load
		return domain.LoadNormal, 0, nil
	}

	hb, online, err := e.heartbeatStore.GetHeartbeat(ctx, nodeID)
	if err != nil {
		e.logger.Error("failed retrieving heartbeat for load estimation", "node_id", nodeID, "err", err)
		return domain.LoadExhausted, 0, fmt.Errorf("heartbeat lookup failed: %w", err)
	}

	if !online || hb == nil {
		return domain.LoadExhausted, 0, ErrNodeOffline
	}

	active := hb.ActiveStreams
	if active < 0 {
		active = 0
	}

	switch {
	case active >= e.cfg.ExhaustedThreshold:
		return domain.LoadExhausted, active, nil
	case active >= e.cfg.SaturatedThreshold:
		return domain.LoadSaturated, active, nil
	case active >= e.cfg.ModerateThreshold:
		return domain.LoadModerate, active, nil
	default:
		return domain.LoadNormal, active, nil
	}
}
