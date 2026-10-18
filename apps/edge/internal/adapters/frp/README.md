# FRP Client Adapter (`apps/edge/internal/adapters/frp`)

## Purpose
Manages the lifecycle of an outbound reverse tunnel via the official Fast Reverse Proxy (`frpc`) managed subprocess. Bypasses ISP CGNAT without requiring inbound open ports, adhering to **Invariant 2** (outbound-only connections from home edge nodes).

## Public Interface
- `Config`: Structure defining relay address, ports, auth token, node ID, and proxy modes (`http` with vhost/subdomain or `tcp`).
  - `Validate() error`: Verifies all fields meet boundary constraints and UUID formats.
  - `GenerateTOML() (string, error)`: Produces official FRP v0.52+ client TOML configuration.
  - `WriteConfigFile() (string, error)`: Writes configuration to a secure temporary file with `0600` permissions.
- `ManagedFRPClient`: Implements `ports.TunnelClient`.
  - `NewManagedFRPClient(cfg Config, logger *slog.Logger) *ManagedFRPClient`
  - `Start(ctx context.Context) error`: Validates configuration, launches `frpc` with supervisory auto-reconnect and exponential backoff.
  - `Stop() error`: Gracefully stops the child process, kills lingering instances, and removes temporary config files.
  - `IsConnected() bool`: Thread-safe query of active tunnel process execution.

## Testing
Run unit tests:
```bash
go test -v ./internal/adapters/frp
```
All tests use a self-contained mock helper process (`TestHelperProcess`) to verify lifecycle, double-start rejection, connection state polling, and graceful cleanup without external binary dependencies.
