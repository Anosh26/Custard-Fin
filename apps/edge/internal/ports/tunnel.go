package ports

import "context"

// TunnelClient manages an outbound reverse tunnel connection to the central relay server,
// allowing remote clients and origin shields to stream media over CGNAT without open inbound ports.
type TunnelClient interface {
	// Start launches the outbound tunnel client process and maintains connection with auto-reconnect.
	Start(ctx context.Context) error
	// Stop terminates the tunnel client process and cleans up any temporary configuration files.
	Stop() error
	// IsConnected returns true if the tunnel client process is currently active.
	IsConnected() bool
}
