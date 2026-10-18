package frp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

var _ ports.TunnelClient = (*ManagedFRPClient)(nil)

// ManagedFRPClient manages the lifecycle of the frpc tunnel client subprocess.
type ManagedFRPClient struct {
	cfg        Config
	logger     *slog.Logger
	mu         sync.RWMutex
	cmd        *exec.Cmd
	configPath string
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	running    bool
	connected  bool

	// Overridable for testing
	lookPathFn func(file string) (string, error)
	cmdBuilder func(ctx context.Context, binPath, configPath string) *exec.Cmd
}

func defaultCmdBuilder(ctx context.Context, binPath, configPath string) *exec.Cmd {
	return exec.CommandContext(ctx, binPath, "-c", configPath)
}

// NewManagedFRPClient constructs a new managed FRP tunnel client.
func NewManagedFRPClient(cfg Config, logger *slog.Logger) *ManagedFRPClient {
	if logger == nil {
		logger = slog.Default()
	}
	return &ManagedFRPClient{
		cfg:        cfg,
		logger:     logger,
		lookPathFn: exec.LookPath,
		cmdBuilder: defaultCmdBuilder,
	}
}

// SetCmdBuilder overrides command execution constructor (used for tests).
func (c *ManagedFRPClient) SetCmdBuilder(fn func(ctx context.Context, binPath, configPath string) *exec.Cmd) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cmdBuilder = fn
}

// Start validates configuration, generates the temporary frpc.toml file,
// resolves the executable, and begins the supervisor process in the background.
func (c *ManagedFRPClient) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return errors.New("frp client is already running")
	}

	if err := c.cfg.Validate(); err != nil {
		return fmt.Errorf("invalid frp configuration: %w", err)
	}

	binPath, err := c.lookPathFn(c.cfg.BinaryPath)
	if err != nil {
		return fmt.Errorf("frpc binary not found at %q (please ensure frpc is installed and in PATH): %w", c.cfg.BinaryPath, err)
	}

	configPath, err := c.cfg.WriteConfigFile()
	if err != nil {
		return fmt.Errorf("failed to create frpc config file: %w", err)
	}
	c.configPath = configPath

	tunnelCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.running = true

	c.wg.Add(1)
	go c.supervise(tunnelCtx, binPath)

	c.logger.Info("managed frpc client supervisor started",
		"server", fmt.Sprintf("%s:%d", c.cfg.ServerAddr, c.cfg.ServerPort),
		"node_id", c.cfg.NodeID,
		"proxy_type", c.cfg.ProxyType,
		"config_path", configPath,
	)

	return nil
}

// supervise runs the subprocess and restarts it with exponential backoff if it exits unexpectedly.
func (c *ManagedFRPClient) supervise(ctx context.Context, binPath string) {
	defer c.wg.Done()

	backoff := 1 * time.Second
	const maxBackoff = 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("frpc supervisor exiting due to context cancellation")
			return
		default:
		}

		c.logger.Info("launching frpc subprocess", "bin", binPath, "config", c.configPath)

		cmd := c.cmdBuilder(ctx, binPath, c.configPath)

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			c.logger.Error("failed to create frpc stdout pipe", "err", err)
			c.sleepWithContext(ctx, backoff)
			backoff = minDuration(backoff*2, maxBackoff)
			continue
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			c.logger.Error("failed to create frpc stderr pipe", "err", err)
			c.sleepWithContext(ctx, backoff)
			backoff = minDuration(backoff*2, maxBackoff)
			continue
		}

		if err := cmd.Start(); err != nil {
			c.logger.Error("failed to start frpc process", "err", err)
			c.sleepWithContext(ctx, backoff)
			backoff = minDuration(backoff*2, maxBackoff)
			continue
		}

		c.mu.Lock()
		c.cmd = cmd
		c.connected = true
		c.mu.Unlock()

		// Stream logs asynchronously
		var pipeWg sync.WaitGroup
		pipeWg.Add(2)
		go c.streamLogs(stdoutPipe, "stdout", &pipeWg)
		go c.streamLogs(stderrPipe, "stderr", &pipeWg)

		// Wait for process exit
		waitErr := cmd.Wait()
		pipeWg.Wait()

		c.mu.Lock()
		c.connected = false
		c.cmd = nil
		c.mu.Unlock()

		if ctx.Err() != nil {
			// Expected termination
			c.logger.Info("frpc process terminated cleanly upon shutdown")
			return
		}

		c.logger.Warn("frpc process exited unexpectedly, scheduling restart",
			"err", waitErr,
			"retry_in", backoff.String(),
		)

		c.sleepWithContext(ctx, backoff)
		backoff = minDuration(backoff*2, maxBackoff)
	}
}

// streamLogs reads output lines from an frpc pipe and emits structured log events.
func (c *ManagedFRPClient) streamLogs(r io.Reader, streamName string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		c.logger.Debug("[frpc]", "stream", streamName, "msg", text)
	}
}

// sleepWithContext pauses execution until duration expires or context is cancelled.
func (c *ManagedFRPClient) sleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

// Stop gracefully stops the subprocess, waits for exit, and removes temporary config files.
func (c *ManagedFRPClient) Stop() error {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return nil
	}

	if c.cancel != nil {
		c.cancel()
	}

	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}

	configPath := c.configPath
	c.running = false
	c.connected = false
	c.mu.Unlock()

	c.wg.Wait()

	if configPath != "" {
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			c.logger.Warn("failed to clean up temporary frpc config", "path", configPath, "err", err)
		}
	}

	c.logger.Info("managed frpc client stopped and cleaned up")
	return nil
}

// IsConnected returns whether the frpc subprocess is currently executing.
func (c *ManagedFRPClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
