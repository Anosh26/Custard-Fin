package heartbeat_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/heartbeat"
	"github.com/anosh26/custardfin/libs/contracts"
)

type mockStreamTracker struct {
	count int
}

func (m *mockStreamTracker) ActiveStreams() int {
	return m.count
}

func TestNewWorker_Validation(t *testing.T) {
	// 1. Missing NodeID
	_, err := heartbeat.NewWorker(heartbeat.Config{
		NodeToken: "token-123",
	})
	if err == nil {
		t.Fatalf("expected error for missing NodeID, got nil")
	}

	// 2. Missing NodeToken
	_, err = heartbeat.NewWorker(heartbeat.Config{
		NodeID: "node-123",
	})
	if err == nil {
		t.Fatalf("expected error for missing NodeToken, got nil")
	}

	// 3. Valid Config
	w, err := heartbeat.NewWorker(heartbeat.Config{
		NodeID:    "node-123",
		NodeToken: "token-123",
	})
	if err != nil {
		t.Fatalf("unexpected error for valid config: %v", err)
	}
	if w == nil {
		t.Fatalf("expected non-nil worker")
	}
}

func TestWorker_SendHeartbeat_Success(t *testing.T) {
	var receivedHeartbeat contracts.Heartbeat
	var receivedToken string
	var reqCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/nodes/node-test-123/heartbeat" {
			t.Errorf("expected path /v1/nodes/node-test-123/heartbeat, got %s", r.URL.Path)
		}

		receivedToken = r.Header.Get("X-Node-Token")

		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedHeartbeat)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	streamTracker := &mockStreamTracker{count: 4}
	gpuCaps := &contracts.GPUCapabilities{NVENC: true}

	w, err := heartbeat.NewWorker(heartbeat.Config{
		CloudBaseURL:  server.URL,
		NodeID:        "node-test-123",
		NodeToken:     "secret-token-xyz",
		Version:       "2.1.0",
		StreamTracker: streamTracker,
		GPUCaps:       gpuCaps,
	})
	if err != nil {
		t.Fatalf("failed to create worker: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := w.SendHeartbeat(ctx, nil); err != nil {
		t.Fatalf("failed to send heartbeat: %v", err)
	}

	if atomic.LoadInt32(&reqCount) != 1 {
		t.Errorf("expected 1 request, got %d", reqCount)
	}
	if receivedToken != "secret-token-xyz" {
		t.Errorf("expected token secret-token-xyz, got %s", receivedToken)
	}
	if receivedHeartbeat.NodeID != "node-test-123" {
		t.Errorf("expected NodeID node-test-123, got %s", receivedHeartbeat.NodeID)
	}
	if receivedHeartbeat.Version != "2.1.0" {
		t.Errorf("expected version 2.1.0, got %s", receivedHeartbeat.Version)
	}
	if receivedHeartbeat.ActiveStreams != 4 {
		t.Errorf("expected 4 active streams, got %d", receivedHeartbeat.ActiveStreams)
	}
	if receivedHeartbeat.GPUCaps == nil || !receivedHeartbeat.GPUCaps.NVENC {
		t.Errorf("expected NVENC GPU capability to be true")
	}
}

func TestWorker_SendHeartbeat_CapClamping(t *testing.T) {
	var receivedStreams int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var hb contracts.Heartbeat
		_ = json.NewDecoder(r.Body).Decode(&hb)
		receivedStreams = hb.ActiveStreams
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Case 1: Active streams exceeds Invariant 10 cap (e.g. 25 streams)
	streamTracker := &mockStreamTracker{count: 25}
	w, err := heartbeat.NewWorker(heartbeat.Config{
		CloudBaseURL:  server.URL,
		NodeID:        "node-cap-test",
		NodeToken:     "token-abc",
		StreamTracker: streamTracker,
	})
	if err != nil {
		t.Fatalf("failed creating worker: %v", err)
	}

	ctx := context.Background()
	if err := w.SendHeartbeat(ctx, nil); err != nil {
		t.Fatalf("expected clamped heartbeat to pass validation, got: %v", err)
	}

	if receivedStreams != contracts.MaxConcurrentStreamsPerNode {
		t.Errorf("expected clamped streams %d, got %d", contracts.MaxConcurrentStreamsPerNode, receivedStreams)
	}

	// Case 2: Negative stream count clamps to 0
	streamTracker.count = -5
	if err := w.SendHeartbeat(ctx, nil); err != nil {
		t.Fatalf("expected clamped negative streams to pass validation, got: %v", err)
	}
	if receivedStreams != 0 {
		t.Errorf("expected negative streams to clamp to 0, got %d", receivedStreams)
	}
}

func TestWorker_SendHeartbeat_Errors(t *testing.T) {
	// Subtest 1: Server returns HTTP 500
	server500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal database failure", http.StatusInternalServerError)
	}))
	defer server500.Close()

	w500, _ := heartbeat.NewWorker(heartbeat.Config{
		CloudBaseURL: server500.URL,
		NodeID:       "node-err-1",
		NodeToken:    "token-1",
	})
	err := w500.SendHeartbeat(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error for 500 response, got nil")
	}

	// Subtest 2: Server is completely unreachable
	wUnreachable, _ := heartbeat.NewWorker(heartbeat.Config{
		CloudBaseURL: "http://127.0.0.1:59999",
		NodeID:       "node-err-2",
		NodeToken:    "token-2",
	})
	err = wUnreachable.SendHeartbeat(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected network error for unreachable cloud, got nil")
	}

	// Subtest 3: Invalid explicit payload validation failure
	invalidHB := &contracts.Heartbeat{
		NodeID:        "", // empty NodeID violates contract
		ActiveStreams: 0,
	}
	err = w500.SendHeartbeat(context.Background(), invalidHB)
	if err == nil {
		t.Fatalf("expected validation error for empty NodeID, got nil")
	}
}

func TestWorker_Lifecycle_StartStopAndTicker(t *testing.T) {
	var mu sync.Mutex
	var beatsReceived int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		beatsReceived++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	w, err := heartbeat.NewWorker(heartbeat.Config{
		CloudBaseURL: server.URL,
		NodeID:       "node-ticker-1",
		NodeToken:    "token-ticker-1",
		Interval:     50 * time.Millisecond, // Fast interval for testing
	})
	if err != nil {
		t.Fatalf("failed creating worker: %v", err)
	}

	ctx := context.Background()

	// 1. Initial Start() triggers immediate heartbeat
	if err := w.Start(ctx); err != nil {
		t.Fatalf("failed to start worker: %v", err)
	}

	// 2. Start called again returns ErrWorkerRunning
	if err := w.Start(ctx); err == nil {
		t.Fatalf("expected ErrWorkerRunning on double start, got nil")
	}

	// Allow ticker to fire at least twice
	time.Sleep(140 * time.Millisecond)

	mu.Lock()
	count := beatsReceived
	mu.Unlock()

	if count < 2 {
		t.Errorf("expected at least 2 heartbeats received, got %d", count)
	}

	// 3. Stop terminates cleanly
	if err := w.Stop(); err != nil {
		t.Fatalf("failed to stop worker: %v", err)
	}

	// 4. Stop called again returns ErrWorkerNotRunning
	if err := w.Stop(); err == nil {
		t.Fatalf("expected ErrWorkerNotRunning on double stop, got nil")
	}
}

func TestWorker_Lifecycle_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	w, err := heartbeat.NewWorker(heartbeat.Config{
		CloudBaseURL: server.URL,
		NodeID:       "node-cancel-1",
		NodeToken:    "token-cancel-1",
		Interval:     50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("failed to create worker: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := w.Start(ctx); err != nil {
		t.Fatalf("failed to start worker: %v", err)
	}

	// Cancel context after brief moment
	time.Sleep(20 * time.Millisecond)
	cancel()

	// Wait briefly for goroutine to exit on ctx.Done()
	time.Sleep(50 * time.Millisecond)

	// Stop cleanly cleans up state
	_ = w.Stop()
}
