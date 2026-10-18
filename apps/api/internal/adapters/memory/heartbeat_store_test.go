package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/libs/contracts"
)

func TestInMemoryHeartbeatStore_Lifecycle(t *testing.T) {
	ctx := context.Background()
	store := memory.NewInMemoryHeartbeatStore()

	nodeID := "node-1111-2222"
	hb := &contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now().UTC(),
		Version:       "0.1.0",
		ActiveStreams: 3,
		GPUCaps: &contracts.GPUCapabilities{
			NVENC: true,
		},
	}

	// 1. Initial state: offline
	online, err := store.IsNodeOnline(ctx, nodeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if online {
		t.Fatalf("expected node to be initially offline")
	}

	// 2. Record heartbeat with 200ms TTL
	if err := store.RecordHeartbeat(ctx, hb, 200*time.Millisecond); err != nil {
		t.Fatalf("failed recording heartbeat: %v", err)
	}

	// 3. Node should be online
	online, err = store.IsNodeOnline(ctx, nodeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !online {
		t.Fatalf("expected node to be online")
	}

	retrieved, found, err := store.GetHeartbeat(ctx, nodeID)
	if err != nil || !found {
		t.Fatalf("failed retrieving active heartbeat: found=%t err=%v", found, err)
	}
	if retrieved.ActiveStreams != 3 {
		t.Errorf("expected 3 active streams, got %d", retrieved.ActiveStreams)
	}
	if retrieved.GPUCaps == nil || !retrieved.GPUCaps.NVENC {
		t.Errorf("expected NVENC GPU cap to be true")
	}

	// 4. Wait past TTL
	time.Sleep(250 * time.Millisecond)

	online, err = store.IsNodeOnline(ctx, nodeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if online {
		t.Fatalf("expected node to be offline after TTL expired")
	}

	retrieved, found, err = store.GetHeartbeat(ctx, nodeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found || retrieved != nil {
		t.Fatalf("expected heartbeat to be purged/unfound after TTL expiry")
	}
}

func TestInMemoryHeartbeatStore_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	store := memory.NewInMemoryHeartbeatStore()

	// 1. Nil heartbeat
	if err := store.RecordHeartbeat(ctx, nil, time.Minute); err == nil {
		t.Errorf("expected error on nil heartbeat, got nil")
	}

	// 2. Stream count exceeds 10-stream cap (Invariant 10)
	invalidHB := &contracts.Heartbeat{
		NodeID:        "node-cap",
		ActiveStreams: 11, // Exceeds 10
	}
	if err := store.RecordHeartbeat(ctx, invalidHB, time.Minute); err == nil {
		t.Errorf("expected error on active streams > 10, got nil")
	}

	// 3. Invalid TTL
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
}

func TestInMemoryHeartbeatStore_Batch(t *testing.T) {
	ctx := context.Background()
	store := memory.NewInMemoryHeartbeatStore()

	nodeA := "node-a"
	nodeB := "node-b"
	nodeC := "node-c"

	_ = store.RecordHeartbeat(ctx, &contracts.Heartbeat{NodeID: nodeA, ActiveStreams: 1}, 500*time.Millisecond)
	_ = store.RecordHeartbeat(ctx, &contracts.Heartbeat{NodeID: nodeB, ActiveStreams: 2}, 50*time.Millisecond)
	// nodeC has no heartbeat

	time.Sleep(80 * time.Millisecond) // nodeB expires, nodeA active

	statuses, err := store.BatchIsNodeOnline(ctx, []string{nodeA, nodeB, nodeC})
	if err != nil {
		t.Fatalf("batch query failed: %v", err)
	}

	if !statuses[nodeA] {
		t.Errorf("expected nodeA to be online")
	}
	if statuses[nodeB] {
		t.Errorf("expected nodeB to be offline after TTL")
	}
	if statuses[nodeC] {
		t.Errorf("expected nodeC to be offline")
	}
}

func TestInMemoryHeartbeatStore_SimulatedError(t *testing.T) {
	ctx := context.Background()
	store := memory.NewInMemoryHeartbeatStore()
	store.SimulateErr = errors.New("redis connection lost")

	hb := &contracts.Heartbeat{NodeID: "node-1", ActiveStreams: 1}
	if err := store.RecordHeartbeat(ctx, hb, time.Minute); err == nil {
		t.Errorf("expected error, got nil")
	}
	if _, err := store.IsNodeOnline(ctx, "node-1"); err == nil {
		t.Errorf("expected error, got nil")
	}
	if _, _, err := store.GetHeartbeat(ctx, "node-1"); err == nil {
		t.Errorf("expected error, got nil")
	}
	if _, err := store.BatchIsNodeOnline(ctx, []string{"node-1"}); err == nil {
		t.Errorf("expected error, got nil")
	}
}
