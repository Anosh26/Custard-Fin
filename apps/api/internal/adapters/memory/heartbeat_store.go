package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

var _ ports.HeartbeatStore = (*InMemoryHeartbeatStore)(nil)

type heartbeatEntry struct {
	hb        contracts.Heartbeat
	expiresAt time.Time
}

// InMemoryHeartbeatStore provides an in-memory implementation of ports.HeartbeatStore for tests.
type InMemoryHeartbeatStore struct {
	mu          sync.RWMutex
	entries     map[string]heartbeatEntry
	SimulateErr error
}

// NewInMemoryHeartbeatStore constructs a new in-memory heartbeat store.
func NewInMemoryHeartbeatStore() *InMemoryHeartbeatStore {
	return &InMemoryHeartbeatStore{
		entries: make(map[string]heartbeatEntry),
	}
}

// RecordHeartbeat records a heartbeat with expiration timestamp.
func (s *InMemoryHeartbeatStore) RecordHeartbeat(ctx context.Context, hb *contracts.Heartbeat, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.SimulateErr != nil {
		return s.SimulateErr
	}
	if hb == nil {
		return errors.New("heartbeat cannot be nil")
	}
	if err := hb.Validate(); err != nil {
		return fmt.Errorf("heartbeat validation failed: %w", err)
	}
	if ttl <= 0 {
		return errors.New("heartbeat ttl must be positive")
	}

	s.entries[hb.NodeID] = heartbeatEntry{
		hb:        *hb,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

// IsNodeOnline checks if the node's heartbeat exists and is unexpired.
func (s *InMemoryHeartbeatStore) IsNodeOnline(ctx context.Context, nodeID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.SimulateErr != nil {
		return false, s.SimulateErr
	}
	if nodeID == "" {
		return false, nil
	}

	entry, exists := s.entries[nodeID]
	if !exists {
		return false, nil
	}
	if time.Now().After(entry.expiresAt) {
		return false, nil
	}
	return true, nil
}

// GetHeartbeat retrieves the active heartbeat payload if present and unexpired.
func (s *InMemoryHeartbeatStore) GetHeartbeat(ctx context.Context, nodeID string) (*contracts.Heartbeat, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.SimulateErr != nil {
		return nil, false, s.SimulateErr
	}
	if nodeID == "" {
		return nil, false, nil
	}

	entry, exists := s.entries[nodeID]
	if !exists {
		return nil, false, nil
	}
	if time.Now().After(entry.expiresAt) {
		return nil, false, nil
	}
	copyHb := entry.hb
	return &copyHb, true, nil
}

// BatchIsNodeOnline checks multiple node online states.
func (s *InMemoryHeartbeatStore) BatchIsNodeOnline(ctx context.Context, nodeIDs []string) (map[string]bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.SimulateErr != nil {
		return nil, s.SimulateErr
	}

	res := make(map[string]bool, len(nodeIDs))
	now := time.Now()
	for _, id := range nodeIDs {
		entry, exists := s.entries[id]
		res[id] = exists && now.Before(entry.expiresAt)
	}
	return res, nil
}
