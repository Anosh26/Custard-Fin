package frp_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/frp"
)

// TestHelperProcess is used by exec.Command to simulate the frpc binary.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	var configPath string
	for i := 0; i < len(args); i++ {
		if args[i] == "-c" && i+1 < len(args) {
			configPath = args[i+1]
			break
		}
	}

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "missing -c config argument: %v\n", os.Args)
		os.Exit(2)
	}

	// Verify config exists
	if _, err := os.Stat(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "config file not found: %s\n", configPath)
		os.Exit(1)
	}

	fmt.Println("2026/10/18 20:00:00 [I] [service.go:100] frpc service started")
	fmt.Println("2026/10/18 20:00:00 [I] [proxy.go:200] [stream-test] start proxy success")

	// Sleep until killed by signal/context
	time.Sleep(30 * time.Second)
	os.Exit(0)
}

func mockCmdBuilder(ctx context.Context, binPath, configPath string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHelperProcess", "--", "-c", configPath)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestManagedFRPClient_ValidationFailure(t *testing.T) {
	invalidCfg := frp.Config{} // Empty
	client := frp.NewManagedFRPClient(invalidCfg, nil)

	err := client.Start(context.Background())
	if err == nil {
		t.Fatalf("expected error on invalid config, got nil")
	}
}

func TestManagedFRPClient_BinaryNotFound(t *testing.T) {
	cfg := validHTTPConfig()
	cfg.BinaryPath = "non_existent_frpc_executable_12345"
	client := frp.NewManagedFRPClient(cfg, nil)

	err := client.Start(context.Background())
	if err == nil {
		t.Fatalf("expected error on non-existent binary, got nil")
	}
}

func TestManagedFRPClient_Lifecycle_WithMockProcess(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "frp_client_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := validHTTPConfig()
	cfg.ConfigDir = tempDir
	cfg.BinaryPath = os.Args[0]

	client := frp.NewManagedFRPClient(cfg, nil)
	client.SetCmdBuilder(mockCmdBuilder)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Start(ctx); err != nil {
		t.Fatalf("failed to start managed client: %v", err)
	}

	// 1. Double start returns error
	if err := client.Start(ctx); err == nil {
		t.Errorf("expected error on second Start() call, got nil")
	}

	// 2. Wait for process to start and verify connection state
	connected := false
	for start := time.Now(); time.Since(start) < 2*time.Second; {
		if client.IsConnected() {
			connected = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !connected {
		t.Errorf("expected client to report connected within 2s, got false")
	}

	// 3. Stop client
	if err := client.Stop(); err != nil {
		t.Fatalf("Stop() returned error: %v", err)
	}

	if client.IsConnected() {
		t.Errorf("expected client to report disconnected after Stop(), got true")
	}

	// Verify temporary config file was cleaned up
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read temp dir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected config file to be cleaned up, found %d files in temp dir", len(entries))
	}

	// 4. Double stop is idempotent
	if err := client.Stop(); err != nil {
		t.Errorf("expected idempotent Stop(), got error: %v", err)
	}
}

func TestManagedFRPClient_ContextCancellation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "frp_cancel_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := validHTTPConfig()
	cfg.ConfigDir = tempDir
	cfg.BinaryPath = os.Args[0]

	client := frp.NewManagedFRPClient(cfg, nil)
	client.SetCmdBuilder(mockCmdBuilder)

	ctx, cancel := context.WithCancel(context.Background())

	if err := client.Start(ctx); err != nil {
		t.Fatalf("failed to start managed client: %v", err)
	}

	connected := false
	for start := time.Now(); time.Since(start) < 2*time.Second; {
		if client.IsConnected() {
			connected = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !connected {
		t.Errorf("expected client to report connected, got false")
	}

	// Cancel the context
	cancel()

	// Wait for cleanup
	time.Sleep(300 * time.Millisecond)
	_ = client.Stop()

	if client.IsConnected() {
		t.Errorf("expected client to be disconnected after context cancellation")
	}
}
