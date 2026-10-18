package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
	redisClient "github.com/redis/go-redis/v9"
)

var _ ports.HeartbeatStore = (*HeartbeatStore)(nil)

// HeartbeatStore manages ephemeral node heartbeats in Redis.
// Invariant 5: Redis is used for heartbeats only.
type HeartbeatStore struct {
	client *redisClient.Client
}

// NewHeartbeatStore constructs a new Redis-backed heartbeat store.
func NewHeartbeatStore(client *redisClient.Client) *HeartbeatStore {
	return &HeartbeatStore{
		client: client,
	}
}

func heartbeatKey(nodeID string) string {
	return "heartbeat:" + nodeID
}

// RecordHeartbeat serializes and writes the heartbeat with TTL.
func (s *HeartbeatStore) RecordHeartbeat(ctx context.Context, hb *contracts.Heartbeat, ttl time.Duration) error {
	if hb == nil {
		return errors.New("heartbeat cannot be nil")
	}
	if err := hb.Validate(); err != nil {
		return fmt.Errorf("heartbeat validation failed: %w", err)
	}
	if ttl <= 0 {
		return errors.New("heartbeat ttl must be positive")
	}

	data, err := json.Marshal(hb)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat json: %w", err)
	}

	return s.client.Set(ctx, heartbeatKey(hb.NodeID), data, ttl).Err()
}

// IsNodeOnline checks whether the node's heartbeat TTL key is present in Redis.
func (s *HeartbeatStore) IsNodeOnline(ctx context.Context, nodeID string) (bool, error) {
	if nodeID == "" {
		return false, nil
	}

	val, err := s.client.Exists(ctx, heartbeatKey(nodeID)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check heartbeat existence in redis: %w", err)
	}

	return val > 0, nil
}

// GetHeartbeat retrieves the active heartbeat payload for a given node.
func (s *HeartbeatStore) GetHeartbeat(ctx context.Context, nodeID string) (*contracts.Heartbeat, bool, error) {
	if nodeID == "" {
		return nil, false, nil
	}

	val, err := s.client.Get(ctx, heartbeatKey(nodeID)).Result()
	if errors.Is(err, redisClient.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to retrieve heartbeat from redis: %w", err)
	}

	var hb contracts.Heartbeat
	if err := json.Unmarshal([]byte(val), &hb); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal heartbeat json: %w", err)
	}

	return &hb, true, nil
}

// BatchIsNodeOnline pipelines multiple key checks in a single round-trip.
func (s *HeartbeatStore) BatchIsNodeOnline(ctx context.Context, nodeIDs []string) (map[string]bool, error) {
	results := make(map[string]bool, len(nodeIDs))
	if len(nodeIDs) == 0 {
		return results, nil
	}

	pipe := s.client.Pipeline()
	cmds := make(map[string]*redisClient.IntCmd, len(nodeIDs))

	for _, id := range nodeIDs {
		cmds[id] = pipe.Exists(ctx, heartbeatKey(id))
	}

	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redisClient.Nil) {
		return nil, fmt.Errorf("failed executing redis pipeline for node liveness: %w", err)
	}

	for id, cmd := range cmds {
		results[id] = cmd.Val() > 0
	}

	return results, nil
}
