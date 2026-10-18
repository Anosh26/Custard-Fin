package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/redis"
	"github.com/anosh26/custardfin/libs/contracts"
)

func TestHeartbeatStore_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	// Construct with nil client to test early validation without calling network
	store := redis.NewHeartbeatStore(nil)

	// 1. Nil heartbeat
	if err := store.RecordHeartbeat(ctx, nil, time.Minute); err == nil {
		t.Errorf("expected error on nil heartbeat, got nil")
	}

	// 2. Invalid heartbeat (missing NodeID)
	invalidHB := &contracts.Heartbeat{
		ActiveStreams: 1,
	}
	if err := store.RecordHeartbeat(ctx, invalidHB, time.Minute); err == nil {
		t.Errorf("expected error on missing node id, got nil")
	}

	// 3. Invalid heartbeat (exceeded 10 concurrent streams cap)
	capExceededHB := &contracts.Heartbeat{
		NodeID:        "node-exceeded",
		ActiveStreams: 12,
	}
	if err := store.RecordHeartbeat(ctx, capExceededHB, time.Minute); err == nil {
		t.Errorf("expected error on active streams > 10, got nil")
	}

	// 4. Invalid TTL
	validHB := &contracts.Heartbeat{
		NodeID:        "node-valid",
		ActiveStreams: 2,
	}
	if err := store.RecordHeartbeat(ctx, validHB, 0); err == nil {
		t.Errorf("expected error on zero TTL, got nil")
	}
	if err := store.RecordHeartbeat(ctx, validHB, -time.Second); err == nil {
		t.Errorf("expected error on negative TTL, got nil")
	}

	// 5. Empty nodeID queries return false/nil without calling redis
	online, err := store.IsNodeOnline(ctx, "")
	if err != nil || online {
		t.Errorf("expected false, nil on empty nodeID, got %t, %v", online, err)
	}

	hb, found, err := store.GetHeartbeat(ctx, "")
	if err != nil || found || hb != nil {
		t.Errorf("expected nil, false, nil on empty nodeID, got %v, %t, %v", hb, found, err)
	}

	batch, err := store.BatchIsNodeOnline(ctx, nil)
	if err != nil || len(batch) != 0 {
		t.Errorf("expected empty map on nil slice, got %v, %v", batch, err)
	}
}
